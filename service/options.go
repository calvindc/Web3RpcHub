package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/calvindc/Web3RpcHub/internal/keys"
	"github.com/calvindc/Web3RpcHub/internal/repository"
	"github.com/calvindc/Web3RpcHub/netwrap"
	"github.com/go-kit/log"
)

type Option func(srv *HubServe) error

// RegRepoPath changes where the replication database and blobs are stored.
func RegRepoPath(path string) Option {
	return func(s *HubServe) error {
		s.repoPath = path
		return nil
	}
}

func RegAppKey(k []byte) Option {
	return func(s *HubServe) error {
		if n := len(k); n != 32 {
			return fmt.Errorf("appKey: need 32 bytes got %d", n)
		}
		s.appKey = k
		return nil
	}
}

func RegNamedKeyPair(name string) Option {
	return func(s *HubServe) error {
		r := repository.New(s.repoPath)
		var err error
		s.keyPair, err = repository.LoadKeyPair(r, name)
		if err != nil {
			return fmt.Errorf("loading named key-pair %q failed: %w", name, err)
		}
		return nil
	}
}

func RegJSONKeyPair(blob string) Option {
	return func(s *HubServe) error {
		var err error
		s.keyPair, err = keys.ParseKeyPair(strings.NewReader(blob))
		if err != nil {
			return fmt.Errorf("JSON KeyPair decode failed: %w", err)
		}
		return nil
	}
}

func RegKeyPair(kp *keys.KeyPair) Option {
	return func(s *HubServe) error {
		s.keyPair = kp
		return nil
	}
}

func RegLogger(log log.Logger) Option {
	return func(s *HubServe) error {
		s.logger = log
		return nil
	}
}

func RegContext(ctx context.Context) Option {
	return func(s *HubServe) error {
		s.servCtx, s.servShutDown = context.RegCancel(ctx)
		return nil
	}
}

func RegDialer(dial netwrap.Dialer) Option {
	return func(s *HubServe) error {
		s.dialer = dial
		return nil
	}
}

func RegNetworkConnTracker(ct hubstate.ConnTracker) Option {
	return func(s *HubServe) error {
		s.networkConnTracker = ct
		return nil
	}
}

func RegPreSecureConnWrapper(cw netwrap.ConnWrapper) Option {
	return func(s *HubServe) error {
		s.preSecureWrappers = append(s.preSecureWrappers, cw)
		return nil
	}
}

func RegPostSecureConnWrapper(cw netwrap.ConnWrapper) Option {
	return func(s *HubServe) error {
		s.postSecureWrappers = append(s.postSecureWrappers, cw)
		return nil
	}
}
