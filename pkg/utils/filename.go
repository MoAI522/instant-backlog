package utils

import (
	"fmt"
	"strings"
)

// GenerateFilename - 指定されたパラメータからファイル名を生成する
// 形式: {ID}_{O/C}_{タイトル}.md
func GenerateFilename(id int, status, title string) string {
	// ステータスの頭文字を取得
	statusChar := "O"
	if strings.ToLower(status) == "close" {
		statusChar = "C"
	}

	// タイトルをスペースからアンダースコアに変換
	safeTitle := strings.ReplaceAll(title, " ", "_")
	// 特殊文字を除去
	safeTitle = sanitizeString(safeTitle)

	return fmt.Sprintf("%d_%s_%s.md", id, statusChar, safeTitle)
}

// ParseFilename - ファイル名からID、ステータス、タイトルを抽出
func ParseFilename(filename string) (int, string, string, error) {
	// 拡張子を除去
	base := strings.TrimSuffix(filename, ".md")

	// 区切り文字で分割
	parts := strings.SplitN(base, "_", 3)
	if len(parts) < 3 {
		return 0, "", "", fmt.Errorf("無効なファイル名形式: %s", filename)
	}

	// IDを解析
	var id int
	_, err := fmt.Sscanf(parts[0], "%d", &id)
	if err != nil {
		return 0, "", "", fmt.Errorf("無効なID: %s", parts[0])
	}

	// ステータスを解析
	status := "Open"
	if parts[1] == "C" {
		status = "Close"
	}

	// タイトルを取得
	title := strings.ReplaceAll(parts[2], "_", " ")

	return id, status, title, nil
}

// sanitizeString - 特殊文字を除去して安全な文字列にする
func sanitizeString(s string) string {
	// ファイル名に使えない文字を置換
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	return replacer.Replace(s)
}
