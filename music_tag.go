package main

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	tag "github.com/gcottom/audiometa"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 写入元数据
func ConcurrentChangeTag(ctx context.Context, cfg *Config, opt *DownloadOption, audioType string) error {
	sem := make(chan struct{}, cfg.DownloadThreads)
	var wg sync.WaitGroup

	// 获取任务队列
	var list []VideoInformationList

	err := LoadJsonFile(cfg.VideoListPath, &list)
	if err != nil {
		return err
	}
	// 遍历下载队列
	for _, video := range list {
		// 转码
		go func(v VideoInformationList) {
			sem <- struct{}{} // 限制并发量
			wg.Add(1)         // 任务 +1

			// 准备参数
			file := cfg.CachePath + "/music/" + strconv.Itoa(v.Cid) + audioType
			songCover := ""
			songName := ""
			songAuthor := ""

			// 元数据个性化写入
			if opt.SongCover {
				songCover = cfg.CachePath + "/cover/" + strconv.Itoa(v.Cid) + ".jpg"
			}
			if opt.SongName {
				songName = v.SongName
			}
			if opt.SongAuthor {
				songAuthor = v.Author
			}

			// 写入歌曲元数据
			err = ChangeTag(file, songName, songCover, songAuthor)
			if err != nil {
				runtime.LogErrorf(ctx, "写入元数据时发生错误：%s\n", err)
			}
			// 下载完成后
			defer func() {
				runtime.LogInfo(ctx, "元数据添加完成")
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知
			}()
		}(video)
	}
	// 等待任务执行完成
	wg.Wait()
	time.Sleep(1 * time.Second)
	return nil
}

// 重命名和移动文件
func ConcurrentChangeName(ctx context.Context, threads int, videolistPath, audioType, aideoSourcePath, aideoDestPath string) error {
	sem := make(chan struct{}, threads)
	var wg sync.WaitGroup

	// 获取任务队列
	var list []VideoInformationList
	err := LoadJsonFile(videolistPath, &list)
	if err != nil {
		return err
	}
	// 遍历下载队列
	for _, video := range list {
		// 转码
		go func(v VideoInformationList) {
			sem <- struct{}{} // 限制并发量
			wg.Add(1)         // 任务 +1

			// 处理音频标题
			NfileName := v.Title
			// 如果是分 P （以分 P 命名为主）
			if v.IsPage {
				NfileName = v.Title + "(" + v.PageTitle + ")"
			}
			// 重命名歌曲文件并移动位置
			err = RenameAndMoveFile(aideoSourcePath+strconv.Itoa(v.Cid)+audioType, aideoDestPath+NfileName+audioType)
			if err != nil {
				return
			}

			// 下载完成后
			defer func() {
				runtime.LogInfo(ctx, "输出完成")
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知
			}()
		}(video)
	}
	// 等待任务执行完成
	wg.Wait()
	time.Sleep(1 * time.Second)
	return nil
}

// 修改 TAG
func ChangeTag(file, songName, coverPath, artist string) error {
	tags, err := tag.OpenTag(file)
	if err != nil {
		return err
	}
	tags.SetTitle(songName) // 歌曲名
	tags.SetArtist(artist)  // 艺术家
	// 写入封面
	if coverPath != "" {
		tags.SetAlbumArtFromFilePath(coverPath)
	}

	// TODO: 将歌曲 tag 数据整理为结构体
	// TODO: 修改作词人，作曲人等，以及自动适配

	// 保存更改
	err = tags.Save()
	if err != nil {
		return err
	}

	return nil
}

// 重命名和移动
func RenameAndMoveFile(sourcePath, destPath string) error {
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		return err
	}
	return nil
}
