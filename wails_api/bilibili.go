package wails_api

import (
	"bili-audio-downloader/bilibili"
	"bili-audio-downloader/config"
	"github.com/tidwall/gjson"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// QueryVideo 查询视频信息
func (a *WailsApi) QueryVideo(bvid string) (bilibili.Video, error) {
	sessdata := ""
	if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
		sessdata = config.Cfg.Account.SESSDATA
	}

	video := new(bilibili.Video)
	err := video.Query(sessdata, bvid)
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return bilibili.Video{}, err
	}
	return *video, err
}

// QueryCollection 查询并返回收藏夹信息
func (a *WailsApi) QueryCollection(favListID string) bilibili.FavList {
	sessdata := ""
	if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
		sessdata = config.Cfg.Account.SESSDATA
	}
	listInf, err := bilibili.GetFavListObj(favListID, sessdata, 1, 1)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取收藏夹内容时出现错误：%s", err)
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return bilibili.FavList{}
	}
	return *listInf
}

// QueryCompilation 查询并返回合集信息
func (a *WailsApi) QueryCompilation(mid, sid int) bilibili.CompliationInformation {
	listInf, err := bilibili.GetCompliationObj(mid, sid, 1, 1)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取合集内容时出现错误：%s", err)
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return bilibili.CompliationInformation{}
	}
	return *listInf
}

// QueryAudio 查询音频信息
func (a *WailsApi) QueryAudio(auid string) (bilibili.Audio, error) {
	audio := new(bilibili.Audio)
	err := audio.Query(auid)
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return bilibili.Audio{}, err
	}
	return *audio, err
}

// QueryProfileVideo 查询音频信息
func (a *WailsApi) QueryProfileVideo(mid string) (int, error) {
	sessdata := ""
	if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
		sessdata = config.Cfg.Account.SESSDATA
	}

	respJson, err := bilibili.GetProfileVideo(mid, "1", "1", sessdata)
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return 0, err
	}
	return int(gjson.Get(respJson, "data.page.count").Int()), err
}
