package main

/// Records example demonstrates how to set DNS records for a domain that's using Namecheap DNS

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

	resp, err := client.DomainsDNS.SetHosts(ctx, &namecheap.DomainsDNSSetHostsArgs{
		Domain: "example.com",
		Records: []namecheap.DomainsDNSHostRecord{
			{
				HostName:   "@",
				RecordType: namecheap.RecordTypeA,
				Address:    "127.0.0.1",
				TTL:        60,
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if !resp.DomainDNSSetHostsResult.IsSuccess {
		log.Fatalln("Failed to update hosts", resp.DomainDNSSetHostsResult.String())
	}
}
