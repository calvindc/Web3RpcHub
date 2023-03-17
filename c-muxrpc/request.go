package c_muxrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/calvindc/Web3RpcHub/c-muxrpc/codec"
)

type RequestEncoding uint

const (
	TypeBinary RequestEncoding = iota
	TypeString
	TypeJSON
)

// Isvalid check the request code
func (rt RequestEncoding) IsValid() bool {
	if rt < 0 {
		return false
	}
	if rt > TypeJSON {
		return false
	}
	return true
}

func (rt RequestEncoding) asCodecFlag() (codec.Flag, error) {
	if !rt.IsValid() {
		return 0, fmt.Errorf("muxrpc: invalid request encoding %d", rt)
	}
	switch rt {
	case TypeBinary:
		return 0, nil
	case TypeString:
		return codec.FlagString, nil
	case TypeJSON:
		return codec.FlagJSON, nil
	default:
		return 0, fmt.Errorf("muxrpc: invalid request encoding %d", rt)
	}
}

// Method defines the name of the endpoint.
type Method []string

// UnmarshalJSON decodes the
func (m *Method) UnmarshalJSON(d []byte) error {
	var newM []string

	err := json.Unmarshal(d, &newM)
	if err != nil {
		var meth string
		err := json.Unmarshal(d, &meth)
		if err != nil {
			return fmt.Errorf("muxrpc/method: error decoding packet: %w", err)
		}
		newM = Method{meth}
	}
	*m = newM
	return nil
}

func (m Method) String() string {
	return strings.Join(m, ".")
}

// CallType is the type of a call
type CallType string

// Request assembles the state of an RPC call
type Request struct {
	// Stream is a legacy adapter for luigi-powered streams
	Stream Stream `json:"-"`

	// Method is the name of the called function
	Method Method `json:"name"`

	// Args contains the call arguments
	RawArgs json.RawMessage `json:"args"`

	// Type is the type of the call, i.e. async, sink, source or duplex
	Type CallType `json:"type"`

	// luigi-less iterators
	sink   *ByteSink
	source *ByteSource

	// same as packet.Req - the numerical identifier for the stream
	id int32

	// used to stop producing more data on this request
	// the calling sight might tell us they had enough of this stream
	abort context.CancelFunc

	remoteAddr net.Addr
	endpoint   *rpc
}

func (req Request) Endpoint() Endpoint {
	return req.endpoint
}
