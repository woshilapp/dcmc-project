package network

import (
	"net"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
)

func ListenConn(conn net.Conn) {
	for {
		data, err := network.ReadMsg(conn)
		if err != nil {
			global.App.Println("[ERRORrt]", err)
			return
		}

		global.App.Println("[Recv Server TCP]", string(data))

		event, err := protocol.Decode(string(data))
		if err != nil {
			global.App.Println("[ERRORdc]", err)
			continue
		}

		err = protocol.VaildateTCPEvent(event...)
		if err != nil {
			global.App.Println("[BADEvent]", err, event)
			continue
		}

		protocol.ExecTCPEvent(global.Serverconn, event...)
	}
}
