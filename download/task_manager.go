package download

import (
	"bili-audio-downloader/bilibili"
	"bili-audio-downloader/utils"
	"github.com/tidwall/gjson"
	"strconv"
)

var DownloadList []DownloadTask

func AddTask(task DownloadTask) {
	DownloadList = append(DownloadList, task)
}

// AddVideoTask 向列表中添加一个视频
func AddVideoTask(sessdata, bvid string, downloadCompilation bool) error {
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
		Title := utils.CheckFileName(video.Meta.Title)
		PageTitle := utils.CheckFileName(video.Videos[i].Part)

		// 处理音频标题（单 P 视频）
		var SongName string
		if total <= 1 {
			// 单集使用视频标题
			SongName, err = utils.ExtractTitle(Title)
			if err != nil {
				SongName = Title
			}
		} else {
			// 多集视频使用分集标题
			SongName, err = utils.ExtractTitle(PageTitle)
			if err != nil {
				SongName = PageTitle
			}
		}

		metaData := MetaData{
			Title:     Title,
			PageTitle: PageTitle,
			PartId:    0,
			SongName:  SongName,
			Author:    video.Up.Name,
			LyricsUrl: "",
		}

		video := NewVideo(video.Bvid, video.Videos[i].Cid, video.Meta.Cover, sessdata, metaData)

		AddTask(video)
	}
	return nil
}

// AddAudio 向列表中添加一个音频项目
func AddAudioTask(sessdata, auid string) error {
	// 查询视频信息
	audio := new(bilibili.Audio)
	err := audio.Query(auid)
	if err != nil {
		return err
	}

	metaData := MetaData{
		Title:     utils.CheckFileName(audio.Meta.Title),
		PageTitle: utils.CheckFileName(audio.Meta.Title),
		PartId:    0,
		SongName:  audio.Meta.Title,
		Author:    audio.Up.Author,
		LyricsUrl: audio.Meta.Lyric,
	}

	audioTask := NewAudio(auid, audio.Meta.Cover, sessdata, metaData)
	AddTask(audioTask)

	return nil
}

// AddCollectionTask 向列表中添加一个收藏夹
func AddCollectionTask(sessdata, favlistId string, count int, downloadCompilation bool) error {
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
				err := AddVideoTask(sessdata, favlist.Data.Medias[j].Bvid, downloadCompilation)
				if err != nil {
					continue
				}
			} else {
				// 添加收藏夹中的音频
				err := AddAudioTask(sessdata, strconv.Itoa(favlist.Data.Medias[j].Id))
				if err != nil {
					continue
				}
			}
		}
	}
	return nil
}

// AddCompilationTask 向列表中添加一个视频合集
func AddCompilationTask(sessdata string, mid, sid, count int, downloadCompilation bool) error {
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
			err := AddVideoTask(sessdata, favlist.Data.Archives[j].Bvid, downloadCompilation)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

// AddProfileVideoTask 向列表中添加个人主页视频
func AddProfileVideoTask(sessdata string, mid, count int, downloadCompilation bool) error {
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
			err := AddVideoTask(sessdata, gjson.Get(respJson, "data.list.vlist."+strconv.Itoa(j)+".bvid").String(), downloadCompilation)
			if err != nil {
				continue
			}
		}
	}

	return nil
}
