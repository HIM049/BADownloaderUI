package main

import (
	"errors"
	"path/filepath"
	"strconv"

	"github.com/myuser/bilibili"
	"github.com/tidwall/gjson"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 获取版本号
func (a *App) GetAppVersion() string {
	return APP_VERSION
}

// 获取主题字符串
func (a *App) GetTheme() (string, error) {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return "", err
	}

	return cfg.Theme, nil
}

// 获取列表中视频数量
func (a *App) GetListCount(path string) int {
	videoList := new(VideoList)
	err := videoList.Get(path)
	if err != nil {
		return 0
	}
	return videoList.Count
}

// 查询视频信息
func (a *App) QueryVideo(bvid string) (bilibili.Video, error) {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return bilibili.Video{}, err
	}
	sessdata := ""
	if cfg.Account.UseAccount && cfg.Account.IsLogin {
		sessdata = cfg.Account.SESSDATA
	}

	video := new(bilibili.Video)
	err = video.Query(sessdata, bvid)
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return bilibili.Video{}, err
	}
	return *video, err
}

// 查询并返回收藏夹信息
func (a *App) QueryCollection(favListID string) bilibili.FavList {
	cfg := new(Config)
	cfg.Get()
	sessdata := ""
	if cfg.Account.UseAccount && cfg.Account.IsLogin {
		sessdata = cfg.Account.SESSDATA
	}
	listInf, err := bilibili.GetFavListObj(favListID, sessdata, 1, 1)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取收藏夹内容时出现错误：%s", err)
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return bilibili.FavList{}
	}
	return *listInf
}

// 查询并返回合集信息
func (a *App) QueryCompilation(mid, sid int) bilibili.CompliationInformation {
	listInf, err := bilibili.GetCompliationObj(mid, sid, 1, 1)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取合集内容时出现错误：%s", err)
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return bilibili.CompliationInformation{}
	}
	return *listInf
}

// 查询音频信息
func (a *App) QueryAudio(auid string) (bilibili.Audio, error) {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return bilibili.Audio{}, err
	}

	audio := new(bilibili.Audio)
	err = audio.Query(auid)
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return bilibili.Audio{}, err
	}
	return *audio, err
}

// 查询音频信息
func (a *App) QueryProfileVideo(mid string) (int, error) {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return 0, err
	}

	sessdata := ""
	if cfg.Account.UseAccount && cfg.Account.IsLogin {
		sessdata = cfg.Account.SESSDATA
	}

	respJson, err := bilibili.GetProfileVideo(mid, "1", "1", sessdata)
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return 0, err
	}
	return int(gjson.Get(respJson, "data.page.count").Int()), err
}

// 创建视频列表
func (a *App) CreatVideoList() error {
	videoList := new(VideoList)
	err := videoList.Save()
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return err
	}
	return nil
}

// 添加单个视频
func (a *App) AddVideoToList(listPath, bvid string, downloadCompilation bool) error {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return err
	}

	videolist := new(VideoList)
	err = videolist.Get(listPath)
	if err != nil {
		return err
	}

	sessdata := ""
	if cfg.Account.IsLogin && cfg.Account.UseAccount {
		sessdata = cfg.Account.SESSDATA
	}

	err = videolist.AddVideo(sessdata, bvid, downloadCompilation)
	if err != nil {
		return err
	}

	videolist.Save(listPath)

	return nil
}

// 添加收藏夹内容
func (a *App) AddCollectionToList(listPath, fid string, count int, downloadCompilation bool) error {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return err
	}

	videoList := new(VideoList)
	err = videoList.Get(listPath)
	if err != nil {
		return err
	}

	sessdata := ""
	if cfg.Account.IsLogin && cfg.Account.UseAccount {
		sessdata = cfg.Account.SESSDATA
	}

	err = videoList.AddCollection(sessdata, fid, count, downloadCompilation)
	if err != nil {
		return err
	}

	err = videoList.Save(listPath)
	if err != nil {
		return err
	}

	return nil
}

// 添加视频合集
func (a *App) AddCompilationToList(listPath string, mid, sid, count int, downloadCompilation bool) error {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return err
	}

	videoList := new(VideoList)
	err = videoList.Get(listPath)
	if err != nil {
		return nil
	}

	sessdata := ""
	if cfg.Account.IsLogin && cfg.Account.UseAccount {
		sessdata = cfg.Account.SESSDATA
	}

	err = videoList.AddCompilation(sessdata, mid, sid, count, downloadCompilation)
	if err != nil {
		return err
	}

	err = videoList.Save(listPath)
	if err != nil {
		return err
	}

	return nil
}

// 添加单个音频
func (a *App) AddAudioToList(listPath, auid string) error {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return err
	}

	videolist := new(VideoList)
	err = videolist.Get(listPath)
	if err != nil {
		return err
	}

	sessdata := ""
	if cfg.Account.IsLogin && cfg.Account.UseAccount {
		sessdata = cfg.Account.SESSDATA
	}

	err = videolist.AddAudio(sessdata, auid)
	if err != nil {
		return err
	}

	videolist.Save(listPath)

	return nil
}

// 添加个人主页视频
func (a *App) AddProfileVideoToList(listPath string, mid, count int, downloadCompilation bool) error {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return err
	}

	videoList := new(VideoList)
	err = videoList.Get(listPath)
	if err != nil {
		return nil
	}

	sessdata := ""
	if cfg.Account.IsLogin && cfg.Account.UseAccount {
		sessdata = cfg.Account.SESSDATA
	}

	err = videoList.AddProfileVideo(sessdata, mid, count, downloadCompilation)
	if err != nil {
		return err
	}

	err = videoList.Save(listPath)
	if err != nil {
		return err
	}

	return nil
}

