package wails_api

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/download"
)

// GetListCount 获取列表中视频数量
func (a *WailsApi) GetListCount() int {
	return len(download.DownloadList)
}

//// CreatVideoList 创建视频列表
//func (a *WailsApi) CreatVideoList() error {
//	videoList := new(services.VideoList)
//	err := videoList.Save()
//	if err != nil {
//		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
//		return err
//	}
//	return nil
//}

// AddVideoToList 添加单个视频
func (a *WailsApi) AddVideoToList(listPath, bvid string, downloadCompilation bool) error {
	sessdata := ""
	if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
		sessdata = config.Cfg.Account.SESSDATA
	}

	err := download.AddVideoTask(sessdata, bvid, downloadCompilation)
	if err != nil {
		return err
	}

	return nil
}

// AddCollectionToList 添加收藏夹内容
func (a *WailsApi) AddCollectionToList(listPath, fid string, count int, downloadCompilation bool) error {
	sessdata := ""
	if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
		sessdata = config.Cfg.Account.SESSDATA
	}

	err := download.AddCollectionTask(sessdata, fid, count, downloadCompilation)
	if err != nil {
		return err
	}

	return nil
}

// AddCompilationToList 添加视频合集
func (a *WailsApi) AddCompilationToList(listPath string, mid, sid, count int, downloadCompilation bool) error {
	sessdata := ""
	if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
		sessdata = config.Cfg.Account.SESSDATA
	}

	err := download.AddCompilationTask(sessdata, mid, sid, count, downloadCompilation)
	if err != nil {
		return err
	}

	return nil
}

// AddAudioToList 添加单个音频
func (a *WailsApi) AddAudioToList(listPath, auid string) error {
	sessdata := ""
	if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
		sessdata = config.Cfg.Account.SESSDATA
	}

	err := download.AddAudioTask(sessdata, auid)
	if err != nil {
		return err
	}

	return nil
}

// AddProfileVideoToList 添加个人主页视频
func (a *WailsApi) AddProfileVideoToList(listPath string, mid, count int, downloadCompilation bool) error {
	sessdata := ""
	if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
		sessdata = config.Cfg.Account.SESSDATA
	}

	err := download.AddProfileVideoTask(sessdata, mid, count, downloadCompilation)
	if err != nil {
		return err
	}

	return nil
}

//// LoadVideoList 加载视频列表
//func (a *WailsApi) LoadVideoList(listPath string) (services.VideoList, error) {
//	videoList := new(services.VideoList)
//	err := videoList.Get(listPath)
//	if err != nil {
//		return services.VideoList{}, err
//	}
//	return *videoList, nil
//}

//// SaveVideoList 保存视频列表
//func (a *WailsApi) SaveVideoList(newList services.VideoList, path string) error {
//	err := newList.Save(path)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//// TidyVideoList 删除列表中的废弃项
//func (a *WailsApi) TidyVideoList(listPath string) error {
//	videoList := new(services.VideoList)
//	err := videoList.Get(listPath)
//	if err != nil {
//		return err
//	}
//
//	videoList.Tidy()
//
//	err = videoList.Save(listPath)
//	if err != nil {
//		return err
//	}
//	return nil
//}
