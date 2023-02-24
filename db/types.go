package db

import (
	"errors"

	refs "github.com/ssbc/go-ssb-refs"
)

var ErrNotFound = errors.New("db: object not found")

// Alias
type Alias struct {
	ID        int64
	Name      string
	Feed      refs.FeedRef
	Signature []byte
}
