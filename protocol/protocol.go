package protocol

import (
	"errors"
	"fmt"
	"net"
)

const (
	IntType    = "int"
	StringType = "string"
	BoolType   = "bool"
	FloatType  = "float64"
)

type event struct {
	types []string
	len   int
	exec  func(conn any, addr net.Addr, args ...any) //conn, addr(udp), code, args
}

var events = map[int]event{} //code:event

func typeof(a any) string {
	return fmt.Sprintf("%T", a)
}

// Registered event
func regEvent(code int, exec func(conn any, addr net.Addr, args ...any), types ...string) {
	events[code] = event{
		types: types,
		len:   len(types) + 1,
		exec:  exec,
	}
}

// Vaildate Event
func VaildateEvent(a ...any) error {
	if len(a) < 1 {
		return errors.New("null event")
	}

	if typeof(a[0]) != IntType {
		return errors.New("illegal event")
	}

	e := a[0].(int)
	event, exist := events[e]
	if !exist {
		return errors.New("unregistered event")
	}

	if len(a) != event.len {
		return errors.New("len mismatch")
	}

	for i := 0; i < len(a)-1; i++ {
		if typeof(a[i+1]) != event.types[i] {
			return errors.New("type mismatch")
		}
	}

	return nil
}

// Before exec an event, it must be vaildated
func execEvent(conn any, addr net.Addr, a ...any) error {
	e := events[a[0].(int)]

	go e.exec(conn, addr, a...)
	return nil
}

// quick functions
func RegTCPEvent(code int, exec func(net.Conn, ...any), types ...string) {
	wrapper := func(conn any, _ net.Addr, args ...any) {
		exec(conn.(net.Conn), args...)
	}
	regEvent(code, wrapper, types...)
}

func RegUDPEvent(code int, exec func(*net.UDPConn, net.Addr, ...any), types ...string) {
	wrapper := func(conn any, addr net.Addr, args ...any) {
		exec(conn.(*net.UDPConn), addr, args...)
	}
	regEvent(code, wrapper, types...)
}

func ExecTCPEvent(conn net.Conn, a ...any) error {
	return execEvent(conn, &net.UDPAddr{}, a...)
}

func ExecUDPEvent(sock *net.UDPConn, addr net.Addr, a ...any) error {
	return execEvent(sock, addr, a...)
}
