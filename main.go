package main

import (
	"flag"
)

func main() {
	// 定义命令行参数
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "配置文件路径")
	flag.StringVar(&configPath, "config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	config, err := loadConfig(configPath)
	if err != nil {
		panic(err)
	}

	// 获取下载目录（默认为"dl_dicts"）
	downloadDir := "dl_dicts"
	if config.DOWNLOAD_DIR != "" {
		downloadDir = config.DOWNLOAD_DIR
	}

	// 下载并处理每个字典文件
	for _, dict := range config.TARGET_DICT {
		println("正在下载:", dict.Name, "...")
		if err := downloadAndModify(dict, downloadDir); err != nil {
			panic(err)
		}
	}
}
