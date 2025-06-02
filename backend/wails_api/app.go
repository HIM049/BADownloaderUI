package wails_api

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/ffmpeg"
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
