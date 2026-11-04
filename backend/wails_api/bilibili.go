package wails_api

import (
	"bili-audio-downloader/backend/biliapi"
	"bili-audio-downloader/backend/config"
	"errors"
	"strconv"
)

// QueryVideo 查询视频信息
func (w *WailsApi) QueryVideo(bvid string) (QueryResult, error) {
	view, err := biliapi.GetVideoView(bvid, accountSessdata())
	if err != nil {
		return QueryResult{}, err
	}
	return QueryResult{Title: view.Title, Cover: view.Pic, Count: 1, UpName: view.Owner.Name}, nil
}

// QueryCollection 查询收藏夹信息
func (w *WailsApi) QueryCollection(favListID string) (QueryResult, error) {
	favlist, err := biliapi.GetFavList(favListID, accountSessdata(), 1, 1)
	if err != nil {
		return QueryResult{}, err
	}
	info := favlist.Info
	return QueryResult{Title: info.Title, Cover: info.Cover, Count: info.MediaCount, UpName: info.Upper.Name}, nil
}

// QueryCompilation 查询合集信息
func (w *WailsApi) QueryCompilation(mid, sid int) (QueryResult, error) {
	season, err := biliapi.GetSeasonArchives(mid, sid, 1, 1)
	if err != nil {
		return QueryResult{}, err
	}
	meta := season.Meta
	return QueryResult{Title: meta.Name, Cover: meta.Cover, Count: meta.Total}, nil
}

// QueryAudio 查询音频信息
func (w *WailsApi) QueryAudio(auid string) (QueryResult, error) {
	audio, err := biliapi.GetAudioInfo(auid)
	if err != nil {
		return QueryResult{}, err
	}
	return QueryResult{Title: audio.Title, Cover: audio.Cover, Count: 1, UpName: audio.Author}, nil
}

// QueryProfileVideo 查询用户投稿视频数量
func (w *WailsApi) QueryProfileVideo(mid string) (QueryResult, error) {
	midInt, err := strconv.Atoi(mid)
	if err != nil {
		return QueryResult{}, err
	}
	space, err := biliapi.GetSpaceVideos(midInt, accountSessdata(), 1, 1)
	if err != nil {
		return QueryResult{}, err
	}
	return QueryResult{Title: "用户个人空间", Count: space.Page.Count}, nil
}

// GetUsersCollect 获取用户创建的收藏夹
func (w *WailsApi) GetUsersCollect() (FolderList, error) {
	mid, _ := strconv.Atoi(config.Cfg.Account.DedeUserID)
	folders, err := biliapi.GetCreatedFolders(mid, config.Cfg.Account.SESSDATA)
	if err != nil {
		return FolderList{}, err
	}
	return FolderList{UserMid: mid, Count: folders.Count, List: folders.List}, nil
}

// GetUserInf 获取已登录用户的信息
func (w *WailsApi) GetUserInf() (AccountInfo, error) {
	if !config.Cfg.Account.IsLogin {
		return AccountInfo{}, errors.New("用户未登录")
	}
	nav, err := biliapi.GetNav(config.Cfg.Account.SESSDATA)
	if err != nil {
		return AccountInfo{}, err
	}
	return AccountInfo{Avatar: nav.Face, Name: nav.Uname}, nil
}

// GetFavCollect 获取收藏的收藏夹
func (w *WailsApi) GetFavCollect(pn int) (FolderList, error) {
	mid, _ := strconv.Atoi(config.Cfg.Account.DedeUserID)
	folders, err := biliapi.GetCollectedFolders(mid, config.Cfg.Account.SESSDATA, 20, pn)
	if err != nil {
		return FolderList{}, err
	}
	return FolderList{UserMid: mid, Count: folders.Count, List: folders.List}, nil
}
