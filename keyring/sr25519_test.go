package keyring

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/itering/scale.go/utiles"
	"github.com/stretchr/testify/assert"
)

func TestSr25519SignAndVerify(t *testing.T) {
	kp, err := GenerateSr25519Keypair()
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("helloworld")
	sig, err := kp.Sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	pub := kp.Public().(*Sr25519PublicKey)
	assert.Equal(t, "68656c6c6f776f726c64", hex.EncodeToString(msg))
	ok := pub.Verify(msg, sig)
	if !ok {
		t.Fatal("Fail: did not verify sr25519 sig")
	}
}

func TestSr25519PublicKeys(t *testing.T) {
	kp, err := GenerateSr25519Keypair()
	if err != nil {
		t.Fatal(err)
	}

	priv := kp.Private().(*Sr25519PrivateKey)
	kp2, err := NewSr25519Keypair(priv.key)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(kp.Public(), kp2.Public()) {
		t.Fatalf("Fail: pubkeys do not match got %x expected %x", kp2.Public(), kp.Public())
	}
}

func TestSr25519EncodeAndDecodePrivateKey(t *testing.T) {
	kp, err := GenerateSr25519Keypair()
	if err != nil {
		t.Fatal(err)
	}

	enc := kp.Private().Encode()
	res := new(Sr25519PrivateKey)
	err = res.Decode(enc)
	if err != nil {
		t.Fatal(err)
	}

	exp := kp.Private().(*Sr25519PrivateKey).key.Encode()
	if !reflect.DeepEqual(res.key.Encode(), exp) {
		t.Fatalf("Fail: got %x expected %x", res.key.Encode(), exp)
	}
}

func TestSr25519EncodeAndDecodePublicKey(t *testing.T) {
	kp, err := GenerateSr25519Keypair()
	if err != nil {
		t.Fatal(err)
	}

	enc := kp.Public().Encode()
	res := new(Sr25519PublicKey)
	err = res.Decode(enc)
	if err != nil {
		t.Fatal(err)
	}

	exp := kp.Public().(*Sr25519PublicKey).key.Encode()
	if !reflect.DeepEqual(res.key.Encode(), exp) {
		t.Fatalf("Fail: got %v expected %v", res.key, exp)
	}
}

func TestSr25519FromSeed(t *testing.T) {
	sk, err := NewSr25519FromSeed(utiles.HexToBytes("0x0aff680b436f6f5622f4a8030148dc4b712f02bb3b96e3dcc21ebbaeade51811"))
	if err != nil {
		t.Fatal(err)
	}

	pk := sk.Public()
	pkEncode := pk.Encode()
	var mss []byte
	mss = append(mss, pkEncode[:]...)
	exp := "3c753c1d5859b082aa23cc7c1dc27b529c4a301dec1c06a6f650c0547901c43f"
	if !reflect.DeepEqual(utiles.BytesToHex(mss), exp) {
		t.Fatalf("Fail: got %s expected %s", utiles.BytesToHex(mss), exp)
	}
	// sig
	msg := []byte("helloworld")
	sig, err := sk.Sign(msg)
	if err != nil {
		t.Fatal(err)
	}
	ok := pk.Verify(msg, sig)
	if !ok {
		t.Fatal("Fail: did not verify sr25519 sig")
	}
}
