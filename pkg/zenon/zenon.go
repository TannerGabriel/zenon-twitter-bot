package zenon

import (
	"github.com/zenon-wiki/go-zdk/client"
	"github.com/zenon-wiki/go-zdk/zdk"
	"log"
)

// CreateZenonZdk creates a Zenon zdk using the provided URL
func CreateZenonZdk(url string) *zdk.Zdk {
	rpc, err := client.NewClient(url)

	if err != nil {
		log.Fatal(err)
	}
	return zdk.NewZdk(rpc)
}
