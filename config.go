package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type DictTarget struct {
	Name       string `yaml:"name"`
	URL        string `yaml:"url"`
	RemoteRepo string `yaml:"remote_repo,omitempty"`
	RemotePath string `yaml:"remote_path,omitempty"`
	MergeMode  string `yaml:"merge_mode,omitempty"`
}

type DictConfig struct {
	TARGET_DICT  []DictTarget `yaml:"TARGET_DICT"`
	DOWNLOAD_DIR string       `yaml:"DOWNLOAD_DIR"`
	REMOTE_REPO  string       `yaml:"REMOTE_REPO,omitempty"`
}

func loadConfig(path string) (*DictConfig, error) {
	// 读取配置文件
	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config DictConfig
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}
	if len(config.TARGET_DICT) == 0 {
		return nil, fmt.Errorf("TARGET_DICT 不能为空")
	}
	seen := make(map[string]bool)
	for _, dict := range config.TARGET_DICT {
		if strings.TrimSpace(dict.Name) == "" || dict.Name == "." || dict.Name == ".." || filepath.Base(dict.Name) != dict.Name || strings.ContainsAny(dict.Name, `/\\`) {
			return nil, fmt.Errorf("词库 name 必须是文件名: %q", dict.Name)
		}
		if seen[dict.Name] {
			return nil, fmt.Errorf("词库 name 重复: %s", dict.Name)
		}
		seen[dict.Name] = true
		if strings.TrimSpace(dict.URL) == "" {
			return nil, fmt.Errorf("词库 %s 缺少 url", dict.Name)
		}
	}

	return &config, nil
}
