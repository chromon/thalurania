package pubsub

import (
	"chalurania/service/log"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

// redis 连接池
type RedisPool struct {
	Pool *redis.Pool
}

func NewRedisPool(addr string, db int, pwd string) *RedisPool {
	pool := &redis.Pool{
		MaxIdle: 10,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr, redis.DialPassword(pwd), redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	log.Info.Printf("new redis pool at %s", addr)

	redisPool := &RedisPool{
		Pool: pool,
	}
	return redisPool
}

func (r *RedisPool) Publish(channel, message string) (int, error) {
	c := r.Pool.Get()
	defer c.Close()
	n, err := redis.Int(c.Do("PUBLISH", channel, message))
	if err != nil {
		return 0, fmt.Errorf("redis publish %s %s, err: %v", channel, message, err)
	}
	return n, nil
}

func (r *RedisPool) Subscribe(ctx context.Context, consume func(redis.Message) error, channel ...string) error {
	psc := redis.PubSubConn{Conn: r.Pool.Get()}

	log.Info.Printf("redis pubsub subscribe channel: %v", channel)
	if err := psc.Subscribe(redis.Args{}.AddFlat(channel)...); err != nil {
		return err
	}
	done := make(chan error, 1)
	// start a new goroutine to receive message
	go func() {
		defer psc.Close()
		for {
			switch msg := psc.Receive().(type) {
			case error:
				done <- fmt.Errorf("redis pubsub receive err: %v", msg)
				return
			case redis.Message:
				if err := consume(msg); err != nil {
					done <- err
					return
				}
			case redis.Subscription:
				if msg.Count == 0 {
					// all channels are unsubscribed
					done <- nil
					return
				}
			}
		}
	}()

	// health check
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			if err := psc.Unsubscribe(); err != nil {
				return fmt.Errorf("redis pubsub unsubscribe err: %v", err)
			}
			return nil
		case err := <-done:
			return err
		case <-tick.C:
			if err := psc.Ping(""); err != nil {
				return err
			}
		}
	}

	//return nil
}

func (r *RedisPool) Close() {
	err := r.Pool.Close()
	if err != nil {
		log.Error.Println("close redis pool err:", err)
	}
}
