package namecheap

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
)

// API Docs: https://www.namecheap.com/support/api/methods/users/get-pricing/

/*
ProductType		|	Product Category	|	ActionName							|	ProductName
-----------------------------------------------------------------------------------------------
DOMAIN			|	DOMAINS				|	REGISTER,RENEW,REACTIVATE,TRANSFER	| 	COM
SSLCERTIFICATE	|	COMODO				|	PURCHASE,RENEW						| 	INSTANTSSL
*/

const (
	ProductTypeDomain         = "DOMAIN"
	ProductTypeSSLCertificate = "SSLCERTIFICATE"

	ProductCategoryDomains = "DOMAINS"
	ProductCategoryComodo  = "COMODO"

	ActionNameRegister   = "REGISTER"
	ActionNameRenew      = "RENEW"
	ActionNameReactivate = "REACTIVATE"
	ActionNameTransfer   = "TRANSFER"

	ActionNamePurchaseSSL = "PURCHASE"
	ActionNameRenewSSL    = "RENEW"

	ProductNameDomain     = "COM"
	ProductNameInstantSSL = "INSTANTSSL"
)

/*
ProductType	String	20	Yes	Product Type to get pricing information
ProductCategory	String	20	No	Specific category within a product type
PromotionCode	String	20	No	Promotional (coupon) code for the user
ActionName	String	20	No	Specific action within a product type
ProductName	String	20	No	The name of the product within a product type
*/
type UserGetPricingArgs struct {
	ProductType     string // Required
	ProductCategory string // Optional
	PromotionCode   string // Optional
	ActionName      string // Optional
	ProductName     string // Optional
}

type UserGetPricingResponse struct {
	XMLName xml.Name `xml:"ApiResponse"`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *UserGetPricingCommandResponse `xml:"CommandResponse"`
}

type UserGetPricingCommandResponse struct {
	UserGetPricingResult UserGetPricingResult `xml:"UserGetPricingResult"`
}

/*
ProductType Name	The type of product
ProductCategory Name	Category type of the product
Product Name	Name of the product
Duration	The duration of the product
DurationType	The duration type of the product
Price	Indicates Final price (it can be from regular, userprice, special price,promo price, tier price)
RegularPrice	Indicates regular price
YourPrice	The user’s price for the product
CouponPrice	Price with coupon enabled
Currency	Currency in which the price is listed
*/
// type UserGetPricingResult struct {
// 	ProductType     ProductType     `xml:"ProductType"`
// 	ProductCategory ProductCategory `xml:"ProductCategory"`
// 	ProductName     ProductName     `xml:"ProductName"`
// 	Duration        int             `xml:"Duration"`
// 	DurationType    string          `xml:"DurationType"`
// 	Price           float64         `xml:"Price"`
// 	RegularPrice    float64         `xml:"RegularPrice"`
// 	YourPrice       float64         `xml:"YourPrice"`
// 	CouponPrice     float64         `xml:"CouponPrice"`
// 	Currency        string          `xml:"Currency"`
// }

type UserGetPricingResult struct {
	ProductType ProductType `xml:"ProductType"`
}

type ProductType struct {
	Name            string            `xml:"Name,attr"`
	ProductCategory []ProductCategory `xml:"ProductCategory"`
}

type ProductCategory struct {
	Name    string    `xml:"Name,attr"`
	Product []Product `xml:"Product"`
}

type Price struct {
	Duration     string `xml:"Duration,attr"`
	DurationType string `xml:"DurationType,attr"`
	Price        string `xml:"Price,attr"`
	RegularPrice string `xml:"RegularPrice,attr"`
	YourPrice    string `xml:"YourPrice,attr"`
	CouponPrice  string `xml:"CouponPrice,attr"`
	Currency     string `xml:"Currency,attr"`
}

type Product struct {
	Name  string  `xml:"Name,attr"`
	Price []Price `xml:"Price"`
}

// Checks the availability of domains
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/check/
func (us *UsersService) GetPricing(ctx context.Context, args UserGetPricingArgs) (*UserGetPricingResult, error) {
	var resp UserGetPricingResponse

	params := userGetPricingArgsToParams(args)
	params["Command"] = "namecheap.users.getPricing"

	_, err := us.client.DoXML(ctx, params, &resp)
	if err != nil {
		return nil, err
	}

	var apiErr error
	if resp.Errors != nil && len(resp.Errors) > 0 {
		errMessages := []string{}
		for _, e := range resp.Errors {
			errMessages = append(errMessages, fmt.Sprintf("%s (%s)", e.Message, e.Number))
		}
		apiErr = fmt.Errorf("%s", strings.Join(errMessages, "; "))
	}

	return &resp.CommandResponse.UserGetPricingResult, apiErr
}

func userGetPricingArgsToParams(args UserGetPricingArgs) map[string]string {
	return map[string]string{
		"ProductType":     string(args.ProductType),
		"ProductCategory": string(args.ProductCategory),
		"PromotionCode":   args.PromotionCode,
		"ActionName":      string(args.ActionName),
		"ProductName":     string(args.ProductName),
	}
}
