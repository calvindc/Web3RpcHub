package mainimpl

import (
	"os/user"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"os"

	"flag"

	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

var (
	development = false

	Appversion = ""
	Commit     = ""
	Goversion  = ""
	Builddate  = ""
)
var (
	log kitlog.Logger
)

var (
	flagprintversion bool
)

/*var app = cli.App{
	Name:    os.Args[0],
	Usage:   "web3' peer communication management",
	Version: Version,

	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "print-version", Usage: SvrCfg_PrintVersion_I},
		&cli.StringFlag{Name: "secret-handsharke-key", Usage: SvrCfg_SecretHandsharkeKey_I, Value: SvrCfg_SecretHandsharkeKey},
		&cli.StringFlag{Name: "listen-addr-shsmux", Usage: SvrCofg_ListenAddrShsMux_I, Value: SvrCofg_ListenAddrShsMux},
		&cli.StringFlag{Name: "listen-addr-http", Usage: SvrCfg_ListenAddrHttp, Value: SvrCfg_ListenAddrHttp},
		&cli.BoolFlag{Name: "enable-unixsock", Usage: SvrCfg_EnableUnixSock_I},
		&cli.StringFlag{Name: "repo-dir", Usage: SvrCfg_RepoDir, Value: SvrCfg_RepoDir_I},
		&cli.StringFlag{Name: "listen-addr-metrics-pprof", Usage: SvrCfg_ListenAddrMetricsPprof, Value: SvrCfg_ListenAddrMetricsPprof},
		&cli.StringFlag{Name: "https-domain", Usage: SvrCfg_HttpsDomain_I, Value: SvrCfg_HttpsDomain},
		&cli.StringFlag{Name: "hub-mode", Usage: SvrCfg_HubMode_I,Value:SvrCfg_HubMode},
		&cli.BoolFlag{Name: "aliases-as-subdomains", Usage: SvrCfg_AliasesAsSubdomains_I},
	},
	Before:
}*/

func checkAndLog(err error) {
	if err != nil {
		level.Error(log).Log("event", "fatal error", "err", err)

	}
}

func initRunConfig() error {
	u, err := user.Current()
	checkAndLog(err)
	if err != nil {
		level.Error(log).Log("event", "fatal error", "err", err)
	}
	flag.BoolVar(&SvrCfg_PrintVersion, "print-version", SvrCfg_PrintVersion, SvrCfg_PrintVersion_I)
	flag.StringVar(&SvrCfg_SecretHandsharkeKey, "secret-handsharke-key", SvrCfg_SecretHandsharkeKey, SvrCfg_SecretHandsharkeKey_I)
	flag.StringVar(&SvrCofg_ListenAddrShsMux, "listen-addr-shsmux", SvrCofg_ListenAddrShsMux, SvrCofg_ListenAddrShsMux_I)
	flag.StringVar(&SvrCfg_ListenAddrHttp, "listen-addr-http", SvrCfg_ListenAddrHttp, SvrCfg_ListenAddrHttp_I)
	flag.BoolVar(&SvrCfg_EnableUnixSock, "enable-unixsock", SvrCfg_EnableUnixSock, SvrCfg_EnableUnixSock_I)
	flag.StringVar(&SvrCfg_RepoDir, "repo-dir", filepath.Join(u.HomeDir, SvrCfg_RepoDir), SvrCfg_RepoDir_I)
	flag.StringVar(&SvrCfg_LogDir, "log-dir", SvrCfg_LogDir, SvrCfg_LogDir_I)
	flag.StringVar(&SvrCfg_ListenAddrMetricsPprof, "listen-addr-metrics-pprof", SvrCfg_ListenAddrMetricsPprof, SvrCfg_ListenAddrMetricsPprof_I)
	flag.StringVar(&SvrCfg_HttpsDomain, "https-domain", SvrCfg_HttpsDomain, SvrCfg_HttpsDomain_I)
	flag.Func("hub-mode", SvrCfg_HubMode_I, SvrCfg_HubMode)
	flag.BoolVar(&SvrCfg_AliasesAsSubdomains, "aliases-as-subdomains", SvrCfg_AliasesAsSubdomains, SvrCfg_AliasesAsSubdomains_I)
	flag.Parse()
	if SvrCfg_LogDir != "" {
		logDir := filepath.Join(SvrCfg_RepoDir, SvrCfg_LogDir)
		os.MkdirAll(logDir, 0700)
		logFileName := fmt.Sprintf("%s-%s.log", filepath.Base(os.Args[0]), time.Now().Format("2006-01-02_15-04"))
		logFile, err := os.Create(filepath.Join(logDir, logFileName))
		if err != nil {
			panic(err)
		}
		log = kitlog.NewJSONLogger(kitlog.NewSyncWriter(logFile))
	} else {
		log = kitlog.NewLogfmtLogger(os.Stderr)
	}

	if SvrCfg_PrintVersion {
		level.Info(log).Log("version", Appversion, "commit", Commit, "goVersion", Goversion, "builddate", Builddate)
	}

	if SvrCfg_HttpsDomain == "" {
		if !development {
			return fmt.Errorf("https-domain can't be empty. See '%s -h' for a full list of options", os.Args[0])
		}
		SvrCfg_HttpsDomain = "localhost"
	}

	_, portMuxRPC, err := net.SplitHostPort(SvrCofg_ListenAddrShsMux)
	if err != nil {
		return fmt.Errorf("invalid muxrpc listener: %w", err)
	}
	_, err = net.LookupPort("tcp", portMuxRPC)
	if err != nil {
		return fmt.Errorf("invalid tcp port for muxrpc listener: %w", err)
	}

	_, portHTTPStr, err := net.SplitHostPort(SvrCfg_ListenAddrHttp)
	if err != nil {
		return fmt.Errorf("invalid http listener: %w", err)
	}
	_, err = net.LookupPort("tcp", portHTTPStr)
	if err != nil {
		return fmt.Errorf("invalid tcp port for muxrpc listener: %w", err)
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = base64.StdEncoding.DecodeString(SvrCfg_SecretHandsharkeKey)
	if err != nil {
		return fmt.Errorf("secret-handshake appkey is invalid base64: %w", err)
	}

	if SvrCfg_ListenAddrMetricsPprof != "" {
		go func() {
			level.Debug(log).Log("starting", "metrics", "addr", SvrCfg_ListenAddrMetricsPprof)
			err := http.ListenAndServe(SvrCfg_ListenAddrMetricsPprof, nil)
			checkAndLog(err)
		}()
	}
	return nil
}

func Runhubsvr() error {
	initRunConfig()
	return nil
}
