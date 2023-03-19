package typemux

import (
	"context"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
)

var _ SinkHandler = (*SinkFunc)(nil)

// SinkFunc is a utility to fulfill SinkHandler just as a function, not a type with the named method
type SinkFunc func(context.Context, *cmuxrpc.Request, *cmuxrpc.ByteSource) error

// HandleSink implements the sink handler for the function
func (sf SinkFunc) HandleSink(ctx context.Context, r *cmuxrpc.Request, src *cmuxrpc.ByteSource) error {
	return sf(ctx, r, src)
}

// SinkHandler initiates a 'sink' call. The handler receives data from the peer through the passed source
type SinkHandler interface {
	HandleSink(context.Context, *cmuxrpc.Request, *cmuxrpc.ByteSource) error
}

type sinkStub struct {
	h SinkHandler
}

func (hm sinkStub) HandleCall(ctx context.Context, req *cmuxrpc.Request) {
	// TODO: check call type

	src, err := req.ResponseSource()
	if err != nil {
		req.CloseWithError(err)
		return
	}

	err = hm.h.HandleSink(ctx, req, src)
	if err != nil {
		req.CloseWithError(err)
		return
	}
}

func (hm sinkStub) HandleConnect(ctx context.Context, edp cmuxrpc.Endpoint) {
	if ch, ok := hm.h.(cmuxrpc.ConnectHandler); ok {
		ch.HandleConnect(ctx, edp)
	}
}
