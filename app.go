package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

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

	// 程序初始化
	runtime.LogInfo(a.ctx, "正在创建文件夹")

	cfg := GetConfig()
	_ = os.MkdirAll(cfg.DownloadPath, 0755)
	_ = os.MkdirAll(cfg.CachePath, 0755)
	_ = os.MkdirAll(cfg.CachePath+"/music", 0755)
	_ = os.MkdirAll(cfg.CachePath+"/cover", 0755)
}

// 程序关闭时
func (a *App) shutdown(ctx context.Context) {
	// 清理缓存
	cfg := GetConfig()
	os.RemoveAll(cfg.CachePath)
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

type DownloadOption struct {
	SongName   bool `json:"song_name"`
	SongCover  bool `json:"song_cover"`
	SongAuthor bool `json:"song_author"`
}

// 执行下载操作时
func (a *App) StartDownload(DownOpt DownloadOption) {
	// fmt.Println(DownOpt)
	cfg := GetConfig()

	// 下载歌曲
	runtime.LogInfo(a.ctx, "开始下载")
	err := DownloadList(a.ctx, &cfg)
	if err != nil {
		runtime.LogError(a.ctx, "下载错误："+err.Error())
	}
	runtime.LogInfo(a.ctx, "下载完成")

	func() {
		fmt.Println("aaaa?")
	}()

	// 写入歌曲元数据
	runtime.LogInfo(a.ctx, "开始写入元数据")
	err = ConcurrentChangeTag(&cfg, &DownOpt, ".m4a")
	if err != nil {
		runtime.LogError(a.ctx, "写入歌曲元数据时发生错误："+err.Error())
	}
	runtime.LogInfo(a.ctx, "元数据写入完成")

	// // 改名并输出到下载文件夹
	// runtime.LogInfo(a.ctx, "开始输出")
	// err = ConcurrentChangeName(cfg.DownloadThreads, cfg.VideoListPath, ".m4a", cfg.CachePath+"/music/", cfg.DownloadPath+"/")

	// if err != nil {
	// 	runtime.LogError(a.ctx, "输出文件时发生错误："+err.Error())
	// }

}

func (a *App) MakeUpEditor() {
	cfg := GetConfig()

	// 获取系统默认的文本编辑器
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "notepad"
	}

	// 创建命令
	cmd := exec.Command(editor, cfg.VideoListPath)

	// 启动命令
	err := cmd.Start()
	if err != nil {
		fmt.Println("无法启动编辑器:", err)
		return
	}
}
