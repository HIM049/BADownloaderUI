package main

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
	ConvertFormat bool   `json:"convert_format"`
	DownloadPath  string `json:"download_path"`
	CachePath     string `json:"cache_path"`
	VideoListPath string `json:"videolist_path"`
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

// 初始化设置
func (cfg *Config) init() {
	*cfg = Config{
		ConfigVersion: 1,
		DeleteCache:   true,
		Theme:         "lightPink",
		DownloadConfig: DownloadConfig{
			DownloadThreads: 5,
			RetryCount:      10,
		},
		FileConfig: FileConfig{
			ConvertFormat: Checkffmpeg(),
			DownloadPath:  "./Download",
			CachePath:     "./Cache",
			VideoListPath: "./Cache/video_list.json",
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
			*cfg = file
			return nil
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
