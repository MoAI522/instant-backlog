package fileops

import (
	"fmt"
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

	// ID別の最新Issueを管理するマップ
	issueMap := make(map[int]*models.Issue)

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		issue, err := parser.ParseIssueFile(filePath)
		if err != nil {
			// エラーログを出力して続行することも可能
			fmt.Printf("警告: Issueファイルの解析に失敗しました %s: %v\n", file.Name(), err)
			continue
		}

		// 既に同じIDのIssueが存在する場合は、最新のステータスを持つほうを採用
		existingIssue, exists := issueMap[issue.ID]
		if !exists {
			issueMap[issue.ID] = issue
			fmt.Printf("Issue ID=%d を登録しました (Status=%s)\n", issue.ID, issue.Status)
		} else {
			// すでに同一IDが存在する場合は警告を出す
			fmt.Printf("警告: Issue ID=%d が重複しています。現在: Status=%s, 新規: Status=%s\n",
				issue.ID, existingIssue.Status, issue.Status)
			// ファイル名からCloseが優先されるようにすると、より安全
			if issue.Status == "Close" {
				issueMap[issue.ID] = issue
				fmt.Printf("更新: Issue ID=%d の最新ステータスとして %s を採用しました\n",
					issue.ID, issue.Status)
			}
		}
	}

	// マップから配列に変換
	var issues []*models.Issue
	for _, issue := range issueMap {
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
