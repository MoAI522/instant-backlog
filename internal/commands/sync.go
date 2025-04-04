package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/fileops"
	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
	"github.com/moai/instant-backlog/pkg/utils"
)

// UpdateEpicStatusBasedOnIssues - Epicのステータスを関連するIssueの状態に基づいて更新する
func UpdateEpicStatusBasedOnIssues(cfg *config.Config) error {
	fmt.Println("Epicステータスの更新を開始します...")
	// すべてのIssueを読み込む
	issues, err := fileops.ReadAllIssues(cfg.IssuesDir)
	if err != nil {
		return fmt.Errorf("Issueの読み込みに失敗しました: %w", err)
	}

	// IDごとにIssueをグループ化
	issuesByID := make(map[int][]*models.Issue)
	for _, issue := range issues {
		issuesByID[issue.ID] = append(issuesByID[issue.ID], issue)
	}

	// 重複チェック
	for id, issueGroup := range issuesByID {
		if len(issueGroup) > 1 {
			fmt.Printf("警告: Issue ID=%d に複数のファイルが見つかりました（%d件）\n", id, len(issueGroup))
			for _, i := range issueGroup {
				fmt.Printf("  - ID=%d, Title=%s, Status=%s\n", i.ID, i.Title, i.Status)
			}
		}
	}

	// すべてのEpicを読み込む
	epics, err := fileops.ReadAllEpics(cfg.EpicDir)
	if err != nil {
		return fmt.Errorf("Epicの読み込みに失敗しました: %w", err)
	}

	// ステータスが変更されたかどうかを追跡
	statusChanged := false

	// Epic IDごとにIssueをグループ化
	issuesByEpic := make(map[int][]*models.Issue)
	for _, issue := range issues {
		issuesByEpic[issue.Epic] = append(issuesByEpic[issue.Epic], issue)
	}

	// 各Epicについて、関連するIssueのステータスをチェック
	for _, epic := range epics {
		if epic.Status == "Close" {
			// すでにClosedなら何もしない
			continue
		}

		// このEpicに紐づくIssueを取得
		epicIssues := issuesByEpic[epic.ID]

		// 紐づくIssueがない場合はスキップ
		if len(epicIssues) == 0 {
			continue
		}

		// すべてのIssueがCloseか確認
		allClosed := true
		for _, issue := range epicIssues {
			if issue.Status != "Close" {
				allClosed = false
				break
			}
		}

		// すべてのIssueがClosedの場合、Epicも閉じる
		if allClosed {
			old := epic.Status
			epic.Status = "Close"

			// 変更があった場合のみファイルを更新
			if old != epic.Status {
				statusChanged = true
				fmt.Printf("Epic更新: ID=%d, タイトル=%s (ステータス: %s → %s)\n",
					epic.ID, epic.Title, old, epic.Status)

				// 旧ファイル名を生成
				oldFilename := utils.GenerateFilename(epic.ID, old, epic.Title)
				oldFilePath := filepath.Join(cfg.EpicDir, oldFilename)

				// 更新されたEpicを書き込む
				if err := fileops.WriteEpic(cfg.EpicDir, epic); err != nil {
					return fmt.Errorf("Epicの更新に失敗しました: %w", err)
				}

				// 古いファイルを明示的に削除（同じIDの重複ファイルを避けるため）
				if err := os.Remove(oldFilePath); err != nil {
					fmt.Printf("警告: 古いEpicファイルの削除に失敗しました: %v\n", err)
					// 削除に失敗しても進める
				}
			}
		}
	}

	// ステータスが変更された場合のみファイル名を更新する
	if statusChanged {
		fmt.Println("Epicステータス変更によりファイル名を更新します...")
		err = RenameCommand(cfg)
		if err != nil {
			return fmt.Errorf("ファイル名の更新に失敗しました: %w", err)
		}
	}

	return nil
}

// SyncCommand - order.csvとIssueファイルの同期を行う
func SyncCommand(cfg *config.Config) error {
	fmt.Println("===== order.csvの同期を開始します... ====")

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

	// Epicステータスを関連するIssueに基づいて更新
	if err := UpdateEpicStatusBasedOnIssues(cfg); err != nil {
		return fmt.Errorf("Epicステータスの更新に失敗しました: %w", err)
	}

	// Epicステータス変更後に確実にファイル名を更新する
	if err := RenameCommand(cfg); err != nil {
		return fmt.Errorf("ファイル名の更新に失敗しました: %w", err)
	}

	fmt.Println("===== order.csvの同期が完了しました ====")

	return nil
}
