package subcmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/kazdevl/tgen"
	"github.com/urfave/cli/v2"
)

func generateCreateCommand() *cli.Command {
	return &cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "create test code",
		Action:  createAction,
		Flags:   getCommonFlags(),
	}
}

func createAction(cCtx *cli.Context) error {
	templateDir := cCtx.String(TemplateDirFlag)
	// 引数にはファイル名が入る想定
	for i, targetFilePath := range cCtx.Args().Slice() {
		// テンプレート用のパラメータを格納するjsonファイルの作成
		paramFilePath := fmt.Sprintf("%s/param_%d.json", templateDir, i)
		f, err := os.Create(paramFilePath)
		if err != nil {
			return err
		}
		// ツールの使用目的に大量のファイルに対して一度に実行するケースはほぼないため
		// forループの中でdeferの呼び出しを許容している
		defer func() {
			f.Close()
			os.Remove(paramFilePath)
		}()

		// オプションの用意
		options := createCommonFlagOptionsForGotests(cCtx)

		// jsonファイルにテンプレート用のパラメータを入れる
		var parameters []byte
		parameters, err = tgen.CreateParameterWithFilePath(targetFilePath)
		if err != nil {
			fmt.Printf("tgenの実行時にerrorが発生しました。\n既存のgotestsをそのまま利用します。err=%+v\n", err)
		} else {
			options = append(options,
				"-template_params_file="+paramFilePath,
				"-template_dir="+templateDir,
			)
			if _, err = f.Write(parameters); err != nil {
				return err
			}
		}

		// goのテストコードを自動生成するコマンドの呼び出し
		if err = callGotests(options, targetFilePath); err != nil {
			return err
		}
	}
	return nil
}

func callGotests(options []string, targetFilePath string) error {
	// goのテストコードを自動生成するコマンドの呼び出し
	cmdArgs := append(options, targetFilePath)
	cmd := exec.Command(
		"gotests",
		cmdArgs...,
	)

	// 実行結果の表示
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Printf("Stdout:\n%s\n", stdout.String())

	return nil
}
