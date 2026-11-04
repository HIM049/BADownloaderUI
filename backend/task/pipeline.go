package task

import (
	"bili-audio-downloader/backend/biliapi"
	"bili-audio-downloader/backend/constants"
	"bili-audio-downloader/backend/ffmpeg"
	"bili-audio-downloader/backend/utils"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Settings 一次下载使用的配置
type Settings struct {
	CacheDir   string
	RetryCount int
	Threads    int
	Convert    bool
}

// Progress downloadFinish 事件携带的任务信息
type Progress struct {
	Index int    `json:"index"`
	Title string `json:"title"`
}

// Run 并发处理列表中的任务，每个任务结束（含已删除）时发出 downloadFinish 事件
func (l *List) Run(ctx context.Context, s Settings) {
	sem := make(chan struct{}, s.Threads+1)
	var wg sync.WaitGroup

	wails.LogPrintf(ctx, "start list download process")
	for i, t := range l.tasks {
		progress := Progress{Index: i, Title: t.Meta.SongName}
		if t.Deleted {
			wails.EventsEmit(ctx, "downloadFinish", progress)
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer func() {
				<-sem
				wg.Done()
				wails.EventsEmit(ctx, "downloadFinish", progress)
			}()

			if err := t.process(i, s); err != nil {
				wails.LogErrorf(ctx, "task %d failed: %v", i, err)
			}
		}()
	}
	wg.Wait()
}

// process 下载流和封面，按需转码，按选项写入元数据并导出
func (t *Task) process(index int, s Settings) error {
	name := fmt.Sprintf("%d_%s", index, t.cacheID())
	streamPath := filepath.Join(s.CacheDir, "audio", name)
	coverPath := filepath.Join(s.CacheDir, "cover", name+".jpg")

	var format string
	var streamErr, coverErr error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		streamErr = retry(s.RetryCount, func() (err error) {
			format, err = t.downloadStream(streamPath)
			return err
		})
	}()
	go func() {
		defer wg.Done()
		if !t.Option.SongCover {
			return
		}
		coverErr = retry(s.RetryCount, func() error {
			return biliapi.Download(t.CoverURL, coverPath)
		})
	}()
	wg.Wait()
	if err := errors.Join(streamErr, coverErr); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	current := streamPath
	if s.Convert && format == constants.AudioType.M4a {
		converted := streamPath + "_conv" + constants.AudioType.Mp3
		if err := ffmpeg.ConvertToMP3(current, converted); err != nil {
			return fmt.Errorf("convert: %w", err)
		}
		current, format = converted, constants.AudioType.Mp3
	}

	// 元数据步骤同时把 DASH 流封装为目标格式，选项全部关闭时也要执行
	var songName, author string
	if t.Option.SongName {
		songName = t.Meta.SongName
	}
	if t.Option.SongAuthor {
		author = t.Meta.Author
	}
	if !t.Option.SongCover {
		coverPath = ""
	}
	tagged := streamPath + "_meta" + format
	if err := ffmpeg.WriteMetadata(current, tagged, coverPath, songName, author, format); err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	if err := utils.ExportFile(t.Meta.Title, t.Meta.PageTitle, format, index, tagged); err != nil {
		return fmt.Errorf("export: %w", err)
	}
	return nil
}

// downloadStream 获取流地址并下载，返回音频格式；视频有 FLAC 流时优先使用
func (t *Task) downloadStream(dst string) (string, error) {
	url, format := "", constants.AudioType.M4a
	switch t.Kind {
	case KindVideo:
		play, err := biliapi.GetPlayURL(t.Bvid, t.Cid, t.Sessdata)
		if err != nil {
			return "", err
		}
		url = play.AudioURL()
		if flac := play.FlacURL(); flac != "" {
			url, format = flac, constants.AudioType.Flac
		}
	case KindAudio:
		stream, err := biliapi.GetAudioStream(t.Auid, t.Sessdata)
		if err != nil {
			return "", err
		}
		// TODO: stream.Type 为 3 时是 FLAC，目前统一按 m4a 处理
		url = stream.URL()
	}
	if url == "" {
		return "", errors.New("stream url is empty")
	}
	return format, biliapi.Download(url, dst)
}

// retry 执行 fn 直到成功，至少执行一次，失败后等待时间逐次增加
func retry(attempts int, fn func() error) error {
	attempts = max(attempts, 1)
	var err error
	for i := range attempts {
		if i > 0 {
			time.Sleep(time.Duration(i) * time.Second)
		}
		if err = fn(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts: %w", attempts, err)
}
