package config

import (
	"encoding/json"
	"os"
	"samples/internal/components/logger"
	"samples/internal/utils/merge"
)

type Config struct {
	ExternalConfigPath string     `json:"externalConfigPath"`
	Log                LogConfig  `json:"log"`
	Http               HttpConfig `json:"http"`
}

type LogConfig struct {
	DirName       string `json:"dirName"`
	DirMaxSizeMB  int    `json:"dirMaxSizeMB"`
	FileName      string `json:"fileName"`
	FileMaxSizeMB int    `json:"fileMaxSizeMB"`
}

type HttpConfig struct {
	Port           int `json:"port"`
	TimeoutSeconds int `json:"timeoutSeconds"`
}

var Cfg Config

// init 初始化配置
func init() {
	// 获取内部配置
	internalConfig := getConfig("config.json")
	if internalConfig.ExternalConfigPath == "" {
		Cfg = internalConfig
		return
	}
	// 获取外部配置
	externalConfig := getConfig(internalConfig.ExternalConfigPath)
	// 合并配置
	Cfg = merge.MergeNonEmpty(internalConfig, externalConfig)
}

// getConfig 获取配置
func getConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	data, _ = json.Marshal(config)
	logger.Debug("%s配置: %+v", path, string(data))
	return config
}
