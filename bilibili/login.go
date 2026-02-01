package bilibili

// 登录密钥请求返回内容
type GetLoginKeyReturn struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Url        string `json:"url"`
		Qrcode_key string `json:"qrcode_key"`
	}
}

// Used to be GetLoginKey and getLoginKey here. Moved to api.go

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

// Used to be checkLoginStatus here. Moved to api.go

// TODO: 与登录部分整合结构体
type AccountInformation struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

// Used to be GetUserInf here. Moved to api.go
