package services

import (
	"errors"
	"regexp"
)

// 书名号匹配
func ExtractTitle(text string) (string, error) {
	// 定义书名号正则表达式
	reg := regexp.MustCompile(`《(.*?)》`)

	// 查找匹配的字符串
	matches := reg.FindStringSubmatch(text)
	if len(matches) < 2 {
		return "", errors.New("无法找到合适的书名号")
	}

	// 返回匹配的书名号内容
	return matches[1], nil
}

// 剔除方括号内容
func rejectSquareBreckets(text string) string {
	reg := regexp.MustCompile(`【(.*?)】`)
	return reg.ReplaceAllString(text, "")
}

// 剔除圆括号内容
func rejectBreckets(text string) string {
	regCN := regexp.MustCompile(`（(.*?)）`)
	regEN := regexp.MustCompile(`/((.*?)/)`)

	return regEN.ReplaceAllString(regCN.ReplaceAllString(text, ""), "")
}

// TODO: 整合结构体？& 闭包处理
