package biliapi

import "strconv"

// SeasonArchives 合集信息和一页视频
type SeasonArchives struct {
	Archives []struct {
		Bvid  string `json:"bvid"`
		Pic   string `json:"pic"`
		Title string `json:"title"`
	} `json:"archives"`
	Meta struct {
		Cover       string `json:"cover"`
		Description string `json:"description"`
		Name        string `json:"name"`
		Total       int    `json:"total"`
	} `json:"meta"`
}

// GetSeasonArchives 获取合集信息和第 pn 页视频
func GetSeasonArchives(mid, sid, ps, pn int) (SeasonArchives, error) {
	return getData[SeasonArchives]("https://api.bilibili.com/x/polymer/web-space/seasons_archives_list",
		WithParams(map[string]string{
			"mid":       strconv.Itoa(mid),
			"season_id": strconv.Itoa(sid),
			"page_size": strconv.Itoa(ps),
			"page_num":  strconv.Itoa(pn),
		}),
	)
}
