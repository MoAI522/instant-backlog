package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moai/instant-backlog/internal/commands"
	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/parser"
)

// 埋め込みテンプレート機能に関するテスト
// テンプレートディレクトリが見つからない場合、埋め込みテンプレートを使用する
func TestEmbeddedTemplateInit(t *testing.T) {
	// テスト用のディレクトリを作成
	tempDir := filepath.Join(os.TempDir(), "instant-backlog-embed-test")
	defer os.RemoveAll(tempDir) // テスト終了後にクリーンアップ

	// テスト用ディレクトリが存在する場合は削除して再作成
	if _, err := os.Stat(tempDir); err == nil {
		os.RemoveAll(tempDir)
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}

	// 設定を作成（存在しないテンプレートパスを指定）
	cfg := config.NewConfig()
	cfg.TemplatePath = filepath.Join(tempDir, "non-existent-template")

	// initコマンドを実行
	err := commands.InitCommand(cfg, tempDir)
	if err != nil {
		t.Fatalf("initコマンドの実行に失敗しました: %v", err)
	}

	// 必要なディレクトリとファイルが作成されているか確認
	projectsDir := filepath.Join(tempDir, "projects")
	epicDir := filepath.Join(projectsDir, "epic")
	issuesDir := filepath.Join(projectsDir, "issues")
	orderCSVPath := filepath.Join(projectsDir, "order.csv")
	readmePath := filepath.Join(projectsDir, "README.md")

	// ディレクトリの確認
	if !dirExists(projectsDir) {
		t.Errorf("projectsディレクトリが作成されていません: %s", projectsDir)
	}
	if !dirExists(epicDir) {
		t.Errorf("epicディレクトリが作成されていません: %s", epicDir)
	}
	if !dirExists(issuesDir) {
		t.Errorf("issuesディレクトリが作成されていません: %s", issuesDir)
	}

	// 必須ファイルの確認
	if !fileExists(orderCSVPath) {
		t.Errorf("order.csvファイルが作成されていません: %s", orderCSVPath)
	}
	if !fileExists(readmePath) {
		t.Errorf("README.mdファイルが作成されていません: %s", readmePath)
	}

	// テンプレートが正しく展開されているか確認
	// Epicファイルの確認
	epicFiles, err := filepath.Glob(filepath.Join(epicDir, "*.md"))
	if err != nil {
		t.Fatalf("Epicファイルの検索に失敗しました: %v", err)
	}
	if len(epicFiles) == 0 {
		t.Errorf("埋め込みテンプレートからEpicファイルが展開されていません")
	} else {
		// 少なくとも1つのEpicファイルがあれば内容を確認
		epicFilePath := epicFiles[0]
		epic, err := parser.ParseEpicFile(epicFilePath)
		if err != nil {
			t.Fatalf("Epicファイルの解析に失敗しました: %v", err)
		}
		if epic.ID <= 0 || epic.Title == "" || epic.Status == "" {
			t.Errorf("Epicファイルの内容が不正です: %+v", epic)
		}
	}

	// Issueファイルの確認
	issueFiles, err := filepath.Glob(filepath.Join(issuesDir, "*.md"))
	if err != nil {
		t.Fatalf("Issueファイルの検索に失敗しました: %v", err)
	}
	if len(issueFiles) == 0 {
		t.Errorf("埋め込みテンプレートからIssueファイルが展開されていません")
	} else {
		// 少なくとも1つのIssueファイルがあれば内容を確認
		issueFilePath := issueFiles[0]
		issue, err := parser.ParseIssueFile(issueFilePath)
		if err != nil {
			t.Fatalf("Issueファイルの解析に失敗しました: %v", err)
		}
		if issue.ID <= 0 || issue.Title == "" || issue.Status == "" || issue.Epic <= 0 || issue.Estimate <= 0 {
			t.Errorf("Issueファイルの内容が不正です: %+v", issue)
		}
	}

	// order.csvの確認
	orderItems, err := parser.ReadOrderCSV(orderCSVPath)
	if err != nil {
		t.Fatalf("order.csvの読み込みに失敗しました: %v", err)
	}
	if len(orderItems) == 0 {
		t.Errorf("order.csvが空です")
	} else {
		// 少なくとも1つのアイテムがあれば内容を確認
		orderItem := orderItems[0]
		if orderItem.ID <= 0 || orderItem.Title == "" || orderItem.Epic <= 0 || orderItem.Estimate <= 0 {
			t.Errorf("order.csvの内容が不正です: %+v", orderItem)
		}
	}
}

// テンプレート優先順位のテスト
// 設定指定のテンプレート > 実行ファイルと同じディレクトリのテンプレート > リポジトリのテンプレート > 埋め込みテンプレート
func TestTemplatePriorityOrder(t *testing.T) {
	// テスト用のディレクトリを作成
	tempDir := filepath.Join(os.TempDir(), "instant-backlog-priority-test")
	defer os.RemoveAll(tempDir) // テスト終了後にクリーンアップ

	// テスト用ディレクトリが存在する場合は削除して再作成
	if _, err := os.Stat(tempDir); err == nil {
		os.RemoveAll(tempDir)
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}

	// カスタムテンプレートディレクトリを作成
	customTemplateDir := filepath.Join(tempDir, "custom-template", "projects")
	if err := os.MkdirAll(filepath.Join(customTemplateDir, "epic"), 0755); err != nil {
		t.Fatalf("カスタムテンプレートディレクトリの作成に失敗しました: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(customTemplateDir, "issues"), 0755); err != nil {
		t.Fatalf("カスタムテンプレートディレクトリの作成に失敗しました: %v", err)
	}

	// カスタムREADMEファイルを作成（これでカスタムテンプレートが使用されたことを確認）
	customReadmeContent := "# カスタムテンプレートからのREADME"
	if err := os.WriteFile(filepath.Join(customTemplateDir, "README.md"), []byte(customReadmeContent), 0644); err != nil {
		t.Fatalf("カスタムREADMEファイルの作成に失敗しました: %v", err)
	}

	// 設定を作成（カスタムテンプレートパスを指定）
	cfg := config.NewConfig()
	cfg.TemplatePath = customTemplateDir

	// initコマンドを実行
	err := commands.InitCommand(cfg, tempDir)
	if err != nil {
		t.Fatalf("initコマンドの実行に失敗しました: %v", err)
	}

	// READMEの内容が期待通りか確認（カスタムテンプレートが使用されていることを確認）
	readmePath := filepath.Join(tempDir, "projects", "README.md")
	if !fileExists(readmePath) {
		t.Errorf("README.mdファイルが作成されていません: %s", readmePath)
		return
	}

	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("README.mdファイルの読み込みに失敗しました: %v", err)
	}

	if string(readmeContent) != customReadmeContent {
		t.Errorf("カスタムテンプレートが使用されていません。期待値: %s, 実際: %s",
			customReadmeContent, string(readmeContent))
	}
}
