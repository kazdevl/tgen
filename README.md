# tgen
モックのメソッド定義+テストケースを含めたテストコードの自動生成が可能なツールです。

以下の形式のテストコードが自動生成されます。
以下は、testdata/target/target.goのUpdateToRandomNameメソッドのテストコードを自動生成したものになります。

```go
func TestSampleService_UpdateToRandomName(t *testing.T) {
	type fields struct {
		SampleRepository func(ctrl *gomock.Controller) repository.IFSampleRepository
		SampleClient     func(ctrl *gomock.Controller) thirdparty.IFSampleClient
	}
	type args struct {
		i int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "異常: 29行目のif文",
			fields: fields{
				SampleClient: func(ctrl *gomock.Controller) IFSampleClient {
					mock := thirdparty.NewMockIFSampleClient(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GenrateRandomName().Return(nil)
					return mock
				},
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetLastSaveTime(nil).Return(nil)
					return mock
				},
			},
		},
		{
			name: "正常",
			fields: fields{
				SampleClient: func(ctrl *gomock.Controller) IFSampleClient {
					mock := thirdparty.NewMockIFSampleClient(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GenrateRandomName().Return(nil)
					return mock
				},
				SampleRepository: func(ctrl *gomock.Controller) IFSampleRepository {
					mock := repository.NewMockIFSampleRepository(ctrl)
					// TODO embed expected args and return values
					mock.EXPECT().GetLastSaveTime(nil).Return(nil)
					mock.EXPECT().Update(nil, nil).Return(nil)
					return mock
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			s := &SampleService{
				SampleRepository: tt.fields.SampleRepository(ctrl),
				SampleClient:     tt.fields.SampleClient(ctrl),
			}
			assert.True(t, errors.Is(tt.wantErr, s.UpdateToRandomName(tt.args.i)))
		})
	}
}
```

[gotests](https://github.com/cweill/gotests)を内部的に利用しています。

※ このツールを、そのままプロジェクトで活用できることは保証しません！
利用するテンプレートの更新や、自動生成物の微修正が必要になる可能性は十二分にあります。
## Installation
```shell
go install github.com/kazdevl/tgen/cmd/tgen@latest
```

## Preparation
このリポジトリのtemplateディレクトリの.tmpl一覧を利用するプロジェクトに用意してください。
※ 独自で用意していただくことも可能です。

内部的にgotestsを利用していますので、事前にgotestsをインストールして、パスを通してください。

## Usage
```shell
tgen create [オプション] テスト対象のファイル...
```
上記のコマンドで基本的にはテスト対象のファイルに存在する全ての関数やメソッドのテストコードが自動生成されます。
以下に利用できるオプション一覧を示します。
```
--only value          指定した正規表現に合致する関数もしくはメソッドに対してテストを生成する
--exported            公開されている関数もしくはメソッドに対してテストを生成する。onlyよりも優先される (default: false)
--excl value          指定した正規表現に合致しない関数もしくはメソッドに対してテストを生成する。onlyとexportedよりも優先される
--template_dir value  テストの生成に利用するテンプレートのディレクトリへのパス (default: "template")
-i                    エラーメッセージにテストの入力を出力するか (default: true)
--parallel            サブテストを並行実行するテストコードを出力する (default: false)
--help, -h            show help (default: false)
```

- 任意の関数やメソッドを除くテストコードの自動生成
```shell
tgen create -excl="New.*" -template_dir="tmpl" testdata/target/target.go
```
- 任意の関数やメソッドのみのテストコードの自動生成
```shell
tgen create -only="is.*" testdata/target/target.go
```
- 公開されている関数やメソッドのみのテストコードの自動生成
```shell
tgen create -exported testdata/target/target.go
```

- 公開されている関数やメソッドを除くテストコードの自動生成
```shell
tgen create -exported -excl="New.*" testdata/target/target.go
```

## Constraints
1. テスト対象のファイルには、そのテスト対象のメソッドと、それを持つ構造体の定義が一緒に含まれている必要があります
```go
// テスト対象ファイルにSampleの構造体がない場合は制約違反
func (s *Sample) Sample() {
    ....
}
```
2. テスト対象のファイルに、複数の構造体のメソッドが混在してはいけない
```go
type SampleOne struct {}
type SampleTwo struct {}

// テスト対象ファイルに複数の構造体のメソッドが複数定義されているとエラー
func (so *SampleOne) SampleOne() {
    ....
}
func (st *SampleTwo) SampleOne() {
    ....
}
```

上記の制約に反したり、tgenによるテスト対象ファイルの解析時にエラーが発生した場合は通常のgotestsによるテストコードの自動生成に切り替わります。

## About TestCase
tgenで自動生成するテストケースについては、特定の要素を持つif文ごとにテストケースをします。

特定の要素は、以下の内容になります。
1. ifの条件が2項演算式ではない
2. 2項演算式のxの型が*ast.Identではない
3. 2項演算式のyの型が*ast.Identではない
3. 2項演算式の式がerr != nilではない

※上記の特定要素については、現状調査中であり、今後ブラッシュアップするつもです。

## UnSupported
### switch文の対応
現状ではswitch内に存在する全てのモック定義が生成されてしまう。
また、テストケースも分割されない。
```go
// 以下の場合、一つのテストケースに全てのモック定義が含まれることになる。
func (s *Sample) Sample(i int) error {
    switch i {
        case 1:
            return s.Repositorty.SampleOne()
        case 2:
            return s.Repositorty.SampleTwo()
        case 3:
            return s.Repositorty.SampleThree()
    }
    return nil
}
```

### return文に、モック化するメソッドが複数含まれているとき
現状では、テスト対象のメソッドの戻り値の数をreturn文に存在するモック化するメソッドの戻り値の数としている。
```go
// 以下の場合、モック定義について、
// Stringの戻り値が3・IntAndErrorの戻り値も3になってしまう。
func (s *Sample) Sample(i int) (string, int, error) {
    return s.Repositorty.String(), s.Repository.IntAndError()
}
```

