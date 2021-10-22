package eval

import (
	gofmt "fmt"
	"os"
)

func NewFmtObj() *FmtObj {
	ret := &FmtObj{}
	SetGlobalObj(fmt_name, ret)

	return ret
}

const (
	FMT_OBJ  = "FMT_OBJ"
	fmt_name = "fmt"
)

type FmtObj struct {
}

func (f *FmtObj) Inspect() string  { return "<" + fmt_name + ">" }
func (f *FmtObj) Type() ObjectType { return FMT_OBJ }

func (f *FmtObj) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "errorf":
		return f.Errorf(line, scope, args...)
	case "print":
		return f.Print(line, scope, args...)
	case "printf":
		return f.Printf(line, scope, args...)
	case "println":
		return f.Println(line, scope, args...)
	case "sprint":
		return f.Sprint(line, scope, args...)
	case "sprintf":
		return f.Sprintf(line, scope, args...)
	case "sprintln":
		return f.Sprintln(line, scope, args...)
	case "fprint":
		return f.Fprint(line, scope, args...)
	case "fprintf":
		return f.Fprintf(line, scope, args...)
	case "fprintln":
		return f.Fprintln(line, scope, args...)
	}
	return NewError(line, NOMETHODERROR, method, f.Type())
}

func (f *FmtObj) Errorf(line string, scope *Scope, args ...Object) Object {
	if len(args) < 1 {
		return NewError(line, ARGUMENTERROR, ">0", len(args))
	}

	formatObj, ok := args[0].(*String)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "errorf", "*String", args[0].Type())
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	err := gofmt.Errorf(formatObj.String, wrapped...)
	return NewError(line, GENERICERROR, err.Error())
}

func (f *FmtObj) Print(line string, scope *Scope, args ...Object) Object {
	if len(args) == 0 {
		n, err := gofmt.Fprint(scope.Writer)
		if err != nil {
			return NewNil(err.Error())
		}
		return NewInteger(int64(n))
	}

	format, wrapped := correctPrintResult(false, args...)
	n, err := gofmt.Fprintf(scope.Writer, format, wrapped...)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (f *FmtObj) Printf(line string, scope *Scope, args ...Object) Object {
	if len(args) < 1 {
		return NewError(line, ARGUMENTERROR, ">0", len(args))
	}

	formatObj, ok := args[0].(*String)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "printf", "*String", args[0].Type())
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	formatStr := formatObj.String
	if len(subArgs) == 0 {
		if REPLColor {
			formatStr = "\033[1;" + colorMap["STRING"] + "m" + formatStr + "\033[0m"
		}
	}
	n, err := gofmt.Fprintf(scope.Writer, formatStr, wrapped...)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (f *FmtObj) Println(line string, scope *Scope, args ...Object) Object {
	if len(args) == 0 {
		n, err := gofmt.Fprintln(scope.Writer)
		if err != nil {
			return NewNil(err.Error())
		}
		return NewInteger(int64(n))
	}

	format, wrapped := correctPrintResult(true, args...)
	n, err := gofmt.Fprintf(scope.Writer, format, wrapped...)
	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (f *FmtObj) Sprint(line string, scope *Scope, args ...Object) Object {
	if len(args) == 0 {
		return NewString("")
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	ret := gofmt.Sprint(wrapped...)
	return NewString(ret)
}

func (f *FmtObj) Sprintf(line string, scope *Scope, args ...Object) Object {
	if len(args) < 1 {
		return NewError(line, ARGUMENTERROR, ">0", len(args))
	}

	format, ok := args[0].(*String)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "sprintf", "*String", args[0].Type())
	}

	subArgs := args[1:]
	wrapped := make([]interface{}, len(subArgs))
	for i, v := range subArgs {
		wrapped[i] = &Formatter{Obj: v}
	}

	ret := gofmt.Sprintf(format.String, wrapped...)
	return NewString(ret)
}

func (f *FmtObj) Sprintln(line string, scope *Scope, args ...Object) Object {
	if len(args) == 0 {
		ret := gofmt.Sprintln()
		return NewString(ret)
	}

	wrapped := make([]interface{}, len(args))
	for i, v := range args {
		wrapped[i] = &Formatter{Obj: v}
	}

	ret := gofmt.Sprintln(wrapped...)
	return NewString(ret)
}

func (f *FmtObj) Fprint(line string, scope *Scope, args ...Object) Object {
	if len(args) < 2 {
		return NewError(line, ARGUMENTERROR, ">=2", len(args))
	}

	w, ok := args[0].(Writable)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "fprint", "Writable", args[0].Type())
	}

	subArgs := args[1:]

	var n int
	var err error

	writer := w.(Writable).IOWriter()
	if writer == os.Stdout || writer == os.Stderr { //output to stdout or stderr
		format, wrapped := correctPrintResult(false, subArgs...)
		n, err = gofmt.Fprintf(scope.Writer, format, wrapped...)
	} else {
		wrapped := make([]interface{}, len(subArgs))
		for i, v := range subArgs {
			wrapped[i] = &Formatter{Obj: v}
		}
		n, err = gofmt.Fprint(writer, wrapped...)
	}

	if err != nil {
		return NewNil(err.Error())
	}

	return NewInteger(int64(n))
}

func (f *FmtObj) Fprintf(line string, scope *Scope, args ...Object) Object {
	if len(args) < 2 {
		return NewError(line, ARGUMENTERROR, ">=2", len(args))
	}

	w, ok := args[0].(Writable)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "fprintf", "Writable", args[0].Type())
	}

	formatObj, ok := args[1].(*String)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "second", "fprintf", "*String", args[1].Type())
	}

	var n int
	var err error

	writer := w.(Writable).IOWriter()
	if len(args) > 2 { //has format
		if writer == os.Stdout || writer == os.Stderr { //output to stdout or stderr
			return f.Printf(line, scope, args[1:]...)
		} else {
			subArgs := args[2:]
			wrapped := make([]interface{}, len(subArgs))
			for i, v := range subArgs {
				wrapped[i] = &Formatter{Obj: v}
			}
			n, err = gofmt.Fprintf(writer, formatObj.String, wrapped...)
		}
	} else { //only string with no format, e.g. fmt.fprintf(stdout, "Hello world\n")
		formatStr := formatObj.String
		if writer == os.Stdout || writer == os.Stderr { //output to stdout or stderr
			return f.Printf(line, scope, args[1:]...)
		}
		n, err = gofmt.Fprintf(writer, formatStr)
	}

	if err != nil {
		return NewNil(err.Error())
	}
	return NewInteger(int64(n))
}

func (f *FmtObj) Fprintln(line string, scope *Scope, args ...Object) Object {
	if len(args) < 2 {
		return NewError(line, ARGUMENTERROR, ">=2", len(args))
	}

	w, ok := args[0].(Writable)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "fprintln", "Writable", args[0].Type())
	}

	subArgs := args[1:]

	var n int
	var err error

	writer := w.(Writable).IOWriter()
	if writer == os.Stdout || writer == os.Stderr { //output to stdout or stderr
		format, wrapped := correctPrintResult(true, subArgs...)
		n, err = gofmt.Fprintf(scope.Writer, format, wrapped...)
	} else {
		wrapped := make([]interface{}, len(subArgs))
		for i, v := range subArgs {
			wrapped[i] = &Formatter{Obj: v}
		}

		n, err = gofmt.Fprintln(writer, wrapped...)
	}

	if err != nil {
		return NewNil(err.Error())
	}
	return NewInteger(int64(n))
}
