package biliapi

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mixinKeyEncTab = []int{
	46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
	33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
	61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
	36, 20, 34, 44, 52,
}

// wbiKeys 缓存 wbi 签名密钥，十分钟刷新一次
var wbiKeys struct {
	sync.Mutex
	img, sub  string
	updatedAt time.Time
}

// signWbi 为 URL 的查询参数添加 wts 和 w_rid 签名
func signWbi(rawURL string) (string, error) {
	urlObj, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	imgKey, subKey, err := getWbiKeys()
	if err != nil {
		return "", err
	}

	query := urlObj.Query()
	params := map[string]string{}
	for k, v := range query {
		params[k] = v[0]
	}
	for k, v := range encWbi(params, imgKey, subKey) {
		query.Set(k, v)
	}
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}

// getWbiKeys 返回缓存的密钥，过期时从导航接口刷新，刷新失败时沿用旧密钥
func getWbiKeys() (string, string, error) {
	wbiKeys.Lock()
	defer wbiKeys.Unlock()

	if wbiKeys.img != "" && time.Since(wbiKeys.updatedAt) < 10*time.Minute {
		return wbiKeys.img, wbiKeys.sub, nil
	}

	imgKey, subKey, err := fetchWbiKeys()
	if err != nil {
		if wbiKeys.img != "" {
			return wbiKeys.img, wbiKeys.sub, nil
		}
		return "", "", err
	}
	wbiKeys.img, wbiKeys.sub, wbiKeys.updatedAt = imgKey, subKey, time.Now()
	return imgKey, subKey, nil
}

// fetchWbiKeys 从导航接口的图片地址中解析密钥，未登录时接口同样返回密钥
func fetchWbiKeys() (string, string, error) {
	nav, err := GetNav("")
	var apiErr *APIError
	if err != nil && !(errors.As(err, &apiErr) && apiErr.Code == CodeNotLogin) {
		return "", "", fmt.Errorf("fetch wbi keys: %w", err)
	}

	imgKey := keyFromURL(nav.WbiImg.ImgURL)
	subKey := keyFromURL(nav.WbiImg.SubURL)
	if len(imgKey+subKey) < len(mixinKeyEncTab) {
		return "", "", errors.New("fetch wbi keys: invalid key length")
	}
	return imgKey, subKey, nil
}

// keyFromURL 取 URL 中不含扩展名的文件名
func keyFromURL(rawURL string) string {
	base := path.Base(rawURL)
	return strings.TrimSuffix(base, path.Ext(base))
}

func encWbi(params map[string]string, imgKey, subKey string) map[string]string {
	mixinKey := getMixinKey(imgKey + subKey)
	params["wts"] = strconv.FormatInt(time.Now().Unix(), 10)

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for k, v := range params {
		params[k] = sanitizeString(v)
	}

	query := url.Values{}
	for _, k := range keys {
		query.Set(k, params[k])
	}

	hash := md5.Sum([]byte(query.Encode() + mixinKey))
	params["w_rid"] = hex.EncodeToString(hash[:])
	return params
}

func getMixinKey(orig string) string {
	var str strings.Builder
	for _, v := range mixinKeyEncTab {
		str.WriteByte(orig[v])
	}
	return str.String()[:32]
}

func sanitizeString(s string) string {
	unwantedChars := []string{"!", "'", "(", ")", "*"}
	for _, char := range unwantedChars {
		s = strings.ReplaceAll(s, char, "")
	}
	return s
}
