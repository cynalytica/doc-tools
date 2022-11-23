package toc

import (
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/cynalytica/doc-tools/internal/flags"
	"github.com/cynalytica/doc-tools/internal/utils"
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
	if err := app.Run([]string{"", "-l", "../../resources/test/project1", "-t", "Test Title", "-r", "/the/ttthheee/", "-r", "/Head/Noggin/", "toc"}); err != nil {
		t.Fatal(err)
	}
}
