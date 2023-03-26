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

/*var app = cli.App{
	Name:    os.Args[0],
	Usage:   "web3' peer communication management",
	Version: Version,

	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "print-version", Usage: config.SvrCfg_PrintVersion_I},
		&cli.StringFlag{Name: "secret-handsharke-key", Usage: config.SvrCfg_SecretHandsharkeKey_I, Value: config.SvrCfg_SecretHandsharkeKey},
		&cli.StringFlag{Name: "listen-addr-shsmux", Usage: SvrCofg_ListenAddrShsMux_I, Value: SvrCofg_ListenAddrShsMux},
		&cli.StringFlag{Name: "listen-addr-http", Usage: config.SvrCfg_ListenAddrHttp, Value: config.SvrCfg_ListenAddrHttp},
		&cli.BoolFlag{Name: "enable-unixsock", Usage: config.SvrCfg_EnableUnixSock_I},
		&cli.StringFlag{Name: "repo-dir", Usage: config.SvrCfg_RepoDir, Value: config.SvrCfg_RepoDir_I},
		&cli.StringFlag{Name: "listen-addr-metrics-pprof", Usage: config.SvrCfg_ListenAddrMetricsPprof, Value: config.SvrCfg_ListenAddrMetricsPprof},
		&cli.StringFlag{Name: "https-domain", Usage: config.SvrCfg_HttpsDomain_I, Value: config.SvrCfg_HttpsDomain},
		&cli.StringFlag{Name: "hub-mode", Usage: config.SvrCfg_HubMode_I,Value:config.SvrCfg_HubMode},
		&cli.BoolFlag{Name: "aliases-as-subdomains", Usage: config.SvrCfg_AliasesAsSubdomains_I},
	},
	Before:
}*/

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

		time.Sleep(3 * time.Second)
		err := hubsrv.ShotDown()
		CheckAndLog(err)

		time.Sleep(3 * time.Second)
		os.Exit(0)
	}()

	// HTTPS重定向和Cryptographic Service Provider
	secureMiddleware := SetupSecureMiddleware(config.SvrCfg_HttpsDomain)

	// HTTP请求限流
	httpRateLimiter, err := ThrottleHttp()
	if err != nil {
		return err
	}

	webHandler := &http.ServeMux{}
	var httpHandler http.Handler //add all /mux_user_api_handler(1,2,3...)
	httpHandler = httpRateLimiter.RateLimit(webHandler)
	httpHandler = secureMiddleware.Handler(httpHandler)
	httpHandler = hubsrv.Network.WebsockHandler(httpHandler)

	level.Info(log).Log(
		"event", "serving", "running...", httpLis.Addr().String())
	/*fmt.Println("12")
	go func() {
		fiveDays := 3 * time.Second
		ticker := time.NewTicker(fiveDays)
		for range ticker.C {
			fmt.Println(time.Now().String())
		}
	}()

	for {
		time.Sleep(time.Second)
	}*/

	/*fiveDays := 3 * time.Second
	ticker := time.NewTicker(fiveDays)
	go func() { // server might not restart as often

		count := 0
		for {
			t := <-ticker.C
			count++
			fmt.Println("now=", t, "count=", count)

			if count == 5 {
				fmt.Println("now=%w", t)
				ticker.Stop()
				//runtime.Goexit()
				os.Exit(1)
			}
		}
	}()
	for {
		time.Sleep(time.Second)
	}*/

	/*ticker := time.NewTicker(time.Second * 2)
	c := make(chan bool)
	go func() {
		time.Sleep(time.Second * 7)
		c <- true
	}()
	for {
		select {
		case <-c:
			fmt.Println("completed")
			return nil

		case tm := <-ticker.C:
			fmt.Println(tm)
		}
	}*/

	/*c := make(chan string)
	go func() {
		for {
			time.Sleep(time.Second)
			c <- time.Now().String()
		}
	}()
	for {
		select {
		case str := <-c:
			fmt.Println("reveice ", str)
		case <-time.After(time.Second * 2):
			fmt.Println("time out")
			c <- "1"
		}
	}*/

	/*c := make(chan string, 10)
	go func() {
		time.Sleep(time.Second)
		for {
			<-c
			fmt.Println("delete data:")
			time.Sleep(5 * time.Second)
		}
	}()
	for {
		select {
		case c <- "-1":

			fmt.Println("add data: ", len(c))
			//os.Exit(1)
			time.Sleep(time.Second * 1)

		default:
			fmt.Println("zi yuan hao jin")
			//time.Sleep(time.Second)
		}///home/ubuntu/ssb-server-js/photon-ssb
	}*/
	return nil
}
