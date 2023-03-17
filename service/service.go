package service

import (
	"context"

	"sync"

	"net"

	"fmt"
	"os/user"
	"path/filepath"

	"encoding/base64"
	"os"

	"github.com/calvindc/Web3RpcHub/config"
	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/hubstat"
	"github.com/calvindc/Web3RpcHub/internal/frequently"
	"github.com/calvindc/Web3RpcHub/internal/keys"
	"github.com/calvindc/Web3RpcHub/internal/network"
	"github.com/calvindc/Web3RpcHub/internal/netwrap"
	"github.com/calvindc/Web3RpcHub/internal/repository"
	"github.com/calvindc/Web3RpcHub/internal/signalbridge"
	"github.com/go-kit/kit/log"
	"go.cryptoscope.co/muxrpc/v2/typemux"
)

type HubServe struct {
	StateManager *hubstat.HubNetManager //网络节点管理

	logger       log.Logger      //日志调度
	servCtx      context.Context //hub service的上下文
	servShutDown context.CancelFunc
	closers      frequently.Closer

	closed      bool
	closedMutex *sync.Mutex
	closeErr    error

	Network      network.Network //8个网络基本接口服务,包括conn-Tracker
	appKey       []byte
	listenAddr   net.Addr //支持各种网络连接协议
	wsAddr       string
	dialer       netwrap.Dialer
	netInfo      network.HubEndpoint
	loadUnixSock bool
	repo         repository.Interface
	repoPath     string
	keyPair      *keys.KeyPair //web3r协议规范

	networkConnTracker network.ConnTracker   //hub的连接Tracker
	preSecureWrappers  []netwrap.ConnWrapper //所有呼叫过的连接管理器
	postSecureWrappers []netwrap.ConnWrapper //目前有效的连接管理器

	public typemux.HandlerMux //public路由 muxrpc.HandlerMux
	master typemux.HandlerMux //master路由 muxrpc.HandlerMux

	Members        db.MembersService          //成员管理
	DeniedKeys     db.DeniedKeysService       //屏蔽管理
	Aliases        db.AliasesService          //alia管理
	authWithToken  db.AuthWithTokenService    //访问web3r的token管理
	authWithBirdge *signalbridge.SignalBridge //http访问验证session过程
	Config         db.HubConfig               //hub配置管理(hub隐私和系统语言)
}

func StartHubServ(hMembers db.MembersService, hDeniedKeys db.DeniedKeysService, hAlias db.AliasesService,
	hAuthWithToken db.AuthWithTokenService, hAuthWithBirdge *signalbridge.SignalBridge,
	hConfig db.HubConfig, hNetInfo network.HubEndpoint, opts ...Option) (*HubServe, error) {
	var svr HubServe

	svr.closedMutex = new(sync.Mutex)
	svr.Members = hMembers
	svr.DeniedKeys = hDeniedKeys
	svr.Aliases = hAlias
	svr.authWithToken = hAuthWithToken
	svr.authWithBirdge = hAuthWithBirdge
	svr.netInfo = hNetInfo

	for x, optx := range opts {
		if err := optx(&svr); err != nil {
			return nil, fmt.Errorf("error applying option #%d: %w", x, err)
		}
	}

	if svr.repoPath == "" {
		u, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("error getting info on current user: %w", err)
		}
		svr.repoPath = filepath.Join(u.HomeDir, config.SvrCfg_RepoDir)
	}

	svr.repo = repository.New(svr.repoPath)

	if svr.appKey == nil {
		shk, err := base64.StdEncoding.DecodeString(config.SvrCfg_SecretHandsharkeKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode default appkey: %w", err)
		}
		svr.appKey = shk
	}

	if svr.dialer == nil {
		svr.dialer = netwrap.Dial
	}

	var err error
	svr.listenAddr, err = net.ResolveTCPAddr("tcp", svr.netInfo.ListenAddressMUXRPC)
	if err != nil {
		return nil, err
	}

	if svr.logger == nil {
		logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
		svr.logger = logger
	}

	svr.StateManager = hubstat.NewHubNetManager(svr.logger)

	svr.public = typemux.New(log.With(svr.logger, "mux", "public"))
	svr.master = typemux.New(log.With(svr.logger, "mux", "master"))

	if svr.servCtx == nil {
		svr.servCtx, svr.servShutDown = context.WithCancel(context.Background())
	}

	svr.netInfo.HubID = svr.keyPair.Feed

	//svr.initHanders()

}
