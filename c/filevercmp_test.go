package main

import (
	"testing"
)

/* set of well sorted examples */
var examples = []string{
	"",
	".",
	"..",
	".0",
	".9",
	".A",
	".Z",
	".a~",
	".a",
	".b~",
	".b",
	".z",
	".zz~",
	".zz",
	".zz.~1~",
	".zz.0",
	".\x01",
	".\x01.txt",
	".\x01x",
	".\x01x\x01",
	".\x01.0",
	"0",
	"9",
	"A",
	"Z",
	"a~",
	"a",
	"a.b~",
	"a.b",
	"a.bc~",
	"a.bc",
	"a+",
	"a.",
	"a..a",
	"a.+",
	"b~",
	"b",
	"gcc-c++-10.fc9.tar.gz",
	"gcc-c++-10.fc9.tar.gz.~1~",
	"gcc-c++-10.fc9.tar.gz.~2~",
	"gcc-c++-10.8.12-0.7rc2.fc9.tar.bz2",
	"gcc-c++-10.8.12-0.7rc2.fc9.tar.bz2.~1~",
	"glibc-2-0.1.beta1.fc10.rpm",
	"glibc-common-5-0.2.beta2.fc9.ebuild",
	"glibc-common-5-0.2b.deb",
	"glibc-common-11b.ebuild",
	"glibc-common-11-0.6rc2.ebuild",
	"libstdc++-0.5.8.11-0.7rc2.fc10.tar.gz",
	"libstdc++-4a.fc8.tar.gz",
	"libstdc++-4.10.4.20040204svn.rpm",
	"libstdc++-devel-3.fc8.ebuild",
	"libstdc++-devel-3a.fc9.tar.gz",
	"libstdc++-devel-8.fc8.deb",
	"libstdc++-devel-8.6.2-0.4b.fc8",
	"nss_ldap-1-0.2b.fc9.tar.bz2",
	"nss_ldap-1-0.6rc2.fc8.tar.gz",
	"nss_ldap-1.0-0.1a.tar.gz",
	"nss_ldap-10beta1.fc8.tar.gz",
	"nss_ldap-10.11.8.6.20040204cvs.fc10.ebuild",
	"z",
	"zz~",
	"zz",
	"zz.~1~",
	"zz.0",
	"zz.0.txt",
	"\x01",
	"\x01.txt",
	"\x01x",
	"\x01x\x01",
	"\x01.0",
	"#\x01.b#",
	"#.b#",
}

/*
Sets of examples that should all sort equally.  Each set is

	terminated by nil.
*/
var equals = [][]string{
	{
		"a",
		"a0",
		"a0000",
	}, {
		"a\x01c-27.txt",
		"a\x01c-027.txt",
		"a\x01c-00000000000000000000000000000000000000000000000000000027.txt",
	}, {
		".a\x01c-27.txt",
		".a\x01c-027.txt",
		".a\x01c-00000000000000000000000000000000000000000000000000000027.txt",
	}, {
		"a\x01c-",
		"a\x01c-0",
		"a\x01c-00",
	}, {
		".a\x01c-",
		".a\x01c-0",
		".a\x01c-00",
	}, {
		"a\x01c-0.txt",
		"a\x01c-00.txt",
	}, {
		".a\x01c-1\x01.txt",
		".a\x01c-001\x01.txt",
	},
}

func sign(n int) int {
	if n < 0 {
		return -1
	}
	if n > 0 {
		return 1
	}
	return 0
}

func signName(n int) string {
	if n < 0 {
		return "negative"
	}
	if n > 0 {
		return "positive"
	}
	return "zero"
}

func test_filevercmp_sorted_list(t *testing.T, list []string) {
	n := len(list)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cmp := filevercmp(list[i], list[j])
			if sign(cmp) != sign(i-j) {
				t.Logf(
					"cmp=%#v, expected %v, list[%d]=%#v, list[%d]=%#v",
					cmp, signName(i-j),
					i, list[i],
					j, list[j],
				)
			}
		}
	}
}

func Test_filevercmp(t *testing.T) {
	test_filevercmp_sorted_list(t, examples)
	for _, list := range equals {
		n := len(list)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				cmp := filevercmp(list[i], list[j])
				if cmp != 0 {
					t.Logf("cmp=%#v, expected 0, v1=%#v, v2=%#v", cmp, list[i], list[j])
				}
			}
		}
	}
}
