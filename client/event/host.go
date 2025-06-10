package event

import (
	"net"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/terminal"
)

func InitHostEvent() {
	protocol.RegTCPEvent(112, handleCreateRoom, protocol.IntType)
	protocol.RegTCPEvent(121, handleNewPeer, protocol.IntType)
}

func handleCreateRoom(conn net.Conn, args ...any) {
	id := args[1].(int)
	global.CurrRoom.Id = uint32(id)

	terminal.SetPrompt(global.App, global.CurrRoom.Name+">")
	global.App.Println("Created room, id:", id)
}

func handleNewPeer(conn net.Conn, args ...any) {
	global.App.Println("New peer enter room, trying to connect...")
}
