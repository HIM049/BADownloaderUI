package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// 全局版本号
const APP_VERSION string = "4.0.0-pre"

func main() {
	// Create an instance of the app structure
	app := &App{}

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "BiliAudioDownloader " + APP_VERSION,
		Width:  1024,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:        true, // 无边框窗口
		DisableResize:    true, // 禁用缩放
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		// OnStartup:        beforeRunFunc,
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

}
