package repl

import (
	"bytes"
	"fmt"
	"io"
	"magpie/ast"
	"magpie/eval"
	"magpie/lexer"
	"magpie/parser"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
)

var magpieKeywords = []string{
	"fn", "let", "true", "false", "if", "else",
	"elif", "return", "import", "and", "or", "struct", "do", "while",
	"break", "continue", "for", "in", "where", "grep", "map", "case",
	"is", "try", "catch", "finally", "throw", "qw", "unless", "spawn",
	"enum", "defer", "nil", "class", "new", "this", "parent", "property",
	"get", "set", "static", "public", "private", "protected", "interface", "default",
	"from", "select", "group", "into", "orderby", "join", "on", "equals", "by", "ascending", "descending",
	"async", "await", "service",
}

//Note: we should put the longest operators first.
var magpieOperators = []string{
	"+=", "-=", "*=", "/=", "%=", "^=",
	"++", "--",
	"&&", "||",
	"<<", ">>",
	"->", "=>",
	"==", "!=", "<=", ">=", "=~", "!~",
	"+", "-", "*", "/", "%", "^",
	"(", ")", "{", "}", "[", "]",
	"=", "<", ">",
	"!", "&", "|", ".",
	",", "?", ":", ";",
}

var colors = map[liner.Category]liner.Color{
	liner.NumberType:   liner.COLOR_YELLOW,
	liner.KeywordType:  liner.COLOR_MAGENTA,
	liner.StringType:   liner.COLOR_CYAN,
	liner.CommentType:  liner.COLOR_GREEN,
	liner.OperatorType: liner.COLOR_RED,
	liner.IdentType:    liner.COLOR_WHITE,
}

const PROMPT = "magpie>> "
const CONT_PROMPT = "... " // continue prompt

func Start(out io.Writer, color bool) {
	history := filepath.Join(os.TempDir(), ".magpie_history")
	l := liner.NewLiner()
	defer l.Close()

	l.SetCtrlCAborts(true)
	l.SetMultiLineMode(true)

	if color {
		l.SetSyntaxHighlight(color) //use syntax highlight or not
		l.RegisterKeywords(magpieKeywords)
		l.RegisterOperators(magpieOperators)
		l.RegisterColors(colors)
	}

	if f, err := os.Open(history); err == nil {
		l.ReadHistory(f)
		f.Close()
	}

	if color && !l.IsInwinConsole() {
		eval.REPLColor = true
	}
	scope := eval.NewScope(nil, os.Stdout)
	wd, err := os.Getwd()
	if err != nil {
		io.WriteString(out, err.Error())
		os.Exit(1)
	}

	// var tmplines []string
	var lex *lexer.Lexer
	var p *parser.Parser
	var program *ast.Program
	for {
		if line, err := l.Prompt(PROMPT); err == nil {
			if line == "exit" || line == "quit" {
				if f, err := os.Create(history); err == nil {
					l.WriteHistory(f)
					f.Close()
				}
				break
			}

			tmpline := strings.TrimSpace(line)
			if len(tmpline) == 0 || tmpline[0] == '#' { //empty line or single comment line
				continue
			} else {
				//check if the line is a valid expression or statement
				lex = lexer.New("", tmpline)
				p = parser.New(lex, wd)
				program = p.ParseProgram()
				if len(p.Errors()) == 0 { // no error
					eval.Eval(program, scope)
					l.AppendHistory(tmpline)
					continue
				} else {
					if program == nil {
						printParserErrors(out, p.Errors())
						continue
					} else if len(program.Statements) == 0 { //it's an 'import' statement
						var errFlag bool
						for _, importItem := range program.Imports {
							if importItem.Program == nil { // error
								errFlag = true
								break
							}
						}

						if errFlag {
							printParserErrors(out, p.Errors())
						}
						continue
					}

					var buf bytes.Buffer
					fmt.Fprintln(&buf, line)
					for {
						if line, err := l.Prompt(CONT_PROMPT); err == nil {
							fmt.Fprintln(&buf, line)

							text := string(buf.Bytes())
							lex = lexer.New("", text)
							p = parser.New(lex, wd)
							program = p.ParseProgram()
							if len(p.Errors()) == 0 { // no error
								eval.Eval(program, scope)
								l.AppendHistory(strings.Replace(text, "\n", "", -1))
								break
							} else {
								continue
							}
						} else if err == liner.ErrPromptAborted { //CTRL-C pressed
							break
						}
					}
				}
			}
		} else if err == liner.ErrPromptAborted { //CTRL-C pressed
			if f, err := os.Create(history); err == nil {
				l.WriteHistory(f)
				f.Close()
			}
			break
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
