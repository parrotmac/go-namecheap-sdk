package namecheap

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsDNSSetHosts(t *testing.T) {
	fakeResponse := `
		<?xml version="1.0" encoding="utf-8"?>
		<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
			<Errors />
			<Warnings />
			<RequestedCommand>namecheap.domains.dns.sethosts</RequestedCommand>
			<CommandResponse Type="namecheap.domains.dns.setHosts">
				<DomainDNSSetHostsResult Domain="domain.net" EmailType="MX" IsSuccess="true">
					<Warnings />
				</DomainDNSSetHostsResult>
			</CommandResponse>
			<Server>PHX01SBAPIEXT05</Server>
			<GMTTimeDifference>--4:00</GMTTimeDifference>
			<ExecutionTime>0.854</ExecutionTime>
		</ApiResponse>
	`

	t.Run("request_command", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.DomainsDNS.SetHosts(context.TODO(), &DomainsDNSSetHostsArgs{
			Domain: "domain.net",
		})
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, "namecheap.domains.dns.setHosts", sentBody.Get("Command"))
	})

	t.Run("request_data_correct_args_mapping", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.DomainsDNS.SetHosts(context.TODO(), &DomainsDNSSetHostsArgs{
			Domain:    "domain.net",
			EmailType: EmailTypeForward,
			Flag:      UInt8(100),
			Tag:       "issue",
		})
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, "domain", sentBody.Get("SLD"))
		assert.Equal(t, "net", sentBody.Get("TLD"))
		assert.Equal(t, "FWD", sentBody.Get("EmailType"))
		assert.Equal(t, "100", sentBody.Get("Flag"))
		assert.Equal(t, "issue", sentBody.Get("Tag"))
	})

	t.Run("request_data_correct_mx_records_mapping", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.DomainsDNS.SetHosts(context.TODO(), &DomainsDNSSetHostsArgs{
			Domain:    "domain.net",
			EmailType: "MX",
			Records: []DomainsDNSHostRecord{
				{
					RecordType: RecordTypeA,
					HostName:   "@",
					Address:    "10.11.12.13",
					TTL:        1800,
				},
				{
					RecordType: RecordTypeMX,
					HostName:   "mail",
					Address:    "super-mail.com",
					TTL:        1800,
					MXPref:     UInt8(10),
				},
			},
		})
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, RecordTypeA, sentBody.Get("RecordType1"))
		assert.Equal(t, "@", sentBody.Get("HostName1"))
		assert.Equal(t, "10.11.12.13", sentBody.Get("Address1"))
		assert.Equal(t, "1800", sentBody.Get("TTL1"))

		assert.Equal(t, RecordTypeMX, sentBody.Get("RecordType2"))
		assert.Equal(t, "mail", sentBody.Get("HostName2"))
		assert.Equal(t, "super-mail.com", sentBody.Get("Address2"))
		assert.Equal(t, "1800", sentBody.Get("TTL2"))
		assert.Equal(t, "10", sentBody.Get("MXPref2"))
	})

	t.Run("request_data_correct_mxe_records_mapping", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.DomainsDNS.SetHosts(context.TODO(), &DomainsDNSSetHostsArgs{
			Domain:    "domain.net",
			EmailType: EmailTypeMXE,
			Records: []DomainsDNSHostRecord{
				{
					RecordType: RecordTypeMXE,
					HostName:   "mail",
					Address:    "10.11.12.13",
					TTL:        1800,
					MXPref:     UInt8(10),
				},
			},
		})
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, RecordTypeMXE, sentBody.Get("RecordType1"))
		assert.Equal(t, "mail", sentBody.Get("HostName1"))
		assert.Equal(t, "10.11.12.13", sentBody.Get("Address1"))
		assert.Equal(t, "1800", sentBody.Get("TTL1"))
		assert.Equal(t, "10", sentBody.Get("MXPref1"))
	})

	t.Run("request_data_correct_url_record", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.DomainsDNS.SetHosts(context.TODO(), &DomainsDNSSetHostsArgs{
			Domain: "domain.net",
			Records: []DomainsDNSHostRecord{
				{
					RecordType: RecordTypeURL,
					HostName:   "redirect",
					Address:    "https://domain.com",
				},
			},
		})
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, RecordTypeURL, sentBody.Get("RecordType1"))
		assert.Equal(t, "redirect", sentBody.Get("HostName1"))
		assert.Equal(t, "https://domain.com", sentBody.Get("Address1"))
	})

	t.Run("request_data_correct_url301_record", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.DomainsDNS.SetHosts(context.TODO(), &DomainsDNSSetHostsArgs{
			Domain: "domain.net",
			Records: []DomainsDNSHostRecord{
				{
					RecordType: RecordTypeURL301,
					HostName:   "redirect",
					Address:    "https://domain.com",
				},
			},
		})
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, RecordTypeURL301, sentBody.Get("RecordType1"))
		assert.Equal(t, "redirect", sentBody.Get("HostName1"))
		assert.Equal(t, "https://domain.com", sentBody.Get("Address1"))
	})

	t.Run("request_data_correct_frame_record", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.DomainsDNS.SetHosts(context.TODO(), &DomainsDNSSetHostsArgs{
			Domain: "domain.net",
			Records: []DomainsDNSHostRecord{
				{
					RecordType: RecordTypeFrame,
					HostName:   "redirect",
					Address:    "https://domain.com",
				},
			},
		})
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, RecordTypeFrame, sentBody.Get("RecordType1"))
		assert.Equal(t, "redirect", sentBody.Get("HostName1"))
		assert.Equal(t, "https://domain.com", sentBody.Get("Address1"))
	})

	t.Run("request_data_correct_CAA_iodef_record", func(t *testing.T) {
		var sentBody url.Values

		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			body, _ := ioutil.ReadAll(request.Body)
			query, _ := url.ParseQuery(string(body))
			sentBody = query
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		_, err := client.DomainsDNS.SetHosts(context.TODO(), &DomainsDNSSetHostsArgs{
			Domain: "domain.net",
			Records: []DomainsDNSHostRecord{
				{
					RecordType: RecordTypeCAA,
					HostName:   "@",
					Address:    "0 iodef http://domain.com",
				},
			},
		})
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, RecordTypeCAA, sentBody.Get("RecordType1"))
		assert.Equal(t, "@", sentBody.Get("HostName1"))
		assert.Equal(t, "0 iodef http://domain.com", sentBody.Get("Address1"))
	})

	var errorCases = []struct {
		Name          string
		Args          *DomainsDNSSetHostsArgs
		ExpectedError string
	}{
		{
			Name: "request_data_error_incorrect_domain",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "dom",
			},
			ExpectedError: "invalid domain: incorrect format",
		},
		{
			Name: "request_data_error_bad_email_type",
			Args: &DomainsDNSSetHostsArgs{
				Domain:    "domain.net",
				EmailType: "BAD_TYPE",
			},
			ExpectedError: "invalid EmailType value: BAD_TYPE",
		},
		{
			Name: "request_data_error_bad_tag",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Tag:    "BAD_TAG",
			},
			ExpectedError: "invalid Tag value: BAD_TAG",
		},
		{
			Name: "request_data_error_no_hostname",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: "CNAME", Address: "domain.com"},
				},
			},
			ExpectedError: "Records[0].HostName is required",
		},
		{
			Name: "request_data_error_no_recordtype",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{HostName: "@", Address: "domain.com"},
				},
			},
			ExpectedError: "Records[0].RecordType is required",
		},
		{
			Name: "request_data_error_bad_recordtype",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: "BAD", HostName: "@", Address: "domain.com"},
				},
			},
			ExpectedError: "invalid Records[0].RecordType value: BAD",
		},
		{
			Name: "request_data_error_too_low_ttl",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeCNAME, HostName: "@", Address: "domain.com", TTL: 59},
				},
			},
			ExpectedError: "invalid Records[0].TTL value: 59",
		},
		{
			Name: "request_data_error_too_big_ttl",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeCNAME, HostName: "@", Address: "domain.com", TTL: 60_001},
				},
			},
			ExpectedError: "invalid Records[0].TTL value: 60001",
		},
		{
			Name: "request_data_error_no_address",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeCNAME, HostName: "@"},
				},
			},
			ExpectedError: "Records[0].Address is required",
		},
		{
			Name: "request_data_error_email_type_mx_without_records",
			Args: &DomainsDNSSetHostsArgs{
				EmailType: EmailTypeMX,
				Domain:    "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeCNAME, HostName: "@", Address: "domain.com", TTL: 1800},
				},
			},
			ExpectedError: "minimum 1 MX record required for MX EmailType",
		},
		{
			Name: "request_data_error_email_type_mxe_without_record",
			Args: &DomainsDNSSetHostsArgs{
				EmailType: EmailTypeMXE,
				Domain:    "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeCNAME, HostName: "@", Address: "domain.com", TTL: 1800},
				},
			},
			ExpectedError: "one MXE record required for MXE EmailType",
		},
		{
			Name: "request_data_error_email_type_nil_with_mx",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeMX, HostName: "mail", Address: "mail.domain.com", MXPref: UInt8(10)},
				},
			},
			ExpectedError: "Records[0].RecordType MX is not allowed for blank EmailType",
		},
		{
			Name: "request_data_error_email_type_fwd_with_mx",
			Args: &DomainsDNSSetHostsArgs{
				Domain:    "domain.net",
				EmailType: EmailTypeForward,
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeMX, HostName: "mail", Address: "mail.domain.com", MXPref: UInt8(10)},
				},
			},
			ExpectedError: "Records[0].RecordType MX is not allowed for EmailType=FWD",
		},
		{
			Name: "request_data_error_email_type_nil_with_mxe",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeMXE, HostName: "mail", Address: "10.11.12.13"},
				},
			},
			ExpectedError: "Records[0].RecordType MXE is not allowed for blank EmailType",
		},
		{
			Name: "request_data_error_email_type_fwd_with_mxe",
			Args: &DomainsDNSSetHostsArgs{
				Domain:    "domain.net",
				EmailType: EmailTypeForward,
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeMXE, HostName: "mail", Address: "10.11.12.13"},
				},
			},
			ExpectedError: "Records[0].RecordType MXE is not allowed for EmailType=FWD",
		},
		{
			Name: "request_data_error_two_mxe_records",
			Args: &DomainsDNSSetHostsArgs{
				Domain:    "domain.net",
				EmailType: EmailTypeMXE,
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeMXE, HostName: "mail", Address: "10.11.12.13"},
					{RecordType: RecordTypeMXE, HostName: "mail2", Address: "10.11.12.14"},
				},
			},
			ExpectedError: "one MXE record required for MXE EmailType",
		},
		{
			Name: "request_data_error_no_mxpref_for_mx_record",
			Args: &DomainsDNSSetHostsArgs{
				Domain:    "domain.net",
				EmailType: EmailTypeMX,
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeMX, HostName: "mail", Address: "mail.domain.com"},
				},
			},
			ExpectedError: "Records[0].MXPref is nil but required for MX record type",
		},
		{
			Name: "request_data_error_no_protocol_prefix_for_url_record",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeURL, HostName: "mail", Address: "domain.com"},
				},
			},
			ExpectedError: `Records[0].Address "domain.com" must contain a protocol prefix for URL record`,
		},
		{
			Name: "request_data_error_no_protocol_prefix_for_url301_record",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeURL301, HostName: "mail", Address: "domain.com"},
				},
			},
			ExpectedError: `Records[0].Address "domain.com" must contain a protocol prefix for URL301 record`,
		},
		{
			Name: "request_data_error_no_protocol_prefix_for_frame_record",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeFrame, HostName: "mail", Address: "domain.com"},
				},
			},
			ExpectedError: `Records[0].Address "domain.com" must contain a protocol prefix for FRAME record`,
		},
		{
			Name: "request_data_error_no_protocol_prefix_for_caa_iodef_record",
			Args: &DomainsDNSSetHostsArgs{
				Domain: "domain.net",
				Records: []DomainsDNSHostRecord{
					{RecordType: RecordTypeCAA, HostName: "@", Address: "0 iodef domain.com"},
				},
			},
			ExpectedError: `Records[0].Address "0 iodef domain.com" must contain a protocol prefix for CAA iodef record`,
		},
	}

	for _, errorCase := range errorCases {
		t.Run(errorCase.Name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte(fakeResponse))
			}))
			defer mockServer.Close()

			client := setupClient(nil)
			client.BaseURL = mockServer.URL

			_, err := client.DomainsDNS.SetHosts(context.TODO(), errorCase.Args)

			assert.EqualError(t, err, errorCase.ExpectedError)
		})
	}
}
