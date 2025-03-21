package commands

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/watcher"
)

// デフォルトのデバウンス時間（ミリ秒）
const defaultDebounceTime = 500 * time.Millisecond

// WatchCommand - 指定したプロジェクトパスの監視を開始
func WatchCommand(cfg *config.Config, projectPath string) error {
	// プロジェクトパスが指定されていない場合はcfgから取得
	if projectPath == "" {
		projectPath = cfg.ProjectsDir
	}

	// パスが相対パスの場合は絶対パスに変換
	if !filepath.IsAbs(projectPath) {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("カレントディレクトリの取得に失敗しました: %w", err)
		}
		projectPath = filepath.Join(currentDir, projectPath)
	}

	// プロジェクトパスの検証
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("指定されたプロジェクトパスが存在しません: %s", projectPath)
	}

	// issuesディレクトリの検証
	issuesDir := filepath.Join(projectPath, "issues")
	if _, err := os.Stat(issuesDir); os.IsNotExist(err) {
		return fmt.Errorf("issuesディレクトリが存在しません: %s", issuesDir)
	}

	// ウォッチマネージャーの取得
	manager := watcher.GetManager()

	// 既に監視中かチェック
	if manager.IsWatching(projectPath) {
		return fmt.Errorf("プロジェクト '%s' は既に監視中です", projectPath)
	}

	// 監視の開始
	err := manager.StartWatching(projectPath, defaultDebounceTime)
	if err != nil {
		return fmt.Errorf("監視の開始に失敗しました: %w", err)
	}

	fmt.Printf("\nプロジェクト '%s' の監視を開始しました\n", projectPath)
	fmt.Println("監視を停止するには Ctrl+C を押すか、別のターミナルで 'ib unwatch' を実行してください")

	// シグナルハンドリング（Ctrl+Cでの終了）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ブロッキング処理
	<-sigChan

	// 監視の停止
	fmt.Println("\n===== 監視を停止しています... =====")
	if err := manager.StopWatching(projectPath); err != nil {
		return fmt.Errorf("監視の停止に失敗しました: %w", err)
	}

	return nil
}
