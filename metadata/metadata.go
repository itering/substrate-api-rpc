package metadata

import (
	"github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/util"
	"strings"
)

type RuntimeRaw struct {
	Spec int
	Raw  string
}

var (
	latestSpec      = -1
	RuntimeMetadata = make(map[int]*types.MetadataStruct)
	Decoder         *scalecodec.MetadataDecoder
)

func Latest(runtime *RuntimeRaw) *types.MetadataStruct {
	if latestSpec != -1 {
		d := RuntimeMetadata[latestSpec]
		return d
	}
	if runtime == nil {
		return nil
	}
	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(runtime.Raw))
	_ = m.Process()

	Decoder = &m
	latestSpec = runtime.Spec

	RuntimeMetadata[latestSpec] = &m.Metadata
	return RuntimeMetadata[latestSpec]
}

func Process(runtime *RuntimeRaw) *types.MetadataStruct {
	if runtime == nil {
		return nil
	}
	if d, ok := RuntimeMetadata[runtime.Spec]; ok {
		return d
	}

	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(runtime.Raw))
	_ = m.Process()

	RuntimeMetadata[runtime.Spec] = &m.Metadata

	return RuntimeMetadata[runtime.Spec]
}

func RegNewMetadataType(spec int, coded string) *types.MetadataStruct {
	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(coded))
	_ = m.Process()

	RuntimeMetadata[spec] = &m.Metadata

	if spec > latestSpec {
		latestSpec = spec
	}
	return RuntimeMetadata[spec]
}

func ModuleStorageMapType(m *types.MetadataStruct, section, method string) (string, *types.StorageType) {
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
