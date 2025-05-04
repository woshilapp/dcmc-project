package protocol

import (
	"errors"
	"fmt"
	"net"
)

var (
	IntType    = "int"
	StringType = "string"
	BoolType   = "bool"
	FloatType  = "float64"
)

type event struct {
	types []string
	len   int
	exec  func(net.Conn, ...any) //conn, code, args
}

var events = map[int]event{} //code:event

func typeof(a any) string {
	return fmt.Sprintf("%T", a)
}

// Registered event
func RegEvent(code int, exec func(net.Conn, ...any), types ...string) {
	events[code] = event{
		types: types,
		len:   len(types) + 1,
		exec:  exec,
	}
}

// Vaildate event
func VaildateEvent(a ...any) error {
	if len(a) < 1 {
		return errors.New("null event")
	}

	if typeof(a[0]) != IntType {
		return errors.New("illegal event")
	}

	e := a[0].(int)

	if _, exist := events[e]; !exist {
		return errors.New("unregistered event")
	}

	if len(a) != events[e].len {
		return errors.New("len mismatch")
	}

	for i := 0; i < len(a)-1; i++ {
		if typeof(a[i+1]) != events[e].types[i] {
			return errors.New("type mismatch")
		}
	}

	return nil
}

// Before exec an event, it must be vaildated
func ExecEvent(conn net.Conn, a ...any) {
	e := events[a[0].(int)]

	go e.exec(conn, a...)
}
