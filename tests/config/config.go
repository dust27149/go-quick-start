package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Log LogConfig `json:"log"`
}

type LogConfig struct {
	DirName  string `json:"dirName"`
	FileName string `json:"fileName"`
}

var Cfg Config

func init() {

	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &Cfg); err != nil {
		panic(err)
	}

	cfgJSON, _ := json.Marshal(Cfg)
	fmt.Printf("配置加载成功: %+v\n", string(cfgJSON))
}
