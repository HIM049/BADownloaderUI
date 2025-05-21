package main

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/services"
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
	config.InitConfig()

	downloadPath := config.Cfg.GetDownloadPath()
	cachePath := config.Cfg.GetCachePath()
	err2 := os.MkdirAll(downloadPath, 0755)
	err3 := os.MkdirAll(cachePath, 0755)
	err5 := os.MkdirAll(cachePath+"/audio", 0755)
	err6 := os.MkdirAll(cachePath+"/cover", 0755)
	if err2 != nil ||
		err3 != nil ||
		err5 != nil ||
		err6 != nil {
		wails.LogFatal(a.ctx, "Initialize Folder Failed")
	} else {
		wails.LogInfo(a.ctx, "Initialize Folder Successful")
	}

	// 检查版本更新
	version, err := services.CheckUpdate(constants.APP_VERSION)
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
	if config.Cfg.DeleteCache {
		os.RemoveAll(config.Cfg.GetCachePath())
	}
}
