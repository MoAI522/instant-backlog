// test/filename_organization_test.go
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moai/instant-backlog/internal/commands"
	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
	"github.com/moai/instant-backlog/pkg/utils"
)

/**
 * ユーザーストーリー4：マークダウンファイルのファイル名を整理できること
 *
 * ユーザーはマークダウンファイルのファイル名を統一された命名規則に従って整理したい。
 * ファイル名は「ID_ステータス_タイトル.md」の形式で、以下の特徴がある：
 * - ID: 一意の識別子（数字）
 * - ステータス: O（Open）またはC（Close）の1文字
 * - タイトル: 課題やエピックのタイトル
 *
 * renameコマンドを実行することで、不正な名前のファイルが正しい名前に修正される。
 * これにより、ファイル名からステータスや内容を一目で把握できるようになる。
 */
func TestOrganizeMarkdownFilenames(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のEpicとIssueを作成（不正なファイル名でファイル作成）
	epicContent := &models.Epic{
		ID:      1,
		Title:   "テストエピック",
		Status:  "Open",
		Content: "これはテスト用のエピックです。",
	}
	epicMd, _ := parser.GenerateMarkdown(epicContent, epicContent.Content)
	epicPath := filepath.Join(cfg.EpicDir, "incorrect_epic_name.md")
	if err := os.WriteFile(epicPath, epicMd, 0644); err != nil {
		t.Fatalf("テスト用Epicファイルの作成に失敗しました: %v", err)
	}

	issueContent := &models.Issue{
		ID:       2,
		Title:    "テストタスク",
		Status:   "Open",
		Epic:     1,
		Estimate: 3,
		Content:  "これはテスト用のタスクです。",
	}
	issueMd, _ := parser.GenerateMarkdown(issueContent, issueContent.Content)
	issuePath := filepath.Join(cfg.IssuesDir, "incorrect_issue_name.md")
	if err := os.WriteFile(issuePath, issueMd, 0644); err != nil {
		t.Fatalf("テスト用Issueファイルの作成に失敗しました: %v", err)
	}

	// renameコマンドを実行
	err := commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// 正しいファイル名に変更されているか確認
	correctEpicName := utils.GenerateFilename(epicContent.ID, epicContent.Status, epicContent.Title)
	correctEpicPath := filepath.Join(cfg.EpicDir, correctEpicName)
	if !fileExists(correctEpicPath) {
		t.Errorf("Epicファイル名が正しく更新されていません: %s", correctEpicPath)
	}
	if fileExists(epicPath) {
		t.Errorf("不正なEpicファイル名がまだ残っています: %s", epicPath)
	}

	correctIssueName := utils.GenerateFilename(issueContent.ID, issueContent.Status, issueContent.Title)
	correctIssuePath := filepath.Join(cfg.IssuesDir, correctIssueName)
	if !fileExists(correctIssuePath) {
		t.Errorf("Issueファイル名が正しく更新されていません: %s", correctIssuePath)
	}
	if fileExists(issuePath) {
		t.Errorf("不正なIssueファイル名がまだ残っています: %s", issuePath)
	}
}

/**
 * 特殊文字を含むファイル名の正規化テスト
 *
 * タイトルに特殊文字が含まれる場合でも、適切にファイル名が生成され、
 * ファイルシステムで扱える形式に正規化されることを確認します。
 */
