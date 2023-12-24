package main

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 创建并保存视频列表
func (a *App) MakeAndSaveList(favlistID string, downloadCount int, downloadCompilation bool) error {
	cfg := GetConfig(a.ctx)
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

// 视频数据结构
type VideoInformationList struct {
	Bvid      string `json:"bvid"`
	Cid       int    `json:"cid"`
	Title     string `json:"title"`
	SongName  string `json:"song_name"`
	Author    string `json:"author"`
	Cover     string `json:"cover"`
	Videos    int    `json:"videos"`
	ListID    int    `json:"list_id"`
	IsPage    bool   `json:"is_page"`
	PageTitle string `json:"page_title"`
	PageID    int    `json:"page_id"`
}

// 创建视频任务列表
func makeVideoList(a *App, favlistID string, downloadCount int, downloadCompilation bool) (*[]VideoInformationList, error) {
	var videoList []VideoInformationList

	// 请求收藏夹基础数据，初始化循环
	favlist, err := GetFavListObj(favlistID, 1, 1)
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
		favlist, err := GetFavListObj(favlistID, pageSize, i+1)
		if err != nil {
			return nil, err
		}
		for j, listVideo := range favlist.Data.Medias {
			// 获取当前视频详细信息
			videoInf, err := GetVideoPageInformationObj(listVideo.Bvid)
			if err != nil {
				// 视频失效
				runtime.LogErrorf(a.ctx, "获取 "+listVideo.Bvid+" 信息时发生错误: %s", err)
				continue
			}

			// 处理音频标题（单 P 视频）
			songName, err := ExtractTitle(CheckFileName(videoInf.Data.Title))
			if err != nil {
				// 如果无法判断标题
				songName = CheckFileName(videoInf.Data.Title)
			}

			// 分 P 判断
			if videoInf.Data.Videos > 1 && downloadCompilation {
				// 如果是多 P
				for _, pages := range videoInf.Data.Pages {

					// 处理音频标题（分 P 视频）
					songName, err = ExtractTitle(CheckFileName(pages.Part))
					if err != nil {
						// 如果无法判断标题
						songName = CheckFileName(pages.Part)
					}

					// 填充 Page 数据
					videoPage := VideoInformationList{
						Bvid:      videoInf.Data.Bvid,
						Cid:       pages.Cid,
						Title:     CheckFileName(videoInf.Data.Title),
						SongName:  songName,
						Author:    videoInf.Data.Owner.Name,
						Cover:     videoInf.Data.Pic,
						Videos:    videoInf.Data.Videos,
						ListID:    (i * 20) + j,
						IsPage:    true,
						PageTitle: CheckFileName(pages.Part),
						PageID:    pages.Page,
					}
					// 组合数据
					videoList = append(videoList, videoPage)
				}
			} else {
				// 如果是单 P
				video := VideoInformationList{
					Bvid:     videoInf.Data.Bvid,
					Cid:      videoInf.Data.Cid,
					Title:    CheckFileName(videoInf.Data.Title),
					SongName: songName,
					Author:   videoInf.Data.Owner.Name,
					Cover:    videoInf.Data.Pic,
					Videos:   videoInf.Data.Videos,
					ListID:   (i * 20) + j,
					IsPage:   false,
				}
				// 组合数据
				videoList = append(videoList, video)
			}

		}
	}

	return &videoList, nil
}
