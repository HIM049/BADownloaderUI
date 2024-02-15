package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
)

type Video struct {
	Bvid string `json:"bvid"`
	Meta struct {
		Title      string `json:"title"`       // 视频标题
		Cover      string `json:"cover"`       // 封面
		Author     string `json:"author"`      // 作者
		LyricsPath string `json:"lyrics_path"` // 歌词
	}
	Up struct {
		Mid    int    `json:"mid"`    // UP MID
		Name   string `json:"name"`   // UP 昵称
		Avatar string `json:"avatar"` // UP 头像
	}
	Videos []Videos
}
type Videos struct {
	Cid  int    `json:"cid"`
	Part string `json:"part"` // 分集名称
	Meta struct {
		SongName string `json:"song_name"` // 歌名
	}
	Stream struct {
		Audio struct {
			Id      int    `json:"id"`       // 音质代码
			BaseUrl string `json:"base_url"` // 音频流
		}
		Flac struct {
			Id      int    `json:"id"`       // 音质代码
			BaseUrl string `json:"base_url"` // 音频流
		}
	}
}

// 以 BVID 为单位请求视频详细信息
func (v *Video) BvQuery() error {
	json, err := getVideoPageInformation(v.Bvid)
	if err != nil {
		return err
	}

	// 错误检查
	if checkObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}

	// 将信息写入结构体
	v.Meta.Title = gjson.Get(json, "data.title").String()     // 视频标题
	v.Meta.Cover = gjson.Get(json, "data.pic").String()       // 视频封面
	v.Up.Mid = int(gjson.Get(json, "data.owner.mid").Int())   // UP MID
	v.Up.Name = gjson.Get(json, "data.owner.name").String()   // UP 昵称
	v.Up.Avatar = gjson.Get(json, "data.owner.face").String() // UP 头像

	// 根据分 P 数量写入对应信息
	for i := 0; i < int(gjson.Get(json, "data.videos").Int()); i++ {

		// 单个分集视频信息
		videos := Videos{
			Cid:  int(gjson.Get(json, "data.pages."+strconv.Itoa(i)+".cid").Int()),
			Part: gjson.Get(json, "data.pages."+strconv.Itoa(i)+".part").String(),
		}
		v.Videos = append(v.Videos, videos)
	}

	return nil
}

// 请求视频详细信息
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/info.md
// TODO：重新添加字幕信息
func getVideoPageInformation(bvid string) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/web-interface/view")

	// 设置 URL 参数
	params.Set("bvid", bvid)

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

// 获取视频流
// TODO：请求前检查数据
func (v *VideoInformationList) GetStream() error {
	// 请求信息
	json, err := getVideoStream(v.Bvid, strconv.Itoa(v.Cid))
	if err != nil {
		return err
	}
	// 错误检查
	if checkObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}
	v.Audio.Audio.Quality = int(gjson.Get(json, "data.dash.audio.0.id").Int())
	v.Audio.Audio.Stream = gjson.Get(json, "data.dash.audio.0.base_url").String()
	v.Audio.Flac.Quality = int(gjson.Get(json, "data.dash.flac.id").Int())
	v.Audio.Flac.Stream = gjson.Get(json, "data.dash.flac.base_url").String()

	return nil
}

