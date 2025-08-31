package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 下载字典文件
func downloadDict(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// 修改字典文件内容
func modifyDictContent(content []byte, dictName string) ([]byte, error) {
	contentStr := string(content)
	filename := filepath.Base(dictName)
	newName := strings.TrimSuffix(filename, ".dict.yaml")

	currentDate := time.Now().Format("2006-01-02")

	endPos := strings.Index(contentStr, "...")
	if endPos == -1 {
		headerEnd := 0
		lineCount := 0
		for i, c := range contentStr {
			if c == '\n' {
				lineCount++
				if lineCount >= 20 {
					headerEnd = i
					break
				}
			}
		}
		if headerEnd > 0 {
			endPos = headerEnd
		} else {
			endPos = len(contentStr)
		}
	}

	modifiedHeader := ""
	foundName := false
	foundVersion := false

	for _, line := range strings.Split(contentStr[:endPos], "\n") {
		trimmedLine := strings.TrimSpace(line)

		if !foundName && strings.HasPrefix(trimmedLine, "name:") {
			line = "name: " + newName
			foundName = true
		}

		if !foundVersion && strings.HasPrefix(trimmedLine, "version:") {
			line = "version: \"" + currentDate + "\""
			foundVersion = true
		}

		modifiedHeader += line + "\n"
	}

	return []byte(modifiedHeader + contentStr[endPos:]), nil
}

func downloadAndModify(dict struct {
	Name       string `yaml:"name"`
	URL        string `yaml:"url"`
	RemoteRepo string `yaml:"remote_repo,omitempty"`
	RemotePath string `yaml:"remote_path,omitempty"`
}, downloadDir string) error {
	// 下载文件
	content, err := downloadDict(dict.URL)
	if err != nil {
		return err
	}

	// 修改内容
	modifiedContent, err := modifyDictContent(content, dict.Name)
	if err != nil {
		return err
	}

	// 创建下载目录（如果不存在）
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		if err := os.Mkdir(downloadDir, 0755); err != nil {
			return err
		}
	}

	// 保存文件
	filename := filepath.Base(dict.Name)
	filePath := filepath.Join(downloadDir, filename)
	if err := os.WriteFile(filePath, modifiedContent, 0644); err != nil {
		return err
	}
	println("已下载并修改:", filePath)
	return nil
}
