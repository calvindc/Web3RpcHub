package netwrap

import (
	"net"

	"github.com/pkg/errors"
)

// ConnWrapper 通讯连接管理器,用于封装网络连接
type ConnWrapper func(net.Conn) (net.Conn, error)

// Dialer 同Dial(), 用于替代dialers (socks proxy)
type Dialer func(net.Addr, ...ConnWrapper) (net.Conn, error)

// Dial首先打开与提供的addr的网络连接，然后应用所有传递的连接wrappers
func Dial(addr net.Addr, wrappers ...ConnWrapper) (net.Conn, error) {
	origConn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		return nil, errors.Wrap(err, "[netwrap/listen]: error dialing")
	}
	conn := origConn
	for _, cw := range wrappers {
		conn, err = cw(conn)
		if err != nil {
			origConn.Close()
			return nil, errors.Wrap(err, "[netwrap/listen]: error wrapping connection")
		}
	}

	return conn, nil
}
