package zenon

import (
	zmq "github.com/pebbe/zmq4"
	"github.com/zenon-network/go-zenon/common/types"
	"github.com/zenon-wiki/go-zdk/client"
	"github.com/zenon-wiki/go-zdk/zdk"
	"log"
)

type Event struct {
	MessageType string `json:"type"`
	Id          string `json:"id"`
}

type Votes struct {
	Yes     int `json:"yes"`
	No      int `json:"no"`
	Abstain int `json:"abstain"`
}

type ProjectEventData struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	Created     int      `json:"created"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Znn         float64  `json:"znn"`
	Qsr         float64  `json:"qsr"`
	Status      int      `json:"status"`
	Votes       Votes    `json:"votes"`
	Phases      []string `json:"phases"`
}

type NewProject struct {
	MessageType string           `json:"type"`
	Id          string           `json:"id"`
	Data        ProjectEventData `json:"data"`
}

type ProjectStatusUpdated struct {
	MessageType string     `json:"type"`
	Id          string     `json:"id"`
	Pid         types.Hash `json:"pid"`
	OldStatus   int        `json:"old"`
	NewStatus   int        `json:"new"`
}

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

func CreateZenonZdk(url string) *zdk.Zdk {
	rpc, err := client.NewClient(url)

	if err != nil {
		log.Fatal(err)
	}
	return zdk.NewZdk(rpc)
}
