package internal

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
)

// GetAnalysisResult ASTから値を抽出し、テンプレートのパラメータ用の構造体を生成する
func GetAnalysisResult(astF *ast.File, fset *token.FileSet, packageTypes *types.Package) (*TestFile, error) {
	v := new(TestFile)
	targetStructName, abbreviationName, err := extractTargetStructName(astF)
	if err != nil {
		return nil, err
	}
	fm, err := extractTargetStructInfo(packageTypes, targetStructName)
	if err != nil {
		return nil, err
	}
	v.FieldMap = fm
	inspect := inspector.New([]*ast.File{astF})
	v.TargetMethodTesCasesMap = extractTargetMethodTestCasesMap(fset, inspect, abbreviationName, targetStructName, fm)
	return v, nil
}

// extractTargetStructName テスト対象のメソッドを持つ構造体の名前などを抽出する
// 戻り値
// structName: 構造体名
// abbreviationName: メソッドのレシーバー名
func extractTargetStructName(src *ast.File) (structName, abbreviationName string, err error) {
	recvTypeNameMap := make(map[string]string)
	for _, decl := range src.Decls {
		fDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fDecl.Recv == nil {
			continue
		}
		abbreviation := fDecl.Recv.List[0].Names[0].Name
		var recvTypeName string
		recvTypeName, err = extractRecvTypeName(fDecl.Recv.List[0].Type)
		if err != nil {
			return "", "", err
		}
		recvTypeNameMap[recvTypeName] = abbreviation
	}
	if len(recvTypeNameMap) != 1 {
		return "", "", errors.New("テスト対象のメソッドを持つ構造体が一意に定まりません")
	}
	for k, v := range recvTypeNameMap {
		return k, v, nil
	}
	return "", "", nil
}

// extractTargetStructInfo テスト対象のメソッドを持つ構造体の情報(フィールド)を抽出する
func extractTargetStructInfo(packageTypes *types.Package, targetStructName string) (fieldMap map[string]*FieldInfo, err error) {
	if packageTypes.Scope() == nil {
		err = errors.New("構造体の型情報を読み取れていません")
		return
	}
	structObj := packageTypes.Scope().Lookup(targetStructName)
	if structObj == nil {
		err = errors.New("対象の構造体の情報を読み取れていません")
		return
	}
	structUnderLyingType, ok := structObj.Type().Underlying().(*types.Struct)
	if !ok {
		err = errors.New("読み取った構造体の型は構造体ではありません")
		return
	}
	fieldMap = make(map[string]*FieldInfo, structUnderLyingType.NumFields())
	for i := 0; i < structUnderLyingType.NumFields(); i++ {
		field := structUnderLyingType.Field(i)
		typeName := field.Type().String()
		typeNameIndex := strings.LastIndex(typeName, "/") + 1
		if strings.Contains(typeName, "command-line-arguments.") {
			typeNameIndex += len("command-line-arguments.")
		}
		typeName = typeName[typeNameIndex:]
		packageName := ""
		if packageNameIndex := strings.Index(typeName, "."); packageNameIndex != -1 {
			packageName = typeName[:packageNameIndex]
			typeName = typeName[packageNameIndex+1:]
		}

		fieldMap[field.Name()] = &FieldInfo{
			IsInterface:            strings.Contains(field.Type().Underlying().String(), "interface{"),
			PackageName:            packageName,
			TypeName:               typeName,
			UpperCamelCaseTypeName: strings.ToUpper(typeName[0:1]) + typeName[1:],
		}
	}
	return
}

// extractTargetMethodTestCasesMap 各テスト対象のメソッドにおけるテストケース一覧を抽出する
func extractTargetMethodTestCasesMap(fset *token.FileSet, inspect *inspector.Inspector, abbreviationName, targetStructName string, fieldMap map[string]*FieldInfo) map[string][]*TestCase {
	targetMethodTestCaseMap := map[string][]*TestCase{}
	targetMethodIfPositionsMap := make(map[string][]token.Pos, 0)
	targetMethodDepMethodsMap := make(map[string][]IFDepMethod, 0)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.ReturnStmt)(nil),
		(*ast.IfStmt)(nil),
	}

	methodName := ""
	methodReturnNum := 0
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		switch n := node.(type) {
		case *ast.FuncDecl:
			var isSuccess bool
			methodName, methodReturnNum, isSuccess = extractValuesFromFuncDecl(n, targetStructName)
			if !isSuccess {
				return
			}
		case *ast.AssignStmt:
			callExpr, ok := n.Rhs[0].(*ast.CallExpr)
			if !ok {
				return
			}
			depMethod, isSuccess := extractDepMethodFromCallExpr(callExpr, abbreviationName, methodReturnNum, fieldMap)
			if isSuccess {
				targetMethodDepMethodsMap[methodName] = append(targetMethodDepMethodsMap[methodName], depMethod)
			}
		case *ast.ReturnStmt:
			for _, result := range n.Results {
				callExpr, ok := result.(*ast.CallExpr)
				if !ok {
					return
				}
				depMethod, isSuccess := extractDepMethodFromCallExpr(callExpr, abbreviationName, methodReturnNum, fieldMap)
				if isSuccess {
					targetMethodDepMethodsMap[methodName] = append(targetMethodDepMethodsMap[methodName], depMethod)
				}
			}
		case *ast.IfStmt:
			callExpr, ok := n.Cond.(*ast.CallExpr)
			if ok {
				depMethod, isSuccess := extractDepMethodFromCallExpr(callExpr, abbreviationName, methodReturnNum, fieldMap)
				if isSuccess {
					targetMethodDepMethodsMap[methodName] = append(targetMethodDepMethodsMap[methodName], depMethod)
				}
			}
			pos, isSuccess := extractPositionFromIfStmt(n)
			if !isSuccess {
				return
			}
			targetMethodIfPositionsMap[methodName] = append(targetMethodIfPositionsMap[methodName], pos)
		}
	})

	for k, v := range targetMethodDepMethodsMap {
		targetMethodTestCaseMap[k] = getTestCases(fset, targetMethodIfPositionsMap[k], v)
	}

	return targetMethodTestCaseMap
}

