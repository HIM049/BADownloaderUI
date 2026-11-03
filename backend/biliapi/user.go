package biliapi

// CodeNotLogin 账号未登录
const CodeNotLogin = -101

// Nav 导航栏信息，未登录时返回 CodeNotLogin，但仍包含 wbi 密钥
type Nav struct {
	IsLogin bool   `json:"isLogin"`
	Mid     int    `json:"mid"`
	Face    string `json:"face"`
	Uname   string `json:"uname"`
	WbiImg  struct {
		ImgURL string `json:"img_url"`
		SubURL string `json:"sub_url"`
	} `json:"wbi_img"`
}

// GetNav 获取导航栏信息
// https://socialsisteryi.github.io/bilibili-API-collect/docs/login/login_info.html
func GetNav(sessdata string) (Nav, error) {
	return getData[Nav]("https://api.bilibili.com/x/web-interface/nav", WithSESSDATA(sessdata))
}
