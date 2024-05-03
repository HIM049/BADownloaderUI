package bilibili

import (
	"errors"
	"io"
	"net/http"
	"strconv"
)

// 请求视频详细信息
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/info.md
// TODO：重新添加字幕信息
func GetVideoPageInformation(bvid, sessdata string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/web-interface/view", nil)
	if err != nil {
		return "", err
	}

	// 添加 Cookie 到请求头
	if sessdata != "" {
		req.Header.Add("Cookie", "SESSDATA="+sessdata)
	}

	// 设置 URL 参数
	q := req.URL.Query()
	q.Add("bvid", bvid)
	req.URL.RawQuery = q.Encode()

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
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

// 获取视频流
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/videostream_url.md#%E8%8E%B7%E5%8F%96%E8%A7%86%E9%A2%91%E6%B5%81%E5%9C%B0%E5%9D%80_web%E7%AB%AF
// TODO：更换新链，增加 wbi 签名验证
func GetVideoStream(bvid, cid, sessdata string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/player/playurl", nil)
	if err != nil {
		return "", err
	}

	// 添加 Cookie 到请求头
	if sessdata != "" {
		req.Header.Add("Cookie", "SESSDATA="+sessdata)
	}

	// 设置 URL 参数
	q := req.URL.Query()
	q.Add("bvid", bvid)
	q.Add("cid", cid)
	q.Add("fnval", "16")
	req.URL.RawQuery = q.Encode()

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
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
