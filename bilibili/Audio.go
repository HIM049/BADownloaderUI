package bilibili

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

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
