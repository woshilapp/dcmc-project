package punching

import (
	"net"
	"time"

	"github.com/woshilapp/dcmc-project/client/network"
	"github.com/woshilapp/dcmc-project/protocol"
)

func NoticePunching(punch_id int, host_conn net.Conn, peer_conn net.Conn) {
	timer := time.NewTimer(60 * time.Second)

	msgh, _ := protocol.Encode(122, punch_id, peer_conn.RemoteAddr().String())
	errh := network.WriteMsg(host_conn, []byte(msgh))
	if errh != nil {
		return
	}

	msgp, _ := protocol.Encode(122, punch_id, host_conn.RemoteAddr().String())
	errp := network.WriteMsg(peer_conn, []byte(msgp))
	if errp != nil {
		return
	}

	<-timer.C // wait punching

	host_conn.Close()
	peer_conn.Close()
}
