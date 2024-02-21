package main

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetVideoList() []VideoInformationList {
	cfg := GetConfig(a.ctx)

	var list []VideoInformationList
	err := LoadJsonFile(cfg.VideoListPath, &list)
	if err != nil {
		runtime.LogErrorf(a.ctx, "读取视频列表时发生错误：%s", err)
	}

	return list
}

func (a *App) SaveVideoList(data any) error {
	cfg := GetConfig(a.ctx)

	err := SaveJsonFile(cfg.VideoListPath, data)
	if err != nil {
		return err
	}
	return nil
}

// 创建并保存视频列表
func (a *App) MakeAndSaveList(favlistID string, downloadCount int, downloadCompilation bool) error {
	cfg := GetConfig(a.ctx)
	data, err := makeVideoList(a, downloadCount, downloadCompilation)
	if err != nil {
		return err
	}
	err = SaveJsonFile(cfg.VideoListPath, data)
	if err != nil {
		return err
	}
	return nil
}

// 转换视频信息为列表信息
func (video *Video) TransToVideoInfList() *[]VideoInformationList {
	var VideoList []VideoInformationList
	for _, part := range video.Videos {
		var list VideoInformationList
		list.Bvid = video.Bvid
		list.Cid = part.Cid
		list.Title = CheckFileName(video.Meta.Title)
		list.PageTitle = CheckFileName(part.Part)
		list.Format = AudioType.m4a
		// 音频流
		list.Audio.Audio.Quality = part.Stream.Audio.Id
		list.Audio.Audio.Stream = part.Stream.Audio.BaseUrl
		list.Audio.Flac.Quality = part.Stream.Flac.Id
		list.Audio.Flac.Stream = part.Stream.Flac.BaseUrl
		// 元数据
		list.Meta.Cover = video.Meta.Cover
		list.Meta.Author = video.Up.Name
		// list.Meta.Lyrics_path =

		// 处理音频标题（单 P 视频）
		var SongName string
		SongName, err := ExtractTitle(list.PageTitle)
		if err != nil {
			// 如果无法判断标题
			SongName = list.Title
		}
		list.Meta.SongName = SongName
		VideoList = append(VideoList, list)

	}
	return &VideoList
}

// 视频数据结构
type VideoInformationList struct {
	Bvid      string `json:"bvid"`
	Cid       int    `json:"cid"`
	Title     string `json:"title"`
	PageTitle string `json:"page_title"`
	// Videos    int    `json:"videos"`
	// ListID    int    `json:"list_id"`
	// IsPage    bool   `json:"is_page"`
	// PageID    int    `json:"page_id"`
	Format string `json:"format"`
	Audio  AudioInformation
	Meta   MetaInformation
}

type AudioInformation struct {
	Audio struct {
		Quality int    `json:"quality"`
		Stream  string `json:"stream"`
	}
	Flac struct {
		Quality int    `json:"quality"`
		Stream  string `json:"stream"`
	}
}

type MetaInformation struct {
	SongName    string `json:"song_name"`
	Cover       string `json:"cover"`
	Author      string `json:"author"`
	Lyrics_path string `json:"lyrics_path"`
}

// 创建视频任务列表
func makeVideoList(a *App, downloadCount int, downloadCompilation bool) (*[]VideoInformationList, error) {
	var videoList []VideoInformationList

	// 请求收藏夹基础数据，初始化循环
	favlist, err := GetFavListObj(1, 1)
	if err != nil {
		return nil, err
	}
	// 计算下载页数
	var pageCount int
	if downloadCount == 0 {
		// 如果下载数量为 0 （全部下载）
		downloadCount = favlist.Data.Info.Media_count
		pageCount = downloadCount / 20
	} else {
		// 计算下载页数
		pageCount = downloadCount / 20
	}
	// 非完整页面
	if downloadCount%20 != 0 {
		pageCount++
	}

	// 主循环
	for i := 0; i < pageCount; i++ {
		pageSize := 20

		// 处理非完整尾页
		if i+1 == pageCount && downloadCount%20 != 0 {
			pageSize = downloadCount % 20
		}

		// 获取当前分页信息
		favlist, err := GetFavListObj(20, i+1)
		if err != nil {
			return nil, err
		}
		// 遍历分页
		for j := 0; j < pageSize; j++ {
			// 获取当前视频详细信息
			video := new(Video)
			video.Bvid = favlist.Data.Medias[j].Bvid
			err := video.BvQuery("")
			if err != nil {
				// 视频失效
				runtime.LogErrorf(a.ctx, "获取 "+video.Bvid+" 信息时发生错误: %s", err)
				continue
			}

			// 组合数据
			videoList = append(videoList, *video.TransToVideoInfList()...)
		}
	}

	return &videoList, nil
}
