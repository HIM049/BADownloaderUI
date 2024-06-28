package bilibili

import (
	"errors"
	"io"
	"net/http"
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
func (v *Video) Query(sessdata, bvid string) error {
	json, err := GetVideoPageInformation(bvid, sessdata)
	if err != nil {
		return err
	}

	// 错误检查
	if CheckObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}

	// 将信息写入结构体
	v.Bvid = bvid
	v.Meta.Title = gjson.Get(json, "data.title").String()                        // 视频标题
	v.Meta.Cover = gjson.Get(json, "data.pic").String()                          // 视频封面
	v.Meta.LyricsPath = gjson.Get(json, "data.subtitle.0.subtitle_url").String() // 字幕获取（临时）
	v.Up.Mid = int(gjson.Get(json, "data.owner.mid").Int())                      // UP MID
	v.Up.Name = gjson.Get(json, "data.owner.name").String()                      // UP 昵称
	v.Up.Avatar = gjson.Get(json, "data.owner.face").String()                    // UP 头像

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
