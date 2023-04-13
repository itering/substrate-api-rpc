package keyring

import "C"
import (
	"errors"
	"strings"

	sr25519 "github.com/ChainSafe/go-schnorrkel"
	"github.com/itering/scale.go/utiles"
)

type Sr25519 struct {
	priv *Sr25519Keypair
}

type Keypair interface {
	Sign(msg []byte) ([]byte, error)
	Public() PublicKey
	Private() PrivateKey
}

type PublicKey interface {
	Verify(msg, sig []byte) bool
	Encode() []byte
	Decode([]byte) error
}

type PrivateKey interface {
	Sign(msg []byte) ([]byte, error)
	Public() (PublicKey, error)
	Encode() []byte
	Decode([]byte) error
}

// NewSr25519 create new sr25519 keyring
func NewSr25519(seed string) *Sr25519 {
	pk, err := NewSr25519FromSeed(utiles.HexToBytes(seed))
	if err != nil {
		panic(err)
	}
	return &Sr25519{priv: pk}
}

// PublicKey return sr25519 public key
func (s *Sr25519) PublicKey() string {
	return utiles.BytesToHex(s.priv.Public().Encode())
}

// Type return keyring type
func (s *Sr25519) Type() Category {
	return Sr25519Type
}

// Sign message by sr25519
// if message has 0x prefix, will decode hex string to bytes
// else will use string bytes
func (s *Sr25519) Sign(message string) string {
	msgBytes := []byte(message)
	if strings.HasPrefix(message, "0x") {
		msgBytes = utiles.HexToBytes(message)
	}
	sig, err := s.priv.Sign(msgBytes)
	if err != nil {
		return ""
	}
	return utiles.BytesToHex(sig)
}

const SecretKeySize = 64

// SigningContext is the context for signatures used or created with substrate
var SigningContext = []byte("substrate")

// Sr25519Keypair is a sr25519 public-private keypair
type Sr25519Keypair struct {
	public  *Sr25519PublicKey
	private *Sr25519PrivateKey
}

type Sr25519PublicKey struct {
	key *sr25519.PublicKey
}

type Sr25519PrivateKey struct {
	key *sr25519.SecretKey
}

// NewSr25519Keypair returns a Sr25519Keypair given a schnorrkel secret key
func NewSr25519Keypair(priv *sr25519.SecretKey) (*Sr25519Keypair, error) {
	pub, err := priv.Public()
	if err != nil {
		return nil, err
	}

	return &Sr25519Keypair{
		public:  &Sr25519PublicKey{key: pub},
		private: &Sr25519PrivateKey{key: priv},
	}, nil
}

// GenerateSr25519Keypair returns a new sr25519 keypair
func GenerateSr25519Keypair() (*Sr25519Keypair, error) {
	priv, pub, err := sr25519.GenerateKeypair()
	if err != nil {
		return nil, err
	}

	return &Sr25519Keypair{
		public:  &Sr25519PublicKey{key: pub},
		private: &Sr25519PrivateKey{key: priv},
	}, nil
}

// Sign uses the keypair to sign the message using the sr25519 signature algorithm
func (kp *Sr25519Keypair) Sign(msg []byte) ([]byte, error) {
	return kp.private.Sign(msg)
}

// Public returns the public key corresponding to this keypair
func (kp *Sr25519Keypair) Public() PublicKey {
	return kp.public
}

// Private returns the private key corresponding to this keypair
func (kp *Sr25519Keypair) Private() PrivateKey {
	return kp.private
}

// Sign uses the private key to sign the message using the sr25519 signature algorithm
func (k *Sr25519PrivateKey) Sign(msg []byte) ([]byte, error) {
	if k.key == nil {
		return nil, errors.New("key is nil")
	}
	t := sr25519.NewSigningContext(SigningContext, msg)
	sig, err := k.key.Sign(t)
	if err != nil {
		return nil, err
	}
	enc := sig.Encode()
	return enc[:], nil
}

// Public returns the public key corresponding to this private key
func (k *Sr25519PrivateKey) Public() (PublicKey, error) {
	if k.key == nil {
		return nil, errors.New("key is nil")
	}
	pub, err := k.key.Public()
	if err != nil {
		return nil, err
	}
	return &Sr25519PublicKey{key: pub}, nil
}

// Encode returns the 32-byte encoding of the private key
func (k *Sr25519PrivateKey) Encode() []byte {
	if k.key == nil {
		return nil
	}
	enc := k.key.Encode()
	return enc[:]
}

// Decode decodes the input bytes into a private key and sets the receiver the decoded key
// Input must be 32 bytes, or else this function will error
func (k *Sr25519PrivateKey) Decode(in []byte) error {
	if len(in) != 32 {
		return errors.New("input to sr25519 private key decode is not 32 bytes")
	}
	b := [32]byte{}
	copy(b[:], in)
	k.key = &sr25519.SecretKey{}
	return k.key.Decode(b)
}

// Verify uses the sr25519 signature algorithm to verify that the message was signed by
// this public key; it returns true if this key created the signature for the message,
// false otherwise
func (k *Sr25519PublicKey) Verify(msg, sig []byte) bool {
	if k.key == nil {
		return false
	}

	b := [64]byte{}
	copy(b[:], sig)

	s := &sr25519.Signature{}
	err := s.Decode(b)
	if err != nil {
		return false
	}

	t := sr25519.NewSigningContext(SigningContext, msg)
	result, err := k.key.Verify(s, t)
	if err != nil {
		return false
	}
	return result
}

// Encode returns the 32-byte encoding of the public key
func (k *Sr25519PublicKey) Encode() []byte {
	if k.key == nil {
		return nil
	}

	enc := k.key.Encode()
	return enc[:]
}

// Decode decodes the input bytes into a public key and sets the receiver the decoded key
// Input must be 32 bytes, or else this function will error
func (k *Sr25519PublicKey) Decode(in []byte) error {
	if len(in) != 32 {
		return errors.New("input to sr25519 public key decode is not 32 bytes")
	}
	b := [32]byte{}
	copy(b[:], in)
	k.key = &sr25519.PublicKey{}
	return k.key.Decode(b)
}

// NewSr25519FromSeed from Secret seed to SecretKey
func NewSr25519FromSeed(seed []byte) (*Sr25519Keypair, error) {
	switch len(seed) {
	case sr25519.MiniSecretKeySize:
		var mss [32]byte
		copy(mss[:], seed)
		ms, err := sr25519.NewMiniSecretKeyFromRaw(mss)
		if err != nil {
			return nil, err
		}
		return NewSr25519Keypair(ms.ExpandEd25519())
	case SecretKeySize:
		var key, nonce [32]byte
		copy(key[:], seed[0:32])
		copy(nonce[:], seed[32:64])
		secret := sr25519.NewSecretKey(key, nonce)
		return NewSr25519Keypair(secret)
	}
	return nil, errors.New("invalid seed length")
}
