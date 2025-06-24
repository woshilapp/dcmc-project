package event

import (
	"fmt"
	"net"

	reuse "github.com/libp2p/go-reuseport"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/client/network"
	netdata "github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/terminal"
)

func InitHostEvent() {
	protocol.RegTCPEvent(112, handleCreateRoom, protocol.IntType)
	protocol.RegTCPEvent(121, handleNewPeer, protocol.IntType)
	protocol.RegTCPEvent(122, handleNoticePunchHost, protocol.IntType, protocol.StringType)
}

func handleCreateRoom(conn net.Conn, args ...any) {
	id := args[1].(int)
	global.CurrRoom.Id = uint32(id)

	terminal.SetPrompt(global.App, global.CurrRoom.Name+">")
	global.App.Println("Created room, id:", id)
}

func handleNewPeer(conn net.Conn, args ...any) {
	global.App.Println("New peer enter room, trying to connect...")

	punch_id := args[1].(int)

	tmp_conn, err := reuse.Dial("tcp", "0.0.0.0:0", global.Serveraddr.String())
	if err != nil {
		fmt.Println("Punch connect server failed")
		return
	}

	go network.HandleConn(tmp_conn)

	peer := &global.TPeers{}
	peer.PunchID = punch_id
	peer.Conn = tmp_conn

	global.Host.Peers = append(global.Host.Peers, peer)

	str, _ := protocol.Encode(302, punch_id)
	netdata.WriteMsg(tmp_conn, []byte(str))
}

func handleNoticePunchHost(conn net.Conn, args ...any) {
	punch_id := args[1].(int)
	peer_addr := args[2].(string)
	punch_type := 0         //0:port, 1:room
	var peer *global.TPeers //for room connection

	for _, p := range global.Host.Peers {
		if punch_id == p.PunchID {
			punch_type = 1
			peer = p
			break
		}
	}

	switch punch_type {
	case 0:
		//port
	case 1:
		peer_conn, err := network.PunchPeer(peer.Conn,
			peer_addr,
			true)
		if err != nil {
			fmt.Println("Connect to a peer failed")
			return
		}

		peer.Conn = peer_conn

		//handle peer conn
		go func() {
			netdata.WriteMsg(peer_conn, []byte("300"))

			for {
				data, err := netdata.ReadMsg(peer_conn)
				if err != nil {
					fmt.Println("Disconnect from a peer")
					//clean up

					peer.Conn.Close()
					//遍历TCPConn和UDPSock关闭链接
					//删除peer, 但是要自己写
					return
				}

				fmt.Println("From Peer recv:", string(data))
			}
		}()
	}
}
