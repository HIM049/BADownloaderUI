package wails_api

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/ffmpeg"
)

// GetAppVersion 获取版本号
func (a *WailsApi) GetAppVersion() string {
	return constants.APP_VERSION
}

// GetTheme 获取主题字符串
func (a *WailsApi) GetTheme() string {
	return config.Cfg.Theme
}

func (a *WailsApi) Checkffmpeg() bool {
	return ffmpeg.CheckExists()
}
