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
	option   Option
	path     Path
	metaData MetaData
}

func NewAudio(auid string, coverUrl, sessdata string, metaData MetaData) *Audio {
	return &Audio{
		auid:     auid,
		coverUrl: coverUrl,
		sessdata: sessdata,
		path: Path{
			StreamPath: fmt.Sprintf("%s/audio/%s", config.Cfg.GetCachePath(), auid),
			CoverPath:  fmt.Sprintf("%s/cover/%s.jpg", config.Cfg.GetCachePath(), auid),
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

func (v *Audio) ConventFormat() error {
	err := ffmpeg.ConvertToMP3(v.path.StreamPath, v.path.StreamPath)
	if err != nil {
		return err
	}
	v.path.OutputFormat = constants.AudioType.Mp3
	return nil
}

func (v *Audio) WriteMetadata() error {
	isMp3 := v.path.OutputFormat == constants.AudioType.Mp3
	err := ffmpeg.WriteMetadata(v.path.StreamPath, v.path.StreamPath, v.path.CoverPath, v.metaData.SongName, v.metaData.Author, isMp3)
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
