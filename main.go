package main

import (
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/services"
	"bili-audio-downloader/backend/wails_api"
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := &App{}
	wailsApi := &wails_api.WailsApi{}

	// Init logger
	customLogger, err := services.NewCustomLogger()
	if err != nil {
		println("Error:", err.Error())
		return
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "BiliAudioDownloader " + constants.APP_VERSION,
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
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			wailsApi.Startup(ctx)
		},
		OnShutdown:         app.shutdown,
		LogLevelProduction: logger.INFO,
		Logger:             customLogger,
		Bind: []interface{}{
			app,
			wailsApi,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

}
