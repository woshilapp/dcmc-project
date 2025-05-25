package event

import (
	"net"

	"github.com/woshilapp/dcmc-project/client/network"
	"github.com/woshilapp/dcmc-project/protocol"
)

func InitPeerEvent() {
	protocol.RegTCPEvent(100, responeServerHello)
}

func responeServerHello(conn net.Conn, args ...any) {
	str, _ := protocol.Encode(200)

	network.WriteMsg(conn, []byte(str))
}
