package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
)

var (
	username   = os.Getenv("NAMECHEAP_USERNAME")
	apiUser    = os.Getenv("NAMECHEAP_API_USER")
	apiKey     = os.Getenv("NAMECHEAP_API_KEY")
	useSandbox = strings.ToLower(os.Getenv("NAMECHEAP_USE_SANDBOX")) != "false"
)

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFunc()

	// Namecheap API requires us to specify the IP we're calling from
	// Note that this IP must also be registered as allowed in the Namecheap dashboard
	clientIP, err := namecheap.LookupClientIP(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	client := namecheap.NewClient(
		namecheap.NewClientOptionsFromEnv(
			clientIP,
		),
	)

	res, err := client.Domains.Check(ctx, []string{"namecheap.com", "might-be-available-i-guess.com", "unsupported.tld"})
	if err != nil {
		log.Println("error while checking domains:", err)
	}
	if res == nil {
		return
	}
	log.Printf("Result: %+v\n", res)
}
