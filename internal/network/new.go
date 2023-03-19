package network

import (
	"net"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
	"github.com/calvindc/Web3RpcHub/internal/keys"
	"github.com/calvindc/Web3RpcHub/internal/netwrap"
	"github.com/go-kit/log"
)

type Options struct {
	Logger log.Logger

	Dialer     netwrap.Dialer
	ListenAddr net.Addr

	KeyPair     *keys.KeyPair
	AppKey      []byte
	MakeHandler func(net.Conn) (cmuxrpc.Handler, error)

	ConnTracker ConnTracker

	// PreSecureWrappers are applied before the shs+boxstream wrapping takes place
	// usefull for accessing the sycall.Conn to apply control options on the socket
	BefreCryptoWrappers []netwrap.ConnWrapper

	// AfterSecureWrappers are applied afterwards, usefull to debug muxrpc content
	AfterSecureWrappers []netwrap.ConnWrapper
}
