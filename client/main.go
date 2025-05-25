package main

import (
	"fmt"
	"net"

	"github.com/woshilapp/dcmc-project/client/event"
	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/client/shell"
	term "github.com/woshilapp/dcmc-project/terminal"
)

func recvUDP() {
	app := global.App

	for {
		buf := make([]byte, 1440)

		n, addr, err := global.Udpsock.ReadFrom(buf)
		if err != nil {
			app.Println("[ERRORru]", err)
		}

		app.Println("[Recv UDP From]", addr.String(), "say:", string(buf[:n]))

	}
}

func main() {
	//Welcome text
	fmt.Println("Hello, world client!")

	global.App = term.NewTerminal("dcmc-project client")

	global.Serveraddr, _ = net.ResolveUDPAddr("udp", "127.0.0.1:7789")
	localaddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	global.Udpsock, _ = net.ListenUDP("udp", localaddr)

	go recvUDP()

	shell.InitCommand()

	if global.Role == 1 {
		event.InitPeerEvent()
	} else {
		event.InitHostEvent()
	}

	global.App.Run()
}
