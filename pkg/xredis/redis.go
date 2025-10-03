package xredis

import (
	"time"

	"github.com/cresendoo/decidash-backend/pkg/errorx"
	radix "github.com/mediocregopher/radix/v3"
)

func NewRedisPool(addr string, pool, db int, auth string) (*radix.Pool, error) {
	connFunc := radix.DefaultConnFunc
	if auth != "" {
		connFunc = func(network, addr string) (radix.Conn, error) {
			return radix.Dial(network, addr,
				radix.DialTimeout(1*time.Minute),
				radix.DialAuthPass(auth),
			)
		}
	}
	redisPool, err := radix.NewPool("tcp", addr, pool, radix.PoolConnFunc(connFunc))
	if err != nil {
		return nil, errorx.Wrap(err)
	}
	return redisPool, nil
}

type RedisPubSub struct {
	pubSub radix.PubSubConn
}

func NewRedisPubSub(network, address, auth string, db int) *RedisPubSub {
	var config = []radix.DialOpt{
		radix.DialTimeout(10 * time.Second),
		radix.DialSelectDB(db),
	}
	if auth != "" {
		config = append(config, radix.DialAuthPass(auth))
	}
	connFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr, config...)
	}
	return &RedisPubSub{pubSub: radix.PersistentPubSub("tcp", address, connFunc)}
}

func (r *RedisPubSub) Subscribe(msgCh chan<- radix.PubSubMessage, channels ...string) error {
	return r.pubSub.Subscribe(msgCh, channels...)
}

func (r *RedisPubSub) Unsubscribe(msgCh chan<- radix.PubSubMessage, channels ...string) error {
	return r.pubSub.Unsubscribe(msgCh, channels...)
}

func (r *RedisPubSub) PSubscribe(msgCh chan<- radix.PubSubMessage, patterns ...string) error {
	return r.pubSub.PSubscribe(msgCh, patterns...)
}

func (r *RedisPubSub) PUnsubscribe(msgCh chan<- radix.PubSubMessage, patterns ...string) error {
	return r.pubSub.PUnsubscribe(msgCh, patterns...)
}

func (r *RedisPubSub) Ping() error {
	return r.pubSub.Ping()
}

func (r *RedisPubSub) Close() error {
	return r.pubSub.Close()
}
