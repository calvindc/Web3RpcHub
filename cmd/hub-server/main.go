package main

import (
	"fmt"
	"os"

	"github.com/calvindc/Web3RpcHub/cmd/hub-server/mainimpl"
)

func main() {
	if err := mainimpl.Runhubsvr(); err != nil {
		fmt.Fscan(os.Stdout, "web3-rpc-hub: %s\n", err)
		os.Exit(1)
	}
}
