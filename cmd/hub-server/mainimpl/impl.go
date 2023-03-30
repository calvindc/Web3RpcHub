package mainimpl

import (
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"fmt"

	"context"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/calvindc/Web3RpcHub/cmuxrpc/debug"
	"github.com/calvindc/Web3RpcHub/config"
	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/db/sqlite"
	"github.com/calvindc/Web3RpcHub/internal/network"
	"github.com/calvindc/Web3RpcHub/internal/repository"
	"github.com/calvindc/Web3RpcHub/internal/signalbridge"
	"github.com/calvindc/Web3RpcHub/models/web/handlers"
	"github.com/calvindc/Web3RpcHub/service"
	log2 "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	development = false

	Appversion = ""
	Commit     = ""
	Goversion  = ""
	Builddate  = ""
)
var (
	log log2.Logger
)

var (
	repoDir          string
	flagprintversion bool
)

type Cfdata struct {
	deAppKey []byte
	portHttp int
}

func CheckAndLog(err error) {
	if err != nil {
		level.Error(log).Log("event", "fatal error", "err", err)

	}
}

func Runhubsvr() error {
	cfdata, err := initRunConfig()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []service.Option{
		service.RegLogger(log),
		service.RegAppKey(cfdata.deAppKey),
		service.RegRepoPath(config.SvrCfg_RepoDir),
		service.RegUNIXSocket(config.SvrCfg_EnableUnixSock),
	}
	if config.SvrCfg_LogDir != "" {
		opts = append(opts, service.RegPostSecureConnWrapper(func(conn net.Conn) (net.Conn, error) {
			parts := strings.Split(conn.RemoteAddr().String(), "|")
			if len(parts) != 2 {
				return conn, nil
			}
			cmuxrpcDumpDir := filepath.Join(config.SvrCfg_RepoDir, config.SvrCfg_LogDir, parts[1], parts[0])
			return debug.WrapDump(cmuxrpcDumpDir, conn)
		}))
	}

	//系统运行性能监视
	if config.SvrCfg_ListenAddrMetricsPprof != "" {
		go func() {
			level.Debug(log).Log("starting", "metrics", "addr", config.SvrCfg_ListenAddrMetricsPprof)
			err := http.ListenAndServe(config.SvrCfg_ListenAddrMetricsPprof, nil)
			CheckAndLog(err)
		}()
	}

	// 新建一个hub-key
	r := repository.New(config.SvrCfg_RepoDir)
	keyPair, err := repository.DefaultKeyPair(r)
	CheckAndLog(err)
	opts = append(opts, service.RegKeyPair(keyPair))

	//hub network endpoint
	networkInfo := network.HubEndpoint{
		HubID:                  keyPair.Feed,
		ListenAddressMUXRPC:    config.SvrCofg_ListenAddrShsMux,
		HttpsDomain:            config.SvrCfg_HttpsDomain,
		HttpsPort:              uint(cfdata.portHttp),
		UseAliasesAsSubdomains: config.SvrCfg_AliasesAsSubdomains,
		Development:            development,
	}

	// setup a db
	hubdb, err := sqlite.OpenDB(r)
	if err != nil {
		return fmt.Errorf("failed to init database: %w", err)
	}

	signal_bridge := signalbridge.NewSignalBridge()
	if config.SvrCfg_HubMode == db.ModeUnknown {
		hubdb.Config.SetPrivacyMode(ctx, config.SvrCfg_HubMode)
	}

	//启动shs+cmuxrpc服务
	hubsrv, err := service.StartHubServ(hubdb.Members, hubdb.DeniedKeys, hubdb.Aliases, hubdb.AuthWitToken,
		signal_bridge, hubdb.Config, networkInfo, opts...)
	if err != nil {
		return fmt.Errorf("failed to instantiate hub server: %w", err)
	}

	//启动HTTP listener
	httpLis, err := net.Listen("tcp", config.SvrCfg_ListenAddrHttp)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		level.Warn(log).Log("event", "killed", "msg", "received signal, shutting down", "signal", sig.String())

		cancel()
		time.Sleep(1 * time.Second)
		hubsrv.ServShutDown()

		time.Sleep(1 * time.Second)
		err = httpLis.Close()
		CheckAndLog(err)

		err := hubsrv.ShotDown()
		CheckAndLog(err)

		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	// HTTPS重定向和Cryptographic Service Provider
	secureMiddleware := SetupSecureMiddleware(config.SvrCfg_HttpsDomain)

	// HTTP请求限流
	httpRateLimiter, err := ThrottleHttp()
	if err != nil {
		return err
	}

	var httpHandler http.Handler //add all /mux_user_api_handler(1,2,3...)
	webHandler, err := handlers.NewWebHandler(
		log2.With(log, "package", "web"),
		repository.New(repoDir),
		networkInfo,
		hubsrv.StateManager,
		hubsrv.Network,
		signal_bridge,
		handlers.Databases{
			Aliases:       hubdb.Aliases,
			AuthFallback:  hubdb.AuthFallback,
			AuthWithToken: hubdb.AuthWitToken,
			Config:        hubdb.Config,
			DeniedKeys:    hubdb.DeniedKeys,
			Invites:       hubdb.Invites,
			Notices:       hubdb.Notices,
			Members:       hubdb.Members,
			PinnedNotices: hubdb.PinnedNotices,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create HTTPdashboard handler: %w", err)
	}
	httpHandler = httpRateLimiter.RateLimit(webHandler)
	httpHandler = secureMiddleware.Handler(httpHandler)
	httpHandler = hubsrv.Network.WebsockHandler(httpHandler)

	// all init was successfull
	level.Info(log).Log(
		"event", "serving",
		"ID", hubsrv.Whoami().String(),
		"shsmuxaddr", config.SvrCofg_ListenAddrShsMux,
		"httpaddr", config.SvrCfg_ListenAddrHttp,
		"version", Appversion, "commit", Commit,
	)

	// 启动hub的http连接服务
	httpSrv := http.Server{
		Addr:              httpLis.Addr().String(),
		ReadHeaderTimeout: 15 * time.Second, //to avoid Slowloris attacks
		WriteTimeout:      1 * time.Minute,
		IdleTimeout:       1 * time.Minute,
		Handler:           httpHandler,
	}
	err = httpSrv.Serve(httpLis)
	if err != nil {
		level.Error(log).Log("event", "http serve failed", "err", err)
	}

	// 启动hub的shs+rpc连接服务
	for {
		err := hubsrv.Network.Serve(ctx)
		if err != nil {
			level.Warn(log).Log("event", "hubsrv node.Serve returned", "err", err)
		}
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			err := hubsrv.ShotDown()
			return err
		default:

		}

	}

	return nil
}
