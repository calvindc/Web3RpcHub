package netwrap

import (
	"net"
	"strings"
)

// Addr is a stack address, representing the protocol stack.
type AddrStack interface {
	net.Addr

	// Head returns the address of the highest element of the protocol stack.
	// This should not return a stack address
	Head() net.Addr

	// Inner returns everything below the highest element. Might be stack address itself.
	Inner() net.Addr
}

// WrapAddr wraps the inner address with head
func WrapAddr(inner, head net.Addr) net.Addr {
	if inner == nil {
		return nil
	}

	return &addr{
		inner: inner,
		head:  head,
	}
}

type addr struct {
	head  net.Addr
	inner net.Addr
}

// Network re-implement net.Addr(name of the network)
func (a *addr) Network() string {
	return a.inner.Network() + "|" + a.head.Network()
}

// String re-implement net.Addr(string form of address)
func (a *addr) String() string {
	return a.inner.String() + "|" + a.head.String()
}

func (a *addr) Inner() net.Addr {
	return a.inner
}

func (a *addr) Head() net.Addr {
	return a.head
}

// GetAddr returns, if available, the concrete address addr with netw == addr.Network() from inside the stack.
// If no such address exists, it returns nil.
func GetAddr(a net.Addr, netw string) net.Addr {
	if a == nil {
		return nil
	}

	for {
		if a.Network() == netw || strings.HasSuffix(a.Network(), "|"+netw) {
			if a, ok := a.(AddrStack); ok {
				return a.Head()
			}
			return a
		}

		b, ok := a.(AddrStack)
		if !ok {
			break
		}

		a = b.Inner()
	}

	return nil
}
