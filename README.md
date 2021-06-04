# go-redisq

Simple Redis queue lib for golang backed by blocking list

## Usage :
1. Declare queue:
    ```go
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
    func handleMessages(i interface{}, err error) {
    
    	if err != nil {
    		log.Error(i, err)
    		return
    	}
    
    	myDto := i.(*MyDTO)
    
    	// .....
    }
    ```
3. Publish messages to the queue:
    ```go
    if err := redisq.Publish(messagesQueue, myDto); err != nil {
        return err
    }
    ```
