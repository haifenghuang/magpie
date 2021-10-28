package eval

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"magpie/ast"
	"magpie/lexer"
	"magpie/message"
	"magpie/parser"
	"magpie/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	LineStep = 5
)

const (
	ADD_BP = iota
	DEL_BP
)

type DbgInfo struct {
	filename string
	line     int  //node's begin line
	entered  bool //true if the same line has been evaluated.
}

type Debugger struct {
	SrcLines      []string
	SrcLinesCache map[string][]string
	DbgInfos      []*DbgInfo

	Functions map[string]*ast.FunctionLiteral

	//for breakpoint
	Breakpoints map[string]bool //key: 'filename:line'

	Node  ast.Node
	Scope *Scope

	Stepping bool

	prevCommand string
	showPrompt  bool
	listLine    int
}

func NewDebugger() *Debugger {
	d := &Debugger{}
	d.SrcLinesCache = make(map[string][]string)
	d.Breakpoints = make(map[string]bool)
	d.showPrompt = true
	d.Stepping = true
	d.prevCommand = ""

	return d
}

// Add a breakpoint at source line
func (d *Debugger) AddBP(filename string, line int) {
	key := fmt.Sprintf("%s:%d", strings.TrimSpace(filename), line)
	d.Breakpoints[key] = true
}

// Delete a breakpoint at source line
func (d *Debugger) DelBP(filename string, line int) {
	key := fmt.Sprintf("%s:%d", strings.TrimSpace(filename), line)
	if _, ok := d.Breakpoints[key]; ok {
		delete(d.Breakpoints, key)
	}
}

// Check if a source line is at a breakpoint
func (d *Debugger) IsBP(filename string, line int) bool {
	key := fmt.Sprintf("%s:%d", strings.TrimSpace(filename), line)
	for k := range d.Breakpoints {
		if k == key {
			return true
		}
	}
	return false
}

func (d *Debugger) SetNodeAndScope(node ast.Node, scope *Scope) {
	d.Node = node
	d.Scope = scope
}

func (d *Debugger) SetDbgInfos(dbgInfos [][]ast.Node) {
	for _, inf := range dbgInfos {
		d.DbgInfos = append(d.DbgInfos, &DbgInfo{filename: inf[0].Pos().Filename, line: inf[0].Pos().Line, entered: false})
	}
}

func (d *Debugger) SetFunctions(functions map[string]*ast.FunctionLiteral) {
	d.Functions = functions
}

func (d *Debugger) ShowBanner() {
	fmt.Println("                                    _     ")
	fmt.Println("   ____ ___   ____ _ ____ _ ____   (_)___ ")
	fmt.Println("  / __ `__ \\ / __ `// __ `// __ \\ / // _ \\")
	fmt.Println(" / / / / / // /_/ // /_/ // /_/ // //  __/")
	fmt.Println("/_/ /_/ /_/ \\__,_/ \\__, // .___//_/ \\___/ ")
	fmt.Println("                  /____//_/             ")
	fmt.Println("")
}

