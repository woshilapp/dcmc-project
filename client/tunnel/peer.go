package tunnel

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/woshilapp/dcmc-project/client/global"
	netdata "github.com/woshilapp/dcmc-project/network"
)

// Handles
func HandleRemotePeer(t *global.Tunnel, conn net.Conn) {
	for {
		if t.Closed {
			return
		}

		id, status, data, err := netdata.TunnelTCPRead(conn)
		if err != nil {
			fmt.Println("RemoteERR")
			return
		}

		peerconn := TGetConnP(t, id)
		if peerconn == nil {
			continue
		}

		switch status {
		case 1:
			_, err = peerconn.Write(data)
			if err != nil {
				fmt.Println("err1111", err)
				TDelConnP(t, id)

				netdata.TunnelTCPWrite(id, 0, t.TCPRemote, []byte{})
			}
		case 0:
			TDelConnP(t, id)
			peerconn.Close()
		}
	}
}

func HandleLocalPeer(t *global.Tunnel, conn net.Conn, id uint32) {
	for {
		if t.Closed {
			return
		}

		buf := make([]byte, 16*1024)

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err222", err)

			TDelConnP(t, id)
			netdata.TunnelTCPWrite(id, 0, t.TCPRemote, []byte{})

			return
		}

		err = netdata.TunnelTCPWrite(id, 1, t.TCPRemote, buf[:n])
		if err != nil {
			fmt.Println("errwrt", err)
			return
		}
	}
}

func ListenLocal(t *global.Tunnel) {
	listenerf, err := net.Listen("tcp", "127.0.0.2:"+strconv.Itoa(int(t.Port)))
	if err != nil {
		fmt.Println("[ERROR]", err)
	}

	for {
		if t.Closed {
			return
		}

		conn, err := listenerf.Accept()
		if err != nil {
			fmt.Println("[ERROR]", err)
		}

		id := TAddConnP(t, conn)

		go HandleLocalPeer(t, conn, id)
	}
}

func HandleUDPRemotePeer(t *global.Tunnel, localconn *net.UDPConn) {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			// not protocol, won't to data layer
			t.UDPRemote.WriteTo([]byte("keepalive"), t.UDPRemoteAddr)
		}
	}()

	for {
		if t.Closed {
			return
		}

		id, _, data, err := netdata.TunnelUDPRead(t.UDPRemote)
		if err != nil {
			if err.Error() == "not dcmc protocol" {
				continue
			}

			fmt.Println("RemoteERR", err)
			return
		}

		peer_addr := UGetAddrP(t, id)
		if peer_addr == nil {
			continue
		}

		_, err = localconn.WriteTo(data, peer_addr)
		if err != nil {
			fmt.Println("udpwerr")
			return
		}
	}
}

func UDPListenLocal(t *global.Tunnel) {
	buf := make([]byte, 64*1024)

	listaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.2:"+strconv.Itoa(int(t.Port)))
	listener, err := net.ListenUDP("udp", listaddr)
	if err != nil {
		fmt.Println("[ERROR]", err)
	}

	go HandleUDPRemotePeer(t, listener)

	for {
		if t.Closed {
			listener.Close()
			return
		}

		n, addr, err := listener.ReadFrom(buf)
		if err != nil {
			fmt.Println("[ERROR]", err)
			return
		}

		var id int

		if id = UGetIDP(t, addr); id == -1 {
			id = int(UAddAddrP(t, addr))
		}

		err = netdata.TunnelUDPWrite(uint32(id), t.UDPRemote, t.UDPRemoteAddr, buf[:n])
		if err != nil {
			fmt.Println("err1123")
			return
		}
	}
}
