package global

import (
	"net"
	"sync"

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

type Tunnel struct {
	Proto         int //1:tcp, 2:udp
	Port          uint16
	PunchID       int
	ID            uint32
	TCPConns      map[uint32]net.Conn //id:conn
	UDPConns      map[uint32]*net.UDPConn
	UDPAddrs      map[uint32]net.Addr
	Lock          sync.RWMutex
	Closed        bool //true = closed
	TCPRemote     net.Conn
	UDPRemote     *net.UDPConn
	UDPRemoteAddr net.Addr
}

type TPeer struct {
	Name     string
	HostConn net.Conn
	Status   uint8 //0: none, 1: connecting, 2: connected
	Tunnels  []*Tunnel
}

type THost struct {
	Passwd   string
	Status   uint8 //0: none, 1: in room
	Peers    []*TPeers
	TCPPorts []uint16
	UDPPorts []uint16
	PunchIDs chan int
	PIDtun   map[int]*Tunnel
	Muted    map[string]int64 //name: unixstamp
	Banned   map[string]int64 //name: unixstamp
	IPBan    map[string]int64 //IP: unixstamp
}

type TPeers struct {
	Conn    net.Conn
	Name    string
	PunchID int
	Tunnels []*Tunnel
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
	PunchIDs: make(chan int, 10),
	PIDtun:   map[int]*Tunnel{},
	Muted:    map[string]int64{},
	Banned:   map[string]int64{},
	IPBan:    map[string]int64{},
}

var Peer TPeer = TPeer{
	Name:    "",
	Status:  0,
	Tunnels: []*Tunnel{},
}
