package application

import (
	"time"

	"github.com/itchyny/timefmt-go"
)

func Strftime(tm *time.Time, fmt string) string {
	if tm == nil {
		return ""
	}
	// "github.com/leekchan/timeutil"
	// return timeutil.Strftime(tm, fmt)
	// has %X, but no %T

	// "github.com/jehiah/go-strftime"	strftime.Format(fmt, *tm)
	// return strftime.Format(fmt, *tm)
	// has no %X, %T, %F, let alone %c, %j etc

	// "github.com/itchyny/timefmt-go"
	return timefmt.Format(*tm, fmt)
	// has %X, %T, %F, %c, %j, etc
	// although it uses goto 26 times (in format.go), it has the best support
}
