package namecheap

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type DomainCreateResponse struct {
	XMLName xml.Name `xml:"ApiResponse"`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainCreateCommandResponse `xml:"CommandResponse"`
}

type DomainCreateCommandResponse struct {
	DomainCreateResult DomainCreateResult `xml:"DomainCreateResult"`
}

/*
Domain	Domain name that you are trying to register.
Registered	Possible responses: True, False. Indicates whether the domain was registered.
ChargedAmount	Total amount charged for registration.
DomainID	Unique integer value that represents the domain.
OrderID	Unique integer value that represents the order.
TransactionID	Unique integer value that represents the transaction.
WhoisguardEnable	Possible responses: True, False. Indicates whether domain privacy protection is enabled for the domain.
NonRealTimeDomain	Possible responses: True, False. Indicates whether the domain registration is instant (real-time) or not.
*/
type DomainCreateResult struct {
	Domain        string `xml:"Domain"`
	Registered    string `xml:"Registered"`
	ChargedAmount string `xml:"ChargedAmount"`
	DomainID      string `xml:"DomainID"`
	OrderID       string `xml:"OrderID"`
	TransactionID string `xml:"TransactionID"`
	Whoisguard    string `xml:"Whoisguard"`
	NonRealTime   string `xml:"NonRealTime"`
}

