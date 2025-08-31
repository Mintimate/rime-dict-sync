package main

import (
	"flag"
	"os"
	"os/exec"
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

	// 如果配置了远程仓库，则进行比较
	if config.REMOTE_REPO != "" {
		println("正在检查文件变化...")
		
		// 克隆远程仓库到临时目录
		remoteRepoDir := "temp_remote_repo"
		if err := cloneRemoteRepo(config.REMOTE_REPO, remoteRepoDir); err != nil {
			println("警告: 无法克隆远程仓库，默认进行更新:", err.Error())
			println("检测到文件变化，需要更新")
			return
		}
		defer os.RemoveAll(remoteRepoDir) // 清理临时目录

		// 检查是否有变化
		hasChanges, err := hasAnyChanges(config, downloadDir, remoteRepoDir)
		if err != nil {
			panic(err)
		}

		if !hasChanges {
			println("所有文件无变化，跳过更新")
			os.Exit(1) // 返回非零退出码表示无需更新
		}

		println("检测到文件变化，需要更新")
	}
}

// 克隆远程仓库
func cloneRemoteRepo(repoURL, targetDir string) error {
	// 如果目录已存在，先删除
	if _, err := os.Stat(targetDir); err == nil {
		if err := os.RemoveAll(targetDir); err != nil {
			return err
		}
	}

	// 克隆仓库
	cmd := exec.Command("git", "clone", repoURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
