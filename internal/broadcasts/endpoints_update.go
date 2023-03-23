package broadcasts

import (
	"io"

	"sync"

	"github.com/calvindc/Web3RpcHub/internal/frequently"
)

type EndpointsEmitter interface {
	Update(members []string) error
	io.Closer
}

// EndpointsBroadcast is an interface for registering one or more Sinks to recieve updates.
type EndpointsBroadcast struct {
	mu    *sync.Mutex
	sinks map[*EndpointsEmitter]struct{}
}

type endpointsSink EndpointsBroadcast

// NewEndpointsEmitter returns the Sink, to write to the broadcaster, and the new
// broadcast instance.
func NewEndpointsEmitter() (EndpointsEmitter, *EndpointsBroadcast) {
	bcst := EndpointsBroadcast{
		mu:    &sync.Mutex{},
		sinks: make(map[*EndpointsEmitter]struct{}),
	}

	return (*endpointsSink)(&bcst), &bcst
}

// Register a Sink for updates to be sent. also returns
func (bcst *EndpointsBroadcast) Register(sink EndpointsEmitter) func() {
	bcst.mu.Lock()
	defer bcst.mu.Unlock()
	bcst.sinks[&sink] = struct{}{}

	return func() {
		bcst.mu.Lock()
		defer bcst.mu.Unlock()
		delete(bcst.sinks, &sink)
		sink.Close()
	}
}

// Pour implements the Sink interface.
func (bcst *endpointsSink) Update(members []string) error {

	bcst.mu.Lock()
	for s := range bcst.sinks {
		err := (*s).Update(members)
		if err != nil {
			delete(bcst.sinks, s)
		}
	}
	bcst.mu.Unlock()

	return nil
}

// Close implements the Sink interface.
func (bcst *endpointsSink) Close() error {
	var sinks []EndpointsEmitter

	bcst.mu.Lock()
	defer bcst.mu.Unlock()

	sinks = make([]EndpointsEmitter, 0, len(bcst.sinks))

	for sink := range bcst.sinks {
		sinks = append(sinks, *sink)
	}

	var (
		wg sync.WaitGroup
		me frequently.ErrorList
	)

	// might be fine without the waitgroup and concurrency

	wg.Add(len(sinks))
	for _, sink_ := range sinks {
		go func(sink EndpointsEmitter) {
			defer wg.Done()

			err := sink.Close()
			if err != nil {
				me.Errs = append(me.Errs, err)
				return
			}
		}(sink_)
	}
	wg.Wait()

	if len(me.Errs) == 0 {
		return nil
	}

	return me
}
