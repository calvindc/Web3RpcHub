package secrethandshake

import (
	"crypto/sha512"

	"filippo.io/edwards25519"
	"golang.org/x/crypto/ed25519"
)

func PrivateKeyToCurve25519(curve25519Private *[32]byte, privateKey ed25519.PrivateKey) {
	h := sha512.New()
	h.Write(privateKey[:32])
	digest := h.Sum(nil)

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(curve25519Private[:], digest)
}

// PublicKeyToCurve25519 converts an Ed25519 public key into the curve25519
// public key that would be generated from the same private key.
func PublicKeyToCurve25519(curveBytes *[32]byte, edBytes ed25519.PublicKey) bool {
	if IsEdLowOrder(edBytes) {
		return false
	}

	edPoint, err := new(edwards25519.Point).SetBytes(edBytes)
	if err != nil {
		return false
	}

	copy(curveBytes[:], edPoint.BytesMontgomery())
	return true
}
