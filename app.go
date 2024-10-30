package main

import (
	"bili-audio-downloader/services"
	"context"
	"os"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// startup is called when the app starts. The context is saved
// so we can call the wails methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 程序初始化
	cfg := new(Config)
	if cfg.Get() != nil {
		wails.LogFatal(a.ctx, "Initialize Config Faild")
	} else {
		wails.LogInfo(a.ctx, "Initialize Config Successful")
	}

	downloadPath := cfg.FileConfig.DownloadPath
	cachePath := cfg.FileConfig.CachePath
	err2 := os.MkdirAll(downloadPath, 0755)
	err3 := os.MkdirAll(cachePath, 0755)
	err4 := os.MkdirAll(cachePath+"/music", 0755)
	err5 := os.MkdirAll(cachePath+"/cover", 0755)
	err6 := os.MkdirAll(cachePath+"/single/cover", 0755)
	err7 := os.MkdirAll(cachePath+"/single/music", 0755)
	if err2 != nil ||
		err3 != nil ||
		err4 != nil ||
		err5 != nil ||
		err6 != nil ||
		err7 != nil {
		wails.LogFatal(a.ctx, "Initialize Folder Faild")
	} else {
		wails.LogInfo(a.ctx, "Initialize Folder Successful")
	}

	// 检查版本更新
	version, err := services.CheckUpdate(APP_VERSION)
	if err != nil {
		wails.LogErrorf(a.ctx, "Check for update Faild: %s", err)
	} else if version == "0" {
		wails.LogInfo(a.ctx, "No software update")
	} else {
		wails.LogInfof(a.ctx, "Found new version: %s", version)

		result, err := wails.MessageDialog(a.ctx, wails.MessageDialogOptions{
			Type:          wails.QuestionDialog,
			Title:         "找到新版本：" + version,
			Message:       "软件有新版本发布了，是否前往下载？",
			DefaultButton: "Yes",
		})

		if err != nil {
			wails.LogError(a.ctx, "弹出更新提示失败")
		}

		wails.LogDebugf(a.ctx, "选择结果：%s", result)

		if result == "Yes" {
			wails.BrowserOpenURL(a.ctx, "https://github.com/HIM049/BADownloaderUI/releases/tag/"+version)
		}

	}
}

// 程序关闭时
func (a *App) shutdown(ctx context.Context) {
	// 清理缓存
	cfg := new(Config)
	cfg.Get()
	if cfg.DeleteCache {
		os.RemoveAll(cfg.FileConfig.CachePath)
	}
}

type DownloadOption struct {
	SongName   bool `json:"song_name"`
	SongCover  bool `json:"song_cover"`
	SongAuthor bool `json:"song_author"`
}
