package toc

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"

	tocLib "github.com/abhinav/goldmark-toc"
	"github.com/adrg/frontmatter"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"

	"github.com/cynalytica/doc-tools/internal/flags"
	"github.com/cynalytica/doc-tools/internal/utils"
)

type manifest struct {
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	Files       []*file  `json:"files"`
	Media       []string `json:"media"`
}

type file struct {
	Title    string `json:"title" yaml:"title"`
	Location string `json:"-"`
	Endpoint string `json:"endpoint"`
	NavOrder int    `json:"order" yaml:"nav_order"`
	Content  string `json:"content,omitempty"`
	TOC      []toc  `json:"toc"`
}

type toc struct {
	Title    string `json:"title"`
	ID       string `json:"id,omitempty"`
	Children []toc  `json:"children,omitempty"`
}

func Run(ctx *cli.Context) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("recover in generate.Run: %s\n%s", r, string(debug.Stack()))
			err = errors.New(msg)
			logrus.Errorf(msg)
		}
	}()
	err = utils.SetUpRegex(ctx)
	if err != nil {
		return err
	}
	// get and format args
	dir := ctx.String(flags.Location)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf(fmt.Sprintf("can't find directory '%s':", dir), err)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	projectPath := filepath.Join(wd, dir)
	title := ctx.String(flags.Title)
	subtitle := ctx.String(flags.Subtitle)
	description := ctx.String(flags.Description)

	// set up manifest file
	manifestFile := manifest{
		Title:       title,
		Subtitle:    subtitle,
		Description: description,
		Files:       make([]*file, 0),
		Media:       make([]string, 0),
	}

	// add media files and md to the manifest
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == "config.yaml" || path == "manifest.json" {
			return err
		}
		if !info.IsDir() {
			path = filepath.ToSlash(path)
			if filepath.Ext(path) == ".md" {
				f := &file{
					TOC:      make([]toc, 0),
					Location: path,
				}
				manifestFile.Files = append(manifestFile.Files, f)
				return nil
			}
			mediaDir, _ := filepath.Split(path)
			if filepath.Base(mediaDir) == "media" {
				manifestFile.Media = append(manifestFile.Media, strings.TrimPrefix(path, filepath.ToSlash(projectPath)))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// parse md files
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	for _, f := range manifestFile.Files {
		fileContent, fErr := os.ReadFile(f.Location)
		if fErr != nil {
			return fErr
		}
		content, fErr := frontmatter.Parse(strings.NewReader(string(fileContent)), &f)
		if fErr != nil {
			return fErr
		}
		if ctx.Bool(flags.IncludeContent) {
			f.Content = string(content)
		}
		f.Endpoint = filepath.ToSlash(filepath.Join(strings.TrimSuffix(filepath.Base(f.Location), filepath.Ext(f.Location))))
		clean := utils.CleanText(content)
		doc := md.Parser().Parse(text.NewReader(clean))
		var seen map[string]struct{} // keeps track of seen slugs to avoid duplicate ids
		// walkFn adapted from https://gist.github.com/artyom/26c2674d459669a38eb8b84f95fa30fb
		// checks if heading node has an associated id, if it doesn't then adds one
		walkFn := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering || n.Kind() != ast.KindHeading {
				return ast.WalkContinue, nil
			}
			if name := slugify(nodeText(n, clean)); name != "" {
				if seen == nil {
					seen = make(map[string]struct{})
				}
				for i := 0; i < 100; i++ {
					var cand string
					if i == 0 {
						cand = name
					} else {
						cand = fmt.Sprintf("%s-%d", name, i)
					}
					if _, ok := seen[cand]; !ok {
						seen[cand] = struct{}{}
						n.SetAttributeString("id", []byte(cand))
						break
					}
				}
			}
			return ast.WalkContinue, nil
		}
		if err := ast.Walk(doc, walkFn); err != nil {
			return err
		}
		// get toc structure and convert it into a toc struct
		tree, fErr := tocLib.Inspect(doc, clean)
		if fErr != nil {
			return fErr
		}
		f.TOC = walkTree(tree)
	}
	sort.Slice(manifestFile.Files, func(i, j int) bool {
		return manifestFile.Files[i].NavOrder < manifestFile.Files[j].NavOrder
	})

	res, err := json.Marshal(manifestFile)
	if err != nil {
		return err
	}

	manifestLocation, err := filepath.Abs(projectPath)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(manifestLocation, "manifest.json"), res, 0644)
	if err != nil {
		return err
	}

	logrus.Info("Successfully created CyRenQL Docs and manifest file")
	return nil
}