// 获取视频流
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/videostream_url.md#%E8%8E%B7%E5%8F%96%E8%A7%86%E9%A2%91%E6%B5%81%E5%9C%B0%E5%9D%80_web%E7%AB%AF
// TODO：更换新链，增加 wbi 签名验证
func getVideoStream(bvid, cid string) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/player/playurl")

	// 设置 URL 参数
	params.Set("bvid", bvid)
	params.Set("cid", cid)
	params.Set("fnval", "16")

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
	if checkObj(obj.Code) {
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

// 获取 用户收藏的视频收藏夹 函数
type userfavoritesCollect struct {
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

func getUserfavoritesCollect(sessdata, pageSize, pageNumber, mid string) (string, error) {
	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 创建一个 GET 请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v3/fav/folder/collected/list", nil)
	if err != nil {
		return "", err
	}

	// 添加 Cookie 到请求头
	req.Header.Add("Cookie", "SESSDATA="+sessdata)

	// 添加参数到请求的查询字符串
	q := req.URL.Query()
	q.Add("ps", pageSize)    // 每页项数
	q.Add("pn", pageNumber)  // 页码
	q.Add("up_mid", mid)     // 用户 mid
	q.Add("platform", "web") // 平台
	req.URL.RawQuery = q.Encode()

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Error:" + strconv.Itoa(resp.StatusCode))
	}

	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString, nil
}

func GetUserFavoritesCollect(account Account, pageSize, pageNumber int) (*userfavoritesCollect, error) {
	var obj userfavoritesCollect
	body, err := getUserfavoritesCollect(account.SESSDATA, strconv.Itoa(pageSize), strconv.Itoa(pageNumber), account.DedeUserID)
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

// 定义一个带有参数和cookie的get请求函数，返回响应和错误
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
	if checkObj(obj.Code) {
		return nil, nil, errors.New(obj.Message)
	}

	return &obj, cookies, nil
}

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

func getFavList(FavListID, ps, pn string) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/v3/fav/resource/list")
	// 设置 URL 参数
	params.Set("media_id", FavListID)
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

func GetFavListObj(ps, pn int) (*FavList, error) {
	var obj FavList
	body, err := getFavList(FavListID, strconv.Itoa(ps), strconv.Itoa(pn))
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

// 用于获取 AUID 音频流信息
type AudioInf struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ID     int    `json:"id"`     // 音频 AUID
		Title  string `json:"title"`  // 音频标题
		Cover  string `json:"cover"`  // 音频封面
		Intro  string `json:"intro"`  // 音频简介
		Lyric  string `json:"lyric"`  // lrc歌词url
		Author string `json:"author"` // 作者名
		Bvid   string `json:"bvid"`   // 关联稿件 BVID
		Cid    int    `json:"cid"`    // 关联稿件 CID
	}
}

func getAudioInf(auid string) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://www.bilibili.com/audio/music-service-c/web/song/info")

	// 设置 URL 参数
	params.Set("sid", auid)

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
func GetAudioInfObj(auid string) (*AudioInf, error) {
	var obj AudioInf
	body, err := getAudioInf(auid)
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}
	if checkObj(obj.Code) {
		return nil, errors.New(obj.Msg)
	}
	return &obj, nil
}

// 用于获取 AUID 音频流
type AudioStream struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Type  int      `json:"type"`  //-1：试听片段（192K） 0：128K 1：192K 2：320K 3：FLAC
		Title string   `json:"title"` // 音频标题
		Cover string   `json:"cover"` // 音频封面
		Cdns  []string `json:"cdns"`  // 音频流列表
	}
}

func getAudio(auid, quality string) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/audio/music-service-c/url")

	// 设置 URL 参数
	params.Set("songid", auid)
	params.Set("quality", quality)
	params.Set("privilege", "2")
	params.Set("mid", "2")
	params.Set("platform", "web")

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
func GetAudioObj(auid, quality string) (*AudioStream, error) {
	var obj AudioStream
	body, err := getAudio(auid, quality)
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}
	if checkObj(obj.Code) {
		return nil, errors.New(obj.Msg)
	}
	return &obj, nil
}

// 用于下载音频流的函数
// 传入流 URL 和文件名
func StreamingDownloader(audioURL, filePathAndName string) error {
	// 先判断文件是否存在，如果存在则跳过下载，否则创建文件
	out, err := os.Create(filePathAndName)
	if err != nil {
		return err
	}
	defer out.Close()

	// 音频流下载函数。接收音频url和文件名。
	client := &http.Client{}
	request, err := http.NewRequest("GET", audioURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("referer", "https://www.bilibili.com")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}
	return nil
}

// 从 URL 下载图片
func SaveFromURL(url string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 发起 HTTP 请求获取图片内容
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 将图片内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
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
func checkObj(code int) bool {
	if code == 0 {
		return false
	} else {
		return true
	}
}
