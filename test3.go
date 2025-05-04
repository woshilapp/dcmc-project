package main

import (
	"github.com/woshilapp/dcmc-project/protocol"
)

// func execsb(conn net.Conn, args ...any) {
// 	fmt.Println("do 100 sb")

// 	for _, v := range args {
// 		fmt.Println("sb num:", v.(int))
// 	}
// }

func main() {
	protocol.Run()

	// protocol.RegEvent(100, protocol.IntType, protocol.BoolType, protocol.StringType)

	// result, err := protocol.VaildateEvent(100, 200, true, "sb")
	// var unconn net.Conn

	// 	unconn := &net.TCPConn{}

	// 	protocol.RegEvent(100, execsb, protocol.IntType)

	// 	event := []any{100, 114514}
	// 	err := protocol.VaildateEvent(event...)
	// 	if err != nil {
	// 		fmt.Println("[ERROR]", err)
	// 		return
	// 	}

	// 	protocol.ExecEvent(unconn, event...)

	// 	time.Sleep(1 * time.Second) //wait
}
