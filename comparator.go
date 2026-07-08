package main

import (
	"os"
	"path/filepath"
	"strings"
)

// 提取文件正文，忽略不同来源生成的头部信息
func extractContentAfterSeparator(content []byte) string {
	contentStr := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.SplitAfter(contentStr, "\n")

	for index, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "..." || strings.Contains(trimmedLine, "此行之后不能写注释") {
			var builder strings.Builder
			for _, bodyLine := range lines[index+1:] {
				builder.WriteString(bodyLine)
			}
			return builder.String()
		}
	}

	if len(lines) > 20 {
		var builder strings.Builder
		for _, bodyLine := range lines[20:] {
			builder.WriteString(bodyLine)
		}
		return builder.String()
	}

	return contentStr
}

// 比较两个文件的内容（忽略头部信息）
func compareFileContent(localPath, remotePath string) (bool, error) {
	// 读取本地文件
	localContent, err := os.ReadFile(localPath)
	if err != nil {
		return false, err
	}

	// 读取远程文件
	remoteContent, err := os.ReadFile(remotePath)
	if err != nil {
		// 如果远程文件不存在，认为有变化
		return false, nil
	}

	// 提取两个文件 "..." 之后的内容
	localDataContent := extractContentAfterSeparator(localContent)
	remoteDataContent := extractContentAfterSeparator(remoteContent)

	// 比较内容是否相同
	return strings.TrimSpace(localDataContent) == strings.TrimSpace(remoteDataContent), nil
}

// 检查是否有任何文件发生变化
func hasAnyChanges(config *DictConfig, downloadDir string, remoteRepoDirs map[string]string) (bool, error) {
	for _, dict := range config.TARGET_DICT {
		localPath := filepath.Join(downloadDir, dict.Name)

		// 确定远程仓库目录和文件路径
		var remoteRepoDir string
		var remotePath string

		if dict.RemoteRepo != "" {
			// 使用字典特定的远程仓库
			remoteRepoDir = remoteRepoDirs[dict.RemoteRepo]
			if dict.RemotePath != "" {
				remotePath = filepath.Join(remoteRepoDir, dict.RemotePath)
			} else {
				remotePath = filepath.Join(remoteRepoDir, dict.Name)
			}
		} else if config.REMOTE_REPO != "" {
			// 使用全局远程仓库
			remoteRepoDir = remoteRepoDirs[config.REMOTE_REPO]
			remotePath = filepath.Join(remoteRepoDir, "dicts", dict.Name)
		} else {
			// 没有配置远程仓库，跳过比较
			println("跳过比较:", dict.Name, "(未配置远程仓库)")
			continue
		}

		isSame, err := compareFileContent(localPath, remotePath)
		if err != nil {
			return false, err
		}

		if !isSame {
			println("检测到变化:", dict.Name)
			return true, nil
		} else {
			println("无变化:", dict.Name)
		}
	}

	return false, nil
}
