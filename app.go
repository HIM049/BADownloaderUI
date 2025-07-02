package main

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/services"
	"context"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

// App struct
type App struct {
	ctx context.Context
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize config
	wails.LogDebug(a.ctx, "Start Initializing Config")
	config.InitConfig(a.ctx)

	// Create folder
	wails.LogDebug(a.ctx, "Start Initializing Folder")
	initFolder(a.ctx)

	// Check software update
	wails.LogDebug(a.ctx, "Start Check Update")
	go checkUpdateAndAlarm(a.ctx)

}

// shutdown is called when the app shutdown
func (a *App) shutdown(ctx context.Context) {
	// 清理缓存
	if config.Cfg.DeleteCache {
		os.RemoveAll(config.Cfg.GetCachePath())
	}
}

// Check update and alarm
func checkUpdateAndAlarm(ctx context.Context) {
	// Check update
	version, err := services.CheckUpdate(constants.APP_VERSION)
	if err != nil {
		wails.LogErrorf(ctx, "Check update faild: %s", err)
	}

	switch version {
	case "-1":
		wails.LogInfo(ctx, "It is special release version, no need to update")
	case "0":
		wails.LogInfo(ctx, "No software update")
	default:
		wails.LogInfof(ctx, "Founded new version: %s", version)
		// Show dialog to user
		{
			result, err := wails.MessageDialog(ctx, wails.MessageDialogOptions{
				Type:          wails.QuestionDialog,
				Title:         "找到新版本：" + version,
				Message:       "软件有新版本发布了，是否前往下载？",
				DefaultButton: "Yes",
			})

			if err != nil {
				wails.LogError(ctx, "Failed to show message dialog: "+err.Error())
			}

			wails.LogDebugf(ctx, "Dialog result：%s", result)

			if result == "Yes" {
				wails.BrowserOpenURL(ctx, "https://github.com/HIM049/BADownloaderUI/releases/tag/"+version)
			}
		}
	}
}

// 初始化必备文件夹
func initFolder(ctx context.Context) {
	downloadPath := config.Cfg.GetDownloadPath()
	cachePath := config.Cfg.GetCachePath()

	var paths = []string{
		downloadPath,
		cachePath,
		cachePath + "/audio",
		cachePath + "/cover",
	}

	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			wails.LogFatalf(ctx, "Initialize Folder Failed: %s", err)
		}
	}
	wails.LogInfo(ctx, "Initialize Folder Successful")
}
