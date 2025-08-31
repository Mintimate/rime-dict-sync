package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type DictConfig struct {
	TARGET_DICT []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"TARGET_DICT"`
	DOWNLOAD_DIR string `yaml:"DOWNLOAD_DIR"`
	REMOTE_REPO  string `yaml:"REMOTE_REPO"`
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

	return &config, nil
}
