package bilibili

import (
	"encoding/json"
)

// 工具函数
// json解析函数
func DecodeJson(jsonFile string, object any) error {
	err := json.Unmarshal([]byte([]byte(jsonFile)), object)
	if err != nil {
		return err
	}
	return nil
}
