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
	if err := app.Run([]string{"pdf", "-l", "../../resources/test/serialguards", "-t", "Test Title", "-s", "Really nice subtitle.", "-d", strings.Repeat("This is a really nice abstract. ", 20), "--version", "1.2.3", "--commit", "abc123def456ggg789", "-o", "../../resources/test/serialguards/media/SerialGuard Documentation.pdf", "-r", "/the/ttthhheee/", "pdf"}); err != nil {
		t.Fatal(err)
	}
}
