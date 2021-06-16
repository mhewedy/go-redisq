# go-redisq

A very minimal queuing library for Golang backed by blocking list in Redis.

## Usage :
1. Declare queue:
    ```go
    import "github.com/mhewedy/go-redisq"
  
    var messagesQueue = redisq.NewQueue("messages.queue", func() *redis.Client {
        return redisClient // of type *redis.Client, could be null at time of registration but should be resolve to value before using the queue.
    })
    ```
2. Register listener for the queue:
    ```go
    func init() {
 	    // the queue object, the handler function and the message type
        redisq.OnMessage(messagesQueue, handleMessages, &MyDTO{})
    }

    // the listener function:
    func handleMessages(i interface{}) error {
 	    myDTO := i.(*MyDTO)
    
 	    // ..... process myDTO
 	
 	    return nil
    }
    ```
3. Publish messages to the queue:
    ```go
    if err := redisq.Publish(messagesQueue, myDto); err != nil {
        return err
    }
    ```
## Run the test:
To run tests, you need to start redis on `localhost:6379`
