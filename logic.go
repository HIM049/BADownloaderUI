package main

import (
	"bili-audio-downloader/bilibili"
	"errors"
	"net/http"
	"os/exec"
	"regexp"
	"syscall"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 登录 bilibili
func (a *App) LoginBilibili() error {
	cfg := new(Config)
	cfg.Get()

	// 获取二维码和请求密钥
	url, key, err := bilibili.GetLoginKey()
	if err != nil {
		return err
	}

	// 生成二维码
	qrcodePath := cfg.FileConfig.CachePath + "/qr.png"
	err = qrcode.WriteFile(url, qrcode.Medium, 256, qrcodePath)
	if err != nil {
		return err
	}

	base64Data, err := bilibili.GetImage(qrcodePath)
	if err != nil {
		return err
	}
	runtime.EventsEmit(a.ctx, "qrcodeStr", base64Data)

	// 请求登录
	cookies, err := func() (*[]*http.Cookie, error) {
		for {
			time.Sleep(2 * time.Second)

			returnObj, cookies, err := bilibili.CheckLoginStatus(key)
			if err != nil {
				return nil, err
			}
			switch returnObj.Data.Code {
			case 0:
				// 登录成功
				runtime.LogDebug(a.ctx, "登录成功")
				runtime.EventsEmit(a.ctx, "loginStatus", "登录成功")
				return cookies, nil
			case 86038:
				// 二维码失效
				runtime.LogDebug(a.ctx, "二维码已失效")
				runtime.EventsEmit(a.ctx, "loginStatus", "二维码已失效")
				return nil, errors.New("二维码已失效")
			case 86090:
				// 扫描成功，待确认
				runtime.LogDebug(a.ctx, "扫描成功，待确认")
				runtime.EventsEmit(a.ctx, "loginStatus", "扫描成功，待确认")
			case 86101:
				// 未扫描
				runtime.LogDebug(a.ctx, "未扫描")
				runtime.EventsEmit(a.ctx, "loginStatus", "请扫描二维码登录")
			}
		}
	}()
	if err != nil {
		return err
	}

	cfg.Account.SESSDATA = (*cookies)[0].Value
	cfg.Account.Bili_jct = (*cookies)[1].Value
	cfg.Account.DedeUserID = (*cookies)[2].Value
	cfg.Account.DedeUserID__ckMd5 = (*cookies)[3].Value
	cfg.Account.Sid = (*cookies)[4].Value
	cfg.Account.IsLogin = true
	cfg.Account.UseAccount = true

	err = cfg.Save()
	if err != nil {
		return err
	}

	return nil
}

// 剔除文件名中的奇怪字符
func CheckFileName(SFileN string) string {
	re := regexp.MustCompile(`[/\$<>?:*|]`)
	newName := re.ReplaceAllString(SFileN, "")
	return newName
}

// 打开文件夹
func OpenFolder(path string) error {
	cmd := exec.Command("cmd", "/c", "start", "", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}
