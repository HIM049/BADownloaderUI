package biliapi

import (
	"encoding/json"
	"fmt"
)

// 扫码登录状态
const (
	QRCodeSuccess    = 0
	QRCodeExpired    = 86038
	QRCodeScanned    = 86090 // 已扫码未确认
	QRCodeNotScanned = 86101
)

// QRCode 登录二维码
type QRCode struct {
	URL string `json:"url"`
	Key string `json:"qrcode_key"`
}

// GetQRCode 申请登录二维码
// https://socialsisteryi.github.io/bilibili-API-collect/docs/login/login_action/QR.html
func GetQRCode() (QRCode, error) {
	return getData[QRCode]("https://passport.bilibili.com/x/passport-login/web/qrcode/generate")
}

// QRCodePoll 扫码状态，登录成功时 Cookies 包含账号 cookie
type QRCodePoll struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Cookies map[string]string `json:"-"`
}

// PollQRCode 查询扫码状态
func PollQRCode(key string) (QRCodePoll, error) {
	resp, err := send(apiClient, "https://passport.bilibili.com/x/passport-login/web/qrcode/poll",
		WithParams(map[string]string{"qrcode_key": key}),
	)
	if err != nil {
		return QRCodePoll{}, err
	}
	defer resp.Body.Close()

	var env envelope[QRCodePoll]
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return QRCodePoll{}, fmt.Errorf("decode response: %w", err)
	}
	if err := env.err(); err != nil {
		return QRCodePoll{}, err
	}

	poll := env.Data
	poll.Cookies = make(map[string]string)
	for _, cookie := range resp.Cookies() {
		poll.Cookies[cookie.Name] = cookie.Value
	}
	return poll, nil
}
