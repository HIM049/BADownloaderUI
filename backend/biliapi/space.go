package biliapi

import "strconv"

// SpaceVideos 用户投稿视频的一页
type SpaceVideos struct {
	List struct {
		Vlist []struct {
			Bvid  string `json:"bvid"`
			Title string `json:"title"`
		} `json:"vlist"`
	} `json:"list"`
	Page struct {
		Count int `json:"count"`
	} `json:"page"`
}

// GetSpaceVideos 获取用户投稿视频的第 pn 页，按发布时间排序
// https://socialsisteryi.github.io/bilibili-API-collect/docs/user/space.html
func GetSpaceVideos(mid int, sessdata string, ps, pn int) (SpaceVideos, error) {
	return getData[SpaceVideos]("https://api.bilibili.com/x/space/wbi/arc/search",
		WithParams(map[string]string{
			"mid":   strconv.Itoa(mid),
			"order": "pubdate",
			"ps":    strconv.Itoa(ps),
			"pn":    strconv.Itoa(pn),
		}),
		WithWbi(),
		WithSESSDATA(sessdata),
	)
}
