package wails_api

import (
	"bili-audio-downloader/backend/config"
)

// ResetConfig 重置设置文件
func (w *WailsApi) ResetConfig() error {
	return config.DefaultConfig().UpdateAndSave()
}

// LoadConfig 读取设置
func (w *WailsApi) LoadConfig() config.Config {
	return config.Cfg
}

// SaveConfig 写入设置
func (w *WailsApi) SaveConfig(cfg config.Config) error {
	return cfg.UpdateAndSave()
}

// RefreshConfig 刷新设置
func (w *WailsApi) RefreshConfig() error {
	return config.LoadConfig(w.ctx)
}
