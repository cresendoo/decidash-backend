package xredis

import radix "github.com/mediocregopher/radix/v3"

func Publish(pool *radix.Pool, channel string, message string) error {
	return pool.Do(radix.Cmd(nil, "PUBLISH", channel, message))
}
