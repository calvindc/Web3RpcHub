package service

import (
	"fmt"
	"net"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/internal/network"
)

//initNetwork opens the shs listener for TCP connections
func (svr *HubServe) initNetwork() error {
	mHandler := func(conn net.Conn) (cmuxrpc.Handler, error) {
		svr.closedMutex.Lock()
		defer svr.closedMutex.Unlock()

		remote, err := network.GetFeedRefFromAddr(conn.RemoteAddr())
		if err != nil {
			return nil, fmt.Errorf("[service/init_network]: expected an address containing an shs-bs addr: %w", err)
		}

		pm, err := svr.Config.GetPrivacyMode(svr.servCtx)
		if err != nil {
			return nil, fmt.Errorf("[service/init_network]: running with unknown privacy mode")
		}

		// if privacy mode is restricted, deny connections from non-members
		if pm == db.ModeRestricted {
			if _, err := svr.Members.GetByFeed(svr.servCtx, remote); err != nil {
				return nil, fmt.Errorf("[service/init_network]: access restricted to members")
			}
		}

		// if feed is in the deny list, deny their connection
		if svr.DeniedKeys.HasFeed(svr.servCtx, remote) {
			return nil, fmt.Errorf("[service/init_network]: this key has been banned")
		}
		// for community + open modes, allow all connections
		return &svr.public, nil
	}

	// tcp+shs
	opts := network.Options{
		Logger:              svr.logger,
		Dialer:              svr.dialer,
		ListenAddr:          svr.listenAddr,
		KeyPair:             svr.keyPair,
		AppKey:              svr.appKey[:],
		MakeHandler:         mHandler,
		ConnTracker:         svr.networkConnTracker,
		BefreCryptoWrappers: svr.preSecureWrappers,
		AfterSecureWrappers: svr.postSecureWrappers,
	}
	var err error
	svr.Network, err = network.New(opts)
	if err != nil {
		return fmt.Errorf("[service/init_network]: failed to create network node: %w", err)
	}

	return nil
}
