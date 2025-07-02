package main

import (
	"bili-audio-downloader/backend/config"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

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

//// 调用保存窗口
//func (a *App) SaveVideoListTo(videolist services.VideoList) error {
//	var FileFilter []wails.FileFilter
//
//	fileFilter := wails.FileFilter{
//		DisplayName: "视频下载列表 (*.json)",
//		Pattern:     "*.json",
//	}
//	FileFilter = append(FileFilter, fileFilter)
//
//	option := wails.SaveDialogOptions{
//		DefaultDirectory: "./",
//		DefaultFilename:  "BAD_VideoList",
//		Title:            "另存视频列表",
//		Filters:          FileFilter,
//	}
//
//	// 弹出对话框
//	path, err := wails.SaveFileDialog(a.ctx, option)
//	if err != nil {
//		return err
//	}
//
//	// 用户取消操作
//	if path == "" {
//		wails.EventsEmit(a.ctx, "error", "未选择保存路径")
//		return nil
//	}
//
//	// 保存列表
//	err = videolist.Save(path)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// 打开下载文件夹
// TODO
func (a *App) OpenDownloadFolader() error {

	//err := OpenFolder(config.Cfg.GetDownloadPath())
	//if err != nil {
	//	return err
	//}

	return nil
}
