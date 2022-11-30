package internal

import "go/token"

// TestFile テスト対象ファイルのASTから抽出した値を格納する構造体
type TestFile struct {
	// テスト対象のメソッドを持つ構造体のフィールド情報を管理
	FieldMap map[string]*FieldInfo
	// 各テスト対象のメソッドのテストケース一覧を管理
	TargetMethodTesCasesMap map[string][]*TestCase
}

// FieldInfo フィールド情報
type FieldInfo struct {
	// インタフェースか否か
	IsInterface bool
	// パッケージ名
	PackageName string
	// 型名
	TypeName string
	// テンプレートのパラメータに用いる型名
	UpperCamelCaseTypeName string
}

// TestCase テストケース
type TestCase struct {
	// テストケースの分岐点となるif文の行数
	Line int
	// 正常系か
	IsSuccessPattern bool
	// 依存しているメソッド一覧(自身のメソッド or mock化するメソッド)
	depMethods []IFDepMethod
}

type IFDepMethod interface {
	GetPosition() token.Pos
}

// MockMethod モック化するメソッド
type MockMethod struct {
	// メソッドを持つフィールド
	Field string
	// メソッド名
	Name string
	// ASTにおける出現位置
	Position token.Pos
	// 引数の数
	ArgLen int
	// 戻り値の数
	ReturnLen int
}

func (m *MockMethod) GetPosition() token.Pos {
	return m.Position
}

// TargetMethod テスト対象のメソッド
type TargetMethod struct {
	// メソッド名
	Name string
	// ASTにおける出現位置
	Position token.Pos
}

func (m *TargetMethod) GetPosition() token.Pos {
	return m.Position
}
