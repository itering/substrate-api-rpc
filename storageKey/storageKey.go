package storageKey

import (
	"strings"

	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/hasher"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/util"
)

type StorageKey struct {
	EncodeKey string
	ScaleType string
}

type Storage struct {
	Prefix string
	Method string
	Type   types.StorageType
}

func EncodeStorageKey(section, method string, args ...string) (storageKey StorageKey) {
	m := metadata.Latest(nil)
	if m == nil {
		return
	}

	method = upperCamel(method)
	prefix, storageType := moduleStorageMapType(m, section, method)
	if storageType == nil {
		return
	}

	mapType := checkoutHasherAndType(storageType)
	if mapType == nil {
		return
	}

	storageKey.ScaleType = mapType.Value

	var hash []byte
	sectionHash := hasher.HashByCryptoName([]byte(upperCamel(prefix)), "Twox128")
	methodHash := hasher.HashByCryptoName([]byte(method), "Twox128")

	hash = append(sectionHash, methodHash[:]...)

	for index, arg := range args {
		hash = append(hash, hasher.HashByCryptoName(util.HexToBytes(arg), mapType.Hasher[index])[:]...)
	}
	storageKey.EncodeKey = util.BytesToHex(hash)
	return
}

type storageOption struct {
	Value  string   `json:"value"`
	Keys   []string `json:"keys"`
	Hasher []string `json:"hasher"`
}

func checkoutHasherAndType(t *types.StorageType) *storageOption {
	option := storageOption{}
	switch t.Origin {
	case "MapType":
		option.Value = t.MapType.Value
		option.Hasher = []string{t.MapType.Hasher}
	case "DoubleMapType":
		option.Value = t.DoubleMapType.Value
		option.Keys = []string{t.DoubleMapType.Key, t.DoubleMapType.Key2}
		option.Hasher = []string{t.DoubleMapType.Hasher, t.DoubleMapType.Key2Hasher}
	case "Map":
		option.Value = t.NMapType.Value
		option.Keys = t.NMapType.KeyVec
		option.Hasher = t.NMapType.Hashers
	default:
		option.Value = *t.PlainType
		option.Hasher = []string{"Twox64Concat"}
	}
	return &option
}

func upperCamel(s string) string {
	if len(s) == 0 {
		return ""
	}
	s = strings.ToUpper(string(s[0])) + string(s[1:])
	return s
}

func moduleStorageMapType(m *metadata.Instant, section, method string) (string, *types.StorageType) {
	modules := m.Metadata.Modules
	for _, value := range modules {
		if strings.EqualFold(strings.ToLower(value.Name), strings.ToLower(section)) {
			for _, storage := range value.Storage {
				if strings.EqualFold(strings.ToLower(storage.Name), strings.ToLower(method)) {
					return value.Prefix, &storage.Type
				}
			}
		}
	}
	return "", nil
}