func extractRecvTypeName(src ast.Expr) (string, error) {
	var name string
	switch recvType := src.(type) {
	case *ast.StarExpr:
		x, ok := recvType.X.(*ast.Ident)
		if !ok {
			return "", errors.New("想定していないデータ形式です")
		}
		name = x.Name
	case *ast.Ident:
		name = recvType.Name
	}
	return name, nil
}

func extractValuesFromFuncDecl(src *ast.FuncDecl, targetStructName string) (methodName string, methodReturnNum int, isSuccess bool) {
	if src.Recv == nil {
		return
	}
	recvTypeName, err := extractRecvTypeName(src.Recv.List[0].Type)
	if err != nil {
		return
	}
	if recvTypeName != targetStructName {
		return
	}
	// 値の抽出
	if src.Type.Results != nil {
		methodReturnNum = len(src.Type.Results.List)
	}
	methodName = src.Name.Name
	isSuccess = true
	return
}

func extractPositionFromIfStmt(src *ast.IfStmt) (pos token.Pos, isSuccess bool) {
	var ok bool
	var binaryExpr *ast.BinaryExpr
	binaryExpr, ok = src.Cond.(*ast.BinaryExpr)
	if !ok {
		return src.Pos(), true
	}

	bodyLen := len(src.Body.List)
	if bodyLen == 1 {
		_, ok = src.Body.List[0].(*ast.ReturnStmt)
		if !ok {
			return 0, false
		}
	}

	var x *ast.Ident
	x, ok = binaryExpr.X.(*ast.Ident)
	if !ok {
		return src.Pos(), true
	}
	var y *ast.Ident
	y, ok = binaryExpr.Y.(*ast.Ident)
	if !ok {
		return src.Pos(), true
	}

	if strings.Contains(strings.ToLower(x.String()), "err") &&
		binaryExpr.Op.String() == "!=" &&
		y.String() == "nil" {
		return 0, false
	}
	return src.Pos(), true
}

func getTestCases(fset *token.FileSet, ifPositions []token.Pos, depMethods []IFDepMethod) []*TestCase {
	if len(ifPositions) == 0 {
		return []*TestCase{{
			Line:             0,
			IsSuccessPattern: true,
			depMethods:       depMethods,
		}}
	}

	testcases := make([]*TestCase, 0, len(ifPositions)+1)
	index := 0
	for i, depMethod := range depMethods {
		if len(ifPositions) == index {
			// 最後のif文の失敗のテストケースの作成後
			return append(testcases, &TestCase{
				Line:             fset.Position(ifPositions[index-1]).Line,
				IsSuccessPattern: true,
				depMethods:       depMethods,
			})
		}
		ifLine := fset.Position(ifPositions[index]).Line
		depMethodLine := fset.Position(depMethod.GetPosition()).Line
		if ifLine < depMethodLine {
			testcases = append(testcases, &TestCase{
				Line:       fset.Position(ifPositions[index]).Line,
				depMethods: depMethods[:i],
			})
			index++
		}
	}

	// if文が全ての依存しているメソッドよりもまだ後の行にある場合
	for ; index <= len(ifPositions); index++ {
		if len(ifPositions) == index {
			return append(testcases, &TestCase{
				Line:             fset.Position(ifPositions[index-1]).Line,
				IsSuccessPattern: true,
				depMethods:       depMethods,
			})
		}
		testcases = append(testcases, &TestCase{
			Line:       fset.Position(ifPositions[index]).Line,
			depMethods: depMethods,
		})
	}

	return testcases
}

func extractDepMethodFromCallExpr(src *ast.CallExpr, targetAbbreviationName string, methodReturnLen int, fieldMap map[string]*FieldInfo) (IFDepMethod, bool) {
	targetMethod, isSuccess := extractTargetMethodFromCallExpr(src, targetAbbreviationName)
	if isSuccess {
		return targetMethod, true
	}

	mockMethod, isSuccess := extractMockMethodFromCallExpr(src, targetAbbreviationName, fieldMap)
	if isSuccess {
		mockMethod.ReturnLen = methodReturnLen
		return mockMethod, true
	}
	return nil, false
}

func extractTargetMethodFromCallExpr(src *ast.CallExpr, targetAbbreviationName string) (*TargetMethod, bool) {
	selectorExpr, ok := src.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}
	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return nil, false
	}
	if targetAbbreviationName != ident.Name {
		return nil, false
	}
	return &TargetMethod{
		Name:     selectorExpr.Sel.Name,
		Position: src.Pos(),
	}, true
}

func extractMockMethodFromCallExpr(src *ast.CallExpr, targetAbbreviationName string, fieldMap map[string]*FieldInfo) (*MockMethod, bool) {
	selectorExpr, ok := src.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}
	selectorExprForField, ok := selectorExpr.X.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}
	ident, ok := selectorExprForField.X.(*ast.Ident)
	if !ok {
		return nil, false
	}
	if targetAbbreviationName != ident.Name {
		return nil, false
	}
	fieldInfo, ok := fieldMap[selectorExprForField.Sel.Name]
	if !ok {
		return nil, false
	}
	if !fieldInfo.IsInterface {
		return nil, false
	}
	return &MockMethod{
		Field:    selectorExprForField.Sel.Name,
		Name:     selectorExpr.Sel.Name,
		Position: src.Pos(),
		ArgLen:   len(src.Args),
	}, true
}
