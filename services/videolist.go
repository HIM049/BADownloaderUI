package services

import (
	"bili-audio-downloader/bilibili"
	"bili-audio-downloader/config"
	"bili-audio-downloader/constants"
	"bili-audio-downloader/utils"
	"errors"
	"strconv"

	"github.com/tidwall/gjson"
)

// 视频列表
type VideoList struct {
	Count int `json:"count"`
	List  []VideoInformation
}

// 视频数据结构
type VideoInformation struct {
	Bvid      string `json:"bvid"`
	Cid       int    `json:"cid"`
	Title     string `json:"title"`
	PageTitle string `json:"page_title"`
	Format    string `json:"format"`
	PartId    int    `json:"part_id"`
	IsAudio   bool   `json:"is_audio"`
	Delete    bool   `json:"delete"`
	Audio     AudioInformation
	Meta      MetaInformation
}
type AudioInformation struct {
	Quality int    `json:"quality"`
	Stream  string `json:"stream"`
}
type MetaInformation struct {
	SongName    string `json:"song_name"`
	Cover       string `json:"cover"`
	Author      string `json:"author"`
	Lyrics_path string `json:"lyrics_path"`
}

// 向列表中添加一个项目
func (list *VideoList) Add(video *VideoInformation) {
	list.List = append(list.List, *video)
	list.Count++
}

// 向列表中添加一个视频
func (VideoList *VideoList) AddVideo(sessdata, bvid string, downloadCompilation bool) error {
	// 查询视频信息
	video := new(bilibili.Video)
	err := video.Query(sessdata, bvid)
	if err != nil {
		return err
	}

	// 处理分集数量
	var total int = 1
	if downloadCompilation {
		total = len(video.Videos)
	}

	// 保存信息
	for i := 0; i < total; i++ {
		var list VideoInformation
		list.Bvid = video.Bvid
		list.Cid = video.Videos[i].Cid
		list.Title = utils.CheckFileName(video.Meta.Title)
		list.PageTitle = utils.CheckFileName(video.Videos[i].Part)
		list.Format = constants.AudioType.M4a
		// 元数据
		list.Meta.Cover = video.Meta.Cover
		list.Meta.Author = video.Up.Name
		list.Delete = false
		// list.Meta.Lyrics_path =

		list.IsAudio = false

		// 处理音频标题（单 P 视频）
		var SongName string
		if total <= 1 {
			// 单集使用视频标题
			SongName, err = utils.ExtractTitle(list.Title)
			if err != nil {
				SongName = list.Title
			}
		} else {
			// 多集视频使用分集标题
			SongName, err = utils.ExtractTitle(list.PageTitle)
			if err != nil {
				SongName = list.PageTitle
			}
		}
		list.Meta.SongName = SongName
		VideoList.Add(&list)
	}
	return nil
}

// 向列表中添加一个音频项目
func (VideoList *VideoList) AddAudio(sessdata, auid string) error {
	// 查询视频信息
	audio := new(bilibili.Audio)
	err := audio.Query(auid)
	if err != nil {
		return err
	}

	aucid, err := strconv.Atoi(auid)
	if err != nil {
		return err
	}

	// 保存信息
	var list VideoInformation
	list.Bvid = auid
	list.Cid = aucid
	list.Title = utils.CheckFileName(audio.Meta.Title)
	list.PageTitle = utils.CheckFileName(audio.Meta.Title)
	list.Format = constants.AudioType.M4a
	// 元数据
	list.Meta.Cover = audio.Meta.Cover
	list.Meta.Author = audio.Up.Author
	list.Meta.Lyrics_path = audio.Meta.Lyric
	list.Meta.SongName = audio.Meta.Title

	list.IsAudio = true
	list.Delete = false

	VideoList.Add(&list)
	return nil
}

// 向列表中添加一个收藏夹
func (VideoList *VideoList) AddCollection(sessdata, favlistId string, count int, downloadCompilation bool) error {
	// 请求收藏夹基础数据，初始化循环
	favlist, err := bilibili.GetFavListObj(favlistId, sessdata, 1, 1)
	if err != nil {
		return err
	}
	// 计算下载页数
	var pageCount int
	if count == 0 {
		// 如果下载数量为 0 （全部下载）
		count = favlist.Data.Info.Media_count
		pageCount = count / 20
	} else {
		// 计算下载页数
		pageCount = count / 20
	}
	// 非完整页面
	if count%20 != 0 {
		pageCount++
	}

	// 主循环
	for i := 0; i < pageCount; i++ {
		// 获取当前分页信息
		favlist, err := bilibili.GetFavListObj(favlistId, sessdata, 20, i+1)
		if err != nil {
			return err
		}

		// 遍历分页
		for j := 0; j < len(favlist.Data.Medias); j++ {

			if favlist.Data.Medias[j].Type == 2 {
				// 添加视频到列表
				err := VideoList.AddVideo(sessdata, favlist.Data.Medias[j].Bvid, downloadCompilation)
				if err != nil {
					continue
				}
			} else {
				// 添加收藏夹中的音频
				err := VideoList.AddAudio(sessdata, strconv.Itoa(favlist.Data.Medias[j].Id))
				if err != nil {
					continue
				}
			}
		}
	}

	return nil
}

