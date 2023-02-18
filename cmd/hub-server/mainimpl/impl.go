package mainimpl

import (
	"os/user"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"os"

	"github.com/urfave/cli"
)

var (
	Version   = ""
	BuildDate = ""
	GoVersion = ""
	GitCommit = ""
)
var (
	log kitlog.Logger
)

var app = cli.App{
	Name:    os.Args[0],
	Usage:   "web3' peer communication management",
	Version: Version,

	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "dd", Usage: SvrCfg_PrintVersion_I},
	},
}

func checkAndLog(err error) {
	if err != nil {
		level.Error(log).Log("event", "fatal error", "err", err)

	}
}
func initRunConfig() {
	_, err := user.Current()
	checkAndLog(err)
	if err != nil {
		level.Error(log).Log("event", "fatal error", "err", err)
	}

}
func Runhubsvr() error {
	initRunConfig()
	return nil
}
