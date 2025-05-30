package shell

import (
	"net"
	"strconv"

	"github.com/desertbit/grumble"
	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/client/network"
	netdata "github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
	term "github.com/woshilapp/dcmc-project/terminal"
)

func InitCommand() {
	term.AddCommand(global.App, "connect", "connect to server", []string{"addr"}, connectToServer)
	term.AddCommand(global.App, "send", "send tcp data to server", []string{"text"}, sendToServer)
	term.AddMultiArgCommand(global.App, "sendencode", "send tcp encoded data to server", "data", sendEncodedToServer)
	term.AddMultiArgCommand(global.App, "sendudpencode", "send udp encoded data to server", "data", sendUDPEncodedToServer)
	term.AddCommand(global.App, "list", "list rooms on server", []string{}, listRoom)
}

func connectToServer(context *grumble.Context) error {
	conn, err := net.Dial("tcp", context.Args.String("addr"))
	if err != nil {
		context.App.Println("[ERRORc]", err)

		return nil
	}

	global.Serverconn = conn

	go network.ListenConn(conn)

	str, _ := protocol.Encode(200)
	netdata.WriteMsg(conn, []byte(str))

	return nil
}

func sendToServer(context *grumble.Context) error {
	// _, err := serverconn.Write([]byte(context.Args.String("text")))
	err := netdata.WriteMsg(global.Serverconn, []byte(context.Args.String("text")))

	if err != nil {
		context.App.Println("[ERRORst]", err)

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
	err := netdata.WriteMsg(global.Serverconn, data)

	if err != nil {
		context.App.Println("[ERRORst]", err)

		return nil
	}

	return nil
}

func sendUDPEncodedToServer(context *grumble.Context) error {
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

	_, err := global.Udpsock.WriteTo(data, global.Serveraddr)

	if err != nil {
		context.App.Println("[ERRORsu]", err)

		return nil
	}

	return nil
}

func listRoom(context *grumble.Context) error {
	global.Roomlist = []global.Room{}

	str, _ := protocol.Encode(202)
	netdata.WriteMsg(global.Serverconn, []byte(str))

	return nil
}
