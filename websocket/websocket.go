package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/itering/substrate-api-rpc/pkg/recws"
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

func Init(nodeName NodeName, options ...Option) (*PoolConn, error) {
	if nodeName == "" {
		nodeName = DefaultNodeName
	}

	var (
		err    error
		maxCap = 25
	)
	if _, ok := wsPool[nodeName]; !ok {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{
				KeepAliveTimeout: 10 * time.Second,
				WriteTimeout:     time.Second,
				ReadTimeout:      time.Second * 2,
				NonVerbose:       true,
				HandshakeTimeout: time.Second}
			for _, o := range options {
				o.Apply(SubscribeConn)
			}

			SubscribeConn.Dial(SubscribeConn.EndPoint, nil)

			if SubscribeConn.MaxPoolCap > 0 {
				maxCap = SubscribeConn.MaxPoolCap
			}

			return SubscribeConn, err
		}
		if wsPool[nodeName], err = NewChannelPool(1, maxCap, factory); err != nil {
			fmt.Println("NewChannelPool", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := wsPool[nodeName].Get()
	return conn, err
}

func Close(nodeName NodeName) {
	if pool, ok := wsPool[nodeName]; ok && pool != nil {
		pool.Close()
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

func WithReadBufferSize(size int) Option {
	return OptionFunc(func(m *recws.RecConn) {
		m.ReadBufferSize = size
	})
}

func WithWriteBufferSize(size int) Option {
	return OptionFunc(func(m *recws.RecConn) {
		m.WriteBufferSize = size
	})
}

func WithEndPoint(endPoint string) Option {
	return OptionFunc(func(m *recws.RecConn) {
		m.EndPoint = endPoint
	})
}

func WithMaxPoolCap(cap int) Option {
	return OptionFunc(func(m *recws.RecConn) {
		m.MaxPoolCap = cap
	})
}
