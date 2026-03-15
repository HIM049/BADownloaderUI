package download

import (
	"bili-audio-downloader/backend/adapter"
	bilibili2 "bili-audio-downloader/backend/adapter/bilibili"
	"bili-audio-downloader/backend/utils"
	"bili-audio-downloader/bilibili"
	"strconv"

	"github.com/tidwall/gjson"
)

var DownloadList []DownloadTask

func ResetTaskList() {
	DownloadList = []DownloadTask{}
}

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

		metaData := adapter.MetaData{
			Title:     Title,
			PageTitle: PageTitle,
			PartId:    i,
			SongName:  SongName,
			Author:    video.Up.Name,
			LyricsUrl: "",
		}

		video := bilibili2.NewVideo(video.Bvid, video.Videos[i].Cid, video.Meta.Cover, sessdata, metaData)

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

	metaData := adapter.MetaData{
		Title:     utils.CheckFileName(audio.Meta.Title),
		PageTitle: utils.CheckFileName(audio.Meta.Title),
		PartId:    0,
		SongName:  audio.Meta.Title,
		Author:    audio.Up.Author,
		LyricsUrl: audio.Meta.Lyric,
	}

	audioTask := bilibili2.NewAudio(auid, audio.Meta.Cover, sessdata, metaData)
	AddTask(audioTask)

	return nil
}

// AddCollectionTask 向列表中添加一个收藏夹
func AddCollectionTask(sessdata, favlistId string, offset, count int, downloadCompilation bool) error {
	// 请求收藏夹基础数据，初始化循环
	favlist, err := bilibili.GetFavListObj(favlistId, sessdata, 1, 1)
	if err != nil {
		return err
	}
	// 计算起始和结束索引
	startIdx := offset - 1
	if startIdx < 0 {
		startIdx = 0
	}
	
	maxCount := favlist.Data.Info.Media_count
	if count == 0 || startIdx+count > maxCount {
		count = maxCount - startIdx
	}
	if count <= 0 {
		return nil
	}

	endIdx := startIdx + count - 1

	startPage := startIdx / 20
	endPage := endIdx / 20

	// 主循环
	for i := startPage; i <= endPage; i++ {
		// 获取当前分页信息
		favlist, err := bilibili.GetFavListObj(favlistId, sessdata, 20, i+1)
		if err != nil {
			return err
		}

		pageStart := 0
		if i == startPage {
			pageStart = startIdx % 20
		}
		
		pageEnd := 19
		if i == endPage {
			pageEnd = endIdx % 20
		}

		// 遍历分页
		for j := pageStart; j <= pageEnd && j < len(favlist.Data.Medias); j++ {

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
func AddCompilationTask(sessdata string, mid, sid, offset, count int, downloadCompilation bool) error {
	// 请求收藏夹基础数据，初始化循环
	favlist, err := bilibili.GetCompliationObj(mid, sid, 1, 1)
	if err != nil {
		return err
	}
	// 计算起始和结束索引
	startIdx := offset - 1
	if startIdx < 0 {
		startIdx = 0
	}
	
	maxCount := favlist.Data.Meta.Total
	if count == 0 || startIdx+count > maxCount {
		count = maxCount - startIdx
	}
	if count <= 0 {
		return nil
	}

	endIdx := startIdx + count - 1

	startPage := startIdx / 20
	endPage := endIdx / 20

	// 主循环
	for i := startPage; i <= endPage; i++ {
		// 获取当前分页信息
		favlist, err := bilibili.GetCompliationObj(mid, sid, 20, i+1)
		if err != nil {
			return err
		}
		
		pageStart := 0
		if i == startPage {
			pageStart = startIdx % 20
		}
		
		pageEnd := 19
		if i == endPage {
			pageEnd = endIdx % 20
		}
		
		// 遍历分页
		for j := pageStart; j <= pageEnd && j < len(favlist.Data.Archives); j++ {
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
func AddProfileVideoTask(sessdata string, mid, offset, count int, downloadCompilation bool) error {
	respJson, err := bilibili.GetProfileVideo(strconv.Itoa(mid), "1", "1", sessdata)
	if err != nil {
		return err
	}

	// 计算起始和结束索引
	startIdx := offset - 1
	if startIdx < 0 {
		startIdx = 0
	}
	
	maxCount := int(gjson.Get(respJson, "data.page.count").Int())
	if count == 0 || startIdx+count > maxCount {
		count = maxCount - startIdx
	}
	if count <= 0 {
		return nil
	}

	endIdx := startIdx + count - 1

	startPage := startIdx / 20
	endPage := endIdx / 20

	// 主循环
	for i := startPage; i <= endPage; i++ {
		// 获取当前分页信息
		respJson, err := bilibili.GetProfileVideo(strconv.Itoa(mid), strconv.Itoa(i+1), "20", sessdata)
		if err != nil {
			return err
		}
		
		pageStart := 0
		if i == startPage {
			pageStart = startIdx % 20
		}
		
		pageEnd := 19
		if i == endPage {
			pageEnd = endIdx % 20
		}
		
		// 遍历分页
		for j := pageStart; j <= pageEnd; j++ {
			bvid := gjson.Get(respJson, "data.list.vlist."+strconv.Itoa(j)+".bvid").String()
			if bvid == "" {
				continue
			}
			// 添加视频到列表
			err := AddVideoTask(sessdata, gjson.Get(respJson, "data.list.vlist."+strconv.Itoa(j)+".bvid").String(), downloadCompilation)
			if err != nil {
				continue
			}
		}
	}

	return nil
}
