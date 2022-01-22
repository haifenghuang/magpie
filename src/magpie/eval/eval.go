package eval

import (
	"bytes"
	"fmt"
	"magpie/ast"
	_ "magpie/lexer"
	"magpie/message"
	"magpie/token"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	TRUE     = &Boolean{Bool: true, Valid: true}
	FALSE    = &Boolean{Bool: false, Valid: true}
	BREAK    = &Break{}
	CONTINUE = &Continue{}
	NIL      = &Nil{}
	EMPTY    = &Optional{Value: NIL}

	// Built-in types which you can extend.
	builtinTypes = []string{
		"integer",  // signed integer
		"uinteger", // unsigned integer
		"float",
		"boolean",
		"string",
		"array",
		"tuple",
		"hash"}
)

var importScope *Scope
var importedCache map[string]Object

var mux sync.Mutex

//REPL with color support
var REPLColor bool

const ServiceHint = "* Running on %s (Press CTRL+C to quit)\n"

var Dbg *Debugger
var MsgHandler *message.MessageHandler

type Context struct {
	N []ast.Node //N: node
	S *Scope     //S: Scope
}

func PanicToError(p interface{}, node ast.Node) error {
	switch e := p.(type) {
	case *Error: //Error Object defined in errors.go file
		return fmt.Errorf("%s - Line:%s", e.Inspect(), node.Pos().Sline())
	case error:
		return fmt.Errorf("%s - Line:%s", e, node.Pos().Sline())
	case string:
		return fmt.Errorf("%s - Line:%s", e, node.Pos().Sline())
	case fmt.Stringer:
		return fmt.Errorf("%s - Line:%s", e.String(), node.Pos().Sline())
	default:
		return fmt.Errorf("unknown error type (%T) - Line:%s", e, node.Pos().Sline())
	}
}

func Eval(node ast.Node, scope *Scope) (val Object) {
	defer func() {
		if r := recover(); r != nil {
			err := PanicToError(r, node)
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			//WHY return NIL? if we do not return 'NIL', we may get something like below:
			//    PANIC=runtime error: invalid memory address or nil pointer
			val = NIL
		}
	}()

	if Dbg != nil {
		Dbg.SetNodeAndScope(node, scope)
		if Dbg.CanStop() {
			MsgHandler.SendMessage(message.Message{Type: message.EVAL_LINE, Body: Context{N: []ast.Node{node}, S: scope}})
		}
	}

	//fmt.Printf("node.Type=%T, node=<%s>, start=%d, end=%d\n", node, node.String(), node.Pos().Line, node.End().Line) //debugging
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, scope)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, scope)
	case *ast.ImportStatement:
		return evalImportStatement(node, scope)
	case *ast.LetStatement:
		return evalLetStatement(node, scope)
	case *ast.ConstStatement:
		return evalConstStatement(node, scope)
	case *ast.ReturnStatement:
		if Dbg != nil {
			MsgHandler.SendMessage(message.Message{Type: message.RETURN, Body: Context{N: []ast.Node{node}, S: scope}})
		}
		return evalReturnStatement(node, scope)
	case *ast.DeferStmt:
		return evalDeferStatement(node, scope)
	case *ast.FunctionStatement:
		return evalFunctionStatement(node, scope)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	case *ast.UIntegerLiteral:
		return evalUIntegerLiteral(node)
	case *ast.FloatLiteral:
		return evalFloatLiteral(node)
	case *ast.StringLiteral:
		return evalStringLiteral(node)
	case *ast.InterpolatedString:
		return evalInterpolatedString(node, scope)
	case *ast.Identifier:
		return evalIdentifier(node, scope)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, scope)
	case *ast.TupleLiteral:
		return evalTupleLiteral(node, scope)
	case *ast.HashLiteral:
		return evalHashLiteral(node, scope)
	case *ast.StructLiteral:
		return evalStructLiteral(node, scope)
	case *ast.EnumLiteral:
		return evalEnumLiteral(node, scope)
	case *ast.RangeLiteral:
		return evalRangeLiteral(node, scope)
	case *ast.EnumStatement:
		return evalEnumStatement(node, scope)
	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, scope)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, scope)
	case *ast.InfixExpression:
		left := Eval(node.Left, scope)
		if left.Type() == ERROR_OBJ {
			return left
		}

		right := Eval(node.Right, scope)
		if right.Type() == ERROR_OBJ {
			return right
		}
		return evalInfixExpression(node, left, right, scope)
	case *ast.PostfixExpression:
		left := Eval(node.Left, scope)
		if left.Type() == ERROR_OBJ {
			return left
		}
		return evalPostfixExpression(left, node, scope)
	case *ast.IfExpression:
		return evalIfExpression(node, scope)
	case *ast.IfMacroStatement:
		return evalIfMacroStatement(node, scope)
	case *ast.UnlessExpression:
		return evalUnlessExpression(node, scope)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, scope)
	case *ast.CallExpression:
		// if Dbg != nil {
		// 	MsgHandler.SendMessage(message.Message{Type: message.CALL, Body: Context{N: []ast.Node{node}, S: scope}})
		// }
		return evalFunctionCall(node, scope)
	case *ast.MethodCallExpression:
		if Dbg != nil {
			MsgHandler.SendMessage(message.Message{Type: message.METHOD_CALL, Body: Context{N: []ast.Node{node}, S: scope}})
		}
		return evalMethodCallExpression(node, scope)
	case *ast.IndexExpression:
		return evalIndexExpression(node, scope)
	case *ast.GrepExpr:
		return evalGrepExpression(node, scope)
	case *ast.MapExpr:
		return evalMapExpression(node, scope)
	case *ast.CaseExpr:
		return evalCaseExpression(node, scope)
	case *ast.DoLoop:
		return evalDoLoopExpression(node, scope)
	case *ast.WhileLoop:
		return evalWhileLoopExpression(node, scope)
	case *ast.ForLoop:
		return evalForLoopExpression(node, scope)
	case *ast.ForEverLoop:
		return evalForEverLoopExpression(node, scope)
	case *ast.ForEachArrayLoop:
		return evalForEachArrayExpression(node, scope)
	case *ast.ForEachDotRange:
		return evalForEachDotRangeExpression(node, scope)
	case *ast.ForEachMapLoop:
		return evalForEachMapExpression(node, scope)
	case *ast.ListComprehension:
		return evalListComprehension(node, scope)
	case *ast.ListRangeComprehension:
		return evalListRangeComprehension(node, scope)
	case *ast.ListMapComprehension:
		return evalListMapComprehension(node, scope)
	case *ast.HashComprehension:
		return evalHashComprehension(node, scope)
	case *ast.HashRangeComprehension:
		return evalHashRangeComprehension(node, scope)
	case *ast.HashMapComprehension:
		return evalHashMapComprehension(node, scope)
	case *ast.BreakExpression:
		return BREAK
	case *ast.ContinueExpression:
		return CONTINUE
	case *ast.ThrowStmt:
		return evalThrowStatement(node, scope)
	case *ast.AssignExpression:
		return evalAssignExpression(node, scope)
	case *ast.RegExLiteral:
		return evalRegExLiteral(node)
	case *ast.TryStmt:
		return evalTryStatement(node, scope)
	case *ast.TernaryExpression:
		return evalTernaryExpression(node, scope)
	case *ast.SpawnStmt:
		return evalSpawnStatement(node, scope)
	case *ast.NilLiteral:
		return NIL
	case *ast.Pipe:
		return evalPipeExpression(node, scope)

	//Class related
	case *ast.ClassStatement:
		return evalClassStatement(node, scope)
	case *ast.ClassLiteral:
		return evalClassLiteral(node, scope)
	case *ast.NewExpression:
		return evalNewExpression(node, scope)
	//using
	case *ast.UsingStmt:
		return evalUsingStatement(node, scope)

	//command expression
	case *ast.CmdExpression:
		return evalCmdExpression(node, scope)

	//linq query expression
	case *ast.QueryExpr:
		return evalLinqQueryExpression(node, scope)

	//await expression
	case *ast.AwaitExpr:
		return evalAwaitExpression(node, scope)

	//service statement
	case *ast.ServiceStatement:
		return evalServiceStatement(node, scope)

	//date time object
	case *ast.DateTimeExpr:
		if node.Pattern == nil {
			return &TimeObj{Tm: time.Now(), Valid: true}
		}

		is := evalInterpolatedString(node.Pattern, scope).(*InterpolatedString)
		var err error
		dt := &TimeObj{Valid: true}
		dt.Tm, err = time.Parse(builtinDate_Normal, is.String.String)
		if err != nil {
			dt.Valid = false
		}
		return dt
	//diamond: <$fobj>
	case *ast.DiamondExpr:
		return evalDiamondExpr(node, scope)
	}

	return nil
}

// Program Evaluation Entry Point Functions, and Helpers:
func evalProgram(program *ast.Program, scope *Scope) (results Object) {
	if importedCache == nil {
		importedCache = make(map[string]Object)
	}

	results = loadImports(program.Imports, scope)
	if results.Type() == ERROR_OBJ {
		return
	}

	for _, statement := range program.Statements {
		results = Eval(statement, scope)
		switch s := results.(type) {
		case *ReturnValue:
			return s.Value
		case *Error:
			return s
		case *Throw:
			//convert ThrowValue to Errors
			return NewError(s.stmt.Pos().Sline(), THROWNOTHANDLED, s.value.Inspect())
		}
	}
	if results == nil {
		return NIL
	}
	return results
}

func loadImports(imports map[string]*ast.ImportStatement, scope *Scope) Object {
	if importScope == nil {
		importScope = NewScope(nil, scope.Writer)
	}
	for _, p := range imports {
		v := Eval(p, scope)
		if v.Type() == ERROR_OBJ {
			return NewError(p.Pos().Sline(), IMPORTERROR, p.ImportPath)
		}
	}
	return NIL
}

// Statements...
func evalImportStatement(i *ast.ImportStatement, scope *Scope) Object {

	mux.Lock()
	defer mux.Unlock()

	// Check the cache
	if cache, ok := importedCache[i.ImportPath]; ok {
		return cache
	}

	imported := &ImportedObject{Name: i.ImportPath, Scope: NewScope(nil, scope.Writer)}

	if _, ok := importScope.Get(i.ImportPath); !ok {
		evalProgram(i.Program, imported.Scope)
		importScope.Set(i.ImportPath, imported)
	}

	//store the evaluated result to cache
	importedCache[i.ImportPath] = imported

	return imported
}

func evalLetStatement(l *ast.LetStatement, scope *Scope) (val Object) {
	if l.DestructingFlag {
		v := Eval(l.Values[0], scope)
		valType := v.Type()
		switch valType {
		case HASH_OBJ:
			h := v.(*Hash)
			for _, item := range l.Names {
				if item.Token.Type == token.UNDERSCORE {
					continue
				}
				found := false
				for _, pair := range h.Pairs {
					if item.String() == pair.Key.Inspect() {
						val = pair.Value
						scope.Set(item.String(), pair.Value)
						found = true
					}
				}
				if !found {
					val = NIL
					scope.Set(item.String(), val)
				}
			}

		case ARRAY_OBJ:
			arr := v.(*Array)
			valuesLen := len(arr.Members)
			for idx, item := range l.Names {
				if idx >= valuesLen { //There are more Names than Values
					if item.Token.Type != token.UNDERSCORE {
						val = NIL
						scope.Set(item.String(), val)
					}
				} else {
					if item.Token.Type == token.UNDERSCORE {
						continue
					}
					val = arr.Members[idx]
					if val.Type() != ERROR_OBJ {
						scope.Set(item.String(), val)
					} else {
						return
					}
				}
			}

		case TUPLE_OBJ:
			tup := v.(*Tuple)
			valuesLen := len(tup.Members)
			for idx, item := range l.Names {
				if idx >= valuesLen { //There are more Names than Values
					if item.Token.Type != token.UNDERSCORE {
						val = NIL
						scope.Set(item.String(), val)
					}
					val = NIL
					scope.Set(item.String(), val)
				} else {
					if item.Token.Type == token.UNDERSCORE {
						continue
					}
					val = tup.Members[idx]
					if val.Type() != ERROR_OBJ {
						scope.Set(item.String(), val)
					} else {
						return
					}
				}
			}

		default:
			return NewError(l.Pos().Sline(), GENERICERROR, "Only Array|Tuple|Hash is allowed!")
		}

		return
	}

	values := []Object{}
	valuesLen := 0
	for _, value := range l.Values {
		val := Eval(value, scope)
		if val.Type() == TUPLE_OBJ {
			tupleObj := val.(*Tuple)
			if tupleObj.IsMulti {
				valuesLen += len(tupleObj.Members)
				for _, tupleItem := range tupleObj.Members {
					values = append(values, tupleItem)
				}
			} else {
				valuesLen += 1
				values = append(values, tupleObj)
			}

		} else {
			valuesLen += 1
			values = append(values, val)
		}
	}

	for idx, item := range l.Names {
		if idx >= valuesLen { //There are more Names than Values
			if item.Token.Type != token.UNDERSCORE {
				val = NIL
				scope.Set(item.String(), val)
			}
		} else {
			if item.Token.Type == token.UNDERSCORE {
				continue
			}
			val = values[idx]
			if val.Type() != ERROR_OBJ {
				scope.Set(item.String(), val)
			} else {
				return
			}
		}
	}

	return
}

func evalConstStatement(c *ast.ConstStatement, scope *Scope) Object {
	for idx, name := range c.Name {
		val := Eval(c.Value[idx], scope)
		if val.Type() == ERROR_OBJ {
			return val
		}
		scope.SetConst(name.Value, val)
	}
	return NIL
}

func evalNumAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	var leftVal float64
	var rightVal float64

	isInt := left.Type() == INTEGER_OBJ && val.Type() == INTEGER_OBJ
	isUInt := left.Type() == UINTEGER_OBJ && val.Type() == UINTEGER_OBJ

	if left.Type() == INTEGER_OBJ {
		leftVal = float64(left.(*Integer).Int64)
	} else if left.Type() == UINTEGER_OBJ {
		leftVal = float64(left.(*UInteger).UInt64)
	} else {
		leftVal = left.(*Float).Float64
	}

	//Check `right`'s type
	if val.Type() == INTEGER_OBJ {
		rightVal = float64(val.(*Integer).Int64)
	} else if val.Type() == UINTEGER_OBJ {
		rightVal = float64(val.(*UInteger).UInt64)
	} else {
		rightVal = val.(*Float).Float64
	}

	var ok bool
	switch a.Token.Literal {
	case "+=":
		result := leftVal + rightVal
		if isInt { //only 'INTEGER + INTEGER'
			ret, ok = scope.Reset(name, NewInteger(int64(result)))
			if ok {
				return
			}
		} else if isUInt { //only 'UINTEGER + UINTEGER'
			ret, ok = scope.Reset(name, NewUInteger(uint64(result)))
			if ok {
				return
			}
		} else {
			ret, ok = checkNumAssign(scope, name, left, val, result)
			if ok {
				return
			}
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())

	case "-=":
		result := leftVal - rightVal
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(result)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(result)))
			if ok {
				return
			}
		} else {
			ret, ok = checkNumAssign(scope, name, left, val, result)
			if ok {
				return
			}
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())

	case "*=":
		result := leftVal * rightVal
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(result)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(result)))
			if ok {
				return
			}
		} else {
			ret, ok = checkNumAssign(scope, name, left, val, result)
			if ok {
				return
			}
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())

	case "/=":
		if rightVal == 0 {
			return NewError(a.Pos().Sline(), DIVIDEBYZERO)
		}

		result := leftVal / rightVal
		//Always return Float
		ret, ok = scope.Reset(name, NewFloat(result))
		if ok {
			return
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())

	case "%=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)%int64(rightVal)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(leftVal)%uint64(rightVal)))
			if ok {
				return
			}
		} else {
			result := math.Mod(leftVal, rightVal)
			ret, ok = checkNumAssign(scope, name, left, val, result)
			if ok {
				return
			}
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())

	case "^=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)^int64(rightVal)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(leftVal)^uint64(rightVal)))
			if ok {
				return
			}
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())

	case "|=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)|int64(rightVal)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(leftVal)|uint64(rightVal)))
			if ok {
				return
			}
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())

	case "&=":
		if isInt {
			ret, ok = scope.Reset(name, NewInteger(int64(leftVal)&int64(rightVal)))
			if ok {
				return
			}
		} else if isUInt {
			ret, ok = scope.Reset(name, NewUInteger(uint64(leftVal)&uint64(rightVal)))
			if ok {
				return
			}
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
	}
	return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
}

func checkNumAssign(scope *Scope, name string, left Object, right Object, result float64) (ret Object, ok bool) {
	if left.Type() == FLOAT_OBJ || right.Type() == FLOAT_OBJ {
		ret, ok = scope.Reset(name, NewFloat(result))
		return
	}

	if (left.Type() == INTEGER_OBJ && right.Type() == UINTEGER_OBJ) ||
		(left.Type() == UINTEGER_OBJ && right.Type() == INTEGER_OBJ) {
		if result > math.MaxInt64 {
			ret, ok = scope.Reset(name, NewUInteger(uint64(result)))
		} else {
			ret, ok = scope.Reset(name, NewInteger(int64(result)))
		}
	}
	return
}

//str[idx] = item
//str += item
func evalStrAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	leftVal := left.(*String).String
	var ok bool

	switch a.Token.Literal {
	case "=":
		switch nodeType := a.Name.(type) {
		case *ast.IndexExpression: //str[idx] = xxx
			index := Eval(nodeType.Index, scope)
			if index == NIL {
				ret = NIL
				return
			}

			var idx int64
			switch o := index.(type) {
			case *Integer:
				idx = o.Int64
			case *UInteger:
				idx = int64(o.UInt64)
			}

			if idx < 0 || idx >= int64(len(leftVal)) {
				return NewError(a.Pos().Sline(), INDEXERROR, idx)
			}

			str := NewString(leftVal[:idx] + val.Inspect() + leftVal[idx+1:])
			ret, ok = scope.Reset(name, str)
			if ok {
				return
			}
		}
	}

	ret, ok = scope.Reset(name, NewString(leftVal+val.Inspect()))
	if ok {
		return
	}
	return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())

}

//array[idx] = item
//array += item
func evalArrayAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	leftVals := left.(*Array).Members

	var ok bool
	switch a.Token.Literal {
	case "+=":
		switch nodeType := a.Name.(type) {
		case *ast.Identifier:
			name = nodeType.Value
			leftVals = append(leftVals, val)
			ret, ok = scope.Reset(name, &Array{Members: leftVals})
			if ok {
				return
			}
		}
	case "=":
		switch nodeType := a.Name.(type) {
		case *ast.IndexExpression: //arr[idx] = xxx
			index := Eval(nodeType.Index, scope)
			if index == NIL {
				ret = NIL
				return
			}

			var idx int64
			switch o := index.(type) {
			case *Integer:
				idx = o.Int64
			case *UInteger:
				idx = int64(o.UInt64)
			}
			if idx < 0 {
				return NewError(a.Pos().Sline(), INDEXERROR, idx)
			}

			if idx < int64(len(leftVals)) { //index is in range
				leftVals[idx] = val
				ret, ok = scope.Reset(name, &Array{Members: leftVals})
				if ok {
					return
				}
			} else { //index is out of range, we auto-expand the array
				for i := int64(len(leftVals)); i < idx; i++ {
					leftVals = append(leftVals, NIL)
				}

				leftVals = append(leftVals, val)
				ret, ok = scope.Reset(name, &Array{Members: leftVals})
				if ok {
					return
				}
			}
		}

		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
	}

	return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
}

func evalTupleAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	//Tuple is an immutable sequence of values
	if a.Token.Literal == "=" { //tuple[idx] = item
		str := fmt.Sprintf("%s[IDX]", TUPLE_OBJ)
		return NewError(a.Pos().Sline(), INFIXOP, str, a.Token.Literal, val.Type())
	}
	return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
}

