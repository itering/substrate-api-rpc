package websocket

import (
	"fmt"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"time"
)

var wsEndPoint = ""

func Init() (*PoolConn, error) {
	var err error
	if wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
			SubscribeConn.Dial(wsEndPoint, nil)
			return SubscribeConn, err
		}
		if wsPool, err = NewChannelPool(1, 10, factory); err != nil {
			fmt.Println("NewChannelPool", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := wsPool.Get()
	return conn, err
}

func RegWSEndPoint(endpoint string) {
	wsEndPoint = endpoint
}

func CloseWsConnection() {
	if wsPool != nil {
		wsPool.Close()
	}
}
