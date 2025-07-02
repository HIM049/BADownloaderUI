package wails_api

import (
	"bili-audio-downloader/backend/adapter"
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/download"
)

// GetListCount 获取列表中视频数量
func (w *WailsApi) GetListCount() int {
	return len(download.DownloadList)
}

// ResetDownloadList 重置列表
func (w *WailsApi) ResetDownloadList() {
	download.ResetTaskList()
}

// GetTaskListAll 获取完整的任务列表信息
func (w *WailsApi) GetTaskListAll() []adapter.TaskInfo {
	var taskList []adapter.TaskInfo

	for i, task := range download.DownloadList {
		info := task.GetTaskInfo()
		info.Index = i
		taskList = append(taskList, *info)
	}

	return taskList
}

// GetTaskListPage 获取一页的任务列表信息
func (w *WailsApi) GetTaskListPage(page int) []adapter.TaskInfo {
	var taskList []adapter.TaskInfo

	const PageSize = 10
	start := page * PageSize
	end := page*PageSize + PageSize

	if end > len(download.DownloadList) {
		end = len(download.DownloadList)
	}

	for i, task := range download.DownloadList[start:end] {
		info := task.GetTaskInfo()
		info.Index = i + start
		taskList = append(taskList, *info)
	}

	return taskList
}

// AddVideoToList 添加单个视频
func (w *WailsApi) AddVideoToList(listPath, bvid string, downloadCompilation bool) error {
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
func (w *WailsApi) AddCollectionToList(listPath, fid string, count int, downloadCompilation bool) error {
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
func (w *WailsApi) AddCompilationToList(listPath string, mid, sid, count int, downloadCompilation bool) error {
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
func (w *WailsApi) AddAudioToList(listPath, auid string) error {
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
func (w *WailsApi) AddProfileVideoToList(listPath string, mid, count int, downloadCompilation bool) error {
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
