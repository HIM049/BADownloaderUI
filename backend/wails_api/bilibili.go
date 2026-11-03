package wails_api

import (
	"bili-audio-downloader/backend/biliapi"
	"bili-audio-downloader/backend/config"
	"errors"
	"strconv"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// QueryVideo 查询视频信息
func (w *WailsApi) QueryVideo(bvid string) (VideoInfo, error) {
	sessdata := accountSessdata()

	view, err := biliapi.GetVideoView(bvid, sessdata)
	if err != nil {
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return VideoInfo{}, err
	}

	var info VideoInfo
	info.Meta = MediaMeta{Title: view.Title, Cover: view.Pic}
	info.Up.Name = view.Owner.Name
	info.Up.Avatar = view.Owner.Face
	return info, nil
}

// QueryCollection 查询并返回收藏夹信息
func (w *WailsApi) QueryCollection(favListID string) FavListInfo {
	sessdata := accountSessdata()
	favlist, err := biliapi.GetFavList(favListID, sessdata, 1, 1)
	if err != nil {
		wails.LogErrorf(w.ctx, "获取收藏夹内容时出现错误：%s", err)
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return FavListInfo{}
	}

	var info FavListInfo
	info.Data.Info.Title = favlist.Info.Title
	info.Data.Info.Cover = favlist.Info.Cover
	info.Data.Info.MediaCount = favlist.Info.MediaCount
	info.Data.Info.Upper.Name = favlist.Info.Upper.Name
	info.Data.Info.Upper.Face = favlist.Info.Upper.Face
	return info
}

// QueryCompilation 查询并返回合集信息
func (w *WailsApi) QueryCompilation(mid, sid int) CompilationInfo {
	season, err := biliapi.GetSeasonArchives(mid, sid, 1, 1)
	if err != nil {
		wails.LogErrorf(w.ctx, "获取合集内容时出现错误：%s", err)
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return CompilationInfo{}
	}

	var info CompilationInfo
	info.Data.Meta.Name = season.Meta.Name
	info.Data.Meta.Cover = season.Meta.Cover
	info.Data.Meta.Total = season.Meta.Total
	return info
}

// QueryAudio 查询音频信息
func (w *WailsApi) QueryAudio(auid string) (AudioInfo, error) {
	audio, err := biliapi.GetAudioInfo(auid)
	if err != nil {
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return AudioInfo{}, err
	}

	var info AudioInfo
	info.Meta = MediaMeta{Title: audio.Title, Cover: audio.Cover}
	info.Up.Author = audio.Author
	return info, nil
}

// QueryProfileVideo 查询用户投稿视频数量
func (w *WailsApi) QueryProfileVideo(mid string) (int, error) {
	sessdata := accountSessdata()

	midInt, err := strconv.Atoi(mid)
	if err != nil {
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return 0, err
	}

	space, err := biliapi.GetSpaceVideos(midInt, sessdata, 1, 1)
	if err != nil {
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return 0, err
	}
	return space.Page.Count, nil
}

// GetUsersCollect 获取用户创建的收藏夹
func (w *WailsApi) GetUsersCollect() FolderList {
	mid, _ := strconv.Atoi(config.Cfg.Account.DedeUserID)
	folders, err := biliapi.GetCreatedFolders(mid, config.Cfg.Account.SESSDATA)
	if err != nil {
		wails.LogErrorf(w.ctx, "获取收藏夹列表失败：%s", err)
		return FolderList{}
	}

	return FolderList{UserMid: mid, Count: folders.Count, List: folders.List}
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
func (w *WailsApi) GetFavCollect(pn int) FolderList {
	mid, _ := strconv.Atoi(config.Cfg.Account.DedeUserID)
	folders, err := biliapi.GetCollectedFolders(mid, config.Cfg.Account.SESSDATA, 20, pn)
	if err != nil {
		wails.LogErrorf(w.ctx, "获取收藏夹列表失败：%s", err)
		return FolderList{}
	}

	return FolderList{UserMid: mid, Count: folders.Count, List: folders.List}
}
