package aliases

import (
	"bytes"
	"crypto/ed25519"

	"github.com/calvindc/Web3RpcHub/refs"
)

// Registration ties an alias to the ID of the user and the HubID it should be registered on
type Registration struct {
	Alias  string
	UserID refs.FeedRef
	HubID  refs.FeedRef
}

// Sign takes the public key (belonging to UserID) and returns the signed confirmation
func (r Registration) Sign(privKey ed25519.PrivateKey) Confirmation {
	var conf Confirmation
	conf.Registration = r
	msg := r.createRegistrationMessage()
	conf.Signature = ed25519.Sign(privKey, msg)
	return conf
}

// createRegistrationMessage returns the string of bytes that should be signed
func (r Registration) createRegistrationMessage() []byte {
	var message bytes.Buffer
	message.WriteString("=hub-alias-registration:")
	message.WriteString(r.HubID.String())
	message.WriteString(":")
	message.WriteString(r.UserID.String())
	message.WriteString(":")
	message.WriteString(r.Alias)
	return message.Bytes()
}

// Confirmation combines a registration with the corresponding signature
type Confirmation struct {
	Registration

	Signature []byte
}

// Verify checks that the confirmation is for the expected hub and from the expected feed
func (c Confirmation) Verify() bool {
	// re-construct the registration
	message := c.createRegistrationMessage()

	// check the signature matches
	return ed25519.Verify(c.UserID.PubKey(), message, c.Signature)
}
