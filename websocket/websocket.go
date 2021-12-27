package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/itering/substrate-api-rpc/pkg/recws"
)

var (
	wsEndPoint = ""
	maxCap     = 25
)

type WsConn interface {
	Dial(urlStr string, reqHeader http.Header)
	IsConnected() bool
	Close()
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, message []byte, err error)
	WriteJSON(v interface{}) error
	ReadJSON(v interface{}) error
	MarkUnusable()
	CloseAndReconnect()
}

func Init(options ...Option) (*PoolConn, error) {
	var err error
	if wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{
				KeepAliveTimeout: 10 * time.Second,
				WriteTimeout:     time.Second,
				ReadTimeout:      time.Second * 2,
				NonVerbose:       true,
				HandshakeTimeout: time.Second}
			SubscribeConn.Dial(wsEndPoint, nil)
			for _, o := range options {
				o.Apply(SubscribeConn)
			}
			return SubscribeConn, err
		}
		if wsPool, err = NewChannelPool(1, maxCap, factory); err != nil {
			fmt.Println("NewChannelPool", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := wsPool.Get()
	return conn, err
}

func SetEndpoint(endpoint string) {
	wsEndPoint = endpoint
}

// SetChannelPoolMaxCap set connection pool max cap
func SetChannelPoolMaxCap(max int) {
	maxCap = max
}

func Close() {
	if wsPool != nil {
		wsPool.Close()
	}
}

type Option interface {
	Apply(*recws.RecConn)
}

type OptionFunc func(*recws.RecConn)

func (f OptionFunc) Apply(conn *recws.RecConn) {
	f(conn)
}

func WithHandshakeTimeout(t time.Duration) Option {
	return OptionFunc(func(m *recws.RecConn) {
		m.HandshakeTimeout = t
	})
}

func WithWriteTimeoutTimeout(t time.Duration) Option {
	return OptionFunc(func(m *recws.RecConn) {
		m.WriteTimeout = t
	})
}

func WithReadTimeoutTimeout(t time.Duration) Option {
	return OptionFunc(func(m *recws.RecConn) {
		m.ReadTimeout = t
	})
}
