package refs

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
)

// BlobRef defines a static binary attachment reference, identified it's hash.
type BlobRef struct {
	hash [32]byte
	rope RefRope
}

// NewBlobRefFromBytes allows to create a blob reference from raw bytes
func NewBlobRefFromBytes(b []byte, algo RefRope) (BlobRef, error) {
	ref := BlobRef{
		rope: algo,
	}
	n := copy(ref.hash[:], b)
	if n != 32 {
		return BlobRef{}, ErrRefLen{rope: ref.rope, n: n}
	}
	return ref, nil
}

// Rope implements the refs.Ref interface
func (br BlobRef) Rope() RefRope {
	return br.rope
}

// CopyHashTo copies the internal hash data somewhere else
// the target needs to have enough space, otherwise an error is returned.
func (br BlobRef) CopyHashTo(b []byte) error {
	if n := len(b); n != len(br.hash) {
		return ErrRefLen{rope: "target", n: n}
	}
	copy(b, br.hash[:])
	return nil
}

// Sigil returns the BlobRef with thealgo sigil &, it's base64 encoded hash and the used algo (currently only sha256)
func (br BlobRef) Sigil() string {
	return fmt.Sprintf("&%s.%s", base64.StdEncoding.EncodeToString(br.hash[:]), br.rope)
}

// ShortSigil returns a truncated version of Sigil()
func (br BlobRef) ShortSigil() string {
	return fmt.Sprintf("<&%s.%s>", base64.StdEncoding.EncodeToString(br.hash[:3]), br.rope)
}

// URI returns the reference in hub-uri form, no matter it's type
func (br BlobRef) URI() string {
	return CanonicalURI{br}.String()
}

// String implements the refs.Ref interface and returns a hub-uri or sigil depending on the type
func (br BlobRef) String() string {
	if br.rope == RefBlobWEB3R {
		return br.Sigil()
	}
	return br.URI()
}

var emptyBlobRef = BlobRef{}

// ParseBlobRef uses ParseRef and checks that it returns a *BlobRef
func ParseBlobRef(str string) (BlobRef, error) {
	if len(str) == 0 {
		return emptyBlobRef, fmt.Errorf("Web3R: blob reference empty")
	}

	split := strings.Split(str[1:], ".")
	if len(split) < 2 {
		return emptyBlobRef, ErrInvalidRef
	}

	raw, err := base64.StdEncoding.DecodeString(split[0])
	if err != nil {
		return emptyBlobRef, fmt.Errorf("blob reference: couldn't parse %q: %s: %w", str, err, ErrInvalidHash)
	}

	if str[0] != '&' {
		return emptyBlobRef, ErrInvalidRefType
	}

	var rope RefRope
	switch RefRope(split[1]) {
	case RefBlobWEB3R:
		rope = RefBlobWEB3R
	default:
		return emptyBlobRef, ErrInvalidRefAlgo
	}
	if n := len(raw); n != 32 {
		return emptyBlobRef, newHashLenError(n)
	}

	newBlob := BlobRef{rope: rope}
	copy(newBlob.hash[:], raw)
	return newBlob, nil
}

// Equal compares two references with each other
func (br BlobRef) Equal(other BlobRef) bool {
	if br.rope != other.rope {
		return false
	}
	return bytes.Equal(br.hash[:], other.hash[:])
}

// IsValid checks if the RefAlgo is known and the length of the data is as expected
func (br BlobRef) IsValid() error {
	if br.rope != RefBlobWEB3R {
		return fmt.Errorf("unknown hash algorithm %q", br.rope)
	}
	if len(br.hash) != 32 {
		return fmt.Errorf("expected hash length 32, got %v", len(br.hash))
	}
	return nil
}

// MarshalText encodes the BlobRef using String()
func (br BlobRef) MarshalText() ([]byte, error) {
	return []byte(br.String()), nil
}

// UnmarshalText uses ParseBlobRef
func (br *BlobRef) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*br = BlobRef{}
		return nil
	}
	newBR, err := ParseBlobRef(string(text))
	if err != nil {
		return fmt.Errorf(" BlobRef/UnmarshalText failed: %w", err)
	}
	*br = newBR
	return nil
}
