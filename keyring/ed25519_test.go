package keyring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEd25519(t *testing.T) {
	seeds := [][2]string{
		{"0xfae8cdffd9f3efbf545d6783eb7794936da0b1ba53eb31fc26cf67dc26021f3a", "ac75facf017d847db23fab2e885130fb1ee7a7375877fea57f042e171db3f632"},
		{"0xc52de3e01f30ca68409ef6b78e8900300acdede3c18097dc2532677f7596498c", "92ccb44c9d97bb3f2ec6f9ff0ae0c8a0ba5bae9519ab976943ff3eac32b6ef1e"},
		{"0x7cb8d571abaf5135999a28cb46d390b9b34434e9cbcbf65d8a831c75a24ae7d7", "25f20a586ef931343361549898d7beb5dec4daeb54e1df40ea5487f535a51bf1"},
	}
	for _, seed := range seeds {
		edKeyRing := NewEd25519(seed[0])
		assert.Equal(t, edKeyRing.Type(), Ed25519Type)
		assert.Equal(t, edKeyRing.PublicKey(), seed[1])
		assert.NotEmpty(t, edKeyRing.Sign(seed[0]))
	}
}
