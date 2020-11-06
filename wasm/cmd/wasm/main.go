package main

import (
	"bytes"
	_ "fmt"
	"strings"
	"syscall/js"

	"magpie/eval"
	"magpie/lexer"
	"magpie/parser"
)

func runCode(this js.Value, i []js.Value) interface{} {
	m := make(map[string]interface{})
	var buf bytes.Buffer

	m["errlines"] = "-1"
	l := lexer.New("", i[0].String())
	p := parser.New(l, "")
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			buf.WriteString(msg + "\n")
		}

		m["errlines"] = strings.Join(p.ErrorLines(), "|")
		m["output"] = buf.String()
		return m
	}

	scope := eval.NewScope(nil, &buf)
	result := eval.Eval(program, scope)
	if (string(result.Type()) == eval.ERROR_OBJ) {
		m["output"] = buf.String() + result.Inspect()
	} else {
		m["output"] = buf.String()
	}

	return m
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("magpie_run_code", js.FuncOf(runCode))
	<-c
}
