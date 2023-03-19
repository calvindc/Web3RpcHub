package typemux

import (
	"context"
	"fmt"

	"github.com/calvindc/Web3RpcHub/cmuxrpc"
	"github.com/go-kit/log"
)

// the handlers passed to typemux are checked bu this muxer and dont need the Handled() function
type handler interface {
	cmuxrpc.CallHandler
	cmuxrpc.ConnectHandler
}

type HandlerMux struct {
	logger log.Logger

	handlers map[string]handler
}

var _ cmuxrpc.Handler = (*HandlerMux)(nil)

func New(log log.Logger) HandlerMux {
	return HandlerMux{
		handlers: make(map[string]handler),
		logger:   log,
	}
}

func (hm *HandlerMux) Handled(m cmuxrpc.Method) bool {
	_, has := hm.handlers[m.String()]
	return has
}

func (hm *HandlerMux) HandleCall(ctx context.Context, req *cmuxrpc.Request) {
	for i := len(req.Method); i > 0; i-- {
		m := req.Method[:i]
		h, ok := hm.handlers[m.String()]
		if ok {
			h.HandleCall(ctx, req)
			return
		}
	}
	req.CloseWithError(fmt.Errorf("no such command: %v", req.Method))
}

func (hm *HandlerMux) HandleConnect(ctx context.Context, edp cmuxrpc.Endpoint) {
	for _, h := range hm.handlers {
		go h.HandleConnect(ctx, edp)
	}
}

// RegisterAsync registers a 'async' call for name method
func (hm *HandlerMux) RegisterAsync(m cmuxrpc.Method, h AsyncHandler) {
	hm.handlers[m.String()] = asyncStub{
		logger: hm.logger,
		h:      h,
	}
}

// RegisterSource registers a 'source' call for name method
func (hm *HandlerMux) RegisterSource(m cmuxrpc.Method, h SourceHandler) {
	hm.handlers[m.String()] = sourceStub{
		// logger: hm.logger,
		h: h,
	}
}

// RegisterSink registers a 'sink' call for name method
func (hm *HandlerMux) RegisterSink(m cmuxrpc.Method, h SinkHandler) {
	hm.handlers[m.String()] = sinkStub{
		// logger: hm.logger,
		h: h,
	}
}

// RegisterDuplex registers a 'sink' call for name method
func (hm *HandlerMux) RegisterDuplex(m cmuxrpc.Method, h DuplexHandler) {
	hm.handlers[m.String()] = duplexStub{
		// logger: hm.logger,
		h: h,
	}
}
