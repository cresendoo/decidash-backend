package xredis

import (
	"fmt"
	"log"
	"testing"
	"time"

	radix "github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/assert"
)

func TestRedisPubSub(t *testing.T) {
	pool, err := NewRedisPool("localhost:6379", 10, 0, "")
	assert.NoError(t, err)
	pubSub := NewRedisPubSub("tcp", "localhost:6379", "", 0)
	assert.NotNil(t, pubSub)

	ch := make(chan radix.PubSubMessage)
	err = pubSub.Subscribe(ch, "test")
	assert.NoError(t, err)

	go func() {
		time.Sleep(1 * time.Second)
		err = Publish(pool, "test", "test")
		assert.NoError(t, err)
	}()

	msg := <-ch
	fmt.Println(msg)
	assert.Equal(t, "test", string(msg.Message))
}

func TestRedisPubSub_Subscribe(t *testing.T) {
	pubSub := NewRedisPubSub("tcp", "localhost:6379", "", 0)
	assert.NotNil(t, pubSub)
	ch := make(chan radix.PubSubMessage)
	err := pubSub.Subscribe(ch, "test")
	assert.NoError(t, err)

	msg := <-ch
	log.Println(msg)
	assert.Equal(t, "test", string(msg.Message))
}

func TestRedisPubSub_Publish(t *testing.T) {
	pool, err := NewRedisPool("localhost:6379", 10, 0, "")
	assert.NoError(t, err)
	err = Publish(pool, "test", "test")
	assert.NoError(t, err)
	log.Println("publish")
}
