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

var Roomlist []Room = []Room{}

var Serverconn net.Conn
var Serveraddr net.Addr
var Udpsock *net.UDPConn
var App *grumble.App
var Role int = 1 //1:peer, 2:host
