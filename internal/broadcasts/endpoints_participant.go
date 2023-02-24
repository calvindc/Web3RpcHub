package broadcasts

import (
	"io"

	"sync"

	"github.com/calvindc/Web3RpcHub/internal/errorslist"
	refs "github.com/ssbc/go-ssb-refs"
)

type participantsSink EndpointsBroadcastParticipants

type EndpointsBroadcastParticipants struct {
	linkedTime   int64
	unlinkedTime int64
	mu           *sync.Mutex
	sinks        map[*EndpointsBroadcastParticipantsEmitter]struct{}
}

type EndpointsBroadcastParticipantsEmitter interface {
	Joined(member refs.FeedRef) error
	Left(member refs.FeedRef) error

	io.Closer
}

// NewEndpointsEmitter 返回Sink，写入广播发送器和新的广播实例
func NewParticipantsEmitter() (EndpointsBroadcastParticipantsEmitter, *EndpointsBroadcastParticipants) {
	bcst := EndpointsBroadcastParticipants{
		mu:    &sync.Mutex{},
		sinks: make(map[*EndpointsBroadcastParticipantsEmitter]struct{}),
	}

	return (*participantsSink)(&bcst), &bcst
}

func (bcst *EndpointsBroadcastParticipants) Register(sink EndpointsBroadcastParticipantsEmitter) func() {
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

func (bcst *participantsSink) Joined(member refs.FeedRef) error {
	bcst.mu.Lock()
	for s := range bcst.sinks {
		err := (*s).Joined(member)
		if err != nil {
			delete(bcst.sinks, s)
		}
	}
	bcst.mu.Unlock()

	return nil
}

func (bcst *participantsSink) Left(member refs.FeedRef) error {
	bcst.mu.Lock()
	for s := range bcst.sinks {
		err := (*s).Left(member)
		if err != nil {
			delete(bcst.sinks, s)
		}
	}
	bcst.mu.Unlock()

	return nil
}

func (bcst *participantsSink) Close() error {
	bcst.mu.Lock()
	defer bcst.mu.Unlock()
	sinks := make([]EndpointsBroadcastParticipantsEmitter, 0, len(bcst.sinks))

	for sink := range bcst.sinks {
		sinks = append(sinks, *sink)
	}

	bcst.mu.Lock()
	defer bcst.mu.Unlock()

	sinks = make([]EndpointsBroadcastParticipantsEmitter, 0, len(bcst.sinks))

	for sink := range bcst.sinks {
		sinks = append(sinks, *sink)
	}

	var (
		wg sync.WaitGroup
		me errorslist.ErrorList
	)

	wg.Add(len(sinks))
	for _, sink_ := range sinks {
		go func(sink ParticipantsEmitter) {
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
