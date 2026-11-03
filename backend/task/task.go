package task

import "strconv"

// Kind 任务来源类型，数值与导出列表中的 type 字段一致
type Kind int

const (
	KindVideo Kind = 1
	KindAudio Kind = 2
)

// MetaData 文件名与音频标签使用的信息
type MetaData struct {
	Title     string
	PageTitle string
	PartId    int
	SongName  string
	Author    string
	LyricsUrl string
}

// Task 一个下载任务，视频分 P 或音频
type Task struct {
	Kind     Kind
	Bvid     string
	Cid      int
	Auid     string
	CoverURL string
	Sessdata string
	Meta     MetaData
	Deleted  bool
}

// TaskInfo 提供给前端和导出列表的任务信息
type TaskInfo struct {
	Index      int
	SongName   string
	SongAuthor string
	CoverUrl   string
	IsDelete   bool

	Type     Kind     `json:"type,omitempty"`
	Bvid     string   `json:"bvid,omitempty"`
	Cid      int      `json:"cid,omitempty"`
	Auid     string   `json:"auid,omitempty"`
	Sessdata string   `json:"sessdata,omitempty"`
	MetaData MetaData `json:"meta_data,omitempty"`
}

// cacheID 返回任务在缓存文件名中使用的标识
func (t *Task) cacheID() string {
	if t.Kind == KindAudio {
		return "au" + t.Auid
	}
	return strconv.Itoa(t.Cid)
}

// info 生成任务信息
func (t *Task) info(index int) TaskInfo {
	return TaskInfo{
		Index:      index,
		SongName:   t.Meta.SongName,
		SongAuthor: t.Meta.Author,
		CoverUrl:   t.CoverURL,
		IsDelete:   t.Deleted,
		Type:       t.Kind,
		Bvid:       t.Bvid,
		Cid:        t.Cid,
		Auid:       t.Auid,
		Sessdata:   t.Sessdata,
		MetaData:   t.Meta,
	}
}

// fromInfo 从任务信息恢复任务，未知类型返回 false
func fromInfo(info TaskInfo, sessdata string) (*Task, bool) {
	if info.Type != KindVideo && info.Type != KindAudio {
		return nil, false
	}
	if info.Sessdata != "" {
		sessdata = info.Sessdata
	}
	return &Task{
		Kind:     info.Type,
		Bvid:     info.Bvid,
		Cid:      info.Cid,
		Auid:     info.Auid,
		CoverURL: info.CoverUrl,
		Sessdata: sessdata,
		Meta:     info.MetaData,
		Deleted:  info.IsDelete,
	}, true
}
