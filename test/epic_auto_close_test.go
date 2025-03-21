// test/epic_auto_close_test.go
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/moai/instant-backlog/internal/commands"
	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
)

/**
 * ユーザーストーリー6：Epicに紐づいたIssueがすべてClosedになった場合、Epicが自動的にCloseになること
 *
 * ユーザーは、Epicに紐づくすべてのIssueが完了（Close）した場合に、
 * 自動的にEpicも完了（Close）状態になることを期待します。
 *
 * この機能により、以下のメリットがあります：
 * - Epicの状態管理が自動化され、手動でステータスを変更する手間が省ける
 * - プロジェクトの進捗状況が正確に反映される
 * - 完了したEpicを視覚的に識別しやすくなる
 *
 * syncコマンドを実行することで、Issueの状態に基づいてEpicの状態が自動的に更新されます。
 */
func TestEpicAutoCloseWhenAllIssuesAreClosed(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のEpicを作成
	createTestEpic(t, cfg, 1, "テストエピック", "Open")

	// このEpicに紐づく複数のIssueを作成
	createTestIssue(t, cfg, 1, "テストタスク1", "Open", 1, 3)
	createTestIssue(t, cfg, 2, "テストタスク2", "Open", 1, 5)

	// 初期状態では、EpicがOpenであることを確認
	epic, err := parser.ParseEpicFile(filepath.Join(cfg.EpicDir, "1_O_テストエピック.md"))
	if err != nil {
		t.Fatalf("Epicファイルの読み込みに失敗しました: %v", err)
	}
	if epic.Status != "Open" {
		t.Fatalf("初期状態でEpicのステータスがOpenではありませんでした")
	}

	// 1つ目のIssueをCloseに変更
	issue1Path := filepath.Join(cfg.IssuesDir, "1_O_テストタスク1.md")
	issue1, err := parser.ParseIssueFile(issue1Path)
	if err != nil {
		t.Fatalf("Issueファイルの読み込みに失敗しました: %v", err)
	}
	issue1.Status = "Close"
	if err := os.WriteFile(issue1Path, []byte(generateIssueMarkdown(issue1)), 0644); err != nil {
		t.Fatalf("Issueの更新に失敗しました: %v", err)
	}

	// Issueファイル名を更新
	err = commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// Issue1がCloseされた段階でsyncを実行
	err = commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// EpicはまだCloseになっていないはず
	epic, err = parser.ParseEpicFile(filepath.Join(cfg.EpicDir, "1_O_テストエピック.md"))
	if err != nil {
		t.Fatalf("Epicファイルの読み込みに失敗しました: %v", err)
	}
	if epic.Status != "Open" {
		t.Errorf("Epicのステータスが不正です、期待値: %s, 実際: %s", "Open", epic.Status)
	}

	// 2つ目のIssueもCloseに変更
	issue2Path := filepath.Join(cfg.IssuesDir, "2_O_テストタスク2.md")
	issue2, err := parser.ParseIssueFile(issue2Path)
	if err != nil {
		t.Fatalf("Issueファイルの読み込みに失敗しました: %v", err)
	}
	issue2.Status = "Close"
	if err := os.WriteFile(issue2Path, []byte(generateIssueMarkdown(issue2)), 0644); err != nil {
		t.Fatalf("Issueの更新に失敗しました: %v", err)
	}

	// Issue2のファイル名を更新
	err = commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// すべてのIssueがCloseになったので、syncを実行
	err = commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// この時点でEpicがCloseに更新されているか確認
	// 注意: ファイル名も変更されているはず
	newEpicPath := filepath.Join(cfg.EpicDir, "1_C_テストエピック.md")
	if !fileExists(newEpicPath) {
		// ディレクトリ内のファイルを確認
		files, _ := os.ReadDir(cfg.EpicDir)
		t.Logf("Epicディレクトリ内のファイル:")
		for _, file := range files {
			t.Logf("- %s", file.Name())
		}
		t.Errorf("Epicファイル名が正しく更新されていません: %s", newEpicPath)
	}

	// 新しいファイル名でEpicを読み込み、ステータスがCloseになっているか確認
	updatedEpic, err := parser.ParseEpicFile(newEpicPath)
	if err != nil {
		t.Fatalf("更新後Epicファイルの読み込みに失敗しました: %v", err)
	}

	if updatedEpic.Status != "Close" {
		t.Errorf("Epicのステータスが正しく更新されていません、期待値: %s, 実際: %s", "Close", updatedEpic.Status)
	}
}

