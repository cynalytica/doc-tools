package flags

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	options = map[string]func(flagFileName string) func(context *cli.Context) (altsrc.InputSourceContext, error){
		".yml":  altsrc.NewYamlSourceFromFlagFunc,
		".yaml": altsrc.NewYamlSourceFromFlagFunc,
		".toml": altsrc.NewTomlSourceFromFlagFunc,
		".ini":  altsrc.NewTomlSourceFromFlagFunc,
		".json": altsrc.NewJSONSourceFromFlagFunc,
	}
)

func HandleConfigFile(ctx *cli.Context) (err error) {
	configFileLocation := ctx.String(ConfigFile)
	if configFileLocation != "" {
		fileExtension := strings.ToLower(filepath.Ext(configFileLocation))
		source, ok := options[fileExtension]
		if !ok {
			return fmt.Errorf("unkown file extension, cannont use as config file")
		}
		flags := ctx.Command.Flags
		if len(flags) < 1 {
			flags = ctx.App.Flags
		}
		err = altsrc.InitInputSourceWithContext(flags, source(ConfigFile))(ctx)
	}
	return err

}

func BuildGlobalFlags() []cli.Flag {
	flags := make([]cli.Flag, 0)
	flags = append(flags, buildHelperFlags()...)
	return flags
}

func buildHelperFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    ConfigFile,
			Aliases: []string{"c"},
			Usage:   "File to read in configuration from rather than command line or ENV vars",
			Hidden:  true,
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    ConfigDebugLevelFlag,
			EnvVars: []string{ConfigDebugLevelEnv},
			Hidden:  true,
		}),
	}
}

func BuildCommonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        Location,
			DefaultText: ".",
			Usage:       "The doc directory location.",
			Aliases:     []string{"l"},
			Required:    true,
		},
		&cli.StringFlag{
			Name:        Title,
			DefaultText: "Documentation",
			Usage:       "The doc title. Will be used in the UI and/or generated PDF.",
			Required:    true,
			Aliases:     []string{"t"},
		},
		&cli.StringFlag{
			Name:    Subtitle,
			Usage:   "The doc subtitle. Will be used in the UI and/or generated PDF.",
			Aliases: []string{"s"},
		},
		&cli.StringSliceFlag{
			Name:    Regex,
			Usage:   "Removes (`/text-to-remove/`) or replaces (`/remove/replace/`) text with regex.",
			Aliases: []string{"r"},
		},
		&cli.GenericFlag{
			Name:  RegexFile,
			Usage: "Location of a file containing line seperated regexs",
			Value: &StringArrayFile{},
		},
		&cli.StringFlag{
			Name:    Description,
			Usage:   "The overall description of the docs. Will appear in the doc's TOC.",
			Aliases: []string{"d"},
		},
	}
}

func BuildTocFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    IncludeContent,
			Usage:   "Whether to include the entire contents of the markdown files in the `manifest.json` file. Not recommended.",
			Aliases: []string{"i"},
		},
	}
}

func BuildPdfFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        Output,
			DefaultText: "./docs.pdf",
			Usage:       "The output directory and file.",
			Aliases:     []string{"o"},
			Required:    true,
		},
		&cli.StringFlag{
			Name:     Version,
			Required: false,
		},
		&cli.StringFlag{
			Name:     CommitHash,
			Required: false,
		},
	}
}
