# go-redisq

Simple Redis queue lib for golang backed by blocking list

## Usage :
1. Declare queue:
    ```go
    import "github.com/mhewedy/go-redisq"
  
    var messagesQueue = redisq.NewQueue("messages.queue", func() *redis.Client {
        return redisClient // of type *redis.Client
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
 	
 	    // Note: the returned error will not be handled and will only be logged
        // if you wish to implement special handling (log to db or so), then 
        // you need to implement the logic yourself
 	    return nil
    }
    ```
3. Publish messages to the queue:
    ```go
    if err := redisq.Publish(messagesQueue, myDto); err != nil {
        return err
    }
    ```
