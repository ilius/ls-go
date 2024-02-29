package lsargs

import (
	"fmt"
	"regexp"
	"strings"

	goopt "github.com/ilius/goopt"
)

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func helpMarkdown() {
	flagRE := regexp.MustCompile(`( -[a-z\-=]+)`)
	flagReplace := func(m string) string {
		return " `" + m[1:] + "`"
	}
	optHelpFormat := func(help string) string {
		help = capitalize(help)
		help = flagRE.ReplaceAllStringFunc(help, flagReplace)
		help = strings.ReplaceAll(help, "ls `-", "`ls -")
		help = strings.ReplaceAll(help, "` `", " ")
		help = strings.ReplaceAll(help, "'", "`")
		help = strings.ReplaceAll(help, "``", "`")
		help = strings.ReplaceAll(help, ";\n", ".\\\n")
		help = strings.ReplaceAll(help, "; ", ".\\\n")
		if !strings.HasSuffix(help, ".") {
			help += "."
		}
		return help
	}
	helpMd := goopt.HelpMarkdown(optHelpFormat)
	helpMd = "#" + strings.ReplaceAll(helpMd, "\n#", "\n##")
	helpMd = strings.ReplaceAll(helpMd, ".`.\n", ".`\n")
	/*
		lines := strings.Split(helpMd, "\n")
		for i, line := range lines {
			if line == "" {
				continue
			}
			if strings.HasPrefix(line, "#") {
				continue
			}
			lines[i] = "\t " + line
		}
		helpMd = strings.Join(lines, "\n")
	*/
	fmt.Println(helpMd)
}
