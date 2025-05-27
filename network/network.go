package network

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
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
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header)
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
