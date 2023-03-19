package typemux

import (
	"context"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var _ AsyncHandler = (*AsyncFunc)(nil)

type AsyncFunc func(context.Context, *cmuxrpc.Request) (interface{}, error)

func (af AsyncFunc) HandleAsync(ctx context.Context, r *cmuxrpc.Request) (interface{}, error) {
	return af(ctx, r)
}

type AsyncHandler interface {
	HandleAsync(context.Context, *cmuxrpc.Request) (interface{}, error)
}

type asyncStub struct {
	logger log.Logger

	h AsyncHandler
}

func (hm asyncStub) HandleCall(ctx context.Context, req *cmuxrpc.Request) {
	// TODO: check call type?

	v, err := hm.h.HandleAsync(ctx, req)
	if err != nil {
		req.CloseWithError(err)
		return
	}

	err = req.Return(ctx, v)
	if err != nil {
		level.Error(hm.logger).Log("evt", "return failed", "err", err, "method", req.Method.String())
	}
}

func (hm asyncStub) HandleConnect(ctx context.Context, edp cmuxrpc.Endpoint) {
	if ch, ok := hm.h.(cmuxrpc.ConnectHandler); ok {
		ch.HandleConnect(ctx, edp)
	}
}
