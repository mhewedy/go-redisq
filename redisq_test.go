package redisq

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
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

	redisClient.RPush(context.Background(), "myqueue", "invalid payload") // error
	_ = Publish(queue, "xyz")                                             // error
	_ = Publish(queue, msg{Data: "Hello"})                                // success

	var c = new(int)

	OnMessage(queue, func(i interface{}) error {

		data := i.(*msg).Data
		fmt.Println("receiving: ", data)

		*c++

		return nil

	}, &msg{})

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, *c)
}
