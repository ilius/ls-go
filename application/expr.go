package application

import (
	"fmt"
	"io/fs"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ilius/expr"
	"github.com/ilius/expr/ast"
	"github.com/ilius/expr/vm"
	"github.com/ilius/go-table"
	"github.com/ilius/ls-go/common"
	"github.com/ilius/ls-go/lstime"
)

var (
	exprTimeFormat     = "2006-01-02 15:04:05.999999999 Z0700"
	exprFloatPrecision = 6
)

func NewExprGetter(colors bool, exprStr string) *ExprGetter {
	return &ExprGetter{
		prog:   compileExpr(exprStr),
		colors: colors,
		// env:
	}
}

func compileExpr(exprStr string) *vm.Program {
	prog, err := expr.Compile(
		exprStr,
		expr.Patch(&patcher{}),
	)
	check(err)
	return prog
}

type patcher struct{}

func (p *patcher) Enter(_ *ast.Node)    {}
func (p *patcher) Exit(node *ast.Node)  {}
func (p *patcher) Visit(node *ast.Node) {}

func parseDuration(s string) time.Duration {
	dur, err := time.ParseDuration(s)
	check(err)
	return dur
}

/*

ls-go --expr 'mtime().After(now().Add(duration("-1h")))' --sort=time --mtime


*/

type ExprGetter struct {
	prog   *vm.Program
	_type  reflect.Type
	colors bool
	// env map[string]any
}

func (f *ExprGetter) evaluateExpr(info FileInfo) (any, error) {
	value, err := expr.Run(f.prog, map[string]any{
		"info": info,
		"now":  *startTime,

		// shortcuts
		"name":     info.Name(),
		"size":     info.Size(),
		"mode":     info.Mode(),
		"basename": info.Basename(),
		"ext":      info.Ext(),
		"dir":      info.Dir(),

		"parsed_name": func() *common.ParsedName {
			return app.FileSystem.SplitExt(info.Name())
		},

		// functions for other columns
		"mtime": info.ModTime,
		"ctime": func() time.Time { return *info.CTime() },
		"atime": func() time.Time { return *info.ATime() },

		// other functions
		"past":   func(tm time.Time) bool { return tm.Before(*startTime) },
		"future": func(tm time.Time) bool { return tm.After(*startTime) },

		"duration": parseDuration,

		"split":      strings.Split,
		"path_split": app.FileSystem.SplitAll,

		"type": reflect.TypeOf,
	})
	if err != nil {
		return nil, err
	}
	// turns out operations don't work on functions
	// because the value type is determined before in expr.Run
	// and will give error if you try to do scalar operatins on a function
	// is there a way to lazily evaluate something?

	switch vt := value.(type) {
	case func() string:
		return vt(), nil
	case func() *time.Time:
		return vt(), nil
	case func() int64:
		return vt(), nil
	case func() fs.FileMode:
		return vt(), nil
	}
	return value, nil
}

func (f *ExprGetter) Type() (reflect.Type, error) {
	if f._type != nil {
		return f._type, nil
	}
	value, err := f.evaluateExpr(&FileInfoImp{
		FileInfo: &FileInfoLow{
			name:    "",
			size:    0,
			mode:    0,
			modTime: *startTime,
			isDir:   false,
			sys:     app.Platform.EmptyFileInfoSys(),
		},
	})
	if err != nil {
		return nil, err
	}
	f._type = reflect.TypeOf(value)
	return f._type, nil
}

func (f *ExprGetter) Alignment() (table.Alignment, error) {
	_type, err := f.Type()
	if err != nil {
		return nil, err
	}
	switch _type.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Uint64, reflect.Uint32:
		return table.AlignmentRight, nil
	case reflect.Float64, reflect.Float32:
		return table.AlignmentRight, nil
	case reflect.Bool:
		return table.AlignmentLeft, nil
	case reflect.String:
		return table.AlignmentLeft, nil
	case reflect.Ptr:
		return table.AlignmentLeft, nil
	case reflect.Slice, reflect.Array, reflect.Struct:
		return table.AlignmentLeft, nil
	}
	return table.AlignmentLeft, nil
}

