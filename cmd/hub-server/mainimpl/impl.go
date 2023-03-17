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

	"github.com/calvindc/Web3RpcHub/config"
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
	flag.BoolVar(&config.SvrCfg_PrintVersion, "print-version", config.SvrCfg_PrintVersion, config.SvrCfg_PrintVersion_I)
	flag.StringVar(&config.SvrCfg_SecretHandsharkeKey, "secret-handsharke-key", config.SvrCfg_SecretHandsharkeKey, config.SvrCfg_SecretHandsharkeKey_I)
	flag.StringVar(&config.SvrCofg_ListenAddrShsMux, "listen-addr-shsmux", config.SvrCofg_ListenAddrShsMux, config.SvrCofg_ListenAddrShsMux_I)
	flag.StringVar(&config.SvrCfg_ListenAddrHttp, "listen-addr-http", config.SvrCfg_ListenAddrHttp, config.SvrCfg_ListenAddrHttp_I)
	flag.BoolVar(&config.SvrCfg_EnableUnixSock, "enable-unixsock", config.SvrCfg_EnableUnixSock, config.SvrCfg_EnableUnixSock_I)
	flag.StringVar(&config.SvrCfg_RepoDir, "repo-dir", filepath.Join(u.HomeDir, config.SvrCfg_RepoDir), config.SvrCfg_RepoDir_I)
	flag.StringVar(&config.SvrCfg_LogDir, "log-dir", config.SvrCfg_LogDir, config.SvrCfg_LogDir_I)
	flag.StringVar(&config.SvrCfg_ListenAddrMetricsPprof, "listen-addr-metrics-pprof", config.SvrCfg_ListenAddrMetricsPprof, config.SvrCfg_ListenAddrMetricsPprof_I)
	flag.StringVar(&config.SvrCfg_HttpsDomain, "https-domain", config.SvrCfg_HttpsDomain, config.SvrCfg_HttpsDomain_I)
	flag.Func("hub-mode", config.SvrCfg_HubMode_I, config.SvrCfg_HubMode)
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
		log = kitlog.NewJSONLogger(kitlog.NewSyncWriter(logFile))
	} else {
		log = kitlog.NewLogfmtLogger(os.Stderr)
	}

	if config.SvrCfg_PrintVersion {
		level.Info(log).Log("version", Appversion, "commit", Commit, "goVersion", Goversion, "builddate", Builddate)
	}

	if config.SvrCfg_HttpsDomain == "" {
		if !development {
			return fmt.Errorf("https-domain can't be empty. See '%s -h' for a full list of options", os.Args[0])
		}
		config.SvrCfg_HttpsDomain = "localhost"
	}

	_, portMuxRPC, err := net.SplitHostPort(config.SvrCofg_ListenAddrShsMux)
	if err != nil {
		return fmt.Errorf("invalid muxrpc listener: %w", err)
	}
	_, err = net.LookupPort("tcp", portMuxRPC)
	if err != nil {
		return fmt.Errorf("invalid tcp port for muxrpc listener: %w", err)
	}

	_, portHTTPStr, err := net.SplitHostPort(config.SvrCfg_ListenAddrHttp)
	if err != nil {
		return fmt.Errorf("invalid http listener: %w", err)
	}
	_, err = net.LookupPort("tcp", portHTTPStr)
	if err != nil {
		return fmt.Errorf("invalid tcp port for muxrpc listener: %w", err)
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = base64.StdEncoding.DecodeString(config.SvrCfg_SecretHandsharkeKey)
	if err != nil {
		return fmt.Errorf("secret-handshake appkey is invalid base64: %w", err)
	}

	if config.SvrCfg_ListenAddrMetricsPprof != "" {
		go func() {
			level.Debug(log).Log("starting", "metrics", "addr", config.SvrCfg_ListenAddrMetricsPprof)
			err := http.ListenAndServe(config.SvrCfg_ListenAddrMetricsPprof, nil)
			checkAndLog(err)
		}()
	}
	return nil
}

func Runhubsvr() error {
	initRunConfig()
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
