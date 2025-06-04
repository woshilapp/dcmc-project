package event

import (
	"net"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/protocol"
)

func InitPeerEvent() {
	protocol.RegTCPEvent(110, handleEnterRoom, protocol.IntType)
	protocol.RegTCPEvent(111, handleRoomFull, protocol.IntType)
}

func handleEnterRoom(conn net.Conn, args ...any) {
	global.CurrRoom = global.Roomlist[args[1].(int)]

	global.App.Println("Enter Room Success, Trying connect to host...")
	global.App.SetPrompt(global.CurrRoom.Name + ">")
}

func handleRoomFull(conn net.Conn, args ...any) {
	global.App.Println(global.Roomlist[args[1].(int)].Name + ": Room is full")
}
