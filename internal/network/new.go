package network

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
	"github.com/calvindc/Web3RpcHub/internal/keys"
	"github.com/calvindc/Web3RpcHub/netwrap"
	"github.com/calvindc/Web3RpcHub/refs"
	"github.com/calvindc/Web3RpcHub/secretstream"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Options struct {
	Logger log.Logger

	Dialer     netwrap.Dialer
	ListenAddr net.Addr

	KeyPair     *keys.KeyPair
	AppKey      []byte
	MakeHandler func(net.Conn) (cmuxrpc.Handler, error)

	ConnTracker ConnTracker

	//在shs+boxstream包装发生之前应用PreSecureWrappers，用于访问sycall.Conne以在套接字上应用控制选项
	BefreCryptoWrappers []netwrap.ConnWrapper

	// AfterSecureWrappers are applied afterwards, usefull to debug muxrpc content
	AfterSecureWrappers []netwrap.ConnWrapper
}

type node struct {
	opts         Options
	log          log.Logger
	listening    chan struct{}
	listenerLock sync.Mutex
	lisClose     sync.Once
	lis          net.Listener

	dialer       netwrap.Dialer
	secretServer *secretstream.Server
	secretClient *secretstream.Client

	connTracker ConnTracker

	remotesLock sync.Mutex
	remotes     map[string]cmuxrpc.Endpoint

	beforeCryptoConnWrappers []netwrap.ConnWrapper
	afterSecureConnWrappers  []netwrap.ConnWrapper

	//web socket
	httpLis     net.Listener
	httpHandler http.Handler
}

func New(opts Options) (Network, error) {
	n := &node{
		opts:                     opts,
		log:                      opts.Logger,
		listening:                make(chan struct{}),
		remotes:                  make(map[string]cmuxrpc.Endpoint),
		beforeCryptoConnWrappers: opts.BefreCryptoWrappers,
		afterSecureConnWrappers:  opts.AfterSecureWrappers,
	}

	if opts.ConnTracker == nil {
		opts.ConnTracker = NewConnTracker()
	}
	n.connTracker = opts.ConnTracker

	if opts.Dialer != nil {
		n.dialer = opts.Dialer
	} else {
		n.dialer = netwrap.Dial
	}

	var err error
	n.secretClient, err = secretstream.NewClient(opts.KeyPair.Pair, opts.AppKey)
	if err != nil {
		return nil, fmt.Errorf("[network/New]: error creating secretstream.Client: %w", err)
	}

	n.secretServer, err = secretstream.NewServer(opts.KeyPair.Pair, opts.AppKey)
	if err != nil {
		return nil, fmt.Errorf("[network/New]: error creating secretstream.Server: %w", err)
	}

	return n, nil
}

//GetEndpointFor feed查找
func (n *node) GetEndpointFor(ref refs.FeedRef) (cmuxrpc.Endpoint, bool) {
	n.remotesLock.Lock()
	defer n.remotesLock.Unlock()

	edp, has := n.remotes[ref.String()]
	return edp, has
}

// addRemote should merge with conntracker
func (n *node) addRemote(edp cmuxrpc.Endpoint) {
	n.remotesLock.Lock()
	defer n.remotesLock.Unlock()
	r, err := GetFeedRefFromAddr(edp.Remote())
	if err != nil {
		panic(err)
	}
	// ref := r.String()
	// if oldEdp, has := n.remotes[ref]; has {
	// n.log.Log("remotes", "previous active", "ref", ref)
	// c := client.FromEndpoint(oldEdp)
	// _, err := c.Whoami()
	// if err == nil {
	// 	// old one still works
	// 	return
	// }
	// }
	// replace with new
	n.remotes[r.String()] = edp
}

// removeRemote should merge with conntracker
func (n *node) removeRemote(edp cmuxrpc.Endpoint) {
	n.remotesLock.Lock()
	defer n.remotesLock.Unlock()
	r, err := GetFeedRefFromAddr(edp.Remote())
	if err != nil {
		panic(err)
	}
	delete(n.remotes, r.String())
}

