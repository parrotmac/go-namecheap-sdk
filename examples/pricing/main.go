package main

import (
	"context"
	"log"
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

	res, err := client.UsersService.GetPricing(ctx, namecheap.UserGetPricingArgs{
		ProductType:     namecheap.ProductTypeDomain,
		ProductCategory: namecheap.ProductCategoryDomains,

		// Optionally, specify the TLD. Omit to list all TLDs (though, it can take quite some time to perform this query).
		// ProductName:     "co",

		ActionName: namecheap.ActionNameRegister,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Result: %+v\n", res)
}
