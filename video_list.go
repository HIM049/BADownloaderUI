package main

import (
	"errors"
	"strconv"

	"github.com/myuser/bilibili"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 读取视频列表
func (a *App) GetVideoList() []VideoInformationList {
	cfg := new(Config)
	cfg.Get()

	var list []VideoInformationList
	err := LoadJsonFile(cfg.VideoListPath, &list)
	if err != nil {
		runtime.LogErrorf(a.ctx, "读取视频列表时发生错误：%s", err)
	}

	return list
}

// 保存视频列表
func (a *App) SaveVideoList(data any) error {
	cfg := new(Config)
	cfg.Get()

	err := SaveJsonFile(cfg.VideoListPath, data)
	if err != nil {
		return err
	}
	return nil
}

// 创建并保存视频列表
func (a *App) MakeAndSaveList(favlistID string, downloadCount int, downloadCompilation bool) error {
	cfg := new(Config)
	cfg.Get()
	data, err := makeVideoList(a, favlistID, downloadCount, downloadCompilation)
	if err != nil {
		return err
	}
	err = SaveJsonFile(cfg.VideoListPath, data)
	if err != nil {
		return err
	}
	return nil
}

// 创建并保存视频列表（视频合集）
func (a *App) MakeAndSaveCompList(mid, sid, downloadCount int, downloadCompilation bool) error {
	cfg := new(Config)
	cfg.Get()
	data, err := makeVideoListFromComp(a, mid, sid, downloadCount, downloadCompilation)
	if err != nil {
		return err
	}
	err = SaveJsonFile(cfg.VideoListPath, data)
	if err != nil {
		return err
	}
	return nil
}

// 视频数据结构
type VideoInformationList struct {
	Bvid      string `json:"bvid"`
	Cid       int    `json:"cid"`
	Title     string `json:"title"`
	PageTitle string `json:"page_title"`
	Format    string `json:"format"`
	Audio     AudioInformation
	Meta      MetaInformation
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

// 获取视频流
// TODO：请求前检查数据
func (v *VideoInformationList) GetStream(sessdata string) error {
	// 请求信息
	json, err := bilibili.GetVideoStream(v.Bvid, strconv.Itoa(v.Cid), sessdata)
	if err != nil {
		return err
	}
	// 错误检查
	if CheckObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}
	v.Audio.Audio.Quality = int(gjson.Get(json, "data.dash.audio.0.id").Int())
	v.Audio.Audio.Stream = gjson.Get(json, "data.dash.audio.0.base_url").String()
	v.Audio.Flac.Quality = int(gjson.Get(json, "data.dash.flac.id").Int())
	v.Audio.Flac.Stream = gjson.Get(json, "data.dash.flac.base_url").String()

	return nil
}

type Video struct {
	Bvid string `json:"bvid"`
	Meta struct {
		Title      string `json:"title"`       // 视频标题
		Cover      string `json:"cover"`       // 封面
		Author     string `json:"author"`      // 作者
		LyricsPath string `json:"lyrics_path"` // 歌词
	}
	Up struct {
		Mid    int    `json:"mid"`    // UP MID
		Name   string `json:"name"`   // UP 昵称
		Avatar string `json:"avatar"` // UP 头像
	}
	Videos []Videos
}
type Videos struct {
	Cid  int    `json:"cid"`
	Part string `json:"part"` // 分集名称
	Meta struct {
		SongName string `json:"song_name"` // 歌名
	}
	Stream struct {
		Audio struct {
			Id      int    `json:"id"`       // 音质代码
			BaseUrl string `json:"base_url"` // 音频流
		}
		Flac struct {
			Id      int    `json:"id"`       // 音质代码
			BaseUrl string `json:"base_url"` // 音频流
		}
	}
}

// 以 BVID 为单位请求视频详细信息
func (v *Video) BvQuery(sessdata string, downloadCompilation bool) error {
	json, err := bilibili.GetVideoPageInformation(v.Bvid, sessdata)
	if err != nil {
		return err
	}

	// 错误检查
	if CheckObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}

	// 将信息写入结构体
	v.Meta.Title = gjson.Get(json, "data.title").String()     // 视频标题
	v.Meta.Cover = gjson.Get(json, "data.pic").String()       // 视频封面
	v.Up.Mid = int(gjson.Get(json, "data.owner.mid").Int())   // UP MID
	v.Up.Name = gjson.Get(json, "data.owner.name").String()   // UP 昵称
	v.Up.Avatar = gjson.Get(json, "data.owner.face").String() // UP 头像

	// 分 p 总数
	total := int(gjson.Get(json, "data.videos").Int())
	if !downloadCompilation {
		total = 1
	}
	// 根据分 P 数量写入对应信息
	for i := 0; i < total; i++ {

		// 单个分集视频信息
		videos := Videos{
			Cid:  int(gjson.Get(json, "data.pages."+strconv.Itoa(i)+".cid").Int()),
			Part: gjson.Get(json, "data.pages."+strconv.Itoa(i)+".part").String(),
		}
		v.Videos = append(v.Videos, videos)
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

// 创建视频任务列表
func makeVideoList(a *App, favlistId string, downloadCount int, downloadCompilation bool) (*[]VideoInformationList, error) {
	var videoList []VideoInformationList
	cfg := new(Config)
	cfg.Get()
	sessdata := ""
	if cfg.Account.UseAccount && cfg.Account.IsLogin {
		sessdata = cfg.Account.SESSDATA
	}
	// 请求收藏夹基础数据，初始化循环
	favlist, err := bilibili.GetFavListObj(favlistId, sessdata, 1, 1)
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
		favlist, err := bilibili.GetFavListObj(favlistId, sessdata, 20, i+1)
		if err != nil {
			return nil, err
		}
		// 遍历分页
		for j := 0; j < pageSize; j++ {
			// 获取当前视频详细信息
			video := new(Video)
			video.Bvid = favlist.Data.Medias[j].Bvid
			err := video.BvQuery("", downloadCompilation)
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

// 创建视频任务列表（合集）
func makeVideoListFromComp(a *App, mid, sid, downloadCount int, downloadCompilation bool) (*[]VideoInformationList, error) {
	var videoList []VideoInformationList

	// 请求收藏夹基础数据，初始化循环
	favlist, err := bilibili.GetCompliationObj(mid, sid, 1, 1)
	if err != nil {
		return nil, err
	}
	// 计算下载页数
	var pageCount int
	if downloadCount == 0 {
		// 如果下载数量为 0 （全部下载）
		downloadCount = favlist.Data.Meta.Total
		pageCount = downloadCount / 20
	} else {
		// 计算下载页数
		pageCount = downloadCount / 20
	}
	// 非完整页面
	if downloadCount%20 != 0 {
		pageCount++
	}

	runtime.LogInfof(a.ctx, "即将下载%d个视频，%d页", downloadCount, pageCount)
	// 主循环
	for i := 0; i < pageCount; i++ {
		pageSize := 20

		// 处理非完整尾页
		if i+1 == pageCount && downloadCount%20 != 0 {
			pageSize = downloadCount % 20
		}

		// 获取当前分页信息
		favlist, err := bilibili.GetCompliationObj(mid, sid, 20, i+1)
		if err != nil {
			return nil, err
		}
		// 遍历分页
		for j := 0; j < pageSize; j++ {
			// 获取当前视频详细信息
			video := new(Video)
			video.Bvid = favlist.Data.Archives[j].Bvid
			err := video.BvQuery("", downloadCompilation)
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
