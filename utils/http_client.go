package utils

import (
	"fmt"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func CreateHttpClient(proxies string) (tlsClient.HttpClient, error) {
	options := []tlsClient.HttpClientOption{
		tlsClient.WithClientProfile(profiles.Chrome_133_PSK),
		tlsClient.WithRandomTLSExtensionOrder(),
		tlsClient.WithInsecureSkipVerify(),
		tlsClient.WithTimeoutSeconds(30),
	}
	if proxies != "" {
		options = append(options, tlsClient.WithProxyUrl(fmt.Sprintf("http://%s", proxies)))
	}

	client, err := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
	if err != nil {
		Logger{}.Error("Failed to create Http Client: %s", err)
		return nil, err
	}

	return client, nil
}

func CookiesToHeader(allCookies map[string][]*http.Cookie) string {
	var cookieStrs []string
	for _, cookies := range allCookies {
		for _, cookie := range cookies {
			cookieStrs = append(cookieStrs, cookie.Name+"="+cookie.Value)
		}
	}
	return strings.Join(cookieStrs, "; ")
}

// HeaderPair represents a header key-value pair
type HeaderPair struct {
	Key   string
	Value string
}

// RequestConfig contains all possible options for making a request
type RequestConfig struct {
	Method  string
	URL     string
	Body    io.Reader
	Headers []HeaderPair
}

// DefaultConfig returns common Twitter request headers and their order
func DefaultConfig() RequestConfig {
	return RequestConfig{
		Headers: []HeaderPair{
			{Key: "accept", Value: "*/*"},
			{Key: "accept-encoding", Value: "gzip, deflate, br"},
			{Key: "content-type", Value: "application/x-www-form-urlencoded"},
			{Key: "origin", Value: "https://twitter.com"},
			{Key: "sec-ch-ua-mobile", Value: "?0"},
			{Key: "sec-ch-ua-platform", Value: `"Windows"`},
			{Key: "sec-fetch-dest", Value: "empty"},
			{Key: "sec-fetch-mode", Value: "cors"},
			{Key: "sec-fetch-site", Value: "same-origin"},
			{Key: "user-agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"},
		},
	}
}

// MakeRequest handles HTTP requests with proper header ordering and error handling
func MakeRequest(client tlsClient.HttpClient, config RequestConfig) ([]byte, *http.Response, error) {
	req, err := http.NewRequest(config.Method, config.URL, config.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Set headers and maintain order
	req.Header = make(http.Header)
	var headerOrder []string
	for _, header := range config.Headers {
		req.Header.Set(header.Key, header.Value)
		headerOrder = append(headerOrder, header.Key)
	}

	req.Header[http.HeaderOrderKey] = headerOrder
	req.Header[http.PHeaderOrderKey] = []string{":authority", ":method", ":path", ":scheme"}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to read response body: %w", err)
	}

	return bodyBytes, resp, nil
}
