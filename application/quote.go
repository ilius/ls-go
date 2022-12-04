package application

import (
	"fmt"
	"strings"

	. "github.com/ilius/ls-go/common"
)

func quoteFileName(in string) string {
	switch app.QuotingStyle {
	case E_none:
		return in
	case E_literal:
		return quoteLiteral(in)
	case E_shell:
		return quoteShellGeneric(in, false, false)
	case E_shell_always:
		return quoteShellGeneric(in, false, true)
	case E_shell_escape:
		return quoteShellGeneric(in, true, false)
	case E_shell_escape_always:
		return quoteShellGeneric(in, true, true)
	case E_c:
		return quoteC(in)
	case E_escape:
		return quoteEscape(in)
	}
	return in
}

// with "shell" and "literal", newline and tab is replaced with question mark!

func quoteLiteral(in string) string {
	in = strings.Replace(in, "\t", app.QuestionMark, -1)
	in = strings.Replace(in, "\n", app.QuestionMark, -1)
	return in
}

func quoteShellGeneric(in string, escape bool, always bool) string {
	// back-quote acts like space
	// special means space or back-quote

	// `escape` applies to newline, tab and possibly other "non-graphic" characters
	// in ls.c lookup var qmark_funny_chars
	// comments in ls.c says qmark_funny_chars is independent from quoting style
	// but I cannot confirm that from the code!

	special := false
	double_q := 0
	single_q := 0
	tab := 0
	newline := 0
	for _, c := range in {
		switch c {
		case ' ', '`':
			special = true
		case '"':
			double_q++
		case '\'':
			single_q++
		case '\t':
			tab++
		case '\n':
			newline++
		}
	}
	if !(always || special || double_q > 0 || single_q > 0 || tab > 0 || newline > 0) {
		return in
	}
	useDoubleQ := false
	if single_q > 0 {
		if double_q > 0 {
			in = strings.Replace(in, "'", `'\''`, single_q)
		} else {
			useDoubleQ = true
		}
	}
	if escape {
		if tab > 0 {
			in = strings.Replace(in, "\t", `'$'\t''`, tab)
		}
		if newline > 0 {
			in = strings.Replace(in, "\n", `'$'\n''`, newline)
		}
	} else {
		if tab > 0 {
			in = strings.Replace(in, "\t", app.QuestionMark, tab)
		}
		if newline > 0 {
			in = strings.Replace(in, "\n", app.QuestionMark, newline)
		}
	}
	// so now special is true, (double_q > 0 || single_q > 0) is true
	if useDoubleQ {
		return `"` + in + `"`
	}
	return "'" + in + "'"
}

func quoteC(in string) string {
	return fmt.Sprintf("%#v", in)
	// b, err := json.Marshal(in)
	// check(err)
	// return string(b)
}

func quoteEscape(in string) string {
	in = strings.Replace(in, `\`, `\\`, -1)
	in = strings.Replace(in, " ", `\ `, -1)
	in = strings.Replace(in, "\t", `\n`, -1)
	in = strings.Replace(in, "\n", `\n`, -1)
	return in
}
