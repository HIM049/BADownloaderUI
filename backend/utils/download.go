package utils

import (
	"bili-audio-downloader/backend/constants"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

type FileName struct {
	Title    string
	Subtitle string
	Quality  string
	ID       int
	Format   string
}

// ParseFileName 解析文件名模板
func ParseFileName(text string) (*template.Template, error) {
	return template.New("filename").Parse(text)
}

// ExportFile 用模板生成文件名，并把文件移动到下载目录
func ExportFile(tmpl *template.Template, downloadDir, title, subtitle, outputFormat string, listid int, currentPath string) error {
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
	var output bytes.Buffer
	if err := tmpl.Execute(&output, fileName); err != nil {
		return err
	}

	// 重命名歌曲文件并移动位置
	return moveFile(currentPath, filepath.Join(downloadDir, output.String()))
}

// moveFile 移动文件，无法直接重命名时（如跨盘）改为复制后删除源文件
func moveFile(src, dst string) error {
	renameErr := os.Rename(src, dst)
	if renameErr == nil {
		return nil
	}
	if err := copyFile(src, dst); err != nil {
		return errors.Join(renameErr, err)
	}
	return os.Remove(src)
}

// copyFile 复制文件内容，失败时删除不完整的目标文件
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	if closeErr := out.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		os.Remove(dst)
	}
	return err
}
