package models

import (
	"bili-audio-downloader/bilibili"
	"bili-audio-downloader/config"
	"bili-audio-downloader/constants"
	"bili-audio-downloader/utils"
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

type MetaData struct {
	Title     string
	PageTitle string
	PartId    int
	SongName  string
	Author    string
	LyricsUrl string
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
			ConventPath:  fmt.Sprintf("%s/convented/%d", config.Cfg.GetCachePath(), cid),
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

//func (v *Video) ConventToMp3() error {
//	services.ConventFile(v.path.StreamPath)
//}

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

//func (v *Video) getFileName() error {
//	// 处理文件名结构体
//	fileName := new(services.FileName)
//	fileName.Title = v.metaData.Title
//	fileName.Subtitle = v.metaData.PageTitle
//	fileName.ID = v.listId
//	fileName.Quality = "" // TODO
//
//	// 处理模板和生成文件名
//	tmpl, err := template.New("filename").Parse(config.Cfg.FileConfig.FileNameTemplate)
//	if err != nil {
//		return err
//	}
//
//	var output bytes.Buffer
//	err = tmpl.Execute(&output, fileName)
//	if err != nil {
//		return err
//	}
//	return nil
//}
