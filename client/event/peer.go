package event

import (
	"fmt"
	"net"

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
}

func handleEnterRoom(conn net.Conn, args ...any) {
	global.CurrRoom = global.Roomlist[args[1].(int)]

	global.App.Println("Enter Room Success, Trying connect to host...")
	terminal.SetPrompt(global.App, global.CurrRoom.Name+">")
}

func handleRoomFull(conn net.Conn, args ...any) {
	global.App.Println(global.Roomlist[args[1].(int)].Name + ": Room is full")
}

func handlePunchHostID(conn net.Conn, args ...any) {
	punch_id := args[1].(int)

	global.Peer.Status = 1

	tmp_conn, err := net.Dial("tcp", global.Serveraddr.String())
	if err != nil {
		fmt.Println("Punch connect server failed")
		return
	}

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
			}
		}()

	case 2:
		//port
	}
}