func evalHashAssignExpression(a *ast.AssignExpression, name string, left Object, scope *Scope, val Object) (ret Object) {
	leftHash := left.(*Hash)

	var ok bool
	switch a.Token.Literal {
	case "+=":
		if _, ok := val.(*Hash); !ok { //must be hash type
			return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
		}

		rightHash := val.(*Hash)
		for _, hk := range rightHash.Order { //hk:hash key
			pair, _ := rightHash.Pairs[hk]
			leftHash.Push(a.Pos().Sline(), pair.Key, pair.Value)
		}
		ret, ok = scope.Reset(name, leftHash)
		if ok {
			return
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
	case "-=":
		_, ok := val.(Hashable)
		if !ok {
			return NewError(a.Pos().Sline(), KEYERROR, val.Type())
		}
		leftHash.Pop(a.Pos().Sline(), val)
		ret, ok = scope.Reset(name, leftHash)
		if ok {
			return
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
	case "=":
		switch nodeType := a.Name.(type) {
		case *ast.IndexExpression: //hashObj[key] = val
			key := Eval(nodeType.Index, scope)
			leftHash.Push(a.Pos().Sline(), key, val)
			return leftHash
		case *ast.Identifier: //hashObj.key = val
			key := strings.Split(a.Name.String(), ".")[1]
			keyObj := NewString(key)
			leftHash.Push(a.Pos().Sline(), keyObj, val)
			return leftHash
		}
		return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
	}

	return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
}

func evalStructAssignExpression(a *ast.AssignExpression, scope *Scope, val Object) (retVal Object) {
	strArr := strings.Split(a.Name.String(), ".")
	var aObj Object
	var aVal Object
	var ok bool
	if aObj, ok = scope.Get(strArr[0]); !ok {
		return NewError(a.Pos().Sline(), UNKNOWNIDENT, strArr[0])
	}

	st, ok := aObj.(*Struct)
	if !ok {
		return NewError(a.Pos().Sline(), UNKNOWNIDENT, strArr[0])
	}

	if aVal, ok = st.Scope.Get(strArr[1]); !ok {
		return NewError(a.Pos().Sline(), UNKNOWNIDENT, strArr[1])
	}

	structScope := st.Scope

	if a.Token.Literal == "=" {
		v, ok := structScope.Reset(strArr[1], retVal)
		if ok {
			return v
		}
		return NewError(a.Pos().Sline(), UNKNOWNIDENT, a.Name.String())
	}

	switch aVal.Type() {
	case INTEGER_OBJ, UINTEGER_OBJ, FLOAT_OBJ:
		retVal = evalNumAssignExpression(a, strArr[1], aVal, structScope, val)
		st.Scope = structScope
		return
	case STRING_OBJ:
		retVal = evalStrAssignExpression(a, strArr[1], aVal, structScope, val)
		st.Scope = structScope
		return
	case ARRAY_OBJ:
		retVal = evalArrayAssignExpression(a, strArr[1], aVal, structScope, val)
		st.Scope = structScope
		return
	case HASH_OBJ:
		retVal = evalHashAssignExpression(a, strArr[1], aVal, structScope, val)
		st.Scope = structScope
		return
	}

	return NewError(a.Pos().Sline(), INFIXOP, aVal.Type(), a.Token.Literal, val.Type())
}

//instanceObj[x] = xxx
//instanceObj[x,y] = xxx
func evalClassIndexerAssignExpression(a *ast.AssignExpression, obj Object, indexExpr *ast.IndexExpression, val Object, scope *Scope) Object {
	instanceObj := obj.(*ObjectInstance)

	var num int
	switch o := indexExpr.Index.(type) {
	case *ast.ClassIndexerExpression:
		num = len(o.Parameters)
	default:
		num = 1
	}

	propName := "this" + fmt.Sprintf("%d", num)

	//check if the Indexer is static
	if instanceObj.IsStatic(propName, ClassPropertyKind) {
		return NewError(a.Pos().Sline(), INDEXERSTATICERROR, instanceObj.Class.Name)
	}

	p := instanceObj.GetProperty(propName)
	if p != nil {
		//no setter or setter block is empty, e.g. 'property xxx { set; }'
		if p.Setter == nil || len(p.Setter.Body.Statements) == 0 {
			return NewError(a.Pos().Sline(), INDEXERUSEERROR, instanceObj.Class.Name)
		} else {
			newScope := NewScope(instanceObj.Scope, nil)
			newScope.Set("value", val)

			switch o := indexExpr.Index.(type) {
			case *ast.ClassIndexerExpression:
				for i, v := range o.Parameters {
					index := Eval(v, scope)
					newScope.Set(p.Indexes[i].Value, index)
				}
			default:
				index := Eval(indexExpr.Index, scope)
				newScope.Set(p.Indexes[0].Value, index)
			}

			results := Eval(p.Setter.Body, newScope)
			if results.Type() == RETURN_VALUE_OBJ {
				return results.(*ReturnValue).Value
			}
			return results
		}
	} else {
		return NewError(a.Pos().Sline(), INDEXNOTFOUNDERROR, instanceObj.Class.Name)
	}
}

func evalAssignExpression(a *ast.AssignExpression, scope *Scope) (val Object) {
	val = Eval(a.Value, scope)
	if val.Type() == ERROR_OBJ {
		return val
	}

	if strings.Contains(a.Name.String(), ".") {
		switch o := a.Name.(type) {
		case *ast.MethodCallExpression:
			obj := Eval(o.Object, scope)
			if obj.Type() == ERROR_OBJ {
				return obj
			}

			switch m := obj.(type) {
			case *Hash:
				switch c := o.Call.(type) {
				case *ast.Identifier:
					//e.g.
					//doc = {"one": {"two": { "three": [1, 2, 3,] }}}
					//doc.one.two.three = 44
					m.Push(a.Pos().Sline(), NewString(c.Value), val)
					return
				case *ast.IndexExpression:
					//e.g.
					//doc = {"one": {"two": { "three": [1, 2, 3,] }}}
					//doc.one.two.three[2] = 44
					leftVal := m.Get(a.Pos().Sline(), NewString(c.Left.String()))
					indexVal := Eval(c.Index, scope)
					switch v := leftVal.(type) {
					case *Hash:
						v.Push(a.Pos().Sline(), indexVal, val)
					case *Array:
						v.Set(a.Pos().Sline(), indexVal, val)
					case *Tuple:
						str := fmt.Sprintf("%s[IDX]", TUPLE_OBJ)
						return NewError(a.Pos().Sline(), INFIXOP, str, "=", val.Type())
					}
					return NIL
				}
			}
		}
		var aObj Object
		var ok bool

		strArr := strings.Split(a.Name.String(), ".")
		if aObj, ok = scope.Get(strArr[0]); !ok {
			return reportTypoSuggestions(a.Pos().Sline(), scope, strArr[0])
			//return NewError(a.Pos().Sline(), UNKNOWNIDENT, strArr[0])3
		}

		if aObj.Type() == ENUM_OBJ { //it's enum type
			return NewError(a.Pos().Sline(), GENERICERROR, "Enum value cannot be reassigned!")
		} else if aObj.Type() == HASH_OBJ { //e.g. hash.key = value
			return evalHashAssignExpression(a, strArr[0], aObj, scope, val)
		} else if aObj.Type() == INSTANCE_OBJ { //e.g. this.var = xxxx
			instanceObj := aObj.(*ObjectInstance)

			//			//get variable's modifier level
			//			ml := instanceObj.GetModifierLevel(strArr[1], ClassMemberKind) //ml:modifier level
			//			if ml == ast.ModifierPrivate {
			//				return NewError(a.Pos().Sline(), CLSMEMBERPRIVATE, strArr[1], instanceObj.Class.Name)
			//			}

			//check if it's a property
			p := instanceObj.GetProperty(strArr[1])
			if p == nil { //not property, return value from scope
				// check if it's a static variable
				if instanceObj.IsStatic(strArr[1], ClassMemberKind) {
					return NewError(a.Pos().Sline(), MEMBERUSEERROR, strArr[1], instanceObj.Class.Name)
				}
				instanceObj.Scope.Set(strArr[1], val)
			} else {
				// check if it's a static property
				if instanceObj.IsStatic(strArr[1], ClassPropertyKind) {
					return NewError(a.Pos().Sline(), PROPERTYUSEERROR, strArr[1], instanceObj.Class.Name)
				}

				if p.Setter == nil { //property xxx { get; }
					_, ok := instanceObj.Scope.Get(strArr[1])
					if !ok { //it's the first time assignment
						instanceObj.Scope.Set(strArr[1], val)
					} else {
						return NewError(a.Pos().Sline(), PROPERTYUSEERROR, strArr[1], instanceObj.Class.Name)
					}
				} else {
					if len(p.Setter.Body.Statements) == 0 { // property xxx { set; }
						instanceObj.Scope.Set("_"+strArr[1], val)
					} else {
						newScope := NewScope(instanceObj.Scope, nil)
						newScope.Set("value", val)
						results := Eval(p.Setter.Body, newScope)
						if results.Type() == RETURN_VALUE_OBJ {
							val = results.(*ReturnValue).Value
						}
					}
				}
			}
			return
		} else if aObj.Type() == CLASS_OBJ { //e.g. parent.var = xxxx
			clsObj := aObj.(*Class)

			//check if it's a property
			p := clsObj.GetProperty(strArr[1])
			if p == nil { //not property
				// check if it's a static member
				if !clsObj.IsStatic(strArr[1], ClassMemberKind) {
					return NewError(a.Pos().Sline(), MEMBERUSEERROR, strArr[1], clsObj.Name)
				}
			} else {
				// check if it's a static property
				if !clsObj.IsStatic(strArr[1], ClassPropertyKind) {
					return NewError(a.Pos().Sline(), PROPERTYUSEERROR, strArr[1], clsObj.Name)
				}
			}

			thisObj, _ := scope.Get("this")
			if thisObj != nil {
				if thisObj.Type() == INSTANCE_OBJ { //'this' refers to 'ObjectInstance' object
					//Check if `thisObj` instance's scope could find `strArr[1]`
					_, ok = thisObj.(*ObjectInstance).Scope.Get(strArr[1])
					if ok {
						thisObj.(*ObjectInstance).Scope.Set(strArr[1], val)
						return
					} else {
						// Why this 'else' branch, please see below example:
						// class Dog {
						//	static let misc = 12
						//
						//	fn MethodA() {
						//		printf("Hello\n")
						//		Dog.misc = 20
						//	}
						// }
						//
						// let dogObj = new Dog("doggy")
						// printf("Dog.misc=%d\n", Dog.misc)
						// dogObj.MethodA();
						// printf("Dog.misc=%d\n", Dog.misc)
						//
						// Here, when we call `dogObj.MethodA`, the 'thisObj' will refer to 'ObjectInstance' object.
						// But when we assign '20' to 'Dog.misc', we should set '20' to clsObj's Scope, not
						// instance object's scope
						clsObj.Scope.Set(strArr[1], val)
						return
					}
				}
			}

			if p == nil { //not property
				clsObj.Scope.Set(strArr[1], val)
			} else {
				if p.Setter == nil { //property xxx { get; }
					_, ok := clsObj.Scope.Get(strArr[1])
					if !ok { //it's the first time assignment
						clsObj.Scope.Set(strArr[1], val)
					} else {
						return NewError(a.Pos().Sline(), PROPERTYUSEERROR, strArr[1], clsObj.Name)
					}
				} else {
					if len(p.Setter.Body.Statements) == 0 { // property xxx { set; }
						clsObj.Scope.Set("_"+strArr[1], val)
					} else {
						newScope := NewScope(clsObj.Scope, nil)
						newScope.Set("value", val)
						results := Eval(p.Setter.Body, newScope)
						if results.Type() == RETURN_VALUE_OBJ {
							val = results.(*ReturnValue).Value
						}
					}
				}
			}
			return
		}

		return evalStructAssignExpression(a, scope, val)
	}

	var name string
	switch nodeType := a.Name.(type) {
	case *ast.Identifier:
		name = nodeType.Value
	case *ast.IndexExpression:
		switch nodeType.Left.(type) {
		case *ast.Identifier:
			name = nodeType.Left.(*ast.Identifier).Value

			//check if it's a class indexer assignment, e.g. 'clsObj[index] = xxx'
			if aObj, ok := scope.Get(name); ok {
				if aObj.Type() == INSTANCE_OBJ {
					return evalClassIndexerAssignExpression(a, aObj, nodeType, val, scope)
				}
			}
		case *ast.IndexExpression:
			leftVal := Eval(nodeType.Left, scope)
			indexVal := Eval(nodeType.Index, scope)
			//fmt.Printf("leftVal.Value=%v, leftVal.Type=%T, leftVal=%s\n", leftVal, leftVal, leftVal.Inspect())
			//fmt.Printf("indexVal.Value=%v, indexVal.Type=%T, indexVal=%s\n", indexVal, indexVal, indexVal.Inspect())
			switch v := leftVal.(type) {
			case *String:
				v.Set(a.Pos().Sline(), indexVal, val)
			case *Hash:
				v.Push(a.Pos().Sline(), indexVal, val)
			case *Array:
				v.Set(a.Pos().Sline(), indexVal, val)
			case *Tuple:
				str := fmt.Sprintf("%s[IDX]", TUPLE_OBJ)
				return NewError(a.Pos().Sline(), INFIXOP, str, "=", val.Type())
			}
			return
		}
	}

	if a.Token.Literal == "=" {
		switch nodeType := a.Name.(type) {
		case *ast.Identifier:
			name := nodeType.Value

			// class Dog {
			// 	static let age = 12
			// 	fn Walk() {
			// 		age = 20
			// 	}
			// }
			//
			// let dogObj = new Dog()
			// dogObj.Walk();
			//
			// In function 'Walk', we assign 'age' to 20, but 'age' is a static variable,
			// we need to use 'Dog.age', not bare 'age', so the below check code.
			thisObj, _ := scope.Get("this")
			if thisObj != nil {
				if thisObj.Type() == INSTANCE_OBJ { //'this' refers to 'ObjectInstance' object
					_, ok := thisObj.(*ObjectInstance).Scope.Get(name)
					if ok {
						v, ok2 := scope.Reset(name, val)
						if ok2 {
							return v
						}
					} else {
						thisObjClass := thisObj.(*ObjectInstance).Class
						_, ok := thisObjClass.Scope.Get(name)
						if ok {
							return NewError(a.Pos().Sline(), UNKNOWNIDENTEX, name, thisObjClass.Name+"."+name)
						}
					}
				}
			}

			//check if it's a readonly variable
			if scope.IsReadOnly(name) {
				return NewError(a.Pos().Sline(), CONSTNOTASSIGNERROR, name)
			}

			v, ok := scope.Reset(name, val)
			if ok {
				return v
			}
			return reportTypoSuggestions(a.Pos().Sline(), scope, a.Name.String())
			//return NewError(a.Pos().Sline(), UNKNOWNIDENT, a.Name.String())
		}
	}

	// Check if the variable exists or not
	var left Object
	var ok bool
	if left, ok = scope.Get(name); !ok {
		return reportTypoSuggestions(a.Pos().Sline(), scope, name)
		//return NewError(a.Pos().Sline(), UNKNOWNIDENT, name)
	}

	switch left.Type() {
	case INTEGER_OBJ, UINTEGER_OBJ, FLOAT_OBJ:
		val = evalNumAssignExpression(a, name, left, scope, val)
		return
	case STRING_OBJ:
		val = evalStrAssignExpression(a, name, left, scope, val)
		return
	case ARRAY_OBJ:
		val = evalArrayAssignExpression(a, name, left, scope, val)
		return
	case HASH_OBJ:
		val = evalHashAssignExpression(a, name, left, scope, val)
		return
	case TUPLE_OBJ:
		val = evalTupleAssignExpression(a, name, left, scope, val)
		return
	}

	return NewError(a.Pos().Sline(), INFIXOP, left.Type(), a.Token.Literal, val.Type())
}

func evalReturnStatement(r *ast.ReturnStatement, scope *Scope) Object {
	ret := &ReturnValue{Value: NIL, Values: []Object{}}
	for _, value := range r.ReturnValues {
		ret.Values = append(ret.Values, Eval(value, scope))
	}

	// for old campatibility
	if len(ret.Values) > 0 {
		ret.Value = ret.Values[0]
	}

	return ret
}

func evalDeferStatement(d *ast.DeferStmt, scope *Scope) Object {
	frame := scope.CurrentFrame()
	if frame == nil {
		return NewError(d.Pos().Sline(), DEFERERROR)
	}

	if frame.CurrentCall == nil {
		return NewError(d.Pos().Sline(), DEFERERROR)
	}

	switch d.Call.(type) {
	case *ast.CallExpression:
		callExp := d.Call.(*ast.CallExpression)
		closure := func() {
			evalFunctionCall(callExp, scope)
		}
		frame.defers = append(frame.defers, closure)
	case *ast.MethodCallExpression:
		callExp := d.Call.(*ast.MethodCallExpression)
		closure := func() {
			evalMethodCallExpression(callExp, scope)
		}
		frame.defers = append(frame.defers, closure)
	}

	return NIL
}

func evalThrowStatement(t *ast.ThrowStmt, scope *Scope) Object {
	throwObj := Eval(t.Expr, scope)
	if throwObj.Type() == ERROR_OBJ {
		return throwObj
	}

	return &Throw{stmt: t, value: throwObj}
}

// Booleans
func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Literals
func evalIntegerLiteral(i *ast.IntegerLiteral) Object {
	return NewInteger(i.Value)
}

func evalUIntegerLiteral(i *ast.UIntegerLiteral) Object {
	return NewUInteger(i.Value)
}

func evalFloatLiteral(f *ast.FloatLiteral) Object {
	return NewFloat(f.Value)
}

func evalStringLiteral(s *ast.StringLiteral) Object {
	return NewString(s.Value)
}

func evalInterpolatedString(is *ast.InterpolatedString, scope *Scope) Object {
	s := &InterpolatedString{String: &String{Valid: true}, RawValue: is.Value, Expressions: is.ExprMap}
	s.Interpolate(scope)
	return s
}

func evalArrayLiteral(a *ast.ArrayLiteral, scope *Scope) Object {
	if a.CreationCount == nil {
		return &Array{Members: evalArgs(a.Members, scope)}
	}

	var i int64
	ret := &Array{}
	for i = 0; i < a.CreationCount.Value; i++ {
		ret.Members = append(ret.Members, NIL)
	}

	return ret
}

func evalTupleLiteral(t *ast.TupleLiteral, scope *Scope) Object {
	return &Tuple{Members: evalArgs(t.Members, scope)}
}

func evalRegExLiteral(re *ast.RegExLiteral) Object {
	regExpression, err := regexp.Compile(re.Value)
	if err != nil {
		return NewError(re.Pos().Sline(), INVALIDARG)
	}

	return &RegEx{RegExp: regExpression, Value: re.Value}
}

func evalIdentifier(i *ast.Identifier, scope *Scope) Object {
	//Get from global scope first
	if obj, ok := GetGlobalObj(i.String()); ok {
		return obj
	}

	val, ok := scope.Get(i.String())
	if !ok {
		if val, ok = importScope.Get(i.String()); !ok {
			return reportTypoSuggestions(i.Pos().Sline(), scope, i.Value)
		}
	}
	if i, ok := val.(*InterpolatedString); ok {
		i.Interpolate(scope)
		return i
	}

	return val
}

func evalHashLiteral(hl *ast.HashLiteral, scope *Scope) Object {
	innerScope := NewScope(scope, nil)

	hash := NewHash()
	for _, key := range hl.Order {
		var k Object
		switch key.(type) {
		case *ast.Identifier: //It's an identifier, so it's a bare word.
			/* e.g. h = {A: "xxxx"}
			Here when evaluate the hash key 'A', it will evaluate to NIL_OBJ, because it's a bare word and
			is an identifier, so we need to treat it as string. that is, we want it to become:
			    h = {"A": "xxxx"}
			*/

			if _, ok := scope.Get(key.(*ast.Identifier).Value); !ok {
				t := key.(*ast.Identifier).Value
				k = NewString(t)
				innerScope.Set(t, k)
			}
		default:
			k = Eval(key, innerScope)
		}

		if k.Type() == ERROR_OBJ {
			return k
		}

		value, _ := hl.Pairs[key]
		v := Eval(value, innerScope)
		if v.Type() == ERROR_OBJ {
			return v
		}
		hash.Push(hl.Pos().Sline(), k, v)
	}
	return hash
}

func evalStructLiteral(s *ast.StructLiteral, scope *Scope) Object {
	structScope := NewScope(nil, scope.Writer)
	for key, value := range s.Pairs {
		if ident, ok := key.(*ast.Identifier); ok {
			aObj := Eval(value, scope)
			structScope.Set(ident.String(), aObj)
		} else {
			return NewError(s.Pos().Sline(), KEYERROR, "IDENT")
		}
	}
	return &Struct{Scope: structScope, methods: make(map[string]*Function)}
}

func evalEnumStatement(enumStmt *ast.EnumStatement, scope *Scope) Object {
	enumLiteral := evalEnumLiteral(enumStmt.EnumLiteral, scope)
	scope.Set(enumStmt.Name.String(), enumLiteral) //save to scope
	return enumLiteral
}

func evalEnumLiteral(e *ast.EnumLiteral, scope *Scope) Object {
	enumScope := NewScope(nil, scope.Writer)
	for key, value := range e.Pairs {
		if ident, ok := key.(*ast.Identifier); ok {
			aObj := Eval(value, scope)
			enumScope.Set(ident.String(), aObj)
		} else {
			return NewError(e.Pos().Sline(), KEYERROR, "IDENT")
		}
	}
	return &Enum{Scope: enumScope}
}

func evalRangeLiteral(r *ast.RangeLiteral, scope *Scope) Object {
	startIdx := Eval(r.StartIdx, scope)
	endIdx := Eval(r.EndIdx, scope)

	return evalRangeExpression(r, startIdx, endIdx, scope)
}

func evalFunctionStatement(FnStmt *ast.FunctionStatement, scope *Scope) Object {
	fnObj := evalFunctionLiteral(FnStmt.FunctionLiteral, scope)
	fn := fnObj.(*Function)

	if !FnStmt.IsServiceAnno {
		processClassAnnotation(FnStmt.Annotations, scope, FnStmt.Pos().Sline(), fn)
	}
	scope.Set(FnStmt.Name.String(), fnObj) //save to scope

	return fnObj
}

func evalFunctionLiteral(fl *ast.FunctionLiteral, scope *Scope) Object {
	fn := &Function{Literal: fl, Variadic: fl.Variadic, Scope: scope, Async: fl.Async}

	if fl.Values != nil { //check for default values
		for _, item := range fl.Parameters {
			if _, ok := fl.Values[item.String()]; !ok { //if not has default value, then continue
				continue
			}
			val := Eval(fl.Values[item.String()], scope)
			if val.Type() != ERROR_OBJ {
				fn.Scope.Set(item.String(), val)
			} else {
				return val
			}
		}
	}

	return fn
}

// Prefix expression for User Defined Operator
func evalPrefixExpressionUDO(p *ast.PrefixExpression, right Object, scope *Scope) Object {
	if fn, ok := scope.Get(p.Operator); ok {
		f := fn.(*Function)
		// set functions's parameters
		scope.Set(f.Literal.Parameters[0].String(), right)
		r := Eval(f.Literal.Body, scope)
		if r.Type() == ERROR_OBJ {
			return r
		}

		if obj, ok := r.(*ReturnValue); ok {
			// if function returns multiple-values
			// returns a tuple instead.
			if len(obj.Values) > 1 {
				return &Tuple{Members: obj.Values, IsMulti: true}
			}
			return obj.Value
		}
		return r
	}
	return NewError(p.Pos().Sline(), PREFIXOP, p, right.Type())
}

// Prefix expression for Meta-Operators
func evalMetaOperatorPrefixExpression(p *ast.PrefixExpression, right Object, scope *Scope) Object {
	if right.Type() != ARRAY_OBJ {
		return NewError(p.Pos().Sline(), PREFIXOP, p, right.Type())
	}

	//convert prefix operator to infix operator,
	//Because 'evalNumberInfixExpression' function need a InfixExpression
	infixExp := &ast.InfixExpression{Token: p.Token, Operator: p.Operator, Right: p.Right}

	members := right.(*Array).Members
	if len(members) == 0 {
		return NewInteger(0)
	}

	result := members[0]
	var leftIsNum bool
	var leftIsStr bool

	_, leftIsNum = result.(Number)
	if !leftIsNum {
		_, leftIsStr = result.(*String)
		if !leftIsStr {
			return NewError(p.Pos().Sline(), METAOPERATORERROR)
		}
	}

	for i := 1; i < len(members); i++ {
		var rightIsNum bool
		var rightIsStr bool
		_, rightIsNum = members[i].(Number)
		if !rightIsNum {
			_, rightIsStr = members[i].(*String)
			if !rightIsStr {
				return NewError(p.Pos().Sline(), METAOPERATORERROR)
			}
		}

		if leftIsStr || rightIsStr {
			result = evalMixedTypeInfixExpression(infixExp, result, members[i])
		} else {
			result = evalNumberInfixExpression(infixExp, result, members[i])
		}

		_, leftIsNum = result.(Number)
		if !leftIsNum {
			_, leftIsStr = result.(*String)
			if !leftIsStr {
				return NewError(p.Pos().Sline(), METAOPERATORERROR)
			}
		}

	} // end for

	return result
}

// Prefix expressions, e.g. `!true, -5`
func evalPrefixExpression(p *ast.PrefixExpression, scope *Scope) Object {
	right := Eval(p.Right, scope)
	if right.Type() == ERROR_OBJ {
		return right
	}

	//User Defined Operator
	if p.Token.Type == token.UDO {
		return evalPrefixExpressionUDO(p, right, scope)
	}

	if isMetaOperators(p.Token.Type) {
		return evalMetaOperatorPrefixExpression(p, right, scope)
	}

	if right.Type() == INSTANCE_OBJ {
		/* e.g. p.operator = '-':
		class vector {
			let x;
			let y;
			fn init (a, b) {
				this.x = a
				this.y = b
			}
			fn -() {
				return new Vector(-x,-y)
			}
		}
		v1 = new vector(3,4)
		v2 = -v1
		*/
		instanceObj := right.(*ObjectInstance)
		method := instanceObj.GetMethod(p.Operator)
		if method != nil {
			switch method.(type) {
			case *Function:
				newScope := NewScope(instanceObj.Scope, nil)
				args := []Object{right}
				return evalFunctionDirect(method, args, instanceObj, newScope, nil)
			case *BuiltinMethod:
				//do nothing for now
			}
		}
		return NewError(p.Pos().Sline(), PREFIXOP, p, right.Type())
	}

	switch p.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "+":
		switch right.Type() {
		case STRING_OBJ: //convert string to number
			var n int64
			var err error

			var content = right.(*String).String
			if strings.HasPrefix(content, "0b") {
				n, err = strconv.ParseInt(content[2:], 2, 64)
			} else if strings.HasPrefix(content, "0x") {
				n, err = strconv.ParseInt(content[2:], 16, 64)
			} else if strings.HasPrefix(content, "0o") {
				n, err = strconv.ParseInt(content[2:], 8, 64)
			} else {
				if len(content) == 0 {
					return NewInteger(0)
				}

				n, err = strconv.ParseInt(content, 10, 64)
				if err != nil {
					// Check if it is a float
					var f float64
					f, err = strconv.ParseFloat(content, 64)
					if err == nil {
						return NewFloat(f)
					}
				}
			}
			if err != nil {
				return NewError(p.Pos().Sline(), PREFIXOP, p, right.Type())
			}
			return NewInteger(n)
		case BOOLEAN_OBJ: //convert boolean to string
			var b = right.(*Boolean)
			if !b.Valid {
				return NewString("false")
			}

			if b.Bool {
				return NewString("true")
			}
			return NewString("false")
		default:
			return right
		}
	case "-":
		switch right.Type() {
		case INTEGER_OBJ:
			i := right.(*Integer)
			return NewInteger(-i.Int64)
			//bug : we need to return a new 'Integer' object, we should not change the original 'Integer' object.
			//i.Int64 = -i.Int64
			//return i
		case UINTEGER_OBJ:
			i := right.(*UInteger)
			if i.UInt64 == 0 {
				return i
			} else {
				return NewError(p.Pos().Sline(), PREFIXOP, p, right.Type())
			}
		case FLOAT_OBJ:
			f := right.(*Float)
			return NewFloat(-f.Float64)
			//bug : we need to return a new 'Float' object, we should not change the original 'Float' object.
			//f.Float64 = -f.Float64
			//return f
		}

	case "++":
		return evalIncrementPrefixOperatorExpression(p, right, scope)
	case "--":
		return evalDecrementPrefixOperatorExpression(p, right, scope)
	}
	return NewError(p.Pos().Sline(), PREFIXOP, p, right.Type())
}

func evalIncrementPrefixOperatorExpression(p *ast.PrefixExpression, right Object, scope *Scope) Object {
	switch right.Type() {
	case INTEGER_OBJ:
		rightObj := right.(*Integer)
		rightVal := rightObj.Int64
		scope.Reset(p.Right.String(), NewInteger(rightVal+1))
		return NewInteger(rightVal + 1)
	case UINTEGER_OBJ:
		rightObj := right.(*UInteger)
		rightVal := rightObj.UInt64
		scope.Reset(p.Right.String(), NewUInteger(rightVal+1))
		return NewUInteger(rightVal + 1)
	case FLOAT_OBJ:
		rightObj := right.(*Float)
		rightVal := rightObj.Float64
		scope.Reset(p.Right.String(), NewFloat(rightVal+1))
		return NewFloat(rightVal + 1)
	default:
		return NewError(p.Pos().Sline(), PREFIXOP, p.Operator, right.Type())
	}
}

func evalDecrementPrefixOperatorExpression(p *ast.PrefixExpression, right Object, scope *Scope) Object {
	switch right.Type() {
	case INTEGER_OBJ:
		rightObj := right.(*Integer)
		rightVal := rightObj.Int64
		scope.Reset(p.Right.String(), NewInteger(rightVal-1))
		return NewInteger(rightVal - 1)
	case UINTEGER_OBJ:
		rightObj := right.(*UInteger)
		rightVal := rightObj.UInt64
		scope.Reset(p.Right.String(), NewUInteger(rightVal-1))
		return NewUInteger(rightVal - 1)
	case FLOAT_OBJ:
		rightObj := right.(*Float)
		rightVal := rightObj.Float64
		scope.Reset(p.Right.String(), NewFloat(rightVal-1))
		return NewFloat(rightVal - 1)
	default:
		return NewError(p.Pos().Sline(), PREFIXOP, p.Operator, right.Type())
	}
}

// Helper for evaluating Bang(!) expressions. Coerces truthyness based on object presence.
func evalBangOperatorExpression(right Object) Object {
	return nativeBoolToBooleanObject(!IsTrue(right))
}

// Infix expression for User Defined Operator
func evalInfixExpressionUDO(p *ast.InfixExpression, left Object, right Object, scope *Scope) Object {
	if fn, ok := scope.Get(p.Operator); ok {
		f := fn.(*Function)
		// set functions two parameters
		scope.Set(f.Literal.Parameters[0].String(), left)
		scope.Set(f.Literal.Parameters[1].String(), right)
		r := Eval(f.Literal.Body, scope)
		if r.Type() == ERROR_OBJ {
			return r
		}

		if obj, ok := r.(*ReturnValue); ok {
			// if function returns multiple-values
			// returns a tuple instead.
			if len(obj.Values) > 1 {
				return &Tuple{Members: obj.Values, IsMulti: true}
			}
			return obj.Value
		}
		return r
	}

	return NewError(p.Pos().Sline(), INFIXOP, left.Type(), p.Operator, right.Type())
}

// Infix expression for Meta-Operators
func evalMetaOperatorInfixExpression(p *ast.InfixExpression, left Object, right Object, scope *Scope) Object {
	//1. [1,2,3] ~+ [4,5,6] = [1+4, 2+5, 3+6]
	//2. [1,2,3] ~+ 4 = [1+4, 2+4, 3+4]
	//left must be an array
	if left.Type() != ARRAY_OBJ {
		return NewError(p.Pos().Sline(), INFIXOP, left.Type(), p.Operator, right.Type())
	}

	leftMembers := left.(*Array).Members
	leftNumLen := len(leftMembers)

	//right could be an array or a number
	var rightMembers []Object
	_, rightIsNum := right.(Number)
	if rightIsNum {
		for i := 0; i < leftNumLen; i++ {
			rightMembers = append(rightMembers, right)
		}
	} else {
		if right.Type() == ARRAY_OBJ {
			rightMembers = right.(*Array).Members
		} else {
			return NewError(p.Pos().Sline(), INFIXOP, left.Type(), p.Operator, right.Type())
		}
	}
	rightNumLen := len(rightMembers)

	if leftNumLen != rightNumLen {
		return NewError(p.Pos().Sline(), GENERICERROR, "Number of items not equal for Meta-Operators!")
	}

	resultArr := &Array{}
	if leftNumLen == 0 {
		return resultArr
	}

	for idx, item := range leftMembers {
		var leftIsNum, rightIsNum bool
		var leftIsStr, rightIsStr bool
		_, leftIsNum = item.(Number)
		if !leftIsNum {
			_, leftIsStr = item.(*String)
			if !leftIsStr {
				return NewError(p.Pos().Sline(), METAOPERATORERROR)
			}
		}

		_, rightIsNum = rightMembers[idx].(Number)
		if !rightIsNum {
			_, rightIsStr = rightMembers[idx].(*String)
			if !rightIsStr {
				return NewError(p.Pos().Sline(), METAOPERATORERROR)
			}
		}

		var result Object
		if leftIsNum && rightIsNum {
			result = evalNumberInfixExpression(p, item, rightMembers[idx])
		} else if leftIsStr && rightIsStr {
			result = evalStringInfixExpression(p, item, rightMembers[idx])
		} else if leftIsStr || rightIsStr {
			result = evalMixedTypeInfixExpression(p, item, rightMembers[idx])
		} else {
			return NewError(p.Pos().Sline(), METAOPERATORERROR)
		}
		resultArr.Members = append(resultArr.Members, result)
	} // end for

	return resultArr

}

// Evaluate infix expressions, e.g 1 + 2, a == 5, true == true, etc...
func evalInfixExpression(node *ast.InfixExpression, left, right Object, scope *Scope) Object {
	if isGoObj(left) {
		left = GoValueToObject(left.(*GoObject).obj)
	}
	if isGoObj(right) {
		right = GoValueToObject(right.(*GoObject).obj)
	}

	//User Defined Operator
	if node.Token.Type == token.UDO {
		return evalInfixExpressionUDO(node, left, right, scope)
	}

	//Null-Coalescing Operator(??)
	if node.Token.Type == token.QUESTIONMM {
		if left.Type() == NIL_OBJ {
			return right
		}
		return left
	}

	if isMetaOperators(node.Token.Type) {
		return evalMetaOperatorInfixExpression(node, left, right, scope)
	}

	// Check if left is 'Writable'
	if _, ok := left.(Writable); ok { //There are two Writeables in magpie: FileObject, HttpResponseWriter.
		if node.Operator == ">>" { // '>>' is refered as 'extraction operator'. e.g.
			// Left is a file object
			if left.Type() == FILE_OBJ { // FileObject is also readable
				//    let a;
				//    stdin >> a
			}

			//right should be an identifier
			var rightVar *ast.Identifier
			var ok bool
			if rightVar, ok = node.Right.(*ast.Identifier); !ok { //not an identifier
				return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
			}
			f := left.(*FileObject)
			if f.File == os.Stdin {
				// 8192 is enough?
				ret := f.Read(node.Pos().Sline(), NewInteger(8192))
				scope.Set(rightVar.String(), ret)
				return ret
			} else {
				ret := f.ReadLine(node.Pos().Sline())
				scope.Set(rightVar.String(), ret)
				return ret
			}
		}

		if node.Operator == "<<" { // '<<' is refered as 'insertion operator'
			if f, ok := left.(*FileObject); ok { // It's a FileOject
				f.Write(node.Pos().Sline(), NewString(right.Inspect()))
				//Here we return left, so we can chain multiple '<<'.
				// e.g.
				//     stdout << "hello " << "world!"
				return left
			}
			if httpResp, ok := left.(*HttpResponseWriter); ok { // It's a HttpResponseWriter
				httpResp.Write(node.Pos().Sline(), NewString(right.Inspect()))
				return left
			}
		}
	}

	_, leftIsNum := left.(Number)
	_, rightIsNum := right.(Number)
	//hasNumArg := leftIsNum || rightIsNum

	//Note :Here the 'switch's order is important, if you change the order, it will evaluate differently
	//e.g. 1 + [2,3] + "45" = [1,2,3,"45"](it's an array), if you change
	//`case (left.Type() == ARRAY_OBJ || right.Type() == ARRAY_OBJ)` to a lower order in the case, it will
	//return [1,2,3]"45"(that is a string)
	switch {
	case node.Operator == "and" || node.Operator == "&&":
		leftCond := objectToNativeBoolean(left)
		if leftCond == false {
			return FALSE
		}

		rightCond := objectToNativeBoolean(right)
		return nativeBoolToBooleanObject(leftCond && rightCond)
	case node.Operator == "or" || node.Operator == "||":
		leftCond := objectToNativeBoolean(left)
		if leftCond == true {
			return TRUE
		}

		rightCond := objectToNativeBoolean(right)
		return nativeBoolToBooleanObject(leftCond || rightCond)
	case leftIsNum && rightIsNum:
		return evalNumberInfixExpression(node, left, right)
	case (left.Type() == ARRAY_OBJ || right.Type() == ARRAY_OBJ):
		return evalArrayInfixExpression(node, left, right, scope)
	case (left.Type() == TUPLE_OBJ || right.Type() == TUPLE_OBJ):
		return evalTupleInfixExpression(node, left, right, scope)
	case (left.Type() == TIME_OBJ || right.Type() == TIME_OBJ):
		return evalTimeInfixExpression(node, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(node, left, right)
	case (left.Type() == STRING_OBJ || right.Type() == STRING_OBJ):
		return evalMixedTypeInfixExpression(node, left, right)
	case (left.Type() == HASH_OBJ && right.Type() == HASH_OBJ):
		return evalHashInfixExpression(node, left, right)
	case left.Type() == INSTANCE_OBJ:
		return evalInstanceInfixExpression(node, left, right)
	case node.Operator == "==":
		if isGoObj(left) || isGoObj(right) { // if it's GoObject
			ret := compareGoObj(left, right)
			if ret {
				return TRUE
			}
			return FALSE
		}

		if left.Type() != right.Type() {
			return FALSE
		}

		//Here we need to special handling for `Boolean` object. Because most of the time `BOOLEAN` will
		//return TRUE and FALSE. But sometimes we have to returns a new `Boolean` object,
		//Here we need to compare `Boolean.Bool` or else when we using
		//   if (aBool == true)
		//it will return false, but actually aBool is true.
		if left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ {
			l := left.(*Boolean)
			r := right.(*Boolean)
			if l.Bool == r.Bool {
				return TRUE
			}
			return FALSE
		}

		if left.Type() == NIL_OBJ && right.Type() == NIL_OBJ { //(s == nil) should return true if s is nil
			return TRUE
		}

		return nativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		if isGoObj(left) || isGoObj(right) { // if it's GoObject
			ret := compareGoObj(left, right)
			if ret {
				return FALSE
			}
			return TRUE
		}

		if left.Type() != right.Type() {
			return TRUE
		}

		if left.Type() == BOOLEAN_OBJ && right.Type() == BOOLEAN_OBJ {
			l := left.(*Boolean)
			r := right.(*Boolean)
			if l.Bool != r.Bool {
				return TRUE
			}
			return FALSE
		}

		return nativeBoolToBooleanObject(left != right)
	}

	return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
}

func isMetaOperators(tokenType token.TokenType) bool {
	return tokenType == token.TILDEPLUS || // ~+
		tokenType == token.TILDEMINUS || // ~-
		tokenType == token.TILDEASTERISK || // ~*
		tokenType == token.TILDESLASH || // ~/
		tokenType == token.TILDEMOD || // ~%
		tokenType == token.TILDECARET // ~^
}

func objectToNativeBoolean(o Object) bool {
	if r, ok := o.(*ReturnValue); ok {
		o = r.Value
	}
	switch obj := o.(type) {
	case *Boolean:
		return obj.Bool
	case *Nil:
		return false
	case *Integer:
		if obj.Int64 == 0 {
			return false
		}
		return true
	case *UInteger:
		if obj.UInt64 == 0 {
			return false
		}
		return true
	case *Float:
		if obj.Float64 == 0.0 {
			return false
		}
		return true
	case *Array:
		if len(obj.Members) == 0 {
			return false
		}
		return true
	case *Tuple:
		if len(obj.Members) == 0 {
			return false
		}
		return true
	case *Hash:
		if len(obj.Pairs) == 0 {
			return false
		}
		return true
	case *GoObject:
		goObj := obj
		tmpObj := GoValueToObject(goObj.obj)
		return objectToNativeBoolean(tmpObj)
	default:
		return true
	}
}

func evalNumberInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	var leftVal float64
	var rightVal float64

	isInt := left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ
	isUInt := left.Type() == UINTEGER_OBJ && right.Type() == UINTEGER_OBJ

	if left.Type() == INTEGER_OBJ {
		leftVal = float64(left.(*Integer).Int64)
	} else if left.Type() == UINTEGER_OBJ {
		leftVal = float64(left.(*UInteger).UInt64)
	} else {
		leftVal = left.(*Float).Float64
	}

	if right.Type() == INTEGER_OBJ {
		rightVal = float64(right.(*Integer).Int64)
	} else if right.Type() == UINTEGER_OBJ {
		rightVal = float64(right.(*UInteger).UInt64)
	} else {
		rightVal = right.(*Float).Float64
	}

	switch node.Operator {
	case "**", "~^":
		val := math.Pow(leftVal, rightVal)
		if isInt {
			return NewInteger(int64(val))
		} else if isUInt {
			return NewUInteger(uint64(val))
		} else {
			return checkNumInfix(left, right, val)
		}
	case "&":
		if isInt {
			val := int64(leftVal) & int64(rightVal)
			return NewInteger(int64(val))
		} else if isUInt {
			val := uint64(leftVal) & uint64(rightVal)
			return NewUInteger(uint64(val))
		}
		return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
	case "|":
		if isInt {
			val := int64(leftVal) | int64(rightVal)
			return NewInteger(int64(val))
		} else if isUInt {
			val := uint64(leftVal) | uint64(rightVal)
			return NewUInteger(uint64(val))
		}
		return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
	case "^":
		if isInt {
			val := int64(leftVal) ^ int64(rightVal)
			return NewInteger(int64(val))
		} else if isUInt {
			val := uint64(leftVal) ^ uint64(rightVal)
			return NewUInteger(uint64(val))
		}
		return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
	case "+", "~+":
		val := leftVal + rightVal
		if isInt {
			return NewInteger(int64(val))
		} else if isUInt {
			return NewUInteger(uint64(val))
		} else {
			return checkNumInfix(left, right, val)
		}
	case "-", "~-":
		val := leftVal - rightVal
		if isInt {
			return NewInteger(int64(val))
		} else if isUInt {
			return NewUInteger(uint64(val))
		} else {
			return checkNumInfix(left, right, val)
		}
	case "*", "~*":
		val := leftVal * rightVal
		if isInt {
			return NewInteger(int64(val))
		} else if isUInt {
			return NewUInteger(uint64(val))
		} else {
			return checkNumInfix(left, right, val)
		}
	case "/", "~/":
		if rightVal == 0 {
			return NewError(node.Pos().Sline(), DIVIDEBYZERO)
		}
		val := leftVal / rightVal
		//Should Always return float
		return NewFloat(val)
	case "%", "~%":
		if isInt {
			return NewInteger(int64(leftVal) % int64(rightVal))
		} else if isUInt {
			return NewUInteger(uint64(leftVal) % uint64(rightVal))
		}
		return NewFloat(math.Mod(leftVal, rightVal))
	case ">>":
		if isInt {
			aRes := uint64(leftVal) >> uint64(rightVal)
			return NewInteger(int64(aRes)) //NOTE: CAST MAYBE NOT CORRECT
		} else if isUInt {
			aRes := uint64(leftVal) >> uint64(rightVal)
			return NewUInteger(uint64(aRes))
		}
	case "<<":
		if isInt {
			aRes := uint64(leftVal) << uint64(rightVal)
			return NewInteger(int64(aRes)) //NOTE: CAST MAYBE NOT CORRECT
		} else if isUInt {
			aRes := uint64(leftVal) << uint64(rightVal)
			return NewUInteger(uint64(aRes))
		}

	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
	}

	return NIL
}

func checkNumInfix(left Object, right Object, val float64) Object {
	if (left.Type() == INTEGER_OBJ && right.Type() == UINTEGER_OBJ) ||
		(left.Type() == UINTEGER_OBJ && right.Type() == INTEGER_OBJ) {
		if val > math.MaxInt64 {
			return NewUInteger(uint64(val))
		} else {
			return NewInteger(int64(val))
		}
	}

	return NewFloat(val)
}

func evalStringInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	l := left.(*String)
	r := right.(*String)

	switch node.Operator {
	case "=~": //match
		matched, _ := regexp.MatchString(r.String, l.String)
		if matched {
			return TRUE
		}
		return FALSE

	case "!~": //not match
		matched, _ := regexp.MatchString(r.String, l.String)
		if matched {
			return FALSE
		}
		return TRUE

	case "==":
		return nativeBoolToBooleanObject(l.String == r.String)
	case "!=":
		return nativeBoolToBooleanObject(l.String != r.String)
	case "+":
		return NewString(l.String + r.String)
	case "<":
		return nativeBoolToBooleanObject(l.String < r.String)
	case "<=":
		return nativeBoolToBooleanObject(l.String <= r.String)
	case ">":
		return nativeBoolToBooleanObject(l.String > r.String)
	case ">=":
		return nativeBoolToBooleanObject(l.String >= r.String)
	}
	return NewError(node.Pos().Sline(), INFIXOP, l.Type(), node.Operator, r.Type())
}

func evalTimeTimeInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	l := left.(*TimeObj)
	r := right.(*TimeObj)

	var b bool
	switch node.Operator {
	case "==":
		b = l.Tm.Equal(r.Tm)
		return NewBooleanObj(b)
	case "!=":
		b = !l.Tm.Equal(r.Tm)
		return NewBooleanObj(b)
	case "<":
		b = l.Tm.Before(r.Tm)
		return NewBooleanObj(b)
	case "<=":
		b = l.Tm.Equal(r.Tm) || l.Tm.Before(r.Tm)
		return NewBooleanObj(b)
	case ">":
		b = l.Tm.After(r.Tm)
		return NewBooleanObj(b)
	case ">=":
		b = l.Tm.Equal(r.Tm) || l.Tm.After(r.Tm)
		return NewBooleanObj(b)
	}
	return NewError(node.Pos().Sline(), INFIXOP, l.Type(), node.Operator, r.Type())
}

/*
	The Duration String support "YMDhms":
		Y:Year    M:Month    D:Day
		h:hour    m:Minute   s:Second

	let dt1 = dt/2018-01-01 12:01:00/ + "-12m"

*/
func evalTimeStringInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	l := left.(*TimeObj)
	r := right.(*String)

	switch node.Operator {
	case "+":
		timeObj, err := ParseDuration(l, r.String)
		if err != nil {
			msg := fmt.Sprintf("Invalid string duration '%s'", r.String)
			return NewError(node.Pos().Sline(), GENERICERROR, msg)
		}
		return timeObj
	}

	return NewError(node.Pos().Sline(), INFIXOP, l.Type(), node.Operator, r.Type())

}

