package bilibili

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/tidwall/gjson"
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

func getFavList(id, ps, pn, sessdata string) (string, error) {

	// 创建请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v3/fav/resource/list", nil)
	if err != nil {
		return "", err
	}

	// 添加 Cookie 到请求头
	if sessdata != "" {
		req.Header.Add("Cookie", "SESSDATA="+sessdata)
	}

	// 设置 URL 参数
	q := req.URL.Query()
	q.Add("media_id", id)    // 每页项数
	q.Add("ps", ps)          // 页码
	q.Add("pn", pn)          // 页码+
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

func GetFavListObj(id, sessdata string, ps, pn int) (*FavList, error) {
	var obj FavList
	body, err := getFavList(id, strconv.Itoa(ps), strconv.Itoa(pn), sessdata)
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

// 获取用户创建的收藏夹
type Collects struct {
	UserMid int `json:"user_mid"`
	Count   int `json:"count"`
	List    []meta
}
type meta struct {
	Id         int    `json:"id"`   // 收藏夹 ID
	Mid        int    `json:"mid"`  // 创建者 MID
	Attr       int    `json:"attr"` // 属性
	Title      string `json:"title"`
	Cover      string `json:"cover"`
	MediaCount int    `json:"media_count"`
}

// 获取用户收藏的收藏夹
func (collects *Collects) GetFavCollect(sessdata string, ps, pn int) error {
	json, err := getUserfavoritesCollect(sessdata, strconv.Itoa(collects.UserMid), strconv.Itoa(ps), strconv.Itoa(pn))
	if err != nil {
		return err
	}

	// 错误检查
	if CheckObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}

	collects.Count = int(gjson.Get(json, "data.count").Int())
	pageCount := collects.Count

	if collects.Count/20 >= pn {
		pageCount = 20
	} else {
		pageCount = collects.Count % 20
	}

	for i := 0; i < pageCount; i++ {
		meta := new(meta)
		meta.Id = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".id").Int())
		meta.Mid = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".mid").Int())
		meta.Attr = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".attr").Int())
		meta.Title = gjson.Get(json, "data.list."+strconv.Itoa(i)+".title").String()
		meta.Cover = gjson.Get(json, "data.list."+strconv.Itoa(i)+".cover").String()
		meta.MediaCount = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".media_count").Int())
		collects.List = append(collects.List, *meta)
	}

	return nil
}

// 获取用户收藏的收藏夹
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

// 获取用户创建的收藏夹
func (collects *Collects) GetUsersCollect(sessdata string) error {
	json, err := getUsersCollect(sessdata, strconv.Itoa(collects.UserMid))
	if err != nil {
		return err
	}

	// 错误检查
	if CheckObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}

	collects.Count = int(gjson.Get(json, "data.count").Int())
	for i := 0; i < collects.Count; i++ {
		meta := new(meta)
		meta.Id = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".id").Int())
		meta.Mid = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".mid").Int())
		meta.Attr = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".attr").Int())
		meta.Title = gjson.Get(json, "data.list."+strconv.Itoa(i)+".title").String()
		meta.MediaCount = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".media_count").Int())
		collects.List = append(collects.List, *meta)
	}

	return nil
}

// 获取用户创建的收藏夹
func getUsersCollect(sessdata, mid string) (string, error) {
	// 创建请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v3/fav/folder/created/list-all", nil)
	if err != nil {
		return "", err
	}

	// 添加 Cookie 到请求头
	if sessdata != "" {
		req.Header.Add("Cookie", "SESSDATA="+sessdata)
	}

	// 设置 URL 参数
	q := req.URL.Query()
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
