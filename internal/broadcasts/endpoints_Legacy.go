package broadcasts

import (
	"io"
	"sync"

	"github.com/calvindc/Web3RpcHub/internal/errorslist"
)

//endpointsSink 接入节点接收管理池
type endpointsLegacySink EndpointsBroadcastLegacy

// EndpointsBroadcast
type EndpointsBroadcastLegacy struct {
	linkedTime   int64
	unlinkedTime int64 //超时主动写入
	mu           *sync.Mutex
	sinks        map[*EndpointsLegacyEmitter]struct{}
}

type EndpointsLegacyEmitter interface {
	Update(members []string) error
	io.Closer
}

// NewEndpointsEmitter 返回Sink，写入广播发送器和新的广播实例
func NewEndpointsEmitter() (EndpointsLegacyEmitter, *EndpointsBroadcastLegacy) {
	bcst := EndpointsBroadcastLegacy{
		mu:    &sync.Mutex{},
		sinks: make(map[*EndpointsLegacyEmitter]struct{}),
	}

	return (*endpointsLegacySink)(&bcst), &bcst
}

func (bcst *EndpointsBroadcastLegacy) Register(sink EndpointsLegacyEmitter) func() {
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

func (bcst *endpointsLegacySink) Update(members []string) error {
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

func (bcst *endpointsLegacySink) Close() error {
	var sinks []EndpointsLegacyEmitter

	bcst.mu.Lock()
	defer bcst.mu.Unlock()

	sinks = make([]EndpointsLegacyEmitter, 0, len(bcst.sinks))

	for sink := range bcst.sinks {
		sinks = append(sinks, *sink)
	}

	var (
		wg   sync.WaitGroup
		errl errorslist.ErrorList
	)

	wg.Add(len(sinks))
	for _, sink_ := range sinks {
		go func(sink EndpointsLegacyEmitter) {
			defer wg.Done()

			err := sink.Close()
			if err != nil {
				errl.Errs = append(errl.Errs, err)
				return
			}
		}(sink_)
	}
	wg.Wait()

	if len(errl.Errs) == 0 {
		return nil
	}
	return errl
}
