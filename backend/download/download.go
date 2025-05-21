package download

import (
	"bili-audio-downloader/backend/config"
	"bili-audio-downloader/backend/constants"
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

type DownloadTask interface {
	SetID(int)
	Download() error
	ConventFormat() error
	WriteMetadata() error
	ExportFile() error
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

type FileName struct {
	Title    string
	Subtitle string
	Quality  string
	ID       int
	Format   string
}

func ExportFile(title, subtitle, outputFormat string, listid int, currentPath string) error {
	quality := "normal"
	if outputFormat == constants.AudioType.Flac {
		quality = "hires"

	}
	// 处理模板和生成文件名
	fileName := FileName{
		Title:    title,
		Subtitle: subtitle,
		Quality:  quality,
		ID:       listid,
		Format:   outputFormat,
	}
	tmpl, err := template.New("filename").Parse(config.Cfg.FileConfig.FileNameTemplate)
	if err != nil {
		return err
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, fileName)
	if err != nil {
		return err
	}

	// 重命名歌曲文件并移动位置
	err = os.Rename(currentPath, filepath.Join(config.Cfg.GetDownloadPath(), output.String()))
	if err != nil {
		return err
	}
	return nil
}
