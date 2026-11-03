package task

// List 下载任务列表
type List struct {
	tasks []*Task
}

// Len 返回任务数量
func (l *List) Len() int {
	return len(l.tasks)
}

// Reset 清空列表
func (l *List) Reset() {
	l.tasks = nil
}

// Add 追加任务
func (l *List) Add(tasks ...*Task) {
	l.tasks = append(l.tasks, tasks...)
}

// All 返回全部任务信息
func (l *List) All() []TaskInfo {
	return l.infos(0, len(l.tasks))
}

// Page 返回第 page 页（从 0 开始）的任务信息
func (l *List) Page(page, size int) []TaskInfo {
	start := page * size
	if start < 0 || start >= len(l.tasks) {
		return []TaskInfo{}
	}
	return l.infos(start, min(start+size, len(l.tasks)))
}

// SetDeleted 设置任务删除状态，越界时忽略
func (l *List) SetDeleted(index int, deleted bool) {
	if t := l.at(index); t != nil {
		t.Deleted = deleted
	}
}

// SetMeta 修改任务的歌曲名和作者，越界时忽略
func (l *List) SetMeta(index int, songName, author string) {
	if t := l.at(index); t != nil {
		t.Meta.SongName = songName
		t.Meta.Author = author
	}
}

// Import 从任务信息恢复任务，未携带 sessdata 的任务使用传入的 sessdata
func (l *List) Import(infos []TaskInfo, sessdata string) {
	for _, info := range infos {
		if t, ok := fromInfo(info, sessdata); ok {
			l.Add(t)
		}
	}
}

func (l *List) at(index int) *Task {
	if index < 0 || index >= len(l.tasks) {
		return nil
	}
	return l.tasks[index]
}

func (l *List) infos(start, end int) []TaskInfo {
	infos := make([]TaskInfo, 0, end-start)
	for i := start; i < end; i++ {
		infos = append(infos, l.tasks[i].info(i))
	}
	return infos
}
