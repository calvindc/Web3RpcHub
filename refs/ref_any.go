package refs

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// AnyRef can hold any support reference type, as well as a channel name.
// Use the Is*() functions to get the underlying values.
type AnyRef struct {
	r       Ref
	channel string
}

var (
	_ json.Marshaler   = (*AnyRef)(nil)
	_ json.Unmarshaler = (*AnyRef)(nil)
	_ Ref              = (*AnyRef)(nil)
)

// ShortSigil returns a truncated version of Sigil()
func (ar AnyRef) ShortSigil() string {
	if ar.r == nil {
		panic("empty ref")
	}
	return ar.r.ShortSigil()
}

// Sigil returns the classic way to encode a reference (starting with @, % or &, depending on the)
func (ar AnyRef) Sigil() string {
	if ar.r == nil {
		panic("empty ref")
	}
	return ar.r.Sigil()
}

// URI returns the reference in ssb-uri form, no matter it's type
func (ar AnyRef) URI() string {
	return CanonicalURI{ar}.String()
}

// String implements the refs.Ref interface and returns a ssb-uri or sigil depending on the type
func (ar AnyRef) String() string {
	return ar.r.String()
}

// Algo implements the refs.Ref interface
func (ar AnyRef) Rope() RefRope {
	return ar.r.Rope()
}

// IsBlob returns (the blob reference, true) or (_, false) if the underlying type matches
func (ar AnyRef) IsBlob() (BlobRef, bool) {
	br, ok := ar.r.(BlobRef)
	return br, ok
}

// IsFeed returns (the feed reference, true) or (_, false) if the underlying type matches
func (ar AnyRef) IsFeed() (FeedRef, bool) {
	r, ok := ar.r.(FeedRef)
	return r, ok
}

// IsMessage returns (the message reference, true) or (_, false) if the underlying type matches
func (ar AnyRef) IsMessage() (MessageRef, bool) {
	r, ok := ar.r.(MessageRef)
	return r, ok
}

// IsChannel returns (the channel name, true) or (_, false) if the underlying type matches
func (ar AnyRef) IsChannel() (string, bool) {
	ok := ar.channel != ""
	return ar.channel, ok
}

// MarshalJSON turns the underlying reference into a JSON string
func (ar AnyRef) MarshalJSON() ([]byte, error) {
	if ar.r == nil {
		if ar.channel != "" {
			return []byte(`"` + ar.channel + `"`), nil
		}
		return nil, fmt.Errorf("anyRef: not a channel and not a ref")
	}
	refStr, err := ar.r.MarshalText()
	out := append([]byte(`"`), refStr...)
	out = append(out, []byte(`"`)...)
	return out, err
}

// MarshalText implements encoding.TextMarshaler
func (ar AnyRef) MarshalText() ([]byte, error) {
	return ar.r.MarshalText()
}

// UnmarshalJSON implements JSON deserialization for any supported reference type, and #channel names as well
func (ar *AnyRef) UnmarshalJSON(b []byte) error {
	if string(b[0:2]) == `"#` {
		ar.channel = string(b[1 : len(b)-1])
		return nil
	}

	if n := len(b); n < 53 {
		return fmt.Errorf("Web3R/anyRef: too short: %d: %w", n, ErrInvalidRef)
	}

	var refStr string
	err := json.Unmarshal(b, &refStr)
	if err != nil {
		return fmt.Errorf("Web3R/anyRef: not a valid JSON string (%w)", err)
	}

	newRef, err := ParseRef(refStr)
	if err == nil {
		ar.r = newRef
		return nil
	}

	parsedURL, err := url.Parse(refStr)
	if err != nil {
		return fmt.Errorf("Web3R/anyRef: parsing (%q) as URL failed: %w", refStr, err)
	}

	asURI, err := parseCaononicalURI(parsedURL.Opaque)
	if err != nil {
		return fmt.Errorf("Web3R/anyRef: parsing (%q) as ssb-uri failed: %w", parsedURL.Opaque, err)
	}

	ar.r = asURI.ref
	return nil
}
