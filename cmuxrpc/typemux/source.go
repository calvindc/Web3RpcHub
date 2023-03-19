package typemux

import (
	"context"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
)

var _ SourceHandler = (*SourceFunc)(nil)

type SourceFunc func(context.Context, *cmuxrpc.Request, *cmuxrpc.ByteSink) error

func (sf SourceFunc) HandleSource(ctx context.Context, r *cmuxrpc.Request, src *cmuxrpc.ByteSink) error {
	return sf(ctx, r, src)
}

// SourceHandler initiates a 'source' call, so the handler is supposed to send a stream of stuff to the peer.
type SourceHandler interface {
	HandleSource(context.Context, *cmuxrpc.Request, *cmuxrpc.ByteSink) error
}

type sourceStub struct {
	h SourceHandler
}

func (hm sourceStub) HandleCall(ctx context.Context, req *cmuxrpc.Request) {
	// TODO: check call type

	w, err := req.ResponseSink()
	if err != nil {
		req.CloseWithError(err)
		return
	}

	err = hm.h.HandleSource(ctx, req, w)
	if err != nil {
		req.CloseWithError(err)
		return
	}
}

func (hm sourceStub) HandleConnect(ctx context.Context, edp cmuxrpc.Endpoint) {
	if ch, ok := hm.h.(cmuxrpc.ConnectHandler); ok {
		ch.HandleConnect(ctx, edp)
	}
}
