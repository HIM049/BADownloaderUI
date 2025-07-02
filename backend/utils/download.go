package utils

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// StreamingDownloader 用于下载音频流的函数
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

type FileName struct {
	Title    string
	Subtitle string
	Quality  string
	ID       int
	Format   string
}

func ExportFile(title, subtitle, outputFormat string, listid int, currentPath string) error {
	quality := "normal"
	if outputFormat == constants.AudioType.Flac {
		quality = "hires"

	}
	// 处理模板和生成文件名
	fileName := FileName{
		Title:    title,
		Subtitle: subtitle,
		Quality:  quality,
		ID:       listid,
		Format:   outputFormat,
	}
	tmpl, err := template.New("filename").Parse(config.Cfg.FileConfig.FileNameTemplate)
	if err != nil {
		return err
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, fileName)
	if err != nil {
		return err
	}

	// 重命名歌曲文件并移动位置
	err = os.Rename(currentPath, filepath.Join(config.Cfg.GetDownloadPath(), output.String()))
	if err != nil {
		return err
	}
	return nil
}
