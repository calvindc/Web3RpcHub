package refs

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Common errors for invalid references
var (
	ErrInvalidRef        = errors.New("web3r: Invalid Ref")
	ErrInvalidRefType    = errors.New("web3r: Invalid Ref Type")
	ErrInvalidRefAlgo    = errors.New("web3r: Invalid Ref Algo")
	ErrInvalidSig        = errors.New("web3r: Invalid Signature")
	ErrInvalidHash       = errors.New("web3r: Invalid Hash")
	ErrUnuspportedFormat = errors.New("web3r: unsupported format")
)

// ErrRefLen is returned when a parsed reference was too short.
type ErrRefLen struct {
	rope RefRope
	n    int
}

func (e ErrRefLen) Error() string {
	return fmt.Sprintf("web3r: Invalid reference len for %s: %d", e.rope, e.n)
}

// NewFeedRefLenError returns a new ErrRefLen error for a feed
func newFeedRefLenError(n int) error {
	return ErrRefLen{rope: RefFeedWEB3R, n: n}
}

func newHashLenError(n int) error {
	return ErrRefLen{rope: RefMessageWEB3R, n: n}
}

// IsMessageUnusable checks if an error is ErrWrongType, ErrMalfromedMsg or *json.SyntaxError
func IsMessageUnusable(err error) bool {
	var errWt ErrWrongType
	if errors.As(err, &errWt) {
		return true
	}

	var errMalMsg ErrMalfromedMsg
	if errors.As(err, &errMalMsg) {
		return true
	}

	return errors.Is(err, &json.SyntaxError{})
}

// ErrMalfromedMsg is returned if a message has invalid values
type ErrMalfromedMsg struct {
	reason string
	m      map[string]interface{}
}

func (emm ErrMalfromedMsg) Error() string {
	s := "ErrMalfromedMsg: " + emm.reason
	if emm.m != nil {
		s += fmt.Sprintf(" %+v", emm.m)
	}
	return s
}

// ErrWrongType is returned if a certain type:value was expected on a message.
type ErrWrongType struct {
	has, want string
}

func (ewt ErrWrongType) Error() string {
	return fmt.Sprintf("ErrWrongType: want: %s has: %s", ewt.want, ewt.has)
}
