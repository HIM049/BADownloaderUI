package services

import (
	"encoding/json"
	"os"
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
