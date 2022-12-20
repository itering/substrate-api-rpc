package test

import (
	"testing"
	"time"

	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
)

func TestInit(t *testing.T) {
	websocket.Init(
		"",
		websocket.WithEndPoint("wss://rpc.polkadot.io"),
		websocket.WithMaxPoolCap(100),
		websocket.WithHandshakeTimeout(5*time.Second),
		websocket.WithWriteTimeoutTimeout(60*time.Second),
		websocket.WithReadTimeoutTimeout(60*time.Second),
		websocket.WithWriteBufferSize(5242880),
	)

	v := &rpc.JsonRpcResult{}
	websocket.SendWsRequest("", nil, v, rpc.ChainGetBlockHash(1, 1))
	t.Log(v)
}

func TestMultiInit(t *testing.T) {
	// 1.default client
	websocket.Init(
		"",
		websocket.WithEndPoint("wss://rpc.polkadot.io"),
		websocket.WithMaxPoolCap(100),
		websocket.WithHandshakeTimeout(5*time.Second),
		websocket.WithWriteTimeoutTimeout(60*time.Second),
		websocket.WithReadTimeoutTimeout(60*time.Second),
		websocket.WithWriteBufferSize(5242880),
	)

	v := &rpc.JsonRpcResult{}
	websocket.SendWsRequest("", nil, v, rpc.ChainGetBlockHash(1, 1))
	t.Log(v)

	// 2.westend client
	const (
		Westend         websocket.NodeName = "westend"
		WestendEndPoint string             = "wss://westend-rpc.polkadot.io"
	)
	websocket.Init(
		Westend,
		websocket.WithEndPoint(WestendEndPoint),
		websocket.WithMaxPoolCap(100),
		websocket.WithHandshakeTimeout(5*time.Second),
		websocket.WithWriteTimeoutTimeout(60*time.Second),
		websocket.WithReadTimeoutTimeout(60*time.Second),
		websocket.WithWriteBufferSize(5242880),
	)

	v2 := &rpc.JsonRpcResult{}
	websocket.SendWsRequest(Westend, nil, v2, rpc.ChainGetBlockHash(1, 1))
	t.Log(v2)
}
