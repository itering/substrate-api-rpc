# substrate-api-rpc

## Contents

- [Install](#Install)
- [Usage](#Usage)
   - [Codec](#Codec)
   - [RPC](#RPC)
- [Contributions](#Contributions)
- [LICENSE](#LICENSE)


## Install

```
go get github.com/itering/substrate-api-rpc
```

## Usage

### Codec

#### Extrinsic Decode

```
metadataRaw := "" // rpc state_getMetadata
specVersion := 0  // rpc chain_getRuntimeVersion
encodeExtrinsic := []string{"0x280402000b10449a7e7301", "0x1c0407005e8b4100"}
decodeExtrinsics, err := substrate.DecodeExtrinsic(encodeExtrinsic, metadata.Process(metadataRaw), specVersion)
```

#### Event Decode

```
metadataRaw := "" // rpc state_getMetadata
specVersion := 0  // rpc chain_getRuntimeVersion
event = "0x080000000000000080e36a09000000000200000001000000000000ca9a3b00000000020000"
substrate.DecodeEvent(event, metadataInstant, specVersion)
```

#### Log Decode

```
logs := ["0x054241424501014a7024ec6c4be378c35c254860d8f4ddc6f9d53ea8ce42ca00bc77c280511f1cb4c93fbd825e3c7dcabb36221372a9b5359c496e095d31afc359bdb9fac45487"]
substrate.DecodeLogDigest(logs)
```

#### Storage Decode

```
raw := "0x2efb"
storage.Decode(raw, "i16", nil)
```


### RPC

#### Substrate RPC 

Example

> state_getMetadata

```
blockHash := ""
rpc.GetMetadataByHash(conn, blockHash)
```

> state_getStorage

```
validatorsRaw, err := rpc.ReadStorage(conn, "Session", "Validators", blockHash)
validatorList := validatorsRaw.ToStringSlice()
```

More information can be viewed https://polkadot.js.org/api/substrate/rpc.html


## Contributions

We welcome contributions of any kind. Issues labeled can be good (first) contributions.

## LICENSE

GPL-3.0