package main

import (
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var AudioType = struct {
	m4a string
	mp3 string
}{m4a: ".m4a", mp3: ".mp3"}

func (a *App) StartDownload(opt DownloadOption) {
	// 初始化参数
	cfg := GetConfig(a.ctx)
	_ = os.MkdirAll(path.Join(cfg.DownloadPath, FavListID), 0755)

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
			// 下载完成后
			defer func() {
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知

				runtime.EventsEmit(a.ctx, "downloadFinish")
			}()

			// 处理音频标题
			finalfileName := v.Title
			// 如果是分 P （以分 P 命名为主）
			if v.IsPage {
				finalfileName += "(" + v.PageTitle + ")"
			}

			//判断是否已下载
			finalFile := path.Join(cfg.DownloadPath, FavListID, finalfileName+AudioType.mp3)
			if IsFileExists(finalFile) {
				runtime.LogInfof(a.ctx, "跳过已下载: %s", finalFile)
				return
			}

			runtime.LogDebugf(a.ctx, "开始下载视频%d", num)
			musicPathAndName := cfg.CachePath + "/music/" + strconv.Itoa(v.Cid)

			// 下载视频
			for i := 0; i < cfg.RetryCount; i++ {
				err := GetAndDownload(v.Bvid, v.Cid, musicPathAndName+AudioType.m4a)
				if err != nil {
					// 下载失败
					runtime.LogErrorf(a.ctx, "(视频%d) 下载视频时出现错误：%s  (重试 %d )", num, err, i+1)
					continue
				}
			}
			runtime.LogDebugf(a.ctx, "(视频%d) 下载视频成功", num)

			if v.Format == AudioType.m4a {
				runtime.LogDebug(a.ctx, "是 M4A 文件")

				// 转码文件
				err = ConventFile(musicPathAndName+AudioType.m4a, musicPathAndName+AudioType.mp3)
				if err != nil {
					runtime.LogErrorf(a.ctx, "转码文件时发生错误：%s", err)
				}
				runtime.LogDebugf(a.ctx, "(视频%d) 转码文件成功", num)
			}

			// 写入元数据
			err = ChangeTag(&cfg, &opt, &v)
			if err != nil {
				runtime.LogErrorf(a.ctx, "(视频%d) 写入元数据时发生错误：%s", num, err)
			}
			runtime.LogDebugf(a.ctx, "(视频%d) 写入元数据成功", num)

			// 输出文件
			err = OutputFile(&cfg, &v, finalfileName+AudioType.mp3)
			if err != nil {
				runtime.LogErrorf(a.ctx, "输出文件时发生错误：%s", err)
			}
			runtime.LogDebugf(a.ctx, "(视频%d) 输出文件成功", num)

		}(video, i)

		go func(v VideoInformationList, num int) {
			// 下载封面图片
			err = SaveFromURL(v.Meta.Cover, cfg.CachePath+"/cover/"+strconv.Itoa(v.Cid)+".jpg")
			if err != nil {
				runtime.LogErrorf(a.ctx, "保存封面时发生错误：%s", err)
			}
			runtime.LogDebugf(a.ctx, "(视频%d) 下载封面成功", num)
		}(video, i)
		time.Sleep(10 * time.Millisecond)
	}
	// 等待任务执行完成
	wg.Wait()
}

func (a *App) AudioDownload(opt DownloadOption, auid, songName, songAuthor, title string) {
	cfg := GetConfig(a.ctx)

	obj, err := GetAudioObj(auid, "2")
	if err != nil {
		runtime.LogErrorf(a.ctx, "获取音频流时发生错误：%s", err)
	}

	// 下载封面图片
	err = SaveFromURL(obj.Data.Cover, cfg.CachePath+"/single/cover/"+auid+".jpg")
	if err != nil {
		runtime.LogErrorf(a.ctx, "保存封面时发生错误：%s", err)
	}
	// 下载媒体流
	err = StreamingDownloader(obj.Data.Cdns[0], cfg.CachePath+"/single/music/"+auid+AudioType.m4a)
	if err != nil {
		runtime.LogErrorf(a.ctx, "保存媒体流时发生错误：%s", err)
	}
	// 写入元数据
	err = SingleChangeTag(&cfg, &opt, auid, songName, songAuthor)
	if err != nil {
		runtime.LogErrorf(a.ctx, "写入元数据时发生错误：%s", err)
	}
	// 输出文件
	err = SingleOutputFile(&cfg, auid, title)
	if err != nil {
		runtime.LogErrorf(a.ctx, "输出文件时发生错误：%s", err)
	}
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
