package global

import (
	"net"

	"github.com/desertbit/grumble"
)

var Serverconn net.Conn
var Serveraddr net.Addr
var Udpsock *net.UDPConn
var App *grumble.App
