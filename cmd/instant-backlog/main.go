package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/moai/instant-backlog/internal/commands"
	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/watcher"
	"github.com/spf13/cobra"
)

func main() {
	// コマンドエグゼキュータを登録
	commands.RegisterCommandExecutor()

	cfg := config.NewConfig()

	// ルートコマンド
	var rootCmd = &cobra.Command{
		Use:     "instant-backlog",
		Aliases: []string{"ib"},
		Short:   "スクラムバックログ管理ツール",
		Long:    `マークダウンファイルを使用してスクラム開発のバックログを管理するシンプルなCLIツール`,
	}

	// syncコマンド
	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "order.csvを同期",
		Long:  `オープンIssueをorder.csvに同期し、クローズIssueを削除します`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.SyncCommand(cfg)
		},
	}

	// renameコマンド
	var renameCmd = &cobra.Command{
		Use:   "rename",
		Short: "ファイル名を更新",
		Long:  `Front Matterの内容に基づいてファイル名を更新します`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return commands.RenameCommand(cfg)
		},
	}

	// watchコマンド
	var watchCmd = &cobra.Command{
		Use:   "watch [project_path]",
		Short: "プロジェクトの監視を開始",
		Long:  `指定したプロジェクト配下のissuesディレクトリを監視し、変更があれば自動でsyncとrenameを実行します`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := ""
			if len(args) > 0 {
				projectPath = args[0]
			}
			
			// シグナルハンドリングのセットアップ
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			
			// ゴルーチンで監視を開始
			errChan := make(chan error, 1)
			go func() {
				errChan <- commands.WatchCommand(cfg, projectPath)
			}()
			
			// シグナルまたはエラーを待機
			select {
			case <-sigChan:
				fmt.Println("\n終了シグナルを受信しました")
				// すべての監視を停止
				watcher.GetManager().StopAll()
				return nil
			case err := <-errChan:
				return err
			}
		},
	}

	// unwatchコマンド
	var unwatchCmd = &cobra.Command{
		Use:   "unwatch [project_path]",
		Short: "プロジェクトの監視を停止",
		Long:  `指定したプロジェクトの監視を停止します。引数がない場合はすべてのプロジェクトの監視を停止します`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// 指定されたプロジェクトの監視を停止
				return commands.UnwatchCommand(cfg, args[0])
			} else {
				// すべてのプロジェクトの監視を停止
				return commands.UnwatchAllCommand(cfg)
			}
		},
	}

	// コマンドをルートコマンドに追加
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(renameCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(unwatchCmd)

	// 実行
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
