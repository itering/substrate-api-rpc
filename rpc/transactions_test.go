package rpc

import (
	"testing"

	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const (
	AliceSeed    = "0xe5be9a5092b81bca64be81d212e7f2f9eba183bb7a90954f7b76361f6edb5c0a"
	BobAccountId = "8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48"
)

var client *Client

func init() {
	websocket.SetEndpoint("wss://rpc.shibuya.astar.network")
	// websocket init and set metadata
	_, err := websocket.Init()
	if err != nil {
		panic(err)
	}
	client = &Client{}
	// metadata
	raw, err := GetMetadataByHash(nil)
	if err != nil {
		panic(err)
	}
	client.SetMetadata(metadata.RegNewMetadataType(92, raw))
}

type SystemAccountInfo struct {
	Data struct {
		Free decimal.Decimal `json:"free"`
	}
}

// Test SignAndSendTransaction function
func Test_SignAndSendTransaction(t *testing.T) {
	client.SetKeyRing(keyring.New(keyring.Sr25519Type, AliceSeed))
	bobAccountRaw, err := ReadStorage(nil, "System", "Account", "", BobAccountId)
	assert.NoError(t, err)

	var bobAccount1 SystemAccountInfo
	bobAccountRaw.ToAny(&bobAccount1)

	signedTransaction, err := client.SignTransaction("Balances", "transfer", map[string]interface{}{"Id": BobAccountId}, 12345)
	assert.NoError(t, err)

	// async, will return until the transaction is in the block, and return the hash of the block
	blockHash, err := client.SendAuthorSubmitAndWatchExtrinsic(signedTransaction)
	assert.NoError(t, err)
	assert.Equal(t, len(blockHash), 66)
	bobAccountRaw, err = ReadStorage(nil, "System", "Account", "", BobAccountId)
	assert.NoError(t, err)

	var bobAccount2 SystemAccountInfo
	bobAccountRaw.ToAny(&bobAccount2)
	assert.Equal(t, bobAccount2.Data.Free.Sub(bobAccount1.Data.Free).String(), "12345")

	// Synchronize, will return the hash of the transaction
	signedTransaction, err = client.SignTransaction("Balances", "transfer", map[string]interface{}{"Id": BobAccountId}, 23456)
	assert.NoError(t, err)
	hash, err := client.SendAuthorSubmitExtrinsic(signedTransaction)
	assert.NoError(t, err)
	assert.Equal(t, len(hash), 66)

	client.p.Close()
}
