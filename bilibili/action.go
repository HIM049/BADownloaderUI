package bilibili

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// 用于下载音频流的函数
// 传入流 URL 和文件名
func StreamingDownloader(audioURL, filePathAndName string) error {
	// 先判断文件是否存在，如果存在则跳过下载，否则创建文件
	out, err := os.Create(filePathAndName)
	if err != nil {
		return err
	}
	defer out.Close()

	// 音频流下载函数。接收音频url和文件名。
	client := &http.Client{}
	request, err := http.NewRequest("GET", audioURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("referer", "https://www.bilibili.com")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}
	return nil
}

// 从 URL 下载图片
func SaveFromURL(url string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 发起 HTTP 请求获取图片内容
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 将图片内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// 工具函数
// json解析函数
func DecodeJson(jsonFile string, object any) error {
	err := json.Unmarshal([]byte([]byte(jsonFile)), object)
	if err != nil {
		return err
	}
	return nil
}
