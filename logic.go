package main

import (
	"bili-audio-downloader/bilibili"
	"encoding/json"
	"errors"
	"net/http"
	"os"
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

// 保存 JSON
func SaveJsonFile(filePath string, theData any) error {
	data, err := json.MarshalIndent(theData, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 读取 JSON
func LoadJsonFile(filePath string, obj interface{}) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, obj)
	if err != nil {
		return err
	}
	return nil
}

// 检查文件是否存在
func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // 文件存在
	}
	if os.IsNotExist(err) {
		return false // 文件不存在
	}
	return false // 其他错误
}

// 剔除文件名中的奇怪字符
func CheckFileName(SFileN string) string {
	re := regexp.MustCompile(`[/\$<>?:*|]`)
	newName := re.ReplaceAllString(SFileN, "")
	return newName
}

// 书名号匹配
func ExtractTitle(input string) (string, error) {
	// 定义书名号正则表达式
	re := regexp.MustCompile(`《(.*?)》`)

	// 查找匹配的字符串
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", errors.New("无法找到合适的书名号")
	}

	// 返回匹配的书名号内容
	return matches[1], nil
}

// 工具函数
// 检查结构体中的状态码
func CheckObj(code int) bool {
	if code == 0 {
		return false
	} else {
		return true
	}
}

// 打开文件夹
func OpenFolder(path string) error {
	cmd := exec.Command("cmd", "/c", "start", "", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}
