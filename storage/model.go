package storage

type RawAuraPreDigest struct {
	SlotNumber int64 `json:"slotNumber"`
}

type ExtrinsicParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ValueRaw string      `json:"valueRaw"`
}

type DecoderLog struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type RawBabePreDigest struct {
	Primary   *RawBabePreDigestPrimary      `json:"primary,omitempty"`
	Secondary *RawBabePreDigestSecondary    `json:"secondary,omitempty"`
	VRF       *RawBabePreDigestSecondaryVRF `json:"VRF,omitempty"`
}

type RawBabePreDigestPrimary struct {
	AuthorityIndex uint   `json:"authorityIndex"`
	SlotNumber     uint64 `json:"slotNumber"`
	Weight         uint   `json:"weight"`
	VrfOutput      string `json:"vrfOutput"`
	VrfProof       string `json:"vrfProof"`
}

type RawBabePreDigestSecondary struct {
	AuthorityIndex uint   `json:"authorityIndex"`
	SlotNumber     uint64 `json:"slotNumber"`
	Weight         uint   `json:"weight"`
}

type RawBabePreDigestSecondaryVRF struct {
	AuthorityIndex uint   `json:"authorityIndex"`
	SlotNumber     uint64 `json:"slotNumber"`
	VrfOutput      string `json:"vrfOutput"`
	VrfProof       string `json:"vrfProof"`
}
