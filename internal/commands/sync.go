package commands

import (
	"fmt"

	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/fileops"
	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
)

// SyncCommand - order.csvとIssueファイルの同期を行う
func SyncCommand(cfg *config.Config) error {
	fmt.Println("order.csvの同期を開始します...")
	
	// 1. すべてのIssueを読み込む
	issues, err := fileops.ReadAllIssues(cfg.IssuesDir)
	if err != nil {
		return fmt.Errorf("Issueの読み込みに失敗しました: %w", err)
	}
	
	// 2. 現在のorder.csvを読み込む
	orderItems, err := parser.ReadOrderCSV(cfg.OrderCSV)
	if err != nil {
		return fmt.Errorf("order.csvの読み込みに失敗しました: %w", err)
	}
	
	// 3. Closeになっているものをorder.csvから削除
	// 4. 新しいOpenのIssueをorder.csvに追加
	var newOrderItems []models.OrderCSVItem
	
	// 既存のマップを作成
	existingIDs := make(map[int]bool)
	for _, issue := range issues {
		if issue.Status == "Open" {
			existingIDs[issue.ID] = true
		}
	}
	
	// クローズされたIssueを除外
	for _, item := range orderItems {
		if existingIDs[item.ID] {
			newOrderItems = append(newOrderItems, item)
			delete(existingIDs, item.ID)
		}
	}
	
	// 新しいOpenのIssueを追加
	for _, issue := range issues {
		if issue.Status == "Open" && existingIDs[issue.ID] {
			newOrderItems = append(newOrderItems, models.OrderCSVItem{
				ID:       issue.ID,
				Title:    issue.Title,
				Epic:     issue.Epic,
				Estimate: issue.Estimate,
			})
		}
	}
	
	// 5. 更新したorder.csvを書き込む
	if err := parser.WriteOrderCSV(cfg.OrderCSV, newOrderItems); err != nil {
		return fmt.Errorf("order.csvの書き込みに失敗しました: %w", err)
	}
	
	fmt.Printf("同期完了: %d件のIssueがorder.csvに保存されました\n", len(newOrderItems))
	return nil
}
