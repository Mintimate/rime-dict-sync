package main

import (
	"crypto/md5"
	"flag"
	"fmt"
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

	// 收集所有需要的远程仓库
	remoteRepos := make(map[string]bool)
	
	// 添加全局远程仓库
	if config.REMOTE_REPO != "" {
		remoteRepos[config.REMOTE_REPO] = true
	}
	
	// 添加字典特定的远程仓库
	for _, dict := range config.TARGET_DICT {
		if dict.RemoteRepo != "" {
			remoteRepos[dict.RemoteRepo] = true
		}
	}

	// 如果有远程仓库配置，则进行比较
	if len(remoteRepos) > 0 {
		println("正在检查文件变化...")
		
		// 克隆所有需要的远程仓库
		remoteRepoDirs := make(map[string]string)
		for repoURL := range remoteRepos {
			repoDir := "temp_remote_repo_" + generateRepoId(repoURL)
			if err := cloneRemoteRepo(repoURL, repoDir); err != nil {
				println("警告: 无法克隆远程仓库", repoURL, "，默认进行更新:", err.Error())
				println("检测到文件变化，需要更新")
				return
			}
			remoteRepoDirs[repoURL] = repoDir
			defer os.RemoveAll(repoDir) // 清理临时目录
		}

		// 检查是否有变化
		hasChanges, err := hasAnyChanges(config, downloadDir, remoteRepoDirs)
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

// 生成仓库URL的唯一标识符
func generateRepoId(repoURL string) string {
	hash := md5.Sum([]byte(repoURL))
	return fmt.Sprintf("%x", hash)[:8]
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
