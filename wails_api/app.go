package wails_api

import (
	"bili-audio-downloader/config"
	"bili-audio-downloader/constants"
)

// GetAppVersion 获取版本号
func (a *WailsApi) GetAppVersion() string {
	return constants.APP_VERSION
}

// GetTheme 获取主题字符串
func (a *WailsApi) GetTheme() string {
	return config.Cfg.Theme
}
