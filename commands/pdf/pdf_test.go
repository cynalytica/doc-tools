package pdf

import (
	"github.com/cynalytica/doc-tools/internal/flags"
	"github.com/cynalytica/doc-tools/internal/utils"
	"github.com/urfave/cli/v2"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	app := cli.NewApp()
	app.Before = func(cCtx *cli.Context) error {
		err := utils.SetUpRegex(cCtx)
		if err != nil {
			return err
		}
		return flags.HandleConfigFile(cCtx)
	}
	app.Flags = append(flags.BuildCommonFlags(), flags.BuildPdfFlags()...)
	app.Commands = []*cli.Command{
		{
			Name:    "generate-pdf",
			Before:  flags.HandleConfigFile,
			Aliases: []string{"pdf"},
			Action:  Create,
		},
	}
	if err := app.Run([]string{"pdf", "-l", "../../resources/test/project1", "-t", "Test Title", "-s", "Really nice subtitle.", "-d", strings.Repeat("This is a really nice abstract. ", 20), "-o", "/home/runner/work/analytics-engine/analytics-engine/ui/src/views/Docs/analytics-engine/media/AnalytICS Engine Documentation.pdf", "-r", "/the/ttthhheee/", "pdf"}); err != nil {
		t.Fatal(err)
	}
}