func (d *Debugger) ProcessCommand() {
	for {
		if !d.showPrompt {
			break
		}

		p := d.Node.Pos()

		/* check if same line has been executed, if so, we need not to show the same line more than once. e.g.
			  println(len("Program end."))
		   Above line have two CallExpressions(println & len),
		   so when we press next, it will show the same line again. we want to avoid this
		*/
		// entered := false
		// for _, inf := range d.DbgInfos {
		// 	if p.Filename == inf.filename && p.Line == inf.line && inf.entered {
		// 		entered = true
		// 		break
		// 	}
		// }
		// if entered {
		// 	break
		// }

		contents, ok := d.SrcLinesCache[p.Filename]
		if ok {
			d.SrcLines = contents
		} else {
			content, _ := ioutil.ReadFile(p.Filename)
			lines := strings.Split(string(content), "\n")
			//pre-append an empty line, so the Lines start with 1, not zero.
			lines = append([]string{""}, lines...)
			d.SrcLinesCache[p.Filename] = lines
			d.SrcLines = lines
		}

		for _, inf := range d.DbgInfos {
			if p.Filename == inf.filename && p.Line == inf.line {
				inf.entered = true
				break
			}
		}

		fmt.Printf("\n%d\t\t%s", p.Line, d.SrcLines[p.Line])
		fmt.Print("\n(magpie) ")

		fmt.Print("\x1b[1m\x1b[36m")

		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		if command == "" && d.prevCommand != "" {
			command = d.prevCommand
		}

		fmt.Print("\x1b[0m")

		d.Stepping = false
		if strings.Compare("c", command) == 0 || strings.Compare("continue", command) == 0 {
			d.prevCommand = command
			break
		} else if strings.Compare("n", command) == 0 || strings.Compare("next", command) == 0 {
			d.prevCommand = command
			d.Stepping = true
			break
		} else if strings.HasPrefix(command, "b ") || strings.HasPrefix(command, "bp ") {
			d.prevCommand = command
			d.processBreakPointCmd(command, ADD_BP)
		} else if strings.HasPrefix(command, "d ") || strings.HasPrefix(command, "del ") {
			d.prevCommand = command
			d.processBreakPointCmd(command, DEL_BP)
		} else if strings.HasPrefix(command, "p ") || strings.HasPrefix(command, "print ") ||
			strings.HasPrefix(command, "e ") || strings.HasPrefix(command, "eval ") {
			d.prevCommand = command
			exp := strings.Split(command, " ")[1:]
			lex := lexer.New("", strings.Join(exp, ""))
			wd, _ := os.Getwd()
			p := parser.New(lex, wd)
			oldLines := d.SrcLines
			oldNode := d.Node
			d.showPrompt = false
			program := p.ParseProgram()
			aval := Eval(program, NewScope(d.Scope, nil))
			fmt.Printf("%s\n\n", aval.Inspect())
			d.SrcLines = oldLines
			d.Node = oldNode
			d.showPrompt = true
		} else if strings.Compare("exit", command) == 0 || strings.Compare("quit", command) == 0 ||
			strings.Compare("bye", command) == 0 || strings.Compare("q", command) == 0 {
			os.Exit(0)
		} else if strings.Compare("l", command) == 0 || strings.Compare("list", command) == 0 {
			if d.listLine == 0 || d.prevCommand != command {
				d.listLine = p.Line
			}

			if d.listLine < len(d.SrcLines) {
				for i := d.listLine; i <= d.listLine+LineStep; i++ {
					if i >= len(d.SrcLines) {
						break
					}
					fmt.Printf("\n%d\t\t%s", i, d.SrcLines[i])
				}
				fmt.Println()
			}

			d.listLine = d.listLine + LineStep + 1
			if d.listLine >= len(d.SrcLines) {
				d.listLine = 0
			}
			d.prevCommand = command
		} else {
			fmt.Printf("Undefined command: '%s'.  Try 'help'.\n", command)
		}
	} //end for
}

//Check if node can be stopped, some nodes cannot be stopped,
//e.g. 'InfixExpression', 'IntegerLiteral'
func (d *Debugger) CanStop() bool {
	flag := false
	switch n := d.Node.(type) {
	case *ast.LetStatement:
		if !n.InClass {
			flag = true
		}
	case *ast.ConstStatement:
		flag = true
	case *ast.ReturnStatement:
		/* we want to stop the 'return' statement to stop only once. for example:
		    fn xxx(x, y) {
                return add(x, y)
		    }

		    if the 'return value(s)' has(have) function call(s), we just stop once, not two.
		*/
		hasCallExpression := false
		for _, value := range n.ReturnValues {
			switch value.(type) {
			case *ast.CallExpression:
				hasCallExpression = true
			}
			if hasCallExpression {
				break
			}
		}

		flag = true
		if hasCallExpression {
			flag = false
		}

	case *ast.DeferStmt:
		flag = true
	case *ast.EnumStatement:
		flag = true
	case *ast.IfExpression:
		flag = true
	case *ast.IfMacroStatement:
		flag = true
	case *ast.UnlessExpression:
		flag = true
	case *ast.CaseExpr:
		flag = true
	case *ast.DoLoop:
		flag = true
	case *ast.WhileLoop:
		flag = true
	case *ast.ForLoop:
		flag = true
	case *ast.ForEverLoop:
		flag = true
	case *ast.ForEachArrayLoop:
		flag = true
	case *ast.ForEachDotRange:
		flag = true
	case *ast.ForEachMapLoop:
		flag = true
	case *ast.BreakExpression:
		flag = true
	case *ast.ContinueExpression:
		flag = true
	case *ast.AssignExpression:
		//  if the assignment expression's value is a 'CallExpression',
		//    we only want to stop once.
		flag = true
		switch n.Value.(type) {
		case *ast.CallExpression:
			flag = false
		}
	case *ast.CallExpression:
		flag = true
	case *ast.TryStmt:
		flag = true
	case *ast.SpawnStmt:
		flag = true
	case *ast.UsingStmt:
		flag = true
	case *ast.QueryExpr:
		flag = true
	case *ast.ServiceStatement:
		flag = true
	default:
		flag = false
	}

	return flag
}

