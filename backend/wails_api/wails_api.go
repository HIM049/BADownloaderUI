package wails_api

import (
	"bili-audio-downloader/backend/download"
	"context"
)

type WailsApi struct {
	ctx context.Context
}

func (w *WailsApi) Startup(ctx context.Context) {
	w.ctx = ctx
}

func (a *WailsApi) DownloadTaskList() {
	download.DownloadTaskList(a.ctx)
}
