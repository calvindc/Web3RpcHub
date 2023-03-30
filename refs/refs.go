package refs

import (
	"encoding"
	"fmt"
	"net/url"
)

// Ref is the abstract interface all reference types should implement.
type Ref interface {
	Rope() RefRope

	// Sigil returns the pre-URI string [ala @foo=.ed25519, %msgkey=.sha256 or &blob=.sha256].
	Sigil() string

	// ShortSigil returns a truncated version of Sigil()
	ShortSigil() string

	// URI prints the reference as a web3r-uri
	URI() string

	fmt.Stringer
	encoding.TextMarshaler
}

// RefRope define a set of known/understood algorithms to reference feeds, messages or blobs
type RefRope string

const (
	RefEncryptCurve    = "ed25519"
	SuffixFeedWEB3R    = ".ed25519"
	SuffixMessageWEB3R = ".sha256"
)

// Some constant identifiers for the supported references
const (
	RefFeedWEB3R    RefRope = "ed25519"
	RefMessageWEB3R RefRope = "sha256"
	RefBlobWEB3R    RefRope = RefMessageWEB3R

	RefFeedBamboo    RefRope = "bamboo"
	RefMessageBamboo RefRope = RefFeedBamboo

	RefFeedBendyButt    RefRope = "bendybutt-v1"
	RefMessageBendyButt RefRope = RefFeedBendyButt

	RefCloakedGroup RefRope = "cloaked"

	RefFeedGabby    RefRope = "gabbygrove-v1" // cbor based chain
	RefMessageGabby RefRope = RefFeedGabby
)

// ParseRef either returns an parsed and understood reference or an error
func ParseRef(str string) (Ref, error) {
	if len(str) == 0 {
		return nil, ErrInvalidRef
	}

	switch string(str[0]) {
	case "@": //feed tag
		return ParseFeedRef(str)
	case "%": //message tag
		return ParseMessageRef(str)
	case "&": //blob tag
		return ParseBlobRef(str)
	default:
		asURL, err := url.Parse(str)
		if err != nil {
			return nil, fmt.Errorf("failed to parse as URL: %s: %w", err, ErrInvalidRefType)
		}
		if asURL.Scheme != "web3r" {
			return nil, fmt.Errorf("expected web3r protocl scheme on URL: %q: %w", str, ErrInvalidRefType)
		}
		asWeb3rURI, err := parseCaononicalURI(asURL.Opaque)
		return asWeb3rURI.ref, err
	}
}
