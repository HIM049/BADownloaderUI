package config

import (
	"bili-audio-downloader/constants"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func migrateConfig(filePath string) error {
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
	config := DefaultConfig()

	// 匹配结构体字段
	config.ConfigVersion = constants.CONFIG_VERSION

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

	config.UpdateAndSave()
	return nil
}
