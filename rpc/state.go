package rpc

import (
	"math/rand"

	"github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
)

// ReadStorage Read substrate storage
func ReadStorage(p websocket.WsConn, module, prefix string, hash string, arg ...string) (r storage.StateStorage, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)
	v := &model.JsonRpcResult{}
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
	v := &model.JsonRpcResult{}
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
