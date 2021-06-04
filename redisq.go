package redisq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

type Queue struct {
	Name     string
	clientFn func() *redis.Client
}

type HandlerFun func(i interface{}, err error)

func NewQueue(name string, clientFn func() *redis.Client) *Queue {
	return &Queue{Name: name, clientFn: clientFn}
}

func Publish(q *Queue, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return q.clientFn().RPush(context.Background(), q.Name, string(b)).Err()
}

func OnMessage(q *Queue, f HandlerFun, messageType interface{}) {
	go onMessage(q, f, messageType)
}

func onMessage(q *Queue, f HandlerFun, messageType interface{}) {
	func() {

		defer func() {
			if err := recover(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "ERROR:", err)
				go onMessage(q, f, messageType)
			}
		}()

		// waiting for qq.clientFn() to be available
		// use to call the onMessage in the init, and when initialize the Redis client in
		// the main function.
		// so unless we wait here, the init function will be called before
		// the main, and then the q.clientFn() (the redis client) will not be available
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for {
			if q.clientFn() != nil || ctx.Err() != nil {
				break
			}
		}

		for {
			var err error

			result, err := q.clientFn().BLPop(context.Background(), 0*time.Second, q.Name).Result()
			payload := result[1]

			if err = json.Unmarshal([]byte(payload), &messageType); err != nil {
				f(messageType, err)
				return
			}

			f(messageType, err)
		}
	}()
}
