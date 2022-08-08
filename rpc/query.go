package rpc

import (
	"fmt"
	"math/rand"

	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
)

// Read substrate storage
func ReadStorage(p websocket.WsConn, module, prefix string, hash string, arg ...string) (r storage.StateStorage, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)
	v := &JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), hash)); err != nil {
		return
	}
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			return "", nil
		}
		return storage.Decode(dataHex, key.ScaleType, nil)
	}
	return r, err

}

func ReadKeysPaged(p websocket.WsConn, module, prefix string) (r []string, scale string, err error) {
	key := storageKey.EncodeStorageKey(module, prefix)
	v := &JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, StateGetKeysPaged(rand.Intn(10000), util.AddHex(key.EncodeKey))); err != nil {
		return
	}
	if keys, err := v.ToInterfaces(); err == nil {
		for _, k := range keys {
			r = append(r, k.(string))
		}
	}
	return r, key.ScaleType, err
}

func GetPaymentQueryInfo(p websocket.WsConn, encodedExtrinsic string) (paymentInfo *PaymentQueryInfo, err error) {
	v := &JsonRpcResult{}
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
	v := &JsonRpcResult{}
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
	v := &JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, StateGetMetadata(rand.Intn(10), hash...)); err != nil {
		return "", err
	}
	return v.ToString()
}

func GetSystemProperties(p websocket.WsConn) (*Properties, error) {
	var t Properties
	v := &JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, SystemProperties(rand.Intn(1000))); err != nil {
		return nil, err
	}
	err := v.ToAnyThing(&t)
	return &t, err
}
