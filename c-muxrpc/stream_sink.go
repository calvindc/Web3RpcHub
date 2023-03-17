package c_muxrpc

import (
	"context"
	"io"
	"sync"

	"github.com/calvindc/Web3RpcHub/c-muxrpc/codec"
)

type ByteSinker interface {
	io.WriteCloser

	CloseWithError(error) error

	SetEncoding(re RequestEncoding)
}

var _ ByteSinker = (*ByteSink)(nil)

type ByteSink struct {
	w         *codec.Writer
	closedMu  sync.Mutex
	closed    error
	streamCtx context.Context
	pkt       codec.Packet
}

func newByteSink(ctx context.Context, w *codec.Writer) *ByteSink {
	return &ByteSink{
		w:         w,
		streamCtx: ctx,
		pkt:       codec.Packet{},
	}
}