/*
   let dt1 = dt/2018-01-01 12:01:00/
   let dt2 = dt/2019-01-01 12:01:00/
   pringln(dt1 <= dt2) # result: true
*/
func evalTimeInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	if left.Type() == TIME_OBJ && right.Type() == TIME_OBJ {
		return evalTimeTimeInfixExpression(node, left, right)
	}

	if left.Type() == TIME_OBJ && right.Type() == STRING_OBJ {
		return evalTimeStringInfixExpression(node, left, right)
	}

	return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
}

func evalMixedTypeInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	if isGoObj(left) {
		left = GoValueToObject(left.(*GoObject).obj)
	}
	if isGoObj(right) {
		right = GoValueToObject(right.(*GoObject).obj)
	}

	switch node.Operator {
	case "+", "~+":
		return NewString(left.Inspect() + right.Inspect())
	case "*", "~*":
		if left.Type() == INTEGER_OBJ {
			integer := left.(*Integer).Int64
			if integer <= 0 {
				return NewString("")
			}
			return NewString(strings.Repeat(right.Inspect(), int(integer)))
		} else if left.Type() == UINTEGER_OBJ {
			uinteger := left.(*UInteger).UInt64
			if uinteger == 0 {
				return NewString("")
			}
			return NewString(strings.Repeat(right.Inspect(), int(uinteger)))
		}
		if right.Type() == INTEGER_OBJ {
			integer := right.(*Integer).Int64
			if integer <= 0 {
				return NewString("")
			}
			return NewString(strings.Repeat(left.Inspect(), int(integer)))
		} else if right.Type() == UINTEGER_OBJ {
			uinteger := right.(*UInteger).UInt64
			if uinteger == 0 {
				return NewString("")
			}
			return NewString(strings.Repeat(left.Inspect(), int(uinteger)))
		}
		return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
	case "==":
		if isGoObj(left) || isGoObj(right) { // if it's GoObject
			ret := compareGoObj(left, right)
			if ret {
				return TRUE
			}
			return FALSE
		}

		if left.Type() != right.Type() {
			return FALSE
		}

		if left.Type() != STRING_OBJ || right.Type() != STRING_OBJ {
			return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
		}

		if left.(*String).String == right.(*String).String {
			return TRUE
		}
		return FALSE

	case "!=":
		if isGoObj(left) || isGoObj(right) { // if it's GoObject
			ret := compareGoObj(left, right)
			if ret {
				return FALSE
			}
			return TRUE
		}

		if left.Type() != right.Type() {
			return TRUE
		}

		if left.Type() != STRING_OBJ || right.Type() != STRING_OBJ {
			return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
		}

		if left.(*String).String != right.(*String).String {
			return TRUE
		}
		return FALSE

	case "=~": //match
		if left.Type() == NIL_OBJ { //nil is not matched with anything
			return FALSE
		}

		var str string
		if left.Type() == INTEGER_OBJ {
			str = fmt.Sprintf("%d", left.(*Integer).Int64)
		} else if left.Type() == UINTEGER_OBJ {
			str = fmt.Sprintf("%d", left.(*UInteger).UInt64)
		} else if left.Type() == FLOAT_OBJ {
			str = fmt.Sprintf("%g", left.(*Float).Float64)
		}
		matched, _ := regexp.MatchString(right.(*String).String, str)
		if matched {
			return TRUE
		}
		return FALSE

	case "!~": //not match
		if left.Type() == NIL_OBJ {
			return TRUE
		}

		var str string
		if left.Type() == INTEGER_OBJ {
			str = fmt.Sprintf("%d", left.(*Integer).Int64)
		} else if left.Type() == UINTEGER_OBJ {
			str = fmt.Sprintf("%d", left.(*UInteger).UInt64)
		} else if left.Type() == FLOAT_OBJ {
			str = fmt.Sprintf("%g", left.(*Float).Float64)
		}
		matched, _ := regexp.MatchString(right.(*String).String, str)
		if matched {
			return FALSE
		}
		return TRUE

	default:
		return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
	}

	//return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
}

//array + item
//item + array
//array + array
//array == array
//array != array
//array << item (<< item)
func evalArrayInfixExpression(node *ast.InfixExpression, left Object, right Object, scope *Scope) Object {
	switch node.Operator {
	case "*": // [1,2,3] * 2 = [1,2,3,1,2,3]
		if left.Type() == ARRAY_OBJ && right.Type() == INTEGER_OBJ {
			leftVals := left.(*Array).Members
			rightVal := right.(*Integer).Int64

			var i int64
			result := &Array{}
			for i = 0; i < rightVal; i++ {
				result.Members = append(result.Members, leftVals...)
			}
			return result
		}
	case "+":
		if left.Type() == ARRAY_OBJ {
			leftVals := left.(*Array).Members

			if right.Type() == ARRAY_OBJ {
				rightVals := right.(*Array).Members
				leftVals = append(leftVals, rightVals...)
			} else {
				leftVals = append(leftVals, right)
			}
			return &Array{Members: leftVals}
		}

		//right is array
		rightVals := right.(*Array).Members
		if left.Type() == ARRAY_OBJ {
			leftVals := left.(*Array).Members
			rightVals = append(rightVals, leftVals...)
			return &Array{Members: rightVals}
		} else {
			ret := &Array{}
			ret.Members = append(ret.Members, left)
			ret.Members = append(ret.Members, rightVals...)
			return ret
		}

	case "==":
		if left.Type() != right.Type() {
			return FALSE
		}

		if left.Type() != ARRAY_OBJ || right.Type() != ARRAY_OBJ {
			return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
		}

		leftVals := left.(*Array).Members
		rightVals := right.(*Array).Members
		if len(leftVals) != len(rightVals) {
			return FALSE
		}

		for i := range leftVals {
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i], scope)
			if !IsTrue(aBool) {
				return FALSE
			}
		}
		return TRUE
	case "!=":
		if left.Type() != right.Type() {
			return TRUE
		}

		if left.Type() != ARRAY_OBJ || right.Type() != ARRAY_OBJ {
			return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
		}
		leftVals := left.(*Array).Members
		rightVals := right.(*Array).Members
		if len(leftVals) != len(rightVals) {
			return TRUE
		}

		for i := range leftVals {
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i], scope)
			if IsTrue(aBool) {
				return TRUE
			}
		}
		return FALSE
	case "<<":
		if left.Type() == ARRAY_OBJ {
			leftVals := left.(*Array).Members
			leftVals = append(leftVals, right)
			left.(*Array).Members = leftVals //Change the array itself
			return left                      //return the original array, so it could be chained by another '<<'
		}
	}
	return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
}

//Almost same as evalArrayInfixExpression
//tuple + item
//item + tuple
//tuple + tuple
//tuple == tuple
//tuple != tuple
func evalTupleInfixExpression(node *ast.InfixExpression, left Object, right Object, scope *Scope) Object {
	switch node.Operator {
	case "+":
		if left.Type() == TUPLE_OBJ {
			leftVals := left.(*Tuple).Members

			if right.Type() == TUPLE_OBJ {
				rightVals := right.(*Tuple).Members
				leftVals = append(leftVals, rightVals...)
			} else {
				leftVals = append(leftVals, right)
			}
			return &Tuple{Members: leftVals}
		}

		//right is array
		rightVals := right.(*Tuple).Members
		if left.Type() == TUPLE_OBJ {
			leftVals := left.(*Tuple).Members
			rightVals = append(rightVals, leftVals...)
			return &Tuple{Members: rightVals}
		} else {
			ret := &Tuple{}
			ret.Members = append(ret.Members, left)
			ret.Members = append(ret.Members, rightVals...)
			return ret
		}

	case "==":
		if left.Type() != right.Type() {
			return FALSE
		}

		if left.Type() != TUPLE_OBJ || right.Type() != TUPLE_OBJ {
			return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
		}

		leftVals := left.(*Tuple).Members
		rightVals := right.(*Tuple).Members
		if len(leftVals) != len(rightVals) {
			return FALSE
		}

		for i := range leftVals {
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i], scope)
			if !IsTrue(aBool) {
				return FALSE
			}
		}
		return TRUE
	case "!=":
		if left.Type() != right.Type() {
			return TRUE
		}

		if left.Type() != TUPLE_OBJ || right.Type() != TUPLE_OBJ {
			return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
		}
		leftVals := left.(*Tuple).Members
		rightVals := right.(*Tuple).Members
		if len(leftVals) != len(rightVals) {
			return TRUE
		}

		for i := range leftVals {
			aBool := evalInfixExpression(node, leftVals[i], rightVals[i], scope)
			if IsTrue(aBool) {
				return TRUE
			}
		}
		return FALSE
	}
	return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
}

//hash + item
//hast == hash
//hash != hash
func evalHashInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	leftHash := left.(*Hash)
	rightHash := right.(*Hash)

	switch node.Operator {
	case "+":
		for _, hk := range rightHash.Order {
			pair, _ := rightHash.Pairs[hk]
			leftHash.Push(node.Pos().Sline(), pair.Key, pair.Value)
		}
		return leftHash
	case "==":
		return nativeBoolToBooleanObject(compareHashObj(leftHash.Pairs, rightHash.Pairs))
	case "!=":
		return nativeBoolToBooleanObject(!compareHashObj(leftHash.Pairs, rightHash.Pairs))
	}
	return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
}

func compareHashObj(left, right map[HashKey]HashPair) bool {
	if len(left) != len(right) {
		return false
	}

	found := 0
	for lk, lv := range left {
		for rk, rv := range right {
			if lk.Value == rk.Value && (lv.Key.Inspect() == rv.Key.Inspect() && lv.Value.Inspect() == rv.Value.Inspect()) {
				found += 1
				continue
			}
		}
	}

	return found == len(left)
}

// for operaotor overloading, e.g.
//    class Vector {
//        fn +(v) { xxxxx }
//    }
//
//    v1 = new Vector()
//    v2 = new Vector()
//    v3 = v1 + v2   //here is the operator overloading, same as 'v3 = v1.+(v2)
func evalInstanceInfixExpression(node *ast.InfixExpression, left Object, right Object) Object {
	instanceObj := left.(*ObjectInstance)

	switch node.Operator {
	case "==":
		if left.Type() != right.Type() {
			return FALSE
		}

		if left == right {
			return TRUE
		}
		return FALSE
	case "!=":
		if left.Type() != right.Type() {
			return TRUE
		}

		if left != right {
			return TRUE
		}
		return FALSE
	}
	//get methods's modifier level
	//	ml := instanceObj.GetModifierLevel(node.Operator, ClassMethodKind) //ml:modifier level
	//	if ml == ast.ModifierPrivate {
	//		return NewError(node.Pos().Sline(), CLSCALLPRIVATE, node.Operator, instanceObj.Class.Name)
	//	}

	method := instanceObj.GetMethod(node.Operator)
	if method != nil {
		switch m := method.(type) {
		case *Function:
			newScope := NewScope(instanceObj.Scope, nil)
			args := []Object{right}
			return evalFunctionDirect(method, args, instanceObj, newScope, nil)
		case *BuiltinMethod:
			args := []Object{right}
			builtinMethod := &BuiltinMethod{Fn: m.Fn, Instance: instanceObj}
			aScope := NewScope(instanceObj.Scope, nil)
			return evalFunctionDirect(builtinMethod, args, instanceObj, aScope, nil)
		}
	}
	return NewError(node.Pos().Sline(), INFIXOP, left.Type(), node.Operator, right.Type())
}

// IF macro statement: #ifdef xxx { block-statements } #else { block-statements }
func evalIfMacroStatement(im *ast.IfMacroStatement, scope *Scope) Object {
	if im.Condition {
		return evalBlockStatements(im.Consequence.Statements, scope)
	} else if im.Alternative != nil {
		return evalBlockStatements(im.Alternative.Statements, scope)
	}
	return NIL
}

func evalIfExpression(ie *ast.IfExpression, scope *Scope) Object {
	//eval "if/else-if" part
	for _, c := range ie.Conditions {
		condition := Eval(c.Cond, scope)
		if condition.Type() == ERROR_OBJ {
			return condition
		}

		if IsTrue(condition) {
			switch o := c.Body.(type) {
			case *ast.BlockStatement:
				return evalBlockStatements(o.Statements, scope)
			}
			return Eval(c.Body, scope)
		}
	}

	//eval "else" part
	if ie.Alternative != nil {
		switch o := ie.Alternative.(type) {
		case *ast.BlockStatement:
			return evalBlockStatements(o.Statements, scope)
		}
		return Eval(ie.Alternative, scope)
	}

	return NIL
}

func evalUnlessExpression(ie *ast.UnlessExpression, scope *Scope) Object {
	condition := Eval(ie.Condition, scope)
	if condition.Type() == ERROR_OBJ {
		return condition
	}

	if !IsTrue(condition) {
		return evalBlockStatements(ie.Consequence.Statements, scope)
	} else if ie.Alternative != nil {
		return evalBlockStatements(ie.Alternative.Statements, scope)
	}

	return NIL
}

func evalDoLoopExpression(dl *ast.DoLoop, scope *Scope) Object {
	newScope := NewScope(scope, nil)

	var e Object
	for {
		e = Eval(dl.Block, newScope)
		if e.Type() == ERROR_OBJ {
			return e
		}

		if _, ok := e.(*Break); ok {
			break
		}
		if _, ok := e.(*Continue); ok {
			continue
		}
		if v, ok := e.(*ReturnValue); ok {
			if v.Value != nil {
				//return v.Value
				return v
			}
			break
		}
	}

	if e == nil || e.Type() == BREAK_OBJ || e.Type() == CONTINUE_OBJ {
		return NIL
	}
	return e
}

func evalWhileLoopExpression(wl *ast.WhileLoop, scope *Scope) Object {
	innerScope := NewScope(scope, nil)

	var result Object
	for {
		condition := Eval(wl.Condition, innerScope)
		switch wl.Condition.(type) {
		case *ast.DiamondExpr:
			innerScope.Set("$_", condition)
		}

		if condition.Type() == ERROR_OBJ {
			return condition
		}

		if !IsTrue(condition) {
			return NIL
		}

		result = Eval(wl.Block, innerScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				/*
					BUG: DO NOT RETURN 'v.Value', instead we should return 'v'.

					If we return 'v.Value' then below code will print '5', not '6'(which is expected)
					let add = fn(x,y){
					    let i = 0
					    while (i++ < 10) {
					        return x * y
					    }
					    return x + y
					}
					println(add(2,3))
				*/

				//return v.Value
				return v
			}
			break
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return NIL
	}
	return result
}

func evalGrepExpression(ge *ast.GrepExpr, scope *Scope) Object {
	aValue := Eval(ge.Value, scope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members: []Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		return NewError(ge.Pos().Sline(), GREPMAPNOTITERABLE)
	}
	if !iterObj.iter() {
		return NewError(ge.Pos().Sline(), GREPMAPNOTITERABLE)
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
	} else if aValue.Type() == LINQ_OBJ {
		linqObj, _ := aValue.(*LinqObj)
		members = linqObj.ToSlice(ge.Pos().Sline()).(*Array).Members
	}

	result := &Array{}

	result.Members = []Object{}

	for _, item := range members {
		//Note: we must opening a new scope, because the variable is different in each iteration.
		//If not, then the next iteration will overwrite the previous assigned variable.
		newSubScope := NewScope(scope, nil)
		newSubScope.Set(ge.Var, item)

		var cond Object

		if ge.Block != nil {
			cond = Eval(ge.Block, newSubScope)
		} else {
			cond = Eval(ge.Expr, newSubScope)
		}

		if IsTrue(cond) {
			result.Members = append(result.Members, item)
		}
	}
	return result
}

func evalMapExpression(me *ast.MapExpr, scope *Scope) Object {
	aValue := Eval(me.Value, scope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members: []Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		return NewError(me.Pos().Sline(), GREPMAPNOTITERABLE)
	}
	if !iterObj.iter() {
		return NewError(me.Pos().Sline(), GREPMAPNOTITERABLE)
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
	} else if aValue.Type() == LINQ_OBJ {
		linqObj, _ := aValue.(*LinqObj)
		members = linqObj.ToSlice(me.Pos().Sline()).(*Array).Members
	}

	result := &Array{}
	result.Members = []Object{}

	for _, item := range members {
		newSubScope := NewScope(scope, nil)
		newSubScope.Set(me.Var, item)

		var r Object
		if me.Block != nil {
			r = Eval(me.Block, newSubScope)
		} else {
			r = Eval(me.Expr, newSubScope)
		}
		if r.Type() == ERROR_OBJ {
			return r
		}

		result.Members = append(result.Members, r)
	}
	return result
}

