package punching

import (
	"net"
	"time"

	"github.com/woshilapp/dcmc-project/client/network"
	"github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/server/global"
)

func CleanPunchSession(punch_id int) {
	counter := 0
	for {
		time.Sleep(5 * time.Second)

		if _, err := global.GetPunchSession(punch_id); err != nil {
			return
		}

		counter += 5

		if counter > 60 {
			global.DeletePunchSession(punch_id)
			return
		}
	}
}

func NoticePunching(punch_id int, host_conn net.Conn, peer_conn net.Conn) {
	timer := time.NewTimer(60 * time.Second) // 1min

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

func NoticeUDPPunching(conn *net.UDPConn, punch_id int, host_addr net.Addr, peer_addr net.Addr) {
	msgh, _ := protocol.Encode(122, punch_id, peer_addr.String())
	_, errh := conn.WriteTo([]byte(msgh), host_addr)
	if errh != nil {
		return
	}

	msgp, _ := protocol.Encode(122, punch_id, host_addr.String())
	_, errp := conn.WriteTo([]byte(msgp), peer_addr)
	if errp != nil {
		return
	}
}
