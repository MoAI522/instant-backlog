package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/watcher"
)

// MockCommandExecutor - テスト用のコマンド実行モック
type MockCommandExecutor struct {
	SyncCalled   bool
	RenameCalled bool
}

// ExecuteSync - モックのSyncCommand実行
func (e *MockCommandExecutor) ExecuteSync(cfg *config.Config) error {
	e.SyncCalled = true
	return nil
}

// ExecuteRename - モックのRenameCommand実行
func (e *MockCommandExecutor) ExecuteRename(cfg *config.Config) error {
	e.RenameCalled = true
	return nil
}

// ファイル監視機能のテスト
func TestFileWatcher(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のデータを作成
	createTestEpic(t, cfg, 1, "テストエピック", "Open")
	createTestIssue(t, cfg, 1, "監視テストタスク", "Open", 1, 3)

	// モックエグゼキュータの作成と登録
	mockExecutor := &MockCommandExecutor{}
	watcher.SetCommandExecutor(mockExecutor)

	// WatchManagerのインスタンスを取得
	manager := watcher.GetManager()

	// 絶対パスを取得
	absPath, err := filepath.Abs(cfg.ProjectsDir)
	if err != nil {
		t.Fatalf("絶対パスの取得に失敗しました: %v", err)
	}

	// 短いデバウンス時間で監視を開始
	err = manager.StartWatching(absPath, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("監視の開始に失敗しました: %v", err)
	}

	// 正しく監視されていることを確認
	if !manager.IsWatching(absPath) {
		t.Errorf("プロジェクトが正しく監視されていません: %s", absPath)
	}

	// ファイル変更をシミュレート
	time.Sleep(200 * time.Millisecond) // 初期化のための待機
	issueFilePath := filepath.Join(cfg.IssuesDir, "1_O_監視テストタスク.md")

	// ファイル内容を読み込む
	content, err := os.ReadFile(issueFilePath)
	if err != nil {
		t.Fatalf("ファイル読み込みに失敗しました: %v", err)
	}

	// 内容を変更
	updatedContent := string(content) + "\n\n追加のコンテンツ"

	// 変更内容を書き込み
	err = os.WriteFile(issueFilePath, []byte(updatedContent), 0644)
	if err != nil {
		t.Fatalf("ファイル書き込みに失敗しました: %v", err)
	}

	// 監視イベントが発火するのを待機
	time.Sleep(500 * time.Millisecond)

	// 監視の停止
	err = manager.StopWatching(absPath)
	if err != nil {
		t.Fatalf("監視の停止に失敗しました: %v", err)
	}

	// 監視が停止していることを確認
	if manager.IsWatching(absPath) {
		t.Errorf("プロジェクトの監視が正しく停止されていません: %s", absPath)
	}

	// コマンドが呼び出されたことを確認
	if !mockExecutor.SyncCalled {
		t.Error("syncコマンドが呼び出されていません")
	}
	if !mockExecutor.RenameCalled {
		t.Error("renameコマンドが呼び出されていません")
	}
}

// 複数プロジェクト監視のテスト
func TestMultiProjectWatcher(t *testing.T) {
	// テスト環境をセットアップ（2つのプロジェクト）
	cfg1, cleanup1 := setupTestEnvironment(t)
	defer cleanup1()

	// 2つ目のプロジェクトは異なるパスを使用
	tempDir2 := filepath.Join(os.TempDir(), fmt.Sprintf("instant-backlog-test-alt-%d", os.Getpid()))
	epicDir2 := filepath.Join(tempDir2, "projects", "epic")
	issuesDir2 := filepath.Join(tempDir2, "projects", "issues")

	// ディレクトリ作成
	if err := os.MkdirAll(epicDir2, 0755); err != nil {
		t.Fatalf("テスト環境のセットアップに失敗しました: %v", err)
	}
	if err := os.MkdirAll(issuesDir2, 0755); err != nil {
		t.Fatalf("テスト環境のセットアップに失敗しました: %v", err)
	}

	// 2つ目のテスト用の設定を作成
	cfg2 := &config.Config{
		ProjectsDir: filepath.Join(tempDir2, "projects"),
		EpicDir:     epicDir2,
		IssuesDir:   issuesDir2,
		OrderCSV:    filepath.Join(tempDir2, "projects", "order.csv"),
	}

	// クリーンアップ関数
	cleanup2 := func() {
		os.RemoveAll(tempDir2)
	}
	defer cleanup2()

	// テスト用のデータを作成
	createTestIssue(t, cfg1, 1, "プロジェクト1タスク", "Open", 1, 3)
	createTestIssue(t, cfg2, 1, "プロジェクト2タスク", "Open", 1, 3)

	// モックエグゼキュータの作成と登録
	mockExecutor := &MockCommandExecutor{}
	watcher.SetCommandExecutor(mockExecutor)

	// WatchManagerのインスタンスを取得
	manager := watcher.GetManager()

	// 絶対パスを取得
	absPath1, _ := filepath.Abs(cfg1.ProjectsDir)
	absPath2, _ := filepath.Abs(cfg2.ProjectsDir)

	// 両方のプロジェクトの監視を開始
	err1 := manager.StartWatching(absPath1, 100*time.Millisecond)
	err2 := manager.StartWatching(absPath2, 100*time.Millisecond)

	if err1 != nil || err2 != nil {
		t.Fatalf("監視の開始に失敗しました: %v, %v", err1, err2)
	}

	// 監視中のプロジェクト数を確認
	projects := manager.GetWatchingProjects()
	if len(projects) != 2 {
		t.Errorf("監視中のプロジェクト数が正しくありません: 期待値=2, 実際=%d", len(projects))
	}

	// すべての監視を停止
	manager.StopAll()

	// すべての監視が停止していることを確認
	projects = manager.GetWatchingProjects()
	if len(projects) != 0 {
		t.Errorf("すべての監視が正しく停止されていません: 残り=%d", len(projects))
	}
}
