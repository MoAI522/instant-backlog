package config

import (
	"os"
	"path/filepath"
)

// Config - アプリケーション設定を表す構造体
type Config struct {
	ProjectsDir string
	EpicDir     string
	IssuesDir   string
	OrderCSV    string
	// テンプレートディレクトリのパス
	TemplatePath string
}

// NewConfig - デフォルト設定で設定構造体を作成
func NewConfig() *Config {
	// ベースディレクトリの取得
	baseDir, err := os.Getwd()
	if err != nil {
		baseDir = "."
	}

	projectsDir := filepath.Join(baseDir, "projects")

	return &Config{
		ProjectsDir: projectsDir,
		EpicDir:     filepath.Join(projectsDir, "epic"),
		IssuesDir:   filepath.Join(projectsDir, "issues"),
		OrderCSV:    filepath.Join(projectsDir, "order.csv"),
	}
}
