package tgen

import (
	"encoding/json"
	"errors"
	"go/parser"
	"go/token"

	"github.com/kazdevl/tgen/internal"
	"golang.org/x/tools/go/packages"
)

// CreateParameterWithFilePath ファイルパスを使って、テンプレートのパラメータを作成する
func CreateParameterWithFilePath(src string) ([]byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		return nil, err
	}
	cfg := &packages.Config{
		Mode: packages.NeedImports | packages.NeedTypes,
	}
	pkgs, err := packages.Load(cfg, src)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, errors.New("pkgsの中身は一つを想定しています")
	}
	base, err := internal.GetAnalysisResult(f, fset, pkgs[0].Types)
	if err != nil {
		return nil, err
	}
	return json.Marshal(internal.CreateTemplateParams(base))
}
