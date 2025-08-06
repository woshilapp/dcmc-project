package event

import (
	"fmt"
	"net"
	"slices"
	"strconv"
	"sync"
	"time"

	reuse "github.com/libp2p/go-reuseport"

	"github.com/woshilapp/dcmc-project/client/global"
	"github.com/woshilapp/dcmc-project/client/network"
	"github.com/woshilapp/dcmc-project/client/tunnel"
	netdata "github.com/woshilapp/dcmc-project/network"
	"github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/terminal"
)

func InitPeerEvent() {
	protocol.RegTCPEvent(110, handleEnterRoom, protocol.IntType)
	protocol.RegTCPEvent(111, handleRoomFull, protocol.IntType)
	protocol.RegTCPEvent(120, handlePunchHostID, protocol.IntType)
	protocol.RegTCPEvent(122, handleNoticePunchPeer, protocol.IntType, protocol.StringType)
	protocol.RegTCPEvent(321, handleReqPwd)
	protocol.RegTCPEvent(320, handleEnterRoomSuc)
	protocol.RegTCPEvent(322, handlePwdError)
	protocol.RegTCPEvent(323, handleNameUsed)
	protocol.RegTCPEvent(324, handleNameList, protocol.StringAnyType)
	protocol.RegTCPEvent(325, handleKicked, protocol.StringType)
	protocol.RegTCPEvent(326, handleBanned, protocol.StringType, protocol.IntType)
	protocol.RegTCPEvent(330, handlePunchPort, protocol.IntType, protocol.IntType, protocol.IntType)
	protocol.RegTCPEvent(331, handleClosePort, protocol.IntType, protocol.IntType)
	protocol.RegTCPEvent(340, handleMsg, protocol.StringType, protocol.StringType)
	protocol.RegTCPEvent(341, handleHostBroadcast, protocol.StringType)
	protocol.RegTCPEvent(342, handleMuted, protocol.StringType, protocol.IntType)
	protocol.RegUDPEvent(122, handleUDPNoticePunchPeer, protocol.IntType, protocol.StringType)
}

func handleEnterRoom(conn net.Conn, args ...any) {
	room_id := args[1].(int)

	for _, r := range global.Roomlist {
		if r.Id == uint32(room_id) {
			global.CurrRoom = r
		}
	}

	global.App.Println("Enter Room Success, Trying connect to host...")
	terminal.SetPrompt(global.App, global.CurrRoom.Name+">")
}

func handleRoomFull(conn net.Conn, args ...any) {
	room_id := args[1].(int)
	room_name := ""

	for _, r := range global.Roomlist {
		if r.Id == uint32(room_id) {
			room_name = r.Name
		}
	}

	global.App.Println(room_name+": Room is full, ID:", room_id)
}

func handlePunchHostID(conn net.Conn, args ...any) {
	punch_id := args[1].(int)

	global.Peer.Status = 1

	tmp_conn, err := reuse.Dial("tcp", "0.0.0.0:0", global.ServerAddr.String())
	if err != nil {
		global.App.Println("Punch connect server failed")
		return
	}

	go network.HandlePunchConn(tmp_conn)

	str, _ := protocol.Encode(203, punch_id)
	netdata.WriteMsg(tmp_conn, []byte(str))

	global.Peer.HostConn = tmp_conn
}

func handleNoticePunchPeer(conn net.Conn, args ...any) {
	switch global.Peer.Status {
	case 1:
		// room
		// punch_id := args[1].(int)
		host_addr := args[2].(string)

		host_conn, err := network.PunchPeer(global.Peer.HostConn, host_addr, false)
		if err != nil {
			global.App.Println("Connect to host failed")

			//clean up
			global.Peer.Status = 0
			global.Peer.HostConn.Close()
			terminal.SetPrompt(global.App, ">")
			global.CurrRoom = global.Room{}
			return
		}

		// global.Peer.Status = 2
		global.Peer.HostConn = host_conn
		//handle host conn
		go func() {
			netdata.WriteMsg(host_conn, []byte("200"))

			for {
				data, err := netdata.ReadMsg(host_conn)
				if err != nil {
					global.App.Println("Disconnect from host")

					//clean up
					global.Peer.Status = 0
					global.Peer.HostConn.Close()
					terminal.SetPrompt(global.App, ">")
					global.CurrRoom = global.Room{}

					for _, t := range global.Peer.Tunnels {
						t.Closed = true

						for _, c := range t.TCPConns {
							c.Close()
						}

						if t.UDPRemote != nil {
							t.UDPRemote.Close()
						}
					}

					return
				}

				// fmt.Println("From Host recv:", string(data))

				network.ProcTCPEvent(host_conn, data)
			}
		}()

	case 2:
		// port
		punch_id := args[1].(int)
		host_addr := args[2].(string)

		var tun *global.Tunnel
		for _, t := range global.Peer.Tunnels {
			if t.PunchID == punch_id {
				tun = t
			}
		}
		if tun == nil {
			return
		}

		// tcp
		host_conn, err := network.PunchPeer(tun.TCPRemote, host_addr, false)
		if err != nil {
			global.App.Println("Punch tcp/" + strconv.Itoa(int(tun.Port)) + " failed")
			return
		}
		tun.TCPRemote = host_conn

		//handle host tunnel
		go tunnel.HandleRemotePeer(tun, host_conn)
		go tunnel.ListenLocal(tun)
	}
}

