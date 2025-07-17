package network

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"

	reuse "github.com/libp2p/go-reuseport"
)

func ProcEvent(conn net.Conn, data []byte) {
	event, err := protocol.Decode(string(data))
	if err != nil {
		global.App.Println("[ERRORdc]", err)
		return
	}

	err = protocol.VaildateTCPEvent(event...)
	if err != nil {
		global.App.Println("[BADEvent]", err, event)
		return
	}

	protocol.ExecTCPEvent(conn, event...)
}

func HandleConn(conn net.Conn) {
	for {
		data, err := network.ReadMsg(conn)
		if err != nil {
			global.App.Println("[ERRORrt]", err)
			global.App.Println("Disconnect from Server")
			return
		}

		global.App.Println("[Recv Server TCP]", string(data))

		ProcEvent(conn, data)
	}
}

func HandlePunchConn(conn net.Conn) {
	for {
		data, err := network.ReadMsg(conn)
		if err != nil {
			return
		}

		global.App.Println("[Recv Server TCP]", string(data))

		ProcEvent(conn, data)
	}
}

func ConnectPeer(conn net.Conn, peer_addr string, conntun chan net.Conn, dietun chan uint8) {
	for {
		select {
		case <-dietun:
			return
		default:
			time.Sleep(500 * time.Millisecond)
			peerConn, err := reuse.DialTimeout("tcp", conn.LocalAddr().String(), peer_addr, 2*time.Second)
			if err != nil {
				fmt.Println("Peer Conn Error:", err)
			} else {
				conntun <- peerConn
				fmt.Println("Peer Connected")
				return
			}
		}
	}
}

func ListenPeer(listener net.Listener, conn net.Conn, conntun chan net.Conn) { //host only
	peerConn, err := listener.Accept()
	if err != nil {
		// fmt.Println("Accept Error:", err)
		//it won't return a conn
		return
	}

	conntun <- peerConn
	// fmt.Println("Accepted Peer Connect")
}

func PunchPeer(conn net.Conn, peer_addr string, isHost bool) (net.Conn, error) {
	timeout := time.Second * 30
	timer := time.NewTimer(timeout)
	conntun := make(chan net.Conn)
	dietun := make(chan uint8, 1)

	if isHost {
		listener, _ := reuse.Listen("tcp", conn.LocalAddr().String())
		go ListenPeer(listener, conn, conntun)
		go ConnectPeer(conn, peer_addr, conntun, dietun)

		go func() {
			time.Sleep(timeout)
			listener.Close()
			conn.Close()
		}()
	} else {
		go ConnectPeer(conn, peer_addr, conntun, dietun)

		go func() {
			time.Sleep(timeout)
			conn.Close()
		}()
	}

	select {
	case peer_conn := <-conntun:
		dietun <- 1
		return peer_conn, nil
	case <-timer.C:
		return nil, errors.New("wait peer timeout")
	}

}
