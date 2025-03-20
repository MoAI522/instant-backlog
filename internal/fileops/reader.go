package fileops

import (
	"os"
	"path/filepath"

	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
)

// ReadAllIssues - 指定ディレクトリからすべてのIssueを読み込む
func ReadAllIssues(directory string) ([]*models.Issue, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var issues []*models.Issue
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		issue, err := parser.ParseIssueFile(filePath)
		if err != nil {
			// エラーログを出力して続行することも可能
			continue
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// ReadAllEpics - 指定ディレクトリからすべてのEpicを読み込む
func ReadAllEpics(directory string) ([]*models.Epic, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var epics []*models.Epic
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		epic, err := parser.ParseEpicFile(filePath)
		if err != nil {
			// エラーログを出力して続行することも可能
			continue
		}

		epics = append(epics, epic)
	}

	return epics, nil
}
