package bilibili

type Video struct {
	Bvid string `json:"bvid"`
	Meta struct {
		Title      string `json:"title"`       // 视频标题
		Cover      string `json:"cover"`       // 封面
		Author     string `json:"author"`      // 作者
		LyricsPath string `json:"lyrics_path"` // 歌词
	}
	Up struct {
		Mid    int    `json:"mid"`    // UP MID
		Name   string `json:"name"`   // UP 昵称
		Avatar string `json:"avatar"` // UP 头像
	}
	Videos []Videos
}
type Videos struct {
	Cid  int    `json:"cid"`
	Part string `json:"part"` // 分集名称
	Meta struct {
		SongName string `json:"song_name"` // 歌名
	}
}
