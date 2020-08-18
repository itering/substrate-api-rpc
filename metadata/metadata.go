package metadata

import (
	"github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/util"
)

type RuntimeRaw struct {
	Spec int
	Raw  string
}

type Instant types.MetadataStruct

var (
	latestSpec      = -1
	RuntimeMetadata = make(map[int]*Instant)
	Decoder         *scalecodec.MetadataDecoder
)

func Latest(runtime *RuntimeRaw) *Instant {
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
	instant := Instant(m.Metadata)
	RuntimeMetadata[latestSpec] = &instant
	return &instant
}

func Process(runtime *RuntimeRaw) *Instant {
	if runtime == nil {
		return nil
	}
	if d, ok := RuntimeMetadata[runtime.Spec]; ok {
		return d
	}

	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(runtime.Raw))
	_ = m.Process()

	instant := Instant(m.Metadata)
	RuntimeMetadata[runtime.Spec] = &instant

	return &instant
}

func RegNewMetadataType(spec int, coded string) *Instant {
	m := scalecodec.MetadataDecoder{}
	m.Init(util.HexToBytes(coded))
	_ = m.Process()

	instant := Instant(m.Metadata)
	RuntimeMetadata[spec] = &instant

	if spec > latestSpec {
		latestSpec = spec
	}
	return &instant
}
