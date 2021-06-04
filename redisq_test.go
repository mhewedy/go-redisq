package redisq

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

type msg struct {
	Data string
}

func TestInvalidMessagePayload(t *testing.T) {

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	queue := NewQueue("myqueue", func() *redis.Client {
		return redisClient
	})

	redisClient.RPush(context.Background(), "myqueue", "invalid payload")

	OnMessage(queue, func(i interface{}) error {

		data := i.(*msg).Data
		fmt.Println("receiving: ", data)

		return nil

	}, &msg{})

	time.Sleep(2 * time.Second)
}
