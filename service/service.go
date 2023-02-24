package service

import (
	"context"

	"sync"

	"net"

	"github.com/calvindc/Web3RpcHub/internal"
	"github.com/calvindc/Web3RpcHub/network"
	"github.com/go-kit/kit/log"
	"go.cryptoscope.co/muxrpc/v2/typemux"
	"go.cryptoscope.co/netwrap"
)

type HubService struct {
	StateManager *network.HubNetManager

	logger  log.Logger
	cc      context.Context
	cf      context.CancelFunc
	closers internal.Closer

	closed      bool
	closedMutex sync.Mutex
	closeErr    error

	Network      network.Network
	appKey       []byte
	listenAddr   net.Addr
	wsAddr       string
	dialer       netwrap.Dialer
	netInfo      network.HubEndpoint
	loadUnixSock bool
	repo         internal.GetPathInterface
	repoPath     string
	keyPair      *internal.KeyPair

	networkConnTracker network.ConnTracker
	preSecureWrappers  []netwrap.ConnWrapper //所有呼叫过的连接管理器
	postSecureWrappers []netwrap.ConnWrapper //目前有效的连接管理器

	public typemux.HandlerMux
	master typemux.HandlerMux
}
