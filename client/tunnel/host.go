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
)

func TCPPunchHost(peer *global.TPeers, port uint16) {
	netdata.WriteMsg(global.Serverconn, []byte("301"))

	punch_id := <-global.Host.PunchIDs
	tun := &global.Tunnel{
		Proto:   1,
		Port:    port,
		PunchID: punch_id,
		Lock:    sync.RWMutex{},
		Closed:  false,
	}
	peer.Tunnels = append(peer.Tunnels, tun)

	tmp_conn, err := reuse.Dial("tcp", "0.0.0.0:0", global.Serveraddr.String())
	if err != nil {
		fmt.Println("Punch connect server failed")
		return
	}

	go network.HandlePunchConn(tmp_conn)

	netdata.WriteMsg(tmp_conn, []byte("302"))

	tun.TCPRemote = tmp_conn
	global.Host.PIDtun[punch_id] = tun
}

func HandleNewPeer(peer *global.TPeers) {
	for _, port := range global.Host.TCPPorts {
		TCPPunchHost(peer, port)
	}

	for _, port := range global.Host.UDPPorts {
		print(port)
	}
}

// Handles
func HandleRemoteHost(t *global.Tunnel, conn net.Conn) {
	for {
		if t.Closed {
			return
		}

		id, data, err := netdata.TunnelTCPRead(conn)
		if err != nil {
			return
		}

		peerconn := TGetConnH(t, id)
		if peerconn == nil {
			peerconn, err := net.Dial("tcp", "127.0.0.1"+strconv.Itoa(int(t.Port)))
			if err != nil {
				fmt.Println("err111")
				continue
			}

			TAddConnH(t, id, peerconn)

			go HandleLocalHost(t, peerconn, id)
		}

		peerconn = TGetConnH(t, id)

		_, err = peerconn.Write(data)
		if err != nil {
			fmt.Println("err222")
			TDelConnH(t, id)
		}
	}
}

func HandleLocalHost(t *global.Tunnel, conn net.Conn, id uint32) {
	for {
		if t.Closed {
			return
		}

		data := []byte{}
		buf := make([]byte, 1024)

		_, err := conn.Read(buf)
		if err != nil {
			TDelConnH(t, id)
			return
		}

		err = netdata.TunnelTCPWrite(id, t.TCPRemote, data)
		if err != nil {
			fmt.Println("err333")
			return
		}
	}
}
