package main

import (
	"strconv"

	"github.com/myuser/bilibili"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// 获取版本号
func (a *App) GetAppVersion() string {
	return APP_VERSION
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

// 加载视频列表
func (a *App) LoadVideoList(listPath string) (VideoList, error) {
	videoList := new(VideoList)
	err := videoList.Get(listPath)
	if err != nil {
		return VideoList{}, err
	}
	return *videoList, nil
}

func (a *App) SaveVideoList(newList VideoList) error {
	err := newList.Save()
	if err != nil {
		return err
	}
	return nil
}

// func (a *App) TestCreatVL() {
// 	fmt.Println("Start Test")

// 	VideoList := new(VideoList)

// 	// 添加单个视频
// 	fmt.Println("Start Add Video")
// 	err := VideoList.AddVideo("", "BV1WK4y1z72p", false)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// 添加收藏夹
// 	fmt.Println("Start Add Collection")
// 	err = VideoList.AddCollection("", "742380048", 5, false)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// 添加合集
// 	fmt.Println("Start Add Compilation")
// 	err = VideoList.AddCompilation("", 449838148, 180937, 5, false)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	VideoList.Save()

// 	fmt.Println(VideoList)
// }

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
func (a *App) SearchSongInformation(auid string) bilibili.AudioInf {
	wails.LogInfo(a.ctx, auid)
	audioInf, err := bilibili.GetAudioInfObj(auid)
	if err != nil {
		wails.LogErrorf(a.ctx, "获取歌曲详情时出现错误：%s", err)
		wails.EventsEmit(a.ctx, "error", "获取歌曲时出错:"+err.Error())
		return bilibili.AudioInf{}
	}
	return *audioInf
}

// 调用 Windows 打开文件窗口
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
	path, err := wails.OpenFileDialog(a.ctx, option)
	if err != nil {
		wails.LogErrorf(a.ctx, err.Error())
		return "", err
	}

	return path, nil
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
