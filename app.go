package main

import (
	"bili-audio-downloader/services"
	"context"
	"os"
	"path/filepath"

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
		wails.LogFatal(a.ctx, "Initialize Config Failed")
	} else {
		wails.LogInfo(a.ctx, "Initialize Config Successful")
	}

	// 创建文件夹
	downloadPath := cfg.FileConfig.DownloadPath
	cachePath := cfg.FileConfig.CachePath
	var paths []string
	paths = append(paths, downloadPath)
	paths = append(paths, cachePath)
	paths = append(paths, filepath.Join(cachePath, "/music"))
	paths = append(paths, filepath.Join(cachePath, "/cover"))
	paths = append(paths, filepath.Join(cachePath, "/single/cover"))
	paths = append(paths, filepath.Join(cachePath, "/single/music"))

	for _, path := range paths {
		pathString, err := filepath.Abs(path)
		if err != nil {
			wails.LogFatal(a.ctx, "Initialize Folder Failed")
			return
		}

		err = os.MkdirAll(pathString, 0755)
		if err != nil {
			wails.LogFatal(a.ctx, "Initialize Folder Failed")
			return
		}
	}
	wails.LogInfo(a.ctx, "Initialize Folder Successful")

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