func (n *node) handleConnection(ctx context.Context, origConn net.Conn, isServer bool, hws ...cmuxrpc.HandlerWrapper) {
	// TODO: overhaul events and logging levels
	conn, err := n.applyConnWrappers(origConn)
	if err != nil {
		origConn.Close()
		level.Error(n.log).Log("msg", "node/Serve: failed to wrap connection", "err", err)
		return
	}

	// TODO: obfuscate the remote address by nulling bytes in it before printing ip and pubkey in full
	remoteRef, err := GetFeedRefFromAddr(conn.RemoteAddr())
	if err != nil {
		conn.Close()
		level.Error(n.log).Log("conn", "not shs authorized", "err", err)
		return
	}
	rLogger := log.With(n.log, "peer", remoteRef.ShortSigil())

	ok, ctx := n.connTracker.OnAccept(ctx, conn)
	if !ok {
		err := conn.Close()
		level.Debug(rLogger).Log("conn", "ignored", "err", err)
		return
	}

	defer func() {
		n.connTracker.OnClose(conn)
		conn.Close()
		origConn.Close()
	}()

	h, err := n.opts.MakeHandler(conn)
	if err != nil {
		level.Debug(rLogger).Log("conn", "mkHandler", "err", err)
		return
	}

	for _, hw := range hws {
		h = hw(h)
	}

	pkr := cmuxrpc.NewPacker(conn)
	filtered := level.NewFilter(n.log, level.AllowInfo())
	edp := cmuxrpc.Handle(pkr, h,
		cmuxrpc.WithContext(ctx),
		cmuxrpc.WithLogger(filtered),
		// _isServer_ defines _are we a server_.
		// the cmuxrpc option asks are we _talking_ to a server > inverted
		cmuxrpc.WithIsServer(!isServer),
	)
	n.addRemote(edp)

	srv := edp.(cmuxrpc.Server)

	err = srv.Serve()
	// level.Warn(n.log).Log("conn", "serve-return", "err", err)
	if err != nil && !isConnBrokenErr(err) && !errors.Is(err, context.Canceled) {
		level.Debug(rLogger).Log("conn", "serve exited", "err", err)
	}
	n.removeRemote(edp)

	// panic("serve exited")
	err = edp.Terminate()
	// level.Error(n.log).Log("conn", "serve-defer-terminate", "err", err)
}

// -------------------------------------------------------------------------

func (n *node) Connect(ctx context.Context, addr net.Addr) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	shsAddr := netwrap.GetAddr(addr, "shs-bs")
	if shsAddr == nil {
		return errors.New("[network/Connect]: expected an address containing an shs-bs addr")
	}

	var pubKey = make(ed25519.PublicKey, ed25519.PublicKeySize)
	if shsAddr, ok := shsAddr.(secretstream.Addr); ok {
		copy(pubKey[:], shsAddr.PubKey)
	} else {
		return errors.New("[network/Connect]: expected shs-bs address to be of type secretstream.Addr")
	}

	conn, err := n.dialer(netwrap.GetAddr(addr, "tcp"), append(n.beforeCryptoConnWrappers,
		n.secretClient.ConnWrapper(pubKey))...)
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		return fmt.Errorf("[network/Connect]: error dialing: %w", err)
	}

	go func(c net.Conn) {
		n.handleConnection(ctx, c, false)
	}(conn)
	return nil
}

