package parser

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/moai/instant-backlog/internal/models"
	"gopkg.in/yaml.v3"
)

// 正規表現でFront Matterを抽出する
// 注: (?s)はdotallモードで.が改行にもマッチする
// 問題があったため、改行の扱いを柔軟にするパターンに修正
var frontMatterRegex = regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)

// ParseIssueFile - Issueファイルを解析してIssue構造体を返す
func ParseIssueFile(filePath string) (*models.Issue, error) {
	// デバッグ: ファイルパスを表示
	fmt.Printf("ParseIssueFile: ファイルパス %s\n", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("ファイル読み込みエラー: %v\n", err)
		return nil, err
	}

	// デバッグ: ファイル内容の先頭部分を表示
	if len(content) > 0 {
		fmt.Printf("ファイル内容の先頭部分: %s\n", string(content[:min(100, len(content))]))
	}

	matches := frontMatterRegex.FindSubmatch(content)
	if len(matches) != 3 {
		fmt.Printf("無効なFront Matter形式が検出されました\n")
		return nil, &InvalidFrontMatterError{FilePath: filePath}
	}

	var issue models.Issue
	err = yaml.Unmarshal(matches[1], &issue)
	if err != nil {
		fmt.Printf("YAMLデシリアライズエラー: %v\n", err)
		return nil, err
	}

	// デバッグ: Unmarshalした値を表示
	fmt.Printf("パースされたIssue: ID=%d, Title=%s, Status=%s\n", issue.ID, issue.Title, issue.Status)

	// FrontMatterではない部分のコンテンツを設定
	issue.Content = strings.TrimSpace(string(matches[2]))

	return &issue, nil
}

// ParseEpicFile - Epicファイルを解析してEpic構造体を返す
func ParseEpicFile(filePath string) (*models.Epic, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	matches := frontMatterRegex.FindSubmatch(content)
	if len(matches) != 3 {
		return nil, &InvalidFrontMatterError{FilePath: filePath}
	}

	var epic models.Epic
	err = yaml.Unmarshal(matches[1], &epic)
	if err != nil {
		return nil, err
	}

	// FrontMatterではない部分のコンテンツを設定
	epic.Content = strings.TrimSpace(string(matches[2]))

	return &epic, nil
}

// GenerateMarkdown - Front Matterとコンテンツからマークダウンテキストを生成
func GenerateMarkdown(frontMatter interface{}, content string) ([]byte, error) {
	// Front Matterをマーシャリング
	frontMatterBytes, err := yaml.Marshal(frontMatter)
	if err != nil {
		return nil, err
	}

	// マークダウン形式で結合
	var buffer bytes.Buffer
	buffer.WriteString("---\n")
	buffer.Write(frontMatterBytes)
	buffer.WriteString("---\n\n")
	buffer.WriteString(content)

	return buffer.Bytes(), nil
}

// カスタムエラー型
type InvalidFrontMatterError struct {
	FilePath string
}

func (e *InvalidFrontMatterError) Error() string {
	return "無効なFront Matter形式: " + e.FilePath
}

// min は2つの整数の小さい方を返す
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
