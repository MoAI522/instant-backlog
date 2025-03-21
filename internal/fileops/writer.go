package fileops

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
	"github.com/moai/instant-backlog/pkg/utils"
)

// WriteIssue - 指定されたIssueをマークダウンファイルに書き込む
func WriteIssue(directory string, issue *models.Issue) error {
	// マークダウンを生成
	mdContent, err := parser.GenerateMarkdown(issue, issue.Content)
	if err != nil {
		return err
	}

	// ファイル名を生成
	filename := utils.GenerateFilename(issue.ID, issue.Status, issue.Title)
	filePath := filepath.Join(directory, filename)

	// ファイルに書き込み
	err = os.WriteFile(filePath, mdContent, 0644)
	fmt.Printf("WriteEpic: ファイル書き込み結果=%v\n", err)
	return err
}

// WriteEpic - 指定されたEpicをマークダウンファイルに書き込む
func WriteEpic(directory string, epic *models.Epic) error {
	fmt.Printf("WriteEpic: Epic ID=%d, Title=%s, Status=%s をファイルに書き込みます\n", epic.ID, epic.Title, epic.Status)
	// マークダウンを生成
	mdContent, err := parser.GenerateMarkdown(epic, epic.Content)
	if err != nil {
		return err
	}

	// ファイル名を生成
	filename := utils.GenerateFilename(epic.ID, epic.Status, epic.Title)
	filePath := filepath.Join(directory, filename)
	fmt.Printf("WriteEpic: 生成したファイル名=%s, ファイルパス=%s\n", filename, filePath)

	// ファイルに書き込み
	err = os.WriteFile(filePath, mdContent, 0644)
	fmt.Printf("WriteEpic: ファイル書き込み結果=%v\n", err)
	return err
}

// RenameFile - ファイル名を変更する（新しいファイル名が必要な場合）
func RenameFile(directory, oldFilename, newFilename string) error {
	if oldFilename == newFilename {
		return nil // 変更不要
	}

	oldPath := filepath.Join(directory, oldFilename)
	newPath := filepath.Join(directory, newFilename)

	return os.Rename(oldPath, newPath)
}
