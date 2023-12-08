package main

import (
	"context"
	"strconv"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 下载列表中歌曲的函数
func DownloadList(ctx context.Context, cfg *Config) error {
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

			sem <- struct{}{} // 给通道中填入数据
			wg.Add(1)         // 任务 +1

			// TODO：下载重试次数可设定

			// 下载重试
			for i := 0; i < 10; i++ {
				err := GetAndDownload(v.Bvid, v.Cid, cfg.CachePath+"/music/"+strconv.Itoa(v.Cid)+".m4a")
				if err == nil {
					// 下载成功
					break
				}
				runtime.LogErrorf(ctx, "下载时出现错误：%s (重试 %d )\n", err, i+1)
			}

			// 下载完成后
			defer func() {
				runtime.LogInfo(ctx, "下载成功")
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知
			}()
		}(video)

		go func(v VideoInformationList) {
			// 下载完成后
			defer func() {
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知
			}()

			sem <- struct{}{} // 给通道中填入数据
			wg.Add(1)         // 任务 +1

			err = SaveFromURL(v.Cover, cfg.CachePath+"/cover/"+strconv.Itoa(v.Cid)+".jpg")
			if err != nil {
				runtime.LogErrorf(ctx, "保存封面时发生错误：%s\n", err)
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
