// test/issue_management_test.go
package test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/moai/instant-backlog/internal/parser"
)

/**
 * ユーザーストーリー1：マークダウンでIssueを管理できること
 *
 * ユーザーは課題（Issue）をマークダウン形式で作成・管理したい。
 * 各Issueには以下の情報が含まれる：
 * - ID: 一意の識別子
 * - タイトル: 課題の簡潔な説明
 * - ステータス: 課題の状態（Open/Close）
 * - エピックID: 関連するエピックの識別子
 * - 見積もり: 課題の作業量の見積もり
 * - 内容: 課題の詳細な説明
 *
 * マークダウンファイルは適切なメタデータを持ち、正しく解析できる必要がある。
 */
func TestManageIssuesWithMarkdown(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のEpicとIssueを作成
	createTestEpic(t, cfg, 1, "テストエピック", "Open")
	createTestIssue(t, cfg, 1, "テストタスク", "Open", 1, 3)

	// Issueファイルが作成されていることを確認
	issueFilePath := filepath.Join(cfg.IssuesDir, "1_O_テストタスク.md")
	if !fileExists(issueFilePath) {
		t.Errorf("Issueファイルが作成されていません: %s", issueFilePath)
	}

	// Issueが正しく読み込めることを確認
	issue, err := parser.ParseIssueFile(issueFilePath)
	if err != nil {
		t.Fatalf("Issueファイルの読み込みに失敗しました: %v", err)
	}

	// メタデータの確認
	if issue.ID != 1 {
		t.Errorf("不正なID, 期待値: %d, 実際: %d", 1, issue.ID)
	}
	if issue.Title != "テストタスク" {
		t.Errorf("不正なタイトル, 期待値: %s, 実際: %s", "テストタスク", issue.Title)
	}
	if issue.Status != "Open" {
		t.Errorf("不正なステータス, 期待値: %s, 実際: %s", "Open", issue.Status)
	}
	if issue.Epic != 1 {
		t.Errorf("不正なEpic ID, 期待値: %d, 実際: %d", 1, issue.Epic)
	}
	if issue.Estimate != 3 {
		t.Errorf("不正な見積もり, 期待値: %d, 実際: %d", 3, issue.Estimate)
	}
}

/**
 * Issueファイルのメタデータと内容の整合性をテスト
 *
 * ファイル内のメタデータとファイル名から解析される情報が一致することを確認します。
 * これにより、ファイル名とファイル内容の整合性が保たれていることを保証します。
 */
func TestIssueMetadataConsistency(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のIssueを作成（異なるパラメータでいくつか作成）
	createTestIssue(t, cfg, 1, "重要タスク", "Open", 1, 5)
	createTestIssue(t, cfg, 2, "緊急タスク", "Open", 1, 8)
	createTestIssue(t, cfg, 3, "完了タスク", "Close", 1, 3)

	// 各Issueファイルを検証
	testCases := []struct {
		id       int
		title    string
		status   string
		epic     int
		estimate int
	}{
		{1, "重要タスク", "Open", 1, 5},
		{2, "緊急タスク", "Open", 1, 8},
		{3, "完了タスク", "Close", 1, 3},
	}

	for _, tc := range testCases {
		fileName := ""
		if tc.status == "Open" {
			fileName = filepath.Join(cfg.IssuesDir, fmt.Sprintf("%d_O_%s.md", tc.id, tc.title))
		} else {
			fileName = filepath.Join(cfg.IssuesDir, fmt.Sprintf("%d_C_%s.md", tc.id, tc.title))
		}

		// ファイルの読み込みとパース
		issue, err := parser.ParseIssueFile(fileName)
		if err != nil {
			t.Fatalf("Issueファイルの読み込みに失敗しました: %v", err)
		}

		// メタデータの検証
		if issue.ID != tc.id {
			t.Errorf("不正なID, 期待値: %d, 実際: %d", tc.id, issue.ID)
		}
		if issue.Title != tc.title {
			t.Errorf("不正なタイトル, 期待値: %s, 実際: %s", tc.title, issue.Title)
		}
		if issue.Status != tc.status {
			t.Errorf("不正なステータス, 期待値: %s, 実際: %s", tc.status, issue.Status)
		}
		if issue.Epic != tc.epic {
			t.Errorf("不正なEpic ID, 期待値: %d, 実際: %d", tc.epic, issue.Epic)
		}
		if issue.Estimate != tc.estimate {
			t.Errorf("不正な見積もり, 期待値: %d, 実際: %d", tc.estimate, issue.Estimate)
		}
	}
}
