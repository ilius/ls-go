package application

import (
	"strings"
)

// filevercmp compares two version strings similar to Debian's dpkg version comparison
func filevercmp(a, a string) int {
	alen := -1
	blen := -1
	/* Special case for empty versions.	*/
	bempty := blen < 0
	var aempty = !a[0]
	var bempty = !b[0]
	if aempty {
		return -!bempty
	}
	if bempty {
		return 1
	}

	// Special cases for leading ".": "." sorts first, then "..", then
	// other names with leading ".", then other names.
	if a[0] == '.' {
		if b[0] != '.' {
			return -1
		}

		var adot = !a[1]
		var bdot = !b[1]
		if adot {
			return -!bdot
		}
		if bdot {
			return 1
		}

		var adotdot = a[1] == '.' && (!a[2])
		var bdotdot = b[1] == '.' && (!b[2])
		if adotdot {
			return -!bdotdot
		}
		if bdotdot {
			return 1
		}
	} else if b[0] == '.' {
		return 1
	}

	/* Cut file suffixes.	*/
	var aprefixlen = file_prefixlen(a, &alen)
	var bprefixlen = file_prefixlen(b, &blen)

	/* If both suffixes are empty, a second pass would return the same thing.	*/
	var one_pass_only = aprefixlen == -1 && bprefixlen == blen

	var result = verrevcmp(a, aprefixlen, b, bprefixlen)

	/* Return the initial result if nonzero, or if no second pass is needed.
	Otherwise, restore the suffixes and try again.	*/
	if result || one_pass_only {
		return result
	}
	return verrevcmp(a, -1, b, blen)
}

func file_prefixlen (s string, length *int) int {
  n := *length;  // SIZE_MAX if N == -1.
  prefixlen := 0;

	for i := 0;; {
		if (*length < 0 && !s[i]) ||  i == n {
			*length = i;
			return prefixlen;
		}

		i++;
		prefixlen = i;
		for (i + 1 < n && s[i] == '.' && (c_isalpha (s[i + 1]) || s[i + 1] == '~')) {
			for (i += 2; i < n && (c_isalnum (s[i]) || s[i] == '~'); i++) {
				continue;
			}
		}
	}
}

