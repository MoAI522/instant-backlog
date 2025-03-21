package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/parser"
	"github.com/moai/instant-backlog/pkg/utils"
)

// RenameCommand - Front Matterの内容に基づいてファイル名を更新
func RenameCommand(cfg *config.Config) error {
	fmt.Println("ファイル名の更新を開始します...")
	fmt.Printf("RenameCommand: EpicDir=%s, IssuesDir=%s\n", cfg.EpicDir, cfg.IssuesDir)
	
	// Epicファイルの更新
	if err := renameEpicFiles(cfg.EpicDir); err != nil {
		return fmt.Errorf("Epicファイルの更新に失敗しました: %w", err)
	}
	
	// Issueファイルの更新
	if err := renameIssueFiles(cfg.IssuesDir); err != nil {
		return fmt.Errorf("Issueファイルの更新に失敗しました: %w", err)
	}
	
	fmt.Println("ファイル名の更新が完了しました")
	return nil
}

// renameEpicFiles - Epicディレクトリ内のファイル名を更新
func renameEpicFiles(directory string) error {
	fmt.Printf("renameEpicFiles: ディレクトリ %s のファイル名を更新します\n", directory)
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	
	// デバッグ情報: ディレクトリ内のファイル一覧を表示
	fmt.Printf("ディレクトリ %s 内のファイル:\n", directory)
	for _, f := range files {
		fmt.Printf("- %s\n", f.Name())
	}
	
	// ディレクトリ内のファイル数をチェック
	if len(files) == 0 {
		fmt.Printf("renameEpicFiles: ディレクトリは空です\n")
		return nil
	}
	
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}
		
		filePath := filepath.Join(directory, file.Name())
		// デバッグ: ファイルパスの表示
		fmt.Printf("解析するファイル: %s\n", filePath)
		
		epic, err := parser.ParseEpicFile(filePath)
		if err != nil {
			fmt.Printf("警告: ファイルの解析に失敗しました %s: %v\n", file.Name(), err)
			continue
		}
		
		// デバッグ: 解析結果の表示
		fmt.Printf("解析結果: ID=%d, Title=%s, Status=%s\n", epic.ID, epic.Title, epic.Status)
		
		// 正しいファイル名を生成
		correctFilename := utils.GenerateFilename(epic.ID, epic.Status, epic.Title)
		fmt.Printf("生成したファイル名: %s\n", correctFilename)
		
		// 現在のファイル名と違う場合は名前変更
		fmt.Printf("renameEpicFiles: ファイル名チェック - 現在: %s, 正しい名前: %s\n", file.Name(), correctFilename)
		if file.Name() != correctFilename {
			newPath := filepath.Join(directory, correctFilename)
			fmt.Printf("リネーム: %s -> %s\n", file.Name(), correctFilename)
			
			// 一時ファイルが既に存在する場合は削除
			if _, err := os.Stat(newPath); err == nil {
				fmt.Printf("警告: 対象ファイルが既に存在します。置き換えます: %s\n", newPath)
				err = os.Remove(newPath)
				fmt.Printf("renameEpicFiles: 既存ファイル削除結果: %v\n", err)
			}
			
			fmt.Printf("renameEpicFiles: リネーム実行: %s -> %s\n", filePath, newPath)
			err = os.Rename(filePath, newPath)
			if err != nil {
				fmt.Printf("警告: ファイルのリネームに失敗しました %s: %v\n", file.Name(), err)
			}
		}
	}
	
	return nil
}

// renameIssueFiles - Issueディレクトリ内のファイル名を更新
func renameIssueFiles(directory string) error {
	fmt.Printf("renameIssueFiles: ディレクトリ %s のファイル名を更新します\n", directory)
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	
	// デバッグ情報: ディレクトリ内のファイル一覧を表示
	fmt.Printf("ディレクトリ %s 内のファイル:\n", directory)
	for _, f := range files {
		fmt.Printf("- %s\n", f.Name())
	}
	
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}
		
		filePath := filepath.Join(directory, file.Name())
		// デバッグ: ファイルパスの表示
		fmt.Printf("解析するファイル: %s\n", filePath)
		
		issue, err := parser.ParseIssueFile(filePath)
		if err != nil {
			fmt.Printf("警告: ファイルの解析に失敗しました %s: %v\n", file.Name(), err)
			continue
		}
		
		// デバッグ: 解析結果の表示
		fmt.Printf("解析結果: ID=%d, Title=%s, Status=%s\n", issue.ID, issue.Title, issue.Status)
		
		// 正しいファイル名を生成
		correctFilename := utils.GenerateFilename(issue.ID, issue.Status, issue.Title)
		fmt.Printf("生成したファイル名: %s\n", correctFilename)
		
		// 現在のファイル名と違う場合は名前変更
		if file.Name() != correctFilename {
			newPath := filepath.Join(directory, correctFilename)
			fmt.Printf("リネーム: %s -> %s\n", file.Name(), correctFilename)
			
			// 一時ファイルが既に存在する場合は削除
			if _, err := os.Stat(newPath); err == nil {
				fmt.Printf("警告: 対象ファイルが既に存在します。置き換えます: %s\n", newPath)
				os.Remove(newPath)
			}
			
			err = os.Rename(filePath, newPath)
			if err != nil {
				fmt.Printf("警告: ファイルのリネームに失敗しました %s: %v\n", file.Name(), err)
			}
		}
	}
	
	return nil
}
