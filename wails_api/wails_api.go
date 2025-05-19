package wails_api

import (
	"bili-audio-downloader/Download"
	"context"
)

type WailsApi struct {
	ctx context.Context
}

func (a *WailsApi) DownloadTaskList() {
	Download.DownloadTaskList(a.ctx)
}
