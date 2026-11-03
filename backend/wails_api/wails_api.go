package wails_api

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/task"
	"context"
)

type WailsApi struct {
	ctx   context.Context
	tasks task.List
}

func (w *WailsApi) Startup(ctx context.Context) {
	w.ctx = ctx
}

func (w *WailsApi) DownloadTaskList() {
	w.tasks.Run(w.ctx, task.Settings{
		CacheDir:   config.Cfg.GetCachePath(),
		RetryCount: config.Cfg.DownloadConfig.RetryCount,
		Threads:    config.Cfg.DownloadConfig.DownloadThreads,
		Convert:    config.Cfg.FileConfig.ConvertFormat,
	})
}

// accountSessdata 返回启用账号且已登录时的 SESSDATA，否则返回空字符串
func accountSessdata() string {
	if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
		return config.Cfg.Account.SESSDATA
	}
	return ""
}