// 加载视频列表
func (a *App) LoadVideoList(listPath string) (VideoList, error) {
	videoList := new(VideoList)
	err := videoList.Get(listPath)
	if err != nil {
		return VideoList{}, err
	}
	return *videoList, nil
}

// 保存视频列表
func (a *App) SaveVideoList(newList VideoList, path string) error {
	err := newList.Save(path)
	if err != nil {
		return err
	}
	return nil
}

// 删除列表中的废弃项
func (a *App) TidyVideoList(listPath string) error {
	videoList := new(VideoList)
	err := videoList.Get(listPath)
	if err != nil {
		return err
	}

	videoList.Tidy()

	err = videoList.Save(listPath)
	if err != nil {
		return err
	}
	return nil
}

// 获取用户创建的收藏夹
func (a *App) GetUsersCollect() bilibili.Collects {
	// 获取设置
	cfg := new(Config)
	cfg.Get()

	// 获取收藏夹列表
	collects := new(bilibili.Collects)
	mid, _ := strconv.Atoi(cfg.Account.DedeUserID)
	collects.UserMid = mid
	err := collects.GetUsersCollect(cfg.Account.SESSDATA)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取收藏夹列表失败：%s", err)
		return bilibili.Collects{}
	}

	return *collects
}

// 获取收藏的收藏夹
func (a *App) GetFavCollect(pn int) bilibili.Collects {
	// 获取设置
	cfg := new(Config)
	cfg.Get()

	// 获取收藏夹列表
	collects := new(bilibili.Collects)
	mid, _ := strconv.Atoi(cfg.Account.DedeUserID)
	collects.UserMid = mid
	err := collects.GetFavCollect(cfg.Account.SESSDATA, 20, pn)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取收藏夹列表失败：%s", err)
		return bilibili.Collects{}
	}

	return *collects
}

// 查询并返回歌曲信息
func (a *App) QuerySongInformation(auid string) (bilibili.Audio, error) {
	audioInf := new(bilibili.Audio)
	err := audioInf.Query(auid)
	if err != nil {
		return bilibili.Audio{}, err
	}
	audioInf.GetStream("")
	return *audioInf, nil
}

// 调用打开文件窗口
func (a *App) OpenFileDialog() (string, error) {
	var FileFilter []wails.FileFilter

	fileFilter := wails.FileFilter{
		DisplayName: "视频下载列表 (*.json)",
		Pattern:     "*.json",
	}
	FileFilter = append(FileFilter, fileFilter)

	option := wails.OpenDialogOptions{
		DefaultDirectory: "./",
		DefaultFilename:  "",
		Title:            "打开本地列表文件",
		Filters:          FileFilter,
	}
	// 弹出对话框
	path, err := wails.OpenFileDialog(a.ctx, option)
	if err != nil {
		wails.LogErrorf(a.ctx, err.Error())
		return "", err
	}

	return path, nil
}

// 调用保存窗口
func (a *App) SaveVideoListTo(videolist VideoList) error {
	var FileFilter []wails.FileFilter

	fileFilter := wails.FileFilter{
		DisplayName: "视频下载列表 (*.json)",
		Pattern:     "*.json",
	}
	FileFilter = append(FileFilter, fileFilter)

	option := wails.SaveDialogOptions{
		DefaultDirectory: "./",
		DefaultFilename:  "BAD_VideoList",
		Title:            "另存视频列表",
		Filters:          FileFilter,
	}

	// 弹出对话框
	path, err := wails.SaveFileDialog(a.ctx, option)
	if err != nil {
		return err
	}

	// 用户取消操作
	if path == "" {
		wails.EventsEmit(a.ctx, "error", "未选择保存路径")
		return nil
	}

	// 保存列表
	err = videolist.Save(path)
	if err != nil {
		return err
	}
	return nil
}

// 获取已登录用户的信息
func (a *App) GetUserInf() (bilibili.AccountInformation, error) {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return bilibili.AccountInformation{}, err
	}

	if !cfg.Account.IsLogin {
		return bilibili.AccountInformation{}, errors.New("用户未登录")
	}
	sessdata := cfg.Account.SESSDATA

	accountInf := new(bilibili.AccountInformation)
	accountInf.GetUserInf(sessdata)

	return *accountInf, nil
}

// 重置设置文件
func (a *App) RefreshConfig() {
	cfg := new(Config)
	cfg.init()
	err := cfg.Save()
	if err != nil {
		wails.LogErrorf(a.ctx, "写入设置文件失败：%s", err)
		wails.EventsEmit(a.ctx, "error", "写入设置时出错:"+err.Error())

	}
}

// 读取设置
func (a *App) LoadConfig() Config {
	cfg := new(Config)
	cfg.Get()
	return *cfg
}

// 写入设置
func (a *App) SaveConfig(cfg Config) {
	err := cfg.Save()
	if err != nil {
		wails.LogErrorf(a.ctx, "写入设置文件失败：%s", err)
		wails.EventsEmit(a.ctx, "error", "写入设置时出错:"+err.Error())
	}
}

// 打开下载文件夹
func (a *App) OpenDownloadFolader() error {
	cfg := new(Config)
	err := cfg.Get()
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(cfg.FileConfig.DownloadPath)
	if err != nil {
		wails.EventsEmit(a.ctx, "error", "错误："+err.Error())
		return err
	}

	err = OpenFolder(absPath)
	if err != nil {
		return err
	}

	return nil
}
