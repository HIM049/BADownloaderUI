package main

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
)

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
