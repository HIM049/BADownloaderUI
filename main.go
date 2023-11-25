package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "BiliAudioDownloader 3.0",
		Width:  520,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 236, G: 236, B: 236, A: 1},
		OnStartup:        beforeRunFunc,
		// OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

	// createFolder()
}

// func createFolder() {
// 	fmt.Println("正在创建文件夹")
// 	cfg := GetConfig()

// 	_ = os.MkdirAll(cfg.DownloadPath, 0755)
// 	_ = os.MkdirAll(cfg.CachePath, 0755)
// 	_ = os.MkdirAll(cfg.CachePath+"/music", 0755)
// 	_ = os.MkdirAll(cfg.CachePath+"/cover", 0755)
// }

// func init() {
// 	fmt.Println("正在创建文件夹")
// 	cfg := GetConfig()

// 	_ = os.MkdirAll(cfg.DownloadPath, 0755)
// 	_ = os.MkdirAll(cfg.CachePath, 0755)
// 	_ = os.MkdirAll(cfg.CachePath+"/music", 0755)
// 	_ = os.MkdirAll(cfg.CachePath+"/cover", 0755)
// }

func beforeRunFunc(ctx context.Context) {
	fmt.Println("正在创建文件夹")
	cfg := GetConfig()

	_ = os.MkdirAll(cfg.DownloadPath, 0755)
	_ = os.MkdirAll(cfg.CachePath, 0755)
	_ = os.MkdirAll(cfg.CachePath+"/music", 0755)
	_ = os.MkdirAll(cfg.CachePath+"/cover", 0755)
}
