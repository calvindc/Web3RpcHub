package netwrap

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testAddr struct {
	net, str string
}

func (a testAddr) Network() string {
	return a.net
}
func (a testAddr) String() string {
	return a.str
}

func TestWrapAddr(t *testing.T) {
	a := assert.New(t)

	tcpAddr := &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8008,
	}
	head := testAddr{
		net: "foo",
		str: "some-info",
	}
	wrappedAddr := WrapAddr(tcpAddr, head)
	unwrappedFoo := GetAddr(wrappedAddr, "foo")
	unwrappedTcp := GetAddr(wrappedAddr, "tcp")

	a.Equal(wrappedAddr.String(), "127.0.0.1:8008|some-info")
	a.Equal(unwrappedFoo.String(), "some-info")
	a.Equal(unwrappedTcp.String(), "127.0.0.1:8008")
}
