package main

type Config struct {
	DownloadPath string `json:"download_path"`
	CachePath    string `json:"cache_path"`
	// TODO: 改为使用 ID 作为列表文件名
	VideoListPath   string `json:"videolist_path"`
	DownloadThreads int    `json:"download_threads"`
	RetryCount      int    `json:"retry_count"`
	ConvertFormat   bool   `json:"convert_format"`
	DeleteCache     bool   `json:"delete_cache"`
	Account         Account
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
		DownloadPath:    "./Download",
		CachePath:       "./Cache",
		VideoListPath:   "./Cache/video_list.json",
		DownloadThreads: 5,
		RetryCount:      10,
		ConvertFormat:   Checkffmpeg(),
		DeleteCache:     true,
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

// 保存设置到文件
func (cfg *Config) Save() error {
	err := SaveJsonFile("./config.json", cfg)
	if err != nil {
		return err
	}
	return nil
}
