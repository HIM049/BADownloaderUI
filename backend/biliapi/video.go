package biliapi

import "strconv"

// VideoView 视频详情
type VideoView struct {
	Bvid  string `json:"bvid"`
	Title string `json:"title"`
	Pic   string `json:"pic"`
	Owner struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"owner"`
	Pages []VideoPage `json:"pages"`
}

// VideoPage 视频分 P
type VideoPage struct {
	Cid  int    `json:"cid"`
	Part string `json:"part"`
}

// GetVideoView 获取视频详情
// https://socialsisteryi.github.io/bilibili-API-collect/docs/video/info.html
func GetVideoView(bvid, sessdata string) (VideoView, error) {
	return getData[VideoView]("https://api.bilibili.com/x/web-interface/view",
		WithParams(map[string]string{"bvid": bvid}),
		WithSESSDATA(sessdata),
	)
}

// PlayURL 视频取流结果，仅保留音频部分
type PlayURL struct {
	Dash struct {
		Audio []struct {
			BaseURL string `json:"base_url"`
		} `json:"audio"`
		Flac *struct {
			Audio *struct {
				BaseURL string `json:"base_url"`
			} `json:"audio"`
		} `json:"flac"`
	} `json:"dash"`
}

// AudioURL 返回第一条普通音频流地址，没有时为空
func (p PlayURL) AudioURL() string {
	if len(p.Dash.Audio) == 0 {
		return ""
	}
	return p.Dash.Audio[0].BaseURL
}

// FlacURL 返回无损音频流地址，没有时为空
func (p PlayURL) FlacURL() string {
	if p.Dash.Flac == nil || p.Dash.Flac.Audio == nil {
		return ""
	}
	return p.Dash.Flac.Audio.BaseURL
}

// GetPlayURL 获取视频流地址
// https://socialsisteryi.github.io/bilibili-API-collect/docs/video/videostream_url.html
func GetPlayURL(bvid string, cid int, sessdata string) (PlayURL, error) {
	return getData[PlayURL]("https://api.bilibili.com/x/player/wbi/playurl",
		WithParams(map[string]string{
			"bvid":  bvid,
			"cid":   strconv.Itoa(cid),
			"fnval": "16",
		}),
		WithWbi(),
		WithSESSDATA(sessdata),
	)
}