func TestFilenameNormalization(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 特殊文字を含むタイトルのIssueを作成
	specialTitles := []string{
		"特殊/文字を含む",
		"記号：*?\"<>|を含む",
		"スペース   を含む",
		"非常に長いタイトルで100文字以上になるようなケースを想定したテストケースとしてこのような長い文字列を使用する場合のテスト",
	}

	for i, title := range specialTitles {
		issueID := i + 1

		// 特殊文字を含むIssueを作成（通常のIssue作成関数を使用）
		issue := &models.Issue{
			ID:       issueID,
			Title:    title,
			Status:   "Open",
			Epic:     1,
			Estimate: 3,
			Content:  "特殊文字テスト用のタスクです。",
		}

		// マークダウンを生成
		issueMd, _ := parser.GenerateMarkdown(issue, issue.Content)

		// 安全なファイル名を生成
		sanitizedTitle := utils.SanitizeFilename(title)
		safePath := filepath.Join(cfg.IssuesDir, "special_"+sanitizedTitle+".md")

		if err := os.WriteFile(safePath, issueMd, 0644); err != nil {
			t.Fatalf("テスト用Issueファイルの作成に失敗しました: %v", err)
		}
	}

	// renameコマンドを実行
	err := commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// 各Issueが正しいファイル名に変更されているか確認
	for i, title := range specialTitles {
		issueID := i + 1
		issue := &models.Issue{
			ID:     issueID,
			Title:  title,
			Status: "Open",
		}

		// 正しいファイル名を生成
		correctName := utils.GenerateFilename(issue.ID, issue.Status, issue.Title)
		correctPath := filepath.Join(cfg.IssuesDir, correctName)

		if !fileExists(correctPath) {
			t.Errorf("特殊文字を含むIssueのファイル名が正しく更新されていません: %s", correctPath)

			// ディレクトリ内のファイルを表示
			files, _ := filepath.Glob(filepath.Join(cfg.IssuesDir, "*.md"))
			for _, f := range files {
				t.Logf("ディレクトリ内のファイル: %s", filepath.Base(f))
			}
		}
	}
}

/**
 * 同一内容の重複ファイルのリネーム処理テスト
 *
 * 同じ内容を持つ複数のファイルが異なるファイル名で存在する場合、
 * 正しく1つのファイルに統合されることを確認します。
 */
func TestDuplicateFileRenaming(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 同一内容のIssueを異なるファイル名で作成
	issue := &models.Issue{
		ID:       10,
		Title:    "重複テストタスク",
		Status:   "Open",
		Epic:     1,
		Estimate: 5,
		Content:  "これは重複テスト用のタスクです。",
	}

	// マークダウンを生成
	issueMd, _ := parser.GenerateMarkdown(issue, issue.Content)

	// 異なるファイル名で同じ内容を保存
	duplicatePaths := []string{
		filepath.Join(cfg.IssuesDir, "duplicate_1.md"),
		filepath.Join(cfg.IssuesDir, "duplicate_2.md"),
		filepath.Join(cfg.IssuesDir, "10_wrong_name.md"),
	}

	for _, path := range duplicatePaths {
		if err := os.WriteFile(path, issueMd, 0644); err != nil {
			t.Fatalf("テスト用重複ファイルの作成に失敗しました: %v", err)
		}
	}

	// renameコマンドを実行
	err := commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("renameコマンドの実行に失敗しました: %v", err)
	}

	// 正しいファイル名のみが残っているか確認
	correctName := utils.GenerateFilename(issue.ID, issue.Status, issue.Title)
	correctPath := filepath.Join(cfg.IssuesDir, correctName)

	if !fileExists(correctPath) {
		t.Errorf("正しいファイル名のファイルが存在しません: %s", correctPath)
	}

	// 元の重複ファイルが削除されているか確認
	for _, path := range duplicatePaths {
		if path != correctPath && fileExists(path) {
			t.Errorf("重複ファイルが削除されていません: %s", path)
		}
	}

	// ディレクトリ内のIssueファイル数が1つだけであることを確認
	issueFiles, _ := filepath.Glob(filepath.Join(cfg.IssuesDir, "*.md"))
	if len(issueFiles) != 1 {
		t.Errorf("重複除去後のIssueファイル数が不正です。期待値: 1, 実際: %d", len(issueFiles))
		for _, f := range issueFiles {
			t.Logf("残存ファイル: %s", filepath.Base(f))
		}
	}
}
