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
	fmt.Println("===== ファイル名の更新を開始します... ====")
	
	// Epicファイルの更新
	if err := renameEpicFiles(cfg.EpicDir); err != nil {
		return fmt.Errorf("Epicファイルの更新に失敗しました: %w", err)
	}
	
	// Issueファイルの更新
	if err := renameIssueFiles(cfg.IssuesDir); err != nil {
		return fmt.Errorf("Issueファイルの更新に失敗しました: %w", err)
	}
	
	fmt.Println("===== ファイル名の更新が完了しました ====")
	return nil
}

// renameEpicFiles - Epicディレクトリ内のファイル名を更新
func renameEpicFiles(directory string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	
	// ディレクトリ内のファイル数をチェック
	if len(files) == 0 {
		return nil
	}
	
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}
		
		filePath := filepath.Join(directory, file.Name())
		
		epic, err := parser.ParseEpicFile(filePath)
		if err != nil {
			fmt.Printf("警告: ファイルの解析に失敗しました %s: %v\n", file.Name(), err)
			continue
		}
		
		// 正しいファイル名を生成
		correctFilename := utils.GenerateFilename(epic.ID, epic.Status, epic.Title)
		
		// 現在のファイル名と違う場合は名前変更
		if file.Name() != correctFilename {
			newPath := filepath.Join(directory, correctFilename)
			fmt.Printf("リネーム: %s → %s\n", file.Name(), correctFilename)
			
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

// renameIssueFiles - Issueディレクトリ内のファイル名を更新
func renameIssueFiles(directory string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".md" {
			continue
		}
		
		filePath := filepath.Join(directory, file.Name())
		
		issue, err := parser.ParseIssueFile(filePath)
		if err != nil {
			fmt.Printf("警告: ファイルの解析に失敗しました %s: %v\n", file.Name(), err)
			continue
		}
		
		// 正しいファイル名を生成
		correctFilename := utils.GenerateFilename(issue.ID, issue.Status, issue.Title)
		
		// 現在のファイル名と違う場合は名前変更
		if file.Name() != correctFilename {
			newPath := filepath.Join(directory, correctFilename)
			fmt.Printf("リネーム: %s → %s\n", file.Name(), correctFilename)
			
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
