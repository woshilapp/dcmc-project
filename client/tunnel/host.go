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

		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("err444", err)
			TDelConnH(t, id)

			netdata.TunnelTCPWrite(id, 0, t.TCPRemote, []byte{})
			return
		}
		// fmt.Println("HRL:", string(buf))

		err = netdata.TunnelTCPWrite(id, 1, t.TCPRemote, buf)
		if err != nil {
			fmt.Println("err333", err)
			return
		}
	}
}
