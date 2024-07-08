package bilibili

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tidwall/gjson"
)

// // 用于获取 AUID 音频流信息
// type audio struct {
// 	Code int    `json:"code"`
// 	Msg  string `json:"msg"`
// 	Data struct {
// 		ID     int    `json:"id"`     // 音频 AUID
// 		Title  string `json:"title"`  // 音频标题
// 		Cover  string `json:"cover"`  // 音频封面
// 		Intro  string `json:"intro"`  // 音频简介
// 		Lyric  string `json:"lyric"`  // lrc歌词url
// 		Author string `json:"author"` // 作者名
// 		Bvid   string `json:"bvid"`   // 关联稿件 BVID
// 		Cid    int    `json:"cid"`    // 关联稿件 CID
// 	}
// }

// // 用于获取 AUID 音频流
// type AudioStream struct {
// 	Code int    `json:"code"`
// 	Msg  string `json:"msg"`
// 	Data struct {
// 		Type  int      `json:"type"`  //-1：试听片段（192K） 0：128K 1：192K 2：320K 3：FLAC
// 		Title string   `json:"title"` // 音频标题
// 		Cover string   `json:"cover"` // 音频封面
// 		Cdns  []string `json:"cdns"`  // 音频流列表
// 	}
// }

type Audio struct {
	Auid string `json:"auid"`
	Meta struct {
		Title string `json:"title"` // 音频标题
		Cover string `json:"cover"` // 音频封面
		Lyric string `json:"lyric"` // lrc歌词url
	}
	Up struct {
		Author string `json:"author"` // 作者名
	}
	Stream struct {
		Type       int    `json:"type"`        //-1：试听片段（192K） 0：128K 1：192K 2：320K 3：FLAC
		StreamLink string `json:"stream_link"` // 音频流列表
	}
}

func (audio *Audio) Query(auid string) error {

	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://www.bilibili.com/audio/music-service-c/web/song/info")

	// 设置 URL 参数
	params.Set("sid", auid)

	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		return err
	}
	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyJson := string(body)
	defer resp.Body.Close()

	audio.Auid = auid
	audio.Meta.Title = gjson.Get(bodyJson, "data.title").String()
	audio.Meta.Cover = gjson.Get(bodyJson, "data.cover").String()
	audio.Meta.Lyric = gjson.Get(bodyJson, "data.lyric").String()
	audio.Up.Author = gjson.Get(bodyJson, "data.author").String()

	return nil
}

func (audio *Audio) GetStream(sessdata string) error {
	// 创建请求
	req, err := http.NewRequest("GET", "https://api.bilibili.com/audio/music-service-c/url", nil)
	if err != nil {
		return err
	}

	// 添加 Cookie 到请求头
	if sessdata != "" {
		req.Header.Add("Cookie", "SESSDATA="+sessdata)
	}

	// 设置 URL 参数
	q := req.URL.Query()
	q.Add("songid", audio.Auid)
	q.Add("quality", "2")
	q.Add("privilege", "2")
	q.Add("mid", "2")
	q.Add("platform", "web")
	req.URL.RawQuery = q.Encode()

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
	if CheckObj(int(gjson.Get(bodyJson, "code").Int())) {
		return errors.New(gjson.Get(bodyJson, "message").String())
	}

	audio.Stream.Type = int(gjson.Get(bodyJson, "data.type").Int())
	audio.Stream.StreamLink = gjson.Get(bodyJson, "data.cdns.0").String()

	return nil
}
