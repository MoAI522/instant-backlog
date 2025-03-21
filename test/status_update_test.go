// test/status_update_test.go
package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moai/instant-backlog/internal/commands"
)

/**
 * ユーザーストーリー5：Issueのステータス変更時にファイル名が更新されること
 *
 * ユーザーがIssueのステータスを変更した際に、ファイル名も自動的に更新されることを期待します。
 * これにより、ファイル名からIssueの状態（Open/Close）が即座に判断できるようになります。
 *
 * 具体的には以下の動作が期待されます：
 * - Issueのステータスが「Open」から「Close」に変更された場合、ファイル名の「_O_」が「_C_」に変わる
 * - Issueのステータスが「Close」から「Open」に変更された場合、ファイル名の「_C_」が「_O_」に変わる
 * - renameコマンドを実行することで、ファイル内容とファイル名の整合性が保たれる
 */
func TestFilenameUpdateOnStatusChange(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// createTestIssue関数でテスト用のOpenステータスのIssueを作成
	createTestIssue(t, cfg, 1, "ステータス変更テスト", "Open", 1, 3)

	// 元のファイル内容を確認
	originalFilePath := filepath.Join(cfg.IssuesDir, "1_O_ステータス変更テスト.md")
	t.Logf("テストファイルパス: %s", originalFilePath)
	if !fileExists(originalFilePath) {
		t.Fatalf("テスト用Issueファイルが作成されていません: %s", originalFilePath)
	}

	// Issueファイルの内容を読み込む
	content, err := os.ReadFile(originalFilePath)
	if err != nil {
		t.Fatalf("ファイル読み込みに失敗しました: %v", err)
	}

	t.Logf("元のファイル内容: %s", string(content))

	// 内容を変更してCloseステータスに更新
	newContent := strings.Replace(string(content), "status: Open", "status: Close", 1)

	// 変更した内容で元のファイルを上書き
	err = os.WriteFile(originalFilePath, []byte(newContent), 0644)
	if err != nil {
		t.Fatalf("ファイル書き込みに失敗しました: %v", err)
	}

	// ファイル内容が書き換えられたことを確認
	updatedContent, err := os.ReadFile(originalFilePath)
	if err != nil {
		t.Fatalf("更新後のファイル読み込みに失敗しました: %v", err)
	}
	if !strings.Contains(string(updatedContent), "status: Close") {
		t.Errorf("ファイル内容が正しく更新されていません")
		t.Logf("実際のファイル内容: %s", string(updatedContent))
	}

	// renameコマンドを実行
	err = commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// 新しいファイル名（Closeステータス）のファイルが作成されているか確認
	newFilePath := filepath.Join(cfg.IssuesDir, "1_C_ステータス変更テスト.md")
	if !fileExists(newFilePath) {
		t.Errorf("ステータス変更後のファイル名が正しく更新されていません: %s", newFilePath)
	}

	// 元のファイル名のファイルが削除されているか確認
	if fileExists(originalFilePath) {
		t.Errorf("元のファイル名がまだ残っています: %s", originalFilePath)
	}
}

/**
 * 複数のステータス変更を連続して行った場合のテスト
 *
 * 複数のIssueで異なるステータス変更を同時に行った場合でも、
 * すべてのファイル名が正しく更新されることを確認します。
 */
