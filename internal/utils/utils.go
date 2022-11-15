package utils

import (
	"github.com/cynalytica/doc-tools/internal/flags"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"regexp"
)

var (
	IdRegex = regexp.MustCompile(" {#([\\w+\\-]*)}")
)

type RegexRmRp struct {
	RmRegex *regexp.Regexp
	Replace string
}

var Regex []RegexRmRp

func SetUpRegex(cCtx *cli.Context) error {
	rm := regexp.MustCompile("/(.*)/")
	rmRp := regexp.MustCompile("/(.*)/(.*)/")
	regexArgs := cCtx.StringSlice(flags.Regex)
	for _, arg := range regexArgs {
		var rmStr string
		var rpStr string
		if match := rmRp.FindStringSubmatch(arg); match != nil {
			rmStr = match[1]
			rpStr = match[2]
		} else if match = rm.FindStringSubmatch(arg); match != nil {
			rmStr = match[1]
			rpStr = ""
		} else {
			logrus.Warnf("couldn't find regex for `%s`, skipping...", arg)
			continue
		}
		re, err := regexp.Compile(rmStr)
		if err != nil {
			logrus.Warnf("`%s` is not valid regex, skipping...", rmStr)
			continue
		}
		Regex = append(Regex, RegexRmRp{
			RmRegex: re,
			Replace: rpStr,
		})
	}
	return nil
}

func CleanText(text []byte) []byte {
	for _, re := range Regex {
		text = re.RmRegex.ReplaceAll(text, []byte(re.Replace))
	}
	return text
}
