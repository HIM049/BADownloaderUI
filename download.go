package main

import (
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) DownloadList(videoListPath string, threads int) error {
	runtime.LogInfo(a.ctx, "开始下载任务列表")
	sem := make(chan struct{}, threads+1)
	var wg sync.WaitGroup
	// var progressBar *pb.ProgressBar

	// 获取任务队列
	var list []VideoInformationList
	err := LoadJsonFile(videoListPath, &list)
	if err != nil {
		return err
	}
	// 设置进度条
	// progressBar = pb.Full.Start(len(list))
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

			err := GetAndDownload(v.Bvid, v.Cid, "C:/Users/HIM/Desktop/Download/"+strconv.Itoa(v.Cid)+".mp3")
			if err != nil {
				runtime.LogError(a.ctx, "下载时出现错误："+err.Error())
			}

		}(video)
	}
	// 等待任务执行完成
	wg.Wait()
	// progressBar.Finish()
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
