package codec

import (
	"fmt"
	"math"
	"strings"
)

const maxBufferSize = math.MaxUint32

// Flag is the first byte of the Header
type Flag byte

type Body []byte

type Packet struct {
	Flag Flag
	Req  int32
	Body Body
}

// Flag bitmasks
const (
	FlagString Flag = 1 << iota // type
	FlagJSON                    // bits
	FlagEndErr
	FlagStream
)

// Header is the wire representation of a packet header
type Header struct {
	Flag Flag
	Len  uint32
	Req  int32
}

func (b Body) String() string {
	return fmt.Sprintf("%s", []byte(b))
}

func (f Flag) Get(g Flag) bool {
	return f&g == g
}

func (f Flag) Set(g Flag) Flag {
	return f | g
}

func (f Flag) Clear(g Flag) Flag {
	return f & ^g
}

func (f Flag) String() string {
	var flags []string

	if f.Get(FlagString) {
		flags = append(flags, "FlagString")
	}
	if f.Get(FlagJSON) {
		flags = append(flags, "FlagJSON")
	}
	if f.Get(FlagStream) {
		flags = append(flags, "FlagStream")
	}
	if f.Get(FlagEndErr) {
		flags = append(flags, "FlagEndErr")
	}

	return "{" + strings.Join(flags, ", ") + "}"
}
