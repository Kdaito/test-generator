package jest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kdaito/test-generator/internal/lib"
	proceccer "github.com/Kdaito/test-generator/internal/proceccer/csv"
)

func Generate(taregetRoot, testCaseRoot string) (int, error) {
	// TODO 以下二つは並行処理にしたい気もする
	targets, err := scanTarget(taregetRoot)

	if err != nil {
		return 0, err
	}

	testCases, err := proceccer.ScanTestCase(testCaseRoot)

	if err != nil {
		return 0, err
	}

	count := 0

	for _, target := range targets {
		fileName := filepath.Base(target.Path)

		// 拡張子を取得する
		fileExt := filepath.Ext(fileName)

		// 拡張子を除いたファイル名を取得
		fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		// 作成するファイル名を生成
		// jsかtsかどうかは対象ファイルにあわせる
		testFileName := fileNameWithoutExt + `.test` + fileExt

		// 対象ファイルが存在するディレクトリを取得する
		parentDir := filepath.Dir(target.Path)

		// 作成するファイルのパスを生成
		testFilePath := filepath.Join(parentDir, testFileName)

		// ファイルを作成
		file, err := os.Create(testFilePath)
		if err != nil {
			return 0, fmt.Errorf("[Failed create test file] path: %s\n[error] %v", testFilePath, err)
		}

		defer file.Close()

		// テストファイルに書き込む内容を作成
		source := buildTestFileSource(target, testCases, fileNameWithoutExt)

		// libフォルダ以下の関数を使ってファイルに出力
		err = lib.WriteTestFile(testFilePath, source)

		if err != nil {
			return 0, err
		}

		count += 1
	}
	return count, nil
}