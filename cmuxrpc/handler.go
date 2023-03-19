package cmuxrpc

import (
	"context"
	"fmt"
)

// https://github.com/maxbrunsfeld/counterfeiter

//go:generate counterfeiter -o fakehandler.go . Handler

// Handler 允许处理连接。call时调用CallHandler，建立连接时调用ConnectHandler
type Handler interface {
	// Handled returns true if the method is handled by the handler
	Handled(Method) bool
	CallHandler
	ConnectHandler
}

type CallHandler interface {
	HandleCall(ctx context.Context, req *Request)
}

type ConnectHandler interface {
	HandleConnect(ctx context.Context, edp Endpoint)
}

type HandlerWrapper func(Handler) Handler

func ApplyHandlerWrappers(h Handler, hws ...HandlerWrapper) Handler {
	for _, hw := range hws {
		h = hw(h)
	}
	return h
}

// HandlerMux来实现Handler接口和重写其三个方法
type HandlerMux struct {
	handlers map[string]Handler
}

// Handled (Handler)
func (hm *HandlerMux) Handled(m Method) bool {
	for _, h := range hm.handlers {
		if h.Handled(m) {
			return true
		}
	}
	return false
}

// HandleCall
func (hm *HandlerMux) HandleCall(ctx context.Context, req *Request) {
	for i := len(req.Method); i > 0; i-- {
		m := req.Method[:i]
		h, ok := hm.handlers[m.String()]
		if ok {
			h.HandleCall(ctx, req)
			return
		}
	}

	req.CloseWithError(fmt.Errorf("no such method: %s", req.Method))
}

// HandleConnect
func (hm *HandlerMux) HandleConnect(ctx context.Context, edp Endpoint) {
	for _, h := range hm.handlers {
		go h.HandleConnect(ctx, edp)
	}
}

var _ Handler = (*HandlerMux)(nil)

// Register add h to handlers
func (hm *HandlerMux) Register(m Method, h Handler) {
	if hm.handlers == nil {
		hm.handlers = make(map[string]Handler)
	}

	hm.handlers[m.String()] = h
}
