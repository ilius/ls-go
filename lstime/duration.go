package lstime

import (
	"fmt"
	"strings"
	"time"
)

const Day = 24 * time.Hour

// example: ls-go --expr 'mtime().Sub(now())'
func FormatDuration(dur time.Duration) string {
	past := false
	if dur < 0 {
		past = true
		dur = -dur
	}
	parts := []string{}
	days := dur / Day
	if days == 1 {
		parts = append(parts, "1 day")
		dur -= Day
	} else if days > 1 {
		parts = append(parts, fmt.Sprintf("%d days", days))
		dur -= days * Day
	}
	if dur >= time.Second {
		hours := dur / time.Hour
		dur -= hours * time.Hour
		minutes := dur / time.Minute
		dur -= minutes * time.Minute
		parts = append(parts, fmt.Sprintf(
			"%.2d:%.2d:%.2d",
			hours,
			minutes,
			int(dur.Seconds()),
		))
	}
	// colors.Duration.Days
	// colors.Duration.Hours
	// colors.Duration.Minutes
	// colors.Duration.Seconds
	// colors.Duration.Colon
	str := strings.Join(parts, ", ")
	if past {
		str += " ago"
	}
	return str
}
