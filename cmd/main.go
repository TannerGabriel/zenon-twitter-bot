package main

import (
	"context"
	"encoding/json"
	"github.com/TannerGabriel/zenon-twitter-bot/pkg/twitter"
	"github.com/TannerGabriel/zenon-twitter-bot/pkg/zenon"
	"github.com/go-zeromq/zmq4"
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

	zenonClient := zenon.CreateZenonZdk(os.Getenv("ZENON_URL"))

	//  Prepare our subscriber
	subscriber := zmq4.NewSub(context.Background())
	defer subscriber.Close()

	err = subscriber.Dial(os.Getenv("ZMQ_URL"))
	if err != nil {
		log.Fatalf("could not dial: %v", err)
	}

	err = subscriber.SetOption(zmq4.OptionSubscribe, "")
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	log.Println("Twitter bot started! Waiting for messages.")

	// Read the data from the subscription
	for {
		//  Read message contents
		content, err := subscriber.Recv()
		if err != nil {
			log.Fatalf("could not receive message: %v", err)
		}

		data := &zenon.Event{}
		if err := json.Unmarshal(content.Frames[0], data); err != nil {
			log.Println("No supported Zenon event. Skipping message!")
			log.Println("Error: ", err)
			continue
		}

		var tweetMessage string

		if data.MessageType == "project:new" {
			log.Println("Handling new project event")
			newProject := &zenon.ProjectNew{}
			if err := json.Unmarshal(content.Frames[0], newProject); err != nil {
				log.Println("Could not cast content to project:new event")
				continue
			}

			tweetMessage = zenon.HandleNewProject(*newProject)
		} else if data.MessageType == "phase:status-update" {
			log.Println("Handling phase status update event")
			statusUpdate := &zenon.PhaseStatusUpdate{}
			if err := json.Unmarshal(content.Frames[0], statusUpdate); err != nil {
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
			if err := json.Unmarshal(content.Frames[0], statusUpdate); err != nil {
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
