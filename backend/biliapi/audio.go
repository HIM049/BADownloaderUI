package biliapi

// AudioInfo 音频稿件信息
type AudioInfo struct {
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	Lyric  string `json:"lyric"` // lrc 歌词 url
	Author string `json:"author"`
}

// GetAudioInfo 获取音频稿件信息
func GetAudioInfo(auid string) (AudioInfo, error) {
	return getData[AudioInfo]("https://www.bilibili.com/audio/music-service-c/web/song/info",
		WithParams(map[string]string{"sid": auid}),
	)
}

// AudioStream 音频取流结果
type AudioStream struct {
	Type int      `json:"type"` // -1：试听片段（192K） 0：128K 1：192K 2：320K 3：FLAC
	Cdns []string `json:"cdns"`
}

// URL 返回第一条音频流地址，没有时为空
func (s AudioStream) URL() string {
	if len(s.Cdns) == 0 {
		return ""
	}
	return s.Cdns[0]
}

// GetAudioStream 获取音频流地址
func GetAudioStream(auid, sessdata string) (AudioStream, error) {
	return getData[AudioStream]("https://api.bilibili.com/audio/music-service-c/url",
		WithParams(map[string]string{
			"songid":    auid,
			"quality":   "2",
			"privilege": "2",
			"mid":       "2",
			"platform":  "web",
		}),
		WithSESSDATA(sessdata),
	)
}
