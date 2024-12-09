package bilibili

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
)

// 登录密钥请求返回内容
type GetLoginKeyReturn struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Url        string `json:"url"`
		Qrcode_key string `json:"qrcode_key"`
	}
}

// 获取登录密钥
func GetLoginKey() (string, string, error) {
	var obj GetLoginKeyReturn
	body, err := getLoginKey()
	if err != nil {
		return "", "", err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return "", "", err
	}
	// 错误检查
	if obj.Code == 0 {
		return "", "", errors.New(obj.Message)
	}
	return obj.Data.Url, obj.Data.Qrcode_key, nil
}

// 请求登录密钥
func getLoginKey() (string, error) {
	resp, err := http.Get("https://passport.bilibili.com/x/passport-login/web/qrcode/generate")
	if err != nil {
		return "", err
	}
	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString, nil
}

// 用于检查扫码状态和获取 cookie 的函数
type checkLoginReturn struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Url           string `json:"url"`           // 游戏分站跨域登录 url
		Refresh_token string `json:"refresh_token"` // 刷新 refresh_token
		Timestamp     int    `json:"timestamp"`     // 登录时间
		Code          int    `json:"code"`          // 0：扫码登录成功 86038：二维码已失效 86090：二维码已扫码未确认 86101：未扫码
		Message       string `json:"message"`       // 扫码状态信息
	}
}

// 检查扫码状态
func checkLoginStatus(qrcode_key string) (string, *[]*http.Cookie, error) {
	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 创建一个 GET 请求
	req, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/poll", nil)
	if err != nil {
		return "", nil, err
	}

	// 添加参数到请求的查询字符串
	q := req.URL.Query()
	q.Add("qrcode_key", qrcode_key)
	req.URL.RawQuery = q.Encode()

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New("Error:" + strconv.Itoa(resp.StatusCode))
	}

	// 读取 Set-Cookie 头部信息
	cookies := resp.Cookies()

	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString, &cookies, nil
}

func CheckLoginStatus(qrcode_key string) (*checkLoginReturn, *[]*http.Cookie, error) {
	var obj checkLoginReturn
	body, cookies, err := checkLoginStatus(qrcode_key)
	if err != nil {
		return nil, nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, nil, err
	}
	// 错误检查
	if obj.Code == 0 {
		return nil, nil, errors.New(obj.Message)
	}

	return &obj, cookies, nil
}

// TODO: 与登录部分整合结构体
type AccountInformation struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

// 获取用户信息
// https://socialsisteryi.github.io/bilibili-API-collect/docs/login/login_info.html
func (accountInf *AccountInformation) GetUserInf(sessdata string) error {

	// 创建请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/web-interface/nav", nil)
	if err != nil {
		return err
	}

	// 添加 Cookie 到请求头
	if sessdata != "" {
		req.Header.Add("Cookie", "SESSDATA="+sessdata)
	}

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return errors.New("Error: " + strconv.Itoa(resp.StatusCode))
	}

	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyJson := string(body)

	// 错误检查
	if int(gjson.Get(bodyJson, "code").Int()) == 0 {
		return errors.New(gjson.Get(bodyJson, "message").String())
	}

	accountInf.Avatar = gjson.Get(bodyJson, "data.face").String()
	accountInf.Name = gjson.Get(bodyJson, "data.uname").String()

	return nil
}
