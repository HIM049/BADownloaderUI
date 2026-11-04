package wails_api

import "bili-audio-downloader/backend/biliapi"

// QueryResult 添加前展示的查询结果
type QueryResult struct {
	Title  string `json:"title"`
	Cover  string `json:"cover"`
	Count  int    `json:"count"`
	UpName string `json:"up_name"`
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
