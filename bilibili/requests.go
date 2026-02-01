package bilibili

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type requestOption func(*requestConfig)

type requestConfig struct {
	method  string
	url     string
	params  map[string]string
	headers map[string]string
	cookies map[string]string
	useWbi  bool
}

func defaultRequestConfig(rawURL string) *requestConfig {
	return &requestConfig{
		method:  "GET",
		url:     rawURL,
		params:  make(map[string]string),
		headers: make(map[string]string),
		cookies: make(map[string]string),
		useWbi:  false,
	}
}

// WithMethod sets the HTTP method
func WithMethod(method string) requestOption {
	return func(c *requestConfig) {
		c.method = method
	}
}

// WithParams adds query parameters
func WithParams(params map[string]string) requestOption {
	return func(c *requestConfig) {
		for k, v := range params {
			c.params[k] = v
		}
	}
}

// WithHeaders adds HTTP headers
func WithHeaders(headers map[string]string) requestOption {
	return func(c *requestConfig) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithCookies adds cookies
func WithCookies(cookies map[string]string) requestOption {
	return func(c *requestConfig) {
		for k, v := range cookies {
			c.cookies[k] = v
		}
	}
}

// WithSESSDATA is a helper for SESSDATA cookie
func WithSESSDATA(sessdata string) requestOption {
	return func(c *requestConfig) {
		if sessdata != "" {
			c.cookies["SESSDATA"] = sessdata
		}
	}
}

// WithWbi enables Wbi signing
func WithWbi() requestOption {
	return func(c *requestConfig) {
		c.useWbi = true
	}
}

// Request sends an HTTP request with the given options
func Request(rawURL string, options ...requestOption) ([]byte, error) {
	config := defaultRequestConfig(rawURL)
	for _, option := range options {
		option(config)
	}

	reqURL, err := url.Parse(config.url)
	if err != nil {
		return nil, err
	}

	q := reqURL.Query()
	for k, v := range config.params {
		q.Add(k, v)
	}
	reqURL.RawQuery = q.Encode()

	finalURL := reqURL.String()
	if config.useWbi {
		signedUrl, err := WbiSignURLParams(finalURL)
		if err != nil {
			return nil, errors.New("Wbi Sign Error: " + err.Error())
		}
		finalURL = signedUrl
	}

	req, err := http.NewRequest(config.method, finalURL, nil)
	if err != nil {
		return nil, err
	}

	// Default Headers
	req.Header.Set("User-Agent", GetRandomUA())
	req.Header.Set("Referer", "https://www.bilibili.com")

	for k, v := range config.headers {
		req.Header.Set(k, v)
	}

	for k, v := range config.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error: " + strconv.Itoa(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