/**
 * 複数のEpicとIssueが存在する場合のEpic自動クローズテスト
 *
 * 複数のEpicとそれに紐づくIssueが存在する環境で、
 * 特定のEpicに紐づくIssueがすべてCloseになった場合のみ、
 * そのEpicが自動的にCloseされることを確認します。
 */
func TestMultipleEpicsAutoClose(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 複数のEpicを作成
	createTestEpic(t, cfg, 1, "Epic1", "Open")
	createTestEpic(t, cfg, 2, "Epic2", "Open")

	// 各Epicに紐づくIssueを作成
	// Epic1のIssue
	createTestIssue(t, cfg, 1, "Epic1-Task1", "Open", 1, 3)
	createTestIssue(t, cfg, 2, "Epic1-Task2", "Open", 1, 2)

	// Epic2のIssue
	createTestIssue(t, cfg, 3, "Epic2-Task1", "Open", 2, 5)
	createTestIssue(t, cfg, 4, "Epic2-Task2", "Open", 2, 4)

	// Epic1のIssueをすべてCloseに変更
	for i := 1; i <= 2; i++ {
		issueFilename := fmt.Sprintf("%d_O_Epic1-Task%d.md", i, i)
		issuePath := filepath.Join(cfg.IssuesDir, issueFilename)
		issue, err := parser.ParseIssueFile(issuePath)
		if err != nil {
			t.Fatalf("Issueファイルの読み込みに失敗しました: %v", err)
		}

		issue.Status = "Close"
		if err := os.WriteFile(issuePath, []byte(generateIssueMarkdown(issue)), 0644); err != nil {
			t.Fatalf("Issueの更新に失敗しました: %v", err)
		}
	}

	// ファイル名を更新
	err := commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// syncを実行
	err = commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// Epic1がCloseになっているか確認
	epic1Path := filepath.Join(cfg.EpicDir, "1_C_Epic1.md")
	if !fileExists(epic1Path) {
		t.Errorf("Epic1が自動的にCloseされていません")
	}

	// Epic2はOpenのままか確認
	epic2Path := filepath.Join(cfg.EpicDir, "2_O_Epic2.md")
	if !fileExists(epic2Path) {
		t.Errorf("Epic2が不正に変更されています")
	}

	// Epic2のIssueを1つだけCloseに変更
	issue3Path := filepath.Join(cfg.IssuesDir, "3_O_Epic2-Task1.md")
	issue3, err := parser.ParseIssueFile(issue3Path)
	if err != nil {
		t.Fatalf("Issueファイルの読み込みに失敗しました: %v", err)
	}

	issue3.Status = "Close"
	if err := os.WriteFile(issue3Path, []byte(generateIssueMarkdown(issue3)), 0644); err != nil {
		t.Fatalf("Issueの更新に失敗しました: %v", err)
	}

	// ファイル名を更新
	err = commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// syncを実行
	err = commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// Epic2はまだOpenのままか確認
	if !fileExists(epic2Path) {
		t.Errorf("Epic2が不正に変更されています")
	}
}

/**
 * Issueが存在しないEpicの自動クローズテスト
 *
 * 紐づくIssueが存在しないEpicの場合、
 * 自動クローズ機能が適用されないことを確認します。
 */
func TestEmptyEpicAutoClose(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Issueのないエピックを作成
	createTestEpic(t, cfg, 5, "空のエピック", "Open")

	// syncを実行
	err := commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// エピックのステータスが変わっていないことを確認
	epicPath := filepath.Join(cfg.EpicDir, "5_O_空のエピック.md")
	if !fileExists(epicPath) {
		t.Errorf("空のエピックファイルが不正に変更されています")
	}

	// ステータスがOpenのままであることを確認
	epic, err := parser.ParseEpicFile(epicPath)
	if err != nil {
		t.Fatalf("Epicファイルの読み込みに失敗しました: %v", err)
	}

	if epic.Status != "Open" {
		t.Errorf("空のエピックのステータスが不正に変更されています。期待値: Open, 実際: %s", epic.Status)
	}
}

// テスト用のIssueマークダウンを生成するヘルパー関数
func generateIssueMarkdown(issue *models.Issue) string {
	return `---
id: ` + fmt.Sprintf("%d", issue.ID) + `
title: "` + issue.Title + `"
status: "` + issue.Status + `"
epic: ` + fmt.Sprintf("%d", issue.Epic) + `
estimate: ` + fmt.Sprintf("%d", issue.Estimate) + `
---

これはテスト用のタスクです。
`
}
