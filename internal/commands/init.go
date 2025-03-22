package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/moai/instant-backlog/internal/config"
	"github.com/moai/instant-backlog/internal/embedtemplate"
)

// InitCommand - プロジェクトを初期化するコマンド
func InitCommand(cfg *config.Config, projectPath string) error {
	// プロジェクトパスを設定
	targetPath := projectPath
	if targetPath == "" {
		// プロジェクトパスが指定されていない場合は現在のディレクトリを使用
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("現在のディレクトリを取得できません: %v", err)
		}
		targetPath = currentDir
	}

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("ディレクトリを作成できません: %v", err)
	}

	// プロジェクト内のディレクトリ作成
	projectsDir := filepath.Join(targetPath, "projects")
	epicDir := filepath.Join(projectsDir, "epic")
	issuesDir := filepath.Join(projectsDir, "issues")

	// ディレクトリを作成
	if err := os.MkdirAll(epicDir, 0755); err != nil {
		return fmt.Errorf("epicディレクトリを作成できません: %v", err)
	}
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		return fmt.Errorf("issuesディレクトリを作成できません: %v", err)
	}

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

	fmt.Printf("プロジェクトを初期化しました: %s\n", targetPath)
	fmt.Println("次のステップ:")
	fmt.Println("1. projects/README.mdファイルを参照して運用方法を確認")
	fmt.Println("2. 必要に応じてEpicとIssueを編集")
	fmt.Println("3. 'instant-backlog watch projects'コマンドを実行して自動監視を開始")

	return nil
}

// findRepositoryRoot - リポジトリのルートディレクトリを探索
func findRepositoryRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	// 最大10階層まで遡って.gitディレクトリを探す
	dir := currentDir
	for i := 0; i < 10; i++ {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return dir
		}

		// 親ディレクトリに移動
		parent := filepath.Dir(dir)
		if parent == dir {
			// これ以上遡れない
			break
		}
		dir = parent
	}

	return ""
}

// copyDir - ディレクトリを再帰的にコピー
func copyDir(src, dst string) error {
	// 送信元ディレクトリを開く
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 送信元がディレクトリでない場合はエラー
	if !srcInfo.IsDir() {
		return fmt.Errorf("%s はディレクトリではありません", src)
	}

	// 送信先ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 送信元ディレクトリのエントリを列挙
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 各エントリに対して処理
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// エントリがディレクトリの場合は再帰的にコピー
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// .gitkeepファイルはスキップ
			if strings.HasSuffix(entry.Name(), ".gitkeep") {
				continue
			}

			// ファイルの場合はコピー
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile - ファイルをコピー
func copyFile(src, dst string) error {
	// 送信元ファイルを開く
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// 送信元ファイルの情報を取得
	srcInfo, err := in.Stat()
	if err != nil {
		return err
	}

	// 送信先ファイルを作成
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	// コピー実行
	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}
