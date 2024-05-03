package bilibili

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// 用于获取收藏夹基本信息的函数
// 传入收藏夹 ID ，ps 单页大小， pn 页码
// 获得如下结构体
type FavList struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 错误消息
	Data    struct {
		Info struct { // 收藏夹信息
			Title       string `json:"title"`       // 收藏夹标题
			Cover       string `json:"cover"`       // 收藏夹封面
			Media_count int    `json:"media_count"` // 收藏夹数量
			Upper       struct {
				Name string `json:"name"` // 创建者昵称
				Face string `json:"face"` // 创建者头像 url
			}
		}
		Medias []struct { // 收藏夹中的视频
			Id    int    `json:"id"`    // 稿件 avid
			Type  int    `json:"type"`  // 内容类型 （视频稿件2 音频12 合集21）
			Title string `json:"title"` // 标题
			Cover string `json:"cover"` // 封面 url
			Page  int    `json:"page"`  // 视频分P数
			Bvid  string `json:"bvid"`  // BV 号
		}
	}
}

func getFavList(id, ps, pn string) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/v3/fav/resource/list")
	// 设置 URL 参数
	params.Set("media_id", id)
	params.Set("ps", ps)
	params.Set("platform", "web")
	params.Set("pn", pn)

	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}

	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString, nil
}

func GetFavListObj(id string, ps, pn int) (*FavList, error) {
	var obj FavList
	body, err := getFavList(id, strconv.Itoa(ps), strconv.Itoa(pn))
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}
	// 错误检查
	if checkObj(obj.Code) {
		return nil, errors.New(obj.Message)
	}
	return &obj, nil
}

// 获取 用户收藏的视频收藏夹 函数
type UserfavoritesCollect struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Count int `json:"count"`
		List  []struct {
			Id          int    `json:"id"`          // 收藏夹ml
			Fid         int    `json:"fid"`         // 原始收藏夹mlid
			Mid         int    `json:"mid"`         // 创建用户mid
			Title       string `json:"title"`       // 收藏夹标题
			Cover       string `json:"cover"`       // 收藏夹封面图片url
			Media_count int    `json:"media_count"` // 收藏夹视频数量
		}
	}
}

func getUserfavoritesCollect(sessdata, mid, pageSize, pageNumber string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v3/fav/folder/collected/list", nil)
	if err != nil {
		return "", err
	}

	// 添加 Cookie 到请求头
	if sessdata != "" {
		req.Header.Add("Cookie", "SESSDATA="+sessdata)
	}

	// 设置 URL 参数
	q := req.URL.Query()
	q.Add("ps", pageSize)    // 每页项数
	q.Add("pn", pageNumber)  // 页码
	q.Add("up_mid", mid)     // 用户 mid
	q.Add("platform", "web") // 平台
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

func GetUserFavoritesCollect(sessdata, mid string, pageSize, pageNumber int) (*UserfavoritesCollect, error) {
	var obj UserfavoritesCollect
	body, err := getUserfavoritesCollect(sessdata, mid, strconv.Itoa(pageSize), strconv.Itoa(pageNumber))
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}
	// 错误检查
	if checkObj(obj.Code) {
		return nil, errors.New(obj.Message)
	}

	return &obj, nil
}
