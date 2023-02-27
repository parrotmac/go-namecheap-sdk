package namecheap

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
)

type DomainCheckResponse struct {
	XMLName xml.Name `xml:"ApiResponse"`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainCheckCommandResponse `xml:"CommandResponse"`
}

type DomainCheckCommandResponse struct {
	DomainCheckResult []DomainCheckResult `xml:"DomainCheckResult"`
}

type DomainCheckResult struct {
	Domain                   string `xml:"Domain,attr"`                   // Domain name for which you wish to check availability
	Available                string `xml:"Available,attr"`                // Indicates whether the domain name is available for registration
	IsPremiumName            string `xml:"IsPremiumName,attr"`            // Indicates whether the domain name is premium
	PremiumRegistrationPrice string `xml:"PremiumRegistrationPrice,attr"` // Registration Price for the premium domain
	PremiumRenewalPrice      string `xml:"PremiumRenewalPrice,attr"`      // Renewal price for the premium domain
	PremiumRestorePrice      string `xml:"PremiumRestorePrice,attr"`      // Restore price for the premium domain
	PremiumTransferPrice     string `xml:"PremiumTransferPrice,attr"`     // Transfer price for the premium domain
	IcannFee                 string `xml:"IcannFee,attr"`                 // Fee charged by ICANN
	EapFee                   string `xml:"EapFee,attr"`                   // Purchase fee for the premium domain during Early Access Program (EAP)*

}

// Checks the availability of domains
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/check/
func (ds *DomainsService) Check(ctx context.Context, domainList []string) ([]DomainCheckResult, error) {
	var checkResponse DomainCheckResponse
	params := map[string]string{
		"Command":    "namecheap.domains.check",
		"DomainList": strings.Join(domainList, ","),
	}

	_, err := ds.client.DoXML(ctx, params, &checkResponse)
	if err != nil {
		return nil, err
	}

	var apiErr error
	if checkResponse.Errors != nil && len(checkResponse.Errors) > 0 {
		errMessages := []string{}
		for _, e := range checkResponse.Errors {
			errMessages = append(errMessages, fmt.Sprintf("%s (%s)", e.Message, e.Number))
		}
		apiErr = fmt.Errorf("%s", strings.Join(errMessages, "; "))
	}

	return checkResponse.CommandResponse.DomainCheckResult, apiErr
}
