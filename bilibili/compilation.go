package bilibili

import (
	"errors"
	"io"
	"net/http"
	"strconv"
)

// 用于获取收藏夹基本信息的函数
// 传入收藏夹 ID ，ps 单页大小， pn 页码
type CompliationInformation struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Archives []struct {
			Bvid  string `json:"bvid"`
			Pic   string `json:"pic"`
			Title string `json:"title"`
		}
		Meta struct {
			Cover       string `json:"cover"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Total       int    `json:"total"`
		}
	}
}

func getCompliation(mid, sid, ps, pn string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/polymer/web-space/seasons_archives_list", nil)
	if err != nil {
		return "", err
	}

	// // 添加 Cookie 到请求头
	// if sessdata != "" {
	// 	req.Header.Add("Cookie", "SESSDATA="+sessdata)
	// }
	req.Header.Set("referer", "https://www.bilibili.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0")

	// 设置 URL 参数
	q := req.URL.Query()
	q.Add("mid", mid)
	q.Add("season_id", sid)
	q.Add("page_size", ps)
	q.Add("page_num", pn)
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

func GetCompliationObj(mid, sid, ps, pn int) (*CompliationInformation, error) {
	var obj CompliationInformation
	body, err := getCompliation(strconv.Itoa(mid), strconv.Itoa(sid), strconv.Itoa(ps), strconv.Itoa(pn))
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}

	// 错误检查
	if CheckObj(obj.Code) {
		return nil, errors.New(obj.Message)
	}
	return &obj, nil
}
