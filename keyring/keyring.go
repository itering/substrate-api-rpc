package keyring

type IKeyRing interface {
	Sign(string) string
	PublicKey() string
	Type() Category
}

type Category string

const (
	Ed25519Type Category = "Ed25519"
	Sr25519Type Category = "Sr25519"
)

func New(category Category, seed string) IKeyRing {
	if category == Ed25519Type {
		return NewEd25519(seed)
	} else if category == Sr25519Type {
		return NewSr25519(seed)
	}
	panic("invalid category")
}
