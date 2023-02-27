package main

/// Nameservers example demonstrates how to list information about registered domains, including the configured nameservers.

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
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
		nameservers = append(nameservers, dnsRecords.DomainDNSGetListResult.Nameservers...)

		fmt.Printf("%s: %s [Auto-Renew: %v] [Expired: %v] [Nameservers: %s]\n", domain.ID, domain.Name, domain.AutoRenew, domain.IsExpired, strings.Join(nameservers, ", "))
	}
}
