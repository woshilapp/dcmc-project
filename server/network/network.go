package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/woshilapp/dcmc-project/protocol"
	_ "github.com/woshilapp/dcmc-project/server/event"
	"github.com/woshilapp/dcmc-project/server/global"
	"github.com/woshilapp/dcmc-project/server/punching"
)

const MaxMsgLength = 1024 * 10 //10KB

func WriteMsg(conn net.Conn, data []byte) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))

	fullMsg := append(header, data...)

	_, err := conn.Write(fullMsg)
	return err
}

func ReadMsg(conn net.Conn) ([]byte, error) {
	// set timeout
	// conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	// defer conn.SetReadDeadline(time.Time{}) // reset

	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header)
	if length > MaxMsgLength {
		return nil, errors.New("message too long")
	}

	body := make([]byte, length)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, err
	}

	return body, nil
}

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

	global.UDPconn = udpconn

	buf := make([]byte, 1440)

	n, connaddr, err := udpconn.ReadFrom(buf)
	if err != nil {
		errchan <- err
	}

	event, err := protocol.Decode(string(buf[:n]))
	if err != nil {
		errchan <- err
	}

	go punching.HandleUDPPunch(event, connaddr)
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
		data, err := ReadMsg(conn)
		if err != nil {
			errchan <- err
			deadChan <- conn.RemoteAddr().String()
			break
		}

		fmt.Println("[Recv] From", conn.RemoteAddr().String(), "say:", string(data))

		event, err := protocol.Decode(string(data))
		if err != nil {
			fmt.Println("[Recv] Decode Error:", err)
			continue
		}

		fmt.Println("[Recv] Event", event)

		err = protocol.VaildateEvent(event...)
		if err != nil {
			fmt.Println("[Recv] Bad event:", err)
			continue
		}

		protocol.ExecEvent(conn, event...)
	}
}