//[ x+1 for x in arr <where cond> ]
//[ str for str in strs <where cond> ]
//[ x for x in tuple <where cond> ]
func evalListComprehension(lc *ast.ListComprehension, scope *Scope) Object {
	innerScope := NewScope(scope, nil)
	aValue := Eval(lc.Value, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members: []Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		return NewError(lc.Pos().Sline(), NOTITERABLE)
	}
	if !iterObj.iter() {
		return NewError(lc.Pos().Sline(), NOTITERABLE)
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
	} else if aValue.Type() == LINQ_OBJ {
		linqObj, _ := aValue.(*LinqObj)
		members = linqObj.ToSlice(lc.Pos().Sline()).(*Array).Members
	}

	ret := &Array{}
	var result Object
	for idx, value := range members {
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(lc.Var, value)
		if lc.Cond != nil {
			cond := Eval(lc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(lc.Expr, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		ret.Members = append(ret.Members, result)
	}

	return ret
}

//[ x for x in a..b <where cond> ]
//Almost same as evalForEachDotRangeExpression() function
func evalListRangeComprehension(lc *ast.ListRangeComprehension, scope *Scope) Object {
	innerScope := NewScope(scope, nil)

	startIdx := Eval(lc.StartIdx, innerScope)
	endIdx := Eval(lc.EndIdx, innerScope)

	arr := evalRangeExpression(lc, startIdx, endIdx, scope).(*Array)

	ret := &Array{}
	var result Object
	for idx, value := range arr.Members {
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(lc.Var, value)
		if lc.Cond != nil {
			cond := Eval(lc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(lc.Expr, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		ret.Members = append(ret.Members, result)
	}

	return ret
}

//[ expr for k,v in hash <where cond> ]
func evalListMapComprehension(mc *ast.ListMapComprehension, scope *Scope) Object {
	innerScope := NewScope(scope, nil)
	aValue := Eval(mc.X, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members: []Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		return NewError(mc.Pos().Sline(), NOTITERABLE)
	}
	if !iterObj.iter() {
		return NewError(mc.Pos().Sline(), NOTITERABLE)
	}

	//must be a *Hash, if not, panic
	hash, _ := aValue.(*Hash)

	ret := &Array{}
	var result Object
	for _, hk := range hash.Order {
		pair, _ := hash.Pairs[hk]
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set(mc.Key, pair.Key)
		newSubScope.Set(mc.Value, pair.Value)

		if mc.Cond != nil {
			cond := Eval(mc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(mc.Expr, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		ret.Members = append(ret.Members, result)
	}

	return ret
}

//{ k:v for x in arr <where cond> }
//{ k:v for str in strs <where cond> }
//{ k:v for x in tuple <where cond> }
//Almost same as evalListComprehension
func evalHashComprehension(hc *ast.HashComprehension, scope *Scope) Object {
	innerScope := NewScope(scope, nil)
	aValue := Eval(hc.Value, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members: []Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		return NewError(hc.Pos().Sline(), NOTITERABLE)
	}
	if !iterObj.iter() {
		return NewError(hc.Pos().Sline(), NOTITERABLE)
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
	}

	ret := NewHash()

	for idx, value := range members {
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(hc.Var, value)
		if hc.Cond != nil {
			cond := Eval(hc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		keyResult := Eval(hc.KeyExpr, newSubScope)
		if keyResult.Type() == ERROR_OBJ {
			return keyResult
		}

		valueResult := Eval(hc.ValExpr, newSubScope)
		if valueResult.Type() == ERROR_OBJ {
			return valueResult
		}

		ret.Push(hc.Pos().Sline(), keyResult, valueResult)
		//ret.Pairs[hashable.HashKey()] = HashPair{Key: keyResult, Value: valueResult}
	}

	return ret
}

//{ k:v for x in a..b <where cond> }
//Almost same as evalListRangeComprehension() function
func evalHashRangeComprehension(hc *ast.HashRangeComprehension, scope *Scope) Object {
	innerScope := NewScope(scope, nil)

	startIdx := Eval(hc.StartIdx, innerScope)
	endIdx := Eval(hc.EndIdx, innerScope)

	arr := evalRangeExpression(hc, startIdx, endIdx, scope).(*Array)

	ret := NewHash()

	for idx, value := range arr.Members {
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(hc.Var, value)
		if hc.Cond != nil {
			cond := Eval(hc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		keyResult := Eval(hc.KeyExpr, newSubScope)
		if keyResult.Type() == ERROR_OBJ {
			return keyResult
		}

		valueResult := Eval(hc.ValExpr, newSubScope)
		if valueResult.Type() == ERROR_OBJ {
			return valueResult
		}

		ret.Push(hc.Pos().Sline(), keyResult, valueResult)
		//ret.Pairs[hashable.HashKey()] = HashPair{Key: keyResult, Value: valueResult}
	}

	return ret
}

//{ k:v for k,v in hash <where cond> }
//Almost same as evalListMapComprehension
func evalHashMapComprehension(mc *ast.HashMapComprehension, scope *Scope) Object {
	innerScope := NewScope(scope, nil)
	aValue := Eval(mc.X, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members: []Object{}}
	}

	iterObj, ok := aValue.(Iterable) //must be Iterable
	if !ok {
		return NewError(mc.Pos().Sline(), NOTITERABLE)
	}
	if !iterObj.iter() {
		return NewError(mc.Pos().Sline(), NOTITERABLE)
	}

	//must be a *Hash, if not, panic
	hash, _ := aValue.(*Hash)

	ret := NewHash()
	for _, hk := range hash.Order {
		pair, _ := hash.Pairs[hk]
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set(mc.Key, pair.Key)
		newSubScope.Set(mc.Value, pair.Value)

		if mc.Cond != nil {
			cond := Eval(mc.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		keyResult := Eval(mc.KeyExpr, newSubScope)
		if keyResult.Type() == ERROR_OBJ {
			return keyResult
		}

		valueResult := Eval(mc.ValExpr, newSubScope)
		if valueResult.Type() == ERROR_OBJ {
			return valueResult
		}

		if _, ok := keyResult.(Hashable); ok {
			ret.Push(mc.Pos().Sline(), keyResult, valueResult)
			//ret.Pairs[hashable.HashKey()] = HashPair{Key: keyResult, Value: valueResult}
		} else {
			return NewError(mc.Pos().Sline(), KEYERROR, keyResult.Type())
		}
	}

	return ret
}

func evalCaseExpression(ce *ast.CaseExpr, scope *Scope) Object {
	rv := Eval(ce.Expr, scope) //case expression
	if rv.Type() == ERROR_OBJ {
		return rv
	}

	done := false
	var elseExpr *ast.CaseElseExpr
	for _, item := range ce.Matches {
		if cee, ok := item.(*ast.CaseElseExpr); ok {
			elseExpr = cee //cee: Case'Expr Else part
			continue
		}

		matchExpr := item.(*ast.CaseMatchExpr)
		matchRv := Eval(matchExpr.Expr, NewScope(scope, nil)) //matcher expression
		if matchRv.Type() == ERROR_OBJ {
			return matchRv
		}

		//check 'rv' and 'matchRv' equality, if not equal, then continue
		if !equal(ce.IsWholeMatch, rv, matchRv) {
			continue
		}
		//Eval matcher block
		matcherScope := NewScope(scope, nil)
		rv = Eval(matchExpr.Block, matcherScope)
		if rv.Type() == ERROR_OBJ {
			return rv
		}

		done = true
		break
	}

	if !done && elseExpr != nil {
		elseScope := NewScope(scope, nil)
		rv = Eval(elseExpr.Block, elseScope)
		if rv.Type() == ERROR_OBJ {
			return rv
		}
	}
	return rv
}

func evalForLoopExpression(fl *ast.ForLoop, scope *Scope) Object { //fl:For Loop
	innerScope := NewScope(scope, nil)

	if fl.Init != nil {
		init := Eval(fl.Init, innerScope)
		if init.Type() == ERROR_OBJ {
			return init
		}
	}

	condition := Eval(fl.Cond, innerScope)
	if condition.Type() == ERROR_OBJ {
		return condition
	}

	var result Object
	for IsTrue(condition) {
		newSubScope := NewScope(innerScope, nil)
		result = Eval(fl.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			if fl.Update != nil {
				newVal := Eval(fl.Update, newSubScope) //Before continue, we need to call 'Update' and 'Cond'
				if newVal.Type() == ERROR_OBJ {
					return newVal
				}
			}

			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				//return v.Value
				return v
			}
			break
		}

		if fl.Update != nil {
			newVal := Eval(fl.Update, newSubScope)
			if newVal.Type() == ERROR_OBJ {
				return newVal
			}
		}

		condition = Eval(fl.Cond, newSubScope)
		if condition.Type() == ERROR_OBJ {
			return condition
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return NIL
	}
	return result
}

func evalForEverLoopExpression(fel *ast.ForEverLoop, scope *Scope) Object {
	var e Object
	newScope := NewScope(scope, nil)
	for {
		e = Eval(fel.Block, newScope)
		if e.Type() == ERROR_OBJ {
			return e
		}

		if _, ok := e.(*Break); ok {
			break
		}
		if _, ok := e.(*Continue); ok {
			continue
		}
		if v, ok := e.(*ReturnValue); ok {
			if v.Value != nil {
				//return v.Value
				return v
			}
			break
		}
	}

	if e == nil || e.Type() == BREAK_OBJ || e.Type() == CONTINUE_OBJ {
		return NIL
	}
	return e
}

func evalForEachFileLine(fal *ast.ForEachArrayLoop, scope *Scope) Object {
	for {
		line := Eval(fal.Value, scope)
		if line.Type() == ERROR_OBJ {
			return line
		}

		if line.Type() == NIL_OBJ { //at end-of-line
			break
		}
		scope.Set(fal.Var, line)

		if fal.Cond != nil {
			cond := Eval(fal.Cond, scope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result := Eval(fal.Block, scope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				return v
			}
			break
		}
	}

	return NIL
}

//for item in array
//for item in string
//for item in tuple
//for item in channel
//for item in goObj
//for item in linqObj
//for item in <$fileObj>
func evalForEachArrayExpression(fal *ast.ForEachArrayLoop, scope *Scope) Object { //fal:For Array Loop
	innerScope := NewScope(scope, nil)

	switch fal.Value.(type) {
	case *ast.DiamondExpr:
		return evalForEachFileLine(fal, innerScope)
	}

	aValue := Eval(fal.Value, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members: []Object{}}
	}

	iterObj, ok := aValue.(Iterable)
	if !ok {
		return NewError(fal.Pos().Sline(), NOTITERABLE)
	}
	if !iterObj.iter() {
		return NewError(fal.Pos().Sline(), NOTITERABLE)
	}

	var members []Object
	if aValue.Type() == STRING_OBJ {
		aStr, _ := aValue.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if aValue.Type() == ARRAY_OBJ {
		arr, _ := aValue.(*Array)
		members = arr.Members
	} else if aValue.Type() == TUPLE_OBJ {
		tuple, _ := aValue.(*Tuple)
		members = tuple.Members
	} else if aValue.Type() == LINQ_OBJ {
		linqObj := aValue.(*LinqObj)
		arr := linqObj.ToSlice(fal.Pos().Sline()).(*Array)
		members = arr.Members
	} else if aValue.Type() == GO_OBJ { // GoObject
		goObj := aValue.(*GoObject)
		arr := GoValueToObject(goObj.obj).(*Array)
		members = arr.Members
	} else if aValue.Type() == CHANNEL_OBJ {
		chanObj := aValue.(*ChanObject)
		ret := &Array{}
		var result Object

		idx := 0
		for value := range chanObj.ch {
			scope.Set("$_", NewInteger(int64(idx)))
			idx++
			scope.Set(fal.Var, value)
			result = Eval(fal.Block, scope)
			if result.Type() == ERROR_OBJ {
				return result
			}

			if _, ok := result.(*Break); ok {
				break
			}
			if _, ok := result.(*Continue); ok {
				continue
			}
			if v, ok := result.(*ReturnValue); ok {

				if v.Value != nil {
					return v
				}
				break
			} else {
				ret.Members = append(ret.Members, result)
			}

		} //end for
		if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
			return ret
		}
		return ret
	}

	ret := &Array{}
	var result Object
	for idx, value := range members {
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(fal.Var, value)
		if fal.Cond != nil {
			cond := Eval(fal.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(fal.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {

			if v.Value != nil {
				//ret.Members = append(ret.Members, v.Value)
				return v
				//return v.Value
			}
			break
		} else {
			ret.Members = append(ret.Members, result)
		}
	}

	//Here we need to check `nil`, because if the initial condition is not true, then `for`'s Body will have no
	//chance to execute, the result will be nil
	//this is the reason why we need to check for `BREAK_OBJ` or `CONTINUE_OBJ`:
	//    for i in 5..1 where i > 2 {
	//      if (i == 3) { continue }
	//      putln('i={i}')
	//    }
	//They will output "continue", this is not we expected
	//A LONG TIME HIDDEN BUG!
	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return ret
	}
	return ret
}

//for index, value in string
//for index, value in array
//for index, value in tuple
//for index, value in linqObj
func evalForEachArrayWithIndex(fml *ast.ForEachMapLoop, val Object, scope *Scope) Object {
	var members []Object
	if val.Type() == STRING_OBJ {
		aStr, _ := val.(*String)
		runes := []rune(aStr.String)
		for _, rune := range runes {
			members = append(members, NewString(string(rune)))
		}
	} else if val.Type() == ARRAY_OBJ {
		arr, _ := val.(*Array)
		members = arr.Members
	} else if val.Type() == TUPLE_OBJ {
		tuple, _ := val.(*Tuple)
		members = tuple.Members
	} else if val.Type() == LINQ_OBJ {
		linqObj := val.(*LinqObj)
		arr := linqObj.ToSlice(fml.Pos().Sline()).(*Array)
		members = arr.Members
	}

	ret := &Array{}
	var result Object
	for idx, value := range members {
		newSubScope := NewScope(scope, nil)
		newSubScope.Set(fml.Key, NewInteger(int64(idx)))
		newSubScope.Set(fml.Value, value)
		if fml.Cond != nil {
			cond := Eval(fml.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(fml.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {

			if v.Value != nil {
				return v
				//return v.Value
			}
			break
		} else {
			ret.Members = append(ret.Members, result)
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return ret
	}
	return ret
}

func evalForEachMapExpression(fml *ast.ForEachMapLoop, scope *Scope) Object { //fml:For Map Loop
	innerScope := NewScope(scope, nil)

	aValue := Eval(fml.X, innerScope)
	if aValue.Type() == ERROR_OBJ {
		return aValue
	}

	//first check if it's a Nil object
	if aValue.Type() == NIL_OBJ {
		//return an empty array object
		return &Array{Members: []Object{}}
	}

	iterObj, ok := aValue.(Iterable)
	if !ok {
		return NewError(fml.Pos().Sline(), NOTITERABLE)
	}
	if !iterObj.iter() {
		return NewError(fml.Pos().Sline(), NOTITERABLE)
	}

	//for index, value in arr
	//for index, value in string
	//for index, value in tuple
	//for index, value in linqObj
	if aValue.Type() == STRING_OBJ || aValue.Type() == ARRAY_OBJ || aValue.Type() == TUPLE_OBJ || aValue.Type() == LINQ_OBJ {
		return evalForEachArrayWithIndex(fml, aValue, innerScope)
	}

	hash, _ := aValue.(*Hash)

	ret := &Array{}
	var result Object
	for _, hk := range hash.Order {
		pair, _ := hash.Pairs[hk]
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set(fml.Key, pair.Key)
		newSubScope.Set(fml.Value, pair.Value)

		if fml.Cond != nil {
			cond := Eval(fml.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(fml.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				return v
			}
			break
		} else {
			ret.Members = append(ret.Members, result)
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return ret
	}
	return ret
}

func evalForEachDotRangeExpression(fdr *ast.ForEachDotRange, scope *Scope) Object { //fdr:For Dot Range
	innerScope := NewScope(scope, nil)

	startIdx := Eval(fdr.StartIdx, innerScope)
	//	if startIdx.Type() != INTEGER_OBJ {
	//		return NewError(fdr.Pos().Sline(), RANGETYPEERROR, startIdx.Type())
	//	}

	endIdx := Eval(fdr.EndIdx, innerScope)
	//	if startIdx.Type() != INTEGER_OBJ {
	//		return NewError(fdr.Pos().Sline(), RANGETYPEERROR, endIdx.Type())
	//	}

	arr := evalRangeExpression(fdr, startIdx, endIdx, scope).(*Array)

	ret := &Array{}
	var result Object
	for idx, value := range arr.Members {
		newSubScope := NewScope(innerScope, nil)
		newSubScope.Set("$_", NewInteger(int64(idx)))
		newSubScope.Set(fdr.Var, value)
		if fdr.Cond != nil {
			cond := Eval(fdr.Cond, newSubScope)
			if cond.Type() == ERROR_OBJ {
				return cond
			}

			if !IsTrue(cond) {
				continue
			}
		}

		result = Eval(fdr.Block, newSubScope)
		if result.Type() == ERROR_OBJ {
			return result
		}

		if _, ok := result.(*Break); ok {
			break
		}
		if _, ok := result.(*Continue); ok {
			continue
		}
		if v, ok := result.(*ReturnValue); ok {
			if v.Value != nil {
				//ret.Members = append(ret.Members, v.Value)
				return v
			}
			break
		} else {
			ret.Members = append(ret.Members, result)
		}
	}

	if result == nil || result.Type() == BREAK_OBJ || result.Type() == CONTINUE_OBJ {
		return ret
	}
	return ret
}

// Helper function IsTrue for IF evaluation - neccessity is dubious
func IsTrue(obj Object) bool {
	if b, ok := obj.(*Boolean); ok { //if it is a Boolean Object
		return b.Bool
	}

	if _, ok := obj.(*Nil); ok { //if it is a Nil Object
		return false
	}

	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		switch obj.Type() {
		case INTEGER_OBJ:
			if obj.(*Integer).Int64 == 0 {
				return false
			}
		case UINTEGER_OBJ:
			if obj.(*UInteger).UInt64 == 0 {
				return false
			}
		case FLOAT_OBJ:
			if obj.(*Float).Float64 == 0.0 {
				return false
			}

			//why remove below check? please see below code:
			//    for line in <$f> { println(line) }
			//Here when the line is empty, we should not return false.
			//		case STRING_OBJ:
			//			if len(obj.(*String).String) == 0 {
			//				return false
			//			}
		case ARRAY_OBJ:
			if len(obj.(*Array).Members) == 0 {
				return false
			}
		case HASH_OBJ:
			if len(obj.(*Hash).Pairs) == 0 {
				return false
			}
		case TUPLE_OBJ:
			if len(obj.(*Tuple).Members) == 0 {
				return false
			}
		case GO_OBJ:
			goObj := obj.(*GoObject)
			return goObj.obj != nil
		case OPTIONAL_OBJ:
			optionalObj := obj.(*Optional)
			if optionalObj.Value == NIL {
				return false
			}
		}
		return true
	}
}

// Block Statement Evaluation - The innards of both IF and Function calls
// very similar to parseProgram, but because we need to leave the return
// value wrapped in it's Object, it remains, for now.
func evalBlockStatements(block []ast.Statement, scope *Scope) (results Object) {
	for _, statement := range block {
		results = Eval(statement, scope)
		if results.Type() == ERROR_OBJ {
			return results
		}

		if results != nil && results.Type() == RETURN_VALUE_OBJ || results.Type() == THROW_OBJ {
			return
		}
		if _, ok := results.(*Break); ok {
			return
		}
		if _, ok := results.(*Continue); ok {
			return
		}
	}
	return //do not return NIL, becuase we have already set the 'results'
}

// Eval when a function is _called_, includes fn literal evaluation and calling builtins
func evalFunctionCall(call *ast.CallExpression, scope *Scope) Object {
	fn, ok := scope.Get(call.Function.String())
	if !ok {
		if f, ok := call.Function.(*ast.FunctionLiteral); ok {

			//let add =fn(x,y) { x+y }
			//add(2,3)
			fn = &Function{Literal: f, Scope: scope, Variadic: f.Variadic}
			scope.Set(call.Function.String(), fn)
		} else if idxExpr, ok := call.Function.(*ast.IndexExpression); ok { //index expression
			//let complex={ "add" : fn(x,y){ x+y } }
			//complex["add"](2,3)
			aValue := Eval(idxExpr, scope)
			if aValue.Type() == ERROR_OBJ {
				return aValue
			}

			if aFn, ok := aValue.(*Function); ok { //index expression
				fn = aFn
			} else {
				return reportTypoSuggestions(call.Function.Pos().Sline(), scope, call.Function.String())
				//return NewError(call.Function.Pos().Sline(), UNKNOWNIDENT, call.Function.String())
			}
		} else if builtin, ok := builtins[call.Function.String()]; ok {
			args := evalArgs(call.Arguments, scope)
			//check for errors
			for _, v := range args {
				if v.Type() == ERROR_OBJ {
					return v
				}
			}
			return builtin.Fn(call.Function.Pos().Sline(), scope, args...)
		} else if callExpr, ok := call.Function.(*ast.CallExpression); ok { //call expression
			//let complex={ "add" : fn(x,y){ fn(z) {x+y+z} } }
			//complex["add"](2,3)(4)
			aValue := Eval(callExpr, scope)
			if aValue.Type() == ERROR_OBJ {
				return aValue
			}

			fn = aValue
		} else {
			return reportTypoSuggestions(call.Function.Pos().Sline(), scope, call.Function.String())
			//return NewError(call.Function.Pos().Sline(), UNKNOWNIDENT, call.Function.String())
		}
	}

	if fn.Type() == CLASS_OBJ {
		return NewError(call.Function.Pos().Sline(), CLASSCREATEERROR, call.Function.String())
	}

	f := fn.(*Function)
	if f.Async && call.Awaited {
		aChan := make(chan Object, 1)

		go func() {
			defer close(aChan)
			aChan <- evalFunctionObj(call, f, scope)
		}()

		return <-aChan
	}

	if f.Async {
		go evalFunctionObj(call, f, scope)
		return NIL
	}

	return evalFunctionObj(call, f, scope)
}

func evalFunctionObj(call *ast.CallExpression, f *Function, scope *Scope) Object {
	var thisObj Object
	var ok bool
	//check if it's static function
	thisObj, ok = scope.Get("this")
	if ok {
		if thisObj.Type() == CLASS_OBJ { // 'this' refers to Class object iteself
			if !f.Literal.StaticFlag { //not static
				return NewError(call.Function.Pos().Sline(), CALLNONSTATICERROR)
			}
		} else { // 'this' refers to Class object instance
			//instance method could call static method
		}
	}

	newScope := NewScope(f.Scope, nil)

	//Register this function call in the call stack
	newScope.CallStack.Frames = append(newScope.CallStack.Frames, CallFrame{FuncScope: newScope, CurrentCall: call})

	//Using golang's defer mechanism, before function return, call current frame's defer method
	defer func() {
		frame := newScope.CurrentFrame()
		if len(frame.defers) != 0 {
			frame.runDefers(newScope)
		}

		//After run, must pop the frame
		stack := newScope.CallStack
		stack.Frames = stack.Frames[0 : len(stack.Frames)-1]
	}()

	variadicParam := []Object{}
	args := evalArgs(call.Arguments, scope)
	for i := range call.Arguments {
		//Because of function default values, we need to check `i >= len(args)`
		if f.Variadic && i >= len(f.Literal.Parameters)-1 {
			for j := i; j < len(args); j++ {
				variadicParam = append(variadicParam, args[j])
			}
			break
		} else if i >= len(f.Literal.Parameters) {
			break
		} else {
			newScope.Set(f.Literal.Parameters[i].String(), args[i])
		}
	}

	// Variadic argument is passed as a single array
	// of parameters.
	if f.Variadic {
		newScope.Set(f.Literal.Parameters[len(f.Literal.Parameters)-1].String(), &Array{Members: variadicParam})
		if len(call.Arguments) < len(f.Literal.Parameters) {
			f.Scope.Set("@_", NewInteger(int64(len(f.Literal.Parameters)-1)))
		} else {
			f.Scope.Set("@_", NewInteger(int64(len(call.Arguments))))
		}
	} else {
		f.Scope.Set("@_", NewInteger(int64(len(f.Literal.Parameters))))
	}

	r := Eval(f.Literal.Body, newScope)
	if r.Type() == ERROR_OBJ {
		return r
	}

	if obj, ok := r.(*ReturnValue); ok {
		// if function returns multiple-values
		// returns a tuple instead.
		if len(obj.Values) > 1 {
			return &Tuple{Members: obj.Values, IsMulti: true}
		}
		return obj.Value
	}

	/* If the function call do not end in a 'return' statement. e.g.
	   let add = fn(x, y) {
	       x + y // not 'return x + y'
	   }
	   We need to send EVAL_LINE to the debugger, so we can step into this line,
	   or else we cannot step into it.
	*/
	if Dbg != nil {
		MsgHandler.SendMessage(message.Message{Type: message.EVAL_LINE, Body: Context{N: []ast.Node{call}, S: newScope}})
	}
	return r
}

// Method calls for builtin Objects
func evalMethodCallExpression(call *ast.MethodCallExpression, scope *Scope) Object {
	//First check if is a stanard library object
	str := call.Object.String()
	if obj, ok := GetGlobalObj(str); ok {
		switch o := call.Call.(type) {
		case *ast.IndexExpression: // e.g. 'if gos.Args[0] == "hello" {'
			if arr, ok := GetGlobalObj(str + "." + o.Left.String()); ok {
				return evalArrayIndex(arr.(*Array), o, scope)
			}
		case *ast.Identifier: //e.g. os.O_APPEND
			if i, ok := GetGlobalObj(str + "." + o.String()); ok {
				return i
			} else { //e.g. method call like 'os.environ'
				if obj.Type() == HASH_OBJ { // It's a GoFuncObject
					hash := obj.(*Hash)
					for _, hk := range hash.Order {
						pair, _ := hash.Pairs[hk]
						funcName := pair.Key.(*String).String
						if funcName == o.String() {
							goFuncObj := pair.Value.(*GoFuncObject)
							return goFuncObj.CallMethod(call.Call.Pos().Sline(), scope, o.String())
						}
					}
				} else {
					return obj.CallMethod(call.Call.Pos().Sline(), scope, o.String())
				}
			}
		case *ast.CallExpression: //e.g. method call like 'os.environ()'
			if method, ok := call.Call.(*ast.CallExpression); ok {
				args := evalArgs(method.Arguments, scope)
				if obj.Type() == HASH_OBJ { // It's a GoFuncObject
					hash := obj.(*Hash)
					for _, hk := range hash.Order {
						pair, _ := hash.Pairs[hk]
						funcName := pair.Key.(*String).String
						if funcName == o.Function.String() {
							goFuncObj := pair.Value.(*GoFuncObject)
							return goFuncObj.CallMethod(call.Call.Pos().Sline(), scope, o.Function.String(), args...)
						}
					}
				} else {
					return obj.CallMethod(call.Call.Pos().Sline(), scope, o.Function.String(), args...)
				}
			}
		}
	} else {
		// if 'GetGlobalObj(str)' returns nil, then try below
		// e.g.
		//     eval.RegisterVars("runtime", map[string]interface{}{
		//          "GOOS": runtime.GOOS,
		//       })
		// The eval.RegisterVars will call SetGlobalObj("runtime.GOOS"), so the
		// global scope's name is 'runtime.GOOS', not 'runtime', therefore, the above
		// GetGlobalObj('runtime') will returns false.
		if obj, ok := GetGlobalObj(str + "." + call.Call.String()); ok {
			return obj
		}
	}

	obj := Eval(call.Object, scope)
	if obj.Type() == ERROR_OBJ {
		return obj
	}

	switch m := obj.(type) {
	case *ImportedObject:
		switch o := call.Call.(type) {
		case *ast.Identifier:
			idName := call.Call.String()
			if i, ok := m.Scope.Get(idName); ok {
				if !unicode.IsUpper(rune(idName[0])) {
					return NewError(call.Call.Pos().Sline(), NAMENOTEXPORTED, str, idName)
				} else {
					return i
				}
			} else {
				return reportTypoSuggestions(call.Call.Pos().Sline(), m.Scope, idName)
			}
		case *ast.CallExpression:
			if o.Function.String() == "Scope" {
				return obj.CallMethod(call.Call.Pos().Sline(), m.Scope, "Scope")
			}

			var fnObj Object
			var ok bool
			if fnObj, ok = m.Scope.Get(o.Function.String()); !ok {
				return reportTypoSuggestionsMeth(o.Function.Pos().Sline(), m.Scope, filepath.Base(m.Name), o.Function.String())
			}

			funcName := o.Function.String()
			if !unicode.IsUpper(rune(funcName[0])) {
				return NewError(o.Function.Pos().Sline(), NAMENOTEXPORTED, str, o.Function.String())
			}

			return evalFunctionObj(o, fnObj.(*Function), scope)
			//return evalFunctionCall(o, m.Scope)
		}
	case *Struct:
		switch o := call.Call.(type) {
		case *ast.Identifier:
			if i, ok := m.Scope.Get(call.Call.String()); ok {
				return i
			}
		case *ast.CallExpression:
			args := evalArgs(o.Arguments, scope)
			return obj.CallMethod(call.Call.Pos().Sline(), m.Scope, o.Function.String(), args...)
		}
	case *Enum:
		switch o := call.Call.(type) {
		case *ast.Identifier:
			if i, ok := m.Scope.Get(call.Call.String()); ok {
				return i
			}
		case *ast.CallExpression:
			args := evalArgs(o.Arguments, scope)
			return obj.CallMethod(call.Call.Pos().Sline(), m.Scope, o.Function.String(), args...)
		}
	case *Hash:
		switch o := call.Call.(type) {
		//e.g.:
		//hashObj.key1=10
		//println(hashObj.key1)
		case *ast.Identifier:
			keyObj := NewString(call.Call.String())
			return m.Get(call.Call.Pos().Sline(), keyObj)

		case *ast.CallExpression:
			//we need to get the hash key
			keyStr := strings.Split(call.Call.String(), "(")[0]
			keyObj := NewString(keyStr)
			hashPair, ok := m.Pairs[keyObj.HashKey()]
			if !ok {
				//Check if it's a hash object's builtin method(e.g. hashObj.keys(), hashObj.values())
				if method, ok := call.Call.(*ast.CallExpression); ok {
					args := evalArgs(method.Arguments, scope)
					return obj.CallMethod(call.Call.Pos().Sline(), scope, method.Function.String(), args...)
				}
			}

			//e.g.:
			//hashObj = {}
			//hashObj.str = fn() { return 10 }
			//hashObj.str()

			// we need 'FunctionLiteral' here, so we need to change 'o.Function',
			// because the o.Function's Type is '*ast.Identifier' which is the Hash's key
			o.Function = hashPair.Value.(*Function).Literal
			//return evalFunctionCall(o, scope)   This is a bug: not 'scope'
			return evalFunctionCall(o, hashPair.Value.(*Function).Scope) //should be Function's scope
		case *ast.IndexExpression:
			//e.g.:
			//doc = {"one": {"two": { "three": [1, 2, 3,] }}}
			//printf("doc.one.two.three[2]=%v\n", doc.one.two.three[2])
			leftVal := m.Get(call.Call.Pos().Sline(), NewString(o.Left.String()))
			indexVal := Eval(o.Index, scope)
			switch v := leftVal.(type) {
			case *Hash:
				return leftVal
			case *Array:
				return v.Get(call.Call.Pos().Sline(), indexVal)
			case *Tuple:
				return v.Get(call.Call.Pos().Sline(), indexVal)
			}
			return NIL
		}
	case *ObjectInstance:
		instanceObj := m
		switch o := call.Call.(type) {
		//e.g.: instanceObj.key1
		case *ast.Identifier:
			val, ok := instanceObj.Scope.Get(o.Value)
			if ok {
				// check if it's a static variable
				if instanceObj.IsStatic(o.Value, ClassMemberKind) {
					return NewError(call.Call.Pos().Sline(), MEMBERUSEERROR, o.Value, instanceObj.Class.Name)
				}

				switch val.(type) {
				case *Function: //Function without parameter. e.g. obj.getMonth(), could be called using 'obj.getMonth'
					return evalFunctionDirect(val, []Object{}, instanceObj, instanceObj.Scope, nil)
				default:
					return val
				}
			}

			//See if it's a property
			p := instanceObj.GetProperty(o.Value)
			if p != nil {
				// check if it's a static variable
				if instanceObj.IsStatic(o.Value, ClassPropertyKind) {
					return NewError(call.Call.Pos().Sline(), PROPERTYUSEERROR, o.Value, instanceObj.Class.Name)
				}

				if p.Getter == nil { //property xxx { set; }
					return NewError(call.Call.Pos().Sline(), PROPERTYUSEERROR, o.Value, instanceObj.Class.Name)
				} else {
					if len(p.Getter.Body.Statements) == 0 { //property xxx { get; }
						v, _ := instanceObj.Scope.Get("_" + o.Value)
						//instanceObj.Scope.Set("_" + o.Value, v)
						return v
					} else {
						results := Eval(p.Getter.Body, instanceObj.Scope)
						if results.Type() == RETURN_VALUE_OBJ {
							return results.(*ReturnValue).Value
						}
					}
				}
			}
			reportTypoSuggestions(call.Call.Pos().Sline(), instanceObj.Scope, o.Value)
			//return NewError(call.Call.Pos().Sline(), UNKNOWNIDENT, o.Value)

		case *ast.CallExpression:
			//e.g. instanceObj.method()
			fname := o.Function.String() // get function name

			isStatic := instanceObj.IsStatic(fname, ClassMethodKind)
			if isStatic {
				return NewError(call.Call.Pos().Sline(), GENERICERROR, "Method is static!")
			}

			method := instanceObj.GetMethod(fname)
			if method != nil {
				switch m := method.(type) {
				case *Function:
					newScope := NewScope(instanceObj.Scope, nil)
					args := evalArgs(o.Arguments, newScope)
					return evalFunctionDirect(method, args, instanceObj, newScope, o)

				case *BuiltinMethod:
					builtinMethod := &BuiltinMethod{Fn: m.Fn, Instance: instanceObj}
					aScope := NewScope(instanceObj.Scope, nil)
					args := evalArgs(o.Arguments, aScope)
					return evalFunctionDirect(builtinMethod, args, instanceObj, aScope, nil)
				}
			}
			return NewError(call.Call.Pos().Sline(), NOMETHODERROR, call.String(), obj.Type())
		}
	case *Class:
		clsObj := m
		switch o := call.Call.(type) {
		case *ast.Identifier: //e.g.: classObj.key1
			var val Object
			var ok bool

			thisObj, _ := scope.Get("this")
			if thisObj != nil && thisObj.Type() == INSTANCE_OBJ { //'this' refers to 'ObjectInstance' object
				val, ok = thisObj.(*ObjectInstance).Scope.Get(o.Value)
				if ok {
					return val
				} else {
					//Why this else? Please see below example:
					// class Dog{
					//	static let misc = 12
					//	fn MethodA() {
					//		printf("misc = %v\n", Dog.misc)
					//	}
					// }
					//
					//	dogObj = new Dog()
					//	dogObj.MethodA()
					//
					//When calling dogObj.MethodA, `thisObj` refers to `dogObj` instance,
					//when the above 'if' branch is executed, the `dogObj` instance's scope
					//has no such variable(i.e. misc), because it belongs to `Dog` class's static
					// variable, not instance variable so the 'if' test will fail. For it to walk,
					// we need to check for static variable of 'Dog' class.
					val, ok = clsObj.Scope.Get(o.Value)
					if ok {
						return val
					}
				}
				return NIL
			}

			//check if it's a property
			p := clsObj.GetProperty(o.Value)
			if p == nil { //not property, it's a member
				// check if it's a static member
				if !clsObj.IsStatic(o.Value, ClassMemberKind) {
					return NewError(call.Call.Pos().Sline(), MEMBERUSEERROR, o.Value, clsObj.Name)
				}
				val, ok = clsObj.Scope.Get(o.Value)
			} else {
				// check if it's a static property
				if !clsObj.IsStatic(o.Value, ClassPropertyKind) {
					return NewError(call.Call.Pos().Sline(), PROPERTYUSEERROR, o.Value, clsObj.Name)
				}
				val, ok = clsObj.Scope.Get(o.Value)
			}

			if ok {
				return val
			}
			return NIL

		case *ast.CallExpression: //e.g. classObj.method()
			newScope := clsObj.Scope
			fname := o.Function.String() // get function name

			// Here is the reason for checking 'this'. Below example code explains why:
			//
			//    Class SubClass : ParentClass {
			//        fn init(xxx) {
			//            parent.init(xxx)
			//        }
			//    }
			// subClassObj = new SubClass(xxx)
			//
			// When you instantiate 'subClassObj' using `new`,
			// the interpreter will call `init` method. In the 'init', the 'parent' is a Class object,
			// but 'this' is an instance, so the scope is the instance's scope, not class's scope.
			thisObj, _ := scope.Get("this")
			if thisObj != nil && thisObj.Type() == INSTANCE_OBJ { //'this' refers to 'ObjectInstance' object
				newScope = thisObj.(*ObjectInstance).Scope
			}

			method := clsObj.GetMethod(fname)
			if method != nil {
				isStatic := clsObj.IsStatic(fname, ClassMethodKind)
				if !isStatic {
					objName := str
					if objName != "parent" { // e.g. parent.method(parameters)
						return NewError(call.Call.Pos().Sline(), CALLNONSTATICERROR)
					}
				}

				switch m := method.(type) {
				case *Function:
					args := evalArgs(o.Arguments, scope)
					return evalFunctionDirect(m, args, nil, newScope, o)
				case *BuiltinMethod:
					builtinMethod := &BuiltinMethod{Fn: m.Fn, Instance: nil}
					aScope := NewScope(newScope, nil)
					args := evalArgs(o.Arguments, aScope)
					return evalFunctionDirect(builtinMethod, args, nil, aScope, nil)
				}
			} else {
				return reportTypoSuggestionsMeth(call.Call.Pos().Sline(), scope, clsObj.Name, fname)
				//args := evalArgs(o.Arguments, scope)
				//return clsObj.CallMethod(call.Call.Pos().Sline(), scope, fname, args...)
			}
		}

	default:
		switch o := call.Call.(type) {
		case *ast.Identifier: //e.g. method call like '[1,2,3].first', 'float$to_integer'
			// Check if it's a builtin type extension method, for example: "float$xxx()"
			ok := false
			objType := strings.ToLower(string(obj.Type()))
			for _, prefix := range builtinTypes {
				if strings.HasPrefix(objType, prefix) {
					ok = true
				}
			}
			if ok {
				name := fmt.Sprintf("%s$%s", objType, o.String())
				if fn, ok := scope.Get(name); ok {
					extendScope := NewScope(scope, nil)
					extendScope.Set("self", obj) // Set "self" to be the implicit object.
					results := Eval(fn.(*Function).Literal.Body, extendScope)
					if results.Type() == RETURN_VALUE_OBJ {
						return results.(*ReturnValue).Value
					}
					return results
				} else {
					return obj.CallMethod(call.Call.Pos().Sline(), scope, o.String())
				}
			} else {
				return obj.CallMethod(call.Call.Pos().Sline(), scope, o.String())
			}
		case *ast.CallExpression: //e.g. method call like '[1,2,3].first()', 'float$to_integer()'
			args := evalArgs(o.Arguments, scope)
			// Check if it's a builtin type extension method, for example: "float$xxx()"
			ok := false
			objType := strings.ToLower(string(obj.Type()))
			for _, prefix := range builtinTypes {
				if strings.HasPrefix(objType, prefix) {
					ok = true
				}
			}
			if ok {
				name := fmt.Sprintf("%s$%s", strings.ToLower(string(m.Type())), o.Function.String())
				if fn, ok := scope.Get(name); ok {
					extendScope := extendFunctionScope(fn.(*Function), args)
					extendScope.Set("self", obj) // Set "self" to be the implicit object.

					results := Eval(fn.(*Function).Literal.Body, extendScope)
					if results.Type() == RETURN_VALUE_OBJ {
						return results.(*ReturnValue).Value
					}
					return results
				} else {
					return obj.CallMethod(call.Call.Pos().Sline(), scope, o.Function.String(), args...)
				}
			} else {
				return obj.CallMethod(call.Call.Pos().Sline(), scope, o.Function.String(), args...)
			}
		}
	}

	return NewError(call.Call.Pos().Sline(), NOMETHODERROR, call.String(), obj.Type())

}

func evalArgs(args []ast.Expression, scope *Scope) []Object {
	//TODO: Refactor this to accept the params and args, go ahead and
	// update scope while looping and return the Scope object.
	e := []Object{}
	for _, v := range args {
		item := Eval(v, scope)
		e = append(e, item)
	}
	return e
}

// Index Expressions, i.e. array[0], array[2:4], tuple[3] or hash["mykey"]
func evalIndexExpression(ie *ast.IndexExpression, scope *Scope) Object {
	left := Eval(ie.Left, scope)
	switch iterable := left.(type) {
	case *Array:
		return evalArrayIndex(iterable, ie, scope)
	case *Hash:
		return evalHashKeyIndex(iterable, ie, scope)
	case *String:
		return evalStringIndex(iterable, ie, scope)
	case *Tuple:
		return evalTupleIndex(iterable, ie, scope)
	case *ObjectInstance: //class indexer's getter
		return evalClassInstanceIndexer(iterable, ie, scope)
	}
	return NewError(ie.Pos().Sline(), NOINDEXERROR, left.Type())
}

func evalClassInstanceIndexer(instanceObj *ObjectInstance, ie *ast.IndexExpression, scope *Scope) Object {
	var num int
	switch o := ie.Index.(type) {
	case *ast.ClassIndexerExpression:
		num = len(o.Parameters)
	default:
		num = 1
	}

	propName := "this" + fmt.Sprintf("%d", num)
	p := instanceObj.GetProperty(propName)
	if p != nil {
		//no getter or getter block is empty, e.g. 'property xxx { get; }'
		if p.Getter == nil || len(p.Getter.Body.Statements) == 0 {
			return NewError(ie.Pos().Sline(), INDEXERUSEERROR, instanceObj.Class.Name)
		} else {
			newScope := NewScope(instanceObj.Scope, nil)

			switch o := ie.Index.(type) {
			case *ast.ClassIndexerExpression:
				for i, v := range o.Parameters {
					index := Eval(v, scope)
					newScope.Set(p.Indexes[i].Value, index)
				}
			default:
				index := Eval(ie.Index, scope)
				newScope.Set(p.Indexes[0].Value, index)
			}

			results := Eval(p.Getter.Body, newScope)
			if results.Type() == RETURN_VALUE_OBJ {
				return results.(*ReturnValue).Value
			}
			return results
		}
	}

	return NewError(ie.Pos().Sline(), INDEXNOTFOUNDERROR, instanceObj.Class.Name)
}

func evalStringIndex(str *String, ie *ast.IndexExpression, scope *Scope) Object {
	var idx int64
	length := int64(utf8.RuneCountInString(str.String))
	//length := int64(len(str.String))
	if exp, success := ie.Index.(*ast.SliceExpression); success {
		return evalStringSliceExpression(str, exp, scope)
	}
	index := Eval(ie.Index, scope)
	if index.Type() == ERROR_OBJ {
		return index
	}

	switch o := index.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	default:
		idx = 0
		if IsTrue(index) {
			idx = 1
		}
	}

	if idx >= length || idx < 0 {
		return NewError(ie.Pos().Sline(), INDEXERROR, idx)
	}

	return NewString(string([]rune(str.String)[idx])) //support utf8,not very efficient
	//return &String{String: string(str.String[idx]), Valid:true}  //only support ASCII
}

func evalStringSliceExpression(str *String, se *ast.SliceExpression, scope *Scope) Object {
	var idx int64
	var slice int64

	length := int64(utf8.RuneCountInString(str.String))
	//length := int64(len(str.String))

	startIdx := Eval(se.StartIndex, scope)
	if startIdx.Type() == ERROR_OBJ {
		return startIdx
	}

	switch o := startIdx.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	}
	if idx >= length || idx < 0 {
		return NewError(se.Pos().Sline(), INDEXERROR, idx)
	}

	if se.EndIndex == nil {
		slice = length
	} else {
		slIndex := Eval(se.EndIndex, scope)
		if slIndex.Type() == ERROR_OBJ {
			return slIndex
		}

		switch o := slIndex.(type) {
		case *Integer:
			slice = o.Int64
		case *UInteger:
			slice = int64(o.UInt64)
		}
		if slice >= (length + 1) {
			return NewError(se.Pos().Sline(), SLICEERROR, idx, slice)
		}
		if slice < 0 {
			return NewError(se.Pos().Sline(), SLICEERROR, idx, slice)
		}
	}
	if idx == 0 && slice == length {
		return str
	}

	if slice < idx {
		return NewError(se.Pos().Sline(), SLICEERROR, idx, slice)
	}

	runes := []rune(str.String)
	return NewString(string(runes[idx:slice]))
}

func evalHashKeyIndex(hash *Hash, ie *ast.IndexExpression, scope *Scope) Object {
	var key Object
	switch ie.Index.(type) {
	case *ast.Identifier:
		//check if the identfier is in scope
		if k, ok := scope.Get(ie.Index.String()); ok {
			return hash.Get(ie.Pos().Sline(), k)
		} else {
			/* not in scope, we assume it's a string without quote.
			   for example: hash[a] ==> hash["a"]
			*/
			key = NewString(ie.Index.String())
		}
	default:
		key = Eval(ie.Index, scope)
	}

	if key.Type() == ERROR_OBJ {
		return key
	}
	return hash.Get(ie.Pos().Sline(), key)
}

func evalArraySliceExpression(array *Array, se *ast.SliceExpression, scope *Scope) Object {
	var idx int64
	var slice int64
	length := int64(len(array.Members))

	startIdx := Eval(se.StartIndex, scope)
	if startIdx.Type() == ERROR_OBJ {
		return startIdx
	}

	switch o := startIdx.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	default:
		idx = 0
		if IsTrue(o) {
			idx = 1
		}
	}

	if idx < 0 {
		return NewError(se.Pos().Sline(), INDEXERROR, idx)
	}

	if idx >= length {
		return NIL
	}

	if se.EndIndex == nil {
		slice = length
	} else {
		slIndex := Eval(se.EndIndex, scope)
		if slIndex.Type() == ERROR_OBJ {
			return slIndex
		}

		switch o := slIndex.(type) {
		case *Integer:
			slice = o.Int64
		case *UInteger:
			slice = int64(o.UInt64)
		default:
			slice = 0
			if IsTrue(o) {
				slice = 1
			}
		}
		if slice >= (length+1) || slice < 0 {
			return NewError(se.Pos().Sline(), SLICEERROR, idx, slice)
		}
	}
	if idx == 0 && slice == length {
		return array
	}

	if slice < idx {
		return NewError(se.Pos().Sline(), SLICEERROR, idx, slice)
	}

	if slice == length {
		return &Array{Members: array.Members[idx:]}
	}
	return &Array{Members: array.Members[idx:slice]}
}

func evalArrayIndex(array *Array, ie *ast.IndexExpression, scope *Scope) Object {
	var idx int64
	length := int64(len(array.Members))
	if exp, success := ie.Index.(*ast.SliceExpression); success {
		return evalArraySliceExpression(array, exp, scope)
	}
	index := Eval(ie.Index, scope)
	if index.Type() == ERROR_OBJ {
		return index
	}

	switch o := index.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	default:
		idx = 0
		if IsTrue(index) {
			idx = 1
		}
	}
	if idx < 0 {
		return NewError(ie.Pos().Sline(), INDEXERROR, idx)
	}
	if idx >= length {
		return NIL
	}
	return array.Members[idx]
}

//Almost same as evalArraySliceExpression
func evalTupleSliceExpression(tuple *Tuple, se *ast.SliceExpression, scope *Scope) Object {
	var idx int64
	var slice int64
	length := int64(len(tuple.Members))

	startIdx := Eval(se.StartIndex, scope)
	if startIdx.Type() == ERROR_OBJ {
		return startIdx
	}

	switch o := startIdx.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	}
	if idx < 0 {
		return NewError(se.Pos().Sline(), INDEXERROR, idx)
	}

	if idx >= length {
		return NIL
	}

	if se.EndIndex == nil {
		slice = length
	} else {
		slIndex := Eval(se.EndIndex, scope)
		if slIndex.Type() == ERROR_OBJ {
			return slIndex
		}

		switch o := slIndex.(type) {
		case *Integer:
			slice = o.Int64
		case *UInteger:
			slice = int64(o.UInt64)
		}
		if slice >= (length+1) || slice < 0 {
			return NewError(se.Pos().Sline(), SLICEERROR, idx, slice)
		}
	}
	if idx == 0 && slice == length {
		return tuple
	}

	if slice < idx {
		return NewError(se.Pos().Sline(), SLICEERROR, idx, slice)
	}

	if slice == length {
		return &Tuple{Members: tuple.Members[idx:]}
	}
	return &Tuple{Members: tuple.Members[idx:slice]}
}

//Almost same as evalArrayIndex
func evalTupleIndex(tuple *Tuple, ie *ast.IndexExpression, scope *Scope) Object {
	var idx int64
	length := int64(len(tuple.Members))
	if exp, success := ie.Index.(*ast.SliceExpression); success {
		return evalTupleSliceExpression(tuple, exp, scope)
	}
	index := Eval(ie.Index, scope)
	if index.Type() == ERROR_OBJ {
		return index
	}

	switch o := index.(type) {
	case *Integer:
		idx = o.Int64
	case *UInteger:
		idx = int64(o.UInt64)
	default:
		idx = 0
		if IsTrue(index) {
			idx = 1
		}
	}

	if idx < 0 {
		return NewError(ie.Pos().Sline(), INDEXERROR, idx)
	}
	if idx >= length {
		return NIL
	}
	return tuple.Members[idx]
}

func evalPostfixExpression(left Object, node *ast.PostfixExpression, scope *Scope) Object {
	if left.Type() == INSTANCE_OBJ { //operator overloading
		instanceObj := left.(*ObjectInstance)
		method := instanceObj.GetMethod(node.Operator)
		if method != nil {
			switch method.(type) {
			case *Function:
				newScope := NewScope(instanceObj.Scope, nil)
				args := []Object{left}
				return evalFunctionDirect(method, args, instanceObj, newScope, nil)
			case *BuiltinMethod:
				//do nothing for now
			}
		}
		return NewError(node.Pos().Sline(), POSTFIXOP, node.Operator, left.Type())
	}

	switch node.Operator {
	case "++":
		return evalIncrementPostfixOperatorExpression(node, left, scope)
	case "--":
		return evalDecrementPostfixOperatorExpression(node, left, scope)
	default:
		return NewError(node.Pos().Sline(), POSTFIXOP, node.Operator, left.Type())
	}
}

func evalIncrementPostfixOperatorExpression(node *ast.PostfixExpression, left Object, scope *Scope) Object {
	switch left.Type() {
	case INTEGER_OBJ:
		leftObj := left.(*Integer)
		returnVal := NewInteger(leftObj.Int64)
		scope.Reset(node.Left.String(), NewInteger(leftObj.Int64+1))
		return returnVal
	case UINTEGER_OBJ:
		leftObj := left.(*UInteger)
		returnVal := NewUInteger(leftObj.UInt64)
		scope.Reset(node.Left.String(), NewUInteger(leftObj.UInt64+1))
		return returnVal
	case FLOAT_OBJ:
		leftObj := left.(*Float)
		returnVal := NewFloat(leftObj.Float64)
		scope.Reset(node.Left.String(), NewFloat(leftObj.Float64+1))
		return returnVal
	default:
		return NewError(node.Pos().Sline(), POSTFIXOP, node.Operator, left.Type())
	}
}

func evalDecrementPostfixOperatorExpression(node *ast.PostfixExpression, left Object, scope *Scope) Object {
	switch left.Type() {
	case INTEGER_OBJ:
		leftObj := left.(*Integer)
		returnVal := NewInteger(leftObj.Int64)
		scope.Reset(node.Left.String(), NewInteger(leftObj.Int64-1))
		return returnVal
	case UINTEGER_OBJ:
		leftObj := left.(*UInteger)
		returnVal := NewUInteger(leftObj.UInt64)
		scope.Reset(node.Left.String(), NewUInteger(leftObj.UInt64-1))
		return returnVal
	case FLOAT_OBJ:
		leftObj := left.(*Float)
		returnVal := NewFloat(leftObj.Float64)
		scope.Reset(node.Left.String(), NewFloat(leftObj.Float64-1))
		return returnVal
	default:
		return NewError(node.Pos().Sline(), POSTFIXOP, node.Operator, left.Type())
	}
}

func evalTryStatement(tryStmt *ast.TryStmt, scope *Scope) Object {
	rv := Eval(tryStmt.Try, scope)
	if rv.Type() == ERROR_OBJ {
		return rv
	}

	throwNotHandled := false
	var throwObj Object = NIL
	if rv.Type() == THROW_OBJ {
		throwObj = rv.(*Throw)
		if tryStmt.Catch != nil {
			catchScope := NewScope(scope, scope.Writer)
			if tryStmt.Var != "" {
				catchScope.Set(tryStmt.Var, rv.(*Throw).value)
			}
			rv = evalBlockStatements(tryStmt.Catch.Statements, catchScope) //catch Block
			if rv.Type() == ERROR_OBJ {
				return rv
			}
		} else {
			throwNotHandled = true
		}
	}

	if tryStmt.Finally != nil { //finally will always run
		rv = evalBlockStatements(tryStmt.Finally.Statements, scope)
		if rv.Type() == ERROR_OBJ {
			return rv
		}
	}

	if throwNotHandled {
		return throwObj
	}
	return NIL
}
//Evaluate ternary expression
func evalTernaryExpression(te *ast.TernaryExpression, scope *Scope) Object {
	condition := Eval(te.Condition, scope) //eval condition
	if condition.Type() == ERROR_OBJ {
		return condition
	}

	if IsTrue(condition) {
		return Eval(te.IfTrue, scope)
	} else {
		return Eval(te.IfFalse, scope)
	}
}

func evalSpawnStatement(s *ast.SpawnStmt, scope *Scope) Object {
	newSpawnScope := NewScope(scope, nil)

	switch callExp := s.Call.(type) {
	case *ast.CallExpression:
		go (func() {
			evalFunctionCall(callExp, newSpawnScope)
		})()
	case *ast.MethodCallExpression:
		go (func() {
			evalMethodCallExpression(callExp, newSpawnScope)
		})()
	default:
		return NewError(s.Pos().Sline(), SPAWNERROR)
	}

	return NIL
}

func evalPipeExpression(p *ast.Pipe, scope *Scope) Object {
	left := Eval(p.Left, scope)

	// Convert the type object back to an expression
	// so it can be passed to the FunctionCall arguments.
	argument := obj2Expression(left)
	if argument == nil {
		return NIL
	}

	// The right side operator should be a function.
	switch rightFunc := p.Right.(type) {
	case *ast.MethodCallExpression:
		// Prepend the left-hand interpreted value
		// to the function arguments.
		switch rightFunc.Call.(type) {
		case *ast.Identifier:
			//e.g.
			//x = ["hello", "world"] |> strings.upper    : rightFunc.Call.(type) == *ast.Identifier
			//x = ["hello", "world"] |> strings.upper()  : rightFunc.Call.(type) == *ast.CallExpression
			//so here we convert *ast.Identifier to * ast.CallExpression
			rightFunc.Call = &ast.CallExpression{Token: p.Token, Function: rightFunc.Call}
		}
		rightFunc.Call.(*ast.CallExpression).Arguments = append([]ast.Expression{argument}, rightFunc.Call.(*ast.CallExpression).Arguments...)
		return Eval(rightFunc, scope)
	case *ast.CallExpression:
		rightFunc.Arguments = append([]ast.Expression{argument}, rightFunc.Arguments...)
		return Eval(rightFunc, scope)

	}

	return NIL
}

//class name : parent { block }
//class name (categoryname) { block }
func evalClassStatement(c *ast.ClassStatement, scope *Scope) Object {
	if c.CategoryName != nil { //it's a class category
		clsObj, ok := scope.Get(c.Name.Value)
		if !ok {
			return NewError(c.Pos().Sline(), CLASSCATEGORYERROR, c.Name, c.CategoryName)
		}

		//category only support methods and properties
		cls := clsObj.(*Class)
		for k, f := range c.ClassLiteral.Methods { //f :function
			cls.Methods[k] = Eval(f, scope).(ClassMethod)
		}
		for k, p := range c.ClassLiteral.Properties { //p :property
			cls.Properties[k] = p
		}

		return NIL
	}

	var clsObj Object
	if c.IsAnnotation {
		clsObj = evalClassLiterlForAnno(c.ClassLiteral, scope)
	} else {
		clsObj = evalClassLiteral(c.ClassLiteral, scope)
	}

	scope.Set(c.Name.Value, clsObj) //save to scope

	return NIL
}

//let name = class : parent { block }
func evalClassLiteral(c *ast.ClassLiteral, scope *Scope) Object {
	var parentClass = BASE_CLASS //base class is the root of all classes in magpie
	if c.Parent != "" {

		parent, ok := scope.Get(c.Parent)
		if !ok {
			return NewError(c.Pos().Sline(), PARENTNOTDECL, c.Parent)
		}

		parentClass, ok = parent.(*Class)
		if !ok {
			return NewError(c.Pos().Sline(), NOTCLASSERROR, c.Parent)
		}
	}

	clsObj := &Class{
		Name:       c.Name,
		Parent:     parentClass,
		Members:    c.Members,
		Properties: c.Properties,
		Methods:    make(map[string]ClassMethod, len(c.Methods)),
	}

	tmpClass := clsObj
	classChain := make([]*Class, 0, 3)
	classChain = append(classChain, clsObj)
	for tmpClass.Parent != nil {
		classChain = append(classChain, tmpClass.Parent)
		tmpClass = tmpClass.Parent
	}

	//create a new Class scope
	newScope := NewScope(scope, nil)
	//evaluate the 'Members' fields of class with proper scope.
	for idx := len(classChain) - 1; idx >= 0; idx-- {
		for _, member := range classChain[idx].Members {
			Eval(member, newScope) //evaluate the 'Members' fields of class
		}
		newScope = NewScope(newScope, nil)
	}
	clsObj.Scope = newScope.parentScope
	clsObj.Scope.Set("this", clsObj) //make 'this' refer to class object itself
	clsObj.Scope.Set("parent", parentClass)

	for k, f := range c.Methods {
		clsObj.Methods[k] = Eval(f, scope).(ClassMethod)
	}

	//check if the method has @Override annotation, if so, search
	//the method in parent hierarchical, if not found, then return error.
	for methodName, fnStmt := range c.Methods {
		for _, anno := range fnStmt.Annotations {
			if anno.Name.Value == OVERRIDE_ANNOCLASS.Name {
				if clsObj.Parent.GetMethod(methodName) == nil {
					return NewError(fnStmt.FunctionLiteral.Pos().Sline(), OVERRIDEERROR, methodName, clsObj.Name)
				}
			}
		}
	}

	return clsObj
}

func evalClassLiterlForAnno(c *ast.ClassLiteral, scope *Scope) Object {
	var parentClass = BASE_CLASS //base class is the root of all classes in magpie
	if c.Parent != "" {
		parent, ok := scope.Get(c.Parent)
		if !ok {
			return NewError(c.Pos().Sline(), PARENTNOTDECL, c.Parent)
		}

		parentClass, ok = parent.(*Class)
		if !ok {
			return NewError(c.Pos().Sline(), NOTCLASSERROR, c.Parent)
		}
	}

	if parentClass != BASE_CLASS && !parentClass.IsAnnotation { //parent not annotation
		return NewError(c.Pos().Sline(), PARENTNOTANNOTATION, c.Name, parentClass.Name)
	}

	clsObj := &Class{
		Name:         c.Name,
		Parent:       parentClass,
		Properties:   c.Properties,
		IsAnnotation: true,
	}

	//create a new Class scope
	clsObj.Scope = NewScope(scope, nil)
	clsObj.Scope.Set("this", clsObj) //make 'this' refer to class object itself
	clsObj.Scope.Set("parent", parentClass)

	return clsObj
}

//new classname(parameters)
func evalNewExpression(n *ast.NewExpression, scope *Scope) Object {
	class := Eval(n.Class, scope)
	if class == NIL || class == nil {
		return NewError(n.Pos().Sline(), CLSNOTDEFINE, n.Class)
	}

	clsObj, ok := class.(*Class)
	if !ok {
		return NewError(n.Pos().Sline(), NOTCLASSERROR, n.Class)
	}

	tmpClass := clsObj
	classChain := make([]*Class, 0, 3)
	classChain = append(classChain, clsObj)
	for tmpClass.Parent != nil {
		classChain = append(classChain, tmpClass.Parent)
		tmpClass = tmpClass.Parent
	}

	//create a new Class scope
	newScope := NewScope(scope, nil)
	//evaluate the 'Members' fields of class with proper scope.
	for idx := len(classChain) - 1; idx >= 0; idx-- {
		for _, member := range classChain[idx].Members {
			if !member.StaticFlag {
				Eval(member, newScope) //evaluate the 'Members' fields of class
			}
		}
		newScope = NewScope(newScope, nil)
	}

	instance := &ObjectInstance{Class: clsObj, Scope: newScope.parentScope}
	instance.Scope.Set("this", instance)        //make 'this' refer to instance
	instance.Scope.Set("parent", classChain[1]) //make 'parent' refer to instance's parent

	//Is it has a constructor ?
	init := clsObj.GetMethod("init")
	if init == nil {
		return instance
	}

	args := evalArgs(n.Arguments, scope)
	if len(args) == 1 && args[0].Type() == ERROR_OBJ {
		return args[0]
	}

	ret := evalFunctionDirect(init, args, instance, instance.Scope, nil)
	if ret.Type() == ERROR_OBJ {
		return ret //return the error object
	}
	return instance
}

func processClassAnnotation(Annotations []*ast.AnnotationStmt, scope *Scope, line string, obj Object) {
	for _, anno := range Annotations { //for each annotation
		annoClass, ok := scope.Get(anno.Name.Value)
		if !ok {
			panic(NewError(line, CLSNOTDEFINE, anno.Name.Value))
		}

		annoClsObj := annoClass.(*Class)

		//create the annotation instance
		newScope := NewScope(scope, nil)
		annoInstanceObj := &ObjectInstance{Class: annoClsObj, Scope: newScope}
		annoInstanceObj.Scope.Set("this", annoInstanceObj) //make 'this' refer to annoObj

		switch o := obj.(type) {
		case *Function:
			o.Annotations = append(o.Annotations, annoInstanceObj)
		case *Array:
			o.Members = append(o.Members, annoInstanceObj)
		}

		defaultPropMap := make(map[string]ast.Expression)
		//get all propertis which have default value in the annotation class
		tmpCls := annoClsObj
		for tmpCls != nil {
			for name, item := range tmpCls.Properties {
				if item.Default != nil {
					defaultPropMap[name] = item.Default
				}
			}

			tmpCls = tmpCls.Parent
		}

		//check if the property(which has default value) exists in anno.Attribues
		for name, item := range defaultPropMap {
			if _, ok := anno.Attributes[name]; !ok { //not exists
				anno.Attributes[name] = item
			}
		}

		for k, v := range anno.Attributes { //for each annotation attribute
			val := Eval(v, annoInstanceObj.Scope)
			p := annoClsObj.GetProperty(k)
			if p == nil {
				annoInstanceObj.Scope.Set(k, val)
			} else {
				annoInstanceObj.Scope.Set("_"+k, val)
			}
		}
	}
}

func evalFunctionDirect(fn Object, args []Object, instance *ObjectInstance, scope *Scope, call *ast.CallExpression) Object {
	switch fn := fn.(type) {
	case *Function:
		fn.Instance = instance
		//		if len(args) < len(fn.Literal.Parameters) {
		//			return NewError("", GENERICERROR, "Not enough parameters to call function")
		//		}

		newScope := NewScope(scope, nil)
		variadicParam := []Object{}
		for i := range args {
			//Because of function default values, we need to check `i >= len(args)`
			if fn.Variadic && i >= len(fn.Literal.Parameters)-1 {
				for j := i; j < len(args); j++ {
					variadicParam = append(variadicParam, args[j])
				}
				break
			} else if i >= len(fn.Literal.Parameters) {
				break
			} else {
				newScope.Set(fn.Literal.Parameters[i].String(), args[i])
			}
		}

		// Variadic argument is passed as a single array
		// of parameters.
		if fn.Variadic {
			newScope.Set(fn.Literal.Parameters[len(fn.Literal.Parameters)-1].String(), &Array{Members: variadicParam})
			if len(args) < len(fn.Literal.Parameters) {
				newScope.Set("@_", NewInteger(int64(len(fn.Literal.Parameters)-1)))
			} else {
				newScope.Set("@_", NewInteger(int64(len(args))))
			}
		} else {
			newScope.Set("@_", NewInteger(int64(len(fn.Literal.Parameters))))
		}

		if fn.Async && call.Awaited {
			aChan := make(chan Object, 1)

			go func() {
				defer close(aChan)

				results := Eval(fn.Literal.Body, newScope)
				if obj, ok := results.(*ReturnValue); ok {
					// if function returns multiple-values
					// returns a tuple instead.
					if len(obj.Values) > 1 {
						results = &Tuple{Members: obj.Values, IsMulti: true}
					} else {
						results = obj.Value
					}
				}

				aChan <- results
			}()

			return <-aChan
		}

		if fn.Async {
			go Eval(fn.Literal.Body, newScope)
			return NIL
		}

		//newScope.DebugPrint("    ") //debug
		results := Eval(fn.Literal.Body, newScope)
		if obj, ok := results.(*ReturnValue); ok {
			// if function returns multiple-values
			// returns a tuple instead.
			if len(obj.Values) > 1 {
				return &Tuple{Members: obj.Values, IsMulti: true}
			}
			return obj.Value
		}

		return results
	case *Builtin:
		return fn.Fn("", scope, args...)
	case *BuiltinMethod:
		return fn.Fn("", fn.Instance, scope, args...)
	}

	return NewError("", GENERICERROR, fn.Type()+" is not a function")
}

//evaluate 'using' statement
func evalUsingStatement(u *ast.UsingStmt, scope *Scope) Object {
	//evaluate the assignment expression
	obj := evalAssignExpression(u.Expr, scope)
	if obj.Type() == ERROR_OBJ {
		return obj
	}

	fn := func() {
		if obj.Type() != NIL_OBJ {
			// Check if val is 'Closeable'
			if c, ok := obj.(Closeable); ok {
				//call the 'Close' method of the object
				c.close(u.Pos().Sline())
			}
		}
	}
	defer func() {
		if r := recover(); r != nil { // if there is panic, we need to call fn()
			fn()
		} else { //no panic, we also need to call fn()
			fn()
		}
	}()

	//evaluate the 'using' block statement
	Eval(u.Block, scope)

	return NIL
}

//code copied from https://github.com/abs-lang/abs with modification for windowns command.
func evalCmdExpression(t *ast.CmdExpression, scope *Scope) Object {
	cmd := t.Value
	// Match all strings preceded by
	// a $ or a \$
	re := regexp.MustCompile("(\\\\)?\\$([a-zA-Z_]{1,})")
	cmd = re.ReplaceAllStringFunc(cmd, func(m string) string {
		// If the string starts with a backslash,
		// that's an escape, so we should replace
		// it with the remaining portion of the match.
		// \$VAR becomes $VAR
		if string(m[0]) == "\\" {
			return m[1:]
		}

		// If the string starts with $, then
		// it's an interpolation. Let's
		// replace $VAR with the variable
		// named VAR in the ABS' environment.
		// If the variable is not found, we
		// just dump an empty string
		v, ok := scope.Get(m[1:])

		if !ok {
			return ""
		}

		return v.Inspect()
	})

	var commands []string
	var executor string
	if runtime.GOOS == "windows" {
		commands = []string{"/C", cmd}
		executor = "cmd.exe"
	} else {
		commands = []string{"-c", cmd}
		executor = "bash"
	}

	c := exec.Command(executor, commands...)
	c.Env = os.Environ()
	var out bytes.Buffer
	var stderr bytes.Buffer
	c.Stdin = os.Stdin
	c.Stdout = &out
	c.Stderr = &stderr
	err := c.Run()

	if err != nil {
		return &String{String: stderr.String(), Valid: false}
	}

	return &String{String: strings.Trim(out.String(), "\n"), Valid: true}
}

//========================================================
//               LINQ EVALUATION LOGIC(BEGIN)
//========================================================
// Evaluate linq query expression
func evalLinqQueryExpression(query *ast.QueryExpr, scope *Scope) Object {
	fromExpr := query.From.(*ast.FromExpr)
	queryBodyExpr := query.QueryBody.(*ast.QueryBodyExpr)

	innerScope := NewScope(scope, nil)

	inValue := Eval(fromExpr.Expr, innerScope)
	if inValue.Type() == ERROR_OBJ {
		return inValue
	}

	line := query.Pos().Sline()

	lq := &LinqObj{}
	//=================================================
	// query_expression : from_clause query_body
	//=================================================
	//from_clause : FROM identifier IN expression
	fromObj := lq.FromQuery(line, innerScope, inValue, NewString(fromExpr.Var)).(*LinqObj)

	//query_body : query_body_clause* select_or_group_clause query_continuation?
	tmpLinq := fromObj

	//query_body_clause*
	for _, queryBody := range queryBodyExpr.QueryBody {
		queryBodyExpr := queryBody.(*ast.QueryBodyClauseExpr)

		switch clause := queryBodyExpr.Expr.(type) {
		case *ast.AssignExpression: //let-clause
			assignExp := clause

			fl := constructFuncLiteral("", assignExp.Value, token.ASSIGN, assignExp.Pos())
			fnObj := evalFunctionLiteral(fl, innerScope)
			letVar := assignExp.Name.(*ast.Identifier).Value
			tmpLinq = tmpLinq.Let(line, innerScope, fnObj, NewString(letVar)).(*LinqObj)

		case *ast.FromExpr: // from_clause : FROM identifier IN expression
			innerFrom := clause
			tmpLinq = tmpLinq.FromInner(line, innerScope, NewString(innerFrom.Var)).(*LinqObj)

		case *ast.WhereExpr: // where_clause : WHERE expression
			whereExp := clause
			fl := constructFuncLiteral("", whereExp.Expr, token.WHERE, whereExp.Pos())

			fnObj := evalFunctionLiteral(fl, innerScope)
			tmpLinq = tmpLinq.Where2(line, innerScope, fnObj).(*LinqObj)

		case *ast.JoinExpr:
			fmt.Fprintf(scope.Writer, "JOIN: [NOT IMPLEMENTED]\n")

		case *ast.OrderExpr: // orderby_clause : ORDERBY ordering (','  ordering)*
			orderExpr := clause
			var i int = 0
			var str string
			for _, orderingExpr := range orderExpr.Ordering {
				order := orderingExpr.(*ast.OrderingExpr)

				str = order.Var
				fl := constructFuncLiteral(str, order.Expr, token.ORDERBY, order.Pos())
				fnObj := evalFunctionLiteral(fl, innerScope)

				if order.IsAscending {
					if i > 0 {
						tmpLinq = tmpLinq.ThenBy(line, innerScope, fnObj).(*LinqObj)
					} else {
						tmpLinq = tmpLinq.OrderBy(line, innerScope, fnObj).(*LinqObj)
					}
				} else {
					if i > 0 {
						tmpLinq = tmpLinq.ThenByDescending(line, innerScope, fnObj).(*LinqObj)
					} else {
						tmpLinq = tmpLinq.OrderByDescending(line, innerScope, fnObj).(*LinqObj)
					}
				}
				i++
			}

			//Note Here: we must convert the LinqObj to an array, then from
			//           the converted array, construct a new LinqObj.
			//
			// Is there a better way for doing this???
			arr := tmpLinq.ToOrderedSlice(line).(*Array)
			tmpLinq = tmpLinq.FromQuery(line, innerScope, arr, NewString(str)).(*LinqObj)

		default:
			fmt.Fprintf(scope.Writer, "[NOT IMPLEMENTED]\n")
		}
	}

	//select_or_group_clause
	switch queryBodyExpr.Expr.(type) {
	case *ast.SelectExpr:
		/*
			let selectArr = [1,2,3,4,5,6,7,8,9,10]
			[NORMAL]:
				result = linq.from(selectArr).select(fn(x) {
					x = x + 2
				})
			<=>
			[LINQ]:
				result = from x in selectArr select x + 2
		*/
		selectExp := queryBodyExpr.Expr.(*ast.SelectExpr)
		fl := constructFuncLiteral("", selectExp.Expr, token.SELECT, selectExp.Pos())
		fnObj := evalFunctionLiteral(fl, innerScope)
		return tmpLinq.Select2(line, innerScope, fnObj)

	case *ast.GroupExpr:
		/*
			let groupByArr = [1, 2, 3, 4, 5, 6, 7, 8, 9]
			[NORMAL]:
				result = linq.from(groupByArr).groupBy(
					fn(v) { return v % 2 == 0 },
					fn(v) { return v }
				)
			<=>
			[LINQ]:
				result = from v in groupByArr group v BY v % 2 == 0
		*/
		groupExp := queryBodyExpr.Expr.(*ast.GroupExpr)

		keyFuncLiteral := constructFuncLiteral("", groupExp.ByExpr, token.BY, groupExp.Pos())
		elementFuncLiteral := constructFuncLiteral("", groupExp.GrpExpr, token.GROUP, groupExp.Pos())
		keySelector := evalFunctionLiteral(keyFuncLiteral, innerScope)
		elementSelector := evalFunctionLiteral(elementFuncLiteral, innerScope)
		return tmpLinq.GroupBy2(line, innerScope, keySelector, elementSelector)
	}

	return NIL
}

//construct a FunctionLiteral for use with linq.xxx() function
func constructFuncLiteral(value string, expr ast.Expression, tokenType token.TokenType, pos token.Position) *ast.FunctionLiteral {
	fl := &ast.FunctionLiteral{Parameters: []ast.Expression{}}
	if len(value) > 0 {
		fl.Parameters = append(fl.Parameters, &ast.Identifier{Value: value})
	}

	tok := token.Token{Type: tokenType, Pos: pos}
	fl.Body = &ast.BlockStatement{Token: tok, Statements: []ast.Statement{&ast.ExpressionStatement{Token: tok, Expression: expr}}}

	return fl
}

//========================================================
//               LINQ EVALUATION LOGIC(END)
//========================================================

func evalAwaitExpression(a *ast.AwaitExpr, scope *Scope) Object {
	switch fn := a.Call.(type) {
	case *ast.CallExpression:
		fn.Awaited = true
		return evalFunctionCall(fn, scope)
	case *ast.MethodCallExpression:
		fn.Call.(*ast.CallExpression).Awaited = true
		return evalMethodCallExpression(fn, scope)
	default:
		//should never reach this line, because the parser will check type call type
		return NIL
	}
}

func evalServiceStatement(s *ast.ServiceStatement, scope *Scope) Object {
	//note: map's value is not important
	var routeMap = map[string]bool{
		"url":     true,
		"methods": true,
		"host":    true,
		"schemes": true,
		"headers": true,
		"queries": true,
	}

	svcObj := NewService(s.Addr).(*ServiceObj)

	for _, fnStmt := range s.Methods {
		f := evalFunctionStatement(fnStmt, scope).(*Function)
		var methodArr *Array
		var host *String
		var schemes *String
		var headers *Hash
		var queries *Hash
		var hasUrl bool
		anno := fnStmt.Annotations[0]
		for k, v := range anno.Attributes { //for each annotation attribute
			if _, ok := routeMap[k]; ok {
				if k == "url" {
					hasUrl = true
					val := Eval(v, scope).(*String)
					//fmt.Printf("key=%s, val=%s, val.Type=%s\n", k, val.Inspect(), val.Type())
					svcObj.HandleFunc(s.Pos().Sline(), scope, val, f)
				} else if k == "methods" {
					methodArr = Eval(v, scope).(*Array)
				} else if k == "host" {
					host = Eval(v, scope).(*String)
				} else if k == "schemes" {
					schemes = Eval(v, scope).(*String)
				} else if k == "headers" {
					headers = Eval(v, scope).(*Hash)
				} else if k == "queries" {
					queries = Eval(v, scope).(*Hash)
				}
			} else {
				continue
			}
		}

		if !hasUrl {
			return NewError(s.Pos().Sline(), SERVICENOURLERROR, s.Name.Value, fnStmt.Name.Value)
		}

		if methodArr != nil {
			svcObj.Methods(s.Pos().Sline(), methodArr.Members...)
		}

		if host != nil {
			svcObj.Host(s.Pos().Sline(), host)
		}

		if schemes != nil {
			svcObj.Schemes(s.Pos().Sline(), host)
		}

		if headers != nil {
			svcObj.Headers(s.Pos().Sline(), headers)
		}

		if queries != nil {
			svcObj.Queries(s.Pos().Sline(), queries)
		}
	}

	fmt.Fprintf(scope.Writer, ServiceHint, svcObj.Addr)
	svcObj.Run(s.Pos().Sline(), NewBooleanObj(s.Debug))
	return NIL
}

//private method for evalate 'a..b' expression, and returns an array object
func evalRangeExpression(node ast.Node, startIdx Object, endIdx Object, scope *Scope) Object {
	arr := &Array{}
	switch startIdx.(type) {
	case *Integer:
		startVal := startIdx.(*Integer).Int64

		var endVal int64
		switch o := endIdx.(type) {
		case *Integer:
			endVal = o.Int64
		case *UInteger:
			endVal = int64(o.UInt64)
		default:
			return NewError(node.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ+"|"+UINTEGER_OBJ, endIdx.Type())
		}

		var j int64
		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewInteger(j))
			}
		}
	case *UInteger:
		startVal := startIdx.(*UInteger).UInt64

		var endVal uint64
		switch o := endIdx.(type) {
		case *Integer:
			endVal = uint64(o.Int64)
		case *UInteger:
			endVal = o.UInt64
		default:
			return NewError(node.Pos().Sline(), RANGETYPEERROR, INTEGER_OBJ+"|"+UINTEGER_OBJ, endIdx.Type())
		}

		var j uint64
		if startVal >= endVal {
			for j = startVal; j >= endVal; j = j - 1 {
				arr.Members = append(arr.Members, NewUInteger(j))
			}
		} else {
			for j = startVal; j <= endVal; j = j + 1 {
				arr.Members = append(arr.Members, NewUInteger(j))
			}
		}
	case *String:
		startVal := startIdx.(*String).String
		if endIdx.Type() != STRING_OBJ {
			return NewError(node.Pos().Sline(), RANGETYPEERROR, STRING_OBJ, endIdx.Type())
		}
		endVal := endIdx.(*String).String

		//only support single character with lowercase
		alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"

		//convert to int for easy comparation
		leftByte := []int32(strings.ToLower(startVal))[0]
		rightByte := []int32(strings.ToLower(endVal))[0]
		if leftByte >= rightByte { // z -> a
			for i := len(alphabet) - 1; i >= 0; i-- {
				v := int32(alphabet[i])
				if v <= leftByte && v >= rightByte {
					arr.Members = append(arr.Members, NewString(string(v)))
				}
			}
		} else { // a -> z
			for _, v := range alphabet {
				if v >= leftByte && v <= rightByte {
					arr.Members = append(arr.Members, NewString(string(v)))
				}
			}
		}
	}

	return arr
}

func evalDiamondExpr(d *ast.DiamondExpr, scope *Scope) Object {
	var obj Object
	var ok bool
	if obj, ok = GetGlobalObj(d.Value); !ok {
		obj, ok = scope.Get(d.Value)
		if !ok {
			if obj, ok = importScope.Get(d.Value); !ok {
				return reportTypoSuggestions(d.Pos().Sline(), scope, d.Value)
			}
		}
	}

	if fileObj, ok := obj.(*FileObject); ok {
		return fileObj.ReadLine(d.Pos().Sline())
	} else {
		return NewError(d.Pos().Sline(), DIAMONDOPERERROR, obj.Type())
	}

}

// Convert a Object to an ast.Expression.
func obj2Expression(obj Object) ast.Expression {
	switch value := obj.(type) {
	case *Boolean:
		return &ast.Boolean{Value: value.Bool}
	case *Integer:
		return &ast.IntegerLiteral{Value: value.Int64}
	case *UInteger:
		return &ast.UIntegerLiteral{Value: value.UInt64}
	case *Float:
		return &ast.FloatLiteral{Value: value.Float64}
	case *String:
		return &ast.StringLiteral{Value: value.String}
	case *Nil:
		return &ast.NilLiteral{}
	case *Array:
		array := &ast.ArrayLiteral{}
		for _, v := range value.Members {
			result := obj2Expression(v)
			if result == nil {
				return nil
			}
			array.Members = append(array.Members, result)
		}
		return array
	case *Tuple:
		tuple := &ast.TupleLiteral{}
		for _, v := range value.Members {
			result := obj2Expression(v)
			if result == nil {
				return nil
			}
			tuple.Members = append(tuple.Members, result)
		}
		return tuple
	case *Hash:
		hash := &ast.HashLiteral{}
		hash.Pairs = make(map[ast.Expression]ast.Expression)

		for _, hk := range value.Order { //hk:hash key
			v, _ := value.Pairs[hk]
			key := &ast.StringLiteral{Value: v.Key.Inspect()}
			result := obj2Expression(v.Value)
			if result == nil {
				return nil
			}
			hash.Pairs[key] = result
		}
		return hash
	}

	return nil
}

//Returns true when lhsV and rhsV is same value.
func equal(isWholeMatch bool, lhsV, rhsV Object) bool {
	if lhsV == nil && rhsV == nil {
		return true
	}
	if (lhsV == nil && rhsV != nil) || (lhsV != nil && rhsV == nil) {
		return false
	}

	if lhsV.Type() != rhsV.Type() {
		return false
	}

	if lhsV.Type() == NIL_OBJ {
		if rhsV.Type() == NIL_OBJ {
			return true
		}
	}

	if isGoObj(lhsV) || isGoObj(rhsV) { // if it's GoObject
		return compareGoObj(lhsV, rhsV)
	}

	if lhsV.Type() == STRING_OBJ && rhsV.Type() == STRING_OBJ {
		leftStr := lhsV.(*String).String
		rightStr := rhsV.(*String).String

		if isWholeMatch {
			r := reflect.DeepEqual(lhsV, rhsV)
			if r {
				return true
			} else {
				return false
			}
		} else {
			matched, _ := regexp.MatchString(rightStr, leftStr)
			return matched
		}
	} else if lhsV.Type() == TIME_OBJ && rhsV.Type() == TIME_OBJ {
		return lhsV.(*TimeObj).Tm.Equal(rhsV.(*TimeObj).Tm)
	} else {
		r := reflect.DeepEqual(lhsV, rhsV)
		if r {
			return true
		} else {
			return false
		}
	}
}

func isTryError(o Object) bool {
	if o.Type() == ERROR_OBJ ||
		(o.Type() == NIL_OBJ && o.(*Nil).OptionalMsg != "") ||
		(o.Type() == BOOLEAN_OBJ && o.(*Boolean).Bool == false && o.(*Boolean).OptionalMsg != "") {
		return true
	}
	return false
}

func isGoObj(o Object) bool {
	return o.Type() == GO_OBJ
}

func compareGoObj(left, right Object) bool {
	if left.Type() == GO_OBJ || right.Type() == GO_OBJ {
		var goObj *GoObject
		var another Object
		if left.Type() == GO_OBJ {
			goObj = left.(*GoObject)
			another = right
		} else {
			goObj = right.(*GoObject)
			another = left
		}

		return goObj.Equal(another)
	}

	//left and right both are GoObject
	return left.(*GoObject).Equal(right)
}

// user typo probs, add better error message here, 'Did you mean... ___?'
// Used for reporting IDENTIFIER not found.
func reportTypoSuggestions(line string, scope *Scope, miss string) Object {
	keys := scope.GetKeys()
	found := TypoSuggestions(keys, miss)
	if len(found) != 0 { //found suggestions
		return NewError(line, UNKNOWNIDENTEX, miss, strings.Join(found, ", "))
	} else {
		return NewError(line, UNKNOWNIDENT, miss)
	}
}

// user typo probs, add better error message here, 'Did you mean... ___?'
// Used for reporting METHOD not found.
func reportTypoSuggestionsMeth(line string, scope *Scope, objName string, miss string) Object {
	keys := scope.GetKeys()
	found := TypoSuggestions(keys, miss)
	if len(found) != 0 { //found suggestions
		return NewError(line, NOMETHODERROREX, miss, objName, strings.Join(found, ", "))
	} else {
		return NewError(line, NOMETHODERROR, miss, objName)
	}
}

func extendFunctionScope(fn *Function, args []Object) *Scope {
	fl := fn.Literal
	scope := NewScope(fn.Scope, nil)

	// Set the defaults
	for k, v := range fl.Values {
		scope.Set(k, Eval(v, scope))
	}
	for idx, param := range fl.Parameters {
		if idx < len(args) { // default parameters must be in the last part.
			scope.Set(param.(*ast.Identifier).Value, args[idx])
		}
	}

	return scope
}
