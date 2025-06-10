package main

import (
	"fmt"
	"net"
	"time"

	"github.com/woshilapp/dcmc-project/protocol"
)

func execsb(conn net.Conn, args ...any) {
	fmt.Println("do 100 sb")

	for _, v := range args {
		fmt.Println("sb num and str:", v.(int))
	}
}

func main1() {
	protocol.Run()

	// protocol.RegTCPEvent(100, protocol.IntType, protocol.BoolType, protocol.StringType)

	// err := protocol.VaildateTCPEvent(100, 200, true, "sb")
	var unconn net.Conn

	// unconn := &net.TCPConn{}

	protocol.RegTCPEvent(100, execsb, protocol.IntType, protocol.StringAnyType)

	// event := []any{100, 114514, true}
	event := []any{100, 114514, "sdasasd", "sadasd", "sada"}
	err := protocol.VaildateTCPEvent(event...)
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	protocol.ExecTCPEvent(unconn, event...)

	time.Sleep(1 * time.Second) //wait
}
