package main

import (
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 下载列表中歌曲的函数（参数直接读取 config ）
func (a *App) DownloadList() error {
	cfg := GetConfig()
	runtime.LogInfo(a.ctx, "开始下载任务列表")
	sem := make(chan struct{}, cfg.DownloadThreads+1)
	var wg sync.WaitGroup

	// 获取任务队列
	var list []VideoInformationList

	err := LoadJsonFile(cfg.VideoListPath, &list)
	if err != nil {
		return err
	}
	// 遍历下载队列
	for _, video := range list {

		go func(v VideoInformationList) {
			// fmt.Println("调用下载")
			// 下载完成后
			defer func() {
				// progressBar.Increment()
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知
			}()

			sem <- struct{}{} // 给通道中
			wg.Add(1)         // 任务 +1

			err := GetAndDownload(v.Bvid, v.Cid, cfg.CachePath+"/music/"+strconv.Itoa(v.Cid))
			if err != nil {
				runtime.LogError(a.ctx, "下载时出现错误："+err.Error())
			}

		}(video)
	}
	// 等待任务执行完成
	wg.Wait()
	return nil
}

// 获取并下载媒体流
func GetAndDownload(bvid string, cid int, filePathAndName string) error {
	// 获取 B 站视频流地址
	video, err := GetVideoObj(bvid, cid)
	if err != nil {
		return err
	}
	// 下载媒体流
	err = StreamingDownloader(video.Data.Dash.Audio[0].BaseUrl, filePathAndName)
	if err != nil {
		return err
	}
	return nil
}
