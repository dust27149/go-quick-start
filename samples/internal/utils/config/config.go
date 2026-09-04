package config

import (
	"encoding/json"
	"net"
	"os"
	"samples/internal/components/logger"
	"samples/internal/utils/merge"
	"strings"
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
	Port           int      `json:"port"`
	TimeoutSeconds int      `json:"timeoutSeconds"`
	IpWhiteList    []string `json:"ipWhiteList"`
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

// IsAllowedIP 判断请求的 IP 是否在允许访问的范围内
func IsAllowedIP(remoteAddr string) bool {
	remoteIP := remoteAddr
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		remoteIP = host
	}
	remoteIP = strings.TrimSpace(remoteIP)
	if strings.EqualFold(remoteIP, "localhost") {
		remoteIP = "127.0.0.1"
	}

	ip := net.ParseIP(remoteIP)
	if ip == nil {
		return false
	}

	for _, allowedIP := range Cfg.Http.IpWhiteList {
		allowedIP = strings.TrimSpace(allowedIP)
		if strings.EqualFold(allowedIP, "localhost") {
			allowedIP = "127.0.0.1"
		}

		if strings.Contains(allowedIP, "/") {
			_, network, err := net.ParseCIDR(allowedIP)
			if err == nil && network.Contains(ip) {
				return true
			}
			continue
		}

		if allowed := net.ParseIP(allowedIP); allowed != nil && allowed.Equal(ip) {
			return true
		}
	}
	return false
}
