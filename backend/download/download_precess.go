package download

import (
	"bili-audio-downloader/backend/config"
	"context"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
	"sync"
	"time"
)

func DownloadTaskList(ctx context.Context) {

	// 初始化参数
	//sessdata := ""
	//if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
	//	sessdata = config.Cfg.Account.SESSDATA
	//}

	sem := make(chan struct{}, config.Cfg.DownloadConfig.DownloadThreads+1)
	var wg sync.WaitGroup

	//// 格式判断
	//audioType := constants.AudioType.M4a
	//if config.Cfg.FileConfig.ConvertFormat {
	//	audioType = constants.AudioType.Mp3
	//}

	wails.LogPrintf(ctx, "start download task list")
	// 遍历下载队列
	for i, task := range DownloadList {
		// 并发函数
		sem <- struct{}{} // 给通道中填入数据
		wg.Add(1)         // 任务 +1
		// 下载完成后
		defer func() {
			<-sem     // 释放一个并发槽
			wg.Done() // 发出任务完成通知

			wails.EventsEmit(ctx, "downloadFinish", i)
		}()

		//// 处理文件名结构体
		//fileName := new(services.FileName)
		//fileName.Title = t.Title
		//fileName.Subtitle = v.PageTitle
		//fileName.ID = num
		//fileName.Quality = "hires"

		////判断是否已下载
		//finalFile := path.Join(config.Cfg.GetDownloadPath(), v.Title+audioType)
		//if utils.IsFileExists(finalFile) {
		//	wails.LogInfof(ctx, "跳过已存在的视频: %s", finalFile)
		//	return
		//}

		//musicPathAndName := config.Cfg.GetCachePath() + "/music/" + strconv.Itoa(v.Cid)

		// 下载媒体流和封面
		wails.LogPrintf(ctx, "开始下载: %d", i)
		err := task.Download()
		if err != nil {
			wails.LogErrorf(ctx, "下载失败：%d ERROR: %v", i, err)
			return
		}

		//// 判断文件类型并转码
		//if v.Format == constants.AudioType.M4a && config.Cfg.FileConfig.ConvertFormat {
		//	wails.LogInfof(a.ctx, "(队列%d) 转码为 MP3", num)
		//	v.Format = constants.AudioType.Mp3
		//	fileName.Format = constants.AudioType.Mp3
		//
		//	// 转码文件
		//	err = ConventFile(musicPathAndName+constants.AudioType.M4a, musicPathAndName+constants.AudioType.Mp3)
		//	if err != nil {
		//		wails.LogErrorf(a.ctx, "转码文件时发生错误：%s", err)
		//	} else {
		//		wails.LogInfof(a.ctx, "(队列%d) 转码文件成功", num)
		//	}
		//} else {
		//	wails.LogInfof(a.ctx, "(队列%d) 无需转码", num)
		//	fileName.Format = v.Format
		//}

		//// 写入元数据
		//if v.Format != constants.AudioType.Flac {
		//	fileName.Quality = "normal"
		//	err = services.ChangeTag(&config.Cfg, &opt, &v)
		//	if err != nil {
		//		wails.LogErrorf(a.ctx, "(队列%d) 写入元数据时发生错误：%s", num, err)
		//	} else {
		//		wails.LogInfof(a.ctx, "(队列%d) 写入元数据成功", num)
		//	}
		//}

		//// 输出文件
		//err = services.OutputFile(&config.Cfg, &v, *fileName)
		//if err != nil {
		//	wails.LogErrorf(a.ctx, "输出文件时发生错误：%s", err)
		//} else {
		//	wails.LogInfof(a.ctx, "(队列%d) 输出文件成功", num)
		//}

		time.Sleep(10 * time.Millisecond)
	}
	// 等待任务执行完成
	wg.Wait()

	//// 下载完成后保存列表
	//err = videoList.Save(listPath)
	//if err != nil {
	//	return err
	//}

	return
}
