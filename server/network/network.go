package network

import (
	"fmt"
	"net"

	"github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
	_ "github.com/woshilapp/dcmc-project/server/event"
	"github.com/woshilapp/dcmc-project/server/global"
)

func ListenServer(addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	return listener, nil
}

func ListenUDP(addr string, errchan chan error) {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		errchan <- err
		return
	}

	udpconn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		errchan <- err
		return
	}

	for {
		buf := make([]byte, 1440)

		n, connaddr, err := udpconn.ReadFrom(buf)
		if err != nil {
			errchan <- err
		}

		fmt.Println("[Recv UDP] From", connaddr, "say:", string(buf[:n]))

		event, err := protocol.Decode(string(buf[:n]))
		if err != nil {
			fmt.Println("[Recv UDP] Decode Error:", err)
		}

		fmt.Println("[Recv] Event", event)

		err = protocol.VaildateUDPEvent(event...)
		if err != nil {
			fmt.Println("[Recv] Bad event:", err)
			continue
		}

		protocol.ExecUDPEvent(udpconn, connaddr, event...)
	}
}

func AccpetConn(listener net.Listener, errchan chan error) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			errchan <- err
		}

		go HandleConn(conn, errchan)
	}
}

func cleanRoom(addr chan string) {
	host_addr := <-addr

	for _, r := range global.GetRoomlist() {
		if r.HostConn.RemoteAddr().String() == host_addr {
			global.RemoveRoom(int(r.Id))
		}
	}
}

func HandleConn(conn net.Conn, errchan chan error) {
	fmt.Println("New Conn, Connaddr:", conn.RemoteAddr().String())

	deadChan := make(chan string, 1)
	go cleanRoom(deadChan)

	for {
		data, err := network.ReadMsg(conn)
		if err != nil {
			errchan <- err
			deadChan <- conn.RemoteAddr().String()
			break
		}

		fmt.Println("[Recv TCP] From", conn.RemoteAddr().String(), "say:", string(data))

		event, err := protocol.Decode(string(data))
		if err != nil {
			fmt.Println("[Recv TCP] Decode Error:", err)
			continue
		}

		fmt.Println("[Recv] Event", event)

		err = protocol.VaildateTCPEvent(event...)
		if err != nil {
			fmt.Println("[Recv] Bad event:", err)
			continue
		}

		protocol.ExecTCPEvent(conn, event...)
	}
}
