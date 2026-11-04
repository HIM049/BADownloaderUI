package wails_api

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/utils"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// OpenFileDialog 调用打开文件窗口
func (w *WailsApi) OpenFileDialog() (string, error) {
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
	path, err := wails.OpenFileDialog(w.ctx, option)
	if err != nil {
		wails.LogErrorf(w.ctx, "OpenFileDialog error: %s", err.Error())
		return "", err
	}

	return path, nil
}

// SetDownloadPathDialog 选择并保存下载路径，取消选择时不做修改
func (w *WailsApi) SetDownloadPathDialog() error {
	path, err := wails.OpenDirectoryDialog(w.ctx, wails.OpenDialogOptions{
		DefaultDirectory: "./",
		Title:            "选择下载路径",
	})
	if err != nil || path == "" {
		return err
	}
	config.Cfg.FileConfig.DownloadPath = path
	return config.Cfg.UpdateAndSave()
}

// 打开下载文件夹
func (w *WailsApi) OpenDownloadFolader() error {

	err := utils.OpenFolder(config.Cfg.GetDownloadPath())
	if err != nil {
		return err
	}

	return nil
}