func handleUDPNoticePunchPeer(conn *net.UDPConn, addr net.Addr, args ...any) {
	punch_id := args[1].(int)
	host_addr, _ := net.ResolveUDPAddr("udp", args[2].(string))

	var tun *global.Tunnel
	for _, t := range global.Peer.Tunnels {
		if t.PunchID == punch_id {
			tun = t
		}
	}
	if tun == nil {
		return
	}

	for i := 3; i > 0; i-- {
		tun.UDPRemote.WriteTo([]byte("200"), host_addr)
	}

	tun.UDPRemoteAddr = host_addr

	// handle host
	go tunnel.UDPListenLocal(tun)
}

func handleReqPwd(conn net.Conn, args ...any) {
	global.App.Println("Enter Password by command 'passwd'")

	// str, _ := protocol.Encode(210, pwd)
	// netdata.WriteMsg(conn, []byte(str))
}

func handleEnterRoomSuc(conn net.Conn, args ...any) {
	global.App.Println("Auth Success, Enter Name by command 'name'")
	global.Peer.Status = 2
	// str, _ := protocol.Encode(211, name)
	// netdata.WriteMsg(conn, []byte(str))
}

func handlePwdError(conn net.Conn, args ...any) {
	global.App.Println("Password Incorrect, Please enter the room again.")
}

func handleNameUsed(conn net.Conn, args ...any) {
	global.App.Println("Name already used, Please enter the room again.")
}

func handleKicked(conn net.Conn, args ...any) {
	reason := args[1].(string)
	global.App.Println("You have been kicked, Reason:", reason)
}

func handleBanned(conn net.Conn, args ...any) {
	reason := args[1].(string)
	time := time.Unix(int64(args[2].(int)), 0)
	global.App.Println("You have been banned, Reason:", reason)
	global.App.Println("Unban time:", time)
}

func handleNameList(conn net.Conn, args ...any) {
	names := ""

	for i := 0; i < len(args); i++ {
		if i == 0 {
			continue
		}

		if i == 1 {
			names = names + args[i].(string)
			continue
		}

		names = names + "," + args[i].(string)
	}

	global.App.Println("Players:", names)
}

func handleMsg(conn net.Conn, args ...any) {
	name := args[1].(string)
	msg := args[2].(string)

	global.App.Println("<" + name + ">" + msg)
}

func handleHostBroadcast(conn net.Conn, args ...any) {
	msg := args[1].(string)
	global.App.Println("Host Broadcast:", msg)
}

func handleMuted(conn net.Conn, args ...any) {
	reason := args[1].(string)
	time := time.Unix(int64(args[2].(int)), 0)
	global.App.Println("You have been muted, Reason:", reason)
	global.App.Println("Unmute time:", time)
}

func handlePunchPort(conn net.Conn, args ...any) {
	punch_id := args[1].(int)
	proto := args[2].(int)
	port := args[3].(int)
	strport := strconv.Itoa(port)

	switch proto {
	case 1:
		global.App.Println("Punch remote port tcp/" + strport + " at local 127.0.0.2:" + strport)

		p := &global.Tunnel{
			Port:     uint16(port),
			Proto:    proto,
			PunchID:  punch_id,
			Lock:     sync.RWMutex{},
			Closed:   false,
			TCPConns: map[uint32]net.Conn{},
		}
		global.Peer.Tunnels = append(global.Peer.Tunnels, p)

		tmp_conn, err := reuse.Dial("tcp", "0.0.0.0:0", global.ServerAddr.String())
		if err != nil {
			global.App.Println("Punch connect server failed")
			return
		}

		p.TCPRemote = tmp_conn

		go network.HandlePunchConn(tmp_conn)

		str, _ := protocol.Encode(203, punch_id)
		netdata.WriteMsg(tmp_conn, []byte(str))
	case 2:
		global.App.Println("Punch remote port udp/" + strport + " at local 127.0.0.2:" + strport)

		p := &global.Tunnel{
			Port:     uint16(port),
			Proto:    proto,
			PunchID:  punch_id,
			Lock:     sync.RWMutex{},
			Closed:   false,
			UDPAddrs: map[uint32]net.Addr{},
		}
		global.Peer.Tunnels = append(global.Peer.Tunnels, p)

		localaddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:0")
		sock, err := net.ListenUDP("udp", localaddr)
		if err != nil {
			fmt.Println("UDP sock dial failed", err)
			return
		}

		p.UDPRemote = sock

		go network.HandleUDP(sock)

		str, _ := protocol.Encode(203, punch_id)
		sock.WriteTo([]byte(str), global.ServerUDPAddr)
	}
}

func handleClosePort(conn net.Conn, args ...any) {
	proto := args[2].(int)
	port := args[3].(int)
	strport := strconv.Itoa(port)
	var tun *global.Tunnel
	var tunidx int

	for i, t := range global.Peer.Tunnels {
		if t.Port == uint16(port) && t.Proto == proto {
			tun = t
			tunidx = i
		}
	}
	if tun == nil {
		return
	}

	switch proto {
	case 1:
		global.App.Println("Host Closed Port tcp/" + strport)

		tun.Closed = true
		tun.TCPRemote.Close()
		for _, c := range tun.TCPConns {
			c.Close()
		}

		global.Peer.Tunnels = slices.Delete(global.Peer.Tunnels, tunidx, tunidx+1)
	case 2:
		global.App.Println("Host Closed Port udp/" + strport)

		tun.Closed = true
		tun.UDPRemote.Close()

		global.Peer.Tunnels = slices.Delete(global.Peer.Tunnels, tunidx, tunidx+1)
	}
}
