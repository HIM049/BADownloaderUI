package wails_api

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/ffmpeg"
	"os/exec"
)

// GetAppVersion 获取版本号
func (w *WailsApi) GetAppVersion() string {
	return constants.APP_VERSION
}

// GetTheme 获取主题字符串
func (w *WailsApi) GetTheme() string {
	return config.Cfg.Theme
}

func (w *WailsApi) Checkffmpeg() bool {
	return ffmpeg.CheckExists()
}

// InstallFfmpegByWinget 使用 winget 安装 ffmpeg
func (w *WailsApi) InstallFfmpegByWinget() error {
	cmd := exec.Command("cmd", "/c", "winget", "install", "Gyan.FFmpeg",
		"--silent",
		"--accept-source-agreements",
		"--accept-package-agreements")
	_, err := cmd.CombinedOutput()
	return err
}
