package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"regexp"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) LoginBilibili() error {
	cfg := GetConfig(a.ctx)

	// 获取二维码和请求密钥
	url, key, err := GetLoginKey()
	if err != nil {
		return err
	}

	// 生成二维码
	err = qrcode.WriteFile(url, qrcode.Medium, 256, cfg.CachePath+"/qr.png")
	if err != nil {
		return err
	}

	// 请求登录
	cookies, err := func() (*[]*http.Cookie, error) {
		for {
			time.Sleep(3 * time.Second)
			returnObj, cookies, err := CheckLoginStatus(key)
			if err != nil {
				return nil, err
			}
			switch returnObj.Data.Code {
			case 0:
				// 登录成功
				runtime.LogInfo(a.ctx, "登录成功")
				return cookies, nil
			case 86038:
				// 二维码失效
				runtime.LogInfo(a.ctx, "二维码已失效")
				return nil, errors.New("二维码已失效")
			case 86090:
				// 扫描成功，待确认
				runtime.LogInfo(a.ctx, "扫描成功，待确认")
			case 86101:
				// 未扫描
				runtime.LogInfo(a.ctx, "未扫描")
			}
		}
	}()
	if err != nil {
		return err
	}

	// for _, cookie := range *cookies {
	// 	fmt.Println("Cookie:", cookie.Name, ": ", cookie.Value)
	// }
	cfg.Account.SESSDATA = (*cookies)[0].Value
	cfg.Account.Bili_jct = (*cookies)[1].Value
	cfg.Account.DedeUserID = (*cookies)[2].Value
	cfg.Account.DedeUserID__ckMd5 = (*cookies)[3].Value
	cfg.Account.Sid = (*cookies)[4].Value

	err = SaveConfig(cfg)
	if err != nil {
		return err
	}

	return nil
}

// 查询用户收藏的收藏夹
func (a *App) QueryFavCollect() (*userfavoritesCollect, error) {
	cfg := GetConfig(a.ctx)

	obj, err := GetUserFavoritesCollect(cfg.Account, 20, 1)
	if err != nil {
		return nil, err
	}
	return obj, nil
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
