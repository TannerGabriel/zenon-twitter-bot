package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
)

type TwitterAuth struct {
	ApiKey            string
	ApiKeySecret      string
	AccessToken       string
	AccessTokenSecret string
}

// GetClient is a helper function that will return a twitter client
func GetClient(auth TwitterAuth) (*twitter.Client, error) {
	config := oauth1.NewConfig(auth.ApiKey, auth.ApiKeySecret)
	token := oauth1.NewToken(auth.AccessToken, auth.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Tweet is used create a Twitter tweet with the passed message
func Tweet(client twitter.Client, message string) error {
	_, resp, err := client.Statuses.Update(message, nil)
	if err != nil || resp.Status != "200" {
		log.Printf("error when tweeting: %e", err)
		return err
	}
	return nil
}
