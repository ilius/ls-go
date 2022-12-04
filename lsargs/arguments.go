package lsargs

import (
	"fmt"
	"os"

	goopt "github.com/ilius/goopt"
	. "github.com/ilius/ls-go/common"
)

// declare the struct that holds all the arguments
type Arguments struct {
	All       *bool
	AlmostAll *bool
	Sort      *string
	Size      *bool
	Human     *bool
	SI        *bool
	Bytes     *bool
	Blocks    *bool

	Time      *string
	TimeStyle *string
	FullTime  *bool
	Mtime     *bool
	Ctime     *bool
	Atime     *bool

	Owner         *bool
	Group         *bool
	NoGroup       *bool
	NumericUidGid *bool
	ModeOct       *bool
	Mode          *bool
	Inode         *bool

	Long       *bool
	ExtraLong  *bool
	SingleCol  *bool
	Horizontal *bool
	Vertical   *bool
	Compact    *bool
	Vbar       *bool

	QuotingStyle *string

	Shortcut_literal *bool
	Shortcut_escape  *bool

	Directory *bool
	DirsFirst *bool
	DirsOnly  *bool
	FilesOnly *bool
	HasMode   *string

	Dereference *bool
	Links       *bool
	LinkRel     *bool

	Reverse   *bool
	Stats     *bool
	Icons     *bool
	Nerdfont  *bool
	Recursive *bool
	Find      *string
	Color     *string

	Header   *bool
	NoHeader *bool

	Json      *bool
	JsonArray *bool
	ASCII     *bool
	Csv       *bool
	Html      *bool

	ReadJson *bool

	Minsize *int
	Maxsize *int

	Shortcut_t *bool
	Shortcut_c *bool
	Shortcut_u *bool
	Shortcut_U *bool
	Shortcut_S *bool
	Shortcut_X *bool

	ColorsJson *bool

	Expr  *string
	Where *string

	CpuProfile *string

	Paths []string
}

const time_flag_desc = `Change the default of using modification times; Access time: 'atime', 'access', 'use'; Change time: 'ctime', 'status'; With -l, it determines which time to show; With --sort=time, sort by given time (newest first)`

