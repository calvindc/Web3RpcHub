package refs

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"fmt"
	"strings"
)

// MessageRef defines the content addressed version of a hub message, identified it's hash.
type MessageRef struct {
	hash [32]byte
	rope RefRope
}

// NewMessageRefFromBytes allows to create a message reference from raw bytes
func NewMessageRefFromBytes(b []byte, rope RefRope) (MessageRef, error) {
	fr := MessageRef{
		rope: rope,
	}
	n := copy(fr.hash[:], b)
	if n != 32 { // TODO: support references of different lengths.
		return MessageRef{}, ErrRefLen{rope: fr.rope, n: n}
	}
	return fr, nil
}

// Rope implements the refs.Ref interface
func (mr MessageRef) Rope() RefRope {
	return mr.rope
}

// Equal compares two references with each other
func (mr MessageRef) Equal(other MessageRef) bool {
	if mr.rope != other.rope {
		return false
	}

	return bytes.Equal(mr.hash[:], other.hash[:])
}

// CopyHashTo copies the internal hash data somewhere else
// the target needs to have enough space, otherwise an error is returned.
func (mr MessageRef) CopyHashTo(b []byte) error {
	if len(b) != len(mr.hash) {
		return ErrRefLen{rope: mr.rope, n: len(b)}
	}
	copy(b, mr.hash[:])
	return nil
}

// Sigil returns the MessageRef with the sigil %, it's base64 encoded hash and the used algo (currently only sha256)
func (mr MessageRef) Sigil() string {
	return fmt.Sprintf("%%%s.%s", base64.StdEncoding.EncodeToString(mr.hash[:]), mr.rope)
}

// ShortSigil prints a shortend version of Sigil()
func (mr MessageRef) ShortSigil() string {
	return fmt.Sprintf("<%%%s.%s>", base64.StdEncoding.EncodeToString(mr.hash[:3]), mr.rope)
}

// URI returns the reference in hub-uri form, no matter it's type
func (mr MessageRef) URI() string {
	return CanonicalURI{mr}.String()
}

// String implements the refs.Ref interface and returns a hub-uri or sigil depending on the type
func (mr MessageRef) String() string {
	if mr.rope == RefMessageWEB3R || mr.rope == RefCloakedGroup {
		return mr.Sigil()
	}
	return mr.URI()
}

var (
	_ encoding.TextMarshaler   = (*MessageRef)(nil)
	_ encoding.TextUnmarshaler = (*MessageRef)(nil)
)

// MarshalText implements encoding.TextMarshaler
func (mr MessageRef) MarshalText() ([]byte, error) {
	if mr.rope == RefMessageWEB3R || mr.rope == RefCloakedGroup {
		return []byte(mr.Sigil()), nil
	}
	asURI := CanonicalURI{mr}
	return []byte(asURI.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (mr *MessageRef) UnmarshalText(input []byte) error {
	txt := string(input)

	newRef, err := ParseMessageRef(txt)
	if err != nil {
		return err
	}

	*mr = newRef
	return nil
}

// ParseMessageRef returns a message ref from a string, if it's valid
func ParseMessageRef(str string) (MessageRef, error) {
	if len(str) == 0 {
		return emptyMsgRef, fmt.Errorf("Web3R: msgRef empty")
	}

	split := strings.Split(str[1:], ".")
	if len(split) < 2 {
		asURI, err := parseCaononicalURI(str)
		if err != nil {
			return emptyMsgRef, err
		}

		newRef, ok := asURI.Message()
		if ok {
			return newRef, nil
		}
		return emptyMsgRef, ErrInvalidRef
	}

	raw, err := base64.StdEncoding.DecodeString(split[0])
	if err != nil {
		return emptyMsgRef, fmt.Errorf("msgRef: couldn't parse %q: %s: %w", str, err, ErrInvalidHash)
	}

	if str[0] != '%' {
		return emptyMsgRef, ErrInvalidRefType
	}

	var rope RefRope
	switch RefRope(split[1]) {
	case RefMessageWEB3R:
		rope = RefMessageWEB3R
	case RefMessageGabby:
		rope = RefMessageGabby
	case RefCloakedGroup:
		rope = RefCloakedGroup
	default:
		return emptyMsgRef, ErrInvalidRefAlgo
	}
	if n := len(raw); n != 32 {
		return emptyMsgRef, newHashLenError(n)
	}
	newMsg := MessageRef{rope: rope}
	copy(newMsg.hash[:], raw)
	return newMsg, nil
}

// MessageRefs holds a slice of multiple message references
type MessageRefs []MessageRef

// String turns a slice of references by String()ifiyng and joining them with a comma
func (mr *MessageRefs) String() string {
	var s []string
	for _, r := range *mr {
		s = append(s, r.String())
	}
	return strings.Join(s, ", ")
}

// UnmarshalJSON implements JSON deserialization for a list of message references.
// It also supports empty and `null` as well as a single refrence as a string ("%foo" and ["%foo"] are both returned as the right-hand case)
func (mr *MessageRefs) UnmarshalJSON(text []byte) error {
	if len(text) == 0 {
		*mr = nil
		return nil
	}

	if bytes.Equal([]byte("[]"), text) || bytes.Equal([]byte("null"), text) {
		*mr = nil
		return nil
	}

	if bytes.HasPrefix(text, []byte("[")) && bytes.HasSuffix(text, []byte("]")) {

		elems := bytes.Split(text[1:len(text)-1], []byte(","))
		newArr := make([]MessageRef, len(elems))

		for i, e := range elems {
			var err error
			r := strings.TrimSpace(string(e))
			r = r[1 : len(r)-1] // remove quotes
			newArr[i], err = ParseMessageRef(r)
			if err != nil {
				return fmt.Errorf("messageRefs %d unmarshal failed: %w", i, err)
			}
		}

		*mr = newArr

	} else {
		newArr := make([]MessageRef, 1)

		var err error
		newArr[0], err = ParseMessageRef(string(text[1 : len(text)-1]))
		if err != nil {
			return fmt.Errorf("messageRefs single unmarshal failed: %w", err)
		}

		*mr = newArr
	}
	return nil
}
