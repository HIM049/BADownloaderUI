package adapter

type Option struct {
	SongName     bool `json:"song_name"`
	SongCover    bool `json:"song_cover"`
	SongAuthor   bool `json:"song_author"`
	DownloadFlac bool `json:"download_flac"`
}

type Path struct {
	StreamPath   string
	CoverPath    string
	CurrentPath  string
	OutputName   string
	OutputFormat string
}

type MetaData struct {
	Title     string
	PageTitle string
	PartId    int
	SongName  string
	Author    string
	LyricsUrl string
}

type TaskInfo struct {
	Index      int
	SongName   string
	SongAuthor string
	CoverUrl   string
}
