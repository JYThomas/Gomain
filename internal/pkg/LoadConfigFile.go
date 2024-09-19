package pkg

import (
	"errors"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// 加载配置文件
func LoadConfig() (*ini.File, error) {
	// 获取当前文件夹所在路径
	basePath, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return nil, errors.New("Failed to determine current directory")
	}

	// 拼接配置文件路径
	configPath := filepath.Join(basePath, "../../../Config.ini")

	// 加载配置文件
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, errors.New("Module RAPIDDNS: Config Load Error")
	}
	return cfg, err
}
