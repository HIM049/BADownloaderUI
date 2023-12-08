package main

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 创建并保存视频列表
func (a *App) MakeAndSaveList(favlistID string, downloadCount int, downloadCompilation bool) error {
	cfg := GetConfig()
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

	// 请求收藏夹基础数据，用于初始化循环
	favlist, err := GetFavListObj(favlistID, 1, 1)
	if err != nil {
		return nil, err
	}
	// 计算下载量
	if downloadCount == 0 {
		downloadCount = favlist.Data.Info.Media_count
	}
	// // 设置进度条
	// progressBar := pb.Full.Start(downloadCount)

	// 主循环
	// TODO：修改请求数量为一次 20
	for i := 0; i < downloadCount; i++ {
		// 获取当前分页信息
		favlist, err := GetFavListObj(favlistID, 1, i+1)
		if err != nil {
			return nil, err
		}
		// 获取当前视频详细信息
		videoInf, err := GetVideoPageInformationObj(favlist.Data.Medias[0].Bvid)
		if err != nil {
			// 视频失效
			// fmt.Printf("获取 "+favlist.Data.Medias[0].Bvid+" 信息时发生错误：%s\n", err)
			runtime.LogError(a.ctx, "获取 "+favlist.Data.Medias[0].Bvid+" 信息时发生错误: "+err.Error())
			continue
		}

		// 处理音频标题（单 P 视频）
		songName, err := ExtractTitle(CheckFileName(videoInf.Data.Title))
		if err != nil {
			// 如果无法判断标题
			songName = CheckFileName(videoInf.Data.Title)
		}

		// 如果是多 P
		if videoInf.Data.Videos > 1 && downloadCompilation {
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
					ListID:    i,
					IsPage:    true,
					PageTitle: CheckFileName(pages.Part),
					PageID:    pages.Page,
				}
				// 组合数据
				videoList = append(videoList, videoPage)
			}
		} else {
			// 填充单 P 数据
			video := VideoInformationList{
				Bvid:     videoInf.Data.Bvid,
				Cid:      videoInf.Data.Cid,
				Title:    CheckFileName(videoInf.Data.Title),
				SongName: songName,
				Author:   videoInf.Data.Owner.Name,
				Cover:    videoInf.Data.Pic,
				Videos:   videoInf.Data.Videos,
				ListID:   i,
				IsPage:   false,
			}
			// 组合数据
			videoList = append(videoList, video)
		}

		// // 进度条增加
		// progressBar.Increment()
	}
	// // 取消进度条显示
	// progressBar.Finish()
	return &videoList, nil
}
