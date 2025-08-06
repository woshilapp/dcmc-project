package tunnel

import (
	"net"

	"github.com/woshilapp/dcmc-project/client/global"
)

// [Proto][Method]ConnH

func TAddConnH(t *global.Tunnel, id uint32, conn net.Conn) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	t.TCPConns[id] = conn
}

func TDelConnH(t *global.Tunnel, id uint32) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	if _, exist := t.TCPConns[id]; !exist {
		return
	}

	delete(t.TCPConns, id)
}

func TGetConnH(t *global.Tunnel, id uint32) net.Conn {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	if _, exist := t.TCPConns[id]; !exist {
		return nil
	}

	return t.TCPConns[id]
}

func UAddConnH(t *global.Tunnel, id uint32, conn *net.UDPConn) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	t.UDPConns[id] = conn
}

func UDelConnH(t *global.Tunnel, id uint32) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	if _, exist := t.UDPConns[id]; !exist {
		return
	}

	delete(t.UDPConns, id)
}

func UGetConnH(t *global.Tunnel, id uint32) *net.UDPConn {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	if _, exist := t.UDPConns[id]; !exist {
		return nil
	}

	return t.UDPConns[id]
}

// [Proto][Method]ConnP

func TAddConnP(t *global.Tunnel, conn net.Conn) uint32 {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	t.ID++
	id := t.ID - 1

	t.TCPConns[id] = conn

	return id
}

func TDelConnP(t *global.Tunnel, id uint32) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	if _, exist := t.TCPConns[id]; !exist {
		return
	}

	delete(t.TCPConns, id)
}

func TGetConnP(t *global.Tunnel, id uint32) net.Conn {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	if _, exist := t.TCPConns[id]; !exist {
		return nil
	}

	return t.TCPConns[id]
}

func UAddAddrP(t *global.Tunnel, addr net.Addr) uint32 {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	t.ID++
	id := t.ID - 1

	t.UDPAddrs[id] = addr

	return id
}

func UDelAddrP(t *global.Tunnel, id uint32) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	if _, exist := t.UDPAddrs[id]; !exist {
		return
	}

	delete(t.UDPAddrs, id)
}

func UGetAddrP(t *global.Tunnel, id uint32) net.Addr {
	t.Lock.RLock()
	defer t.Lock.RUnlock()

	if _, exist := t.UDPAddrs[id]; !exist {
		return nil
	}

	return t.UDPAddrs[id]
}

func UGetIDP(t *global.Tunnel, addr net.Addr) int {
	for k, v := range t.UDPAddrs {
		if v == addr {
			return int(k)
		}
	}

	return -1
}
