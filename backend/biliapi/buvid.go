package biliapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// buvid 缓存风控所需的设备标识 cookie
var buvid struct {
	sync.Mutex
	b3, b4     string
	lastFailed time.Time
}

// getBuvid 返回缓存的 buvid3 和 buvid4，缺失时请求获取；失败后一分钟内不再重试
func getBuvid() (string, string) {
	buvid.Lock()
	defer buvid.Unlock()

	if buvid.b3 != "" || time.Since(buvid.lastFailed) < time.Minute {
		return buvid.b3, buvid.b4
	}

	b3, b4, err := fetchBuvid()
	if err != nil {
		buvid.lastFailed = time.Now()
		return "", ""
	}
	buvid.b3, buvid.b4 = b3, b4
	return b3, b4
}

// fetchBuvid 从 spi 接口获取一组新的 buvid3 和 buvid4
func fetchBuvid() (string, string, error) {
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/frontend/finger/spi", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Referer", "https://www.bilibili.com")

	resp, err := apiClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status %d: finger/spi", resp.StatusCode)
	}

	var env envelope[struct {
		B3 string `json:"b_3"`
		B4 string `json:"b_4"`
	}]
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return "", "", err
	}
	if err := env.err(); err != nil {
		return "", "", err
	}
	if env.Data.B3 == "" {
		return "", "", errors.New("empty buvid3")
	}
	return env.Data.B3, env.Data.B4, nil
}
