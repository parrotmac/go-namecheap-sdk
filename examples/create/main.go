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

	res, err := client.Domains.Create(ctx, namecheap.DomainCreateArgs{
		DomainName: "some-domain-name-foobar.xyz",
		Years:      1,

		RegistrantFirstName:     "Jane",
		RegistrantLastName:      "Smith",
		RegistrantAddress1:      "16 Pennsylvania St",
		RegistrantCity:          "Not Lake City",
		RegistrantStateProvince: "UM",
		RegistrantPostalCode:    "02134",
		RegistrantCountry:       "US",
		RegistrantPhone:         "+1.2145550000",
		RegistrantEmailAddress:  "webmaster@example.com",

		TechFirstName:     "Jane",
		TechLastName:      "Smith",
		TechAddress1:      "16 Pennsylvania St",
		TechCity:          "Not Lake City",
		TechStateProvince: "UM",
		TechPostalCode:    "02134",
		TechCountry:       "US",
		TechPhone:         "+1.2145550000",
		TechEmailAddress:  "webmaster@example.com",

		AdminFirstName:     "Jane",
		AdminLastName:      "Smith",
		AdminAddress1:      "16 Pennsylvania St",
		AdminCity:          "Not Lake City",
		AdminStateProvince: "UM",
		AdminPostalCode:    "02134",
		AdminCountry:       "US",
		AdminPhone:         "+1.2145550000",
		AdminEmailAddress:  "webmaster@example.com",

		AuxBillingFirstName:     "Jane",
		AuxBillingLastName:      "Smith",
		AuxBillingAddress1:      "16 Pennsylvania St",
		AuxBillingCity:          "Not Lake City",
		AuxBillingStateProvince: "UM",
		AuxBillingPostalCode:    "02134",
		AuxBillingCountry:       "US",
		AuxBillingPhone:         "+1.2145550000",
		AuxBillingEmailAddress:  "webmaster@example.com",

		Nameservers: "ns1.some-dns-servers.invalid,ns2.other-nameservers-here.example",
	})
	if err != nil {
		log.Println("error while creating/registering domain:", err)
	}
	if res == nil {
		return
	}
	log.Printf("Result: %+v\n", res)
}
