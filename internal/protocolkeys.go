package internal

import (
	"fmt"
	"io"

	"os"
	"path/filepath"

	"encoding/base64"
	"encoding/json"

	"strings"

	refs "github.com/ssbc/go-ssb-refs"
	"go.cryptoscope.co/secretstream/secrethandshake"
)

const (
	RefAlgoFeed       = "ed25519"
	RefAlgoMessage    = "sha256"
	RefAlgoBlob       = RefAlgoMessage
	SuffixAlgoFeed    = ".ed25519"
	SuffixAlgoMessage = ".sha256"
)

type KeyPair struct {
	Feed refs.FeedRef
	Pair secrethandshake.EdKeyPair
}

type protocolSecret struct {
	Curve   string       `json:"curve"`
	ID      refs.FeedRef `json:"id"`
	Public  string       `json:"public"`
	Private string       `json:"private"`
}

func IsValidFeedFormat(r refs.FeedRef) error {
	if r.Algo() != RefAlgoFeed {
		return fmt.Errorf("Keys: unsupported feed format:%s", r.Algo())
	}
	return nil
}

func NewKeyPair(ir io.Reader) (*KeyPair, error) {
	kp, err := secrethandshake.GenEdKeyPair(ir)
	if err != nil {
		return nil, fmt.Errorf("Keys: error building key pair: %w", err)
	}
	feed, err := refs.NewFeedRefFromBytes(kp.Public[:], RefAlgoFeed)
	if err != nil {
		return nil, fmt.Errorf("Keys: error building key pair: %w", err)
	}
	return &KeyPair{Feed: feed, Pair: *kp}, nil
}

// LoadKeyPair 读取key文件
func LoadKeyPair(fname string) (*KeyPair, error) {
	f, err := os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("Keys: load key pair file, couldn't open key file %s: %w", fname, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("Keys: load key pair file, couldn't stat key file %s: %w", fname, err)
	}

	if perms := info.Mode().Perm(); perms != secretPermission {
		return nil, fmt.Errorf("Keys: load key pair file, expected key file permissions %s, but got %s", secretPermission, perms)
	}

	return ParseKeyPair(NoCommentReader(f))
}

var secretPermission = os.FileMode(0600)

// SaveKeyPair 写入key文件
func SaveKeyPair(kp KeyPair, path string) error {
	if err := IsValidFeedFormat(kp.Feed); err != nil {
		return err
	}

	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Keys: failed to mkdir for keypair: %w", err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, secretPermission)
	if err != nil {
		return fmt.Errorf("Keys: failed to create key pair file: %w", err)
	}

	if err := EncodeKeyPairAsJSON(kp, f); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("Keys: save key pair(to close file) failed : %w", err)
	}
	return nil
}

// EncodeKeyPairAsJSON serializes key 文件
func EncodeKeyPairAsJSON(kp KeyPair, w io.Writer) error {
	var sec = protocolSecret{
		Curve:   RefAlgoFeed,
		ID:      kp.Feed,
		Public:  base64.StdEncoding.EncodeToString(kp.Pair.Public[:]) + SuffixAlgoFeed,
		Private: base64.StdEncoding.EncodeToString(kp.Pair.Secret[:]) + SuffixAlgoFeed,
	}
	err := json.NewEncoder(w).Encode(sec)
	if err != nil {
		return fmt.Errorf("Keys: encode key pair as json failed: %w", err)
	}
	return nil
}

// ParseKeyPair 转换json格式的key文件to key pair
func ParseKeyPair(r io.Reader) (*KeyPair, error) {
	var s protocolSecret
	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return nil, fmt.Errorf("Keys: parse JSON decoding failed: %w", err)
	}

	if err := IsValidFeedFormat(s.ID); err != nil {
		return nil, err
	}

	public, err := base64.StdEncoding.DecodeString(strings.TrimSuffix(s.Public, SuffixAlgoFeed))
	if err != nil {
		return nil, fmt.Errorf("Keys: parse base64 decode of public part failed: %w", err)
	}

	private, err := base64.StdEncoding.DecodeString(strings.TrimSuffix(s.Private, SuffixAlgoFeed))
	if err != nil {
		return nil, fmt.Errorf("Keys: parse base64 decode of private part failed: %w", err)
	}

	pair, err := secrethandshake.NewKeyPair(public, private)
	if err != nil {
		return nil, fmt.Errorf("Keys: parse base64 decode of private part failed: %w", err)
	}

	ssbkp := KeyPair{
		Feed: s.ID,
		Pair: *pair,
	}
	return &ssbkp, nil
}
