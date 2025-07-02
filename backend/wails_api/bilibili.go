package wails_api

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/bilibili"
	"errors"
	"github.com/tidwall/gjson"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
	"strconv"
)

// QueryVideo 查询视频信息
func (w *WailsApi) QueryVideo(bvid string) (bilibili.Video, error) {
	sessdata := ""
	if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
		sessdata = config.Cfg.Account.SESSDATA
	}

	video := new(bilibili.Video)
	err := video.Query(sessdata, bvid)
	if err != nil {
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return bilibili.Video{}, err
	}
	return *video, err
}

// QueryCollection 查询并返回收藏夹信息
func (w *WailsApi) QueryCollection(favListID string) bilibili.FavList {
	sessdata := ""
	if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
		sessdata = config.Cfg.Account.SESSDATA
	}
	listInf, err := bilibili.GetFavListObj(favListID, sessdata, 1, 1)
	if err != nil {
		wails.LogErrorf(w.ctx, "获取收藏夹内容时出现错误：%s", err)
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return bilibili.FavList{}
	}
	return *listInf
}

// QueryCompilation 查询并返回合集信息
func (w *WailsApi) QueryCompilation(mid, sid int) bilibili.CompliationInformation {
	listInf, err := bilibili.GetCompliationObj(mid, sid, 1, 1)
	if err != nil {
		wails.LogErrorf(w.ctx, "获取合集内容时出现错误：%s", err)
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return bilibili.CompliationInformation{}
	}
	return *listInf
}

// QueryAudio 查询音频信息
func (w *WailsApi) QueryAudio(auid string) (bilibili.Audio, error) {
	audio := new(bilibili.Audio)
	err := audio.Query(auid)
	if err != nil {
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return bilibili.Audio{}, err
	}
	return *audio, err
}

// QueryProfileVideo 查询音频信息
func (w *WailsApi) QueryProfileVideo(mid string) (int, error) {
	sessdata := ""
	if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
		sessdata = config.Cfg.Account.SESSDATA
	}

	respJson, err := bilibili.GetProfileVideo(mid, "1", "1", sessdata)
	if err != nil {
		wails.EventsEmit(w.ctx, "error", "错误："+err.Error())
		return 0, err
	}
	return int(gjson.Get(respJson, "data.page.count").Int()), err
}

// GetUsersCollect 获取用户创建的收藏夹
func (w *WailsApi) GetUsersCollect() bilibili.Collects {
	// 获取收藏夹列表
	collects := new(bilibili.Collects)
	mid, _ := strconv.Atoi(config.Cfg.Account.DedeUserID)
	collects.UserMid = mid
	err := collects.GetUsersCollect(config.Cfg.Account.SESSDATA)
	if err != nil {
		wails.LogErrorf(w.ctx, "获取收藏夹列表失败：%s", err)
		return bilibili.Collects{}
	}

	return *collects
}

// GetUserInf 获取已登录用户的信息
func (w *WailsApi) GetUserInf() (bilibili.AccountInformation, error) {
	if !config.Cfg.Account.IsLogin {
		return bilibili.AccountInformation{}, errors.New("用户未登录")
	}
	sessdata := config.Cfg.Account.SESSDATA

	accountInf := new(bilibili.AccountInformation)
	err := accountInf.GetUserInf(sessdata)
	if err != nil {
		return bilibili.AccountInformation{}, err
	}

	return *accountInf, nil
}

// GetFavCollect 获取收藏的收藏夹
func (w *WailsApi) GetFavCollect(pn int) bilibili.Collects {
	// 获取收藏夹列表
	collects := new(bilibili.Collects)
	mid, _ := strconv.Atoi(config.Cfg.Account.DedeUserID)
	collects.UserMid = mid
	err := collects.GetFavCollect(config.Cfg.Account.SESSDATA, 20, pn)
	if err != nil {
		wails.LogErrorf(w.ctx, "获取收藏夹列表失败：%s", err)
		return bilibili.Collects{}
	}

	return *collects
}