/*
DomainName	String	70	Yes	Domain name to register
Years	Number	2	Yes	Number of years to register -- Default Value: 2
PromotionCode	String	20	No	Promotional (coupon) code for the domain

RegistrantOrganizationName	String	255	No	Organization of the Registrant user
RegistrantJobTitle	String	255	No	Job title of the Registrant user
RegistrantFirstName	String	255	Yes	First name of the Registrant user
RegistrantLastName	String	255	Yes	Second name of the Registrant user
RegistrantAddress1	String	255	Yes	Address1 of the Registrant user
RegistrantAddress2	String	255	No	Address2 of the Registrant user
RegistrantCity	String	50	Yes	City of the Registrant user
RegistrantStateProvince	String	50	Yes	State/Province of the Registrant user
RegistrantStateProvinceChoice	String	50	No	StateProvinceChoice of the Registrant user
RegistrantPostalCode	String	50	Yes	PostalCode of the Registrant user
RegistrantCountry	String	50	Yes	Country of the Registrant user
RegistrantPhone	String	50	Yes	Phone number in the format +NNN.NNNNNNNNNN
RegistrantPhoneExt	String	50	No	PhoneExt of the Registrant user
RegistrantFax	String	50	No	Fax number in the format +NNN.NNNNNNNNNN
RegistrantEmailAddress	String	255	Yes	Email address of the Registrant user

TechOrganizationName	String	255	No	Organization of the Tech user
TechJobTitle	String	255	No	Job title of the Tech user
TechFirstName	String	255	Yes	First name of the Tech user
TechLastName	String	255	Yes	Second name of the Tech user
TechAddress1	String	255	Yes	Address1 of the Tech user
TechAddress2	String	255	No	Address2 of the Tech user
TechCity	String	50	Yes	City of the Tech user
TechStateProvince	String	50	Yes	State/Province of the Tech user
TechStateProvinceChoice	String	50	No	StateProvinceChoice of the Tech user
TechPostalCode	String	50	Yes	PostalCode of the Tech user
TechCountry	String	50	Yes	Country of the Tech user
TechPhone	String	50	Yes	Phone number in the format +NNN.NNNNNNNNNN
TechPhoneExt	String	50	No	PhoneExt of the Tech user
TechFax	String	50	No	Fax number in the format +NNN.NNNNNNNNNN
TechEmailAddress	String	255	Yes	Email address of the Tech user

AdminOrganizationName	String	255	No	Organization of the Admin user
AdminJobTitle	String	255	No	Job title of the Admin user
AdminFirstName	String	255	Yes	First name of the Admin user
AdminLastName	String	255	Yes	Second name of the Admin user
AdminAddress1	String	255	Yes	Address1 of the Admin user
AdminAddress2	String	255	No	Address2 of the Admin user
AdminCity	String	50	Yes	City of the Admin user
AdminStateProvince	String	50	Yes	State/Province of the Admin user
AdminStateProvinceChoice	String	50	No	StateProvinceChoice of the Admin user
AdminPostalCode	String	50	Yes	PostalCode of the Admin user
AdminCountry	String	50	Yes	Country of the Admin user
AdminPhone	String	50	Yes	Phone number in the format +NNN.NNN NNN NNNN
AdminPhoneExt	String	50	No	PhoneExt of the Admin user
AdminFax	String	50	No	Fax number in the format +NNN.NNNNNNNNNN
AdminEmailAddress	String	255	Yes	Email address of the Admin user

AuxBillingOrganizationName	String	255	No	Organization of the AuxBilling user
AuxBillingJobTitle	String	255	No	Job title of the AuxBilling user
AuxBillingFirstName	String	255	Yes	First name of the AuxBilling user
AuxBillingLastName	String	255	Yes	Second name of the AuxBilling user
AuxBillingAddress1	String	255	Yes	Address1 of the AuxBilling user
AuxBillingAddress2	String	255	No	Address2 of the AuxBilling user
AuxBillingCity	String	50	Yes	City of the AuxBilling user
AuxBillingStateProvince	String	50	Yes	State/Province of the AuxBilling user
AuxBillingStateProvinceChoice	String	50	No	StateProvinceChoice of the AuxBilling user
AuxBillingPostalCode	String	50	Yes	PostalCode of the AuxBilling user
AuxBillingCountry	String	50	Yes	Country of the AuxBilling user
AuxBillingPhone	String	50	Yes	Phone number in the format +NNN.NNNNNNNNNN
AuxBillingPhoneExt	String	50	No	PhoneExt of the AuxBilling user
AuxBillingFax	String	50	No	Fax number in the format +NNN.NNNNNNNNNN
AuxBillingEmailAddress	String	255	Yes	Email address of the AuxBilling user

BillingFirstName	String	255	No	First name of the Billing user
BillingLastName	String	255	No	Second name of the Billing user
BillingAddress1	String	255	No	Address1 of the Billing user
BillingAddress2	String	255	No	Address2 of the Billing user
BillingCity	String	50	No	City of the Billing user
BillingStateProvince	String	50	No	State/Province of the Billing user
BillingStateProvinceChoice	String	50	No	StateProvinceChoice of the Billing user
BillingPostalCode	String	50	No	PostalCode of the Billing user
BillingCountry	String	50	No	Country of the Billing user
BillingPhone	String	50	No	Phone number in the format +NNN.NNNNNNNNNN
BillingPhoneExt	String	50	No	PhoneExt of the Billing user
BillingFax	String	50	No	Fax number in the format +NNN.NNNNNNNNNN
BillingEmailAddress	String	255	No	Email address of the Billing user

IdnCode	String	100	No	Code of Internationalized Domain Name (please refer to the note below)
Extended attributes	String n/a Yes	Required for .us, .eu, .ca, .co.uk, .org.uk, .me.uk, .nu , .com.au, .net.au, .org.au, .es, .nom.es, .com.es, .org.es, .de, .fr TLDs only
Nameservers	String n/a No Comma-separated list of custom nameservers to be associated with the domain name
AddFreeWhoisguard	String	10	No	Adds free domain privacy for the domain -- Default Value: no
WGEnabled	String	10	No	Enables free domain privacy for the domain -- Default Value: no
IsPremiumDomain	Boolean	10	No	Indication if the domain name is premium
PremiumPrice	Currency	20	No	Registration price for the premium domain
EapFee	Currency	20	No	Purchase fee for the premium domain during Early Access Program (EAP)*
*/
type DomainCreateArgs struct {
	DomainName    string // Required
	Years         int    // Required
	PromotionCode string // Optional

	RegistrantOrganizationName    string // Optional
	RegistrantJobTitle            string // Optional
	RegistrantFirstName           string // Required
	RegistrantLastName            string // Required
	RegistrantAddress1            string // Required
	RegistrantAddress2            string // Optional
	RegistrantCity                string // Required
	RegistrantStateProvince       string // Required
	RegistrantStateProvinceChoice string // Optional
	RegistrantPostalCode          string // Required
	RegistrantCountry             string // Required
	RegistrantPhone               string // Required
	RegistrantPhoneExt            string // Optional
	RegistrantFax                 string // Optional
	RegistrantEmailAddress        string // Required

	TechOrganizationName    string // Optional
	TechJobTitle            string // Optional
	TechFirstName           string // Required
	TechLastName            string // Required
	TechAddress1            string // Required
	TechAddress2            string // Optional
	TechCity                string // Required
	TechStateProvince       string // Required
	TechStateProvinceChoice string // Optional
	TechPostalCode          string // Required
	TechCountry             string // Required
	TechPhone               string // Required
	TechPhoneExt            string // Optional
	TechFax                 string // Optional
	TechEmailAddress        string // Required

	AdminOrganizationName    string // Optional
	AdminJobTitle            string // Optional
	AdminFirstName           string // Required
	AdminLastName            string // Required
	AdminAddress1            string // Required
	AdminAddress2            string // Optional
	AdminCity                string // Required
	AdminStateProvince       string // Required
	AdminStateProvinceChoice string // Optional
	AdminPostalCode          string // Required
	AdminCountry             string // Required
	AdminPhone               string // Required
	AdminPhoneExt            string // Optional
	AdminFax                 string // Optional
	AdminEmailAddress        string // Required

	AuxBillingOrganizationName    string // Optional
	AuxBillingJobTitle            string // Optional
	AuxBillingFirstName           string // Required
	AuxBillingLastName            string // Required
	AuxBillingAddress1            string // Required
	AuxBillingAddress2            string // Optional
	AuxBillingCity                string // Required
	AuxBillingStateProvince       string // Required
	AuxBillingStateProvinceChoice string // Optional
	AuxBillingPostalCode          string // Required
	AuxBillingCountry             string // Required
	AuxBillingPhone               string // Required
	AuxBillingPhoneExt            string // Optional
	AuxBillingFax                 string // Optional
	AuxBillingEmailAddress        string // Required

	BillingFirstName           string // Optional
	BillingLastName            string // Optional
	BillingAddress1            string // Optional
	BillingAddress2            string // Optional
	BillingCity                string // Optional
	BillingStateProvince       string // Optional
	BillingStateProvinceChoice string // Optional
	BillingPostalCode          string // Optional
	BillingCountry             string // Optional
	BillingPhone               string // Optional
	BillingPhoneExt            string // Optional
	BillingFax                 string // Optional
	BillingEmailAddress        string // Optional

	IdnCode            string // Optional
	ExtendedAttributes string // Required -- Check docs!
	Nameservers        string // Optional
	AddFreeWhoisguard  string // Optional
	WGEnabled          string // Optional
	IsPremiumDomain    bool   // Optional
	PremiumPrice       string // Optional
	EapFee             string // Optional
}

