package models

// Epic - スクラムバックログのエピックを表す構造体
type Epic struct {
	ID      int    `yaml:"id"`
	Title   string `yaml:"title"`
	Status  string `yaml:"status"` // "Open" または "Close"
	Content string `yaml:"-"`      // Front Matterではない部分のコンテンツ
}
