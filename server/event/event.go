package event

import (
	"net"

	"github.com/woshilapp/dcmc-project/client/network"
	proto "github.com/woshilapp/dcmc-project/protocol"
	"github.com/woshilapp/dcmc-project/server/global"
	"github.com/woshilapp/dcmc-project/server/punching"
)

func init() {
	//hello
	proto.RegEvent(200, handlePeerHello)
	proto.RegEvent(300, handleHostHello)

	//room
	proto.RegEvent(201, handleEnterRoom, proto.IntType)
	proto.RegEvent(202, handleReqRoomlist)
	proto.RegEvent(301, handlePunchPort)
	proto.RegEvent(310, handleRegRoom, proto.StringType, proto.IntType, proto.StringType, proto.BoolType)
	proto.RegEvent(311, handleDeleteRoom, proto.IntType)
	proto.RegEvent(312, handleUpdateRoom, proto.IntType, proto.IntType, proto.StringType)

	//punching
	proto.RegEvent(203, handleReqPunchClient, proto.IntType)
	proto.RegEvent(302, handleReqPunchHost, proto.IntType)
}

// fmt 200 (int)
func handlePeerHello(conn net.Conn, a ...any) {
	str, _ := proto.Encode(100)

	network.WriteMsg(conn, []byte(str))
}

// fmt 300 (int)
func handleHostHello(conn net.Conn, a ...any) {
	str, _ := proto.Encode(100)

	network.WriteMsg(conn, []byte(str))
}

// fmt 201|room_id (int, int)
func handleEnterRoom(conn net.Conn, a ...any) {
	//TODO: tell host new peer and punch_id
	//		tell peer success or not and punch_id
	room, err := global.GetRoom(a[1].(int))
	if err != nil {
		return
	}

	if room.CurrPeer >= room.MaxPeer { // full
		str, _ := proto.Encode(111, room.Id)

		network.WriteMsg(conn, []byte(str))
		return
	}

	// alloc punch session
	punch_id := global.AllocPunchSession()
	go punching.CleanPunchSession(punch_id)

	// tell client
	str, _ := proto.Encode(110, room.Id)
	str1, _ := proto.Encode(120, punch_id)

	network.WriteMsg(conn, []byte(str))
	network.WriteMsg(conn, []byte(str1))

	// tell host
	str, _ = proto.Encode(121, punch_id)
	network.WriteMsg(room.HostConn, []byte(str))
}

// fmt 202 (int)
func handleReqRoomlist(conn net.Conn, a ...any) {
	data := []any{101}

	for _, r := range global.GetRoomlist() {
		room, _ := proto.Encode(int(r.Id), r.Name, r.Description, r.MaxPeer, r.CurrPeer, r.RequiredPwd)
		data = append(data, room)

		if len(data) > 10 {
			str, _ := proto.Encode(data...)
			network.WriteMsg(conn, []byte(str))

			data = []any{101}
		}
	}

	str, _ := proto.Encode(data...)
	network.WriteMsg(conn, []byte(str))

	str, _ = proto.Encode(102)
	network.WriteMsg(conn, []byte(str))
}

// fmt 203|punch_id (int, int)
func handleReqPunchClient(conn net.Conn, a ...any) {
	punch, err := global.GetPunchSession(a[1].(int))
	if err != nil {
		return
	}

	host_conn := punch.HostConn
	global.UpdatePunchSession(int(punch.Id), host_conn, conn)

	if host_conn != nil {
		go punching.NoticePunching(int(punch.Id), conn, host_conn)
		global.DeletePunchSession(int(punch.Id))
	}
}

// fmt 301 (int)
func handlePunchPort(conn net.Conn, a ...any) {
	punch_id := global.AllocPunchSession()
	go punching.CleanPunchSession(punch_id) // haha we still need it

	str, _ := proto.Encode(120, punch_id)
	network.WriteMsg(conn, []byte(str))
}

// fmt 302|punch_id (int, int)
func handleReqPunchHost(conn net.Conn, a ...any) {
	punch, err := global.GetPunchSession(a[1].(int))
	if err != nil {
		return
	}

	peer_conn := punch.PeerConn
	global.UpdatePunchSession(int(punch.Id), conn, peer_conn)

	if peer_conn != nil {
		go punching.NoticePunching(int(punch.Id), conn, peer_conn)
		global.DeletePunchSession(int(punch.Id))
	}
}

// fmt 310|room_name|max_peer|descrpition|need_password (int, string, int, string, bool)
func handleRegRoom(conn net.Conn, a ...any) {
	for _, r := range global.GetRoomlist() { // A host a room
		if r.HostConn.RemoteAddr().String() == conn.RemoteAddr().String() {
			return
		}
	}

	room_id := global.AddRoom(a[1].(string), a[3].(string), conn, a[2].(int), a[4].(bool))
	str, _ := proto.Encode(112, room_id)

	network.WriteMsg(conn, []byte(str))
}

// fmt 311|room_id (int, int)
func handleDeleteRoom(conn net.Conn, a ...any) {
	room, err := global.GetRoom(a[1].(int))
	if err != nil {
		return
	}

	if conn.RemoteAddr().String() != room.HostConn.RemoteAddr().String() { // Avoid bad things
		return
	}

	global.RemoveRoom(a[1].(int))
}

// fmt 312|room_id|curr_peer|descrpition (int, int, int, string)
func handleUpdateRoom(conn net.Conn, a ...any) {
	room, err := global.GetRoom(a[1].(int))
	if err != nil {
		return
	}

	if conn.RemoteAddr().String() != room.HostConn.RemoteAddr().String() { // Avoid bad things
		return
	}

	if a[2].(int) > room.MaxPeer { // impossible
		return
	}

	global.UpdateRoom(a[1].(int), a[2].(int), a[3].(string))
}
