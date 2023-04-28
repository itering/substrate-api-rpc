package keyring

import (
	"crypto"
	"crypto/ed25519"
	"strings"

	"github.com/itering/scale.go/utiles"
)

type Ed25519 struct {
	priv ed25519.PrivateKey
}

// NewEd25519 create new ed25519 keyring
// seed: 0x + 64 hex string
// if seed length not 64, will panic
func NewEd25519(seed string) *Ed25519 {
	return &Ed25519{priv: ed25519.NewKeyFromSeed(utiles.HexToBytes(seed))}
}

// PublicKey return ED25519 public key
func (e *Ed25519) PublicKey() string {
	return utiles.BytesToHex(e.priv.Public().(ed25519.PublicKey))
}

// Type return keyring type
func (e *Ed25519) Type() Category {
	return Ed25519Type
}

// Sign message by ed25519
// if message has 0x prefix, will decode hex string to bytes
// else will use string bytes
func (e *Ed25519) Sign(message string) string {
	msgBytes := []byte(message)
	if strings.HasPrefix(message, "0x") {
		msgBytes = utiles.HexToBytes(message)
	}
	var noHash crypto.Hash
	sig, err := e.priv.Sign(nil, msgBytes, noHash)
	if err != nil {
		return ""
	}
	return utiles.BytesToHex(sig)
}
