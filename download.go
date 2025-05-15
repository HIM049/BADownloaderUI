package main

import (
	"bili-audio-downloader/bilibili"
	"bili-audio-downloader/config"
	"bili-audio-downloader/constants"
	"bili-audio-downloader/services"
	"path"
	"strconv"
	"sync"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) ListDownload(listPath string, opt DownloadOption) error {
	// 初始化参数
	sessdata := ""
	if config.Cfg.Account.UseAccount && config.Cfg.Account.IsLogin {
		sessdata = config.Cfg.Account.SESSDATA
	}

	sem := make(chan struct{}, config.Cfg.DownloadConfig.DownloadThreads+1)
	var wg sync.WaitGroup

	videoList := new(services.VideoList)
	err := videoList.Get(listPath)
	if err != nil {
		wails.LogErrorf(a.ctx, "读取视频列表时发生错误：%s", err)
		return err
	}

	// 格式判断
	audioType := constants.AudioType.M4a
	if config.Cfg.FileConfig.ConvertFormat {
		audioType = constants.AudioType.Mp3
	}

	// 遍历下载队列
	for i, video := range videoList.List {
		// 并发函数
		go func(v services.VideoInformation, num int) {
			sem <- struct{}{} // 给通道中填入数据
			wg.Add(1)         // 任务 +1
			// 下载完成后
			defer func() {
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知

				wails.EventsEmit(a.ctx, "downloadFinish", v.Meta.SongName)
			}()

			// 处理文件名结构体
			fileName := new(FileName)
			fileName.Title = v.Title
			fileName.Subtitle = v.PageTitle
			fileName.ID = num
			fileName.Quality = "hires"

			//判断是否已下载
			finalFile := path.Join(config.Cfg.GetDownloadPath(), v.Title+audioType)
			if IsFileExists(finalFile) {
				wails.LogInfof(a.ctx, "跳过已存在的视频: %s", finalFile)
				return
			}

			musicPathAndName := config.Cfg.GetCachePath() + "/music/" + strconv.Itoa(v.Cid)

			// 下载视频
			for i := 0; i < config.Cfg.DownloadConfig.RetryCount; i++ {

				// 音频下载逻辑
				if v.IsAudio {
					audio := new(bilibili.Audio)
					audio.Auid = v.Bvid
					err := audio.GetStream(sessdata)
					if err != nil {
						wails.LogErrorf(a.ctx, "(队列%d) 获取媒体流时出现错误：%s  (重试 %d )", num, err, i+1)
						continue
					}

					// 下载媒体流
					err = bilibili.StreamingDownloader(audio.Stream.StreamLink, musicPathAndName+constants.AudioType.M4a)
					if err != nil {
						// 下载失败
						wails.LogErrorf(a.ctx, "(队列%d) 下载时出现错误：%s  (重试 %d )", num, err, i+1)
						continue
					} else {
						wails.LogInfof(a.ctx, "(队列%d) 下载视频成功", num)
					}
					break
				}

				err := v.GetStream(sessdata)
				if err != nil {
					// 获取流失败
					wails.LogErrorf(a.ctx, "(队列%d) 获取媒体流时出现错误：%s  (重试 %d )", num, err, i+1)
					continue
				}
				// 下载媒体流
				err = bilibili.StreamingDownloader(v.Audio.Stream, musicPathAndName+v.Format)
				if err != nil {
					// 下载失败
					wails.LogErrorf(a.ctx, "(队列%d) 下载时出现错误：%s  (重试 %d )", num, err, i+1)
					continue
				} else {
					wails.LogInfof(a.ctx, "(队列%d) 下载视频成功", num)
				}

				break
			}

			// 判断文件类型并转码
			if v.Format == constants.AudioType.M4a && config.Cfg.FileConfig.ConvertFormat {
				wails.LogInfof(a.ctx, "(队列%d) 转码为 MP3", num)
				v.Format = constants.AudioType.Mp3
				fileName.Format = constants.AudioType.Mp3

				// 转码文件
				err = ConventFile(musicPathAndName+constants.AudioType.M4a, musicPathAndName+constants.AudioType.Mp3)
				if err != nil {
					wails.LogErrorf(a.ctx, "转码文件时发生错误：%s", err)
				} else {
					wails.LogInfof(a.ctx, "(队列%d) 转码文件成功", num)
				}
			} else {
				wails.LogInfof(a.ctx, "(队列%d) 无需转码", num)
				fileName.Format = v.Format
			}

			// 写入元数据
			if v.Format != constants.AudioType.Flac {
				fileName.Quality = "normal"
				err = ChangeTag(&config.Cfg, &opt, &v)
				if err != nil {
					wails.LogErrorf(a.ctx, "(队列%d) 写入元数据时发生错误：%s", num, err)
				} else {
					wails.LogInfof(a.ctx, "(队列%d) 写入元数据成功", num)
				}
			}

			// 输出文件
			err = OutputFile(&config.Cfg, &v, *fileName)
			if err != nil {
				wails.LogErrorf(a.ctx, "输出文件时发生错误：%s", err)
			} else {
				wails.LogInfof(a.ctx, "(队列%d) 输出文件成功", num)
			}

		}(video, i)

		go func(v services.VideoInformation, num int) {
			// 下载封面图片
			err = bilibili.SaveFromURL(v.Meta.Cover, config.Cfg.GetCachePath()+"/cover/"+strconv.Itoa(v.Cid)+".jpg")
			if err != nil {
				wails.LogErrorf(a.ctx, "保存封面时发生错误：%s", err)
			} else {
				wails.LogInfof(a.ctx, "(队列%d) 下载封面成功", num)
			}
		}(video, i)
		time.Sleep(10 * time.Millisecond)
	}
	// 等待任务执行完成
	wg.Wait()

	// 下载完成后保存列表
	err = videoList.Save(listPath)
	if err != nil {
		return err
	}

	return nil
}
