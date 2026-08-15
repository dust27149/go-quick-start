package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Log  LogConfig  `json:"log"`
	Http HttpConfig `json:"http"`
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
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &Cfg); err != nil {
		panic(err)
	}
}
