package event

import (
	"net"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/protocol"
)

func InitHostEvent() {
	protocol.RegTCPEvent(112, handleCreateRoom, protocol.IntType)
	protocol.RegTCPEvent(121, handleNewPeer, protocol.IntType)
}

func handleCreateRoom(conn net.Conn, args ...any) {
	global.App.Println("Created room, id:", args[1].(int))
}

func handleNewPeer(conn net.Conn, args ...any) {
	global.App.Println("New peer enter room, trying to connect...")
}
