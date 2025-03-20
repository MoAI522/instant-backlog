package utils

import (
	"fmt"

	"github.com/moai/instant-backlog/internal/models"
)

// ValidateIssue - Issueのバリデーションを行う
func ValidateIssue(issue *models.Issue) error {
	if issue.ID <= 0 {
		return fmt.Errorf("ID は正の整数でなければなりません")
	}
	
	if issue.Title == "" {
		return fmt.Errorf("タイトルは必須です")
	}
	
	if issue.Status != "Open" && issue.Status != "Close" {
		return fmt.Errorf("ステータスは 'Open' または 'Close' でなければなりません")
	}
	
	if issue.Epic <= 0 {
		return fmt.Errorf("Epic ID は正の整数でなければなりません")
	}
	
	if issue.Estimate < 0 {
		return fmt.Errorf("見積もりポイントは非負の整数でなければなりません")
	}
	
	return nil
}

// ValidateEpic - Epicのバリデーションを行う
func ValidateEpic(epic *models.Epic) error {
	if epic.ID <= 0 {
		return fmt.Errorf("ID は正の整数でなければなりません")
	}
	
	if epic.Title == "" {
		return fmt.Errorf("タイトルは必須です")
	}
	
	if epic.Status != "Open" && epic.Status != "Close" {
		return fmt.Errorf("ステータスは 'Open' または 'Close' でなければなりません")
	}
	
	return nil
}
