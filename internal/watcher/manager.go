package watcher

import (
	"fmt"
	"sync"
	"time"
)

// WatchManager - 複数プロジェクトの監視を管理するシングルトンマネージャー
type WatchManager struct {
	watchers map[string]*ProjectWatcher // プロジェクトパスをキーとしたウォッチャーマップ
	mu       sync.Mutex                 // スレッドセーフ操作のためのミューテックス
}

// グローバルな監視マネージャーインスタンス
var (
	instance *WatchManager
	once     sync.Once
)

// GetManager - シングルトンパターンでWatchManagerのインスタンスを取得
func GetManager() *WatchManager {
	once.Do(func() {
		instance = &WatchManager{
			watchers: make(map[string]*ProjectWatcher),
		}
	})
	return instance
}

// StartWatching - 指定したプロジェクトパスの監視を開始
func (m *WatchManager) StartWatching(projectPath string, debounceTime time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 既に監視中の場合はエラー
	if _, exists := m.watchers[projectPath]; exists {
		return fmt.Errorf("プロジェクト '%s' は既に監視中です", projectPath)
	}

	// 新しいプロジェクトウォッチャーを作成
	watcher, err := NewProjectWatcher(projectPath, debounceTime)
	if err != nil {
		return fmt.Errorf("ウォッチャーの作成に失敗しました: %w", err)
	}

	// 監視開始
	if err := watcher.Start(); err != nil {
		return fmt.Errorf("監視の開始に失敗しました: %w", err)
	}

	// マップに追加
	m.watchers[projectPath] = watcher
	return nil
}

// StopWatching - 指定したプロジェクトパスの監視を停止
func (m *WatchManager) StopWatching(projectPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 監視中でない場合はエラー
	watcher, exists := m.watchers[projectPath]
	if !exists {
		return fmt.Errorf("プロジェクト '%s' は監視されていません", projectPath)
	}

	// 監視停止
	if err := watcher.Stop(); err != nil {
		return fmt.Errorf("監視の停止に失敗しました: %w", err)
	}

	// マップから削除
	delete(m.watchers, projectPath)
	return nil
}

// IsWatching - 指定したプロジェクトパスが監視中かどうかを確認
func (m *WatchManager) IsWatching(projectPath string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.watchers[projectPath]
	return exists
}

// GetWatchingProjects - 現在監視中のすべてのプロジェクトパスを取得
func (m *WatchManager) GetWatchingProjects() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	projects := make([]string, 0, len(m.watchers))
	for path := range m.watchers {
		projects = append(projects, path)
	}
	return projects
}

// StopAll - すべてのプロジェクトの監視を停止
func (m *WatchManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for path, watcher := range m.watchers {
		if err := watcher.Stop(); err != nil {
			fmt.Printf("警告: プロジェクト '%s' の監視停止に失敗しました: %v\n", path, err)
		}
	}

	// マップをクリア
	m.watchers = make(map[string]*ProjectWatcher)
}
