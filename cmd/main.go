package main

import (
	"encoding/json"
	"fmt"
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

	zenonClient := zenon.CreateZenonZdk(os.Getenv("ZENON_URL"))
	subscriber, err := zenon.CreateZmqClient(os.Getenv("ZMQ_URL"), "")
	defer subscriber.Close()
	if err != nil {
		log.Fatalf("Error creating zmq client: %e", err)
	}

	// Read the data from the subscription
	for {
		//  Read message contents
		content, _ := subscriber.Recv(0)

		data := &zenon.Event{}
		if err := json.Unmarshal([]byte(content), data); err != nil {
			log.Println("No supported Zenon event. Skipping message!")
			log.Println("Error: ", err)
			continue
		}

		if data.MessageType == "project:new" {
			log.Println("Handling new project event")
			newProject := &zenon.NewProject{}
			if err := json.Unmarshal([]byte(content), newProject); err != nil {
				log.Println("Could not cast content to project:new event")
				continue
			}

			tweetMessage := fmt.Sprintf(`New A-Z project: %s
								
								Requested funds:
								%f ZNN %f QSR

								%s`,
				newProject.Data.Name, newProject.Data.Znn, newProject.Data.Qsr, newProject.Data.Url,
			)

			if err := twitter.Tweet(*twitterClient, tweetMessage); err != nil {
				log.Printf("Error while sending Tweet: %e", err)
			}
		} else if data.MessageType == "phase:status-update" {
			log.Println("Handling phase status update event")
			statusUpdate := &zenon.ProjectStatusUpdated{}
			if err := json.Unmarshal([]byte(content), statusUpdate); err != nil {
				log.Println("Could not cast content to phase:status-update event")
				continue
			}

			project, err := zenonClient.Embedded.Accelerator.GetProjectById(statusUpdate.Pid)
			if err != nil {
				log.Println("Error while fetching project: ", err)
			}

			tweetMessage := fmt.Sprintf(`Project has been accepted: %s
								
								Votes:
								Yes %d, No %d

								%s`,
				project.Name, project.Votes.Yes, project.Votes.No, project.Url,
			)

			if err := twitter.Tweet(*twitterClient, tweetMessage); err != nil {
				log.Printf("Error while sending Tweet: %e", err)
			}
		}
	}
}
