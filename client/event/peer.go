package event

import (
	"net"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/protocol"
)

func InitPeerEvent() {
	protocol.RegTCPEvent(101, recvRoomList, protocol.StringAnyType)
	protocol.RegTCPEvent(102, printRoomList)
}

func recvRoomList(conn net.Conn, args ...any) {
	if len(args) == 1 {
		return
	}

	for i := 1; i < len(args); i++ {
		roomdata, err := protocol.Decode(args[i].(string))
		if err != nil {
			return
		}

		global.Roomlist = append(global.Roomlist, global.Room{
			Id:          uint32(roomdata[0].(int)),
			Name:        roomdata[1].(string),
			Description: roomdata[2].(string),
			MaxPeer:     roomdata[3].(int),
			CurrPeer:    roomdata[4].(int),
			RequiredPwd: roomdata[5].(bool),
		})
	}
}

func printRoomList(conn net.Conn, args ...any) {
	for _, v := range global.Roomlist {
		global.App.Println(v.Id, v.Name, v.Description, v.MaxPeer, v.CurrPeer, v.RequiredPwd)
	}
}
