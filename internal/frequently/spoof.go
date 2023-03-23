package frequently

import (
	"errors"
	"net"

	"github.com/calvindc/Web3RpcHub/netwrap"
	"github.com/calvindc/Web3RpcHub/secretstream"
)

type SpoofedConn struct {
	net.Conn
	spoofedRemote net.Addr
}

func (sc SpoofedConn) RemoteAddr() net.Addr {
	return sc.spoofedRemote
}

// SpoofRemoteAddress wraps the connection with the passed reference
func SpoofRemoteAddress(pubKey []byte) netwrap.ConnWrapper {
	if len(pubKey) != 32 {
		return func(_ net.Conn) (net.Conn, error) {
			return nil, errors.New("invalid public key length")
		}
	}
	var spoofedAddr secretstream.Addr
	spoofedAddr.PubKey = pubKey
	return func(c net.Conn) (net.Conn, error) {
		sc := SpoofedConn{
			Conn:          c,
			spoofedRemote: netwrap.WrapAddr(c.RemoteAddr(), spoofedAddr),
		}
		return sc, nil
	}
}
