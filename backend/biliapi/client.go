package biliapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// apiClient 供所有接口请求共用
var apiClient = &http.Client{
	Timeout: 30 * time.Second,
}

// downloadClient 不设总超时，避免大文件下载被中断
var downloadClient = &http.Client{
	Transport: newDownloadTransport(),
}

func newDownloadTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.ResponseHeaderTimeout = 30 * time.Second
	return t
}

// APIError 接口返回非 0 状态码时的错误
type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("bilibili api error %d: %s", e.Code, e.Message)
}

// envelope 接口返回的通用外层结构
type envelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Msg     string `json:"msg"`
	Data    T      `json:"data"`
}

// err 状态码非 0 时返回 APIError，部分接口的错误信息在 msg 字段
func (e envelope[T]) err() error {
	if e.Code == 0 {
		return nil
	}
	msg := e.Message
	if msg == "" {
		msg = e.Msg
	}
	return &APIError{Code: e.Code, Message: msg}
}

type RequestOption func(*requestConfig)

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

// WithMethod 设置请求方法
func WithMethod(method string) RequestOption {
	return func(c *requestConfig) {
		c.method = method
	}
}

// WithParams 添加查询参数
func WithParams(params map[string]string) RequestOption {
	return func(c *requestConfig) {
		for k, v := range params {
			c.params[k] = v
		}
	}
}

// WithHeaders 添加请求头
func WithHeaders(headers map[string]string) RequestOption {
	return func(c *requestConfig) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithCookies 添加 cookie
func WithCookies(cookies map[string]string) RequestOption {
	return func(c *requestConfig) {
		for k, v := range cookies {
			c.cookies[k] = v
		}
	}
}

// WithSESSDATA 非空时添加 SESSDATA cookie
func WithSESSDATA(sessdata string) RequestOption {
	return func(c *requestConfig) {
		if sessdata != "" {
			c.cookies["SESSDATA"] = sessdata
		}
	}
}

// WithWbi 对请求参数进行 wbi 签名
func WithWbi() RequestOption {
	return func(c *requestConfig) {
		c.useWbi = true
	}
}

// request 发送请求并返回响应体
func request(rawURL string, options ...RequestOption) ([]byte, error) {
	resp, err := send(apiClient, rawURL, options...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// getData 发送请求并解开外层结构；状态码非 0 时同时返回 data 和 APIError
func getData[T any](rawURL string, options ...RequestOption) (T, error) {
	var env envelope[T]
	body, err := request(rawURL, options...)
	if err != nil {
		return env.Data, err
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return env.Data, fmt.Errorf("decode response: %w", err)
	}
	return env.Data, env.err()
}

// Download 通过临时文件把 rawURL 的内容保存到 dst
func Download(rawURL, dst string, options ...RequestOption) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	resp, err := send(downloadClient, rawURL, options...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmp := dst + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		os.Remove(tmp)
		return err
	}

	return os.Rename(tmp, dst)
}

// send 按选项构建并发送请求，只返回状态码为 200 的响应
func send(client *http.Client, rawURL string, options ...RequestOption) (*http.Response, error) {
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

	// 接口请求附带设备标识 cookie 以通过风控
	if _, ok := config.cookies["buvid3"]; !ok && reqURL.Host == "api.bilibili.com" {
		if b3, b4 := getBuvid(); b3 != "" {
			config.cookies["buvid3"] = b3
			config.cookies["buvid4"] = b4
		}
	}

	finalURL := reqURL.String()
	if config.useWbi {
		signedUrl, err := signWbi(finalURL)
		if err != nil {
			return nil, err
		}
		finalURL = signedUrl
	}

	req, err := http.NewRequest(config.method, finalURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Referer", "https://www.bilibili.com")

	for k, v := range config.headers {
		req.Header.Set(k, v)
	}

	for k, v := range config.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, reqURL.Host+reqURL.Path)
	}

	return resp, nil
}
