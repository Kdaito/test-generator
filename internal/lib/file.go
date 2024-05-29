package lib

import (
	"bufio"
	"fmt"
	"os"
)

func WriteTestFile(filePath string, lines []string) error {
	// ファイルを作成
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("[Failed create test file] path: %s\n[error] %v", filePath, err)
	}

	// 関数終了時にファイルを閉じる
	defer file.Close()

	// bufio.Writerを使用してファイルに書き込む
	writer := bufio.NewWriter(file)

	// 1行ずつファイルに書き込む
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("[Failed write test file line] path: %s, line: %s\n[error] %v", filePath, line, err)
		}
	}

	// バッファをフラッシュして書き込みを確定させる
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("[Failed flash buffer after create test file] path: %s\n[error] %v", filePath, err)
	}

	return nil
}