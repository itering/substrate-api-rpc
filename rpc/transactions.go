package rpc

import (
	"errors"
	"fmt"
	"github.com/itering/scale.go/utiles"
	"math/rand"

	gorilla "github.com/gorilla/websocket"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/substrate-api-rpc/hasher"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
)

const TxVersionInfo = "84"

var (
	InvalidMetadataErr      = errors.New("invalid Metadata")
	InvalidCallErr          = errors.New("invalid call name or module name")
	InvalidCallArgsErr      = errors.New("invalid call args length")
	NotSetKeyRingErr        = errors.New("not set keyring")
	NetworkErr              = errors.New("network error")
	SubscribeTransactionErr = errors.New("subscribe transaction failed")
)

// SendAuthorSubmitExtrinsic send extrinsic
// will call rpc author_submitExtrinsic
func (cl *Client) SendAuthorSubmitExtrinsic(signedExtrinsic string) (string, error) {
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(cl.p, v, AuthorSubmitExtrinsic(rand.Intn(1000), signedExtrinsic)); err != nil {
		return "", err
	}
	return v.ToString()
}

// SendAuthorSubmitAndWatchExtrinsic send extrinsic and watch
// will call rpc author_submitAndWatchExtrinsic
func (cl *Client) SendAuthorSubmitAndWatchExtrinsic(signedExtrinsic string) (string, error) {
	v := &model.JsonRpcResult{}
	var p *websocket.PoolConn
	var err error
	if cl.p == nil {
		if p, err = websocket.Init(); err != nil {
			return "", nil
		}
		defer p.Close() // nolint: errcheck
		cl.p = p.Conn
	}
	if err = cl.p.WriteMessage(gorilla.TextMessage, AuthorSubmitAndWatchExtrinsic(rand.Intn(1000), signedExtrinsic)); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		return "", fmt.Errorf("websocket send error: %v", err)
	}
	var retry int
	for {
		if err = cl.p.ReadJSON(v); err != nil {
			if p != nil {
				p.MarkUnusable()
			}
			return "", fmt.Errorf("websocket read error: %v", err)
		}
		if err = v.CheckErr(); err != nil {
			return "", err
		}
		if retry > 10 {
			break
		}
		if v.Method == "author_extrinsicUpdate" {
			if AuthorExtrinsicUpdate := v.ToAuthorExtrinsicUpdate(); AuthorExtrinsicUpdate != nil && AuthorExtrinsicUpdate.InBlock != nil {
				return *AuthorExtrinsicUpdate.InBlock, nil
			}
		}
		retry++
	}
	return "", SubscribeTransactionErr
}

// SignTransaction sign transaction
// p: websocket connection
// keyRing: keyring
// moduleName: module name
// callName: call name
// args: call args
// return: transaction hex
func (cl *Client) SignTransaction(moduleName, callName string, args ...interface{}) (string, error) {
	// check metadata
	if cl.metadata == nil {
		// use latest metadata
		m := metadata.Latest(nil)
		if m == nil {
			return "", InvalidMetadataErr
		}
		cl.metadata = m
	}

	// check keyring
	if cl.keyRing == nil {
		return "", NotSetKeyRingErr
	}

	// find call
	call := cl.metadata.FindCallCallName(moduleName, callName)
	if call == nil {
		return "", InvalidCallErr
	}

	// check call args
	if len(call.Args) != len(args) {
		return "", InvalidCallArgsErr
	}

	metadataStruct := types.MetadataStruct(*cl.metadata)
	opt := &types.ScaleDecoderOption{Metadata: &metadataStruct}

	// build params
	var params []scalecodec.ExtrinsicParam
	for _, v := range args {
		params = append(params, scalecodec.ExtrinsicParam{Value: v})
	}

	// encode call
	encodeCall := types.EncodeWithOpt("Call", map[string]interface{}{"call_index": call.Lookup, "params": params}, opt)
	// build extrinsic
	genericExtrinsic := &scalecodec.GenericExtrinsic{
		VersionInfo: TxVersionInfo,
		Signer:      map[string]interface{}{"Id": cl.keyRing.PublicKey()},
		Era:         "00", Nonce: int(GetSystemAccountNextIndex(cl.p, cl.keyRing.PublicKey())),
		Params:   params,
		CallCode: call.Lookup,
	}
	if util.StringInSlice("ChargeAssetTxPayment", metadataStruct.Extrinsic.SignedIdentifier) {
		genericExtrinsic.SignedExtensions = map[string]interface{}{"ChargeAssetTxPayment": map[string]interface{}{"tip": 0, "asset_id": nil}}
	}

	// build payload
	payload, err := cl.buildExtrinsicPayload(encodeCall, genericExtrinsic)
	if err != nil {
		return "", err
	}

	// if payload length > 256, Blake256 hash payload
	if len(util.HexToBytes(payload)) > 256 {
		payload = util.BytesToHex(hasher.HashByCryptoName(util.HexToBytes(payload), "Blake2_256"))
	}
	// sign payload
	genericExtrinsic.SignatureRaw = map[string]interface{}{string(cl.keyRing.Type()): cl.keyRing.Sign(util.AddHex(payload))}

	// send extrinsic will return hash
	encodedExtrinsic, err := genericExtrinsic.Encode(opt)
	if err != nil {
		return "", err
	}
	return util.AddHex(encodedExtrinsic), nil
}

// buildExtrinsicPayload build extrinsic payload
func (cl *Client) buildExtrinsicPayload(encodeCall string, genericExtrinsic *scalecodec.GenericExtrinsic) (string, error) {
	genesisHash, err := GetChainGetBlockHash(cl.p, 0)
	if err != nil {
		return "", NetworkErr
	}
	version := GetStateGetRuntimeVersion(cl.p, "")
	if version == nil {
		return "", NetworkErr
	}
	data := encodeCall
	data = data + types.Encode("EraExtrinsic", genericExtrinsic.Era)   // era
	data = data + types.Encode("Compact<U32>", genericExtrinsic.Nonce) // nonce
	if len(cl.metadata.Extrinsic.SignedIdentifier) > 0 && utiles.SliceIndex("ChargeTransactionPayment", cl.metadata.Extrinsic.SignedIdentifier) > -1 {
		data = data + types.Encode("Compact<Balance>", genericExtrinsic.Tip) // tip
	}
	for identifier, extension := range genericExtrinsic.SignedExtensions {
		for _, ext := range cl.metadata.Extrinsic.SignedExtensions {
			if ext.Identifier == identifier {
				data = data + types.Encode(ext.TypeString, extension)
			}
		}
	}
	data = data + types.Encode("U32", version.SpecVersion)        // specVersion
	data = data + types.Encode("U32", version.TransactionVersion) // transactionVersion
	data = data + util.TrimHex(types.Encode("Hash", genesisHash)) // genesisHash
	data = data + util.TrimHex(types.Encode("Hash", genesisHash)) // blockHash
	return data, nil
}
