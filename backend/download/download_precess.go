package download

import (
	"bili-audio-downloader/backend/config"
	"context"
	"sync"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

func DownloadTaskList(ctx context.Context) {

	sem := make(chan struct{}, config.Cfg.DownloadConfig.DownloadThreads+1)
	var wg sync.WaitGroup

	wails.LogPrintf(ctx, "start list download process")
	// 遍历下载队列
	for i, task := range DownloadList {
		go func() {
			// 并发函数
			sem <- struct{}{} // 给通道中填入数据
			wg.Add(1)         // 任务 +1
			// 下载完成后
			defer func() {
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知

				wails.EventsEmit(ctx, "downloadFinish", i)
			}()

			task.SetID(i)

			// // 判断是否已下载
			// finalFile := path.Join(config.Cfg.GetDownloadPath(), v.Title+audioType)
			// if utils.IsFileExists(finalFile) {
			// 	wails.LogInfof(ctx, "跳过已存在的视频: %s", finalFile)
			// 	return
			// }

			// 下载媒体流和封面
			wails.LogPrintf(ctx, "Downloading file: %d", i)
			err := task.Download()
			if err != nil {
				wails.LogErrorf(ctx, "Failed to download (file%d): %v", i, err)
				return
			}

			// 转码文件
			wails.LogPrintf(ctx, "Converting format: %d", i)
			if config.Cfg.FileConfig.ConvertFormat {
				err = task.ConventFormat()
				if err != nil {
					wails.LogErrorf(ctx, "Failed to convert file (file%d): %v", i, err)
					return
				}
			}

			wails.LogPrintf(ctx, "Writing metadata: %d", i)
			err = task.WriteMetadata()
			if err != nil {
				wails.LogErrorf(ctx, "Failed to write metadata (file%d): %v", i, err)
				return
			}

			wails.LogPrintf(ctx, "Finishing download task: %d", i)
			err = task.ExportFile()
			if err != nil {
				wails.LogErrorf(ctx, "Failed to export file (file%d): %v", i, err)
				return
			}
		}()

		time.Sleep(10 * time.Millisecond)
	}
	// 等待任务执行完成
	wg.Wait()

	//// 下载完成后保存列表
	//err = videoList.Save(listPath)
	//if err != nil {
	//	return err
	//}

}
