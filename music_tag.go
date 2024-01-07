package main

import (
	"os"
	"strconv"

	tag "github.com/gcottom/audiometa"
)

// 修改 TAG
func ChangeTag(cfg *Config, opt *DownloadOption, v *VideoInformationList, audioType string) error {

	// 准备参数
	file := cfg.CachePath + "/music/" + strconv.Itoa(v.Cid) + audioType
	songCover := cfg.CachePath + "/cover/" + strconv.Itoa(v.Cid) + ".jpg"
	songName := v.Meta.SongName
	songAuthor := v.Meta.Author

	// 打开歌曲元数据
	tags, err := tag.OpenTag(file)
	if err != nil {
		return err
	}

	// 封面
	if opt.SongCover {
		tags.SetAlbumArtFromFilePath(songCover)
	}
	// 歌曲名
	if opt.SongName {
		tags.SetTitle(songName)
	}
	// 艺术家
	if opt.SongAuthor {
		tags.SetArtist(songAuthor)
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

// 输出文件
func OutputFile(cfg *Config, v *VideoInformationList, audioType string) error {
	// 处理音频标题
	NfileName := v.Title
	// 如果是分 P （以分 P 命名为主）
	if v.IsPage {
		NfileName = v.Title + "(" + v.PageTitle + ")"
	}

	sourcePath := cfg.CachePath + "/music/" + strconv.Itoa(v.Cid) + audioType
	destPath := cfg.DownloadPath + "/" + NfileName + audioType
	// 重命名歌曲文件并移动位置
	err := RenameAndMoveFile(sourcePath, destPath)
	if err != nil {
		return err
	}
	return nil
}

// 修改 TAG
func SingleChangeTag(cfg *Config, opt *DownloadOption, auid, songName, songAuthor, audioType string) error {

	// 准备参数
	file := cfg.CachePath + "/single/music/" + auid + audioType
	songCover := cfg.CachePath + "/single/cover/" + auid + ".jpg"

	// 打开歌曲元数据
	tags, err := tag.OpenTag(file)
	if err != nil {
		return err
	}

	// 封面
	if opt.SongCover {
		tags.SetAlbumArtFromFilePath(songCover)
	}
	// 歌曲名
	if opt.SongName {
		tags.SetTitle(songName)
	}
	// 艺术家
	if opt.SongAuthor {
		tags.SetArtist(songAuthor)
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

// 输出文件
func SingleOutputFile(cfg *Config, uuid, Title, audioType string) error {

	sourcePath := cfg.CachePath + "/single/music/" + uuid + audioType
	destPath := cfg.DownloadPath + "/" + Title + audioType
	// 重命名歌曲文件并移动位置
	err := RenameAndMoveFile(sourcePath, destPath)
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
