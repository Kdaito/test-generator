package jest

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Kdaito/test-generator/internal/types"
)

// 対象フォルダから、テスト対象オブジェクトを生成する
func scanTarget(root string) ([]*types.TestTarget, error) {
	var testTargets []*types.TestTarget

	// 対象フォルダを探索する
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("[Failed access path] path: %s \n error: %v", path, err)
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".ts")) && !strings.Contains(path, ".test.") && !strings.Contains(path, "node_modules") {
			fmt.Println(path)
			res, err := scanTargetMethods(path)

			if err != nil {
				return err
			}

			testTargets = append(testTargets, res)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return testTargets, nil
}


// 対象ファイルから、テスト実装対象関数を探す
func scanTargetMethods(path string) (*types.TestTarget, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("[Failed open target file] path: %s \n[error] %v", path, err)
	}

	defer file.Close()

	// テスト作成対象関数
	var targetFunctionNames []string

	scanner := bufio.NewScanner(file)
	var prevLine string

	functionRegex := regexp.MustCompile(`function\s+(\w+)\s*\(|const\s+(\w+)\s*=\s*\(`)

	// テスト対象のメソッドを探し、メソッド名を取得する
	for scanner.Scan() {
		line := scanner.Text()

		// 関数を発見した場合
		if functionRegex.MatchString(line) {
			matches := functionRegex.FindStringSubmatch(line)

			funcName := ""

			if matches[1] != "" {
				funcName = matches[1] // 通常の関数名
			} else if matches[2] != "" {
				funcName = matches[2] // アロー関数名
			}

			// 関数の前文のコメントに `// Test`があった場合、実装対象
			if strings.TrimSpace(prevLine) == "// Test" {
				targetFunctionNames = append(targetFunctionNames, funcName)
			}
		}
		prevLine = line
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("[Failed scan target file] path: %s \n[error]: %v", path, err)
	}

	return &types.TestTarget{TargetFuncNames: targetFunctionNames, Path: path}, nil
}