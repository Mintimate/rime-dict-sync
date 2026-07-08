package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const mergeModeKeepRemoteHeader = "keep_remote_header"

// 下载字典文件
func downloadDict(url string) ([]byte, error) {
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("download %s failed: %s", url, resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func dictNameFromPath(dictName string) string {
	filename := filepath.Base(dictName)
	return strings.TrimSuffix(filename, ".dict.yaml")
}

func splitAfterLine(content string, match func(string) bool) (string, string, bool) {
	lines := strings.SplitAfter(content, "\n")
	var header strings.Builder
	var body strings.Builder
	found := false

	for _, line := range lines {
		if found {
			body.WriteString(line)
			continue
		}

		header.WriteString(line)
		if match(strings.TrimSpace(line)) {
			found = true
		}
	}

	return header.String(), body.String(), found
}

// 修改字典文件内容
func modifyDictContent(content []byte, dictName string) ([]byte, error) {
	contentStr := string(content)
	newName := dictNameFromPath(dictName)

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

	var modifiedHeader strings.Builder
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

		modifiedHeader.WriteString(line)
		modifiedHeader.WriteByte('\n')
	}

	return []byte(modifiedHeader.String() + contentStr[endPos:]), nil
}

func mergeWithRemoteHeader(upstreamContent, remoteContent []byte, dictName string) ([]byte, error) {
	upstreamStr := string(upstreamContent)
	remoteStr := string(remoteContent)

	if strings.Contains(upstreamStr, "\n---\n") || strings.HasPrefix(upstreamStr, "---\n") {
		remotePreamble, _, remoteFound := splitAfterLine(remoteStr, func(line string) bool {
			return line == "---"
		})
		if !remoteFound {
			return nil, fmt.Errorf("remote yaml header separator not found for %s", dictName)
		}

		_, upstreamAfterStart, upstreamFound := splitAfterLine(upstreamStr, func(line string) bool {
			return line == "---"
		})
		if !upstreamFound {
			return nil, fmt.Errorf("upstream yaml header separator not found for %s", dictName)
		}

		upstreamMetadata, upstreamBody, upstreamEndFound := splitAfterLine(upstreamAfterStart, func(line string) bool {
			return line == "..."
		})
		if !upstreamEndFound {
			return nil, fmt.Errorf("upstream yaml body separator not found for %s", dictName)
		}

		var metadata strings.Builder
		for _, line := range strings.SplitAfter(upstreamMetadata, "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "name:") {
				metadata.WriteString("name: ")
				metadata.WriteString(dictNameFromPath(dictName))
				metadata.WriteByte('\n')
			} else {
				metadata.WriteString(line)
			}
		}

		return []byte(remotePreamble + metadata.String() + upstreamBody), nil
	}

	remoteHeader, _, remoteFound := splitAfterLine(remoteStr, func(line string) bool {
		return strings.Contains(line, "此行之后不能写注释")
	})
	if !remoteFound {
		return nil, fmt.Errorf("remote table header marker not found for %s", dictName)
	}

	_, upstreamBody, upstreamFound := splitAfterLine(upstreamStr, func(line string) bool {
		return strings.Contains(line, "此行之后不能写注释")
	})
	if !upstreamFound {
		return nil, fmt.Errorf("upstream table header marker not found for %s", dictName)
	}

	return []byte(remoteHeader + upstreamBody), nil
}

func downloadAndModify(dict DictTarget, downloadDir string, remoteContent []byte) error {
	// 下载文件
	content, err := downloadDict(dict.URL)
	if err != nil {
		return err
	}

	// 修改内容
	var modifiedContent []byte
	if dict.MergeMode == mergeModeKeepRemoteHeader {
		if len(remoteContent) == 0 {
			return fmt.Errorf("merge mode %s requires remote content for %s", dict.MergeMode, dict.Name)
		}
		modifiedContent, err = mergeWithRemoteHeader(content, remoteContent, dict.Name)
	} else {
		modifiedContent, err = modifyDictContent(content, dict.Name)
	}
	if err != nil {
		return err
	}

	// 创建下载目录（如果不存在）
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
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
