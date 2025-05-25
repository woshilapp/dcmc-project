package global

import (
	"net"

	"github.com/desertbit/grumble"
)

var Serverconn net.Conn
var Serveraddr net.Addr
var Udpsock *net.UDPConn
var App *grumble.App
var Role int = 1 //1:peer, 2:host
