package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func UpdateConfig(filePath string) error {
	// 打开 JSON 文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 将 JSON 文件反序列化为 map
	var configMap map[string]interface{}
	err = json.Unmarshal(data, &configMap)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	config := new(Config)
	config.init()

	// 匹配结构体字段
	config.ConfigVersion = CONFIG_VERSION

	if v, ok := configMap["delete_cache"].(bool); ok {
		config.DeleteCache = v
	}

	if v, ok := configMap["theme"].(string); ok {
		config.Theme = v
	}

	if downloadConfig, ok := configMap["download_config"].(map[string]interface{}); ok {
		if v, ok := downloadConfig["download_threads"].(float64); ok {
			config.DownloadConfig.DownloadThreads = int(v)
		}
		if v, ok := downloadConfig["retry_count"].(float64); ok {
			config.DownloadConfig.RetryCount = int(v)
		}
	}

	if fileConfig, ok := configMap["file_config"].(map[string]interface{}); ok {
		if v, ok := fileConfig["convert_format"].(bool); ok {
			config.FileConfig.ConvertFormat = v
		}
		if v, ok := fileConfig["file_name_template"].(string); ok {
			config.FileConfig.FileNameTemplate = v
		}
		if v, ok := fileConfig["download_path"].(string); ok {
			config.FileConfig.DownloadPath = v
		}
		if v, ok := fileConfig["cache_path"].(string); ok {
			config.FileConfig.CachePath = v
		}
		if v, ok := fileConfig["videolist_path"].(string); ok {
			config.FileConfig.VideoListPath = v
		}
	}

	if account, ok := configMap["account"].(map[string]interface{}); ok {
		if v, ok := account["is_login"].(bool); ok {
			config.Account.IsLogin = v
		}
		if v, ok := account["use_account"].(bool); ok {
			config.Account.UseAccount = v
		}
		if v, ok := account["sessdata"].(string); ok {
			config.Account.SESSDATA = v
		}
		if v, ok := account["bili_jct"].(string); ok {
			config.Account.Bili_jct = v
		}
		if v, ok := account["dede_user_id"].(string); ok {
			config.Account.DedeUserID = v
		}
		if v, ok := account["dede_user_id__ck_md5"].(string); ok {
			config.Account.DedeUserID__ckMd5 = v
		}
		if v, ok := account["sid"].(string); ok {
			config.Account.Sid = v
		}
	}

	// 保存设置
	err = config.Save()
	if err != nil {
		return err
	}

	return nil
}

// 初始化设置
func (cfg *Config) init() {
	*cfg = Config{
		ConfigVersion: CONFIG_VERSION,
		DeleteCache:   true,
		Theme:         "lightPink",
		DownloadConfig: DownloadConfig{
			DownloadThreads: 5,
			RetryCount:      10,
		},
		FileConfig: FileConfig{
			ConvertFormat:    Checkffmpeg(),
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
}

// 读取设置文件
func (cfg *Config) Get() error {
	for {
		// 判断设置文件是否已经存在
		if !IsFileExists("./config.json") {
			// 文件不存在
			file := new(Config)
			file.init()
			err := file.Save()
			if err != nil {
				return err
			}
		} else {
			// 文件已存在
			var file Config
			err := LoadJsonFile("./config.json", &file)
			if err != nil {
				return err
			}
			if file.ConfigVersion == CONFIG_VERSION {
				*cfg = file
				return nil
			}
			if file.ConfigVersion < CONFIG_VERSION {
				err := UpdateConfig("./config.json")
				if err != nil {
					return err
				} else {
					continue
				}
			}
			if file.ConfigVersion > CONFIG_VERSION {
				cfg.init()
				err := cfg.Save()
				if err != nil {
					return err
				}
			}
		}
	}
}

func (cfg *Config) Check() {
	if cfg.ConfigVersion != CONFIG_VERSION {
		cfg.ConfigVersion = CONFIG_VERSION
		cfg.Save()
	}
}

// 保存设置到文件
func (cfg *Config) Save() error {
	err := SaveJsonFile("./config.json", cfg)
	if err != nil {
		return err
	}
	return nil
}
