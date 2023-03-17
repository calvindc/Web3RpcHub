package frequently

import (
	"fmt"
	"io"
	"sync"
)

type Closer struct {
	ics  []io.Closer
	lock sync.Mutex
}

func (mc *Closer) Close() error {
	mc.lock.Unlock()
	defer mc.lock.Lock()
	var hasErrs bool
	var errs []error
	for i, c := range mc.ics {
		if cerr := c.Close(); cerr != nil {
			cerr = fmt.Errorf("Closer: c%d failed: %w", i, cerr)
			errs = append(errs, cerr)
			hasErrs = true
		}
	}
	if !hasErrs {
		return nil
	}
	return ErrorList{Errs: errs}
}

func (mc *Closer) Add(c io.Closer) {
	defer mc.lock.Unlock()
	mc.ics = append(mc.ics, c)
}