func (f *ExprGetter) Value(item any) (any, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	return f.evaluateExpr(info)
}

func (f *ExprGetter) ValueBool(info FileInfo) (bool, error) {
	value, err := f.evaluateExpr(info)
	if err != nil {
		return false, nil
	}
	switch vt := value.(type) {
	case bool:
		return vt, nil
	case int:
		return vt != 0, nil
	case int8:
		return vt != 0, nil
	case int16:
		return vt != 0, nil
	case int32:
		return vt != 0, nil
	case int64:
		return vt != 0, nil
	case uint:
		return vt != 0, nil
	case uint8:
		return vt != 0, nil
	case uint16:
		return vt != 0, nil
	case uint32:
		return vt != 0, nil
	case uint64:
		return vt != 0, nil
	case float64:
		return vt != 0, nil
	case float32:
		return vt != 0, nil
	case string:
		return vt != "", nil
	case []byte:
		return len(vt) != 0, nil
	case []rune:
		return len(vt) != 0, nil
	}
	vv := reflect.ValueOf(value)
	switch vv.Kind() {
	case reflect.Ptr:
		return !vv.IsNil(), nil
	case reflect.Slice:
		return vv.Len() != 0, nil
	case reflect.Array:
		return vv.Len() != 0, nil
	case reflect.Struct:
		return true, nil
	case reflect.Bool:
		return vv.Bool(), nil
	case reflect.Int:
		return vv.Int() != 0, nil
	}
	return false, fmt.Errorf("ValueBool: unknow type %T", value)
}

func (f *ExprGetter) MustValueBool(info FileInfo) bool {
	b, err := f.ValueBool(info)
	check(err)
	return b
}

func (f *ExprGetter) ValueString(colName string, item any) (string, error) {
	info, ok := item.(FileInfo)
	if !ok {
		return "", fmt.Errorf("Value: invalid type %T, must be FileInfo", item)
	}
	value, err := f.evaluateExpr(info)
	if err != nil {
		return "", err
	}
	return app.FormatValue(colName, fmt.Sprintf("%v", value))
}

func (f *ExprGetter) Format(item any, value any) (string, error) {
	if !f.colors {
		return fmt.Sprintf("%v", value), nil
	}
	switch vt := value.(type) {
	case string:
		return app.Colorize(vt, colors.Expr.String), nil
	case int64:
		return app.Colorize(strconv.FormatInt(vt, 10), colors.Expr.Integer), nil
	case int:
		return app.Colorize(strconv.FormatInt(int64(vt), 10), colors.Expr.Integer), nil
	case uint:
		return app.Colorize(strconv.FormatInt(int64(vt), 10), colors.Expr.Integer), nil
	case uint64:
		return app.Colorize(strconv.FormatInt(int64(vt), 10), colors.Expr.Integer) + Reset, nil
	case int32:
		return app.Colorize(strconv.FormatInt(int64(vt), 10), colors.Expr.Integer) + Reset, nil
	case uint32:
		return app.Colorize(strconv.FormatInt(int64(vt), 10), colors.Expr.Integer) + Reset, nil
	case float64:
		return app.Colorize(strconv.FormatFloat(vt, 'f', exprFloatPrecision, 64), colors.Expr.Float), nil
	case *time.Time:
		if vt == nil {
			return "", nil
		}
		return app.Colorize(vt.Format(exprTimeFormat), colors.Expr.Time), nil
	case time.Time:
		return app.Colorize(vt.Format(exprTimeFormat), colors.Expr.Time), nil
	case time.Duration:
		return lstime.FormatDuration(vt), nil
	}
	return fmt.Sprintf("%v", value), nil
}
