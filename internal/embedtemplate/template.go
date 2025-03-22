// internal/embedtemplate/template.go

package embedtemplate

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed "template/projects"
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
		data, err := fs.ReadFile(projectsDir, path)
		if err != nil {
			return fmt.Errorf("ファイル読み込みエラー %s: %v", path, err)
		}

		// ファイルを書き込み
		return os.WriteFile(targetPath, data, 0644)
	})
}
