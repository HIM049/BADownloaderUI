package config

import (
	"bili-audio-downloader/backend/constants"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// configPath 配置文件位置
const configPath = "./config.json"

type Config struct {
	ConfigVersion  int            `json:"config_version"`
	DeleteCache    bool           `json:"delete_cache"`
	Theme          string         `json:"theme"`
	DownloadConfig DownloadConfig `json:"download_config"`
	FileConfig     FileConfig     `json:"file_config"`
	Account        Account        `json:"account"`
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

// InitConfig 启动时读取配置，失败时退出
func InitConfig(ctx context.Context) {
	if err := LoadConfig(ctx); err != nil {
		wails.LogFatalf(ctx, "Failed to load config: %v", err)
	}
}

// LoadConfig 读取配置文件到 Cfg
// 文件不存在时写入默认配置；内容无法解析时备份为 .bak 并使用默认配置；旧版本配置会升级版本号并保存
func LoadConfig(ctx context.Context) error {
	cfg := DefaultConfig()
	needSave := false

	data, err := os.ReadFile(configPath)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		wails.LogWarning(ctx, "Config file not found, create a new one")
		needSave = true
	case err != nil:
		return fmt.Errorf("failed to read config file: %w", err)
	default:
		cfg.ConfigVersion = 0
		if err := json.Unmarshal(data, cfg); err != nil {
			wails.LogErrorf(ctx, "Invalid config file, backup to %s.bak and use default: %v", configPath, err)
			if err := os.Rename(configPath, configPath+".bak"); err != nil {
				return fmt.Errorf("failed to backup config file: %w", err)
			}
			cfg = DefaultConfig()
			needSave = true
		}
	}

	switch {
	case cfg.ConfigVersion < constants.CONFIG_VERSION:
		wails.LogInfof(ctx, "Upgrade config version from %d to %d", cfg.ConfigVersion, constants.CONFIG_VERSION)
		cfg.ConfigVersion = constants.CONFIG_VERSION
		needSave = true
	case cfg.ConfigVersion > constants.CONFIG_VERSION:
		wails.LogInfo(ctx, "Config version is higher than current version")
	}

	if needSave {
		return cfg.UpdateAndSave()
	}
	cfg.normalize()
	Cfg = *cfg
	return nil
}

// UpdateAndSave 修正非法值后写入配置文件，并替换当前使用的 Cfg
func (cfg *Config) UpdateAndSave() error {
	cfg.normalize()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp := configPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmp, configPath); err != nil {
		return err
	}

	Cfg = *cfg
	return nil
}

// normalize 将空值或非法值替换为默认值
func (cfg *Config) normalize() {
	def := DefaultConfig()
	fc := &cfg.FileConfig

	if _, err := template.New("filename").Parse(fc.FileNameTemplate); strings.TrimSpace(fc.FileNameTemplate) == "" || err != nil {
		fc.FileNameTemplate = def.FileConfig.FileNameTemplate
	}
	if strings.TrimSpace(fc.DownloadPath) == "" {
		fc.DownloadPath = def.FileConfig.DownloadPath
	}
	if strings.TrimSpace(fc.CachePath) == "" {
		fc.CachePath = def.FileConfig.CachePath
	}
	cfg.DownloadConfig.DownloadThreads = max(cfg.DownloadConfig.DownloadThreads, 0)
	cfg.DownloadConfig.RetryCount = max(cfg.DownloadConfig.RetryCount, 0)
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		ConfigVersion: constants.CONFIG_VERSION,
		DeleteCache:   true,
		Theme:         "lightBlue",
		DownloadConfig: DownloadConfig{
			DownloadThreads: 5,
			RetryCount:      10,
		},
		FileConfig: FileConfig{
			ConvertFormat:    false,
			FileNameTemplate: "{{.ID}}_{{.Title}}({{.Subtitle}})_{{.Quality}}{{.Format}}",
			DownloadPath:     "./Download",
			CachePath:        "./Cache",
		},
	}
}

// GetDownloadPath 返回下载目录的绝对路径
func (cfg *Config) GetDownloadPath() string {
	return absOrRaw(cfg.FileConfig.DownloadPath)
}

// GetCachePath 返回缓存目录的绝对路径
func (cfg *Config) GetCachePath() string {
	return absOrRaw(cfg.FileConfig.CachePath)
}

// absOrRaw 返回绝对路径，转换失败时返回原路径
func absOrRaw(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}
	return path
}
