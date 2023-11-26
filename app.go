package main

import (
	"context"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	beforeRunFunc()
}

// Greet returns a greeting for the given name
// func (a *App) Greet(name string) string {
// 	return fmt.Sprintf("Hello %s, It's show time!", name)
// }

// func (a *App) CheckFavListInformation(favID string) string {
// 	favPackage, err := GetFavListObj(favID, 1, 1)
// 	if err != nil {
// 		wails.LogError(a.ctx, err.Error())
// 	}
// 	return favPackage.Data.Info.Title
// }

func beforeRunFunc() {
	fmt.Println("正在创建文件夹")
	cfg := GetConfig()

	_ = os.MkdirAll(cfg.DownloadPath, 0755)
	_ = os.MkdirAll(cfg.CachePath, 0755)
	_ = os.MkdirAll(cfg.CachePath+"/music", 0755)
	_ = os.MkdirAll(cfg.CachePath+"/cover", 0755)
}

// 查询并返回收藏夹信息
func (a *App) SearchFavListInformation(favListID string) FavList {
	listInf, err := GetFavListObj(favListID, 1, 1)
	if err != nil {
		runtime.LogError(a.ctx, "获取收藏夹内容时出现错误："+err.Error())
		// fmt.Printf("获取收藏夹内容时出现错误：%s", err)
		return FavList{}
	}
	return *listInf
}
