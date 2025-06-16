package wails_api

import (
	"bili-audio-downloader/backend/config"
	"github.com/spf13/viper"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// ResetConfig 重置设置文件
func (w *WailsApi) ResetConfig() {
	cfg := config.DefaultConfig()
	err := cfg.UpdateAndSave()
	if err != nil {
		wails.LogErrorf(w.ctx, "写入设置文件失败：%s", err)
		wails.EventsEmit(w.ctx, "error", "写入设置时出错:"+err.Error())
		// TODO 增加统一前端错误提示接口
	}
}

// LoadConfig 读取设置
func (w *WailsApi) LoadConfig() config.Config {
	return config.Cfg
}

// SaveConfig 写入设置
func (w *WailsApi) SaveConfig(cfg config.Config) {
	err := cfg.UpdateAndSave()
	if err != nil {
		wails.LogErrorf(w.ctx, "写入设置文件失败：%s", err)
		wails.EventsEmit(w.ctx, "error", "写入设置时出错:"+err.Error())
	}
}

// RefreshConfig 刷新设置
func (w *WailsApi) RefreshConfig() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