func (d *Debugger) MessageReceived(msg message.Message) {
	ctx := msg.Body.(Context)

	msgType := msg.Type
	switch msgType {
	case message.EVAL_LINE:
		line := ctx.N[0].Pos().Line
		filename := filepath.Base(ctx.N[0].Pos().Filename)
		if d.Stepping {
			d.ProcessCommand()
		} else if d.IsBP(filename, line) {
			fmt.Printf("\nBreakpoint hit at '%s:%d'\n", filename, line)
			d.ProcessCommand()
		}

	case message.CALL:
		// c := ctx.N[0].(*ast.CallExpression)
		// fn := c.Function.String()
		// for funcName, f := range d.Functions {
		// 	if fn == funcName {
		// 		fmt.Printf("\nEnter function '%s' at line %d\n", fn, f.StmtPos().Line)
		// 		break
		// 	}
		// }
	case message.METHOD_CALL:
		// mc := ctx.N[0].(*ast.MethodCallExpression)
		// obj := mc.Object.String()
		// if call, ok := mc.Call.(*ast.CallExpression); ok {
		// 	fn := call.Function.String()
		// 	for funcName, f := range d.Functions {
		// 		if fn == funcName {
		// 			fmt.Printf("\nEnter function '%s.%s' at line %d\n", obj, fn, f.StmtPos().Line)
		// 			break
		// 		}
		// 	}
		// }

	case message.RETURN:
		// r := ctx.N[0].(*ast.ReturnStatement)
		// line := r.Pos().Line
		// for funcName, f := range d.Functions {
		// 	if line >= f.Pos().Line && line <= f.End().Line {
		// 		fmt.Printf("Function '%s' returns\n\n", funcName)
		// 		break
		// 	}
		// }
	}
}

func (d *Debugger) processBreakPointCmd(command string, add_or_del int) {
	p := d.Node.Pos()

	arr := strings.Split(command, " ")
	if len(arr) < 2 {
		fmt.Println("Line number expected.")
	} else {
		//get filename & line/function separator
		filename, breakTxt := getCommandTxt(arr[1:], p)

		line, err := strconv.Atoi(breakTxt)
		if err == nil {
			if line <= 0 {
				fmt.Println("Line number must greater than zero.")
			} else {
				if add_or_del == ADD_BP {
					d.AddBP(filename, line)
				} else {
					d.DelBP(filename, line)
				}
			}
		} else {
			funcName := breakTxt
			var f *ast.FunctionLiteral
			var ok bool
			if f, ok = d.Functions[funcName]; !ok {
				fmt.Println("Function name not found.")
			} else {
				baseName := filepath.Base(f.Pos().Filename)
				if baseName == filename {
					if add_or_del == ADD_BP {
						d.AddBP(filename, line)
					} else {
						d.DelBP(filename, f.StmtPos().Line)
					}
				} else {
					fmt.Println("Function name not found.")
				}
			}
		}
	}
}

//returns 'filename, line/func'
func getCommandTxt(command []string, pos token.Position) (string, string) {
	breakInfTxt := strings.Join(command, " ")
	breakInfTxt = strings.ReplaceAll(breakInfTxt, " ", "") //remove all spaces

	//get filename & line/function separator
	var filename string
	var breakTxt string
	idx := strings.Index(breakInfTxt, ":")
	if idx == -1 {
		filename = filepath.Base(pos.Filename)
		breakTxt = breakInfTxt
	} else {
		filename = breakInfTxt[:idx]
		breakTxt = breakInfTxt[idx+1:]
	}

	return filename, breakTxt
}
