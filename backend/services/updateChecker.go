package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

// 通过 GitHub 检查程序更新
// string 为 "0" 代表没有更新，有更新时该位为最新版本号
func CheckUpdate(currentVersion string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", "HIM049/BADownloaderUI")
	resp, err := http.Get(url)
	if err != nil {
		return "0", errors.New(fmt.Sprintln("无法获取最新版本:", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求失败，状态码:", resp.StatusCode)
		return "0", errors.New(fmt.Sprintln("请求失败，状态码:", resp.StatusCode))
	}

	bodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0", errors.New(fmt.Sprintln("读取请求体失败:", err))
	}

	latestVersion := gjson.Get(string(bodyString), "tag_name").String()

	// 比较版本号
	if latestVersion > currentVersion {
		return latestVersion, nil
	}
	return "0", nil
}
