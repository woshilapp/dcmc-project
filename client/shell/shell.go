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
	term.AddCommand(global.App, "connect", "connect to server",
		[]string{"addr"}, "",
		connectToServer)

	term.AddCommand(global.App, "send", "send tcp data to server",
		[]string{"text"}, "",
		sendToServer)

	term.AddCommand(global.App, "sendencode", "send tcp encoded data to server",
		[]string{}, "data",
		sendEncodedToServer)

	term.AddCommand(global.App, "sendudpencode", "send udp encoded data to server",
		[]string{}, "data",
		sendUDPEncodedToServer)

	term.AddCommand(global.App, "list", "list rooms on server",
		[]string{}, "",
		listRoom)

	term.AddCommand(global.App, "enter", "enter a room",
		[]string{"id"}, "",
		enterRoom)

	term.AddCommand(global.App, "create", "create a room",
		[]string{"name", "max_peer", "desc"}, "pwd",
		createRoom)
}

func connectToServer(context *grumble.Context) error {
	conn, err := net.Dial("tcp", context.Args.String("addr"))
	if err != nil {
		context.App.Println("[ERRORc]", err)

		return nil
	}

	global.Serverconn = conn

	go network.ListenConn(conn)

	var helloInt int

	if global.Role == 1 {
		helloInt = 200
	} else {
		helloInt = 300
	}

	str, _ := protocol.Encode(helloInt) //send hello to server
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

func enterRoom(context *grumble.Context) error {
	if global.Role != 1 {
		return nil
	}

	id, err := strconv.Atoi(context.Args.String("id"))
	if err != nil {
		global.App.Println("Bad room id")
		return nil
	}

	str, _ := protocol.Encode(201, id)
	netdata.WriteMsg(global.Serverconn, []byte(str))

	return nil
}

func createRoom(context *grumble.Context) error {
	if global.Role != 2 {
		return nil
	}

	maxpeer, err := strconv.Atoi(context.Args.String("max_peer"))
	if err != nil {
		return nil
	}

	arglist := context.Args.StringList("pwd")
	reqpwd := false
	pwd := ""

	if len(arglist) != 0 {
		reqpwd = true
		pwd = arglist[0]
	}

	global.CurrRoom = global.Room{
		Name:        context.Args.String("name"),
		Description: context.Args.String("desc"),
		MaxPeer:     maxpeer,
		CurrPeer:    0,
		RequiredPwd: reqpwd,
		Passwd:      pwd,
	}

	str, _ := protocol.Encode(310, global.CurrRoom.Name,
		global.CurrRoom.MaxPeer,
		global.CurrRoom.Description,
		global.CurrRoom.RequiredPwd,
	)

	netdata.WriteMsg(global.Serverconn, []byte(str))

	return nil
}
