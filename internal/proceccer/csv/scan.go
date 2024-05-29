package proceccer

import (
	"encoding/csv"
	"fmt"
	"github.com/Kdaito/test-generator/internal/types"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ScanTestCase(root string) ([]*types.TestCase, error) {
	var testCases []*types.TestCase

	// 対象フォルダを探索する
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("[Failed access path] path: %s\n[error] %v", path, err)
		}

		// csvファイルを見つけた場合
		if !info.IsDir() && strings.HasSuffix(path, ".csv") {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("[Failed open path] path: %s\n[error] %v", path, err)
			}

			defer file.Close()

			// CSVリーダーを作成
			reader := csv.NewReader(file)

			// CSVを一行ずつ読み込む
			records, err := reader.ReadAll()
			if err != nil {
				return fmt.Errorf("[Failed reading CSV file] path: %s\n[error] %v", path, err)
			}

			for _, record := range records {
				targetFuncName := record[0]
				testName := record[1]

				var checkList []string

				for _, v := range strings.Split(record[2], ",") {
					checkList = append(checkList, v)
				}

				testCases = append(testCases, &types.TestCase{
					TargetFuncName: targetFuncName,
					TestName:       testName,
					CheckList:      checkList,
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return testCases, nil
}
