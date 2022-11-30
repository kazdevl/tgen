package subcmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const (
	OnlyFlag            = "only"
	ExportedFlag        = "exported"
	ExclFlag            = "excl"
	TemplateDirFlag     = "template_dir"
	PrintTestInputsFlag = "i"
	ParallelFlag        = "parallel"
)

func ProvideSubCommands() cli.Commands {
	return cli.Commands{
		generateCreateCommand(),
	}
}

// gotestsコマンドで利用するオプション
func getCommonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name: OnlyFlag, Usage: "指定した正規表現に合致する関数もしくはメソッドに対してテストを生成する。",
		},
		&cli.BoolFlag{
			Name: ExportedFlag, Usage: "公開されている関数もしくはメソッドに対してテストを生成する。onlyよりも優先される", Value: false,
		},
		&cli.StringFlag{
			Name: ExclFlag, Usage: "指定した正規表現に合致しない関数もしくはメソッドに対してテストを生成する。onlyとexportedよりも優先される",
		},
		&cli.StringFlag{
			Name: TemplateDirFlag, Usage: "テストの生成に利用するテンプレートのディレクトリへのパス", Value: "template",
		},
		&cli.BoolFlag{
			Name: PrintTestInputsFlag, Usage: "エラーメッセージにテストの入力を出力するか", Value: false,
		},
		&cli.BoolFlag{
			Name: ParallelFlag, Usage: "サブテストを並行実行するテストコードを出力する", Value: false,
		},
	}
}

func createCommonFlagOptionsForGotests(cCtx *cli.Context) []string {
	options := []string{
		"-w",
		"-all",
		fmt.Sprintf("-exported=%t", cCtx.Bool(ExportedFlag)),
		fmt.Sprintf("-i=%t", cCtx.Bool(PrintTestInputsFlag)),
		fmt.Sprintf("-parallel=%t", cCtx.Bool(ParallelFlag)),
	}
	if onlyFuncs := cCtx.String(OnlyFlag); onlyFuncs != "" {
		options = append(options, fmt.Sprintf("-only=%s", onlyFuncs))
	}
	if exclFuncs := cCtx.String(ExclFlag); exclFuncs != "" {
		options = append(options, fmt.Sprintf("-excl=%s", exclFuncs))
	}
	return options
}
