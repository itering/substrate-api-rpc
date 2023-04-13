package rpc

import (
	"github.com/itering/substrate-api-rpc/keyring"
	"github.com/itering/substrate-api-rpc/metadata"
	"github.com/itering/substrate-api-rpc/websocket"
)

type Client struct {
	p        websocket.WsConn
	metadata *metadata.Instant
	keyRing  keyring.IKeyRing
}

func (cl *Client) SetKeyRing(keyRing keyring.IKeyRing) {
	cl.keyRing = keyRing
}

func (cl *Client) SetMetadata(metadata *metadata.Instant) {
	cl.metadata = metadata
}
