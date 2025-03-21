// test/integration_test.go
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moai/instant-backlog/internal/commands"
	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/fileops"
	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
	"github.com/moai/instant-backlog/pkg/utils"
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

// ユーザーストーリー1：マークダウンでIssueを管理できること
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

// ユーザーストーリー2：マークダウンでEpicを管理できること
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

// ユーザーストーリー3：Issueの優先順位を管理できること
func TestManageIssuePriorities(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のデータを作成
	createTestEpic(t, cfg, 1, "テストエピック", "Open")
	createTestIssue(t, cfg, 1, "高優先度タスク", "Open", 1, 5)
	createTestIssue(t, cfg, 2, "中優先度タスク", "Open", 1, 3)
	createTestIssue(t, cfg, 3, "低優先度タスク", "Open", 1, 1)
	createTestIssue(t, cfg, 4, "クローズタスク", "Close", 1, 2)

	// テスト用のCSVを作成（優先度順）
	initialOrderItems := []models.OrderCSVItem{
		{ID: 1, Title: "高優先度タスク", Epic: 1, Estimate: 5},
		{ID: 2, Title: "中優先度タスク", Epic: 1, Estimate: 3},
		{ID: 3, Title: "低優先度タスク", Epic: 1, Estimate: 1},
		{ID: 4, Title: "クローズタスク", Epic: 1, Estimate: 2}, // 初期CSVにはクローズタスクも含まれている
	}
	createTestOrderCSV(t, cfg, initialOrderItems)

	// syncコマンドを実行
	err := commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// 更新後のCSVを読み込む
	updatedOrderItems, err := parser.ReadOrderCSV(cfg.OrderCSV)
	if err != nil {
		t.Fatalf("order.csvの読み込みに失敗しました: %v", err)
	}

	// クローズされたタスクがCSVから削除されているか確認
	for _, item := range updatedOrderItems {
		if item.ID == 4 {
			t.Errorf("クローズされたタスク(ID=4)がorder.csvから削除されていません")
		}
	}

	// 残りのオープンタスクが正しく含まれているか確認
	expectedIDs := map[int]bool{1: true, 2: true, 3: true}
	for _, item := range updatedOrderItems {
		if !expectedIDs[item.ID] {
			t.Errorf("予期しないタスクID %d がorder.csvに含まれています", item.ID)
		}
		delete(expectedIDs, item.ID)
	}

	if len(expectedIDs) > 0 {
		var missingIDs []string
		for id := range expectedIDs {
			missingIDs = append(missingIDs, fmt.Sprintf("%d", id))
		}
		t.Errorf("期待されるタスクIDs %s がorder.csvに含まれていません", strings.Join(missingIDs, ", "))
	}
}

// ユーザーストーリー4：マークダウンファイルのファイル名を整理できること
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

	// ディレクトリ内のファイル一覧を取得して確認
	files, _ := os.ReadDir(cfg.IssuesDir)
	for _, file := range files {
		t.Logf("リネーム後のディレクトリ内のファイル: %s", file.Name())
	}

	// Closeステータスに対応する正しいファイル名に変更されているか確認
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

// Issueのステータス変更時にファイル名が更新されることを確認するテスト
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
	
	// 2回目のrenameコマンド実行を確認するために再度実行
	err = commands.RenameCommand(cfg)
	if err != nil {
		t.Fatalf("2回目のrenameコマンドの実行に失敗しました: %v", err)
	}

	// ディレクトリ内のファイル一覧を確認
	filesBefore, _ := os.ReadDir(cfg.IssuesDir)
	t.Logf("リネーム前のディレクトリ内容:")
	for _, file := range filesBefore {
		t.Logf("- %s", file.Name())
		// ファイルの内容を表示
		filepath := filepath.Join(cfg.IssuesDir, file.Name())
		fileContent, _ := os.ReadFile(filepath)
		t.Logf("ファイル内容: %s", string(fileContent))
	}

	// renameコマンドを実行
	newFilePath := filepath.Join(cfg.IssuesDir, "1_C_ステータス変更テスト.md")
	if !fileExists(newFilePath) {
		t.Errorf("ステータス変更後のファイル名が正しく更新されていません: %s", newFilePath)
	}
	if fileExists(originalFilePath) {
		t.Errorf("元のファイル名がまだ残っています: %s", originalFilePath)
	}
}

// Epicに紐づいたIssueがすべてClosedになった場合に、EpicがCloseに更新されることをテスト
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
	if err := fileops.WriteIssue(cfg.IssuesDir, issue1); err != nil {
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
	if err := fileops.WriteIssue(cfg.IssuesDir, issue2); err != nil {
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
