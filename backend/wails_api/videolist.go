package wails_api

import (
	"bili-audio-downloader/backend/adapter"
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/download"
	"encoding/json"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

	const PageSize = 20
	length := len(download.DownloadList)
	start := page * PageSize
	end := start + PageSize

	if start >= length {
		return []adapter.TaskInfo{}
	}

	if end > length {
		end = length
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
func (w *WailsApi) AddCollectionToList(listPath, fid string, offset, count int, downloadCompilation bool) error {
	sessdata := ""
	if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
		sessdata = config.Cfg.Account.SESSDATA
	}

	err := download.AddCollectionTask(sessdata, fid, offset, count, downloadCompilation)
	if err != nil {
		return err
	}

	return nil
}

// AddCompilationToList 添加视频合集
func (w *WailsApi) AddCompilationToList(listPath string, mid, sid, offset, count int, downloadCompilation bool) error {
	sessdata := ""
	if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
		sessdata = config.Cfg.Account.SESSDATA
	}

	err := download.AddCompilationTask(sessdata, mid, sid, offset, count, downloadCompilation)
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
func (w *WailsApi) AddProfileVideoToList(listPath string, mid, offset, count int, downloadCompilation bool) error {
	sessdata := ""
	if config.Cfg.Account.IsLogin && config.Cfg.Account.UseAccount {
		sessdata = config.Cfg.Account.SESSDATA
	}

	err := download.AddProfileVideoTask(sessdata, mid, offset, count, downloadCompilation)
	if err != nil {
		return err
	}

	return nil
}

// SetTaskDeleteState 设置任务删除状态
func (w *WailsApi) SetTaskDeleteState(index int, delete bool) {
	if index >= 0 && index < len(download.DownloadList) {
		download.DownloadList[index].SetDelete(delete)
	}
}

// UpdateTaskMeta 更新任务元数据
func (w *WailsApi) UpdateTaskMeta(index int, songName, author string) {
	if index >= 0 && index < len(download.DownloadList) {
		download.DownloadList[index].SetMeta(songName, author)
	}
}

// ExportVideoList 导出任务列表（带对话框）
func (w *WailsApi) ExportVideoList() error {
	path, err := runtime.SaveFileDialog(w.ctx, runtime.SaveDialogOptions{
		Title:           "保存列表",
		DefaultFilename: "video_list.json",
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON Files (*.json)", Pattern: "*.json"},
		},
	})
	if err != nil {
		return err
	}
	if path == "" {
		return nil
	}
	return w.SaveVideoListTo(path)
}

// SaveVideoListTo 导出任务列表
func (w *WailsApi) SaveVideoListTo(path string) error {
	var listToExport []adapter.TaskInfo
	for i, task := range download.DownloadList {
		info := task.GetTaskInfo()
		info.Index = i
		listToExport = append(listToExport, *info)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(listToExport)
	if err != nil {
		return err
	}

	return nil
}
