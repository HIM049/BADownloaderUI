package wails_api

import (
	"bili-audio-downloader/Download"
	"context"
	"fmt"
)

type WailsApi struct {
	ctx context.Context
}

func (w *WailsApi) Startup(ctx context.Context) {
	w.ctx = ctx
	fmt.Println("ctx: ", w.ctx)
}

func (a *WailsApi) DownloadTaskList() {
	Download.DownloadTaskList(a.ctx)
}
