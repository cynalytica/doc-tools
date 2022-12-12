package pdf

import (
	"embed"
	"errors"
	"fmt"
	"github.com/minio/cli"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/cynalytica/doc-tools/internal/flags"
	"github.com/cynalytica/doc-tools/internal/utils"
	"github.com/sirupsen/logrus"
)

//go:embed cyrenql.xml cytemplate.tex
var content embed.FS

type file struct {
	Location string
	Index    int `yaml:"nav_order"`
	Content  []byte
}

func Create(ctx *cli.Context) error {
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
	// we create a bunch of temp files, clean them all up here
	tempFiles := make([]string, 0)
	defer closeFiles(tempFiles)

	// process arguments
	title := ctx.String(flags.Title)
	subtitle := ctx.String(flags.Subtitle)
	abstract := ctx.String(flags.Description)
	dir := ctx.String(flags.Location)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf(fmt.Sprintf("can't find directory '%s':", dir), err)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	projectPath := filepath.Join(wd, dir)
	fileDir := filepath.Join(projectPath, "files")
	outputFile := ctx.String(flags.Output)
	if filepath.Ext(outputFile) != ".pdf" {
		outputFile = outputFile + ".pdf"
	}
	outputFile = filepath.Join(wd, outputFile)

	// turn embedded files into temp files that we can pass to pandoc
	syntaxContent, err := content.ReadFile("cyrenql.xml")
	if err != nil {
		return err
	}
	syntaxFile, err := createTemp("cyrenql.xml", syntaxContent)
	if err != nil {
		return err
	}
	tempFiles = append(tempFiles, syntaxFile.Name())
	templateContent, err := content.ReadFile("cytemplate.tex")
	if err != nil {
		return err
	}
	templateFile, err := createTemp("cytemplate.tex", templateContent)
	if err != nil {
		return err
	}
	tempFiles = append(tempFiles, templateFile.Name())

	// process files - find the nav_order for sort and get content
	// unfortunately we have to make temp files for the md files because we have to remove custom markdown
	files := make([]file, 0)
	err = filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || filepath.Ext(path) != ".md" {
			return err
		}
		if !info.IsDir() {
			f := file{
				Location: path,
			}
			fileContent, fErr := os.ReadFile(path)
			if fErr != nil {
				return fErr
			}
			f.Content = fileContent

			_, fErr = frontmatter.Parse(strings.NewReader(string(fileContent)), &f)
			if fErr != nil {
				return fErr
			}
			files = append(files, f)
		}
		return nil
	})
	// hacky way to make a title page
	titleContent := []byte(fmt.Sprintf(`---
title: %s
subtitle: %s
abstract: %s
---
`, title, subtitle, abstract))
	f, err := createTemp("title.md", titleContent)
	files = append(files, file{
		Location: f.Name(),
		Index:    0,
		Content:  titleContent,
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].Index < files[j].Index
	})

	// now that we have the sorted files we can create temp files with cleaned text
	// we'll pass these files as program args in the next step
	orderedFiles := make([]string, 0, len(files))
	for _, f := range files {
		temp, rErr := createTemp(strings.TrimPrefix(strings.ReplaceAll(f.Location, string(filepath.Separator), "-"), projectPath), utils.CleanText(f.Content))
		if rErr != nil {
			return rErr
		}
		tempFiles = append(tempFiles, temp.Name())
		orderedFiles = append(orderedFiles, temp.Name())
	}

	// run the pandoc command
	args := []string{"-s",
		"--toc",
		"--pdf-engine",
		"pdflatex",
		"--from",
		"markdown+escaped_line_breaks+backtick_code_blocks+pipe_tables+multiline_tables+fenced_code_attributes",
		"--title", fmt.Sprintf("\"%s\"", title),
		"--template", templateFile.Name(),
		"--syntax-definition", syntaxFile.Name(),
		"-o", fmt.Sprintf("\"%s\"", outputFile),
		"--metadata", fmt.Sprintf("\"title=%s\"", title),
		"--metadata", fmt.Sprintf("\"subtitle=%s\"", subtitle),
		"--metadata", fmt.Sprintf("\"abstract=%s\"", abstract)}
	args = append(args, orderedFiles...)
	cmd := exec.Command("pandoc", args...)
	str := cmd.String()
	logrus.Info("Running command: `", str, "`")
	out, err := cmd.Output()

	// deal with the fallout (or success)
	pandocMsg := string(out)
	var errMsg string
	if err == nil {
		if len(pandocMsg) > 0 {
			logrus.Info("pandoc message:", pandocMsg)
		}
		logrus.Info("PDF built successfully.")
		return nil
	} else if exitError, ok := err.(*exec.ExitError); ok {
		errMsg = fmt.Sprintf("pandoc exit error: %s", string(exitError.Stderr))
	} else {
		errMsg = err.Error()
	}
	if len(out) > 0 {
		errMsg = errMsg + " - " + string(out)
	}
	return fmt.Errorf("error building pdf: %s", errMsg)
}

func createTemp(name string, content []byte) (*os.File, error) {
	f, err := os.CreateTemp("", name)
	if err != nil {
		return nil, err
	}
	if _, err = f.Write(content); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	return f, nil
}

func closeFiles(files []string) {
	for _, f := range files {
		err := os.Remove(f)
		if err != nil {
			// no need to panic, just warn
			logrus.Warnf("failed to remove file: %s", f)
		}
	}
}
