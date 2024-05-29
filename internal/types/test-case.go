package types

// テストケースファイルの情報
type TestCase struct {
	TargetFuncName string
	TestName       string
	CheckList      []string
}

// プロジェクトのテスト対象関数情報
type TestTarget struct {
	TargetFuncNames []string
	Path string
}