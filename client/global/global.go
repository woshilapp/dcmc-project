package global

import (
	"net"

	"github.com/desertbit/grumble"
)

type Room struct {
	Id          uint32
	Name        string
	Description string
	MaxPeer     int
	CurrPeer    int
	RequiredPwd bool
}

type TPeer struct {
	Name     string
	HostConn net.Conn
	Status   uint8 //0: none, 1: connecting, 2: connected
	TCPConn  map[uint16]net.Conn
	UDPSock  map[uint16]*net.UDPConn
}

type THost struct {
	Passwd   string
	Status   uint8 //0: none, 1: in room
	Peers    []*TPeers
	TCPPorts []uint16
	UDPPorts []uint16
	Muted    map[string]int64 //name: unixstamp
	Banned   map[string]int64 //name: unixstamp
	IPBan    map[string]int64 //IP: unixstamp
}

type TPeers struct {
	Conn    net.Conn
	Name    string
	PunchID int
	TCPConn map[uint16]net.Conn
	UDPSock map[uint16]*net.UDPConn
	Auth    bool
}

var Roomlist []Room = []Room{}
var CurrRoom Room

var (
	Serverconn net.Conn
	Serveraddr net.Addr
	Udpsock    *net.UDPConn
	App        *grumble.App
	Role       int = 1 //1:peer, 2:host
)

var Host THost = THost{
	Passwd:   "",
	Status:   0,
	Peers:    []*TPeers{},
	TCPPorts: []uint16{},
	UDPPorts: []uint16{},
	Muted:    map[string]int64{},
	Banned:   map[string]int64{},
	IPBan:    map[string]int64{},
}

var Peer TPeer = TPeer{
	Name:    "",
	Status:  0,
	TCPConn: map[uint16]net.Conn{},
	UDPSock: map[uint16]*net.UDPConn{},
}
