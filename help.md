## Flags:

### `--all`, `-a`
Do not ignore entries starting with `.`

### `--almost-all`, `-A`
Do not list implied `.` and `..`

### `--sort=[|none|size|time|extension|kind|inode|links|filesize|mode|name-len]`
Sort by given column instead of basename.

### `--size`, `-s`
Print the size of each file.

### `--human-readable`, `-h`
With `-l` or `-s` / `--size`, print sizes like `1K`, `234M`, `2G`, etc.

### `--si`
Like `--human-readable`, but use powers of 1000, not 1024.

### `--bytes`
Print sizes in bytes.

### `--blocks`
Show allocated number of blocks (like `ls -s`) as a new column.

### `--time=[mtime|ctime|atime|status|change|access|use|modified|accessed]`
Change the default of using modification times.\
Access time: `atime`, `access`, `use`.\
Change time: `ctime`, `status`.\
With `-l`, it determines which time to show.\
With `--sort=time`, sort by given time (newest first).

### `--time-style=`
Time/date format with `-l`.\
See `README.md` for details.

### `--full-time`
Shortcut to `-l --time-style=full-iso`.

### `--mtime`, `--modified`
Include modification time (of file contents).

### `--ctime`, `--changed`
Include change time (of file contents or metadata).

### `--atime`, `--accessed`
Include access time.

### `--owner`
Include owner and group.

### `--group`
Show group (without long mode).

### `--no-group`, `-G`
Hide group name (with `-l`).

### `--numeric-uid-gid`, `-n`
Like `-l`, but list numeric user and group IDs.

### `--perm`, `--mode`
Include permissions for owner, group, and other.

### `--perm-oct`, `--mode-oct`, `--octal-permissions`
Include permissions / mode in octal format.

### `--inode`, `-i`
Print the index number (inode number) of each file.

### `--long`, `-l`
Include size, date, owner, and permissions.

### `--extra-long`
Include all columns.

### `--oneline`, `-1`
Show one file per line.

### `--horizontal`, `-x`
List entries by lines instead of by columns.

### `--vertical`, `--grid`
List entries by columns.

### `--compact`
Try to fit more columns in many-files-per-line modes (vertical/horizontal).

### `--vbar`
Show vertical bars between files in a row, or between columns in `--long` or `--oneline` mode.

### `--quoting-style=[|literal|shell|shell-always|shell-escape|shell-escape-always|c|escape|none]`
Use given quoting style for entry names (overrides QUOTING_STYLE environment variable).

### `--literal`, `-N`
Shortcut to `--quoting-style=literal`.\
Print entry names without quoting.

### `--escape`, `-b`
Shortcut to `--quoting-style=escape`.\
Print C-style escapes for nongraphic characters.

### `--directory`, `--list-dirs`, `-d`
List directories themselves, not their contents.

### `--dirs-first`, `--dir-first`, `--group-directories-first`
Show directories before files.

### `--dirs-only`, `--dir-only`, `--only-dirs`
Only show directories.

### `--files`
Only show files.

### `--has-mode=`
Only show items with mode(permissions) that contains the given octal mode.

### `--dereference`, `-L`
When showing file information for a symbolic link, show information for the file the link references rather than for the link itself.

### `--links`
Show paths for symlinks.

### `--link-rel`
Show symlinks as relative paths if shorter than absolute path.

### `--reverse`, `-r`
Reverse order while sorting.

### `--stats`
Show statistics.

### `--icons`
Show folder icon before directory name.

### `--nerd-font`
Show nerd font glyphs before file names.

### `--recursive`, `-R`
Traverse all directories recursively.

### `--find=`
Filter items with a regexp.

### `--color=[auto||always|y|yes|never|n|no]`
Whether or not to colorize the output.\
`auto` means if stdout connected to a terminal.

### `--header`
Add a header line with `-l` / `--long` or `-1` / `--oneline`.

### `--no-header`
Do not add a header line with `--csv` or `--json-array`.

### `--json`
Print JSON-encoded lines instead of tables (one object per line).

### `--json-array`
Print JSON-encoded lines instead of tables, one array per line.

### `--ascii`
With `--json` and `--json-array`, escape Unicode characters and ensure output is ASCII. In tabular/normal mode, apply this only to file names.

### `--csv`
Print a CSV table.

### `--html`
Print HTML.

### `--read-json`
Read JSON-encoded lines from stdin, instead of looking at filesystem and path arguments.

### `--minsize=0`
Minimum file size (in bytes).

### `--maxsize=0`
Maximum file size (in bytes).

### `-t`
Shortcut to `--sort=time`.\
Sort by time, newest first.\
See `--time`.

### `-c`
Shortcut to `--time=ctime`.\
With `-lt`: sort by, and show, ctime (time of last modification of file status information).\
With `-l`: show ctime and sort by name.\
Otherwise: sort by ctime, newest first.

### `-u`
Shortcut to `--time=use`.\
With `-lt`: sort by, and show, access time.\
With `-l`: show access time and sort by name.\
Otherwise: sort by access time, newest first.

### `-U`
Shortcut to `--sort=none`.\
Do not sort (list entries in directory order).

### `-S`
Shortcut to `--sort=size`.\
Sort by file size, largest first.

### `-X`
Shortcut to `--sort=extension`.\
Sort alphabetically by entry extension.

### `--colors-json`
Print colors in json format and exit.

### `--cpuprofile=`
Write cpu profile to file.

### `--expr=`
An expression to be evaluated as a new column.

### `--where=`
An expression to be evaluated and filter files by.

### `--help-md`
Show help in markdown format.

### `--help`
Show usage message.

### `--version`
Show version.


