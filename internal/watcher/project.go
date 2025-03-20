package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/moai/instant-backlog/internal/config"
)

// ProjectWatcher - 単一プロジェクトの監視を担当する構造体
type ProjectWatcher struct {
	projectPath  string         // 監視対象のプロジェクトパス
	issuesDir    string         // issuesディレクトリのパス
	watcher      *fsnotify.Watcher // fsnotifyのウォッチャー
	debounceTime time.Duration  // デバウンス時間
	stopChan     chan struct{}  // 停止シグナル用のチャネル
	mutex        sync.Mutex     // 並行アクセス用のミューテックス
	isRunning    bool           // 実行中かどうかのフラグ
	lastEventTime time.Time     // 最後のイベント時刻（デバウンス用）
	timer        *time.Timer    // デバウンスタイマー
}

// NewProjectWatcher - 新しいProjectWatcherインスタンスを作成
func NewProjectWatcher(projectPath string, debounceTime time.Duration) (*ProjectWatcher, error) {
	// プロジェクトパスの検証
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("指定されたプロジェクトパスが存在しません: %s", projectPath)
	}

	// issuesディレクトリの検証
	issuesDir := filepath.Join(projectPath, "issues")
	if _, err := os.Stat(issuesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("issuesディレクトリが存在しません: %s", issuesDir)
	}

	return &ProjectWatcher{
		projectPath:  projectPath,
		issuesDir:    issuesDir,
		debounceTime: debounceTime,
		stopChan:     make(chan struct{}),
		isRunning:    false,
	}, nil
}

// Start - 監視を開始
func (pw *ProjectWatcher) Start() error {
	pw.mutex.Lock()
	defer pw.mutex.Unlock()

	if pw.isRunning {
		return fmt.Errorf("ウォッチャーは既に実行中です")
	}

	// fsnotifyウォッチャーの作成
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fsnotifyウォッチャーの作成に失敗しました: %w", err)
	}

	// issuesディレクトリを監視
	err = watcher.Add(pw.issuesDir)
	if err != nil {
		watcher.Close()
		return fmt.Errorf("ディレクトリの監視に失敗しました: %w", err)
	}

	pw.watcher = watcher
	pw.isRunning = true
	pw.timer = time.NewTimer(pw.debounceTime)
	pw.timer.Stop() // 初期状態では停止しておく

	// イベント処理ゴルーチンを起動
	go pw.processEvents()

	fmt.Printf("プロジェクト '%s' の監視を開始しました\n", pw.projectPath)
	return nil
}

// Stop - 監視を停止
func (pw *ProjectWatcher) Stop() error {
	pw.mutex.Lock()
	defer pw.mutex.Unlock()

	if !pw.isRunning {
		return fmt.Errorf("ウォッチャーは実行されていません")
	}

	// 停止シグナルを送信
	close(pw.stopChan)
	
	// fsnotifyウォッチャーを閉じる
	if pw.watcher != nil {
		pw.watcher.Close()
		pw.watcher = nil
	}

	if pw.timer != nil {
		pw.timer.Stop()
	}

	pw.isRunning = false
	fmt.Printf("プロジェクト '%s' の監視を停止しました\n", pw.projectPath)
	return nil
}

// processEvents - ファイルシステムイベントを処理
func (pw *ProjectWatcher) processEvents() {
	for {
		select {
		case <-pw.stopChan:
			// 停止シグナルを受信
			return
			
		case event, ok := <-pw.watcher.Events:
			if !ok {
				// チャネルがクローズされた
				return
			}
			
			// 拡張子が.mdでない場合は無視
			if filepath.Ext(event.Name) != ".md" {
				continue
			}

			// 書き込みイベントを処理
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Chmod) != 0 {
				// デバウンス処理
				pw.mutex.Lock()
				now := time.Now()
				pw.lastEventTime = now
				
				// 既存のタイマーを停止して再設定
				if pw.timer != nil {
					pw.timer.Stop()
				}
				
				pw.timer = time.AfterFunc(pw.debounceTime, func() {
					pw.mutex.Lock()
					// デバウンス時間内に新しいイベントがなかった場合のみ実行
					if time.Since(pw.lastEventTime) >= pw.debounceTime {
						pw.executeCommands()
					}
					pw.mutex.Unlock()
				})
				
				pw.mutex.Unlock()
			}
			
		case err, ok := <-pw.watcher.Errors:
			if !ok {
				// チャネルがクローズされた
				return
			}
			fmt.Printf("監視エラー: %v\n", err)
		}
	}
}

// CommandExecutor - コマンド実行のためのインターフェース
type CommandExecutor interface {
	ExecuteSync(cfg *config.Config) error
	ExecuteRename(cfg *config.Config) error
}

// DefaultCommandExecutor - デフォルトのコマンド実行構造体
type DefaultCommandExecutor struct{}

// ExecuteSync - デフォルトのSyncCommand実行
func (e *DefaultCommandExecutor) ExecuteSync(cfg *config.Config) error {
	fmt.Println("警告: 本番のコマンド実行インスタンスが登録されていません")
	return nil
}

// ExecuteRename - デフォルトのRenameCommand実行
func (e *DefaultCommandExecutor) ExecuteRename(cfg *config.Config) error {
	fmt.Println("警告: 本番のコマンド実行インスタンスが登録されていません")
	return nil
}

// commandExecutor - 現在のコマンド実行インスタンス
var commandExecutor CommandExecutor = &DefaultCommandExecutor{}

// SetCommandExecutor - コマンド実行インスタンスの設定
func SetCommandExecutor(executor CommandExecutor) {
	commandExecutor = executor
}

// executeCommands - 関連コマンドを実行
func (pw *ProjectWatcher) executeCommands() {
	fmt.Printf("ファイル変更を検知しました: %s\n", pw.projectPath)
	
	// 設定オブジェクトの作成
	cfg := &config.Config{
		ProjectsDir: pw.projectPath,
		EpicDir:     filepath.Join(pw.projectPath, "epic"),
		IssuesDir:   pw.issuesDir,
		OrderCSV:    filepath.Join(pw.projectPath, "order.csv"),
	}
	
	// syncコマンドを実行
	fmt.Println("syncコマンドを実行中...")
	if err := commandExecutor.ExecuteSync(cfg); err != nil {
		fmt.Printf("syncコマンドの実行中にエラーが発生しました: %v\n", err)
	}
	
	// renameコマンドを実行
	fmt.Println("renameコマンドを実行中...")
	if err := commandExecutor.ExecuteRename(cfg); err != nil {
		fmt.Printf("renameコマンドの実行中にエラーが発生しました: %v\n", err)
	}
	
	fmt.Println("ファイル変更の処理が完了しました")
}
