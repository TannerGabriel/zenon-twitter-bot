package main

import (
	"encoding/json"
	"github.com/TannerGabriel/zenon-twitter-bot/pkg/twitter"
	"github.com/TannerGabriel/zenon-twitter-bot/pkg/zenon"
	"github.com/joho/godotenv"
	czmq "gopkg.in/zeromq/goczmq.v4"
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

	zenonClient := zenon.CreateZenonZdk(os.Getenv("ZENON_URL"))
	subscriber, err := czmq.NewSub(os.Getenv("ZMQ_URL"), "")
	defer subscriber.Destroy()
	if err != nil {
		log.Fatalf("Error creating zmq client: %e", err)
	}

	// Read the data from the subscription
	for {
		//  Read message contents
		content, _, err := subscriber.RecvFrame()

		if err != nil {
			log.Printf("Error while receiving frame: %e", err)
			continue
		}

		data := &zenon.Event{}
		if err := json.Unmarshal(content, data); err != nil {
			log.Println("No supported Zenon event. Skipping message!")
			log.Println("Error: ", err)
			continue
		}

		var tweetMessage string

		if data.MessageType == "project:new" {
			log.Println("Handling new project event")
			newProject := &zenon.ProjectNew{}
			if err := json.Unmarshal(content, newProject); err != nil {
				log.Println("Could not cast content to project:new event")
				continue
			}

			tweetMessage = zenon.HandleNewProject(*newProject)
		} else if data.MessageType == "phase:status-update" {
			log.Println("Handling phase status update event")
			statusUpdate := &zenon.PhaseStatusUpdate{}
			if err := json.Unmarshal([]byte(content), statusUpdate); err != nil {
				log.Println("Could not cast content to phase:status-update event")
				continue
			}

			tweetMessage, err = zenon.HandlePhaseStatusUpdate(*statusUpdate, zenonClient)
			if err != nil {
				log.Printf("Error while handling phase:status-update event: %e", err)
				continue
			}
		} else if data.MessageType == "project:status-update" {
			log.Println("Handling project status update event")
			statusUpdate := &zenon.ProjectStatusUpdate{}
			if err := json.Unmarshal(content, statusUpdate); err != nil {
				log.Println("Could not cast content to project:status-update event")
				continue
			}

			tweetMessage, err = zenon.HandleProjectStatusUpdate(*statusUpdate, zenonClient)
			if err != nil {
				log.Printf("Error while handling project:status-update event: %e", err)
				continue
			}
		}

		if tweetMessage != "" {
			if err := twitter.Tweet(*twitterClient, tweetMessage); err != nil {
				log.Printf("Error while sending Tweet: %e", err)
			}
		}
	}
}
