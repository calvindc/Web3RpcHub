package mainimpl

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/calvindc/Web3RpcHub/config"

	"github.com/calvindc/Web3RpcHub/db"
	log2 "go.mindeco.de/log"
	"go.mindeco.de/log/level"
)

func initRunConfig() (*Cfdata, error) {
	u, err := user.Current()
	CheckAndLog(err)
	if err != nil {
		level.Error(log).Log("event", "fatal error", "err", err)
	}
	flag.BoolVar(&config.SvrCfg_PrintVersion, "print-version", config.SvrCfg_PrintVersion, config.SvrCfg_PrintVersion_I)
	flag.StringVar(&config.SvrCfg_SecretHandsharkeKey, "secret-handsharke-key", config.SvrCfg_SecretHandsharkeKey, config.SvrCfg_SecretHandsharkeKey_I)
	flag.StringVar(&config.SvrCofg_ListenAddrShsMux, "listen-addr-shsmux", config.SvrCofg_ListenAddrShsMux, config.SvrCofg_ListenAddrShsMux_I)
	flag.StringVar(&config.SvrCfg_ListenAddrHttp, "listen-addr-http", config.SvrCfg_ListenAddrHttp, config.SvrCfg_ListenAddrHttp_I)
	flag.BoolVar(&config.SvrCfg_EnableUnixSock, "enable-unixsock", config.SvrCfg_EnableUnixSock, config.SvrCfg_EnableUnixSock_I)
	flag.StringVar(&config.SvrCfg_RepoDir, "repo-dir", filepath.Join(u.HomeDir, config.SvrCfg_RepoDir), config.SvrCfg_RepoDir_I)
	flag.StringVar(&config.SvrCfg_LogDir, "log-dir", config.SvrCfg_LogDir, config.SvrCfg_LogDir_I)
	flag.StringVar(&config.SvrCfg_ListenAddrMetricsPprof, "listen-addr-metrics-pprof", config.SvrCfg_ListenAddrMetricsPprof, config.SvrCfg_ListenAddrMetricsPprof_I)
	flag.StringVar(&config.SvrCfg_HttpsDomain, "https-domain", config.SvrCfg_HttpsDomain, config.SvrCfg_HttpsDomain_I)
	flag.Func("hub-mode", config.SvrCfg_HubMode_I, func(val string) error {
		ppm := db.ParsePrivacyMode(val)
		if err := ppm.IsValid(); err != nil {
			return err
		}
		config.SvrCfg_HubMode = ppm
		return nil
	})
	flag.BoolVar(&config.SvrCfg_AliasesAsSubdomains, "aliases-as-subdomains", config.SvrCfg_AliasesAsSubdomains, config.SvrCfg_AliasesAsSubdomains_I)
	flag.Parse()

	if config.SvrCfg_LogDir != "" {
		logDir := filepath.Join(config.SvrCfg_RepoDir, config.SvrCfg_LogDir)
		os.MkdirAll(logDir, 0700)
		logFileName := fmt.Sprintf("%s-%s.log", filepath.Base(os.Args[0]), time.Now().Format("2006-01-02_15-04"))
		logFile, err := os.Create(filepath.Join(logDir, logFileName))
		if err != nil {
			panic(err)
		}
		log = log2.NewJSONLogger(log2.NewSyncWriter(logFile))
	} else {
		log = log2.NewLogfmtLogger(os.Stderr)
	}

	if config.SvrCfg_PrintVersion {
		level.Info(log).Log("version", Appversion, "commit", Commit, "goVersion", Goversion, "builddate", Builddate)
	}

	if config.SvrCfg_HttpsDomain == "" {
		if !development {
			return nil, fmt.Errorf("https-domain can't be empty. See '%s -h' for a full list of options", os.Args[0])
		}
		config.SvrCfg_HttpsDomain = "localhost"
	}

	_, portMuxRPC, err := net.SplitHostPort(config.SvrCofg_ListenAddrShsMux)
	if err != nil {
		return nil, fmt.Errorf("invalid muxrpc listener: %w", err)
	}
	_, err = net.LookupPort("tcp", portMuxRPC)
	if err != nil {
		return nil, fmt.Errorf("invalid tcp port for muxrpc listener: %w", err)
	}

	_, portHTTPStr, err := net.SplitHostPort(config.SvrCfg_ListenAddrHttp)
	if err != nil {
		return nil, fmt.Errorf("invalid http listener: %w", err)
	}
	portHTTP, err := net.LookupPort("tcp", portHTTPStr)
	if err != nil {
		return nil, fmt.Errorf("invalid tcp port for muxrpc listener: %w", err)
	}

	dshs, err := base64.StdEncoding.DecodeString(config.SvrCfg_SecretHandsharkeKey)
	if err != nil {
		return nil, fmt.Errorf("secret-handshake appkey is invalid base64: %w", err)
	}

	cf := &Cfdata{
		deAppKey: dshs,
		portHttp: portHTTP,
	}
	return cf, nil
}
