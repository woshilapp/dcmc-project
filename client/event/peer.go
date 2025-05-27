package event

import (
	"net"

	netdata "github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
)

func InitPeerEvent() {
	protocol.RegTCPEvent(100, responeServerHello)
}

func responeServerHello(conn net.Conn, args ...any) {
	str, _ := protocol.Encode(200)

	netdata.WriteMsg(conn, []byte(str))
}
