package main

import (
	"chalurania/service/pubsub"
	"context"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	r := pubsub.NewRedisPool("127.0.0.1:6379", 0, "")
	consume := func(msg redis.Message) error {
		log.Printf("recv msg: %s", msg.Data)
		return nil
	}
	for i := 0; i <10; i++{
		log.Printf("-------------- %d -----------------", i)
		ctx, cancel := context.WithCancel(context.Background())
		go func(){
			if err := r.Subscribe(ctx, consume, "channel"); err != nil {
				log.Println("subscribe err:", err)
			}
		}()
		time.Sleep(time.Second)
		_, err:= r.Publish("channel", "hello, " + strconv.Itoa(i))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		cancel()
	}
	forever := make(chan struct{})
	<-forever
}