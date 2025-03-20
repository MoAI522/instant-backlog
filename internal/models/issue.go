package models

// Issue - スクラムバックログの課題を表す構造体
type Issue struct {
	ID       int    `yaml:"id"`
	Title    string `yaml:"title"`
	Status   string `yaml:"status"` // "Open" または "Close"
	Epic     int    `yaml:"epic"`   // 関連するEpicのID
	Estimate int    `yaml:"estimate"`
	Content  string `yaml:"-"` // Front Matterではない部分のコンテンツ
}

// OrderCSVItem - order.csvに保存される項目
type OrderCSVItem struct {
	ID       int    `csv:"id"`
	Title    string `csv:"title"`
	Epic     int    `csv:"epic"`
	Estimate int    `csv:"estimate"`
}
