package download

import (
	"bili-audio-downloader/backend/adapter"
	"bili-audio-downloader/backend/adapter/bilibili"
	"bili-audio-downloader/backend/biliapi"
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/utils"
	"fmt"
	"strconv"
)

var DownloadList []DownloadTask

func ResetTaskList() {
	DownloadList = []DownloadTask{}
}

func AddTask(task DownloadTask) {
	DownloadList = append(DownloadList, task)
}

// ImportTasks 从导出的任务信息中恢复任务列表
func ImportTasks(tasks []adapter.TaskInfo) {
	for _, info := range tasks {
		sessdata := info.Sessdata
		if sessdata == "" {
			if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
				sessdata = config.Cfg.Account.SESSDATA
			}
		}

		switch info.Type {
		case 1: // Video
			video := bilibili.NewVideo(info.Bvid, info.Cid, info.CoverUrl, sessdata, info.MetaData)
			video.SetDelete(info.IsDelete)
			AddTask(video)
		case 2: // Audio
			audio := bilibili.NewAudio(info.Auid, info.CoverUrl, sessdata, info.MetaData)
			audio.SetDelete(info.IsDelete)
			AddTask(audio)
		}
	}
}

// AddVideoTask 向列表中添加一个视频
func AddVideoTask(sessdata, bvid string, downloadCompilation bool) error {
	// 查询视频信息
	view, err := biliapi.GetVideoView(bvid, sessdata)
	if err != nil {
		return err
	}
	if len(view.Pages) == 0 {
		return fmt.Errorf("video %s has no pages", bvid)
	}

	// 处理分集数量
	var total int = 1
	if downloadCompilation {
		total = len(view.Pages)
	}

	// 保存信息
	for i := 0; i < total; i++ {
		Title := utils.CheckFileName(view.Title)
		PageTitle := utils.CheckFileName(view.Pages[i].Part)

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
			Author:    view.Owner.Name,
			LyricsUrl: "",
		}

		video := bilibili.NewVideo(bvid, view.Pages[i].Cid, view.Pic, sessdata, metaData)

		AddTask(video)
	}
	return nil
}

// AddAudio 向列表中添加一个音频项目
func AddAudioTask(sessdata, auid string) error {
	// 查询音频信息
	info, err := biliapi.GetAudioInfo(auid)
	if err != nil {
		return err
	}

	metaData := adapter.MetaData{
		Title:     utils.CheckFileName(info.Title),
		PageTitle: utils.CheckFileName(info.Title),
		PartId:    0,
		SongName:  info.Title,
		Author:    info.Author,
		LyricsUrl: info.Lyric,
	}

	audioTask := bilibili.NewAudio(auid, info.Cover, sessdata, metaData)
	AddTask(audioTask)

	return nil
}

// AddCollectionTask 向列表中添加一个收藏夹
func AddCollectionTask(sessdata, favlistId string, offset, count int, downloadCompilation bool) error {
	// 请求收藏夹基础数据，初始化循环
	favlist, err := biliapi.GetFavList(favlistId, sessdata, 1, 1)
	if err != nil {
		return err
	}
	// 计算起始和结束索引
	startIdx := offset - 1
	if startIdx < 0 {
		startIdx = 0
	}

	maxCount := favlist.Info.MediaCount
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
		favlist, err := biliapi.GetFavList(favlistId, sessdata, 20, i+1)
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
		for j := pageStart; j <= pageEnd && j < len(favlist.Medias); j++ {

			if favlist.Medias[j].Type == biliapi.FavMediaVideo {
				// 添加视频到列表
				err := AddVideoTask(sessdata, favlist.Medias[j].Bvid, downloadCompilation)
				if err != nil {
					continue
				}
			} else {
				// 添加收藏夹中的音频
				err := AddAudioTask(sessdata, strconv.Itoa(favlist.Medias[j].ID))
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
	season, err := biliapi.GetSeasonArchives(mid, sid, 1, 1)
	if err != nil {
		return err
	}
	// 计算起始和结束索引
	startIdx := offset - 1
	if startIdx < 0 {
		startIdx = 0
	}

	maxCount := season.Meta.Total
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
		season, err := biliapi.GetSeasonArchives(mid, sid, 20, i+1)
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
		for j := pageStart; j <= pageEnd && j < len(season.Archives); j++ {
			// 添加视频到列表
			err := AddVideoTask(sessdata, season.Archives[j].Bvid, downloadCompilation)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

// AddProfileVideoTask 向列表中添加个人主页视频
func AddProfileVideoTask(sessdata string, mid, offset, count int, downloadCompilation bool) error {
	space, err := biliapi.GetSpaceVideos(mid, sessdata, 1, 1)
	if err != nil {
		return err
	}

	// 计算起始和结束索引
	startIdx := offset - 1
	if startIdx < 0 {
		startIdx = 0
	}

	maxCount := space.Page.Count
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
		space, err := biliapi.GetSpaceVideos(mid, sessdata, 20, i+1)
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
		for j := pageStart; j <= pageEnd && j < len(space.List.Vlist); j++ {
			bvid := space.List.Vlist[j].Bvid
			if bvid == "" {
				continue
			}
			// 添加视频到列表
			err := AddVideoTask(sessdata, bvid, downloadCompilation)
			if err != nil {
				continue
			}
		}
	}

	return nil
}
