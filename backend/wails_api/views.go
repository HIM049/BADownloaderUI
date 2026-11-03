package wails_api

import "bili-audio-downloader/backend/biliapi"

// 以下结构保持前端现有的字段名，由接口返回的数据转换而来

// MediaMeta 标题和封面
type MediaMeta struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
}

// VideoInfo 视频信息
type VideoInfo struct {
	Meta MediaMeta `json:"Meta"`
	Up   struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"Up"`
}

// AudioInfo 音频信息
type AudioInfo struct {
	Meta MediaMeta `json:"Meta"`
	Up   struct {
		Author string `json:"author"`
	} `json:"Up"`
}

// FavListInfo 收藏夹信息
type FavListInfo struct {
	Data struct {
		Info struct {
			Title      string `json:"title"`
			Cover      string `json:"cover"`
			MediaCount int    `json:"media_count"`
			Upper      struct {
				Name string `json:"name"`
				Face string `json:"face"`
			} `json:"Upper"`
		} `json:"Info"`
	} `json:"Data"`
}

// CompilationInfo 合集信息
type CompilationInfo struct {
	Data struct {
		Meta struct {
			Name  string `json:"name"`
			Cover string `json:"cover"`
			Total int    `json:"total"`
		} `json:"Meta"`
	} `json:"Data"`
}

// FolderList 收藏夹列表
type FolderList struct {
	UserMid int                 `json:"user_mid"`
	Count   int                 `json:"count"`
	List    []biliapi.FavFolder `json:"List"`
}

// AccountInfo 已登录用户信息
type AccountInfo struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}
