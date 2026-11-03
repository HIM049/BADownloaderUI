package biliapi

import "strconv"

// 收藏夹内容类型
const (
	FavMediaVideo = 2
	FavMediaAudio = 12
)

// FavList 收藏夹信息和一页内容
type FavList struct {
	Info struct {
		Title      string `json:"title"`
		Cover      string `json:"cover"`
		MediaCount int    `json:"media_count"`
		Upper      struct {
			Name string `json:"name"`
			Face string `json:"face"`
		} `json:"upper"`
	} `json:"info"`
	Medias []FavMedia `json:"medias"`
}

// FavMedia 收藏夹中的一项内容
type FavMedia struct {
	ID    int    `json:"id"` // 视频为 avid，音频为 auid
	Type  int    `json:"type"`
	Title string `json:"title"`
	Bvid  string `json:"bvid"`
}

// GetFavList 获取收藏夹信息和第 pn 页内容
// https://socialsisteryi.github.io/bilibili-API-collect/docs/fav/list.html
func GetFavList(mediaID, sessdata string, ps, pn int) (FavList, error) {
	return getData[FavList]("https://api.bilibili.com/x/v3/fav/resource/list",
		WithParams(map[string]string{
			"media_id": mediaID,
			"ps":       strconv.Itoa(ps),
			"pn":       strconv.Itoa(pn),
			"platform": "web",
		}),
		WithSESSDATA(sessdata),
	)
}

// FavFolders 收藏夹列表
type FavFolders struct {
	Count int         `json:"count"`
	List  []FavFolder `json:"list"`
}

// FavFolder 收藏夹列表中的一项
type FavFolder struct {
	ID         int    `json:"id"`
	Mid        int    `json:"mid"`
	Attr       int    `json:"attr"`
	Title      string `json:"title"`
	Cover      string `json:"cover"`
	MediaCount int    `json:"media_count"`
}

// GetCreatedFolders 获取用户创建的全部收藏夹
func GetCreatedFolders(mid int, sessdata string) (FavFolders, error) {
	return getData[FavFolders]("https://api.bilibili.com/x/v3/fav/folder/created/list-all",
		WithParams(map[string]string{
			"up_mid":   strconv.Itoa(mid),
			"platform": "web",
		}),
		WithSESSDATA(sessdata),
	)
}

// GetCollectedFolders 获取用户收藏的收藏夹和合集
func GetCollectedFolders(mid int, sessdata string, ps, pn int) (FavFolders, error) {
	return getData[FavFolders]("https://api.bilibili.com/x/v3/fav/folder/collected/list",
		WithParams(map[string]string{
			"up_mid":   strconv.Itoa(mid),
			"ps":       strconv.Itoa(ps),
			"pn":       strconv.Itoa(pn),
			"platform": "web",
		}),
		WithSESSDATA(sessdata),
	)
}
