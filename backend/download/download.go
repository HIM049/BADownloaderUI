package download

import (
	"bili-audio-downloader/backend/config"
	"fmt"
)

type DownloadTask interface {
	Download() error
}

type Option struct {
	SongName     bool `json:"song_name"`
	SongCover    bool `json:"song_cover"`
	SongAuthor   bool `json:"song_author"`
	DownloadFlac bool `json:"download_flac"`
}

type Path struct {
	StreamPath   string
	CoverPath    string
	ConventPath  string
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

func (p *Path) GetOutputPath() string {
	return fmt.Sprintf("%s.%s%s", config.Cfg.GetDownloadPath(), p.OutputName, p.OutputFormat)
}
