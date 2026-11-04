package wails_api

import (
	"bili-audio-downloader/backend/task"
	"encoding/json"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// GetListCount 获取列表中视频数量
func (w *WailsApi) GetListCount() int {
	return w.tasks.Len()
}

// ResetDownloadList 重置列表
func (w *WailsApi) ResetDownloadList() {
	w.tasks.Reset()
}

// GetTaskListAll 获取完整的任务列表信息
func (w *WailsApi) GetTaskListAll() []task.TaskInfo {
	return w.tasks.All()
}

// GetTaskListPage 获取一页的任务列表信息
func (w *WailsApi) GetTaskListPage(page int) []task.TaskInfo {
	return w.tasks.Page(page, 20)
}

// AddVideoToList 添加单个视频
func (w *WailsApi) AddVideoToList(bvid string, downloadCompilation bool, option task.Option) error {
	return w.tasks.AddVideo(accountSessdata(), bvid, downloadCompilation, option)
}

// AddCollectionToList 添加收藏夹内容
func (w *WailsApi) AddCollectionToList(fid string, offset, count int, downloadCompilation bool, option task.Option) error {
	return w.tasks.AddFavList(accountSessdata(), fid, offset, count, downloadCompilation, option)
}

// AddCompilationToList 添加视频合集
func (w *WailsApi) AddCompilationToList(mid, sid, offset, count int, downloadCompilation bool, option task.Option) error {
	return w.tasks.AddSeason(accountSessdata(), mid, sid, offset, count, downloadCompilation, option)
}

// AddAudioToList 添加单个音频
func (w *WailsApi) AddAudioToList(auid string, option task.Option) error {
	return w.tasks.AddAudio(accountSessdata(), auid, option)
}

// AddProfileVideoToList 添加个人主页视频
func (w *WailsApi) AddProfileVideoToList(mid, offset, count int, downloadCompilation bool, option task.Option) error {
	return w.tasks.AddSpaceVideos(accountSessdata(), mid, offset, count, downloadCompilation, option)
}

// SetTaskDeleteState 设置任务删除状态
func (w *WailsApi) SetTaskDeleteState(index int, delete bool) {
	w.tasks.SetDeleted(index, delete)
}

// UpdateTaskMeta 更新任务元数据
func (w *WailsApi) UpdateTaskMeta(index int, songName, author string) {
	w.tasks.SetMeta(index, songName, author)
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
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(w.tasks.All())
	if err != nil {
		return err
	}

	return nil
}

// ImportVideoList 导入任务列表
func (w *WailsApi) ImportVideoList() error {
	path, err := runtime.OpenFileDialog(w.ctx, runtime.OpenDialogOptions{
		Title: "打开列表",
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

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var tasks []task.TaskInfo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		return err
	}

	sessdata := accountSessdata()
	w.tasks.Import(tasks, sessdata)
	runtime.EventsEmit(w.ctx, "refreshVideoList")
	return nil
}
