package wails_api

import (
	"bili-audio-downloader/backend/biliapi"
	"bili-audio-downloader/backend/config"
	"encoding/base64"
	"errors"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// LoginBilibili 登录 bilibili
func (w *WailsApi) LoginBilibili() error {
	// 获取二维码和请求密钥
	qr, err := biliapi.GetQRCode()
	if err != nil {
		return err
	}

	// 生成二维码
	png, err := qrcode.Encode(qr.URL, qrcode.Medium, 256)
	if err != nil {
		return err
	}
	runtime.EventsEmit(w.ctx, "qrcodeStr", base64.StdEncoding.EncodeToString(png))

	// 请求登录
	cookies, err := func() (map[string]string, error) {
		for {
			time.Sleep(2 * time.Second)

			poll, err := biliapi.PollQRCode(qr.Key)
			if err != nil {
				return nil, err
			}
			switch poll.Code {
			case biliapi.QRCodeSuccess:
				runtime.LogDebug(w.ctx, "登录成功")
				runtime.EventsEmit(w.ctx, "loginStatus", "登录成功")
				return poll.Cookies, nil
			case biliapi.QRCodeExpired:
				runtime.LogDebug(w.ctx, "二维码已失效")
				runtime.EventsEmit(w.ctx, "loginStatus", "二维码已失效")
				return nil, errors.New("二维码已失效")
			case biliapi.QRCodeScanned:
				runtime.LogDebug(w.ctx, "扫描成功，待确认")
				runtime.EventsEmit(w.ctx, "loginStatus", "扫描成功，待确认")
			case biliapi.QRCodeNotScanned:
				runtime.LogDebug(w.ctx, "未扫描")
				runtime.EventsEmit(w.ctx, "loginStatus", "请扫描二维码登录")
			}
		}
	}()
	if err != nil {
		return err
	}

	if cookies["SESSDATA"] == "" {
		return errors.New("登录返回中缺少 SESSDATA")
	}

	config.Cfg.Account.SESSDATA = cookies["SESSDATA"]
	config.Cfg.Account.Bili_jct = cookies["bili_jct"]
	config.Cfg.Account.DedeUserID = cookies["DedeUserID"]
	config.Cfg.Account.DedeUserID__ckMd5 = cookies["DedeUserID__ckMd5"]
	config.Cfg.Account.Sid = cookies["sid"]
	config.Cfg.Account.IsLogin = true
	config.Cfg.Account.UseAccount = true

	err = config.Cfg.UpdateAndSave()
	if err != nil {
		return err
	}

	return nil
}