func TestMultipleStatusChanges(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 異なるステータスのIssueを複数作成
	createTestIssue(t, cfg, 1, "OpenからCloseへ", "Open", 1, 3)
	createTestIssue(t, cfg, 2, "CloseからOpenへ", "Close", 1, 2)
	createTestIssue(t, cfg, 3, "変更なし_Open", "Open", 1, 1)
	createTestIssue(t, cfg, 4, "変更なし_Close", "Close", 1, 4)

	// 元のファイルパスを保存
	file1Path := filepath.Join(cfg.IssuesDir, "1_O_OpenからCloseへ.md")
	file2Path := filepath.Join(cfg.IssuesDir, "2_C_CloseからOpenへ.md")
	file3Path := filepath.Join(cfg.IssuesDir, "3_O_変更なし_Open.md")
	file4Path := filepath.Join(cfg.IssuesDir, "4_C_変更なし_Close.md")

	// すべてのファイルが存在することを確認
	for _, path := range []string{file1Path, file2Path, file3Path, file4Path} {
		if !fileExists(path) {
			t.Fatalf("テスト用Issueファイルが作成されていません: %s", path)
		}
	}

	// 各ファイルの内容を読み込み、ステータスを変更
	// ファイル1: OpenからCloseへ
	content1, _ := os.ReadFile(file1Path)
	newContent1 := strings.Replace(string(content1), "status: Open", "status: Close", 1)
	os.WriteFile(file1Path, []byte(newContent1), 0644)

	// ファイル2: CloseからOpenへ
	content2, _ := os.ReadFile(file2Path)
	newContent2 := strings.Replace(string(content2), "status: Close", "status: Open", 1)
	os.WriteFile(file2Path, []byte(newContent2), 0644)

	// ファイル3,4はそのまま

	// renameコマンドを実行
	err := commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// 更新後の期待されるファイルパス
	expectedFile1Path := filepath.Join(cfg.IssuesDir, "1_C_OpenからCloseへ.md")
	expectedFile2Path := filepath.Join(cfg.IssuesDir, "2_O_CloseからOpenへ.md")
	// ファイル3,4は変更なし

	// ステータスが変更されたファイルの名前が更新されているか確認
	if !fileExists(expectedFile1Path) || fileExists(file1Path) {
		t.Errorf("ファイル1のステータス変更後のファイル名が正しく更新されていません")
	}

	if !fileExists(expectedFile2Path) || fileExists(file2Path) {
		t.Errorf("ファイル2のステータス変更後のファイル名が正しく更新されていません")
	}

	// 変更がないファイルは元のままか確認
	if !fileExists(file3Path) {
		t.Errorf("変更のないファイル3が不正に変更されています")
	}

	if !fileExists(file4Path) {
		t.Errorf("変更のないファイル4が不正に変更されています")
	}
}

/**
 * ステータス変更と同時にタイトル変更を行った場合のテスト
 *
 * Issueのステータスとタイトルの両方を同時に変更した場合でも、
 * ファイル名が正しく更新されることを確認します。
 */
func TestStatusAndTitleChange(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のIssueを作成
	createTestIssue(t, cfg, 1, "元のタイトル", "Open", 1, 3)

	// 元のファイルパス
	originalPath := filepath.Join(cfg.IssuesDir, "1_O_元のタイトル.md")
	if !fileExists(originalPath) {
		t.Fatalf("テスト用Issueファイルが作成されていません: %s", originalPath)
	}

	// ファイルの内容を読み込む
	content, err := os.ReadFile(originalPath)
	if err != nil {
		t.Fatalf("ファイル読み込みに失敗しました: %v", err)
	}

	// ステータスとタイトルの両方を変更
	newContent := strings.Replace(string(content), "status: Open", "status: Close", 1)
	newContent = strings.Replace(newContent, "title: 元のタイトル", "title: 新しいタイトル", 1)

	// 変更した内容で上書き
	err = os.WriteFile(originalPath, []byte(newContent), 0644)
	if err != nil {
		t.Fatalf("ファイル書き込みに失敗しました: %v", err)
	}

	// renameコマンドを実行
	err = commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// 更新後の期待されるファイルパス
	expectedPath := filepath.Join(cfg.IssuesDir, "1_C_新しいタイトル.md")

	// ファイル名が正しく更新されているか確認
	if !fileExists(expectedPath) {
		t.Errorf("ステータスとタイトル変更後のファイル名が正しく更新されていません: %s", expectedPath)
	}

	if fileExists(originalPath) {
		t.Errorf("元のファイル名がまだ残っています: %s", originalPath)
	}

	// 更新されたファイルの内容を確認
	updatedContent, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("更新後のファイル読み込みに失敗しました: %v", err)
	}

	if !strings.Contains(string(updatedContent), "status: Close") ||
		!strings.Contains(string(updatedContent), "title: 新しいタイトル") {
		t.Errorf("ファイル内容が正しく更新されていません")
		t.Logf("実際のファイル内容: %s", string(updatedContent))
	}
}
