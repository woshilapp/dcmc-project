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

	term.AddCommand(global.App, "sendpeer", "send tcp data to peer(s)",
		[]string{"text"}, "",
		sendToPeer)

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

	term.AddCommand(global.App, "passwd", "set a passwd for room",
		[]string{"pwd"}, "",
		setPasswd)

	term.AddCommand(global.App, "clearpwd", "clear the passwd",
		[]string{}, "",
		clrPasswd)

	term.AddCommand(global.App, "name", "set the name",
		[]string{"name"}, "",
		setName)

	term.AddCommand(global.App, "namelist", "list names on host",
		[]string{}, "",
		reqNamelist)

	term.AddCommand(global.App, "msg", "send msg",
		[]string{"msg"}, "",
		sendMsg)
}

func connectToServer(context *grumble.Context) error {
	conn, err := net.Dial("tcp", context.Args.String("addr"))
	if err != nil {
		context.App.Println("[ERRORc]", err)

		return nil
	}

	global.Serverconn = conn
	global.Serveraddr, _ = net.ResolveTCPAddr("tcp", context.Args.String("addr"))

	go network.HandleConn(conn)

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

func sendToPeer(context *grumble.Context) error {
	if global.Role == 1 {
		err := netdata.WriteMsg(global.Peer.HostConn, []byte(context.Args.String("text")))

		if err != nil {
			context.App.Println("[ERRORsp]", err)

			return nil
		}

		return nil
	} else {
		for _, p := range global.Host.Peers {
			err := netdata.WriteMsg(p.Conn, []byte(context.Args.String("text")))

			if err != nil {
				context.App.Println("[ERRORsp]", err)

				return nil
			}
		}
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
	}

	global.Host.Passwd = pwd
	global.Host.Status = 1

	str, _ := protocol.Encode(310, global.CurrRoom.Name,
		global.CurrRoom.MaxPeer,
		global.CurrRoom.Description,
		global.CurrRoom.RequiredPwd,
	)

	netdata.WriteMsg(global.Serverconn, []byte(str))

	return nil
}

func setPasswd(context *grumble.Context) error {
	switch global.Role {
	case 1:
		pwd := context.Args.String("pwd")
		str, _ := protocol.Encode(210, pwd)
		netdata.WriteMsg(global.Peer.HostConn, []byte(str))
	case 2:
		if global.Host.Status == 0 {
			return nil
		}

		global.Host.Passwd = context.Args.String("pwd")
		global.CurrRoom.RequiredPwd = true

		str, _ := protocol.Encode(312, global.CurrRoom.Id, global.CurrRoom.Description, true)

		netdata.WriteMsg(global.Serverconn, []byte(str))
	}

	return nil
}

func clrPasswd(context *grumble.Context) error {
	if global.Host.Status == 0 {
		return nil
	}

	global.Host.Passwd = ""
	global.CurrRoom.RequiredPwd = false

	str, _ := protocol.Encode(312, global.CurrRoom.Id, global.CurrRoom.Description, false)

	netdata.WriteMsg(global.Serverconn, []byte(str))

	return nil
}

func setName(context *grumble.Context) error {
	if global.Role != 1 {
		return nil
	}

	name := context.Args.String("name")

	str, _ := protocol.Encode(211, name)
	netdata.WriteMsg(global.Peer.HostConn, []byte(str))

	return nil
}

func reqNamelist(context *grumble.Context) error {
	switch global.Role {
	case 1:
		str, _ := protocol.Encode(212)
		netdata.WriteMsg(global.Peer.HostConn, []byte(str))
	case 2:
		str := ""
		for _, p := range global.Host.Peers {
			if p.Name == "" {
				continue
			}

			if str == "" {
				str = str + p.Name
				continue
			}

			str = str + p.Name + ","
		}

		global.App.Println("Players:", str)
	}

	return nil
}

func sendMsg(context *grumble.Context) error {
	msg := context.Args.String("msg")

	switch global.Role {
	case 1:
		str, _ := protocol.Encode(230, msg)
		netdata.WriteMsg(global.Peer.HostConn, []byte(str))
	case 2:
		str := "<Host>" + msg

		for _, p := range global.Host.Peers {
			if p.Name != "" {
				netdata.WriteMsg(p.Conn, []byte(str))
			}
		}
	}

	return nil
}
