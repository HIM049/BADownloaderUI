package bilibili

import (
	"errors"
	"io"
	"net/http"
	"strconv"
)

// 获取用户投稿列表
// https://socialsisteryi.github.io/bilibili-API-collect/docs/user/space.html#%E6%9F%A5%E8%AF%A2%E7%94%A8%E6%88%B7%E6%8A%95%E7%A8%BF%E8%A7%86%E9%A2%91%E6%98%8E%E7%BB%86
func GetProfileVideo(mid, pn, ps, sessdata string) (string, error) {
	// 创建请求
	request, err := http.NewRequest("GET", "https://api.bilibili.com/x/space/wbi/arc/search", nil)
	if err != nil {
		return "", err
	}

	// 设置 URL 参数
	q := request.URL.Query()
	q.Add("mid", mid)
	q.Add("order", "pubdate")
	q.Add("pn", pn)
	q.Add("ps", ps)
	request.URL.RawQuery = q.Encode()

	signedUrl, err := WbiSignURLParams(request.URL.String())
	if err != nil {
		return "", errors.New("Wbi Sign Error: " + err.Error())
	}

	signedRequest, err := http.NewRequest("GET", signedUrl, nil)
	if err != nil {
		return "", errors.New("New Signed Request Error: " + err.Error())
	}

	signedRequest.Header.Set("referer", "https://www.bilibili.com")
	signedRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0")

	// 添加 Cookie 到请求头
	if sessdata != "" {
		signedRequest.Header.Add("Cookie", "SESSDATA="+sessdata)
	}

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(signedRequest)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Error: " + strconv.Itoa(resp.StatusCode))
	}

	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	return bodyString, nil
}
