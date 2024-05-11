package main

import (
	"context"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	runtime.LogInfo(a.ctx, "正在初始化文件夹")

	cfg := new(Config)
	if cfg.Get() != nil {
		runtime.LogInfo(a.ctx, "初始化文件夹失败")

	}
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
	cfg := new(Config)
	cfg.Get()
	if cfg.DeleteCache {
		os.RemoveAll(cfg.CachePath)
	}
}

type DownloadOption struct {
	SongName   bool `json:"song_name"`
	SongCover  bool `json:"song_cover"`
	SongAuthor bool `json:"song_author"`
}
