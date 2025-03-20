package parser

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/moai/instant-backlog/internal/models"
)

// ReadOrderCSV - order.csvを読み込みOrderCSVItemスライスを返す
func ReadOrderCSV(filePath string) ([]models.OrderCSVItem, error) {
	// ファイルが存在しない場合は空のスライスを返す
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []models.OrderCSVItem{}, nil
	}

	// CSVファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// CSVをパース
	var orderItems []models.OrderCSVItem
	if err := gocsv.UnmarshalFile(file, &orderItems); err != nil {
		return nil, err
	}

	return orderItems, nil
}

// WriteOrderCSV - OrderCSVItemスライスをorder.csvに書き込む
func WriteOrderCSV(filePath string, orderItems []models.OrderCSVItem) error {
	// CSVファイルを作成/上書き
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// CSVにマーシャリング
	return gocsv.MarshalFile(&orderItems, file)
}
