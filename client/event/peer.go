package event

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	reuse "github.com/libp2p/go-reuseport"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/client/network"
	netdata "github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/terminal"
)

func InitPeerEvent() {
	protocol.RegTCPEvent(110, handleEnterRoom, protocol.IntType)
	protocol.RegTCPEvent(111, handleRoomFull, protocol.IntType)
	protocol.RegTCPEvent(120, handlePunchHostID, protocol.IntType)
	protocol.RegTCPEvent(122, handleNoticePunchPeer, protocol.IntType, protocol.StringType)
	protocol.RegTCPEvent(321, handleReqPwd)
	protocol.RegTCPEvent(320, handleEnterRoomSuc)
	protocol.RegTCPEvent(322, handlePwdError)
	protocol.RegTCPEvent(323, handleNameUsed)
	protocol.RegTCPEvent(324, handleNameList, protocol.StringAnyType)
	protocol.RegTCPEvent(325, handleKicked, protocol.StringType)
	protocol.RegTCPEvent(326, handleBanned, protocol.StringType, protocol.IntType)
	protocol.RegTCPEvent(340, handleMsg, protocol.StringType, protocol.StringType)
	protocol.RegTCPEvent(341, handleHostBroadcast, protocol.StringType)
	protocol.RegTCPEvent(342, handleMuted, protocol.StringType, protocol.IntType)
}

func handleEnterRoom(conn net.Conn, args ...any) {
	room_id := args[1].(int)

	for _, r := range global.Roomlist {
		if r.Id == uint32(room_id) {
			global.CurrRoom = r
		}
	}

	global.App.Println("Enter Room Success, Trying connect to host...")
	terminal.SetPrompt(global.App, global.CurrRoom.Name+">")
}

func handleRoomFull(conn net.Conn, args ...any) {
	room_id := args[1].(int)
	room_name := ""

	for _, r := range global.Roomlist {
		if r.Id == uint32(room_id) {
			room_name = r.Name
		}
	}

	global.App.Println(room_name+": Room is full, ID:", room_id)
}

func handlePunchHostID(conn net.Conn, args ...any) {
	punch_id := args[1].(int)

	global.Peer.Status = 1

	tmp_conn, err := reuse.Dial("tcp", "0.0.0.0:0", global.Serveraddr.String())
	if err != nil {
		fmt.Println("Punch connect server failed")
		return
	}

	go network.HandlePunchConn(tmp_conn)

	str, _ := protocol.Encode(203, punch_id)
	netdata.WriteMsg(tmp_conn, []byte(str))

	global.Peer.HostConn = tmp_conn
}

func handleNoticePunchPeer(conn net.Conn, args ...any) {
	switch global.Peer.Status {
	case 1:
		// punch_id := args[1].(int)
		host_addr := args[2].(string)

		host_conn, err := network.PunchPeer(global.Peer.HostConn, host_addr, false)
		if err != nil {
			fmt.Println("Connect to host failed")

			//clean up
			global.Peer.Status = 0
			global.Peer.HostConn.Close()
			terminal.SetPrompt(global.App, ">")
			global.CurrRoom = global.Room{}
			return
		}

		global.Peer.Status = 2
		global.Peer.HostConn = host_conn
		//handle host conn
		go func() {
			netdata.WriteMsg(host_conn, []byte("200"))

			for {
				data, err := netdata.ReadMsg(host_conn)
				if err != nil {
					fmt.Println("Disconnect from host")
					return
				}

				fmt.Println("From Host recv:", string(data))

				network.ProcEvent(data)
			}
		}()

	case 2:
		//port
	}
}

func handleReqPwd(conn net.Conn, args ...any) {
	input := bufio.NewReader(os.Stdin)

	global.App.Print("Enter Password:") //prompt
	data, _, _ := input.ReadLine()
	pwd := string(data)

	str, _ := protocol.Encode(210, pwd)
	netdata.WriteMsg(conn, []byte(str))
}

func handleEnterRoomSuc(conn net.Conn, args ...any) {
	input := bufio.NewReader(os.Stdin)

	global.App.Print("Auth Success, Enter Name:") //prompt
	data, _, _ := input.ReadLine()
	name := string(data)

	str, _ := protocol.Encode(211, name)
	netdata.WriteMsg(conn, []byte(str))
}

func handlePwdError(conn net.Conn, args ...any) {
	global.App.Println("Password Incorrect, Please enter the room again.")
}

func handleNameUsed(conn net.Conn, args ...any) {
	global.App.Println("Name already used, Please enter the room again.")
}

func handleKicked(conn net.Conn, args ...any) {
	reason := args[1].(string)
	global.App.Println("You have been kicked, Reason:", reason)
}

func handleBanned(conn net.Conn, args ...any) {
	reason := args[1].(string)
	time := time.Unix(int64(args[2].(int)), 0)
	global.App.Println("You have been banned, Reason:", reason)
	global.App.Println("Unban time:", time)
}

func handleNameList(conn net.Conn, args ...any) {
	names := ""

	for i := 0; i < len(args); i++ {
		if i == 0 {
			continue
		}

		if i == 1 {
			names = names + args[i].(string)
		}

		names = names + "," + args[i].(string)
	}

	global.App.Println("Players:", names)
}

func handleMsg(conn net.Conn, args ...any) {
	name := args[1].(string)
	msg := args[2].(string)

	global.App.Println("<" + name + ">" + msg)
}

func handleHostBroadcast(conn net.Conn, args ...any) {
	msg := args[1].(string)
	global.App.Println("Host Broadcast:", msg)
}

func handleMuted(conn net.Conn, args ...any) {
	reason := args[1].(string)
	time := time.Unix(int64(args[2].(int)), 0)
	global.App.Println("You have been muted, Reason:", reason)
	global.App.Println("Unmute time:", time)
}
