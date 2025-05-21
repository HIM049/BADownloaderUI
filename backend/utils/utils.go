package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"syscall"
)

// CheckFileName 剔除文件名中的奇怪字符
func CheckFileName(SFileN string) string {
	re := regexp.MustCompile(`[/$<>?:*|]`)
	newName := re.ReplaceAllString(SFileN, "")
	return newName
}

// ExtractTitle 书名号匹配
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

// SaveJsonFile 保存 JSON
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

// LoadJsonFile 读取 JSON
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

// IsFileExists 检查文件是否存在
func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// 复用命令行参数函数
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	// 跨平台隐藏控制台窗口（仅 Windows 需要）
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("command '%s' failed: %v", name, err)
	}

	return out.String(), nil
}
