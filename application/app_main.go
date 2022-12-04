package application

import (
	"encoding/json"
	"io"
	"os"
	"runtime/pprof"
	"time"

	"github.com/ilius/ls-go/lsargs"
)

var (
	startTime *time.Time

	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr

	args = lsargs.New()

	app *Application
)

func Run(rawArgs []string) {
	now := time.Now()
	startTime = &now

	// parse the arguments and populate the struct
	args.Parse(rawArgs, VERSION)

	{
		jsonStr := os.Getenv("LSGO_COLORS")
		if jsonStr != "" {
			err := json.Unmarshal([]byte(jsonStr), &colors)
			check(err)
		}
	}

	if *args.ColorsJson {
		printColorsJson()
		os.Exit(0)
	}

	if *args.CpuProfile != "" {
		// runtime.SetCPUProfileRate(500)
		f, err := os.Create(*args.CpuProfile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	app = NewApplication()
	tableSpec := app.PostParse(args)

	colorsEnable, err := app.Terminal.ColorsEnabled(*args.Color)
	check(err)
	stdout, stderr = app.Platform.OutputAndError(colorsEnable)

	app.ListMain(tableSpec)

	app.PrintErrors()
	app.Exit()
}
