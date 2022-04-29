package main

/// Nameservers example demonstrates how to list information about registered domains, including the configured nameservers.

import (
	"context"
	"fmt"
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

	domainListing, err := client.Domains.GetList(ctx, &namecheap.DomainsGetListArgs{
		// Filter on `.com` domains
		SearchTerm: ".com",
	})
	if err != nil {
		log.Fatalln("couldn't list domains", err)
	}

	for _, domain := range domainListing.Domains {
		var nameservers []string
		dnsRecords, err := client.DomainsDNS.GetList(ctx, domain.Name)
		if err != nil {
			log.Fatalln(err)
		}
		for _, rec := range dnsRecords.DomainDNSGetListResult.Nameservers {
			nameservers = append(nameservers, rec)
		}
		fmt.Printf("%s: %s [Auto-Renew: %v] [Expired: %v] [Nameservers: %s]\n", domain.ID, domain.Name, domain.AutoRenew, domain.IsExpired, strings.Join(nameservers, ", "))
	}
}
