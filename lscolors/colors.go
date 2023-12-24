package lscolors

import (
	"fmt"
	"strconv"
)

// see the color codes
// http://i.stack.imgur.com/UQVe5.png

const (
	DEFAULT = "_default"
	SELF    = "_self"
)

type Style struct {
	Fg   uint8 `json:"fg,omitempty"`
	Bg   uint8 `json:"bg,omitempty"`
	Bold bool  `json:"bold,omitempty"`
}

func (s *Style) SetFg(code uint8) *Style {
	s.Fg = code
	return s
}

func (s *Style) SetBg(code uint8) *Style {
	s.Bg = code
	return s
}

func (s *Style) SetBold() *Style {
	s.Bold = true
	return s
}

func (s *Style) S() string {
	st := ""
	// color code 0 is black, but code 16 is also black
	// besides we never set color black as fg or bg
	if s.Fg > 0 {
		st += "\x1b[38;5;" + strconv.FormatUint(uint64(s.Fg), 10) + "m"
	}
	if s.Bg > 0 {
		st += "\x1b[48;5;" + strconv.FormatUint(uint64(s.Bg), 10) + "m"
	}
	if s.Bold {
		st += "\x1b[1m"
	}
	return st
}

var defaultStyle = &Style{}

type StyleMap map[string]*Style

func (sm StyleMap) Get(key string) *Style {
	s := sm[key]
	if s != nil {
		return s
	}
	return sm.Default()
}

func (sm StyleMap) Default() *Style {
	s := sm[DEFAULT]
	if s != nil {
		return s
	}
	return defaultStyle
}

// Fg wraps an 8-bit foreground color code in the ANSI escape sequence
func Fg(code uint8) *Style {
	return &Style{Fg: code}
}

// Bg wraps an 8-bit background color code in the ANSI escape sequence
func Bg(code uint8) *Style {
	return &Style{Bg: code}
}

// Rgb2code converts RGB values (up to 5) to an 8-bit color code
/*
func Rgb2code(r uint8, g uint8, b uint8) uint8 {
	code := 36*r + 6*g + b + 16
	if code < 16 || 231 < code {
		panic(fmt.Errorf("Invalid RGB values (%d, %d, %d)", r, g, b))
	}
	return code
}
*/

// Gray converts a scalar of "grayness" to an 8-bit color code
func Gray(lightness uint8) uint8 {
	code := lightness + 232
	if code < 232 || 255 < code {
		panic(fmt.Errorf("Invalid lightness value (%d) for gray", lightness))
	}
	return code
}

// FgGray converts a scalar of "grayness" to an ANSI-escaped foreground 8-bit color code
func FgGray(lightness uint8) *Style {
	return Fg(Gray(lightness))
}

// BgGray converts a scalar of "grayness" to an ANSI-escaped background 8-bit color code
func BgGray(lightness uint8) *Style {
	return Bg(Gray(lightness))
}

type LightDark struct {
	Light *Style `json:"light"`
	Dark  *Style `json:"dark"`
}

// returns (mainColor, accentColor)
func (ld LightDark) Get() (*Style, *Style) {
	return ld.Light, ld.Dark
}

type DirColors struct {
	Name       *Style `json:"name"`
	Ext        *Style `json:"ext"`
	HiddenName *Style `json:"hidden_name"`
	HiddenExt  *Style `json:"hidden_ext"`
}

type LinkColors struct {
	Name    *Style `json:"name"`
	NameDir *Style `json:"name_dir"`
	Arrow   *Style `json:"arrow"`
	Path    *Style `json:"path"`
	Broken  *Style `json:"broken"`
}

type TimeColors struct {
	Year        *Style `json:"year"`
	Number      *Style `json:"number"`
	NumberColon *Style `json:"number_colon"`
	NumberSlash *Style `json:"number_slash"`
	Word        *Style `json:"word"`
}

// PermColor holds color mappings for users and groups
type PermColors struct {
	User  StyleMap `json:"user"`
	Group StyleMap `json:"group"`
	Other StyleMap `json:"other"`
}

type FolderHeaderColors struct {
	Arrow      *Style `json:"arrow"`
	Main       *Style `json:"main"`
	Slash      *Style `json:"slash"`
	LastFolder *Style `json:"last_folder"`
	Error      *Style `json:"error"`
}

type StatsColors struct {
	Text   *Style `json:"text"`
	Number *Style `json:"number"`
	MS     *Style `json:"ms"`
}

type LightDarkMap map[string]LightDark

type ExprColors struct {
	String  *Style `json:"string"`
	Integer *Style `json:"integer"`
	Float   *Style `json:"float"`
	Time    *Style `json:"time"`
}

type TabularColors struct {
	FolderHeader FolderHeaderColors `json:"folder_header"`
	TableHeader  *Style             `json:"table_header"`
}

type HtmlColors struct {
	Default      *Style             `json:"default"`
	FolderHeader FolderHeaderColors `json:"folder_header"`
	TableHeader  *Style             `json:"table_header"`
}

type Colors struct {
	File LightDarkMap `json:"file"`
	Dir  DirColors    `json:"dir"`
	Link LinkColors   `json:"link"`
	Size StyleMap     `json:"size"`
	Time TimeColors   `json:"time"`
	Perm PermColors   `json:"perm"`
	Expr ExprColors   `json:"expr"`

	Tabular *TabularColors `json:"tabular"`
	Html    *HtmlColors    `json:"html"`

	Device *Style `json:"device"`
	Socket *Style `json:"socket"`
	Pipe   *Style `json:"pipe"`

	Stats StatsColors `json:"stats"`
}
