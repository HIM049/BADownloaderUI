package bilibili

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
)

// 读取图片的函数
func GetImage(ImgPath string) (string, error) {
	// 打开图片
	img, err := os.Open(ImgPath)
	if err != nil {
		return "", err
	}
	defer img.Close()

	// 读取图片
	data, err := io.ReadAll(img)
	if err != nil {
		return "", err
	}

	// 编码为 base64
	base64Data := base64.StdEncoding.EncodeToString(data)

	return base64Data, nil
}

// 工具函数
// json解析函数
func decodeJson(jsonFile string, object any) error {
	err := json.Unmarshal([]byte([]byte(jsonFile)), object)
	if err != nil {
		return err
	}
	return nil
}

// 工具函数
// 检查结构体中的状态码
func CheckObj(code int) bool {
	if code == 0 {
		return false
	} else {
		return true
	}
}
