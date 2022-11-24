package pdf

import (
	"fmt"
	"github.com/cynalytica/doc-tools/internal/flags"
	"github.com/cynalytica/doc-tools/internal/utils"
	"github.com/urfave/cli/v2"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	fmt.Println("DELETE ME")
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
	if err := app.Run([]string{"pdf", "-l", "../../resources/test/project1", "-t", "Test Title", "-s", "Really nice subtitle.", "-d", strings.Repeat("And this is a really nice abstract. ", 20), "-o", "../../resources/test/project1/Cool Documentation.pdf", "-r", "/the/ttthhheee/", "pdf"}); err != nil {
		t.Fatal(err)
	}
}
