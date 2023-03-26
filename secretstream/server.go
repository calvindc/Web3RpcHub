package secretstream

import (
	"fmt"
	"net"
	"time"

	"github.com/calvindc/Web3RpcHub/netwrap"
	"github.com/calvindc/Web3RpcHub/secretstream/boxstream"
	"github.com/calvindc/Web3RpcHub/secretstream/secrethandshake"
)

// Server can create net.Listeners
type Server struct {
	keyPair secrethandshake.EdKeyPair
	appKey  []byte
}

// NewServer returns a Server which uses the passed keyPair and appKey
func NewServer(keyPair secrethandshake.EdKeyPair, appKey []byte) (*Server, error) {
	return &Server{keyPair: keyPair, appKey: appKey}, nil
}

// ListenerWrapper returns a listener wrapper.
func (s *Server) ListenerWrapper() netwrap.ListenerWrapper {
	return netwrap.NewListenerWrapper(s.Addr(), s.ConnWrapper())
}

// ConnWrapper returns a connection wrapper.
func (s *Server) ConnWrapper() netwrap.ConnWrapper {
	return func(conn net.Conn) (net.Conn, error) {
		state, err := secrethandshake.NewServerState(s.appKey, s.keyPair)
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		errconn := make(chan error)
		go func() {
			errconn <- secrethandshake.ServerShake(state, conn)
			close(errconn)
		}()

		select {
		case err := <-errconn:
			if err != nil {
				return nil, err
			}
		case <-time.After(connServerWaitTimeout):
			return nil, fmt.Errorf("[secretstream/Server.ConnWrapper]: handshake timeout")
		}

		enKey, enNonce := state.GetBoxstreamEncKeys()
		deKey, deNonce := state.GetBoxstreamDecKeys()

		remote := state.Remote()
		boxed := &Conn{
			boxer:   boxstream.NewBoxer(conn, &enNonce, &enKey),
			unboxer: boxstream.NewUnboxer(conn, &deNonce, &deKey),
			conn:    conn,
			local:   s.keyPair.Public[:],
			remote:  remote[:],
		}

		return boxed, nil
	}
}

// Addr returns the shs-bs address of the server.
func (s *Server) Addr() net.Addr {
	return Addr{s.keyPair.Public[:]}
}
