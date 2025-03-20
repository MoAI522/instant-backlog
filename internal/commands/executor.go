package commands

import (
	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/watcher"
)

// CommandExecutorImpl - watcher.CommandExecutor の実装
type CommandExecutorImpl struct{}

// ExecuteSync - SyncCommand を実行
func (e *CommandExecutorImpl) ExecuteSync(cfg *config.Config) error {
	return SyncCommand(cfg)
}

// ExecuteRename - RenameCommand を実行
func (e *CommandExecutorImpl) ExecuteRename(cfg *config.Config) error {
	return RenameCommand(cfg)
}

// RegisterCommandExecutor - コマンド実行インスタンスを登録
func RegisterCommandExecutor() {
	watcher.SetCommandExecutor(&CommandExecutorImpl{})
}
