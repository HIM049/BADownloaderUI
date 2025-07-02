package download

import (
	"bili-audio-downloader/backend/adapter"
)

type DownloadTask interface {
	SetID(int)
	Download() error
	ConventFormat() error
	WriteMetadata() error
	ExportFile() error
	GetTaskInfo() *adapter.TaskInfo
}
