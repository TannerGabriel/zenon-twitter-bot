package zenon

import (
	zmq "github.com/pebbe/zmq4"
	"log"
)

// CreateZmqClient creates a Zmq SUB and subscribes to the filtered topics
func CreateZmqClient(url string, filter string) (*zmq.Socket, error) {
	//  Prepare our subscriber
	subscriber, _ := zmq.NewSocket(zmq.SUB)

	// Subscribe to ZMQ server
	err := subscriber.Connect(url)
	if err != nil {
		log.Printf("error creating zmq client: %e", err)
		return nil, err
	}

	// Set subscribe event
	err = subscriber.SetSubscribe(filter)
	if err != nil {
		log.Fatal(err)
	}

	return subscriber, nil
}
