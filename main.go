package main

import (
	"bili-audio-downloader/services"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

// 全局版本号
const APP_VERSION string = "4.8.3"
const CONFIG_VERSION int = 2

func main() {
	// Create an instance of the app structure
	app := &App{}

	// Init logger
	customLogger, err := services.NewCustomLogger()
	if err != nil {
		println("Error:", err.Error())
		return
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "BiliAudioDownloader " + APP_VERSION,
		Width:  1024,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:     true, // 无边框窗口
		DisableResize: true, // 窗口尺寸
		Windows: &windows.Options{
			IsZoomControlEnabled: false, // 页面缩放比例
		},
		BackgroundColour:   &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:          app.startup,
		OnShutdown:         app.shutdown,
		LogLevelProduction: logger.INFO,
		Logger:             customLogger,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

}
