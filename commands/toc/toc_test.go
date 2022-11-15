package toc

import (
	"github.com/cynalytica/doc-tools/internal/flags"
	"github.com/cynalytica/doc-tools/internal/utils"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestRun(t *testing.T) {
	app := cli.NewApp()
	app.Before = func(cCtx *cli.Context) error {
		err := utils.SetUpRegex(cCtx)
		if err != nil {
			return err
		}
		return flags.HandleConfigFile(cCtx)
	}
	app.Flags = append(flags.BuildCommonFlags(), flags.BuildTocFlags()...)
	app.Commands = []*cli.Command{
		{
			Name:    "generate-toc",
			Before:  flags.HandleConfigFile,
			Aliases: []string{"toc"},
			Action:  Run,
		},
	}
	if err := app.Run([]string{"toc", "-l", "../../resources/test/project1", "-t", "Test Title", "-r", "/the/ttthheee/", "toc"}); err != nil {
		t.Fatal(err)
	}
}
