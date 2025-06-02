package download

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/ffmpeg"
	"bili-audio-downloader/backend/utils"
	"bili-audio-downloader/bilibili"
	"errors"
	"fmt"
	"sync"
)

type Audio struct {
	auid     string
	coverUrl string
	sessdata string
	listId   int
	option   Option
	path     Path
	metaData MetaData
}

func (a *Audio) SetID(id int) {
	a.listId = id
}

func NewAudio(auid string, coverUrl, sessdata string, metaData MetaData) *Audio {
	return &Audio{
		auid:     auid,
		coverUrl: coverUrl,
		sessdata: sessdata,
		path: Path{
			StreamPath:   fmt.Sprintf("%s/audio/%s", config.Cfg.GetCachePath(), auid),
			CoverPath:    fmt.Sprintf("%s/cover/%s.jpg", config.Cfg.GetCachePath(), auid),
			CurrentPath:  fmt.Sprintf("%s/audio/%s", config.Cfg.GetCachePath(), auid),
			OutputFormat: constants.AudioType.M4a,
		},
		metaData: metaData,
	}
}

func (a *Audio) Download() error {
	var wg sync.WaitGroup
	errorResults := make(chan error, 2)

	wg.Add(2)
	// 下载流媒体

	go func() {
		defer wg.Done()

		var err error
		for i := 0; i < config.Cfg.DownloadConfig.RetryCount; i++ {
			err = a.downloadStream(a.path.StreamPath)
			if err != nil {
				err = errors.New(fmt.Sprintf("failed to download video stream: %v (retry %d)", err, i))
				continue

			}
			if !utils.IsFileExists(a.path.StreamPath) {
				err = errors.New(fmt.Sprintf("failed to download video stream: %s (retry %d)", "file not exitsts", i))
				continue
			}
			break
		}
		if err != nil {
			errorResults <- err
		}

	}()

	// 下载封面
	go func() {
		defer wg.Done()

		var err error
		for i := 0; i < config.Cfg.DownloadConfig.RetryCount; i++ {
			err = utils.SaveFromURL(a.coverUrl, a.path.CoverPath)
			if err != nil {
				err = errors.New(fmt.Sprintf("failed to download video stream: %v (retry %d)", err, i))
				continue
			}
			if !utils.IsFileExists(a.path.CoverPath) {
				err = errors.New(fmt.Sprintf("failed to download video stream: %s (retry %d)", "file not exitsts", i))
				continue
			}
			break
		}
		if err != nil {
			errorResults <- err
		}

	}()

	wg.Wait()
	close(errorResults)

	if len(errorResults) > 0 {
		return <-errorResults
	}

	return nil
}

func (a *Audio) ConventFormat() error {
	newPath := fmt.Sprintf("%s.c", a.path.StreamPath)
	err := ffmpeg.ConvertToMP3(a.path.CurrentPath, newPath)
	if err != nil {
		return err
	}
	a.path.OutputFormat = constants.AudioType.Mp3
	a.path.CurrentPath = newPath
	return nil
}

func (a *Audio) WriteMetadata() error {
	newPath := fmt.Sprintf("%s.meta", a.path.StreamPath)
	err := ffmpeg.WriteMetadata(a.path.CurrentPath, newPath, a.path.CoverPath, a.metaData.SongName, a.metaData.Author, a.path.OutputFormat)
	if err != nil {
		return err
	}
	a.path.CurrentPath = newPath
	return nil
}

func (a *Audio) ExportFile() error {
	err := ExportFile(a.metaData.Title, a.metaData.PageTitle, a.path.OutputFormat, a.listId, a.path.CurrentPath)
	if err != nil {
		return err
	}
	return nil
}

func (a *Audio) downloadStream(streamPath string) error {
	// 请求流地址
	audio := new(bilibili.Audio)
	audio.Auid = a.auid
	err := audio.GetStream(a.sessdata)
	if err != nil {
		return errors.New(fmt.Sprintf("get audio stream error: %s", err))
	}

	// 通过流地址下载
	err = utils.StreamingDownloader(audio.Stream.StreamLink, streamPath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to download stream: %s", err))
	}
	return nil
}

func (a *Audio) GetTaskInfo() *TaskInfo {
	taskInfo := TaskInfo{
		SongName:   a.metaData.SongName,
		SongAuthor: a.metaData.Author,
		CoverUrl:   a.coverUrl,
	}
	return &taskInfo
}
