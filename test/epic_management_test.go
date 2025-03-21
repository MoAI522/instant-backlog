// test/epic_management_test.go
package test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/moai/instant-backlog/internal/parser"
)

/**
 * ユーザーストーリー2：マークダウンでEpicを管理できること
 *
 * ユーザーは大きな目標や機能単位（Epic）をマークダウン形式で作成・管理したい。
 * 各Epicには以下の情報が含まれる：
 * - ID: 一意の識別子
 * - タイトル: エピックの簡潔な説明
 * - ステータス: エピックの状態（Open/Close）
 * - 内容: エピックの詳細な説明
 *
 * Epicはプロジェクト全体の大きな目標を表し、複数のIssueをまとめる役割を持つ。
 * マークダウンファイルは適切なメタデータを持ち、正しく解析できる必要がある。
 */
func TestManageEpicsWithMarkdown(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のEpicを作成
	createTestEpic(t, cfg, 1, "テストエピック", "Open")

	// Epicファイルが作成されていることを確認
	epicFilePath := filepath.Join(cfg.EpicDir, "1_O_テストエピック.md")
	if !fileExists(epicFilePath) {
		t.Errorf("Epicファイルが作成されていません: %s", epicFilePath)
	}

	// Epicが正しく読み込めることを確認
	epic, err := parser.ParseEpicFile(epicFilePath)
	if err != nil {
		t.Fatalf("Epicファイルの読み込みに失敗しました: %v", err)
	}

	// メタデータの確認
	if epic.ID != 1 {
		t.Errorf("不正なID, 期待値: %d, 実際: %d", 1, epic.ID)
	}
	if epic.Title != "テストエピック" {
		t.Errorf("不正なタイトル, 期待値: %s, 実際: %s", "テストエピック", epic.Title)
	}
	if epic.Status != "Open" {
		t.Errorf("不正なステータス, 期待値: %s, 実際: %s", "Open", epic.Status)
	}
}

/**
 * 複数のEpicの管理と状態の検証
 *
 * 複数のEpicを異なるステータスで作成し、それぞれが正しく管理されることを確認します。
 * これにより、同時に複数のEpicを管理する機能が正しく動作することを保証します。
 */
func TestMultipleEpicManagement(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 複数のテスト用Epicを作成
	testEpics := []struct {
		id     int
		title  string
		status string
	}{
		{1, "進行中エピック", "Open"},
		{2, "完了エピック", "Close"},
		{3, "新規エピック", "Open"},
	}

	for _, epic := range testEpics {
		createTestEpic(t, cfg, epic.id, epic.title, epic.status)
	}

	// 各Epicファイルを検証
	for _, epic := range testEpics {
		fileSuffix := "_O_"
		if epic.status == "Close" {
			fileSuffix = "_C_"
		}
		epicFilePath := filepath.Join(cfg.EpicDir, fmt.Sprintf("%d%s%s.md", epic.id, fileSuffix, epic.title))

		// ファイルの存在を確認
		if !fileExists(epicFilePath) {
			t.Errorf("Epicファイルが作成されていません: %s", epicFilePath)
			continue
		}

		// 読み込みとパース
		parsedEpic, err := parser.ParseEpicFile(epicFilePath)
		if err != nil {
			t.Fatalf("Epicファイルの読み込みに失敗しました: %v", err)
		}

		// メタデータの検証
		if parsedEpic.ID != epic.id {
			t.Errorf("不正なID, 期待値: %d, 実際: %d", epic.id, parsedEpic.ID)
		}
		if parsedEpic.Title != epic.title {
			t.Errorf("不正なタイトル, 期待値: %s, 実際: %s", epic.title, parsedEpic.Title)
		}
		if parsedEpic.Status != epic.status {
			t.Errorf("不正なステータス, 期待値: %s, 実際: %s", epic.status, parsedEpic.Status)
		}
	}
}

/**
 * Epicとそれに関連するIssueの関係の検証
 *
 * Epicに紐づくIssueの関連性を確認し、正しく管理されることを検証します。
 * Epicとそれに紐づくIssueの関係が一貫性を持ち、適切に管理されることを保証します。
 */
func TestEpicToIssueRelationship(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のEpicを作成
	epicID := 5
	epicTitle := "関連性テストエピック"
	createTestEpic(t, cfg, epicID, epicTitle, "Open")

	// このEpicに関連するIssueを作成
	relatedIssues := []struct {
		id       int
		title    string
		estimate int
	}{
		{10, "関連タスク1", 3},
		{11, "関連タスク2", 5},
		{12, "関連タスク3", 8},
	}

	for _, issue := range relatedIssues {
		createTestIssue(t, cfg, issue.id, issue.title, "Open", epicID, issue.estimate)
	}

	// 各Issueが正しくEpicに紐づいているか確認
	for _, issue := range relatedIssues {
		issueFilePath := filepath.Join(cfg.IssuesDir, fmt.Sprintf("%d_O_%s.md", issue.id, issue.title))

		// ファイルの読み込みとパース
		parsedIssue, err := parser.ParseIssueFile(issueFilePath)
		if err != nil {
			t.Fatalf("Issueファイルの読み込みに失敗しました: %v", err)
		}

		// Epic IDの検証
		if parsedIssue.Epic != epicID {
			t.Errorf("Issue %d が正しいEpic IDに紐づいていません, 期待値: %d, 実際: %d",
				issue.id, epicID, parsedIssue.Epic)
		}
	}
}
