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
	config.InitConfig(a.ctx)
	initFolder(a.ctx)
	checkUpdateAndAlarm(a.ctx)

}

// 程序关闭时
func (a *App) shutdown(ctx context.Context) {
	// 清理缓存
	if config.Cfg.DeleteCache {
		os.RemoveAll(config.Cfg.GetCachePath())
	}
}

// 检查更新并提示
func checkUpdateAndAlarm(ctx context.Context) {
	// 检查版本更新
	version, err := services.CheckUpdate(constants.APP_VERSION)
	if err != nil {
		wails.LogErrorf(ctx, "Check update faild: %s", err)
	} else if version == "0" {
		wails.LogInfo(ctx, "No software update")
	} else {
		wails.LogInfof(ctx, "Founded new version: %s", version)

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
