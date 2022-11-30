package internal

import (
	"encoding/json"
	"strings"
)

// TemplateParams jsonに渡すパラメータ
type TemplateParams struct {
	// テスト対象メソッドを持つ構造体のフィールドの情報
	FieldMap map[string]*FieldInfo
	// テスト対象メソッドごとのテストケースの情報
	TargetMethodTesCasesMap map[string][]*UpdateTestCase
}

func (t *TemplateParams) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

// UpdateTestCase テンプレートのパラメータ用に更新されたテストケース
type UpdateTestCase struct {
	// テストケースの分岐点となる行数
	Line int
	// 正常系のテストケースか否か
	IsSuccessPattern bool
	// テストケース内で利用されている各フィールドのメソッド群
	DepMethodsInField map[string][]*TemplateMockMethod
}

// TemplateMockMethod テンプレートのパラメータ用のmock化するメソッド
type TemplateMockMethod struct {
	// メソッド名
	Name string
	// ASTにおけるメソッドの位置
	Position int
	// 引数(nil*引数の数)
	Arg string
	// 戻り値(nil*戻り値の数)
	Return string
}

// CreateTemplateParams テスト対象のファイルから抽出した情報(*TestFile)を元にテンプレートのパラメータを返す
func CreateTemplateParams(t *TestFile) *TemplateParams {
	resolvedTargetMethods := make(map[string][]*MockMethod)

	v := new(TemplateParams)
	v.FieldMap = t.FieldMap
	v.TargetMethodTesCasesMap = make(map[string][]*UpdateTestCase, len(t.TargetMethodTesCasesMap))
	for targetMethodName, methodTestCases := range t.TargetMethodTesCasesMap {
		for _, testCase := range methodTestCases {
			uTestCase := new(UpdateTestCase)
			uTestCase.Line = testCase.Line
			uTestCase.IsSuccessPattern = testCase.IsSuccessPattern
			uTestCase.DepMethodsInField = map[string][]*TemplateMockMethod{}
			for _, depMethod := range testCase.depMethods {
				switch method := depMethod.(type) {
				case *TargetMethod:
					mockMethods, ok := resolvedTargetMethods[method.Name]
					if ok {
						inputTemplateMockMethods(mockMethods, uTestCase.DepMethodsInField)
					} else {
						resolvedMockMethods := resolveToMockMethods([]IFDepMethod{method}, t.TargetMethodTesCasesMap, resolvedTargetMethods)
						resolvedTargetMethods[method.Name] = resolvedMockMethods
						inputTemplateMockMethods(resolvedMockMethods, uTestCase.DepMethodsInField)
					}
				case *MockMethod:
					uTestCase.DepMethodsInField[method.Field] = append(
						uTestCase.DepMethodsInField[method.Field],
						&TemplateMockMethod{
							Name:     method.Name,
							Position: int(method.Position),
							Arg:      createNumberOfNilString(method.ArgLen),
							Return:   createNumberOfNilString(method.ReturnLen),
						},
					)
				}
			}
			v.TargetMethodTesCasesMap[targetMethodName] = append(v.TargetMethodTesCasesMap[targetMethodName], uTestCase)
		}
	}
	return v
}

// inputTemplateMockMethods mockメソッド一覧をテンプレートのパラメータに変換して格納する
func inputTemplateMockMethods(src []*MockMethod, dest map[string][]*TemplateMockMethod) {
	for _, mockMethod := range src {
		dest[mockMethod.Field] = append(
			dest[mockMethod.Field],
			&TemplateMockMethod{
				Name:     mockMethod.Name,
				Position: int(mockMethod.Position),
				Arg:      createNumberOfNilString(mockMethod.ArgLen),
				Return:   createNumberOfNilString(mockMethod.ReturnLen),
			},
		)
	}
}

// resolveToMockMethods テスト対象のメソッドから、mock化するメソッド一覧を抽出する
// 引数
// src: テスト対象のメソッド・mock化するメソッドが含まれている
// targetMethodTestCaseMap: 各テスト対象のメソッドにおけるテストケース一覧, テスト対象のメソッドにおける全てのメソッドを抽出するのに利用する
// dest: 抽出済みのテスト対象メソッドにおけるmock化するメソッド一覧, 同じメソッドの抽出を無くす為に利用
func resolveToMockMethods(src []IFDepMethod, targetMethodTestCaseMap map[string][]*TestCase, dest map[string][]*MockMethod) []*MockMethod {
	results := make([]*MockMethod, 0, len(src))
	for _, depMethod := range src {
		switch method := depMethod.(type) {
		case *TargetMethod:
			mockMethods, ok := dest[method.Name]
			if ok {
				results = append(results, mockMethods...)
				continue
			}

			testcases := targetMethodTestCaseMap[method.Name]
			if len(testcases) == 0 {
				continue
			}
			successPattern := testcases[len(testcases)-1]
			targetMethodMockMethods := resolveToMockMethods(successPattern.depMethods, targetMethodTestCaseMap, dest)
			results = append(results, targetMethodMockMethods...)
			dest[method.Name] = targetMethodMockMethods
		case *MockMethod:
			results = append(results, method)
		}
	}
	return results
}

// createNumberOfNilString 指定数のnilの文字列を作成する
// mockメソッドの引数と戻り値の初期値を埋めるのに利用する
func createNumberOfNilString(num int) string {
	if num == 0 {
		return ""
	}
	if num == 1 {
		return "nil"
	}
	nils := make([]string, 0, num)
	for i := 0; i < num; i++ {
		nils = append(nils, "nil")
	}
	return strings.Join(nils, ",")
}