// TODO
// birth time: birth, creation;

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func New() *Arguments {
	return &Arguments{
		All: goopt.Flag(
			[]string{"--all", "-a"},
			nil,
			"Do not ignore entries starting with '.'",
			"",
		),
		AlmostAll: goopt.Flag(
			[]string{"--almost-all", "-A"},
			nil,
			"Do not list implied '.' and '..'",
			"",
		),
		Sort: goopt.Alternatives(
			[]string{"--sort"},
			// first item must be "" because that is the value we get when flag is not passed
			[]string{
				"",
				S_NONE,
				S_SIZE,
				S_TIME,
				S_EXTENSION,
				S_KIND,
				S_INODE,
				S_LINKS,
				S_FILESIZE,
				S_MODE,
				S_NAME_LEN,
			},
			"Sort by given column instead of basename",
		),
		Size: goopt.Flag(
			[]string{"--size", "-s"},
			nil,
			"Print the size of each file",
			"",
		),
		Human: goopt.Flag(
			[]string{"--human-readable", "-h"},
			nil,
			"With -l or -s / --size, print sizes like '1K', '234M', '2G', etc",
			"",
		),
		SI: goopt.Flag(
			[]string{"--si"},
			nil,
			"Like --human-readable, but use powers of 1000, not 1024",
			"",
		),
		Bytes: goopt.Flag(
			[]string{"--bytes"},
			nil,
			"Print sizes in bytes",
			"",
		),
		Blocks: goopt.Flag(
			[]string{"--blocks"},
			nil,
			"Show allocated number of blocks (like ls -s) as a new column",
			"",
		),
		Time: goopt.Alternatives(
			[]string{"--time"},
			[]string{
				"mtime", "ctime", "atime", // main names
				"status", "change", "access", "use", // supported by ls
				"modified", "accessed", // used by exa
			},
			time_flag_desc,
		),
		TimeStyle: goopt.String(
			[]string{"--time-style"},
			"",
			"Time/date format with -l; See 'README.md' for details.",
		),
		FullTime: goopt.Flag(
			[]string{"--full-time"},
			nil,
			"Shortcut to -l --time-style=full-iso",
			"",
		),
		Mtime: goopt.Flag(
			[]string{"--mtime", "--modified"},
			nil,
			"Include modification time (of file contents)",
			"",
		),
		Ctime: goopt.Flag(
			[]string{"--ctime", "--changed"},
			nil,
			"Include change time (of file contents or metadata)",
			"",
		),
		Atime: goopt.Flag(
			[]string{"--atime", "--accessed"},
			nil,
			"Include access time",
			"",
		),
		Owner: goopt.Flag(
			[]string{"--owner"},
			nil,
			"Include owner and group",
			"",
		),
		Group: goopt.Flag(
			[]string{"--group"},
			nil,
			"Show group (without long mode)",
			"",
		),
		NoGroup: goopt.Flag(
			[]string{"--no-group", "-G"},
			nil,
			"Hide group name (with -l)",
			"",
		),
		NumericUidGid: goopt.Flag(
			[]string{"-n", "--numeric-uid-gid"},
			nil,
			"like -l, but list numeric user and group IDs",
			"",
		),
		Mode: goopt.Flag(
			[]string{"--perm", "--mode"},
			nil,
			"Include permissions for owner, group, and other",
			"",
		),
		ModeOct: goopt.Flag(
			[]string{
				"--perm-oct",
				"--mode-oct",
				"--octal-permissions",
				// --oct works, but not if we have 2 flags starting with --oct
				// even if both belong to the same parameter / goopt.Flag
				// this needs fixing in goopt
			},
			nil,
			"Include permissions / mode in octal format",
			"",
		),
		Inode: goopt.Flag(
			[]string{"--inode", "-i"},
			nil,
			"Print the index number (inode number) of each file",
			"",
		),
		Long: goopt.Flag(
			[]string{"--long", "-l"},
			nil,
			"Include size, date, owner, and permissions",
			"",
		),
		ExtraLong: goopt.Flag(
			[]string{"--extra-long"},
			nil,
			"Include all columns",
			"",
		),
		SingleCol: goopt.Flag(
			[]string{"-1", "--oneline"},
			// commands that have --oneline as long flag for -1
			// 		git log --oneline
			// 		exa --oneline
			// ls has no long bool flag for it,
			// except `ls --format=single` or `ls --format=single-column`
			nil,
			"Show one file per line",
			"",
		),
		// like ls -x or ls --format=horizontal
		Horizontal: goopt.Flag(
			[]string{"--horizontal", "-x"},
			nil,
			"list entries by lines instead of by columns",
			"",
		),
		// like ls -C or ls --format=vertical
		Vertical: goopt.Flag(
			[]string{
				"--vertical",
				"--grid", // like exa
			},
			nil,
			"list entries by columns",
			"",
		),
		Compact: goopt.Flag(
			[]string{"--compact"},
			nil,
			"try to fit more columns in many-files-per-line modes (vertical/horizontal)",
			"",
		),
		Vbar: goopt.Flag(
			[]string{"--vbar"},
			nil,
			"show vertical bars between files in a row, or between columns in '--long' or '--oneline' mode",
			"",
		),
		QuotingStyle: goopt.Alternatives(
			[]string{"--quoting-style"},
			[]string{
				"", // default, must be first
				E_literal,
				E_shell,
				E_shell_always,
				E_shell_escape,
				E_shell_escape_always,
				E_c,
				E_escape,
				E_none,
			},
			"use given quoting style for entry names (overrides QUOTING_STYLE environment variable)",
		),
		Shortcut_literal: goopt.Flag(
			[]string{"--literal", "-N"},
			nil,
			"Shortcut to --quoting-style=literal; Print entry names without quoting",
			"",
		),
		Shortcut_escape: goopt.Flag(
			[]string{"--escape", "-b"},
			nil,
			"Shortcut to --quoting-style=escape; Print C-style escapes for nongraphic characters",
			"",
		),
		Directory: goopt.Flag(
			[]string{"--directory", "-d", "--list-dirs"},
			nil,
			"List directories themselves, not their contents",
			"",
		),

		DirsFirst: goopt.Flag(
			[]string{"--dirs-first", "--dir-first", "--group-directories-first"},
			nil,
			"Show directories before files",
			"",
		),
		DirsOnly: goopt.Flag(
			[]string{"--dirs-only", "--dir-only", "--only-dirs"},
			nil,
			"Only show directories",
			"",
		),
		FilesOnly: goopt.Flag(
			[]string{"--files"},
			nil,
			"Only show files",
			"",
		),
		HasMode: goopt.String(
			[]string{"--has-mode"},
			"",
			"Only show items with mode(permissions) that contains the given octal mode",
		),

		Dereference: goopt.Flag(
			[]string{"--dereference", "-L"},
			nil,
			"When showing file information for a symbolic link, show information for the file the link references rather than for the link itself",
			"",
		),
		Links: goopt.Flag(
			[]string{"--links"},
			nil,
			"Show paths for symlinks",
			"",
		),
		LinkRel: goopt.Flag(
			[]string{"--link-rel"},
			nil,
			"Show symlinks as relative paths if shorter than absolute path",
			"",
		),
		Reverse: goopt.Flag(
			[]string{"--reverse", "-r"},
			nil,
			"Reverse order while sorting",
			"",
		),
		Stats: goopt.Flag(
			[]string{"--stats"},
			nil,
			"Show statistics",
			"",
		),
		Icons: goopt.Flag(
			[]string{"--icons"},
			nil,
			"Show folder icon before directory name",
			"",
		),
		Nerdfont: goopt.Flag(
			[]string{"--nerd-font"},
			nil,
			"Show nerd font glyphs before file names",
			"",
		),
		Recursive: goopt.Flag(
			[]string{
				"--recursive", "-R",
				// "--recurse",
				// if you add --recurse, then --re, --rec, --recu --recur --recurs
				// will not work anymore
				// gotta make the goopt smarter about this
			},
			nil,
			"Traverse all directories recursively",
			"",
		),
		Find: goopt.String(
			[]string{"--find"},
			"",
			"Filter items with a regexp",
		),
		Color: goopt.Alternatives(
			[]string{"--color"},
			[]string{
				"auto", // default, must be first
				"",     // to allow --color=
				"always", "y", "yes",
				"never", "n", "no",
			},
			"Whether or not to colorize the output; 'auto' means if stdout connected to a terminal",
		),
		Header: goopt.Flag(
			[]string{"--header"},
			nil,
			"Add a header line with '-l' / '--long' or '-1' / '--oneline'",
			"",
		),
		NoHeader: goopt.Flag(
			[]string{"--no-header"},
			nil,
			"Do not add a header line with '--csv' or '--json-array'",
			"",
		),
		Json: goopt.Flag(
			[]string{"--json"},
			nil,
			"Print JSON-encoded lines instead of tables (one object per line)",
			"",
		),
		JsonArray: goopt.Flag(
			[]string{"--json-array"},
			nil,
			"Print JSON-encoded lines instead of tables, one array per line",
			"",
		),
		ASCII: goopt.Flag(
			[]string{"--ascii"},
			nil,
			"With --json and --json-array, escape Unicode characters and ensure output is ASCII. In tabular/normal mode, apply this only to file names.",
			"",
		),
		Csv: goopt.Flag(
			[]string{"--csv"},
			nil,
			"Print a CSV table",
			"",
		),
		Html: goopt.Flag(
			[]string{"--html"},
			nil,
			"Print HTML",
			"",
		),

		ReadJson: goopt.Flag(
			[]string{"--read-json"},
			nil,
			"Read JSON-encoded lines from stdin, instead of looking at filesystem and path arguments",
			"",
		),

		Minsize: goopt.Int(
			[]string{"--minsize"},
			0,
			"minimum file size (in bytes)",
		),
		Maxsize: goopt.Int(
			[]string{"--maxsize"},
			0,
			"maximum file size (in bytes)",
		),

		Shortcut_t: goopt.Flag(
			[]string{"-t"},
			nil,
			`Shortcut to --sort=time; Sort by time, newest first; See --time`,
			"",
		),

		Shortcut_c: goopt.Flag(
			[]string{"-c"},
			nil,
			`Shortcut to --time=ctime; With -lt: sort by, and show, ctime (time of last modification of file status information); With -l: show ctime and sort by name; Otherwise: sort by ctime, newest first`,
			"",
		),

		Shortcut_u: goopt.Flag(
			[]string{"-u"},
			nil,
			`Shortcut to --time=use; With -lt: sort by, and show, access time; With -l: show access time and sort by name; Otherwise: sort by access time, newest first`,
			"",
		),

		Shortcut_U: goopt.Flag(
			[]string{"-U"},
			nil,
			`Shortcut to --sort=none; Do not sort (list entries in directory order)`,
			"",
		),

		Shortcut_S: goopt.Flag(
			[]string{"-S"},
			nil,
			`Shortcut to --sort=size; Sort by file size, largest first`,
			"",
		),

		Shortcut_X: goopt.Flag(
			[]string{"-X"},
			nil,
			`Shortcut to --sort=extension; Sort alphabetically by entry extension`,
			"",
		),

		ColorsJson: goopt.Flag(
			[]string{"--colors-json"},
			nil,
			`Print colors in json format and exit`,
			"",
		),

		CpuProfile: goopt.String(
			[]string{"--cpuprofile"},
			"",
			"Write cpu profile to file",
		),

		Expr: goopt.String(
			[]string{"--expr"},
			"",
			"An expression to be evaluated as a new column",
		),
		Where: goopt.String(
			[]string{"--where"},
			"",
			"An expression to be evaluated and filter files by",
		),
	}
}

func (args *Arguments) Parse(rawArgs []string, version string) {
	help_md := goopt.Flag([]string{"--help-md"}, nil, "show help in markdown format", "")

	// remove -h as Shortcut to --help, because we want -h as --human-readable
	goopt.SetHelpFlags([]string{"--help"})

	// set version for --version
	goopt.Version = version

	goopt.Parse(rawArgs, nil)

	if *help_md {
		helpMarkdown()
		os.Exit(0)
	}

	if *args.Horizontal && *args.Vertical {
		fmt.Fprintln(os.Stderr, "Conflicting flags: --horizontal and --vertical")
		os.Exit(0)
	}

	paths := goopt.Args
	if len(paths) == 0 {
		paths = []string{"."}
	}
	args.Paths = paths
}
