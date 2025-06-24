package protocol

import (
	"errors"
	"fmt"
	"net"
)

const (
	IntType       = "int"
	StringType    = "string"
	BoolType      = "bool"
	FloatType     = "float64"
	StringAnyType = "stringany" //any number string
)

type tcpevent struct {
	types []string
	len   int
	exec  func(conn net.Conn, args ...any) //conn, code, args
}

type udpevent struct {
	types []string
	len   int
	exec  func(conn *net.UDPConn, addr net.Addr, args ...any) //conn, addr, code, args
}

var tcpevents = map[int]tcpevent{} //code:TCPevent
var udpevents = map[int]udpevent{} //code:UDPevent

func typeof(a any) string {
	return fmt.Sprintf("%T", a)
}

// Registered event
func RegTCPEvent(code int, exec func(net.Conn, ...any), types ...string) {
	tcpevents[code] = tcpevent{
		types: types,
		len:   len(types) + 1,
		exec:  exec,
	}
}

func RegUDPEvent(code int, exec func(*net.UDPConn, net.Addr, ...any), types ...string) {
	udpevents[code] = udpevent{
		types: types,
		len:   len(types) + 1,
		exec:  exec,
	}
}

// Vaildate Event
func VaildateTCPEvent(a ...any) error {
	if len(a) < 1 {
		return errors.New("null event")
	}

	if typeof(a[0]) != IntType {
		return errors.New("illegal event")
	}

	e := a[0].(int)
	event, exist := tcpevents[e]
	if !exist {
		return errors.New("unregistered event")
	}

	if len(a) != event.len && event.len != 1 {
		for i, v := range event.types {
			if v != "stringany" && i == event.len-2 {
				return errors.New("len mismatch")
			}
		}
	}

	for i := 0; i < len(a)-1; i++ {
		if event.types[i] == "stringany" {
			for j := i; j < len(a)-1; j++ {
				if typeof(a[j+1]) != "string" {
					return errors.New("type mismatch")
				}
			}

			return nil
		}

		if typeof(a[i+1]) != event.types[i] {
			return errors.New("type mismatch")
		}
	}

	return nil
}

func VaildateUDPEvent(a ...any) error {
	if len(a) < 1 {
		return errors.New("null event")
	}

	if typeof(a[0]) != IntType {
		return errors.New("illegal event")
	}

	e := a[0].(int)
	event, exist := udpevents[e]
	if !exist {
		return errors.New("unregistered event")
	}

	if len(a) != event.len {
		for i, v := range event.types {
			if v != "stringany" && i == event.len-1 {
				return errors.New("len mismatch")
			}
		}
	}

	for i := 0; i < len(a)-1; i++ {
		if event.types[i] == "stringany" {
			for j := i; j < len(a)-1; j++ {
				if typeof(a[j+1]) != "string" {
					return errors.New("type mismatch")
				}
			}

			return nil
		}

		if typeof(a[i+1]) != event.types[i] {
			return errors.New("type mismatch")
		}
	}

	return nil
}

// Before exec an event, it must be vaildated
func ExecTCPEvent(conn net.Conn, a ...any) error {
	e := tcpevents[a[0].(int)]

	go e.exec(conn, a...)
	return nil
}

func ExecUDPEvent(sock *net.UDPConn, addr net.Addr, a ...any) error {
	e := udpevents[a[0].(int)]

	go e.exec(sock, addr, a...)
	return nil
}
