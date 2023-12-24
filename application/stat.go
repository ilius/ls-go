package application

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func printStats(colorsEnable bool, numFiles int, numDirs int) {
	if !colorsEnable {
		printStatsNoColor(numFiles, numDirs)
		return
	}
	c := colors.Stats
	duration := time.Since(*startTime)
	milliSeconds := float64(duration.Microseconds()) / 1000.0
	statStrings := []string{
		app.Colorize(strconv.FormatInt(int64(numDirs), 10), c.Number),
		app.Colorize("dirs", c.Text),
		app.Colorize(strconv.FormatInt(int64(numFiles), 10), c.Number),
		app.Colorize("files", c.Text),
		app.Colorize(strconv.FormatFloat(milliSeconds, 'f', 2, 64), c.MS),
		app.Colorize("ms", c.Text),
	}
	fmt.Fprintln(stdout, strings.Join(statStrings, " "))
}

func printStatsNoColor(numFiles int, numDirs int) {
	duration := time.Since(*startTime)
	milliSeconds := float64(duration.Microseconds()) / 1000.0
	statStrings := []string{
		strconv.FormatInt(int64(numDirs), 10),
		"dirs",
		strconv.FormatInt(int64(numFiles), 10),
		"files",
		strconv.FormatFloat(milliSeconds, 'f', 2, 64),
		"ms",
	}
	fmt.Fprintln(stdout, strings.Join(statStrings, " "))
}
