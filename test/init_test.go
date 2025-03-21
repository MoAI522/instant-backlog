// test/init_test.go
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/moai/instant-backlog/internal/commands"
	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/parser"
)

// initコマンドのテスト - プロジェクトを初期化できること
func TestInitCommand(t *testing.T) {
	// テスト用のディレクトリを作成
	tempDir := filepath.Join(os.TempDir(), "instant-backlog-init-test")
	defer os.RemoveAll(tempDir) // テスト終了後にクリーンアップ

	// テスト用ディレクトリが存在する場合は削除して再作成
	if _, err := os.Stat(tempDir); err == nil {
		os.RemoveAll(tempDir)
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}

	// 設定を作成
	cfg := config.NewConfig()

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

	// テンプレートが正しくコピーされているか確認
	// Epicファイルの確認
	epicFiles, err := filepath.Glob(filepath.Join(epicDir, "*.md"))
	if err != nil {
		t.Fatalf("Epicファイルの検索に失敗しました: %v", err)
	}
	if len(epicFiles) == 0 {
		t.Errorf("Epicファイルがコピーされていません")
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
		t.Errorf("Issueファイルがコピーされていません")
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

// ディレクトリの存在を確認するヘルパー関数
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// 空のディレクトリに対するinitコマンドテスト
func TestInitCommandOnEmptyDirectory(t *testing.T) {
	// テスト用のディレクトリを作成（何も入っていない状態）
	emptyDir := filepath.Join(os.TempDir(), "instant-backlog-empty-test")
	defer os.RemoveAll(emptyDir)

	// テスト用ディレクトリが存在する場合は削除して再作成
	if _, err := os.Stat(emptyDir); err == nil {
		os.RemoveAll(emptyDir)
	}
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}

	// 設定を作成
	cfg := config.NewConfig()

	// initコマンドを実行
	err := commands.InitCommand(cfg, emptyDir)
	if err != nil {
		t.Fatalf("空のディレクトリに対するinitコマンドの実行に失敗しました: %v", err)
	}

	// プロジェクト構造が正しく初期化されたか確認
	if !dirExists(filepath.Join(emptyDir, "projects")) ||
		!dirExists(filepath.Join(emptyDir, "projects", "epic")) ||
		!dirExists(filepath.Join(emptyDir, "projects", "issues")) ||
		!fileExists(filepath.Join(emptyDir, "projects", "order.csv")) ||
		!fileExists(filepath.Join(emptyDir, "projects", "README.md")) {
		t.Errorf("空のディレクトリにプロジェクト構造が正しく初期化されていません")
	}
}

// 引数なしのinitコマンドテスト（カレントディレクトリを使用）
func TestInitCommandWithoutArgument(t *testing.T) {
	// 元のカレントディレクトリを保存
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("カレントディレクトリの取得に失敗しました: %v", err)
	}
	defer os.Chdir(originalDir) // テスト終了時に元に戻す

	// テスト用の一時ディレクトリを作成してそこに移動
	tempDir := filepath.Join(os.TempDir(), "instant-backlog-current-test")
	defer os.RemoveAll(tempDir)

	if _, err := os.Stat(tempDir); err == nil {
		os.RemoveAll(tempDir)
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}

	// テスト用のテンプレートディレクトリを作成
	templateDir := filepath.Join(tempDir, "template", "projects")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("テンプレート用ディレクトリの作成に失敗しました: %v", err)
	}

	// テスト用のテンプレートファイルを作成
	epicDir := filepath.Join(templateDir, "epic")
	if err := os.MkdirAll(epicDir, 0755); err != nil {
		t.Fatalf("epic用ディレクトリの作成に失敗しました: %v", err)
	}

	issuesDir := filepath.Join(templateDir, "issues")
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		t.Fatalf("issues用ディレクトリの作成に失敗しました: %v", err)
	}

	// 必要最小限のテンプレートファイルを作成
	err = os.WriteFile(filepath.Join(templateDir, "README.md"), []byte("# テスト用README"), 0644)
	if err != nil {
		t.Fatalf("テンプレートファイルの作成に失敗しました: %v", err)
	}
	err = os.WriteFile(filepath.Join(templateDir, "order.csv"), []byte("id,title,epic,estimate\n"), 0644)
	if err != nil {
		t.Fatalf("テンプレートファイルの作成に失敗しました: %v", err)
	}

	// テンプレートの基本ファイルを作成
	epicContent := `---
id: 1
title: "テストエピック"
status: "Open"
---

テスト用エピック
`
	err = os.WriteFile(filepath.Join(epicDir, "1_O_テストエピック.md"), []byte(epicContent), 0644)
	if err != nil {
		t.Fatalf("テンプレートエピックファイルの作成に失敗しました: %v", err)
	}

	issueContent := `---
id: 1
title: "テストタスク"
status: "Open"
epic: 1
estimate: 3
---

テスト用タスク
`
	err = os.WriteFile(filepath.Join(issuesDir, "1_O_テストタスク.md"), []byte(issueContent), 0644)
	if err != nil {
		t.Fatalf("テンプレートイシューファイルの作成に失敗しました: %v", err)
	}

	// テストディレクトリに移動
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("テストディレクトリへの移動に失敗しました: %v", err)
	}

	// 設定を作成（テンプレートパスを明示的に設定）
	cfg := config.NewConfig()
	cfg.TemplatePath = templateDir // 設定にテンプレートパスを追加

	// 引数なしでinitコマンドを実行（カレントディレクトリが使用される）
	err = commands.InitCommand(cfg, "")
	if err != nil {
		t.Fatalf("引数なしのinitコマンドの実行に失敗しました: %v", err)
	}

	// プロジェクト構造が正しく初期化されたか確認
	if !dirExists(filepath.Join(tempDir, "projects")) ||
		!dirExists(filepath.Join(tempDir, "projects", "epic")) ||
		!dirExists(filepath.Join(tempDir, "projects", "issues")) ||
		!fileExists(filepath.Join(tempDir, "projects", "order.csv")) ||
		!fileExists(filepath.Join(tempDir, "projects", "README.md")) {
		t.Errorf("カレントディレクトリにプロジェクト構造が正しく初期化されていません")
	}
}

