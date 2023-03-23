package service

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
	"github.com/calvindc/Web3RpcHub/internal/frequently"
	"github.com/calvindc/Web3RpcHub/internal/repository"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// WithUNIXSocket enables listening for muxrpc connections on a unix socket files ($repo/socket).
// This socket is not encrypted or authenticated since access to it is mediated by filesystem ownership.
func WithUNIXSocket(yes bool) Option {
	return func(s *HubServe) error {
		s.loadUnixSock = yes
		return nil
	}
}

// creates the UNIX socket file listener for local usage
func (s *HubServe) initUnixSock() error {
	if s.keyPair == nil {
		return fmt.Errorf("roomsrv/unixsock: keypair is nil. please use unixSocket with LateOption")
	}
	spoofWrapper := frequently.SpoofRemoteAddress(s.keyPair.Feed.PubKey())

	r := repository.New(s.repoPath)
	sockPath := r.GetPath("socket")

	// local clients (not using network package because we don't want conn limiting or advertising)
	c, err := net.Dial("unix", sockPath)
	if err == nil {
		c.Close()
		return fmt.Errorf("roomsrv: repo already in use, socket accepted connection")
	}
	os.Remove(sockPath)
	os.MkdirAll(filepath.Dir(sockPath), 0700)

	uxLis, err := net.Listen("unix", sockPath)
	if err != nil {
		return err
	}
	s.closers.Add(uxLis)

	go func() {

	acceptLoop:
		for {
			c, err := uxLis.Accept()
			if err != nil {
				if nerr, ok := err.(*net.OpError); ok {
					if nerr.Err.Error() == "use of closed network connection" {
						return
					}
				}

				level.Warn(s.logger).Log("event", "unix sock accept failed", "err", err)
				continue
			}

			wc, err := spoofWrapper(c)
			if err != nil {
				c.Close()
				continue
			}
			for _, w := range s.postSecureWrappers {
				var err error
				wc, err = w(wc)
				if err != nil {
					level.Warn(s.logger).Log("err", err)
					c.Close()
					continue acceptLoop
				}
			}

			go func(conn net.Conn) {
				defer conn.Close()

				pkr := cmuxrpc.NewPacker(conn)

				edp := cmuxrpc.Handle(pkr, &s.master,
					cmuxrpc.WithContext(s.servCtx),
					cmuxrpc.WithLogger(log.NewNopLogger()),
				)

				srv := edp.(cmuxrpc.Server)
				if err := srv.Serve(); err != nil {
					level.Warn(s.logger).Log("conn", "serve exited", "err", err, "peer", conn.RemoteAddr())
				}
				edp.Terminate()

			}(wc)
		}
	}()
	return nil

}
