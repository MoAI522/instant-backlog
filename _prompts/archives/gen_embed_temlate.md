file-system を使用して、ファイルを直接読み書きして作業を進めてください。
回答は日本語で行ってください。
プロジェクトのディレクトリは"\\wsl.localhost\Ubuntu-22.04\home\moai\instant-backlog"です
まずプロジェクトの全容を把握してください。

このプロジェクトに対し、以下の方針で仕様変更の開発を行います。

````
# instant-backlogのテンプレート埋め込み実装方針

## 現状の課題

現在のinstant-backlogは、初期化時にテンプレートファイルを以下の優先順位で検索しています：

1. 設定ファイルで指定されたテンプレートパス
2. 実行ファイルと同じディレクトリの`template/projects`
3. リポジトリルートの`template/projects`

この方法では、ユーザーが自分でテンプレートを用意しなくても良い柔軟性がありますが、デフォルトテンプレートをバイナリに含めていないため、リポジトリに含まれるテンプレートが必ず必要となります。

## 実装目標

- テンプレートをバイナリに埋め込み、ユーザーが何も準備せずに`init`コマンドを実行できるようにする
- 現在のカスタマイズ性（ユーザー独自のテンプレートを使用可能）を維持する
- リポジトリ内のテンプレートファイルをそのまま維持し、開発時のカスタマイズを容易にする

## 実装方針

Go 1.16から導入された標準の`go:embed`機能を使用して、テンプレートファイルをバイナリに埋め込みます。

### 1. 埋め込みテンプレート用のパッケージ作成

```go
// internal/embedtemplate/template.go

package embedtemplate

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed ../../template/projects
var TemplateFS embed.FS

// ExtractTemplate は埋め込みテンプレートを指定したディレクトリに展開します
func ExtractTemplate(targetDir string) error {
	// 埋め込みファイルシステムから template/projects ディレクトリを取得
	projectsDir, err := fs.Sub(TemplateFS, "template/projects")
	if err != nil {
		return fmt.Errorf("埋め込みテンプレートにアクセスできません: %v", err)
	}

	// 再帰的にファイルを抽出
	return fs.WalkDir(projectsDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ターゲットパスを作成
		targetPath := filepath.Join(targetDir, path)

		// ディレクトリの場合は作成
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// .gitkeepファイルはスキップ
		if filepath.Base(path) == ".gitkeep" {
			return nil
		}

		// ファイルを読み込み
		data, err := projectsDir.ReadFile(path)
		if err != nil {
			return fmt.Errorf("ファイル読み込みエラー %s: %v", path, err)
		}

		// ファイルを書き込み
		return os.WriteFile(targetPath, data, 0644)
	})
}
````

### 2. init.go の修正

既存のテンプレート探索ロジックを維持しつつ、外部テンプレートが見つからない場合は埋め込みテンプレートを使用するよう修正します。

```go
// internal/commands/init.go の修正部分

import (
	// 既存のインポート
	"github.com/moai/instant-backlog/internal/embedtemplate"
)

// InitCommand - プロジェクトを初期化するコマンド
func InitCommand(cfg *config.Config, projectPath string) error {
	// 既存のコード（プロジェクトパスの設定など）...

	// テンプレートパスを取得
	var templateDir string
	templateFound := false

	// 設定からテンプレートパスを取得
	if cfg.TemplatePath != "" {
		templateDir = cfg.TemplatePath
		if _, err := os.Stat(templateDir); err == nil {
			templateFound = true
		}
	}

	// 設定に指定がなく、テンプレートが見つからない場合は既存のロジックを使用
	if !templateFound {
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("実行ファイルのパスを取得できません: %v", err)
		}

		// 実行ファイルと同じディレクトリにあるtemplate/projectsディレクトリを使用
		templateDir = filepath.Join(filepath.Dir(execPath), "template", "projects")

		if _, err := os.Stat(templateDir); err == nil {
			templateFound = true
		}
	}

	// リポジトリのrootディレクトリを探索
	if !templateFound {
		rootDir := findRepositoryRoot()
		if rootDir != "" {
			templateDir = filepath.Join(rootDir, "template", "projects")
			if _, err := os.Stat(templateDir); err == nil {
				templateFound = true
			}
		}
	}

	// 外部テンプレートが見つかった場合はそれをコピー
	if templateFound {
		// テンプレートファイルをコピー
		if err := copyDir(templateDir, projectsDir); err != nil {
			return fmt.Errorf("テンプレートファイルのコピーに失敗しました: %v", err)
		}
	} else {
		// 埋め込みテンプレートを使用
		if err := embedtemplate.ExtractTemplate(projectsDir); err != nil {
			return fmt.Errorf("埋め込みテンプレートの展開に失敗しました: %v", err)
		}
	}

	// 残りの既存コード...
	return nil
}
```

## メリット

1. **簡単な初期化**: ユーザーは環境に関係なく、すぐに`init`コマンドを実行できる
2. **柔軟性の維持**: カスタムテンプレートの優先度が高いため、カスタマイズ性は維持される
3. **一貫性**: バイナリに含まれるテンプレートとリポジトリのテンプレートが同じ
4. **標準ライブラリの使用**: サードパーティ依存なしで実装可能

## 注意点

1. **Go 1.16 以上が必要**: `go:embed`を使うには Go 1.16 以上が必要
2. **バイナリサイズの増加**: テンプレートファイルの分だけバイナリサイズが増加する
3. **ビルド時の考慮**: テンプレートが更新されたら再ビルドが必要

## 実装ステップ

1. `internal/embedtemplate` パッケージを作成
2. `init.go` のテンプレート探索ロジックを修正
3. ビルドプロセスを更新して埋め込みテンプレートを含める

これらの変更により、ユーザーエクスペリエンスを向上させつつ、既存の柔軟性も維持できます。

```

以下のタスクを遂行してください。

- この機能を実装してください。
- この機能を使用するユーザーストーリーを考え、統合テストを追加してください。
- 既存のテストのユーザーストーリーを確認し、必要に応じてテストごと修正してください。
- README.mdを更新してください。
- ONBOARDING.mdを更新してください。

edit_fileでファイルを編集する際は、まず変更内容をartifactに記載し、その後でその内容をedit_fileで反映してください。
```
