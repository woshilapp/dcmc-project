package global

import (
	"errors"
	"net"
	"sync"
)

type Room struct {
	Id          uint32
	Name        string
	Description string
	HostConn    net.Conn
	MaxPeer     int
	CurrPeer    int
	RequiredPwd bool
}

type Punch struct {
	Id       uint32
	HostConn net.Conn //tCP
	HostAddr net.Addr //UDP
	PeerConn net.Conn //TCP
	PeerAddr net.Addr //UDP
}

var room_id = 0
var punch_id = 0

var UDPconn *net.UDPConn

var (
	punchlock sync.RWMutex
	roomlock  sync.RWMutex
	punchssn  = make(map[uint32]*Punch) //punch_id:*Punch
	roomlist  = make(map[uint32]*Room)  //room_id:*Room
)

func GetRoomlist() (r []Room) {
	roomlock.RLock()
	defer roomlock.RUnlock()

	for _, v := range roomlist {
		r = append(r, *v)
	}

	return r
}

func GetRoom(id int) (Room, error) {
	roomlock.RLock()
	defer roomlock.RUnlock()

	if _, exist := roomlist[uint32(id)]; exist {
		return *roomlist[uint32(id)], nil
	}

	return Room{}, errors.New("doesn't exist")
}

func AddRoom(name string, description string, host_conn net.Conn, max_peer int, req_pwd bool) int {
	roomlock.Lock()
	defer roomlock.Unlock()

	id := uint32(room_id)

	roomlist[id] = &Room{
		Id:          id,
		Name:        name,
		Description: description,
		HostConn:    host_conn,
		MaxPeer:     max_peer,
		CurrPeer:    0,
		RequiredPwd: req_pwd,
	}

	room_id++

	return room_id - 1
}

func RemoveRoom(id int) {
	roomlock.Lock()
	defer roomlock.Unlock()

	delete(roomlist, uint32(id))
}

func UpdateRoom(id int, curr_peer int, description string) {
	roomlock.Lock()
	defer roomlock.Unlock()

	roomlist[uint32(id)].CurrPeer = curr_peer
	roomlist[uint32(id)].Description = description
}

func AllocPunchSession() int {
	punchlock.Lock()
	defer punchlock.Unlock()

	punchssn[uint32(punch_id)] = &Punch{
		Id: uint32(punch_id),
	}

	punch_id++
	return punch_id - 1
}

func GetPunchSession(id int) (Punch, error) {
	punchlock.RLock()
	defer punchlock.RUnlock()

	if _, exist := punchssn[uint32(id)]; exist {
		return *punchssn[uint32(id)], nil
	}

	return Punch{}, errors.New("doesn't exist")
}

func UpdatePunchSession(id int, host_conn net.Conn, peer_conn net.Conn) {
	punchlock.Lock()
	defer punchlock.Unlock()

	punchssn[uint32(id)].HostConn = host_conn
	punchssn[uint32(id)].PeerConn = peer_conn
}

func UpdateUDPPunchSession(id int, host_addr net.Addr, peer_addr net.Addr) {
	punchlock.Lock()
	defer punchlock.Unlock()

	punchssn[uint32(id)].HostAddr = host_addr
	punchssn[uint32(id)].PeerAddr = peer_addr
}

func DeletePunchSession(id int) {
	punchlock.Lock()
	defer punchlock.Unlock()

	delete(punchssn, uint32(id))
}
