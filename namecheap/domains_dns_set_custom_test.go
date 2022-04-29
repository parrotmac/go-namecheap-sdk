package namecheap

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainsDNSSetCustom(t *testing.T) {
	fakeResponse := `<?xml version="1.0" encoding="utf-8"?>
		<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
			<Errors />
			<Warnings />
			<RequestedCommand>namecheap.domains.dns.setcustom</RequestedCommand>
			<CommandResponse Type="namecheap.domains.dns.setCustom">
				<DomainDNSSetCustomResult Domain="domain.net" Updated="true" />
			</CommandResponse>
			<Server>PHX01SBAPIEXT06</Server>
			<GMTTimeDifference>--4:00</GMTTimeDifference>
			<ExecutionTime>2.599</ExecutionTime>
		</ApiResponse>`

	fakeNameservers := []string{"dns1.nameserver.com", "dns2.nameserver.com"}

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

		_, err := client.DomainsDNS.SetCustom(context.TODO(), "domain.net", fakeNameservers)
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, "namecheap.domains.dns.setCustom", sentBody.Get("Command"))
	})

	t.Run("request_data_domain", func(t *testing.T) {
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

		_, err := client.DomainsDNS.SetCustom(context.TODO(), "domain.net", fakeNameservers)
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, "net", sentBody.Get("TLD"))
		assert.Equal(t, "domain", sentBody.Get("SLD"))
	})

	t.Run("request_data_nameservers", func(t *testing.T) {
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

		_, err := client.DomainsDNS.SetCustom(context.TODO(), "domain.net", fakeNameservers)
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		expectedNameservers := strings.Join(fakeNameservers, ",")

		assert.Equal(t, expectedNameservers, sentBody.Get("Nameservers"))
	})

	t.Run("correct_parsing_result_attributes", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(fakeResponse))
		}))
		defer mockServer.Close()

		client := setupClient(nil)
		client.BaseURL = mockServer.URL

		result, err := client.DomainsDNS.SetCustom(context.TODO(), "domain.net", fakeNameservers)
		if err != nil {
			t.Fatal("Unable to get domains", err)
		}

		assert.Equal(t, "domain.net", result.DomainDNSSetCustomResult.Domain)
		assert.Equal(t, true, result.DomainDNSSetCustomResult.Updated)
	})

	errorCases := []struct {
		Nameservers []string
	}{
		{Nameservers: []string{}},
		{Nameservers: []string{"name.server"}},
	}

	for _, errorCase := range errorCases {
		t.Run("request_data_error_"+strconv.Itoa(len(errorCase.Nameservers))+"_nameservers", func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte(fakeResponse))
			}))
			defer mockServer.Close()

			client := setupClient(nil)
			client.BaseURL = mockServer.URL

			_, err := client.DomainsDNS.SetCustom(context.TODO(), "domain.net", errorCase.Nameservers)

			assert.EqualError(t, err, "invalid nameservers: must contain minimum two items")
		})
	}
}
