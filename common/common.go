package common

type SizeFormat uint8

const (
	SizeFormatLegacy  = SizeFormat(0)
	SizeFormatMetric  = SizeFormat(1)
	SizeFormatInteger = SizeFormat(3)
)

// column names
const (
	C_Inode      = "inode"
	C_ModeOct    = "mode_oct"
	C_Mode       = "mode"
	C_HardLinks  = "hard_links"
	C_Owner      = "owner"
	C_Group      = "group"
	C_Blocks     = "blocks"
	C_Size       = "size"
	C_MTime      = "mtime"
	C_CTime      = "ctime"
	C_ATime      = "atime"
	C_Name       = "name"
	C_LinkTarget = "link_target"
)

// quoting styles
const (
	E_none                = "none"
	E_literal             = "literal"
	E_locale              = "locale"
	E_shell               = "shell"
	E_shell_always        = "shell-always"
	E_shell_escape        = "shell-escape"
	E_shell_escape_always = "shell-escape-always"
	E_c                   = "c"
	E_escape              = "escape"
)
