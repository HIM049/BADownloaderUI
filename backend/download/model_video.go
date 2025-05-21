package download

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/ffmpeg"
	"bili-audio-downloader/backend/utils"
	"bili-audio-downloader/bilibili"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"sync"
)

type Video struct {
	bvid     string
	cid      int
	coverUrl string
	sessdata string
	listId   int
	option   Option
	path     Path
	metaData MetaData
}

func NewVideo(bvid string, cid int, coverUrl, sessdata string, metaData MetaData) *Video {
	return &Video{
		bvid:     bvid,
		cid:      cid,
		coverUrl: coverUrl,
		sessdata: sessdata,
		path: Path{
			StreamPath:   fmt.Sprintf("%s/music/%d", config.Cfg.GetCachePath(), cid),
			CoverPath:    fmt.Sprintf("%s/cover/%d.jpg", config.Cfg.GetCachePath(), cid),
			CurrentPath:  fmt.Sprintf("%s/music/%d", config.Cfg.GetCachePath(), cid),
			OutputFormat: constants.AudioType.M4a,
		},
		metaData: metaData,
	}
}

func (v *Video) SetID(id int) {
	v.listId = id
}

func (v *Video) Download() error {
	var wg sync.WaitGroup
	errorResults := make(chan error, 2)

	wg.Add(2)
	// 下载流媒体

	go func() {
		defer wg.Done()

		var err error
		for i := 0; i < config.Cfg.DownloadConfig.RetryCount; i++ {
			err = v.downloadStream(v.path.StreamPath)
			if err != nil {
				err = errors.New(fmt.Sprintf("failed to download video stream: %v (retry %d)", err, i))
				continue
			}
			if !utils.IsFileExists(v.path.StreamPath) {
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
			err = utils.SaveFromURL(v.coverUrl, v.path.CoverPath)
			if err != nil {
				err = errors.New(fmt.Sprintf("failed to download video stream: %v (retry %d)", err, i))
				continue
			}
			if !utils.IsFileExists(v.path.CoverPath) {
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
	fmt.Println("等待结束")
	close(errorResults)

	if len(errorResults) > 0 {
		return <-errorResults
	}
	fmt.Println("下载完成")

	return nil
}

func (v *Video) downloadStream(streamPath string) error {
	// 请求流地址
	json, err := bilibili.GetVideoStream(v.bvid, strconv.Itoa(v.cid), v.sessdata)
	if gjson.Get(json, "code").Int() != 0 {
		return errors.New("return value is not 0")
	}

	// 获取流地址
	stream := gjson.Get(json, "data.dash.audio.0.base_url").String()
	flacStream := gjson.Get(json, "data.dash.flac.audio.base_url").String()
	if v.option.DownloadFlac && flacStream != "" {
		stream = flacStream
		v.path.OutputFormat = constants.AudioType.Flac
	}

	if stream == "" {
		return errors.New("return stream is empty")
	}

	// 通过流地址下载
	err = utils.StreamingDownloader(stream, streamPath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to download stream: %s", err))
	}
	return nil
}

func (v *Video) ConventFormat() error {
	// 根据目标格式转码
	newPath := fmt.Sprintf("%s.c", v.path.StreamPath)
	if v.path.OutputFormat == constants.AudioType.M4a {
		err := ffmpeg.ConvertToMP3(v.path.CurrentPath, newPath)
		if err != nil {
			return err
		}
		v.path.OutputFormat = constants.AudioType.Mp3
	} else {
		err := ffmpeg.ConvertToMP3(v.path.CurrentPath, newPath)
		if err != nil {
			return err
		}
		v.path.OutputFormat = constants.AudioType.Flac
	}
	v.path.CurrentPath = newPath
	return nil
}

func (v *Video) WriteMetadata() error {
	newPath := fmt.Sprintf("%s.meta", v.path.StreamPath)
	err := ffmpeg.WriteMetadata(v.path.CurrentPath, newPath, v.path.CoverPath, v.metaData.SongName, v.metaData.Author, v.path.OutputFormat)
	if err != nil {
		return err
	}
	v.path.CurrentPath = newPath
	return nil
}

func (v *Video) ExportFile() error {
	err := ExportFile(v.metaData.Title, v.metaData.PageTitle, v.path.OutputFormat, v.listId, v.path.CurrentPath)
	if err != nil {
		return err
	}
	return nil
}
