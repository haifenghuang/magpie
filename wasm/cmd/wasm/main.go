package main

import (
	"bytes"
	_ "fmt"
	"syscall/js"

	"magpie/eval"
	"magpie/lexer"
	"magpie/parser"
)

func runCode(this js.Value, i []js.Value) interface{} {
	m := make(map[string]interface{})
	var buf bytes.Buffer

	m["errline"] = -1
	l := lexer.New("", i[0].String())
	p := parser.New(l, "")
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		buf.WriteString("parser errors:\n")
		for _, msg := range p.Errors() {
			buf.WriteString("\t" + msg + "\n")
		}
		errLines := p.ErrorLines()

		m["errline"] = errLines[0]
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
