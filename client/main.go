package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/desertbit/grumble"
	"github.com/woshilapp/dcmc-project/client/network"
	"github.com/woshilapp/dcmc-project/protocol"
	term "github.com/woshilapp/dcmc-project/terminal"
)

var serverconn net.Conn

func makeServerConn(context *grumble.Context) error {
	conn, err := net.Dial("tcp", context.Args.String("addr"))
	if err != nil {
		context.App.Println("[ERROR]", err)

		return nil
	}

	serverconn = conn

	go func() {
		for {
			// buf := make([]byte, 1024)
			// n, err := conn.Read(buf)
			// if err != nil {
			// 	context.App.Println("[ERROR]", err)

			// 	return
			// }

			data, err := network.ReadMsg(conn)
			if err != nil {
				context.App.Println("[ERROR]", err)

				return
			}

			context.App.Println("[Recv Server]", string(data))
		}
	}()

	return nil
}

func sendToServer(context *grumble.Context) error {
	// _, err := serverconn.Write([]byte(context.Args.String("text")))
	err := network.WriteMsg(serverconn, []byte(context.Args.String("text")))

	if err != nil {
		context.App.Println("[ERROR]", err)

		return nil
	}

	return nil
}

func sendEncodedToServer(context *grumble.Context) error {
	args := []any{}
	for _, v := range context.Args.StringList("data") {
		if vi, err := strconv.Atoi(v); err == nil { //int
			args = append(args, vi)
		} else if vf, err := strconv.ParseFloat(v, 64); err == nil { //float
			args = append(args, vf)
		} else if v == "true" { //bool
			args = append(args, true)
		} else if v == "false" { //bool
			args = append(args, false)
		} else { //string
			args = append(args, v)
		}
	}

	tmp, _ := protocol.Encode(args...)
	data := []byte(tmp)

	// _, err := serverconn.Write(data)
	err := network.WriteMsg(serverconn, data)

	if err != nil {
		context.App.Println("[ERROR]", err)

		return nil
	}

	return nil
}

func main() {
	//Welcome text
	fmt.Println("Hello, world client!")

	app := term.NewTerminal("dcmc-project client")

	//command example

	// term.AddCommand(app, "fku", "fku mom", []string{}, func(c *grumble.Context) error {
	// 	app.Println("cnm666")
	// 	return nil
	// })

	// term.AddCommand(app, "fking", "fking anyone", []string{"who"}, func(c *grumble.Context) error {
	// 	go func() {
	// 		time.Sleep(1 * time.Second)
	// 		app.Println("[FUCKER] fking", c.Args.String("who"))
	// 	}()

	// 	return nil
	// })

	term.AddCommand(app, "connect", "connect to server", []string{"addr"}, makeServerConn)

	term.AddCommand(app, "send", "send data to server", []string{"text"}, sendToServer)

	term.AddMultiArgCommand(app, "sendencode", "send encoded data to server", "data", sendEncodedToServer)

	app.Run()
}