// Serve starts the network listener and configured resources like local discovery.
// Canceling the passed context makes the function return. Defers take care of stopping these resources.
func (n *node) Serve(ctx context.Context, wrappers ...cmuxrpc.HandlerWrapper) error {
	// TODO: make multiple listeners (localhost:8008 should not restrict or kill connections)
	lisWrap := netwrap.NewListenerWrapper(n.secretServer.Addr(), append(n.opts.BefreCryptoWrappers, n.secretServer.ConnWrapper())...)
	var err error

	n.listenerLock.Lock()
	n.lis, err = netwrap.Listen(n.opts.ListenAddr, lisWrap)
	if err != nil {
		n.listenerLock.Unlock()
		return fmt.Errorf("[network/Serve]: error creating listener: %w", err)
	}
	n.lisClose = sync.Once{} // reset once
	close(n.listening)
	n.listenerLock.Unlock()

	defer func() { // refresh listener to re-call
		n.lisClose.Do(func() {
			n.listenerLock.Lock()
			n.lis.Close()
			n.lis = nil
			n.listenerLock.Unlock()
		})
		n.listening = make(chan struct{})
	}()

	// accept in a goroutine so that we can react to context cancel and close the listener
	newConn := make(chan net.Conn)
	go func() {
		defer close(newConn)
		for {
			n.listenerLock.Lock()
			if n.lis == nil {
				n.listenerLock.Unlock()
				return
			}
			n.listenerLock.Unlock()
			conn, err := n.lis.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					// yikes way of handling this
					// but means this needs to be restarted anyway
					return
				}

				continue
			}

			newConn <- conn
		}
	}()

	defer level.Debug(n.log).Log("event", "network listen loop exited")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case conn := <-newConn:
			if conn == nil {
				return nil
			}
			go n.handleConnection(ctx, conn, true, wrappers...)
		}
	}
}

// GetListenAddr waits for Serve() to be called!
func (n *node) GetListenAddr() net.Addr {
	_, ok := <-n.listening
	if !ok {
		return n.lis.Addr()
	}
	level.Error(n.log).Log("msg", "listener not ready")
	return nil
}

// GetAllEndpoints should merge with conntracker
func (n *node) GetAllEndpoints() []EndpointStat {
	n.remotesLock.Lock()
	defer n.remotesLock.Unlock()

	var stats []EndpointStat

	for ref, edp := range n.remotes {
		id, _ := refs.ParseFeedRef(ref)
		remote := edp.Remote()
		ok, durr := n.connTracker.Active(remote)
		if !ok {
			continue
		}
		stats = append(stats, EndpointStat{
			ID:       &id,
			Addr:     remote,
			Since:    durr,
			Endpoint: edp,
		})
	}
	return stats
}

func (n *node) GetConnTracker() ConnTracker {
	return n.connTracker
}

func (n *node) HandleHTTP(h http.Handler) {
	n.httpHandler = h
}

func (n *node) Close() error {
	if n.httpLis != nil {
		err := n.httpLis.Close()
		if err != nil {
			return fmt.Errorf("[network/Close]: failed to close http listener: %w", err)
		}
	}
	n.listenerLock.Lock()
	defer n.listenerLock.Unlock()
	if n.lis != nil {
		var closeErr error
		n.lisClose.Do(func() {
			closeErr = n.lis.Close()
		})
		if closeErr != nil && !strings.Contains(closeErr.Error(), "use of closed network connection") {
			return fmt.Errorf("[network/Close]:: network node failed to close it's listener: %w", closeErr)
		}
	}

	n.remotesLock.Lock()
	defer n.remotesLock.Unlock()
	for addr, edp := range n.remotes {
		if err := edp.Terminate(); err != nil {
			n.log.Log("event", "failed to terminate endpoint", "addr", addr, "err", err)
		}
	}

	if cnt := n.connTracker.Count(); cnt > 0 {
		n.log.Log("event", "warning", "msg", "still open connections", "count", cnt)
		n.connTracker.CloseAll()
	}

	return nil
}

func (n *node) applyConnWrappers(conn net.Conn) (net.Conn, error) {
	for i, cw := range n.afterSecureConnWrappers {
		var err error
		conn, err = cw(conn)
		if err != nil {
			return nil, fmt.Errorf("[network/applyConnWrappers]:error applying connection wrapper #%d: %w", i, err)
		}
	}
	return conn, nil
}
