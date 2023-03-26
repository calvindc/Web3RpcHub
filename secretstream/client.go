package secretstream

import (
	"fmt"
	"net"
	"time"

	"github.com/calvindc/Web3RpcHub/netwrap"
	"github.com/calvindc/Web3RpcHub/secretstream/boxstream"
	"github.com/calvindc/Web3RpcHub/secretstream/secrethandshake"
)

// Client can dial secret-handshake server endpoints
type Client struct {
	appKey []byte
	kp     secrethandshake.EdKeyPair
}

// NewClient creates a new Client with the passed keyPair and appKey
func NewClient(kp secrethandshake.EdKeyPair, appKey []byte) (*Client, error) {
	// TODO: consistancy check
	return &Client{
		appKey: appKey,
		kp:     kp,
	}, nil
}

// ConnWrapper returns a connection wrapper for the client.
func (c *Client) ConnWrapper(pubKey []byte) netwrap.ConnWrapper {
	return func(conn net.Conn) (net.Conn, error) {
		state, err := secrethandshake.NewClientState(c.appKey, c.kp, pubKey)
		if err != nil {
			return nil, err
		}

		errconn := make(chan error)
		go func() {
			errconn <- secrethandshake.ClientShack(state, conn)
			close(errconn)
		}()

		select {
		case err := <-errconn:
			if err != nil {
				return nil, err
			}
		case <-time.After(connClientWaitTimeout):
			return nil, fmt.Errorf("[secretstream/Client.ConnWrapper]: handshake timeout")
		}

		enKey, enNonce := state.GetBoxstreamEncKeys()
		deKey, deNonce := state.GetBoxstreamDecKeys()

		boxed := &Conn{
			boxer:   boxstream.NewBoxer(conn, &enNonce, &enKey),
			unboxer: boxstream.NewUnboxer(conn, &deNonce, &deKey),
			conn:    conn,
			local:   c.kp.Public[:],
			remote:  state.Remote(),
		}

		return boxed, nil
	}
}
