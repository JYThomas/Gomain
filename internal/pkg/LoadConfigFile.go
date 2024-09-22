package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// 从当前工作目录向上查找直到找到Gomain目录
func FindGomainDir(startPath string) (string, error) {
	for {
		gomainDir := filepath.Join(startPath, "Gomain")

		if _, err := os.Stat(gomainDir); !os.IsNotExist(err) {
			return gomainDir, nil
		}

		// 获取父目录
		parentDir := filepath.Dir(startPath)
		if parentDir == startPath { // 达到根目录
			return "", errors.New("Error: No Gomain Dir")
		}

		startPath = parentDir
	}
}

// 加载配置文件
func LoadConfig() (*ini.File, error) {
	// 获取当前工作目录
	startPath, err := filepath.Abs(".")
	if err != nil {
		return nil, errors.New("Load Config File Error: fail to find current File")
	}

	// 查找Gomain目录
	gomainDir, err := FindGomainDir(startPath)
	if err != nil {
		return nil, fmt.Errorf("Load Config File Error: Fail to find Gomain dir: %v", err)
	}

	// 配置文件相对路径
	configRelPath := "internal/Config.ini"
	configPath := filepath.Join(gomainDir, configRelPath)

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Load Config File Error: Config File Not Exist: %s", configPath)
	}

	// 加载配置文件
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("Load Config File Error: %v", err)
	}

	return cfg, nil
}
