package service

import (
	"net"

	"fmt"

	"github.com/calvindc/Web3RpcHub/db"
	"github.com/calvindc/Web3RpcHub/internal/network"
	"go.cryptoscope.co/muxrpc/v2"
)

func (svr *HubServe) initNetwork() error {

	mHandle := func(conn net.Conn) (muxrpc.Handler, error) {
		svr.closedMutex.Lock()
		defer svr.closedMutex.Unlock()

		remote, err := network.GetFeedRefFromAddr(conn.RemoteAddr())
		if err != nil {
			return nil, fmt.Errorf("sbot: expected an address containing an shs-bs addr: %w", err)
		}

		pm, err := svr.Config.GetPrivacyMode(svr.servCtx)
		if err != nil {
			return nil, fmt.Errorf("running with unknown privacy mode")
		}

		// if privacy mode is restricted, deny connections from non-members
		if pm == db.ModeRestricted {
			if _, err := svr.Members.GetByFeed(svr.servCtx, remote); err != nil {
				return nil, fmt.Errorf("access restricted to members")
			}
		}

		// if feed is in the deny list, deny their connection
		if svr.DeniedKeys.HasFeed(svr.servCtx, remote) {
			return nil, fmt.Errorf("this key has been banned")
		}
		// for community + open modes, allow all connections
		return &svr.public, nil
	}

	//opts := network.
	return nil
}
