package refs

import (
	"bytes"
	"crypto/ed25519"
	"encoding"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// FeedRef defines a publickey as ID with a specific algorithm (currently only ed25519)
type FeedRef struct {
	id   [32]byte
	rope RefRope
}

// NewFeedRefFromByteRefs creats a feed reference directly from some bytes
func NewFeedRefFromBytes(b []byte, rope RefRope) (FeedRef, error) {
	fr := FeedRef{
		rope: rope,
	}
	n := copy(fr.id[:], b)
	if n != 32 {
		return FeedRef{}, ErrRefLen{rope: fr.rope, n: n}
	}
	return fr, nil
}

// PubKey returns the crypto/ed25519 public key representation
func (fr FeedRef) PubKey() ed25519.PublicKey {
	return fr.id[:]
}

// Rope implements the refs.Ref interface
func (fr FeedRef) Rope() RefRope {
	return fr.rope
}

// Equal compares two references with each other
func (fr FeedRef) Equal(b FeedRef) bool {
	if fr.rope != b.rope {
		return false
	}
	return bytes.Equal(fr.id[:], b.id[:])
}

// Sigil returns the FeedRef as a string with the sigil @, it's base64 encoded hash and the used algo
func (fr FeedRef) Sigil() string {
	return fmt.Sprintf("@%s.%s", base64.StdEncoding.EncodeToString(fr.id[:]), fr.rope)
}

// ShortSigil returns a truncated version of Sigil()
func (fr FeedRef) ShortSigil() string {
	return fmt.Sprintf("<@%s.%s>", base64.StdEncoding.EncodeToString(fr.id[:3]), fr.rope)
}

// URI returns the reference in uri form, no matter it's type
func (fr FeedRef) URI() string {
	return CanonicalURI{fr}.String()
}

// String implements the refs.Ref interface and returns a uri or sigil depending on the type
func (fr FeedRef) String() string {
	if fr.rope == RefFeedWEB3R {
		return fr.Sigil()
	}
	return fr.URI()
}

var (
	_ encoding.TextMarshaler   = (*FeedRef)(nil)
	_ encoding.TextUnmarshaler = (*FeedRef)(nil)
)

// MarshalText implements encoding.TextMarshaler
func (fr FeedRef) MarshalText() ([]byte, error) {
	if fr.rope == RefFeedWEB3R {
		return []byte(fr.Sigil()), nil
	}
	asURI := CanonicalURI{fr}
	return []byte(asURI.String()), nil
}

// UnmarshalText uses ParseFeedRef
func (fr *FeedRef) UnmarshalText(input []byte) error {
	txt := string(input)

	newRef, err := ParseFeedRef(txt)
	if err == nil {
		*fr = newRef
		return nil
	}

	asURI, err := parseCaononicalURI(txt)
	if err != nil {
		return err
	}

	newFeedRef, ok := asURI.Feed()
	if !ok {
		return fmt.Errorf("Web3R uri is not a feed ref: %s", asURI.Kind())
	}

	*fr = newFeedRef
	return nil
}

var (
	emptyFeedRef = FeedRef{}
	emptyMsgRef  = MessageRef{}
)

// ParseFeedRef uses ParseRef and checks that it returns a *FeedRef
// e.g: @CwBkqhcOueI7hPsJJKy3Bp6uN2Mexo+jRO+hwQSYjLE=.ed25519
func ParseFeedRef(str string) (FeedRef, error) {
	if len(str) == 0 {
		return emptyFeedRef, fmt.Errorf("Web3R: feedRef empty")
	}

	split := strings.Split(str[1:], ".")
	if len(split) < 2 {
		asURL, err := url.Parse(str)
		if err != nil {
			return emptyFeedRef, fmt.Errorf("failed to parse as URL: %s: %w", err, ErrInvalidRef)
		}
		if asURL.Scheme != "web3r" {
			return emptyFeedRef, fmt.Errorf("expected web3 protocol scheme on URL: %q: %w", str, ErrInvalidRef)
		}
		asURI, err := parseCaononicalURI(asURL.Opaque)
		if err != nil {
			return emptyFeedRef, err
		}
		feedRef, ok := asURI.Feed()
		if !ok {
			return emptyFeedRef, fmt.Errorf("Web3RURI is not a feed ref")
		}
		return feedRef, nil
	}

	raw, err := base64.StdEncoding.DecodeString(split[0])
	if err != nil {
		return emptyFeedRef, fmt.Errorf("feedRef: couldn't parse %q: %s: %w", str, err, ErrInvalidHash)
	}

	if str[0] != '@' {
		return emptyFeedRef, ErrInvalidRefType
	}

	var rope RefRope
	switch RefRope(split[1]) {
	case RefFeedWEB3R:
		rope = RefFeedWEB3R
	case RefFeedGabby:
		rope = RefFeedGabby
	case RefFeedBendyButt:
		rope = RefFeedBendyButt
	default:
		return emptyFeedRef, fmt.Errorf("unhandled feed algorithm: %s: %w", str, ErrInvalidRefAlgo)
	}

	if n := len(raw); n != 32 {
		return emptyFeedRef, newFeedRefLenError(n)
	}

	newRef := FeedRef{rope: rope}
	copy(newRef.id[:], raw)
	return newRef, nil

}
