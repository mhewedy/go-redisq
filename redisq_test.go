package redisq

import (
	"context"
	"errors"
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

	queue := NewQueue("TestInvalidMessagePayload", func() *redis.Client {
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

func TestDeadLetterQueue(t *testing.T) {

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	queue := NewQueue("TestDeadLetterQueue", func() *redis.Client {
		return redisClient
	})

	_ = Publish(queue, msg{Data: "Hello"})

	OnMessage(queue, func(i interface{}) error {
		return errors.New("failed to process message")
	}, &msg{})

	// -----

	var c = new(int)
	dlQueue := NewQueue("TestDeadLetterQueue.dl.queue", func() *redis.Client {
		return redisClient
	})

	OnMessage(dlQueue, func(i interface{}) error {

		data := i.(*msg).Data
		fmt.Println("found in dl queue: ", data)

		*c++

		return nil

	}, &msg{})

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, *c)
}
