package eval

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"io/ioutil"
	"os"
	_ "reflect"
	"unicode/utf8"
)

const (
	JSON_OBJ  = "JSON_OBJ"
	json_name = "json"
)

type Json struct {
}

func NewJsonObj() Object {
	ret := &Json{}
	SetGlobalObj(json_name, ret)

	return ret
}

func (j *Json) Inspect() string  { return "<" + json_name + ">" }
func (j *Json) Type() ObjectType { return JSON_OBJ }

func (j *Json) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "marshal", "toJson", "stringify":
		return j.Marshal(line, args...)
	case "unmarshal", "fromJson", "parse":
		return j.UnMarshal(line, args...)
	case "indent":
		return j.Indent(line, args...)
	case "read": // read from a string
		return j.Read(line, args...)
	case "readFile":
		return j.ReadFile(line, args...)
	case "writeFile":
		return j.WriteFile(line, args...)
	}
	panic(NewError(line, NOMETHODERROR, method, j.Type()))
}

func (j *Json) Marshal(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	switch args[0].(type) {
	case *Integer:
		value := args[0].(*Integer)
		res, err := value.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewString(string(res))
	case *UInteger:
		value := args[0].(*UInteger)
		res, err := value.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewString(string(res))
	case *Float:
		value := args[0].(*Float)
		res, err := value.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewString(string(res))
	case *String:
		value := args[0].(*String)
		res, err := value.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewString(string(res))
	case *Boolean:
		value := args[0].(*Boolean)
		res, err := value.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewString(string(res))
	case *Array:
		value := args[0].(*Array)
		res, err := value.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewString(string(res))
	case *Tuple:
		value := args[0].(*Tuple)
		res, err := value.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewString(string(res))
	case *Hash:
		value := args[0].(*Hash)
		res, err := value.MarshalJSON()
		if err != nil {
			return NewNil(err.Error())
		}
		return NewString(string(res))
	default:
		panic(NewError(line, JSONERROR))
	}
	return NIL
}

func (j *Json) UnMarshal(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	jsonStr, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "unmarshal", "*String", args[0].Type()))
	}

	in := []byte(jsonStr.String)
	b := bytes.TrimSpace(in)
	r, _ := utf8.DecodeRune(b)

	if r == '[' { // array
		a := &Array{}
		err := a.UnmarshalJSON(b)
		if err != nil {
			return NewNil(err.Error())
		}
		return a
	} else if r == '{' { // hash
		h := NewHash()
		err := h.UnmarshalJSON(b)
		if err != nil {
			return NewNil(err.Error())
		}
		return h
	} else { //simple types, e.g. number, string
		var val interface{}
		err := json.Unmarshal(b, &val)
		if err != nil {
			return NewNil(err.Error())
		}
		ret, err := unmarshalJsonObject(val)
		if err != nil {
			return NewNil(err.Error())
		}
		return ret
	}
}

func (j *Json) Indent(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		panic(NewError(line, ARGUMENTERROR, "1|2", len(args)))
	}

	val, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "indent", "*String", args[0].Type()))
	}

	var indent string = "\t"
	if len(args) == 2 {
		str, ok := args[1].(*String)
		if !ok {
			panic(NewError(line, PARAMTYPEERROR, "second", "indent", "*String", args[1].Type()))
		}
		indent = str.String
	}

	var out bytes.Buffer
	b := []byte(val.String)

	err := json.Indent(&out, b, "", indent)
	if err != nil {
		return NewNil(err.Error())
	}
	return NewString(out.String())
}

func (j *Json) Read(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "read", "*String", args[0].Type()))
	}

	return j.UnMarshal(line, strObj)
}

func (j *Json) ReadFile(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	strObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "readFile", "*String", args[0].Type()))
	}

	byteValue, err := ioutil.ReadFile(strObj.String)
	if err != nil {
		return NewNil(err.Error())
	}

	return j.UnMarshal(line, NewString(string(byteValue)))
}

func (j *Json) WriteFile(line string, args ...Object) Object {
	if len(args) != 2 && len(args) != 3 {
		panic(NewError(line, ARGUMENTERROR, "2|3", len(args)))
	}

	fileNameObj, ok := args[0].(*String)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "writeFile", "*String", args[0].Type()))
	}

	permObj, ok := args[2].(*Integer)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "third", "writeFile", "*Integer", args[2].Type()))
	}

	v := j.Marshal(line, args[1])
	if v.Type() == NIL_OBJ {
		return v
	}

	err := ioutil.WriteFile(fileNameObj.String, []byte(v.(*String).String), os.FileMode(permObj.Int64))
	if err != nil {
		return NewFalseObj(err.Error())
	}

	return TRUE
}
