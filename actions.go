package main

import (
	"github.com/myuser/bilibili"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetAppVersion() string {
	return APP_VERSION
}

// 查询并返回收藏夹信息
func (a *App) SearchFavListInformation(favListID string) bilibili.FavList {
	listInf, err := bilibili.GetFavListObj(favListID, 1, 1)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取收藏夹内容时出现错误：%s", err)
		return bilibili.FavList{}
	}
	return *listInf
}

// 查询并返回收藏夹信息
func (a *App) SearchCompListInformation(mid, sid int) bilibili.CompliationInformation {
	listInf, err := bilibili.GetCompliationObj(mid, sid, 1, 1)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取合集内容时出现错误：%s", err)
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
	}
}
