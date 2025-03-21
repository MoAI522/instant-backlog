// test/test_utils.go
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/fileops"
	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
)

// テスト用の一時ディレクトリを作成
func setupTestEnvironment(t *testing.T) (*config.Config, func()) {
	t.Helper()

	// 一時ディレクトリを作成
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("instant-backlog-test-%d", os.Getpid()))
	epicDir := filepath.Join(tempDir, "projects", "epic")
	issuesDir := filepath.Join(tempDir, "projects", "issues")

	// ディレクトリ作成
	if err := os.MkdirAll(epicDir, 0755); err != nil {
		t.Fatalf("テスト環境のセットアップに失敗しました: %v", err)
	}
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		t.Fatalf("テスト環境のセットアップに失敗しました: %v", err)
	}

	// テスト用の設定を作成
	cfg := &config.Config{
		ProjectsDir: filepath.Join(tempDir, "projects"),
		EpicDir:     epicDir,
		IssuesDir:   issuesDir,
		OrderCSV:    filepath.Join(tempDir, "projects", "order.csv"),
	}

	// クリーンアップ関数
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return cfg, cleanup
}

// テスト用のEpicファイルを作成
func createTestEpic(t *testing.T, cfg *config.Config, id int, title, status string) {
	t.Helper()

	epic := &models.Epic{
		ID:      id,
		Title:   title,
		Status:  status,
		Content: fmt.Sprintf("これはテスト用のエピック %d です。", id),
	}

	if err := fileops.WriteEpic(cfg.EpicDir, epic); err != nil {
		t.Fatalf("テスト用Epicの作成に失敗しました: %v", err)
	}
}

// テスト用のIssueファイルを作成
func createTestIssue(t *testing.T, cfg *config.Config, id int, title, status string, epicID, estimate int) {
	t.Helper()

	issue := &models.Issue{
		ID:       id,
		Title:    title,
		Status:   status,
		Epic:     epicID,
		Estimate: estimate,
		Content:  fmt.Sprintf("これはテスト用のタスク %d です。", id),
	}

	if err := fileops.WriteIssue(cfg.IssuesDir, issue); err != nil {
		t.Fatalf("テスト用Issueの作成に失敗しました: %v", err)
	}
}

// テスト用CSVファイルを作成
func createTestOrderCSV(t *testing.T, cfg *config.Config, items []models.OrderCSVItem) {
	t.Helper()

	if err := parser.WriteOrderCSV(cfg.OrderCSV, items); err != nil {
		t.Fatalf("テスト用CSVの作成に失敗しました: %v", err)
	}
}

// ファイルの存在を確認
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ディレクトリの存在を確認するヘルパー関数
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
