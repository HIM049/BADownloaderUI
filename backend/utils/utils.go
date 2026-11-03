package utils

import (
	"errors"
	"regexp"
	"strings"
)

// CheckFileName 剔除文件名中的奇怪字符
func CheckFileName(SFileN string) string {
	// 英文双引号转中文双引号
	SFileN = strings.ReplaceAll(SFileN, "\"", "“")
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
