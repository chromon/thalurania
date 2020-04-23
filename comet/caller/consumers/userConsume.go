package consumers

import "github.com/gomodule/redigo/redis"

func UserConsume(msg redis.Message) error {
	return nil
}