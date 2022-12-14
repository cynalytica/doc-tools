package toc

import (
	tocLib "github.com/abhinav/goldmark-toc"
	"github.com/cynalytica/doc-tools/internal/utils"
	"github.com/yuin/goldmark/ast"
	"regexp"
	"strings"
	"unicode"
)

var re = regexp.MustCompile("!\\{link(.*)}")

// nodeText walks node and extracts plain text from it and its descendants,
// effectively removing all markdown syntax
func nodeText(node ast.Node, src []byte) string {
	var b strings.Builder
	fn := func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindText:
			if t, ok := n.(*ast.Text); ok {
				text := t.Text(src)
				b.Write(text)
			}
		}
		return ast.WalkContinue, nil
	}
	if err := ast.Walk(node, fn); err != nil {
		return ""
	}
	return re.ReplaceAllString(b.String(), "$1")
}

func slugify(text string) string {
	if match := utils.IdRegex.FindStringSubmatch(text); match != nil && len(match) > 1 && len(match[1]) > 0 {
		return match[1]
	}
	// don't want to have dupe IDs, so prepend "docs-" for autogenerated IDs
	// dupe prevention for autogenerated IDs is handled down the pipeline
	text = "docs-" + text
	f := func(r rune) rune {
		switch {
		case r == '-' || r == '_':
			return r
		case unicode.IsSpace(r):
			return '-'
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			return unicode.ToLower(r)
		}
		return -1
	}
	return strings.Map(f, text)
}

func handleLayer(item *tocLib.Item) toc {
	ref := toc{
		Children: make([]toc, 0),
	}
	item.Title = utils.IdRegex.ReplaceAll(item.Title, []byte{})
	ref.Title = re.ReplaceAllString(string(item.Title), "$1")
	ref.ID = string(item.ID)
	ref.Children = make([]toc, 0)
	for _, child := range item.Items {
		if child.Title == nil {
			continue
		}
		ref.Children = append(ref.Children, handleLayer(child))
	}
	return ref
}

func walkTree(tree *tocLib.TOC) []toc {
	tocArr := make([]toc, 0)
	for _, item := range tree.Items {
		if string(item.Title) == "Table of Contents" && len(item.Items) == 0 {
			continue
		}
		tocArr = append(tocArr, handleLayer(item))
	}
	return tocArr
}
