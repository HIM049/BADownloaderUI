package config

import (
	"bili-audio-downloader/constants"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	ConfigVersion  int            `json:"config_version"`
	DeleteCache    bool           `json:"delete_cache"`
	Theme          string         `json:"theme"`
	DownloadConfig DownloadConfig `json:"download_config"`
	FileConfig     FileConfig     `json:"file_config"`
	Account        Account
}

type DownloadConfig struct {
	DownloadThreads int `json:"download_threads"`
	RetryCount      int `json:"retry_count"`
}

type FileConfig struct {
	ConvertFormat    bool   `json:"convert_format"`
	FileNameTemplate string `json:"file_name_template"`
	DownloadPath     string `json:"download_path"`
	CachePath        string `json:"cache_path"`
	VideoListPath    string `json:"videolist_path"`
}

type Account struct {
	IsLogin           bool   `json:"is_login"`
	UseAccount        bool   `json:"use_account"`
	SESSDATA          string `json:"sessdata"`
	Bili_jct          string `json:"bili_jct"`
	DedeUserID        string `json:"dede_user_id"`
	DedeUserID__ckMd5 string `json:"dede_user_id__ck_md5"`
	Sid               string `json:"sid"`
}

var Cfg Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found: ", err)
			newConfig := DefaultConfig()
			newConfig.UpdateAndSave()
			fmt.Println("Created a new config")
		} else {
			log.Fatalf("Failed to read config file: %v", err)
		}
	}

	// 初始化嵌套结构体
	downloadCfg := DownloadConfig{
		DownloadThreads: viper.GetInt("download_config.download_threads"),
		RetryCount:      viper.GetInt("download_config.retry_count"),
	}

	fileCfg := FileConfig{
		ConvertFormat:    viper.GetBool("file_config.convert_format"),
		FileNameTemplate: viper.GetString("file_config.file_name_template"),
		DownloadPath:     viper.GetString("file_config.download_path"),
		CachePath:        viper.GetString("file_config.cache_path"),
		VideoListPath:    viper.GetString("file_config.videolist_path"),
	}

	account := Account{
		IsLogin:           viper.GetBool("Account.is_login"),
		UseAccount:        viper.GetBool("Account.use_account"),
		SESSDATA:          viper.GetString("Account.sessdata"),
		Bili_jct:          viper.GetString("Account.bili_jct"),
		DedeUserID:        viper.GetString("Account.dede_user_id"),
		DedeUserID__ckMd5: viper.GetString("Account.dede_user_id__ck_md5"),
		Sid:               viper.GetString("Account.sid"),
	}

	// 初始化主配置结构体
	Cfg = Config{
		ConfigVersion: viper.GetInt("config_version"),
		DeleteCache:   viper.GetBool("delete_cache"),
		Theme:         viper.GetString("theme"),
		// Debug:        viper.GetBool("debug"),
		DownloadConfig: downloadCfg,
		FileConfig:     fileCfg,
		Account:        account,
	}

	// 检查配置文件版本
	if Cfg.ConfigVersion != constants.CONFIG_VERSION {
		if Cfg.ConfigVersion < constants.CONFIG_VERSION {
			err := migrateConfig("./config.json")
			if err != nil {
				log.Fatalf("Failed to migrate config: %v\n", err)
			}
		} else {
			fmt.Println("Config version is higher than current version")
		}
	}
}

func (cfg *Config) UpdateAndSave() error {
	viper.Set("config_version", cfg.ConfigVersion)
	viper.Set("delete_cache", cfg.DeleteCache)
	viper.Set("theme", cfg.Theme)

	viper.Set("download_config.download_threads", cfg.DownloadConfig.DownloadThreads)
	viper.Set("download_config.retry_count", cfg.DownloadConfig.RetryCount)

	viper.Set("file_config.convert_format", cfg.FileConfig.ConvertFormat)
	viper.Set("file_config.file_name_template", cfg.FileConfig.FileNameTemplate)
	viper.Set("file_config.download_path", cfg.FileConfig.DownloadPath)
	viper.Set("file_config.cache_path", cfg.FileConfig.CachePath)
	viper.Set("file_config.videolist_path", cfg.FileConfig.VideoListPath)

	viper.Set("Account.is_login", cfg.Account.IsLogin)
	viper.Set("Account.use_account", cfg.Account.UseAccount)
	viper.Set("Account.sessdata", cfg.Account.SESSDATA)
	viper.Set("Account.bili_jct", cfg.Account.Bili_jct)
	viper.Set("Account.dede_user_id", cfg.Account.DedeUserID)
	viper.Set("Account.dede_user_id__ck_md5", cfg.Account.DedeUserID__ckMd5)
	viper.Set("Account.sid", cfg.Account.Sid)

	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

// 初始化设置
func DefaultConfig() *Config {
	cfg := Config{
		ConfigVersion: constants.CONFIG_VERSION,
		DeleteCache:   true,
		Theme:         "lightPink",
		DownloadConfig: DownloadConfig{
			DownloadThreads: 5,
			RetryCount:      10,
		},
		FileConfig: FileConfig{
			ConvertFormat:    false, // TODO
			FileNameTemplate: "{{.ID}}_{{.Title}}({{.Subtitle}})_{{.Quality}}.{{.Format}}",
			DownloadPath:     "./Download",
			CachePath:        "./Cache",
			VideoListPath:    "./Cache/video_list.json",
		},
		Account: Account{
			IsLogin:           false,
			UseAccount:        false,
			SESSDATA:          "",
			Bili_jct:          "",
			DedeUserID:        "",
			DedeUserID__ckMd5: "",
			Sid:               "",
		},
	}
	return &cfg
}

func (cfg *Config) GetDownloadPath() string {
	path, err := filepath.Abs(cfg.FileConfig.DownloadPath)
	if err != nil {
		log.Fatalln("Failed to get abs path: ", err)
	}
	return path
}

func (cfg *Config) GetCachePath() string {
	path, err := filepath.Abs(cfg.FileConfig.CachePath)
	if err != nil {
		log.Fatalln("Failed to get abs path: ", err)
	}
	return path
}

func (cfg *Config) GetVideolistPath() string {
	path, err := filepath.Abs(cfg.FileConfig.VideoListPath)
	if err != nil {
		log.Fatalln("Failed to get abs path: ", err)
	}
	return path
}
