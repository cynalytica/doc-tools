package main

import (
	"context"
	"fmt"
	"github.com/cynalytica/doc-tools/commands/pdf"
	"github.com/cynalytica/doc-tools/commands/toc"
	"github.com/cynalytica/doc-tools/internal/flags"
	"github.com/cynalytica/doc-tools/internal/meta"
	"github.com/cynalytica/doc-tools/internal/utils"
	"github.com/urfave/cli/v2"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"os"
)

var cancelFunc context.CancelFunc

func sigHandler(ctx context.Context, sigHandle <-chan os.Signal) {
	for {
		select {
		case <-sigHandle:
			cancelFunc()
			return
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	var ctx context.Context
	sighandle := make(chan os.Signal, 1)
	signal.Notify(sighandle, syscall.SIGINT, syscall.SIGTERM)
	cnx := context.Background()
	ctx, cancelFunc = context.WithCancel(cnx)
	app := cli.NewApp()
	app.Name = meta.Name
	app.Version = meta.Version
	app.Authors = []*cli.Author{{Name: meta.Vendor, Email: meta.VendorMail}}
	app.Usage = meta.Usage
	app.Copyright = meta.Vendor
	app.Before = func(cCtx *cli.Context) error {
		err := utils.SetUpRegex(cCtx)
		if err != nil {
			return err
		}
		return flags.HandleConfigFile(cCtx)
	}
	app.Flags = flags.BuildGlobalFlags()
	app.UseShortOptionHandling = true
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "generate-toc",
			Before:  flags.HandleConfigFile,
			Aliases: []string{"toc"},
			Action:  toc.Run,
			Flags:   append(flags.BuildCommonFlags(), flags.BuildTocFlags()...),
		},
		{
			Name:    "generate-pdf",
			Before:  flags.HandleConfigFile,
			Aliases: []string{"pdf"},
			Action:  pdf.Create,
			Flags:   append(flags.BuildCommonFlags(), flags.BuildPdfFlags()...),
		},
		{
			Name:    "version",
			Before:  flags.HandleConfigFile,
			Aliases: []string{"v"},
			Action: func(cCtx *cli.Context) error {
				fmt.Printf("%s v%s (%s)", meta.Name, meta.Version, meta.CommitHash)
				return nil
			},
		},
	}
	go sigHandler(ctx, sighandle)
	if err := app.RunContext(ctx, os.Args); err != nil {
		logrus.Fatal(err)
	}
}
