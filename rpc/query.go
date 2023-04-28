package rpc

import (
	"fmt"
	"math/rand"

	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/util/ss58"
	"github.com/itering/substrate-api-rpc/websocket"
)

func GetPaymentQueryInfo(p websocket.WsConn, encodedExtrinsic string) (paymentInfo *model.PaymentQueryInfo, err error) {
	v := &model.JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, SystemPaymentQueryInfo(rand.Intn(10000), util.AddHex(encodedExtrinsic))); err != nil {
		return
	}
	paymentInfo = v.ToPaymentQueryInfo()
	if paymentInfo == nil {
		return nil, fmt.Errorf("get PaymentQueryInfo error")
	}
	return
}

func ReadStorageByKey(p websocket.WsConn, key storageKey.StorageKey, hash string) (r storage.StateStorage, err error) {
	v := &model.JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, StateGetStorage(rand.Intn(10000), key.EncodeKey, hash)); err != nil {
		return
	}
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			return storage.StateStorage(""), nil
		}
		return storage.Decode(dataHex, key.ScaleType, nil)
	}
	return
}

func GetMetadataByHash(p websocket.WsConn, hash ...string) (string, error) {
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, StateGetMetadata(rand.Intn(10), hash...)); err != nil {
		return "", err
	}
	return v.ToString()
}

func GetSystemProperties(p websocket.WsConn) (*model.Properties, error) {
	var t model.Properties
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, SystemProperties(rand.Intn(1000))); err != nil {
		return nil, err
	}
	err := v.ToAnyThing(&t)
	return &t, err
}

func GetChainGetBlockHash(p websocket.WsConn, blockNum int) (string, error) {
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, ChainGetBlockHash(rand.Intn(1000), blockNum)); err != nil {
		return "", err
	}
	return v.ToString()
}

func GetStateGetRuntimeVersion(p websocket.WsConn, hash string) *model.RuntimeVersion {
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, StateGetRuntimeVersion(rand.Intn(1000), hash)); err != nil {
		return nil
	}
	if err := v.CheckErr(); err != nil {
		return nil
	}
	return v.ToRuntimeVersion()
}

func GetSystemAccountNextIndex(p websocket.WsConn, accountId string) uint64 {
	v := &model.JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, SystemAccountNextIndex(rand.Intn(1000), ss58.Encode(accountId, 42))); err != nil {
		return 0
	}
	if err := v.CheckErr(); err != nil {
		return 0
	}
	return uint64(v.ToFloat64())
}
