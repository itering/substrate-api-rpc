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

### KeyRing

#### Generate KeyPair

```
keyRing := keyring.New(keyring.Sr25519Type, AliceSeed) // sr25519
keyRing := keyring.New(keyring.Ed25519Type, AliceSeed)
```

#### Sign Message

```
keyRing.Sign("hello world") // sign utf-8 message
keyRing.Sign("0xffff") // sign hex message
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

#### Send Extrinsic

```
// set websocket endoint 
websocket.SetEndpoint("wss://shibuya-rpc.dwellir.com")
client = &rpc.Client{}
// init latest metadata
raw, err := GetMetadataByHash(nil)
if err != nil {
    panic(err)
}
// set metadata
client.SetMetadata(metadata.RegNewMetadataType(92, raw))/
// set sr25519 seed
client.SetKeyRing(keyring.New(keyring.Sr25519Type, AliceSeed))
// sign transaction
signedTransaction, err := client.SignTransaction("Balances", "transfer", map[string]interface{}{"Id": BobAccountId}, 12345)

// send transaction async
blockHash, err := client.SendAuthorSubmitAndWatchExtrinsic(signedTransaction)

// send transaction synchronize
blockHash, err := client.SendAuthorSubmitExtrinsic(signedTransaction)
```

More information can be viewed https://polkadot.js.org/api/substrate/rpc.html

## Contributions

We welcome contributions of any kind. Issues labeled can be good (first) contributions.

## LICENSE

Apache-2.0