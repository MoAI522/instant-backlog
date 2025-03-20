package fileops

import (
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
	return os.WriteFile(filePath, mdContent, 0644)
}

// WriteEpic - 指定されたEpicをマークダウンファイルに書き込む
func WriteEpic(directory string, epic *models.Epic) error {
	// マークダウンを生成
	mdContent, err := parser.GenerateMarkdown(epic, epic.Content)
	if err != nil {
		return err
	}

	// ファイル名を生成
	filename := utils.GenerateFilename(epic.ID, epic.Status, epic.Title)
	filePath := filepath.Join(directory, filename)

	// ファイルに書き込み
	return os.WriteFile(filePath, mdContent, 0644)
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
