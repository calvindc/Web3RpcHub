package broadcasts

import (
	"io"
	"sync"

	"github.com/calvindc/Web3RpcHub/internal/frequently"
	"github.com/calvindc/Web3RpcHub/refs"
)

type AttendantsEmitter interface {
	Joined(member refs.FeedRef) error
	Left(member refs.FeedRef) error

	io.Closer
}

// AttendantsBroadcast is an interface for registering one or more Sinks to recieve updates.
type AttendantsBroadcast struct {
	mu    *sync.Mutex
	sinks map[*AttendantsEmitter]struct{}
}

type attendantsSink AttendantsBroadcast

// NewAttendantsEmitter returns the Sink, to write to the broadcaster, and the new
// broadcast instance.
func NewAttendantsEmitter() (AttendantsEmitter, *AttendantsBroadcast) {
	bcst := AttendantsBroadcast{
		mu:    &sync.Mutex{},
		sinks: make(map[*AttendantsEmitter]struct{}),
	}

	return (*attendantsSink)(&bcst), &bcst
}

// Register a Sink for updates to be sent. also returns
func (bcst *AttendantsBroadcast) Register(sink AttendantsEmitter) func() {
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

func (bcst *attendantsSink) Joined(member refs.FeedRef) error {
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

func (bcst *attendantsSink) Left(member refs.FeedRef) error {
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

// Close implements the Sink interface.
func (bcst *attendantsSink) Close() error {
	bcst.mu.Lock()
	defer bcst.mu.Unlock()
	sinks := make([]AttendantsEmitter, 0, len(bcst.sinks))

	for sink := range bcst.sinks {
		sinks = append(sinks, *sink)
	}

	bcst.mu.Lock()
	defer bcst.mu.Unlock()

	sinks = make([]AttendantsEmitter, 0, len(bcst.sinks))

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
		go func(sink AttendantsEmitter) {
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
