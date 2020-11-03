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
	var buf bytes.Buffer

	l := lexer.New("", i[0].String())
	p := parser.New(l, "")
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		buf.WriteString("parser errors:\n")
		for _, msg := range p.Errors() {
			buf.WriteString("\t" + msg + "\n")
		}
		return buf.String()
	}

	scope := eval.NewScope(nil, &buf)
	result := eval.Eval(program, scope)
	if result.Type() == eval.ERROR_OBJ {
		return result.Inspect()
	}

	return buf.String()
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("magpie_run_code", js.FuncOf(runCode))
	<-c
}
