package tunnel

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	reuse "github.com/libp2p/go-reuseport"
	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/client/network"
	netdata "github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
)

func TCPPunchHost(peer *global.TPeers, port uint16) {
	netdata.WriteMsg(global.Serverconn, []byte("301"))

	punch_id := <-global.Host.PunchIDs
	tun := &global.Tunnel{
		Proto:    1,
		Port:     port,
		PunchID:  punch_id,
		Lock:     sync.RWMutex{},
		Closed:   false,
		TCPConns: map[uint32]net.Conn{},
	}
	peer.Tunnels = append(peer.Tunnels, tun)

	str, _ := protocol.Encode(330, punch_id, 1, port)
	netdata.WriteMsg(peer.Conn, []byte(str))

	tmp_conn, err := reuse.Dial("tcp", "0.0.0.0:0", global.Serveraddr.String())
	if err != nil {
		fmt.Println("Punch connect server failed")
		return
	}

	go network.HandlePunchConn(tmp_conn)

	str, _ = protocol.Encode(302, punch_id)
	netdata.WriteMsg(tmp_conn, []byte(str))

	tun.TCPRemote = tmp_conn
	global.Host.PIDtun[punch_id] = tun
}

func UDPPunchHost(peer *global.TPeers, port uint16) {
	netdata.WriteMsg(global.Serverconn, []byte("301"))

	punch_id := <-global.Host.PunchIDs
	tun := &global.Tunnel{
		Proto:    1,
		Port:     port,
		PunchID:  punch_id,
		Lock:     sync.RWMutex{},
		Closed:   false,
		UDPConns: map[uint32]*net.UDPConn{},
	}
	peer.Tunnels = append(peer.Tunnels, tun)

	str, _ := protocol.Encode(330, punch_id, 2, port)
	netdata.WriteMsg(peer.Conn, []byte(str))

	sock, err := net.DialUDP("udp", nil, nil)
	if err != nil {
		fmt.Println("Punch Listen Err")
	}

	go network.HandleUDP(sock)

	str, _ = protocol.Encode(302, punch_id)
	sock.WriteTo([]byte(str), global.Serveraddr)

	tun.UDPRemote = sock
	global.Host.PIDtun[punch_id] = tun
}

func HandleNewPeer(peer *global.TPeers) {
	for _, port := range global.Host.TCPPorts {
		TCPPunchHost(peer, port)
	}

	for _, port := range global.Host.UDPPorts {
		UDPPunchHost(peer, port)
	}
}

// Handles
func HandleRemoteHost(t *global.Tunnel, conn net.Conn) {
	for {
		if t.Closed {
			return
		}

		id, status, data, err := netdata.TunnelTCPRead(conn)
		if err != nil {
			return
		}
		// fmt.Println("HRR:", string(data))

		peerconn := TGetConnH(t, id)
		if peerconn == nil {
			peerconn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(int(t.Port)))
			if err != nil {
				fmt.Println("err111", err)
				continue
			}

			TAddConnH(t, id, peerconn)

			go HandleLocalHost(t, peerconn, id)
		}

		peerconn = TGetConnH(t, id)

		switch status {
		case 1:
			_, err = peerconn.Write(data)
			if err != nil {
				fmt.Println("err222", err)
				TDelConnH(t, id)

				netdata.TunnelTCPWrite(id, 0, t.TCPRemote, []byte{})
			}
		case 0:
			peerconn.Close()
			TDelConnH(t, id)
		}
	}
}

func HandleLocalHost(t *global.Tunnel, conn net.Conn, id uint32) {
	for {
		if t.Closed {
			return
		}

		buf := make([]byte, 16*1024)

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err444", err)
			TDelConnH(t, id)

			netdata.TunnelTCPWrite(id, 0, t.TCPRemote, []byte{})
			return
		}
		// fmt.Println("HRL:", string(buf))

		err = netdata.TunnelTCPWrite(id, 1, t.TCPRemote, buf[:n])
		if err != nil {
			fmt.Println("err333", err)
			return
		}
	}
}

func HandleUDPRemoteHost(t *global.Tunnel, conn *net.UDPConn) {
	localaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(int(t.Port)))

	for {
		if t.Closed {
			return
		}

		id, addr, data, err := netdata.TunnelUDPRead(conn)
		if err != nil {
			return
		}
		if addr != t.UDPRemoteAddr {
			continue
		}
		// fmt.Println("HRR:", string(data))

		peerconn := UGetConnH(t, id)
		if peerconn == nil {
			peerconn, err := net.DialUDP("udp", nil, nil)
			if err != nil {
				fmt.Println("err111", err)
				continue
			}

			UAddConnH(t, id, peerconn)

			go HandleUDPLocalHost(t, peerconn, id)
		}

		peerconn = UGetConnH(t, id)

		_, err = peerconn.WriteTo(data, localaddr)
		if err != nil {
			fmt.Println("err222", err)
			UDelConnH(t, id)
		}
	}
}

func HandleUDPLocalHost(t *global.Tunnel, conn *net.UDPConn, id uint32) {
	localaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(int(t.Port)))

	for {
		if t.Closed {
			return
		}

		buf := make([]byte, 16*1024)

		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("err444", err)
			UDelConnH(t, id)
			return
		}
		if addr != localaddr {
			continue
		}
		// fmt.Println("HRL:", string(buf))

		err = netdata.TunnelUDPWrite(id, t.UDPRemote, t.UDPRemoteAddr, buf[:n])
		if err != nil {
			fmt.Println("err333", err)
			return
		}
	}
}
