package event

import (
	"fmt"
	"net"
	"slices"

	reuse "github.com/libp2p/go-reuseport"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/client/network"
	"github.com/woshilapp/dcmc-project/client/tunnel"
	netdata "github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/terminal"
)

func InitHostEvent() {
	protocol.RegTCPEvent(112, handleCreateRoom, protocol.IntType)
	protocol.RegTCPEvent(120, handleNewPunchID, protocol.IntType)
	protocol.RegTCPEvent(121, handleNewPeer, protocol.IntType)
	protocol.RegTCPEvent(122, handleNoticePunchHost, protocol.IntType, protocol.StringType)
	protocol.RegTCPEvent(200, handlePeerHello)
	protocol.RegTCPEvent(210, handlePeerPasswd, protocol.StringType)
	protocol.RegTCPEvent(211, handlePeerName, protocol.StringType)
	protocol.RegTCPEvent(212, handlePeerReqList)
	protocol.RegTCPEvent(230, handlePeerMsg, protocol.StringType)
	protocol.RegUDPEvent(122, handleUDPNoticePunchHost, protocol.IntType, protocol.StringType)
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

	tmp_conn, err := reuse.Dial("tcp", "0.0.0.0:0", global.ServerAddr.String())
	if err != nil {
		fmt.Println("Punch connect server failed")
		return
	}

	go network.HandlePunchConn(tmp_conn)

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
	var peer_ind int        //for clean thread

	for i, p := range global.Host.Peers {
		if punch_id == p.PunchID {
			punch_type = 1
			peer = p
			peer_ind = i
			break
		}
	}

	switch punch_type {
	case 0:
		tun := global.Host.PIDtun[punch_id]
		peer_conn, err := network.PunchPeer(tun.TCPRemote,
			peer_addr,
			true)
		if err != nil {
			fmt.Println("Connect to a peer failedp")
			return
		}

		tun.TCPRemote = peer_conn
		delete(global.Host.PIDtun, punch_id)

		go tunnel.HandleRemoteHost(tun, peer_conn)
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

					for _, t := range global.Host.Peers[peer_ind].Tunnels {
						switch t.Proto {
						case 1:
							t.TCPRemote.Close()
						case 2:
							t.UDPRemote.Close()
						}

						for _, c := range t.TCPConns {
							c.Close()
						}
						for _, c := range t.UDPConns {
							c.Close()
						}
					}

					global.Host.Peers = slices.Delete(global.Host.Peers, peer_ind, peer_ind+1)

					if peer.Auth {
						global.CurrRoom.CurrPeer--

						str, _ := protocol.Encode(312, global.CurrRoom.Id,
							global.CurrRoom.CurrPeer,
							global.CurrRoom.Description,
							global.CurrRoom.RequiredPwd,
						)
						netdata.WriteMsg(global.ServerConn, []byte(str))
					}

					return
				}

				// fmt.Println("From Peer recv:", string(data))

				network.ProcTCPEvent(peer_conn, data)
			}
		}()
	}
}

func handleUDPNoticePunchHost(conn *net.UDPConn, addr net.Addr, args ...any) {
	punch_id := args[1].(int)
	peer_addr, _ := net.ResolveUDPAddr("udp", args[2].(string))

	tun := global.Host.PIDtun[punch_id]

	for i := 3; i > 0; i-- {
		tun.UDPRemote.WriteTo([]byte("300"), peer_addr)
	}

	tun.UDPRemoteAddr = peer_addr

	delete(global.Host.PIDtun, punch_id)

	go tunnel.HandleUDPRemoteHost(tun, conn)
}

func handlePeerHello(conn net.Conn, args ...any) {
	if global.CurrRoom.RequiredPwd {
		str, _ := protocol.Encode(321)
		netdata.WriteMsg(conn, []byte(str))
		return
	}

	var peer *global.TPeers
	for _, p := range global.Host.Peers {
		if p.Conn == conn {
			peer = p
		}
	}

	peer.Auth = true

	str, _ := protocol.Encode(320)
	netdata.WriteMsg(conn, []byte(str))

	//update status
	global.CurrRoom.CurrPeer++
	str1, _ := protocol.Encode(312, global.CurrRoom.Id,
		global.CurrRoom.CurrPeer,
		global.CurrRoom.Description,
		global.CurrRoom.RequiredPwd,
	)
	netdata.WriteMsg(global.ServerConn, []byte(str1))

	tunnel.HandleNewPeer(peer)
}

func handlePeerPasswd(conn net.Conn, args ...any) {
	if !global.CurrRoom.RequiredPwd {
		return
	}

	pwd := args[1].(string)

	if global.Host.Passwd != pwd {
		str, _ := protocol.Encode(322)
		netdata.WriteMsg(conn, []byte(str))
		return
	}

	var peer *global.TPeers
	for _, p := range global.Host.Peers {
		if p.Conn == conn {
			peer = p
		}
	}

	peer.Auth = true

	str, _ := protocol.Encode(320)
	netdata.WriteMsg(conn, []byte(str))

	//update status
	global.CurrRoom.CurrPeer++
	str1, _ := protocol.Encode(312, global.CurrRoom.Id,
		global.CurrRoom.CurrPeer,
		global.CurrRoom.Description,
		global.CurrRoom.RequiredPwd,
	)
	netdata.WriteMsg(global.ServerConn, []byte(str1))

	tunnel.HandleNewPeer(peer)
}

func handlePeerName(conn net.Conn, args ...any) {
	var peer *global.TPeers

	for _, p := range global.Host.Peers {
		if p.Conn == conn {
			peer = p
		}
	}

	if !peer.Auth || peer.Name != "" {
		return
	}

	name := args[1].(string)
	for _, p := range global.Host.Peers {
		if p.Name == name || name == "Host" {
			str, _ := protocol.Encode(323)
			netdata.WriteMsg(conn, []byte(str))
			return
		}
	}

	peer.Name = name
}

func handlePeerReqList(conn net.Conn, args ...any) {
	data := []any{324}
	for _, p := range global.Host.Peers {
		if p.Name != "" {
			data = append(data, p.Name)
		}
	}

	str, _ := protocol.Encode(data...)
	netdata.WriteMsg(conn, []byte(str))
}

func handlePeerMsg(conn net.Conn, args ...any) {
	msg := args[1].(string)
	var peer *global.TPeers

	for _, p := range global.Host.Peers {
		if p.Conn == conn {
			peer = p
		}
	}

	if peer.Name == "" {
		return
	}

	str, _ := protocol.Encode(340, peer.Name, msg)

	for _, p := range global.Host.Peers {
		if p.Name != "" {
			netdata.WriteMsg(p.Conn, []byte(str))
		}
	}

	global.App.Println("<" + peer.Name + ">" + msg)
}

func handleNewPunchID(conn net.Conn, args ...any) {
	punch_id := args[1].(int)

	global.Host.PunchIDs <- punch_id
}
