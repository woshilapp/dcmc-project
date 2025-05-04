package punching

import (
	"time"

	"github.com/woshilapp/dcmc-project/client/network"
	"github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/server/global"
)

func PunchingThread(punch *global.Punch) {
	ttimer := time.NewTimer(60 * time.Second)

	select {
	case <-punch.NoticePunch: //notice
		to_host, _ := protocol.Encode(122, punch.Id, punch.PeerConn.RemoteAddr().String())
		to_peer, _ := protocol.Encode(122, punch.Id, punch.HostConn.RemoteAddr().String())

		network.WriteMsg(punch.HostConn, []byte(to_host))
		network.WriteMsg(punch.PeerConn, []byte(to_peer))
	case <-ttimer.C: //wait punching
		//clear
		punch.HostConn.Close()
		punch.PeerConn.Close()
		global.DeletePunchSession(int(punch.Id))
		return
	}
}
