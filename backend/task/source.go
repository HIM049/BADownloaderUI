package task

import (
	"bili-audio-downloader/backend/biliapi"
	"bili-audio-downloader/backend/utils"
	"fmt"
	"strconv"
)

// pageSize 批量添加时每页请求的条目数
const pageSize = 20

// AddVideo 向列表中添加一个视频，downloadCompilation 为真时添加全部分 P
func (l *List) AddVideo(sessdata, bvid string, downloadCompilation bool, option Option) error {
	view, err := biliapi.GetVideoView(bvid, sessdata)
	if err != nil {
		return err
	}
	if len(view.Pages) == 0 {
		return fmt.Errorf("video %s has no pages", bvid)
	}

	total := 1
	if downloadCompilation {
		total = len(view.Pages)
	}

	title := utils.CheckFileName(view.Title)
	for i := range total {
		pageTitle := utils.CheckFileName(view.Pages[i].Part)

		// 单 P 使用视频标题作为歌曲名，多 P 使用分集标题
		songName := title
		if total > 1 {
			songName = pageTitle
		}
		if extracted, err := utils.ExtractTitle(songName); err == nil {
			songName = extracted
		}

		l.Add(&Task{
			Kind:     KindVideo,
			Bvid:     bvid,
			Cid:      view.Pages[i].Cid,
			CoverURL: view.Pic,
			Sessdata: sessdata,
			Option:   option,
			Meta: MetaData{
				Title:     title,
				PageTitle: pageTitle,
				PartId:    i,
				SongName:  songName,
				Author:    view.Owner.Name,
			},
		})
	}
	return nil
}

// AddAudio 向列表中添加一个音频
func (l *List) AddAudio(sessdata, auid string, option Option) error {
	info, err := biliapi.GetAudioInfo(auid)
	if err != nil {
		return err
	}

	title := utils.CheckFileName(info.Title)
	l.Add(&Task{
		Kind:     KindAudio,
		Auid:     auid,
		CoverURL: info.Cover,
		Sessdata: sessdata,
		Option:   option,
		Meta: MetaData{
			Title:     title,
			PageTitle: title,
			SongName:  info.Title,
			Author:    info.Author,
			LyricsUrl: info.Lyric,
		},
	})
	return nil
}

// AddFavList 向列表中添加收藏夹内容，单条添加失败时跳过
func (l *List) AddFavList(sessdata, favlistId string, offset, count int, downloadCompilation bool, option Option) error {
	first, err := biliapi.GetFavList(favlistId, sessdata, 1, 1)
	if err != nil {
		return err
	}

	return eachInRange(offset, count, first.Info.MediaCount,
		func(pn int) ([]biliapi.FavMedia, error) {
			page, err := biliapi.GetFavList(favlistId, sessdata, pageSize, pn)
			return page.Medias, err
		},
		func(media biliapi.FavMedia) {
			if media.Type == biliapi.FavMediaVideo {
				_ = l.AddVideo(sessdata, media.Bvid, downloadCompilation, option)
			} else {
				_ = l.AddAudio(sessdata, strconv.Itoa(media.ID), option)
			}
		},
	)
}

// AddSeason 向列表中添加视频合集内容，单条添加失败时跳过
func (l *List) AddSeason(sessdata string, mid, sid, offset, count int, downloadCompilation bool, option Option) error {
	first, err := biliapi.GetSeasonArchives(mid, sid, 1, 1)
	if err != nil {
		return err
	}

	return eachInRange(offset, count, first.Meta.Total,
		func(pn int) ([]string, error) {
			page, err := biliapi.GetSeasonArchives(mid, sid, pageSize, pn)
			bvids := make([]string, 0, len(page.Archives))
			for _, archive := range page.Archives {
				bvids = append(bvids, archive.Bvid)
			}
			return bvids, err
		},
		func(bvid string) {
			_ = l.AddVideo(sessdata, bvid, downloadCompilation, option)
		},
	)
}

// AddSpaceVideos 向列表中添加用户主页投稿视频，单条添加失败时跳过
func (l *List) AddSpaceVideos(sessdata string, mid, offset, count int, downloadCompilation bool, option Option) error {
	first, err := biliapi.GetSpaceVideos(mid, sessdata, 1, 1)
	if err != nil {
		return err
	}

	return eachInRange(offset, count, first.Page.Count,
		func(pn int) ([]string, error) {
			page, err := biliapi.GetSpaceVideos(mid, sessdata, pageSize, pn)
			bvids := make([]string, 0, len(page.List.Vlist))
			for _, video := range page.List.Vlist {
				bvids = append(bvids, video.Bvid)
			}
			return bvids, err
		},
		func(bvid string) {
			if bvid != "" {
				_ = l.AddVideo(sessdata, bvid, downloadCompilation, option)
			}
		},
	)
}

// eachInRange 从第 offset 条（从 1 开始）起按页遍历 count 条，count 为 0 表示到末尾
func eachInRange[T any](offset, count, total int, fetch func(pn int) ([]T, error), handle func(T)) error {
	start := max(offset-1, 0)
	end := total
	if count != 0 {
		end = min(start+count, total)
	}

	for page := start / pageSize; page*pageSize < end; page++ {
		items, err := fetch(page + 1)
		if err != nil {
			return err
		}
		for j, item := range items {
			idx := page*pageSize + j
			if idx >= start && idx < end {
				handle(item)
			}
		}
	}
	return nil
}
