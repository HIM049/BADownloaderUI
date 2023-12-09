package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) StartDownload(opt DownloadOption) {
	// 初始化参数
	cfg := GetConfig(a.ctx)
	audioType := ".m4a"

	sem := make(chan struct{}, cfg.DownloadThreads+1)
	var wg sync.WaitGroup

	// 获取任务队列
	var list []VideoInformationList
	err := LoadJsonFile(cfg.VideoListPath, &list)
	if err != nil {
		runtime.LogErrorf(a.ctx, "读取视频列表时发生错误：%s", err)
	}

	// 遍历下载队列
	for i, video := range list {
		// 并发函数
		go func(v VideoInformationList, num int) {
			sem <- struct{}{} // 给通道中填入数据
			wg.Add(1)         // 任务 +1

			runtime.LogInfof(a.ctx, "开始下载视频%d", num)
			// 下载视频
			for i := 0; i < cfg.RetryCount; i++ {
				err := GetAndDownload(v.Bvid, v.Cid, cfg.CachePath+"/music/"+strconv.Itoa(v.Cid)+".m4a")
				if err != nil {
					// 下载失败
					runtime.LogErrorf(a.ctx, "(视频%d) 下载视频时出现错误：%s  (重试 %d )", num, err, i+1)
					continue
				}
			}
			runtime.LogInfof(a.ctx, "(视频%d) 下载视频成功", num)

			// 下载封面图片
			err = SaveFromURL(v.Cover, cfg.CachePath+"/cover/"+strconv.Itoa(v.Cid)+".jpg")
			if err != nil {
				runtime.LogErrorf(a.ctx, "保存封面时发生错误：%s", err)
			}
			runtime.LogInfof(a.ctx, "(视频%d) 下载封面成功", num)

			// 写入元数据
			err = ChangeTag(&cfg, &opt, &v, audioType)
			if err != nil {
				runtime.LogErrorf(a.ctx, "(视频%d) 写入元数据时发生错误：%s", num, err)
			}
			runtime.LogInfof(a.ctx, "(视频%d) 写入元数据成功", num)

			// 输出文件
			err = OutputFile(&cfg, &v, audioType)
			if err != nil {
				runtime.LogErrorf(a.ctx, "输出文件时发生错误：%s", err)
			}
			runtime.LogInfof(a.ctx, "(视频%d) 输出文件成功", num)

			// 下载完成后
			defer func() {
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知
			}()

		}(video, i)
		time.Sleep(10 * time.Millisecond)
	}
	// 等待任务执行完成
	wg.Wait()
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
