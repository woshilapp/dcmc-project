package network

import (
	"encoding/binary"
	"errors"
	"io"
	"net"

	"github.com/woshilapp/dcmc-project/client/global"
)

const MaxMsgLength = 1024 * 100 //100KB

func WriteMsg(conn net.Conn, data []byte) error {
	header := []byte{0xdc, 0x3c} //DC3C(HEX), 56380(DEC), 2bytes

	lendata := make([]byte, 4)
	binary.BigEndian.PutUint32(lendata, uint32(len(data)))
	header = append(header, lendata...)

	fullMsg := append(header, data...)

	_, err := conn.Write(fullMsg)
	return err
}

func ReadMsg(conn net.Conn) ([]byte, error) {
	header := make([]byte, 6)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	proto := binary.BigEndian.Uint16(header[:2])
	if proto != 56380 {
		return nil, errors.New("not dcmc protocol")
	}

	length := binary.BigEndian.Uint32(header[2:])
	if length > MaxMsgLength {
		println(length)
		return nil, errors.New("message too long")
	}

	body := make([]byte, length)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, err
	}

	return body, nil
}

func TunnelTCPWrite(id uint32, status uint16, conn net.Conn, data []byte) error {
	body := make([]byte, 4)
	binary.BigEndian.PutUint32(body, id)

	statbin := make([]byte, 2)
	binary.BigEndian.PutUint16(statbin, status)
	body = append(body, statbin...)

	body = append(body, data...)

	err := WriteMsg(conn, body)
	if err != nil {
		return err
	}

	return nil
}

func TunnelTCPRead(conn net.Conn) (uint32, uint16, []byte, error) {
	data, err := ReadMsg(conn)
	if err != nil {
		return 0, 0, nil, err
	}

	id := binary.BigEndian.Uint32(data)
	status := binary.BigEndian.Uint16(data[4:6])

	return id, status, data[6:], nil
}

func TunnelUDPWrite(id uint32, conn *net.UDPConn, dest net.Addr, data []byte) error {
	dataf := []byte{0xdc, 0x3c}

	body := make([]byte, 4)
	binary.BigEndian.PutUint32(body, id)
	dataf = append(dataf, body...)

	dataf = append(dataf, data...)

	_, err := conn.WriteTo(dataf, dest)
	if err != nil {
		return err
	}

	return nil
}

func TunnelUDPRead(conn *net.UDPConn) (uint32, net.Addr, []byte, error) {
	data := make([]byte, 64*1024)

	n, addr, err := conn.ReadFrom(data)
	if err != nil {
		return 0, nil, nil, err
	}
	global.App.Println("UDPTUN: ", data[:n])

	header := data[:2]
	proto := binary.BigEndian.Uint16(header[:2])
	if proto != 56380 {
		return 0, nil, nil, errors.New("not dcmc protocol")
	}

	id := binary.BigEndian.Uint32(data[2:6])

	return id, addr, data[6:n], nil
}
