package bilibili

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
