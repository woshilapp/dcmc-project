package punching

import (
	"errors"
	"fmt"
	"net"

	"github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/server/global"
)

func vaildateUDPEvent(a ...any) error {
	if len(a) < 1 {
		return errors.New("null event")
	}

	if t := fmt.Sprintf("%T", a[0]); t != "int" {
		return errors.New("illegal event")
	}

	e := a[0].(int)

	if e != 203 { // peer punch req
		return errors.New("unsupported event")
	} else if e != 302 { // host punch req
		return errors.New("unsupported event")
	}

	return nil
}

func HandleUDPPunch(a []any, addr net.Addr) {
	err := vaildateUDPEvent(a...)
	if err != nil {
		return
	}

	if a[0] == 302 {
		go handleUDPReqPunchHost(a, addr)
	} else {
		go handleUDPReqPunchPeer(a, addr)
	}
}

func handleUDPReqPunchPeer(a []any, addr net.Addr) {
	punch, err := global.GetPunchSession(a[1].(int))
	if err != nil {
		return
	}

	host_addr := punch.HostAddr

	if host_addr != nil {
		go noticeUDPPunching(int(punch.Id), host_addr, addr)
		global.DeletePunchSession(int(punch.Id))
	}
}

func handleUDPReqPunchHost(a []any, addr net.Addr) {
	punch, err := global.GetPunchSession(a[1].(int))
	if err != nil {
		return
	}

	peer_addr := punch.PeerAddr
	global.UpdateUDPPunchSession(int(punch.Id), addr, peer_addr)

	if peer_addr != nil {
		go noticeUDPPunching(int(punch.Id), addr, peer_addr)
		global.DeletePunchSession(int(punch.Id))
	}
}

func noticeUDPPunching(punch_id int, host_addr net.Addr, peer_addr net.Addr) {
	msgh, _ := protocol.Encode(122, punch_id, peer_addr.String())
	_, errh := global.UDPconn.WriteTo([]byte(msgh), host_addr)
	if errh != nil {
		return
	}

	msgp, _ := protocol.Encode(122, punch_id, host_addr.String())
	_, errp := global.UDPconn.WriteTo([]byte(msgp), peer_addr)
	if errp != nil {
		return
	}
}
