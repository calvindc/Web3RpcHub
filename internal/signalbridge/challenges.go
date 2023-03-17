package signalbridge

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/calvindc/Web3RpcHub/internal/refs"
)

// challengeLength 访问hub使用256位的nonces
const challengeLength = 32

func DecodeChallenge(c string) ([]byte, error) {
	challengeBytes, err := base64.URLEncoding.DecodeString(c)
	if err != nil {
		return nil, fmt.Errorf("invalid challenge encoding: %w", err)
	}

	if n := len(challengeBytes); n != challengeLength {
		return nil, fmt.Errorf("invalid challenge length: expected %d but got %d", challengeLength, n)
	}

	return challengeBytes, nil
}

func GenerateChallenge() string {
	buf := make([]byte, challengeLength)
	rand.Read(buf)
	return base64.URLEncoding.EncodeToString(buf)
}

type ClientPayload struct {
	ClientID, ServerID refs.FeedRef

	ClientChallenge string
	ServerChallenge string
}

// recreate the signed message
func (cr ClientPayload) createMessage() []byte {
	var msg bytes.Buffer
	msg.WriteString("=http-auth-sign-in:")
	msg.WriteString(cr.ServerID.String())
	msg.WriteString(":")
	msg.WriteString(cr.ClientID.String())
	msg.WriteString(":")
	msg.WriteString(cr.ServerChallenge)
	msg.WriteString(":")
	msg.WriteString(cr.ClientChallenge)
	return msg.Bytes()
}

// Sign 25519签名
func (cr ClientPayload) Sign(privateKey ed25519.PrivateKey) []byte {
	msg := cr.createMessage()
	return ed25519.Sign(privateKey, msg)
}

// Validate 签名校验
func (cr ClientPayload) Validate(signature []byte) bool {
	msg := cr.createMessage()
	return ed25519.Verify(cr.ClientID.PubKey(), msg, signature)
}
