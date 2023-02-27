package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/namecheap/go-namecheap-sdk/v2/namecheap"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

const (
	ansiReset         = "\033[0m"
	ansiStrikethrough = "\033[9m"
)

func ansiWrap(s string, with string) string {
	return fmt.Sprintf("%s%s%s", with, s, ansiReset)
}

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

	domainsToCheck := []string{
		"namecheap.com",
		"abc.xyz",
		"example.com",
		"example.invalid",
		"some-domain-somebody-could-own.net",
		"some-domain-somebody-could-own.xyz",
	}

	availabilityResults, err := client.Domains.Check(ctx, domainsToCheck)
	if err != nil {
		fmt.Println("Failed to check availability:", err)
	}

	if availabilityResults == nil {
		return
	}

	for _, domain := range availabilityResults {
		if domain.Available != "true" {
			fmt.Printf("%s is not available\n", domain.Domain)
			continue
		}
		if domain.IsPremiumName == "true" {
			fmt.Printf("[PREMIUM] Price for %s: %s\n", domain.Domain, domain.PremiumRegistrationPrice)
		} else {
			domainInfo, err := publicsuffix.Parse(domain.Domain)
			if err != nil {
				fmt.Printf("Failed to parse TLD: %s\n", err)
				continue
			}

			pricing, err := client.UsersService.GetPricing(ctx, namecheap.UserGetPricingArgs{
				ProductType:     namecheap.ProductTypeDomain,
				ProductCategory: namecheap.ProductCategoryDomains,
				ProductName:     domainInfo.TLD,
				ActionName:      namecheap.ActionNameRegister,
			})
			if err != nil {
				log.Fatalln(err)
			}
			p := (*pricing).ProductType.ProductCategory[0].Product[0].Price[0]
			if p.RegularPrice == p.YourPrice {
				fmt.Printf("[REGULAR] Price for %s: %s\n", domainInfo.String(), p.RegularPrice)
			} else {
				fmt.Printf("[SALE] Price for %s: %s -> %s\n", domainInfo.String(), ansiWrap(p.RegularPrice, ansiStrikethrough), p.YourPrice)
			}
		}
	}
}