// 向列表中添加一个视频合集
func (VideoList *VideoList) AddCompilation(sessdata string, mid, sid, count int, downloadCompilation bool) error {
	// 请求收藏夹基础数据，初始化循环
	favlist, err := bilibili.GetCompliationObj(mid, sid, 1, 1)
	if err != nil {
		return err
	}
	// 计算下载页数
	var pageCount int
	if count == 0 {
		// 如果下载数量为 0 （全部下载）
		count = favlist.Data.Meta.Total
		pageCount = count / 20
	} else {
		// 计算下载页数
		pageCount = count / 20
	}
	// 非完整页面
	if count%20 != 0 {
		pageCount++
	}

	// 主循环
	for i := 0; i < pageCount; i++ {
		// 获取当前分页信息
		favlist, err := bilibili.GetCompliationObj(mid, sid, 20, i+1)
		if err != nil {
			return err
		}
		// 遍历分页
		for j := 0; j < len(favlist.Data.Archives); j++ {
			// 添加视频到列表
			err := VideoList.AddVideo(sessdata, favlist.Data.Archives[j].Bvid, downloadCompilation)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

// 向列表中添加个人主页视频
func (VideoList *VideoList) AddProfileVideo(sessdata string, mid, count int, downloadCompilation bool) error {
	respJson, err := bilibili.GetProfileVideo(strconv.Itoa(mid), "1", "1", sessdata)
	if err != nil {
		return err
	}

	// 计算下载页数
	var pageCount int
	if count == 0 {
		// 如果下载数量为 0 （全部下载）
		count = int(gjson.Get(respJson, "data.page.count").Int())
		pageCount = count / 20
	} else {
		// 计算下载页数
		pageCount = count / 20
	}
	// 非完整页面
	if count%20 != 0 {
		pageCount++
	}

	// 主循环
	for i := 0; i < pageCount; i++ {
		pageSize := 20

		// 处理非完整尾页
		if i+1 == pageCount && count%20 != 0 {
			pageSize = count % 20
		}

		// 获取当前分页信息
		respJson, err := bilibili.GetProfileVideo(strconv.Itoa(mid), strconv.Itoa(i+1), "20", sessdata)
		if err != nil {
			return err
		}
		// 遍历分页
		for j := 0; j < pageSize; j++ {
			// 添加视频到列表
			err := VideoList.AddVideo(sessdata, gjson.Get(respJson, "data.list.vlist."+strconv.Itoa(j)+".bvid").String(), downloadCompilation)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

// 读取视频列表
func (VideoList *VideoList) Get(path ...string) error {
	// 指定路径
	filePath := config.Cfg.GetVideolistPath()
	if len(path) > 0 {
		filePath = path[0]
	}

	err := utils.LoadJsonFile(filePath, VideoList)
	if err != nil {
		return err
	}
	return err
}

// 保存视频列表
func (VideoList *VideoList) Save(path ...string) error {
	// 指定路径
	filePath := config.Cfg.GetVideolistPath()
	if len(path) > 0 {
		filePath = path[0]
	}

	err := utils.SaveJsonFile(filePath, VideoList)
	if err != nil {
		return err
	}
	return nil
}

// 获取视频流
// TODO：请求前检查数据
func (v *VideoInformation) GetStream(sessdata string) error {
	// 请求信息
	json, err := bilibili.GetVideoStream(v.Bvid, strconv.Itoa(v.Cid), sessdata)
	if err != nil {
		return err
	}
	// 错误检查
	if gjson.Get(json, "code").Int() != 0 {
		return errors.New(gjson.Get(json, "message").String())
	}

	// 选择音频流
	if gjson.Get(json, "data.dash.flac.audio").String() != "" {
		v.Audio.Quality = int(gjson.Get(json, "data.dash.audio.id").Int())
		v.Audio.Stream = gjson.Get(json, "data.dash.flac.audio.base_url").String()
		v.Format = constants.AudioType.Flac

		return nil
	}
	v.Audio.Quality = int(gjson.Get(json, "data.dash.audio.0.id").Int())
	v.Audio.Stream = gjson.Get(json, "data.dash.audio.0.base_url").String()

	return nil
}

func (videoList *VideoList) Tidy() {
	if len(videoList.List) == 0 {
		return
	}

	result := videoList.List[:0]
	for _, video := range videoList.List {
		if !video.Delete {
			result = append(result, video)
		}
	}
	videoList.List = result
	videoList.Count = len(result)
}
