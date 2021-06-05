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

type HandlerFun func(i interface{}) error

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

// OnMessage calls the passed HandlerFunc
// and in case the HandlerFunc function returned an error,
// it logs the error to stdout and push the original message
// to a queue with name <original_queue>.dl.queue
func OnMessage(q *Queue, f HandlerFun, messageType interface{}) {
	go onMessage(q, f, messageType)
}

func onMessage(q *Queue, f HandlerFun, messageType interface{}) {
	func() {

		// to handle panic cases from inside the HandlerFunc
		// in such case, we start a new goroutine
		defer func() {
			if err := recover(); err != nil {
				logError(err)
				go onMessage(q, f, messageType)
			}
		}()

		// waiting for q.clientFn() to be available (assigned).
		//
		// this introduced to make it okay to register HandlerFunc in packages init method
		// as usually the redis client is being available in the main package
		// which started after all package init methods are being executed.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for {
			if q.clientFn() != nil || ctx.Err() != nil {
				break
			}
		}

		// start the listener loop
		for {
			var err error

			result, err := q.clientFn().BLPop(context.Background(), 0*time.Second, q.Name).Result()
			if err != nil {
				logError(fmt.Errorf("%v (%v)", err, result))
				continue
			}
			if len(result) < 2 {
				logError(fmt.Errorf("redis result call doesn't return valid respones %s\n", result))
				continue
			}
			if err = json.Unmarshal([]byte(result[1]), &messageType); err != nil {
				logError(fmt.Errorf("%v (%v)", err, result[1]))
				continue
			}
			if err = f(messageType); err != nil {
				logError(fmt.Errorf("HandlerFun error: %v (%v)\n", err, result[1]))

				publishToDLQueue(q, result)
				continue
			}
		}
	}()
}

func publishToDLQueue(q *Queue, result []string) {

	dlQueueName := fmt.Sprintf("%s.dl.queue", q.Name)
	if err := q.clientFn().RPush(context.Background(), dlQueueName, result[1]).Err(); err != nil {
		logError(fmt.Errorf("%v (%v)", err, result[1]))
	}
}

func logError(err interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, "[go-redisq]:", err)
}
