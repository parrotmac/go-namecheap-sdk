package namecheap

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/namecheap/go-namecheap-sdk/v2/namecheap/internal/syncretry"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

const (
	namecheapProductionApiUrl = "https://api.namecheap.com/xml.response"
	namecheapSandboxApiUrl    = "https://api.sandbox.namecheap.com/xml.response"
)

type ClientOptions struct {
	UserName   string
	ApiUser    string
	ApiKey     string
	ClientIp   string
	UseSandbox bool
}

type Client struct {
	http   *http.Client
	common service
	sr     *syncretry.SyncRetry

	ClientOptions *ClientOptions
	BaseURL       string

	Domains    DomainsService
	DomainsDNS DomainsDNSService
}

type service struct {
	client *Client
}

func LookupClientIP(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://ipv4.icanhazip.com", nil)
	if err != nil {
		return "", fmt.Errorf("failed to build HTTP request %v", err)
	}

	c := cleanhttp.DefaultClient()

	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch external IP %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("did not get OK status when looking up external IP: %d/%s", resp.StatusCode, resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body when looking up external IP: %v", err)
	}
	return strings.TrimSpace(string(respBody)), nil
}

func NewClientOptionsFromEnv(clientIP string) *ClientOptions {
	return &ClientOptions{
		UserName:   os.Getenv("NAMECHEAP_USERNAME"),
		ApiUser:    os.Getenv("NAMECHEAP_API_USER"),
		ApiKey:     os.Getenv("NAMECHEAP_API_KEY"),
		ClientIp:   clientIP,
		UseSandbox: strings.ToLower(os.Getenv("NAMECHEAP_USE_SANDBOX")) != "false",
	}
}

// NewClient returns a new Namecheap API Client
func NewClient(options *ClientOptions) *Client {
	client := &Client{
		ClientOptions: options,
		http:          cleanhttp.DefaultClient(),
		sr:            syncretry.NewSyncRetry(&syncretry.Options{Delays: []int{1, 5, 15, 30, 50}}),
	}

	client.BaseURL = namecheapProductionApiUrl
	if options.UseSandbox {
		client.BaseURL = namecheapSandboxApiUrl
	}

	client.common.client = client
	client.Domains = (DomainsService)(client.common)
	client.DomainsDNS = (DomainsDNSService)(client.common)

	return client
}

// NewRequest creates a new request with the params
func (c *Client) NewRequest(ctx context.Context, body map[string]string) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL)

	if err != nil {
		return nil, fmt.Errorf("Error parsing base URL: %s", err)
	}

	body["Username"] = c.ClientOptions.UserName
	body["ApiKey"] = c.ClientOptions.ApiKey
	body["ApiUser"] = c.ClientOptions.ApiUser
	body["ClientIp"] = c.ClientOptions.ClientIp

	rBody := encodeBody(body)

	// Build the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBufferString(rBody))

	if err != nil {
		return nil, fmt.Errorf("Error creating request: %s", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(rBody)))

	return req, nil
}

func (c *Client) DoXML(ctx context.Context, body map[string]string, obj interface{}) (*http.Response, error) {
	var requestResponse *http.Response
	err := c.sr.Do(ctx, func() error {
		request, err := c.NewRequest(ctx, body)
		if err != nil {
			return err
		}

		response, err := c.http.Do(request)
		if err != nil {
			return err
		}

		if response.StatusCode == 405 {
			return syncretry.RetryError
		}

		requestResponse = response
		defer response.Body.Close()

		err = decodeBody(response.Body, obj)
		return err
	})

	if err != nil && errors.Is(err, syncretry.RetryAttemptsError) {
		return nil, fmt.Errorf("API retry limit exceeded")
	}

	return requestResponse, err
}

// decodeBody decodes the interface from received XML
func decodeBody(reader io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(reader)
	err := decoder.Decode(&obj)
	if err != nil {
		return fmt.Errorf("unable to parse server response: %s", err)
	}
	return nil
}

// encodeBody converts the map into query string
func encodeBody(body map[string]string) string {
	data := url.Values{}
	for key, val := range body {
		data.Set(key, val)
	}
	return data.Encode()
}

// ParseDomain is a wrapper around publicsuffix.Parse to throw the correct error
func ParseDomain(domain string) (*publicsuffix.DomainName, error) {
	const regDomainString = `^([\-a-zA-Z0-9]+\.+){1,}[a-zA-Z0-9]+$`
	regDomain, err := regexp.Compile(regDomainString)
	if err != nil {
		return nil, err
	}

	if !regDomain.MatchString(domain) {
		return nil, fmt.Errorf("invalid domain: incorrect format")
	}

	parsedDomain, err := publicsuffix.Parse(domain)
	if err != nil {
		return nil, fmt.Errorf("invalid domain: %v", err)
	}

	return parsedDomain, nil
}

// UInt8 is a helper routine that allocates a new uint8 value
// to store v and returns a pointer to it.
func UInt8(v uint8) *uint8 { return &v }
