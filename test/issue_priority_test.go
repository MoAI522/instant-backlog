// test/issue_priority_test.go
package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/moai/instant-backlog/internal/commands"
	"github.com/moai/instant-backlog/internal/models"
	"github.com/moai/instant-backlog/internal/parser"
)

/**
 * ユーザーストーリー3：Issueの優先順位を管理できること
 *
 * ユーザーはIssue（タスク）の優先順位を管理したい。
 * 優先順位はCSVファイルで管理され、以下の要件がある：
 * - CSVにはオープンなIssueのみが含まれる
 * - CSVの順序がIssueの優先順位を表す
 * - クローズされたIssueは自動的にCSVから削除される
 * - syncコマンドを実行することでCSVとIssueの整合性が保たれる
 *
 * この機能により、ユーザーは容易に作業の優先順位を把握し、管理できる。
 */
func TestManageIssuePriorities(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のデータを作成
	createTestEpic(t, cfg, 1, "テストエピック", "Open")
	createTestIssue(t, cfg, 1, "高優先度タスク", "Open", 1, 5)
	createTestIssue(t, cfg, 2, "中優先度タスク", "Open", 1, 3)
	createTestIssue(t, cfg, 3, "低優先度タスク", "Open", 1, 1)
	createTestIssue(t, cfg, 4, "クローズタスク", "Close", 1, 2)

	// テスト用のCSVを作成（優先度順）
	initialOrderItems := []models.OrderCSVItem{
		{ID: 1, Title: "高優先度タスク", Epic: 1, Estimate: 5},
		{ID: 2, Title: "中優先度タスク", Epic: 1, Estimate: 3},
		{ID: 3, Title: "低優先度タスク", Epic: 1, Estimate: 1},
		{ID: 4, Title: "クローズタスク", Epic: 1, Estimate: 2}, // 初期CSVにはクローズタスクも含まれている
	}
	createTestOrderCSV(t, cfg, initialOrderItems)

	// syncコマンドを実行
	err := commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// 更新後のCSVを読み込む
	updatedOrderItems, err := parser.ReadOrderCSV(cfg.OrderCSV)
	if err != nil {
		t.Fatalf("order.csvの読み込みに失敗しました: %v", err)
	}

	// クローズされたタスクがCSVから削除されているか確認
	for _, item := range updatedOrderItems {
		if item.ID == 4 {
			t.Errorf("クローズされたタスク(ID=4)がorder.csvから削除されていません")
		}
	}

	// 残りのオープンタスクが正しく含まれているか確認
	expectedIDs := map[int]bool{1: true, 2: true, 3: true}
	for _, item := range updatedOrderItems {
		if !expectedIDs[item.ID] {
			t.Errorf("予期しないタスクID %d がorder.csvに含まれています", item.ID)
		}
		delete(expectedIDs, item.ID)
	}

	if len(expectedIDs) > 0 {
		var missingIDs []string
		for id := range expectedIDs {
			missingIDs = append(missingIDs, fmt.Sprintf("%d", id))
		}
		t.Errorf("期待されるタスクIDs %s がorder.csvに含まれていません", strings.Join(missingIDs, ", "))
	}
}

/**
 * 優先順位の変更と保持のテスト
 *
 * CSVファイルで優先順位を変更した際に、その順序が保持されることを検証します。
 * syncコマンド実行後も、明示的に変更された優先順位が維持されることを確認します。
 */
func TestPriorityOrderPreservation(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のデータを作成
	createTestEpic(t, cfg, 1, "優先順位テストエピック", "Open")
	createTestIssue(t, cfg, 1, "タスクA", "Open", 1, 2)
	createTestIssue(t, cfg, 2, "タスクB", "Open", 1, 3)
	createTestIssue(t, cfg, 3, "タスクC", "Open", 1, 1)

	// 初期CSVを順序を変えて作成（B->C->A の順）
	initialOrderItems := []models.OrderCSVItem{
		{ID: 2, Title: "タスクB", Epic: 1, Estimate: 3},
		{ID: 3, Title: "タスクC", Epic: 1, Estimate: 1},
		{ID: 1, Title: "タスクA", Epic: 1, Estimate: 2},
	}
	createTestOrderCSV(t, cfg, initialOrderItems)

	// syncコマンドを実行
	err := commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// 更新後のCSVを読み込む
	updatedOrderItems, err := parser.ReadOrderCSV(cfg.OrderCSV)
	if err != nil {
		t.Fatalf("order.csvの読み込みに失敗しました: %v", err)
	}

	// 元の優先順位が保持されているか確認
	if len(updatedOrderItems) != 3 {
		t.Fatalf("CSVのアイテム数が不正です。期待値: 3, 実際: %d", len(updatedOrderItems))
	}

	// 期待する順序
	expectedOrder := []int{2, 3, 1}
	for i, expectedID := range expectedOrder {
		if updatedOrderItems[i].ID != expectedID {
			t.Errorf("優先順位が保持されていません。位置 %d では 期待値: %d, 実際: %d",
				i, expectedID, updatedOrderItems[i].ID)
		}
	}
}

/**
 * 新規Issueが追加された際の優先順位処理のテスト
 *
 * 新しいIssueが作成された際に、それが適切にCSVファイルに追加されることを検証します。
 * 既存の優先順位を維持しながら、新規Issueが末尾に追加されることを確認します。
 */
func TestNewIssuePriorityHandling(t *testing.T) {
	// テスト環境をセットアップ
	cfg, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// テスト用のデータを初期作成
	createTestEpic(t, cfg, 1, "優先順位テストエピック", "Open")
	createTestIssue(t, cfg, 1, "既存タスク1", "Open", 1, 3)
	createTestIssue(t, cfg, 2, "既存タスク2", "Open", 1, 2)

	// 初期CSVを作成
	initialOrderItems := []models.OrderCSVItem{
		{ID: 1, Title: "既存タスク1", Epic: 1, Estimate: 3},
		{ID: 2, Title: "既存タスク2", Epic: 1, Estimate: 2},
	}
	createTestOrderCSV(t, cfg, initialOrderItems)

	// 新しいIssueを追加
	createTestIssue(t, cfg, 3, "新規タスク", "Open", 1, 5)

	// syncコマンドを実行
	err := commands.SyncCommand(cfg)
	if err != nil {
		t.Fatalf("syncコマンドの実行に失敗しました: %v", err)
	}

	// 更新後のCSVを読み込む
	updatedOrderItems, err := parser.ReadOrderCSV(cfg.OrderCSV)
	if err != nil {
		t.Fatalf("order.csvの読み込みに失敗しました: %v", err)
	}

	// 新規Issueが追加され、既存の順序が維持されているか確認
	if len(updatedOrderItems) != 3 {
		t.Fatalf("CSVのアイテム数が不正です。期待値: 3, 実際: %d", len(updatedOrderItems))
	}

	// 元の順序が維持され、新規タスクが末尾に追加されているか確認
	expectedIDs := []int{1, 2, 3}
	for i, expectedID := range expectedIDs {
		if updatedOrderItems[i].ID != expectedID {
			t.Errorf("位置 %d での ID が不正です。期待値: %d, 実際: %d",
				i, expectedID, updatedOrderItems[i].ID)
		}
	}
}