// Registers a given domain name
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/check/
func (ds *DomainsService) Create(ctx context.Context, args DomainCreateArgs) (*DomainCreateResult, error) {
	var resp DomainCreateResponse

	params := domainCreateArgsToParams(args)
	params["Command"] = "namecheap.domains.create"

	_, err := ds.client.DoXML(ctx, params, &resp)
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

	return &resp.CommandResponse.DomainCreateResult, apiErr
}

func domainCreateArgsToParams(args DomainCreateArgs) map[string]string {
	params := map[string]string{
		"DomainName": args.DomainName,
		"Years":      strconv.Itoa(args.Years),

		"RegistrantFirstName":     args.RegistrantFirstName,
		"RegistrantLastName":      args.RegistrantLastName,
		"RegistrantAddress1":      args.RegistrantAddress1,
		"RegistrantCity":          args.RegistrantCity,
		"RegistrantStateProvince": args.RegistrantStateProvince,
		"RegistrantPostalCode":    args.RegistrantPostalCode,
		"RegistrantCountry":       args.RegistrantCountry,
		"RegistrantPhone":         args.RegistrantPhone,
		"RegistrantEmailAddress":  args.RegistrantEmailAddress,

		"TechFirstName":     args.TechFirstName,
		"TechLastName":      args.TechLastName,
		"TechAddress1":      args.TechAddress1,
		"TechCity":          args.TechCity,
		"TechStateProvince": args.TechStateProvince,
		"TechPostalCode":    args.TechPostalCode,
		"TechCountry":       args.TechCountry,
		"TechPhone":         args.TechPhone,
		"TechEmailAddress":  args.TechEmailAddress,

		"AdminFirstName":     args.AdminFirstName,
		"AdminLastName":      args.AdminLastName,
		"AdminAddress1":      args.AdminAddress1,
		"AdminCity":          args.AdminCity,
		"AdminStateProvince": args.AdminStateProvince,
		"AdminPostalCode":    args.AdminPostalCode,
		"AdminCountry":       args.AdminCountry,
		"AdminPhone":         args.AdminPhone,
		"AdminEmailAddress":  args.AdminEmailAddress,

		"AuxBillingFirstName":     args.AuxBillingFirstName,
		"AuxBillingLastName":      args.AuxBillingLastName,
		"AuxBillingAddress1":      args.AuxBillingAddress1,
		"AuxBillingCity":          args.AuxBillingCity,
		"AuxBillingStateProvince": args.AuxBillingStateProvince,
		"AuxBillingPostalCode":    args.AuxBillingPostalCode,
		"AuxBillingCountry":       args.AuxBillingCountry,
		"AuxBillingPhone":         args.AuxBillingPhone,
		"AuxBillingEmailAddress":  args.AuxBillingEmailAddress,

		"IsPremiumDomain": strconv.FormatBool(args.IsPremiumDomain),
	}

	if args.PromotionCode != "" {
		params["PromotionCode"] = args.PromotionCode
	}

	if args.RegistrantOrganizationName != "" {
		params["RegistrantOrganizationName"] = args.RegistrantOrganizationName
	}

	if args.RegistrantJobTitle != "" {
		params["RegistrantJobTitle"] = args.RegistrantJobTitle
	}

	if args.RegistrantAddress2 != "" {
		params["RegistrantAddress2"] = args.RegistrantAddress2
	}

	if args.RegistrantStateProvinceChoice != "" {
		params["RegistrantStateProvinceChoice"] = args.RegistrantStateProvinceChoice
	}

	if args.RegistrantPhoneExt != "" {
		params["RegistrantPhoneExt"] = args.RegistrantPhoneExt
	}

	if args.RegistrantFax != "" {
		params["RegistrantFax"] = args.RegistrantFax
	}

	if args.TechOrganizationName != "" {
		params["TechOrganizationName"] = args.TechOrganizationName
	}

	if args.TechJobTitle != "" {
		params["TechJobTitle"] = args.TechJobTitle
	}

	if args.TechAddress2 != "" {
		params["TechAddress2"] = args.TechAddress2
	}

	if args.TechStateProvinceChoice != "" {
		params["TechStateProvinceChoice"] = args.TechStateProvinceChoice
	}

	if args.TechPhoneExt != "" {
		params["TechPhoneExt"] = args.TechPhoneExt
	}

	if args.TechFax != "" {
		params["TechFax"] = args.TechFax
	}

	if args.AdminOrganizationName != "" {
		params["AdminOrganizationName"] = args.AdminOrganizationName
	}

	if args.AdminJobTitle != "" {
		params["AdminJobTitle"] = args.AdminJobTitle
	}

	if args.AdminAddress2 != "" {
		params["AdminAddress2"] = args.AdminAddress2
	}

	if args.AdminStateProvinceChoice != "" {
		params["AdminStateProvinceChoice"] = args.AdminStateProvinceChoice
	}

	if args.AdminPhoneExt != "" {
		params["AdminPhoneExt"] = args.AdminPhoneExt
	}

	if args.AdminFax != "" {
		params["AdminFax"] = args.AdminFax
	}

	if args.AuxBillingOrganizationName != "" {
		params["AuxBillingOrganizationName"] = args.AuxBillingOrganizationName
	}

	if args.AuxBillingJobTitle != "" {
		params["AuxBillingJobTitle"] = args.AuxBillingJobTitle
	}

	if args.AuxBillingAddress2 != "" {
		params["AuxBillingAddress2"] = args.AuxBillingAddress2
	}

	if args.AuxBillingStateProvinceChoice != "" {
		params["AuxBillingStateProvinceChoice"] = args.AuxBillingStateProvinceChoice
	}

	if args.AuxBillingPhoneExt != "" {
		params["AuxBillingPhoneExt"] = args.AuxBillingPhoneExt
	}

	if args.AuxBillingFax != "" {
		params["AuxBillingFax"] = args.AuxBillingFax
	}

	if args.AuxBillingEmailAddress != "" {
		params["AuxBillingEmailAddress"] = args.AuxBillingEmailAddress
	}

	if args.IdnCode != "" {
		params["IdnCode"] = args.IdnCode
	}

	if args.Nameservers != "" {
		params["Nameservers"] = args.Nameservers
	}

	if args.AddFreeWhoisguard != "" {
		params["AddFreeWhoisguard"] = args.AddFreeWhoisguard
	}

	if args.WGEnabled != "" {
		params["WGEnabled"] = args.WGEnabled
	}

	if args.PremiumPrice != "" {
		params["PremiumPrice"] = args.PremiumPrice
	}

	if args.EapFee != "" {
		params["EapFee"] = args.EapFee
	}

	return params
}
