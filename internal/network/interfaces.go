package network

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/url"

	"time"

	"context"

	"io"
	"net/http"

	"github.com/calvindc/Web3RpcHub/internal/refs"
	"go.cryptoscope.co/muxrpc/v2"
	"go.cryptoscope.co/netwrap"
	"go.cryptoscope.co/secretstream"
)

// Network supplies all network related functionalitiy
type Network interface {
	Connect(ctx context.Context, addr net.Addr) error
	Serve(context.Context, ...muxrpc.HandlerWrapper) error
	GetListenAddr() net.Addr

	GetAllEndpoints() []EndpointStat
	Endpoints

	GetConnTracker() ConnTracker

	// WebsockHandler returns a "middleware" like thing that is able to upgrade a
	// websocket request to a muxrpc connection and authenticate using shs.
	// It calls the next handler if it fails to upgrade the connection to websocket.
	// However, it will error on the request and not call the passed handler
	// if the websocket upgrade is successfull.
	WebsockHandler(next http.Handler) http.Handler

	io.Closer
}

// Endpoints returns the connected endpoint for the passed feed, or false if there is none
type Endpoints interface {
	GetEndpointFor(refs.FeedRef) (muxrpc.Endpoint, bool)
}

// ConnTracker decides if connections should be established and keeps track of them
type ConnTracker interface {
	// Active returns true and since when a peer connection is active
	Active(net.Addr) (bool, time.Duration)

	// OnAccept receives a new connection as an argument.
	// If it decides to accept it, it returns true and a context that will be canceled once it should shut down
	// If it decides to deny it, it returns false (and a nil context)
	OnAccept(context.Context, net.Conn) (bool, context.Context)

	// OnClose notifies the tracker that a connection was closed
	OnClose(conn net.Conn) time.Duration

	// Count returns the number of open connections
	Count() uint

	// CloseAll closes all tracked connections
	CloseAll()
}

// EndpointStat gives some information about a connected peer
type EndpointStat struct {
	ID       *refs.FeedRef
	Addr     net.Addr
	Since    time.Duration
	Endpoint muxrpc.Endpoint
}

// HubEndpoint
type HubEndpoint struct {
	HubID                  refs.FeedRef
	ListenAddressMUXRPC    string //defaults "127.0.0.1:8008"
	HttpsDomain            string
	HttpsPort              uint //defaults 443
	UseAliasesAsSubdomains bool
	Development            bool
}

// HubAddress
// eg: net:110.41.135.27:8008~shs:8p3pnr4zESotFXWFjLPFb8Lc18DJ4NOlUoJ4iREZjag=
func (hed HubEndpoint) HubAddress() string {
	addr, err := net.ResolveTCPAddr("tcp", hed.ListenAddressMUXRPC)
	if err != nil {
		panic(err)
	}
	var hubPubKey = base64.StdEncoding.EncodeToString(hed.HubID.PubKey())
	return fmt.Sprintf("net:%s:%d~shs:%s", hed.HttpsDomain, addr.Port, hubPubKey)
}

// GetFeedRefFromAddr uses netwrap to get the secretstream address and then uses ParseFeedRef
func GetFeedRefFromAddr(addr net.Addr) (refs.FeedRef, error) {
	addr = netwrap.GetAddr(addr, secretstream.NetworkString)
	if addr == nil {
		return refs.FeedRef{}, errors.New("no shs-bs address found")
	}
	ssAddr := addr.(secretstream.Addr)
	return refs.ParseFeedRef(ssAddr.String())
}

func (hed HubEndpoint) URLForAlias(a string) string {
	var u url.URL
	if hed.Development {
		u.Path = "/alias/" + a
		u.Scheme = "http"
		u.Host = fmt.Sprintf("localhost:%d", hed.HttpsPort)
		return u.String()
	}
	u.Scheme = "https"
	if hed.UseAliasesAsSubdomains {
		u.Host = a + "." + hed.HttpsDomain
	} else {
		u.Host = hed.HttpsDomain
		u.Path = "/alias/" + a
	}
	return u.String()
}
