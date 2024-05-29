package jest

import (
	"fmt"

	"github.com/Kdaito/test-generator/internal/lib"
	"github.com/Kdaito/test-generator/internal/types"
)

// buildTestFileSourceは、指定されたテストターゲットとテストケースに基づいてテストファイルのソースコードを生成します。
// 生成されたソースには、インポート文とテスト関数の実装が含まれます。
//
// パラメータ:
//   - targets: テスト対象の関数名を含むTestTarget構造体のポインタ。
//   - testCases: 対象関数のテストケースを含むTestCase構造体のポインタのスライス。
//   - targetFileName: テストファイル作成対象の関数が存在するファイル名
//
// 戻り値:
//   - 生成されたテストファイルのソースコード行を表す文字列のスライス。
//
// この関数は以下の手順を実行します:
//   1. 実装された関数とテストソースを格納するスライスを初期化します。
//   2. 対象関数名を繰り返し処理して、各関数のテスト関数ソースコードを生成します。
//   3. 実装された関数と提供された相対パスに基づいてインポート文を作成します。
//   4. インポート文とテストソースを1つのスライスに結合して返します。
func buildTestFileSource(targets *types.TestTarget, testCases []*types.TestCase, targetFileName string) []string {
	// テストを実装した関数
	impleFuncs := []string{}

	// テストの内容
	testSources := []string{}

	// FIXME このループの中も並行処理したい気持ちがあるk
	for _, funcName := range targets.TargetFuncNames {
		impleFunc, source := buildTestFuncSource(funcName, testCases)
		impleFuncs = append(impleFuncs, impleFunc)
		testSources = append(testSources, source...)
	}

	// インポート文の作成
	importSource := buildImportSource(impleFuncs, targetFileName)

	// 本文全体
	resource := []string{importSource, ""}
	resource = append(resource, testSources...)

	return resource
}

// buildImportSourceは、指定された関数名と相対パスに基づいてインポート文を生成します。
//
// パラメータ:
//   - impleFuncs: テストを実装する関数名のスライス。
//   - targetFileName: テストファイル作成対象の関数が存在するファイル名文字列。
//
// 戻り値:
//   - 生成されたインポート文を表す文字列。
//
// この関数は以下の手順を実行します:
//   1. 空の文字列impleFuncSourceを初期化します。
//   2. impleFuncsスライスの各関数名をフォーマットしてimpleFuncSourceに追加します。
//   3. 2で生成された関数名と相対パスを用いてインポート文を構築し、返します。
func buildImportSource(impleFuncs []string, targetFileName string) string {
	impleFuncSource := ``
	for i, impleFunc := range impleFuncs {
		if i > 0 {
			impleFuncSource += `, `
		}
		impleFuncSource += impleFunc
	}

	return `import { ` + impleFuncSource + ` } from './` + targetFileName + `';`
}

// buildTestFuncSourceは、指定された関数名とテストケースに基づいて、テストファイルのソースコードを生成します。
//
// パラメータ:
//   - funcName: テスト対象の関数名を表す文字列。
//   - testCases: テスト対象の関数に対するテストケースを含むTestCase構造体のポインタのスライス。
//
// 戻り値:
//   - テスト対象の関数名と生成されたテストファイルのソースコード行を表す文字列のスライス。
//     テストケースが存在しない場合、空の文字列とnilを返します。
//
// この関数は以下の手順を実行します:
//   1. 指定された関数名に対応するテストケースを取得します。
//   2. テストケースが存在しない場合、空の文字列とnilを返します。
//   3. テストケースが存在する場合、テストファイルの内容を生成します。
//   4. テストファイルの内容を文字列のスライスとして返します。
func buildTestFuncSource(funcName string, testCases []*types.TestCase) (string, []string) {
	// テストケースを取得する
	targetTestCase := lib.Find(testCases, func(testCase *types.TestCase) bool {
		return testCase.TargetFuncName == funcName
	})

	// テストケースがなかった場合
	if targetTestCase == nil {
		return "", nil
	}

	lines := []string{}

	// テストケースがあった場合、テストファイルの内容を生成する
	descStart := fmt.Sprintf(`describe('%s', () => {`, funcName)
	descEnd := `});`

	lines = append(lines, descStart)
	lines = append(lines, "")

	for i, checkItem := range targetTestCase.CheckList {
		if i != 0 {
			lines = append(lines, "")
		}

		checkItemStart := fmt.Sprintf(`  test('%s', () => {`, checkItem)
		checkItemEnd := `  });`
		lines = append(lines, checkItemStart)
		lines = append(lines, `    `)
		lines = append(lines, checkItemEnd)
	}

	lines = append(lines, descEnd)
	lines = append(lines, "")

	return targetTestCase.TargetFuncName, lines
}
