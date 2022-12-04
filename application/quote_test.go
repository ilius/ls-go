package application

import (
	"testing"

	"github.com/ilius/is/v2"
	. "github.com/ilius/ls-go/common"
)

func Test_quoteFileName(t *testing.T) {
	app = NewApplication()
	defer func() {
		app = nil
	}()

	is := is.New(t)

	test := func(id int, in string, out string) {
		actualOut := quoteFileName(in)
		is.AddMsg(
			"id=%d, in=%#v, quotingStyle=%v",
			id, in, app.QuotingStyle,
		).Equal(actualOut, out)
	}

	app.QuotingStyle = E_literal
	test(0, "test", "test")
	test(1, "a b", "a b")
	test(2, "a, b", "a, b")
	test(3, "a'b", "a'b")
	test(4, `a"b`, `a"b`)
	test(5, `a"'b`, `a"'b`)
	test(6, "a`b", "a`b")
	test(7, "a\tb", "a"+app.QuestionMark+"b")
	test(8, "a\nb", "a"+app.QuestionMark+"b")
	test(-1, "README.md", "README.md")

	app.QuotingStyle = E_shell
	test(0, "test", "test")
	test(1, "a b", "'a b'")
	test(2, "a, b", "'a, b'")
	test(3, "a'b", `"a'b"`)
	test(4, `a"b`, `'a"b'`)
	test(5, `a"'b`, `'a"'\''b'`)
	test(6, "a`b", "'a`b'") // back-quote acts like space
	test(7, "a\tb", `'a`+app.QuestionMark+`b'`)
	test(8, "a\nb", `'a`+app.QuestionMark+`b'`)
	test(-1, "README.md", "README.md")

	app.QuotingStyle = E_shell_always
	test(0, "test", `'test'`)
	test(1, "a b", `'a b'`)
	test(2, "a, b", `'a, b'`)
	test(3, "a'b", `"a'b"`)
	test(4, `a"b`, `'a"b'`)
	test(5, `a"'b`, `'a"'\''b'`)
	test(6, "a`b", "'a`b'") // back-quote acts like space
	test(7, "a\tb", `'a`+app.QuestionMark+`b'`)
	test(8, "a\nb", `'a`+app.QuestionMark+`b'`)
	test(-1, "README.md", `'README.md'`)

	app.QuotingStyle = E_shell_escape
	test(0, "test", "test")
	test(1, "a b", "'a b'")
	test(2, "a, b", "'a, b'")
	test(3, "a'b", `"a'b"`)
	test(4, `a"b`, `'a"b'`)
	test(5, `a"'b`, `'a"'\''b'`)
	test(6, "a`b", "'a`b'") // back-quote acts like space
	test(7, "a\tb", `'a'$'\t''b'`)
	test(8, "a\nb", `'a'$'\n''b'`)
	test(-1, "README.md", "README.md")

	app.QuotingStyle = E_shell_escape_always
	test(0, "test", `'test'`)
	test(1, "a b", "'a b'")
	test(2, "a, b", "'a, b'")
	test(3, "a'b", `"a'b"`)
	test(4, `a"b`, `'a"b'`)
	test(5, `a"'b`, `'a"'\''b'`)
	test(6, "a`b", "'a`b'") // back-quote acts like space
	test(7, "a\tb", `'a'$'\t''b'`)
	test(8, "a\nb", `'a'$'\n''b'`)
	test(-1, "README.md", `'README.md'`)

	app.QuotingStyle = E_c
	test(0, "test", `"test"`)
	test(1, "a b", `"a b"`)
	test(3, "a'b", `"a'b"`)
	test(4, `a"b`, `"a\"b"`)
	test(5, `a"'b`, `"a\"'b"`)
	test(6, "a`b", "\"a`b\"")
	test(7, "a\tb", `"a\tb"`)
	test(8, "a\nb", `"a\nb"`)
	test(-1, "README.md", `"README.md"`)

	app.QuotingStyle = E_escape
	test(0, "test", "test")
	test(1, "a b", `a\ b`)
	test(2, "a, b", `a,\ b`)
	test(3, "a'b", "a'b")
	test(4, `a"b`, `a"b`)
	test(5, `a"'b`, `a"'b`)
	test(6, "a`b", "a`b")
	test(7, "a\tb", `a\nb`)
	test(8, "a\nb", `a\nb`)
	test(9, `a\b`, `a\\b`)
	test(-1, "README.md", "README.md")
}
