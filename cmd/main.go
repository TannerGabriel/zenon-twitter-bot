package main

import (
	"encoding/json"
	"github.com/TannerGabriel/zenon-twitter-bot/pkg/twitter"
	"github.com/TannerGabriel/zenon-twitter-bot/pkg/zenon"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	auth := twitter.TwitterAuth{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ApiKey:            os.Getenv("API_KEY"),
		ApiKeySecret:      os.Getenv("API_KEY_SECRET"),
	}

	twitterClient, err := twitter.GetClient(auth)

	if err != nil {
		log.Fatalf("Error creating twitter client: %e", err)
	}

	log.Println(twitterClient)

	subscriber, err := zenon.CreateZmqClient(os.Getenv("ZMQ_URL"), "")
	defer subscriber.Close()
	if err != nil {
		log.Fatalf("Error creating zmq client: %e", err)
	}

	// Read the data from the subscription
	for {
		//  Read message contents
		content, _ := subscriber.Recv(0)
		log.Println(content)

		data := &zenon.Event{}
		err := json.Unmarshal([]byte(content), data)
		if err != nil {
			log.Println("No supported Zenon event. Skipping message!")
			log.Println("Error: ", err)
			continue
		}

		log.Println(data)

		if data.MessageType == "project:new" {
			// TODO: Handle new project
			log.Println("Handling new project")
		} else if data.MessageType == "phase:status-update" {
			// TODO: Handle phase updated
			log.Println("Handling phase status update")
		}
	}
}
