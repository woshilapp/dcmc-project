package network

import (
	"fmt"
	"net"
	"time"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"

	reuse "github.com/libp2p/go-reuseport"
)

func HandleConn(conn net.Conn) {
	for {
		data, err := network.ReadMsg(conn)
		if err != nil {
			global.App.Println("[ERRORrt]", err)
			return
		}

		global.App.Println("[Recv Server TCP]", string(data))

		event, err := protocol.Decode(string(data))
		if err != nil {
			global.App.Println("[ERRORdc]", err)
			continue
		}

		err = protocol.VaildateTCPEvent(event...)
		if err != nil {
			global.App.Println("[BADEvent]", err, event)
			continue
		}

		protocol.ExecTCPEvent(global.Serverconn, event...)
	}
}

func ConnectPeer(conn net.Conn, peer_addr string, conntun chan net.Conn) {
	for {
		time.Sleep(500 * time.Millisecond)
		peerConn, err := reuse.DialTimeout("tcp", conn.LocalAddr().String(), peer_addr, 2*time.Second)
		if err != nil {
			fmt.Println("Peer Conn Error:", err)
		} else {
			conntun <- peerConn
			fmt.Println("Peer Connected")
		}
	}
}

func ListenPeer(listener net.Listener, conn net.Conn, conntun chan net.Conn) { //host only
	peerConn, err := listener.Accept()
	if err != nil {
		fmt.Println("Accept Error:", err)
	}

	conntun <- peerConn
	fmt.Println("Accepted Peer Connect")
}

func PunchPeer(conn net.Conn, peer_addr string, isHost bool) (net.Conn, error) {
	timeout := time.Second * 30
	conntun := make(chan net.Conn)

	if isHost {
		listener, _ := reuse.Listen("tcp", conn.LocalAddr().String())
		go ListenPeer(listener, conn, conntun)
		go ConnectPeer(conn, peer_addr, conntun)

		go func() {
			time.Sleep(timeout)
			listener.Close()
			conn.Close()
		}()
	} else {
		go ConnectPeer(conn, peer_addr, conntun)

		go func() {
			time.Sleep(timeout)
			conn.Close()
		}()
	}

	peer_conn := <-conntun

	return peer_conn, nil
}
