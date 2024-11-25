package main

import (
	"bytes"
	"html/template"
	"os"
	"path"
	"strconv"

	tag "github.com/gcottom/audiometa"
)

// 修改 TAG
func ChangeTag(cfg *Config, opt *DownloadOption, v *VideoInformation) error {

	// 准备参数
	file := cfg.FileConfig.CachePath + "/music/" + strconv.Itoa(v.Cid) + v.Format
	songCover := cfg.FileConfig.CachePath + "/cover/" + strconv.Itoa(v.Cid) + ".jpg"
	songName := v.Meta.SongName
	songAuthor := v.Meta.Author

	// 打开歌曲元数据
	tags, err := tag.OpenTag(file)
	if err != nil {
		return err
	}

	// 封面
	if opt.SongCover {
		err := tags.SetAlbumArtFromFilePath(songCover)
		if err != nil {
			return err
		}
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

type FileName struct {
	Title    string
	Subtitle string
	Quality  string
	ID       int
	Format   string
	SongName string
	Author   string
}

// 输出文件
func OutputFile(cfg *Config, v *VideoInformation, fileName *FileName) error {
	// 处理模板和生成文件名
	tmpl, err := template.New("filename").Parse(cfg.FileConfig.FileNameTemplate)
	if err != nil {
		return err
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, fileName)
	if err != nil {
		return err
	}

	// 添加路径
	sourcePath := path.Join(cfg.FileConfig.CachePath, "music", strconv.Itoa(v.Cid)+v.Format)
	destPath := path.Join(cfg.FileConfig.DownloadPath, output.String())

	// 重命名歌曲文件并移动位置
	err = os.Rename(sourcePath, destPath)
	if err != nil {
		return err
	}
	return nil
}

// 修改 TAG
func SingleChangeTag(cfg *Config, opt *DownloadOption, auid, songName, songAuthor string) error {

	// 准备参数
	file := cfg.FileConfig.CachePath + "/single/music/" + auid + AudioType.m4a
	songCover := cfg.FileConfig.CachePath + "/single/cover/" + auid + ".jpg"

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
func SingleOutputFile(cfg *Config, uuid, Title string) error {

	sourcePath := path.Join(cfg.FileConfig.CachePath, "single/music", uuid+AudioType.m4a)
	destPath := path.Join(cfg.FileConfig.DownloadPath, Title+AudioType.mp3)

	// 重命名歌曲文件并移动位置
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		return err
	}
	return nil
}
