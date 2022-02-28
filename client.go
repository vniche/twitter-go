package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	API_URL     = "https://api.twitter.com"
	API_V2_PATH = "/2"
)

// Credentials stands for the Twitter API credentails
type Credentials struct {
	APIKey       string
	APIKeySecret string
	Bearer       string
}

// Client stands for the TWitter API HTTP client
type Client struct {
	credentials *Credentials
	BaseURL     *url.URL
	UserAgent   string

	httpClient *http.Client
}

func new() (*Client, error) {
	baseURL, err := url.ParseRequestURI(API_URL)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: http.DefaultClient,
		BaseURL:    baseURL,
	}, nil
}

// WithAPIKey returns a new instance of Twitter API v2 http client with api key and api key secret based authentication
func WithAPIKey(apiKey string, apiKeySecret string) (*Client, error) {
	return nil, fmt.Errorf("WithAPIKey not yet implemented")
}

func (client *Client) SetAPIKey(apiKey string, apiKeySecret string) {
	client.credentials = &Credentials{
		APIKey:       apiKey,
		APIKeySecret: apiKeySecret,
	}
}

// WithBearerToken returns a new instance of Twitter API v2 http client with bearer token (app-only) based authentication
func WithBearerToken(bearerToken string) (*Client, error) {
	client, err := new()
	if err != nil {
		return nil, err
	}

	client.SetBearerToken(bearerToken)

	return client, nil
}

func (client *Client) SetBearerToken(bearerToken string) {
	client.credentials = &Credentials{
		Bearer: bearerToken,
	}
}

func (client *Client) buildRequest(method, path string, body interface{}) (*http.Request, error) {
	// parses request path
	splittedPath := strings.Split(path, "?")
	apiPath := splittedPath[0]

	// assembles request info
	rel := &url.URL{Path: API_V2_PATH + apiPath}
	apiURL := client.BaseURL.ResolveReference(rel)
	var jsonBytes []byte
	var err error
	if body != nil {
		jsonBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, apiURL.String(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	if method == "GET" {
		var params url.Values
		if len(splittedPath) > 1 {
			params, err = url.ParseQuery(splittedPath[1])
			if err != nil {
				return nil, err
			}
		}

		if params.Encode() != "" {
			req.URL.RawQuery = params.Encode()
		}
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.credentials.Bearer)

	return req, nil
}

// StatusCode stands for the resp.Status code index
const StatusCode = 0

func (client *Client) do(ctx context.Context, req *http.Request) ([]byte, error) {
	req = req.WithContext(ctx)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	statusInfo := strings.Split(resp.Status, " ")

	// test for response status code
	status, err := strconv.Atoi(statusInfo[StatusCode])
	if err != nil {
		return nil, err
	} else if status < 200 || status > 299 {
		return nil, fmt.Errorf("request failed with status %d", status)
	}

	return bodyBytes, nil
}

type QueryParameters map[string][]string
