package cmuxrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
)

// ErrSessionTerminated 如果调用Terminate()或连接断开，返回ErrSessionTerminated
var ErrSessionTerminated = errors.New("cmuxrpc: session terminated")

var errSinkClosed = errors.New("cmuxrpc: pour to closed sink")

type ErrNoSuchMethod struct {
	Method Method
}

func (e ErrNoSuchMethod) Error() string {
	return fmt.Sprintf("cmuxrpc: no such command: %s", e.Method)
}

// CallError is returned when a call fails
type CallError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

func (e CallError) Error() string {
	return fmt.Sprintf("cmuxrpc CallError: %s - %s", e.Name, e.Message)
}

func parseError(data []byte) (*CallError, error) {
	var e CallError

	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil, fmt.Errorf("cmuxrpc: failed to unmarshal error packet: %w", err)
	}

	// There are also TypeErrors and numerous other things we might get from this..
	// if e.Name != "Error" {
	// 	return nil, fmt.Errorf(`name is not "Error" but %q`, e.Name)
	// }

	return &e, nil
}

type ErrWrongStreamType struct{ ct CallType }

func (wst ErrWrongStreamType) Error() string {
	return fmt.Sprintf("cmuxrpc: wrong stream type: %s", wst.ct)
}

// IsSinkClosed should be moved to luigi to gether with the error
func IsSinkClosed(err error) bool {
	if err == nil {
		return false
	}

	if err == errSinkClosed {
		return true
	}

	if err == ErrSessionTerminated {
		return true
	}

	if isAlreadyClosed(err) {
		return true
	}

	var ce *CallError
	if errors.As(err, &ce) && ce.Message == "unexpected end of parent stream" {
		return true
	}

	return false
}

func isAlreadyClosed(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, io.EOF) || errors.Is(err, os.ErrClosed) || errors.Is(err, io.ErrClosedPipe) {
		return true
	}

	if sysErr, ok := (err).(*os.PathError); ok {
		if sysErr.Err == os.ErrClosed {
			// fmt.Printf("debug: found syscall err: %T) %s\n", err, err)
			return true
		}
	}

	if opErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := opErr.Err.(*os.SyscallError); ok {
			if errNo, ok := syscallErr.Err.(syscall.Errno); ok {
				if errNo == syscall.EPIPE {
					return true
				}
			}
		}
	}
	return false
}
