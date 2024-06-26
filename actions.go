package main

import (
	"strconv"

	"github.com/myuser/bilibili"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetAppVersion() string {
	return APP_VERSION
}

// 获取用户创建的收藏夹
func (a *App) GetUsersCollect() bilibili.Collects {
	// 获取设置
	cfg := new(Config)
	cfg.Get()

	// 获取收藏夹列表
	collects := new(bilibili.Collects)
	mid, _ := strconv.Atoi(cfg.Account.DedeUserID)
	collects.UserMid = mid
	err := collects.GetUsersCollect(cfg.Account.SESSDATA)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取收藏夹列表失败：%s", err)
		return bilibili.Collects{}
	}

	return *collects
}

// 获取收藏的收藏夹
func (a *App) GetFavCollect(pn int) bilibili.Collects {
	// 获取设置
	cfg := new(Config)
	cfg.Get()

	// 获取收藏夹列表
	collects := new(bilibili.Collects)
	mid, _ := strconv.Atoi(cfg.Account.DedeUserID)
	collects.UserMid = mid
	err := collects.GetFavCollect(cfg.Account.SESSDATA, 20, pn)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取收藏夹列表失败：%s", err)
		return bilibili.Collects{}
	}

	return *collects
}

// 查询并返回收藏夹信息
func (a *App) SearchFavListInformation(favListID string) bilibili.FavList {
	cfg := new(Config)
	cfg.Get()
	sessdata := ""
	if cfg.Account.UseAccount && cfg.Account.IsLogin {
		sessdata = cfg.Account.SESSDATA
	}
	listInf, err := bilibili.GetFavListObj(favListID, sessdata, 1, 1)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取收藏夹内容时出现错误：%s", err)
		runtime.EventsEmit(a.ctx, "error", err.Error())
		return bilibili.FavList{}
	}
	return *listInf
}

// 查询并返回合集信息
func (a *App) SearchCompListInformation(mid, sid int) bilibili.CompliationInformation {
	listInf, err := bilibili.GetCompliationObj(mid, sid, 1, 1)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取合集内容时出现错误：%s", err)
		runtime.EventsEmit(a.ctx, "error", err.Error())
		return bilibili.CompliationInformation{}
	}
	return *listInf
}

// 查询并返回歌曲信息
func (a *App) SearchSongInformation(auid string) bilibili.AudioInf {
	runtime.LogInfo(a.ctx, auid)
	audioInf, err := bilibili.GetAudioInfObj(auid)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取歌曲详情时出现错误：%s", err)
		runtime.EventsEmit(a.ctx, "error", "获取歌曲时出错:"+err.Error())
		return bilibili.AudioInf{}
	}
	return *audioInf
}

// 重置设置文件
func (a *App) RefreshConfig() {
	cfg := new(Config)
	cfg.init()
	err := cfg.Save()
	if err != nil {
		runtime.LogErrorf(a.ctx, "写入设置文件失败：%s", err)
		runtime.EventsEmit(a.ctx, "error", "写入设置时出错:"+err.Error())

	}
}

// 读取设置
func (a *App) LoadConfig() Config {
	cfg := new(Config)
	cfg.Get()
	return *cfg
}

// 写入设置
func (a *App) SaveConfig(cfg Config) {
	err := cfg.Save()
	if err != nil {
		runtime.LogErrorf(a.ctx, "写入设置文件失败：%s", err)
		runtime.EventsEmit(a.ctx, "error", "写入设置时出错:"+err.Error())
	}
}
