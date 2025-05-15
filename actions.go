package main

import (
	"bili-audio-downloader/bilibili"
	"bili-audio-downloader/config"
	"bili-audio-downloader/services"
	"errors"
	"strconv"

	"github.com/spf13/viper"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// GetUsersCollect 获取用户创建的收藏夹
func (a *App) GetUsersCollect() bilibili.Collects {
	// 获取收藏夹列表
	collects := new(bilibili.Collects)
	mid, _ := strconv.Atoi(config.Cfg.Account.DedeUserID)
	collects.UserMid = mid
	err := collects.GetUsersCollect(config.Cfg.Account.SESSDATA)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取收藏夹列表失败：%s", err)
		return bilibili.Collects{}
	}

	return *collects
}

// GetFavCollect 获取收藏的收藏夹
func (a *App) GetFavCollect(pn int) bilibili.Collects {
	// 获取收藏夹列表
	collects := new(bilibili.Collects)
	mid, _ := strconv.Atoi(config.Cfg.Account.DedeUserID)
	collects.UserMid = mid
	err := collects.GetFavCollect(config.Cfg.Account.SESSDATA, 20, pn)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取收藏夹列表失败：%s", err)
		return bilibili.Collects{}
	}

	return *collects
}

// OpenFileDialog 调用打开文件窗口
func (a *App) OpenFileDialog() (string, error) {
	var FileFilter []wails.FileFilter

	fileFilter := wails.FileFilter{
		DisplayName: "视频下载列表 (*.json)",
		Pattern:     "*.json",
	}
	FileFilter = append(FileFilter, fileFilter)

	option := wails.OpenDialogOptions{
		DefaultDirectory: "./",
		DefaultFilename:  "",
		Title:            "打开本地列表文件",
		Filters:          FileFilter,
	}
	// 弹出对话框
	path, err := wails.OpenFileDialog(a.ctx, option)
	if err != nil {
		wails.LogErrorf(a.ctx, err.Error())
		return "", err
	}

	return path, nil
}

func (a *App) SetDownloadPathDialog() {

	option := wails.OpenDialogOptions{
		DefaultDirectory: "./",
		DefaultFilename:  "",
		Title:            "选择下载路径",
	}

	path, err := wails.OpenDirectoryDialog(a.ctx, option)
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
	}

	config.Cfg.FileConfig.DownloadPath = path
	err = config.Cfg.UpdateAndSave()
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
	}

}

// 调用保存窗口
func (a *App) SaveVideoListTo(videolist services.VideoList) error {
	var FileFilter []wails.FileFilter

	fileFilter := wails.FileFilter{
		DisplayName: "视频下载列表 (*.json)",
		Pattern:     "*.json",
	}
	FileFilter = append(FileFilter, fileFilter)

	option := wails.SaveDialogOptions{
		DefaultDirectory: "./",
		DefaultFilename:  "BAD_VideoList",
		Title:            "另存视频列表",
		Filters:          FileFilter,
	}

	// 弹出对话框
	path, err := wails.SaveFileDialog(a.ctx, option)
	if err != nil {
		return err
	}

	// 用户取消操作
	if path == "" {
		wails.EventsEmit(a.ctx, "error", "未选择保存路径")
		return nil
	}

	// 保存列表
	err = videolist.Save(path)
	if err != nil {
		return err
	}
	return nil
}

// 获取已登录用户的信息
func (a *App) GetUserInf() (bilibili.AccountInformation, error) {
	if !config.Cfg.Account.IsLogin {
		return bilibili.AccountInformation{}, errors.New("用户未登录")
	}
	sessdata := config.Cfg.Account.SESSDATA

	accountInf := new(bilibili.AccountInformation)
	accountInf.GetUserInf(sessdata)

	return *accountInf, nil
}

// 重置设置文件
func (a *App) ResetConfig() {
	cfg := config.DefaultConfig()
	err := cfg.UpdateAndSave()
	if err != nil {
		wails.LogErrorf(a.ctx, "写入设置文件失败：%s", err)
		wails.EventsEmit(a.ctx, "error", "写入设置时出错:"+err.Error())

	}
}

// 读取设置
func (a *App) LoadConfig() config.Config {
	return config.Cfg
}

// 写入设置
func (a *App) SaveConfig(cfg config.Config) {
	err := cfg.UpdateAndSave()
	if err != nil {
		wails.LogErrorf(a.ctx, "写入设置文件失败：%s", err)
		wails.EventsEmit(a.ctx, "error", "写入设置时出错:"+err.Error())
	}
}

func (a *App) RefreshConfig() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

// 打开下载文件夹
func (a *App) OpenDownloadFolader() error {

	err := OpenFolder(config.Cfg.GetDownloadPath())
	if err != nil {
		return err
	}

	return nil
}
