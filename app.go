package main

import (
	"context"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	FavListID = ""
)

// App struct
type App struct {
	ctx context.Context
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 程序初始化
	runtime.LogInfo(a.ctx, "正在创建文件夹")

	cfg := GetConfig(a.ctx)
	_ = os.MkdirAll(cfg.DownloadPath, 0755)
	_ = os.MkdirAll(cfg.CachePath, 0755)
	_ = os.MkdirAll(cfg.CachePath+"/music", 0755)
	_ = os.MkdirAll(cfg.CachePath+"/cover", 0755)
	_ = os.MkdirAll(cfg.CachePath+"/single/cover", 0755)
	_ = os.MkdirAll(cfg.CachePath+"/single/music", 0755)
}

// 程序关闭时
func (a *App) shutdown(ctx context.Context) {
	// 清理缓存
	// cfg := GetConfig(a.ctx)
	// os.RemoveAll(cfg.CachePath)
}

// 查询并返回收藏夹信息
func (a *App) SearchFavListInformation(favListID string) FavList {
	FavListID = favListID
	listInf, err := GetFavListObj(1, 1)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取收藏夹内容时出现错误：%s", err)
		return FavList{}
	}
	return *listInf
}

// 查询并返回歌曲信息
func (a *App) SearchSongInformation(auid string) AudioInf {
	runtime.LogInfo(a.ctx, auid)
	audioInf, err := GetAudioInfObj(auid)
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取歌曲详情时出现错误：%s", err)
		return AudioInf{}
	}
	return *audioInf
}

type DownloadOption struct {
	SongName   bool `json:"song_name"`
	SongCover  bool `json:"song_cover"`
	SongAuthor bool `json:"song_author"`
}
