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

		id, data, err := netdata.TunnelTCPRead(conn)
		if err != nil {
			fmt.Println("RemoteERR")
		}

		peerconn := TGetConnP(t, id)
		if peerconn == nil {
			continue
		}

		_, err = peerconn.Write(data)
		if err != nil {
			fmt.Println("err1111")
			TDelConnP(t, id)
		}
	}
}

func HandleLocalPeer(t *global.Tunnel, conn net.Conn, id uint32) {
	for {
		if t.Closed {
			return
		}

		data := []byte{}
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err222")
			TDelConnP(t, id)
			return
		}

		err = netdata.TunnelTCPWrite(id, t.TCPRemote, data)
		if err != nil {
			fmt.Println("errwrt")
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
