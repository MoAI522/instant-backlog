package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/watcher"
)

// UnwatchCommand - 指定したプロジェクトパスの監視を停止
func UnwatchCommand(cfg *config.Config, projectPath string) error {
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

	// ウォッチマネージャーの取得
	manager := watcher.GetManager()

	// 監視中でない場合は通知
	if !manager.IsWatching(projectPath) {
		return fmt.Errorf("プロジェクト '%s' は現在監視されていません", projectPath)
	}

	// 監視の停止
	err := manager.StopWatching(projectPath)
	if err != nil {
		return fmt.Errorf("監視の停止に失敗しました: %w", err)
	}

	fmt.Printf("プロジェクト '%s' の監視を停止しました\n", projectPath)
	return nil
}

// UnwatchAllCommand - すべてのプロジェクトの監視を停止
func UnwatchAllCommand(_ *config.Config) error {
	// ウォッチマネージャーの取得
	manager := watcher.GetManager()

	// 現在監視中のプロジェクトを取得
	projects := manager.GetWatchingProjects()
	if len(projects) == 0 {
		fmt.Println("現在監視中のプロジェクトはありません")
		return nil
	}

	// すべての監視を停止
	manager.StopAll()

	fmt.Printf("すべてのプロジェクト（%d個）の監視を停止しました\n", len(projects))
	return nil
}