// 既存のプロジェクトディレクトリに対するinitコマンドテスト（上書き確認）
func TestInitCommandOnExistingProject(t *testing.T) {
	// テスト用のディレクトリを作成
	existingDir := filepath.Join(os.TempDir(), "instant-backlog-existing-test")
	defer os.RemoveAll(existingDir)

	// テスト用ディレクトリが存在する場合は削除して再作成
	if _, err := os.Stat(existingDir); err == nil {
		os.RemoveAll(existingDir)
	}
	if err := os.MkdirAll(existingDir, 0755); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}

	// あらかじめプロジェクト構造の一部を作成
	if err := os.MkdirAll(filepath.Join(existingDir, "projects", "epic"), 0755); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(existingDir, "projects", "issues"), 0755); err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}

	// テスト用のオリジナルファイルを作成
	originalContent := "これはテスト用のオリジナルファイルです"
	originalFilePath := filepath.Join(existingDir, "projects", "original.txt")
	if err := os.WriteFile(originalFilePath, []byte(originalContent), 0644); err != nil {
		t.Fatalf("テスト用ファイルの作成に失敗しました: %v", err)
	}

	// 設定を作成
	cfg := config.NewConfig()

	// initコマンドを実行
	err := commands.InitCommand(cfg, existingDir)
	if err != nil {
		t.Fatalf("既存プロジェクトに対するinitコマンドの実行に失敗しました: %v", err)
	}

	// プロジェクト構造が正しく初期化されたか確認
	if !dirExists(filepath.Join(existingDir, "projects")) ||
		!dirExists(filepath.Join(existingDir, "projects", "epic")) ||
		!dirExists(filepath.Join(existingDir, "projects", "issues")) ||
		!fileExists(filepath.Join(existingDir, "projects", "order.csv")) ||
		!fileExists(filepath.Join(existingDir, "projects", "README.md")) {
		t.Errorf("既存プロジェクトディレクトリが正しく初期化されていません")
	}

	// オリジナルファイルが保持されているか確認（上書きされていないか）
	if !fileExists(originalFilePath) {
		t.Errorf("オリジナルファイルが削除されています: %s", originalFilePath)
	} else {
		content, err := os.ReadFile(originalFilePath)
		if err != nil {
			t.Fatalf("オリジナルファイルの読み込みに失敗しました: %v", err)
		}
		if string(content) != originalContent {
			t.Errorf("オリジナルファイルの内容が変更されています。期待値: %s, 実際: %s", originalContent, string(content))
		}
	}
}
