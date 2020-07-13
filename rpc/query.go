package rpc

import (
	"fmt"

	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
	"math/rand"

	"github.com/itering/substrate-api-rpc/util"
)

func GetCurrentEra(p *websocket.PoolConn, hash string) (int, error) {
	eraIndex, err := ReadStorage(p, "Staking", "CurrentEra", hash)
	if err != nil {
		return 0, err
	}
	return eraIndex.ToInt(), nil
}

func GetActiveEra(p *websocket.PoolConn, hash string) (int, error) {
	eraIndex, err := ReadStorage(p, "Staking", "ActiveEra", hash)
	if err != nil {
		return 0, err
	}
	if era := eraIndex.ToActiveEraInfo(); era != nil {
		return era.Index, nil
	}
	return 0, fmt.Errorf("decode ActiveEra error")
}

// Read substrate storage
func ReadStorage(p *websocket.PoolConn, module, prefix string, hash string, arg ...string) (r storage.StateStorage, err error) {
	key := storageKey.EncodeStorageKey(module, prefix, arg...)
	v := &JsonRpcResult{}
	if err = websocket.SendWsRequest(p, v, StateGetStorage(rand.Intn(10000), util.AddHex(key.EncodeKey), hash)); err != nil {
		return
	}
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			return storage.StateStorage(""), nil
		}
		return storage.Decode(dataHex, key.ScaleType, nil)
	}
	return r, err

}

func ReadKeysPaged(p *websocket.PoolConn, module, prefix string) (r []string, scale string, err error) {
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

func GetPaymentQueryInfo(p *websocket.PoolConn, encodedExtrinsic string) (paymentInfo *PaymentQueryInfo, err error) {
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

func ReadStorageByKey(p *websocket.PoolConn, key storageKey.StorageKey, hash string) (r storage.StateStorage, err error) {
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

func GetMetadataByHash(p *websocket.PoolConn, hash ...string) (string, error) {
	v := &JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, StateGetMetadata(rand.Intn(10), hash...)); err != nil {
		return "", err
	}
	return v.ToString()
}

func GetFreeBalance(p *websocket.PoolConn, accountId, hash string) (decimal.Decimal, decimal.Decimal, error) {
	var accountValue storage.StateStorage
	var err error
	accountValue, err = ReadStorage(p, "System", "Account", hash, util.TrimHex(accountId))
	if err == nil {
		if account := accountValue.ToAccountInfo(); account != nil {
			return account.Data.Free.Add(account.Data.Reserved), decimal.Zero, nil
		}
	}

	return decimal.Zero, decimal.Zero, err
}

func GetAccountLock(p *websocket.PoolConn, address string) (balance decimal.Decimal, err error) {
	var sv storage.StateStorage
	sv, err = ReadStorage(p, "Balances", "Locks", "", util.TrimHex(address))
	if err == nil {
		if locks := sv.ToBalanceLock(); len(locks) > 0 {
			for _, lock := range locks {
				if lock.Amount.GreaterThanOrEqual(balance) {
					balance = lock.Amount
				}
			}
			return balance, nil
		}
	}
	return
}

func GetValidatorFromSub(p *websocket.PoolConn, hash string) ([]string, error) {
	validators, err := ReadStorage(p, "Session", "Validators", hash)
	if err != nil {
		return []string{}, err
	}
	var r []string
	for _, address := range validators.ToStringSlice() {
		r = append(r, util.TrimHex(address))
	}
	return r, nil
}

func GetSystemProperties(p *websocket.PoolConn) (*Properties, error) {
	var t Properties
	v := &JsonRpcResult{}
	if err := websocket.SendWsRequest(p, v, SystemProperties(rand.Intn(1000))); err != nil {
		return nil, err
	}
	err := v.ToAnyThing(&t)
	return &t, err
}
