package main

import (
	"os"
	"strings"
)

// 提取文件内容中 "..." 之后的部分
func extractContentAfterSeparator(content []byte) string {
	contentStr := string(content)
	
	// 查找 "..." 分隔符
	separatorPos := strings.Index(contentStr, "...")
	if separatorPos == -1 {
		// 如果没有找到 "..."，则查找前20行之后的内容
		lineCount := 0
		for i, c := range contentStr {
			if c == '\n' {
				lineCount++
				if lineCount >= 20 {
					return contentStr[i+1:]
				}
			}
		}
		// 如果行数不足20行，返回原内容
		return contentStr
	}
	
	// 返回 "..." 之后的内容
	return contentStr[separatorPos+3:]
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
func hasAnyChanges(config *DictConfig, downloadDir, remoteRepoDir string) (bool, error) {
	for _, dict := range config.TARGET_DICT {
		localPath := downloadDir + "/" + dict.Name
		remotePath := remoteRepoDir + "/dicts/" + dict.Name
		
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
