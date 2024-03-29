package typemux

import (
	"context"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
)

// DuplexHandler initiates a 'duplex' call. The handler receives data from the peer through the passed source
type DuplexHandler interface {
	HandleDuplex(context.Context, *cmuxrpc.Request, *cmuxrpc.ByteSource, *cmuxrpc.ByteSink) error
}

// DuplexFunc is a utility to fulfill DuplexHandler just as a function, not a type with the named method
type DuplexFunc func(context.Context, *cmuxrpc.Request, *cmuxrpc.ByteSource, *cmuxrpc.ByteSink) error

// HandleDuplex implements the duplex handler for the function
func (sf DuplexFunc) HandleDuplex(ctx context.Context, r *cmuxrpc.Request, src *cmuxrpc.ByteSource, snk *cmuxrpc.ByteSink) error {
	return sf(ctx, r, src, snk)
}

var _ DuplexHandler = (*DuplexFunc)(nil)

type duplexStub struct {
	h DuplexHandler
}

func (hm duplexStub) HandleCall(ctx context.Context, req *cmuxrpc.Request) {
	// TODO: check call type

	r, err := req.ResponseSource()
	if err != nil {
		req.CloseWithError(err)
		return
	}

	w, err := req.ResponseSink()
	if err != nil {
		req.CloseWithError(err)
		return
	}

	err = hm.h.HandleDuplex(ctx, req, r, w)
	if err != nil {
		req.CloseWithError(err)
		return
	}
}

func (hm duplexStub) HandleConnect(ctx context.Context, edp cmuxrpc.Endpoint) {
	if ch, ok := hm.h.(cmuxrpc.ConnectHandler); ok {
		ch.HandleConnect(ctx, edp)
	}
}
