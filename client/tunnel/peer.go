package tunnel

import (
	"fmt"
	"net"
	"strconv"

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
		// fmt.Println("PRR:", string(data))

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

		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err222", err)

			TDelConnP(t, id)
			netdata.TunnelTCPWrite(id, 0, t.TCPRemote, []byte{})

			return
		}
		// fmt.Println("PRL:", string(buf))

		err = netdata.TunnelTCPWrite(id, 1, t.TCPRemote, buf)
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
