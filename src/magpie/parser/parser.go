package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"magpie/ast"
	"magpie/lexer"
	"magpie/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var numMap = map[rune]rune{
	'𝟎': '0', '𝟘': '0', '𝟢': '0', '𝟬': '0', '𝟶': '0', '０': '0',
	'𝟏': '1', '𝟙': '1', '𝟣': '1', '𝟭': '1', '𝟷': '1', '１': '1',
	'𝟐': '2', '𝟚': '2', '𝟤': '2', '𝟮': '2', '𝟸': '2', '２': '2',
	'𝟑': '3', '𝟛': '3', '𝟥': '3', '𝟯': '3', '𝟹': '3', '３': '3',
	'𝟒': '4', '𝟜': '4', '𝟦': '4', '𝟰': '4', '𝟺': '4', '４': '4',
	'𝟓': '5', '𝟝': '5', '𝟧': '5', '𝟱': '5', '𝟻': '5', '５': '5',
	'𝟔': '6', '𝟞': '6', '𝟨': '6', '𝟲': '6', '𝟼': '6', '６': '6',
	'𝟕': '7', '𝟟': '7', '𝟩': '7', '𝟳': '7', '𝟽': '7', '７': '7',
	'𝟖': '8', '𝟠': '8', '𝟪': '8', '𝟴': '8', '𝟾': '8', '８': '8',
	'𝟗': '9', '𝟡': '9', '𝟫': '9', '𝟵': '9', '𝟿': '9', '９': '9',
}

const (
	_ int = iota
	LOWEST
	PIPE
	ASSIGN
	FATARROW
	CONDOR
	CONDAND
	NULLCOALESCING //??
	TERNARY
	EQUALS
	LESSGREATER
	BITOR
	BITXOR
	BITAND
	SHIFTS
	SLICE
	DOTDOT
	SUM
	PRODUCT
	PREFIX
	MATCHING
	CALL
	INDEX
	INCREMENT
)

var precedences = map[token.TokenType]int{
	token.PIPE:       PIPE,
	token.ASSIGN:     ASSIGN,
	token.CONDOR:     CONDOR,
	token.OR:         CONDOR,
	token.AND:        CONDAND,
	token.CONDAND:    CONDAND,
	token.EQ:         EQUALS,
	token.NEQ:        EQUALS,
	token.LT:         LESSGREATER,
	token.LE:         LESSGREATER,
	token.GT:         LESSGREATER,
	token.GE:         LESSGREATER,
	token.UDO:        LESSGREATER, // User defined Operator
	token.BITOR:      BITOR,
	token.BITOR_A:    BITOR,
	token.BITXOR_A:   BITXOR,
	token.BITXOR:     BITXOR,
	token.BITAND_A:   BITAND,
	token.BITAND:     BITAND,
	token.SHIFT_L:    SHIFTS,
	token.SHIFT_R:    SHIFTS,
	token.COLON:      SLICE,
	token.QUESTIONM:  TERNARY,
	token.QUESTIONMM: NULLCOALESCING,
	token.DOTDOT:     DOTDOT,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.PLUS_A:     SUM,
	token.MINUS_A:    SUM,
	token.MOD:        PRODUCT,
	token.MOD_A:      PRODUCT,
	token.ASTERISK:   PRODUCT,
	token.ASTERISK_A: PRODUCT,
	token.SLASH:      PRODUCT,
	token.SLASH_A:    PRODUCT,
	token.POWER:      PRODUCT,
	token.MATCH:      MATCHING,
	token.NOTMATCH:   MATCHING,
	token.LPAREN:     CALL,
	token.DOT:        CALL,
	token.LBRACKET:   INDEX,
	token.INCREMENT:  INCREMENT,
	token.DECREMENT:  INCREMENT,
	token.FATARROW:   FATARROW,

	//Meta-Operators
	token.TILDEPLUS:     SUM,
	token.TILDEMINUS:    SUM,
	token.TILDEASTERISK: PRODUCT,
	token.TILDESLASH:    PRODUCT,
	token.TILDEMOD:      PRODUCT,
	token.TILDECARET:    PRODUCT,
}

// A Mode value is a set of flags (or 0).
// They control the amount of source code parsed and other optional
// parser functionality.
//
type Mode uint

const (
	ParseComments Mode = 1 << iota // parse comments and add them to AST
	Trace                          // print a trace of parsed productions
)

// implement sort interface
type SortByLine []ast.Node

func (d SortByLine) Len() int           { return len(d) }
func (d SortByLine) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d SortByLine) Less(i, j int) bool { return d[i].Pos().Line < d[j].Pos().Line }

var (
	FileLines     []string
	tmpDebugInfos SortByLine
	DebugInfos    SortByLine
)

//group by ast.Node's line number.
func SplitSlice(list []ast.Node) [][]ast.Node {
	sort.Sort(SortByLine(list))
	returnData := make([][]ast.Node, 0)
	i := 0
	var j int
	for {
		if i >= len(list) {
			break
		}

		for j = i + 1; j < len(list) && list[i].Pos().Line == list[j].Pos().Line; j++ {
		}

		returnData = append(returnData, list[i:j])
		i = j
	}

	return returnData
}

type Parser struct {
	// Tracing/debugging
	mode   Mode // parsing mode
	trace  bool // == (mode & Trace != 0)
	indent int  // indentation used for tracing output

	// Comments
	comments    []*ast.CommentGroup
	lineComment *ast.CommentGroup // last line comment

	l      *lexer.Lexer
	errors []string
	path   string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	classMap map[string]bool

	//for debugger use
	Functions map[string]*ast.FunctionLiteral
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func NewWithDoc(l *lexer.Lexer, wd string) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
		path:   wd,
		mode:   ParseComments,
	}
	p.l.SetMode(lexer.ScanComments)

	p.classMap = make(map[string]bool)
	p.Functions = make(map[string]*ast.FunctionLiteral)
	p.registerAction()
	p.nextToken()
	p.nextToken()
	return p
}

func New(l *lexer.Lexer, wd string) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
		path:   wd,
	}

	p.classMap = make(map[string]bool)
	p.Functions = make(map[string]*ast.FunctionLiteral)
	p.registerAction()
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerAction() {
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.UINT, p.parseUIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.PLUS, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.UNLESS, p.parseUnlessExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.DO, p.parseDoLoopExpression)
	p.registerPrefix(token.WHILE, p.parseWhileLoopExpression)
	p.registerPrefix(token.FOR, p.parseForLoopExpression)
	p.registerPrefix(token.GREP, p.parseGrepExpression)
	p.registerPrefix(token.MAP, p.parseMapExpression)
	p.registerPrefix(token.CASE, p.parseCaseExpression)
	p.registerPrefix(token.TRY, p.parseTryStatement)
	p.registerPrefix(token.STRING, p.parseStringLiteralExpression)
	p.registerPrefix(token.REGEX, p.parseRegExLiteralExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayExpression)
	p.registerPrefix(token.LBRACE, p.parseHashExpression)
	p.registerPrefix(token.STRUCT, p.parseStructExpression)
	p.registerPrefix(token.ISTRING, p.parseInterpolatedString)
	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)
	p.registerPrefix(token.INCREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.DECREMENT, p.parsePrefixExpression)
	p.registerPrefix(token.NIL, p.parseNilExpression)
	p.registerPrefix(token.ENUM, p.parseEnumExpression)
	p.registerPrefix(token.QW, p.parseQWExpression)
	p.registerPrefix(token.CLASS, p.parseClassLiteral)
	p.registerPrefix(token.NEW, p.parseNewExpression)
	p.registerPrefix(token.UDO, p.parsePrefixExpression)
	p.registerPrefix(token.CMD, p.parseCommand)

	//Meta-Operators
	p.registerPrefix(token.TILDEPLUS, p.parsePrefixExpression)
	p.registerPrefix(token.TILDEMINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TILDEASTERISK, p.parsePrefixExpression)
	p.registerPrefix(token.TILDESLASH, p.parsePrefixExpression)
	p.registerPrefix(token.TILDEMOD, p.parsePrefixExpression)
	p.registerPrefix(token.TILDECARET, p.parsePrefixExpression)

	//linq query
	p.registerPrefix(token.FROM, p.parseLinqExpression)

	//async & await
	p.registerPrefix(token.ASYNC, p.parseAsyncLiteral)
	p.registerPrefix(token.AWAIT, p.parseAwaitExpression)

	//datetime literal
	p.registerPrefix(token.DATETIME, p.parseDateTime)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.POWER, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.MATCH, p.parseInfixExpression)
	p.registerInfix(token.NOTMATCH, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LE, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.CONDAND, p.parseInfixExpression)
	p.registerInfix(token.CONDOR, p.parseInfixExpression)
	p.registerInfix(token.SHIFT_L, p.parseInfixExpression)
	p.registerInfix(token.SHIFT_R, p.parseInfixExpression)
	p.registerInfix(token.BITAND, p.parseInfixExpression)
	p.registerInfix(token.BITOR, p.parseInfixExpression)
	p.registerInfix(token.BITXOR, p.parseInfixExpression)
	p.registerInfix(token.UDO, p.parseInfixExpression)
	p.registerInfix(token.QUESTIONMM, p.parseInfixExpression)

	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.PLUS_A, p.parseAssignExpression)
	p.registerInfix(token.MINUS_A, p.parseAssignExpression)
	p.registerInfix(token.ASTERISK_A, p.parseAssignExpression)
	p.registerInfix(token.SLASH_A, p.parseAssignExpression)
	p.registerInfix(token.MOD_A, p.parseAssignExpression)
	p.registerInfix(token.BITOR_A, p.parseAssignExpression)
	p.registerInfix(token.BITAND_A, p.parseAssignExpression)
	p.registerInfix(token.BITXOR_A, p.parseAssignExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpressions)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseMethodCallExpression)
	p.registerInfix(token.DOTDOT, p.parseRangeLiteralExpression)
	p.registerInfix(token.QUESTIONM, p.parseTernaryExpression)
	p.registerInfix(token.COLON, p.parseSliceExpression)
	p.registerInfix(token.INCREMENT, p.parsePostfixExpression)
	p.registerInfix(token.DECREMENT, p.parsePostfixExpression)
	p.registerInfix(token.PIPE, p.parsePipeExpression)
	p.registerInfix(token.FATARROW, p.parseFatArrowFunction)

	//Meta-Operators
	p.registerInfix(token.TILDEPLUS, p.parseInfixExpression)
	p.registerInfix(token.TILDEMINUS, p.parseInfixExpression)
	p.registerInfix(token.TILDEASTERISK, p.parseInfixExpression)
	p.registerInfix(token.TILDESLASH, p.parseInfixExpression)
	p.registerInfix(token.TILDEMOD, p.parseInfixExpression)
	p.registerInfix(token.TILDECARET, p.parseInfixExpression)
}

func (p *Parser) ParseProgram() *ast.Program {
	defer func() {
		if r := recover(); r != nil {
			return //Here we just 'return', because the caller will report the errors.
		}
	}()

	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	program.Includes = make(map[string]*ast.IncludeStatement)

	//if the magpie file only have ';', then we should return earlier.
	if p.curTokenIs(token.SEMICOLON) && p.peekTokenIs(token.EOF) {
		return program
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			if include, ok := stmt.(*ast.IncludeStatement); ok {
				includePath := strings.TrimSpace(include.IncludePath.String())
				_, ok := program.Includes[includePath]
				if !ok {
					program.Includes[includePath] = include
				}
			} else {
				program.Statements = append(program.Statements, stmt)
			}
		}
		p.nextToken()
	}

	for _, n := range tmpDebugInfos {
		switch n.(type) {
		case *ast.LetStatement:
			l := n.(*ast.LetStatement)
			if !l.InClass {
				DebugInfos = append(DebugInfos, n)
			}
		case *ast.ConstStatement:
			DebugInfos = append(DebugInfos, n)
		case *ast.ReturnStatement:
			DebugInfos = append(DebugInfos, n)
		case *ast.DeferStmt:
			DebugInfos = append(DebugInfos, n)
		case *ast.EnumStatement:
			DebugInfos = append(DebugInfos, n)
		case *ast.IfExpression:
			DebugInfos = append(DebugInfos, n)
		case *ast.UnlessExpression:
			DebugInfos = append(DebugInfos, n)
		case *ast.CaseExpr:
			DebugInfos = append(DebugInfos, n)
		case *ast.DoLoop:
			DebugInfos = append(DebugInfos, n)
		case *ast.WhileLoop:
			DebugInfos = append(DebugInfos, n)
		case *ast.ForLoop:
			DebugInfos = append(DebugInfos, n)
		case *ast.ForEverLoop:
			DebugInfos = append(DebugInfos, n)
		case *ast.ForEachArrayLoop:
			DebugInfos = append(DebugInfos, n)
		case *ast.ForEachDotRange:
			DebugInfos = append(DebugInfos, n)
		case *ast.ForEachMapLoop:
			DebugInfos = append(DebugInfos, n)
		case *ast.BreakExpression:
			DebugInfos = append(DebugInfos, n)
		case *ast.ContinueExpression:
			DebugInfos = append(DebugInfos, n)
		case *ast.AssignExpression:
			DebugInfos = append(DebugInfos, n)
		case *ast.CallExpression:
			DebugInfos = append(DebugInfos, n)
		case *ast.TryStmt:
			DebugInfos = append(DebugInfos, n)
		case *ast.SpawnStmt:
			DebugInfos = append(DebugInfos, n)
		case *ast.UsingStmt:
			DebugInfos = append(DebugInfos, n)
		case *ast.QueryExpr:
			DebugInfos = append(DebugInfos, n)
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	var ret ast.Statement
	switch p.curToken.Type {
	case token.LET:
		ret = p.parseLetStatement(false, true)
	case token.CONST:
		ret = p.parseConstStatement()
	case token.RETURN:
		ret = p.parseReturnStatement()
	case token.DEFER:
		ret = p.parseDeferStatement()
	case token.SPAWN:
		ret = p.parseSpawnStatement()
	case token.INCLUDE:
		return p.parseIncludeStatement()
	case token.THROW:
		ret = p.parseThrowStatement()
	case token.FUNCTION:
		//if p.peekTokenIs(token.IDENT) { //function statement. e.g. 'fn add(x,y) { xxx }'
		//	return p.parseFunctionStatement()
		//} else {
		//	// if we reach here, it means the "FN" token is
		//	// assumed to be the beginning of an expression.
		//	return nil
		//}

		//Because we need to support operator overloading,like:
		//     fn +(v) { block }
		//so we should not use above code
		return p.parseFunctionStatement()
	case token.CLASS:
		return p.parseClassStatement()
	case token.SERVICE:
		return p.parseServiceStatement()
	case token.ENUM:
		ret = p.parseEnumStatement()
	case token.USING:
		ret = p.parseUsingStatement()
	case token.ASYNC:
		return p.parseAsyncStatement()
	case token.IDENT:
		//if the current token is an 'identifier' and next token is a ',', 
		//then we think it's a multiple assignment, but we threat it as a 'let' statement.
		//otherwise, we just fallthrough.
		if p.peekTokenIs(token.COMMA) {
			return p.parseLetStatement(false, false)
		}
		fallthrough
	default:
		return p.parseExpressionStatement()
	}

	tmpDebugInfos = append(tmpDebugInfos, ret)
	return ret
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

//class classname : parentClass { block }
//class classname (categoryname) { block }  //has category name
//class classname () { block }              //no category name
//class @classname : parentClass { block }  //annotation
func (p *Parser) parseClassStatement() *ast.ClassStatement {
	stmt := &ast.ClassStatement{Token: p.curToken}
	stmt.Doc = p.lineComment

	if p.peekTokenIs(token.AT) { //declare an annotion.
		p.nextToken()
		stmt.IsAnnotation = true
	}

	if !p.expectPeek(token.IDENT) { //classname
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.classMap[stmt.Name.Value] = true

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		if p.peekTokenIs(token.RPAREN) { //the category name is empty
			//create a dummy category name
			tok := token.Token{Type: token.ILLEGAL, Literal: ""}
			stmt.CategoryName = &ast.Identifier{Token: tok, Value: ""}
			p.nextToken() //skip current token
		} else if p.peekTokenIs(token.IDENT) {
			p.nextToken() //skip current token
			stmt.CategoryName = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			p.nextToken()
		} else {
			pos := p.fixPosCol()
			msg := fmt.Sprintf("Syntax Error:%v- Class's category should be followed by an identifier or a ')', got %s instead.", pos, p.peekToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	if stmt.IsAnnotation {
		stmt.ClassLiteral = p.parseClassLiteralForAnno().(*ast.ClassLiteral)
	} else {
		stmt.ClassLiteral = p.parseClassLiteral().(*ast.ClassLiteral)
	}

	stmt.ClassLiteral.Name = stmt.Name.Value

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	stmt.SrcEndToken = p.curToken
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// Run the infix function until the next token has
	// a higher precedence.
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			tmpDebugInfos = append(tmpDebugInfos, leftExp)
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	tmpDebugInfos = append(tmpDebugInfos, leftExp)
	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()

	// if the token is '**', we process it specially. e.g. 3 ** 2 ** 3 = 3 ** (2 ** 3)
	// i.e. Exponent operator '**'' has right-to-left associativity
	if p.curTokenIs(token.POWER) {
		precedence-- 
	}

	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{Token: p.curToken, Left: left, Operator: p.curToken.Literal}
}

// `cmd xx xx`
func (p *Parser) parseCommand() ast.Expression {
	return &ast.CmdExpression{Token: p.curToken, Value: p.curToken.Literal}
}

//func (p *Parser) parseGroupedExpression() ast.Expression {
//	curToken := p.curToken
//	p.nextToken()
//
//	if p.curTokenIs(token.COMMA) {
//		if !p.expectPeek(token.RPAREN) { //empty tuple
//			return nil
//		}
//		ret := &ast.TupleLiteral{Token: curToken, Members: []ast.Expression{}}
//		return ret
//	}
//
//	// NOTE: if previous token is toke.LPAREN, and the current
//	//       token is token.RPAREN, that is an empty parentheses(e.g. '() => 5'),
//	//       we need to return earlier.
//	if curToken.Type == token.LPAREN && p.curTokenIs(token.RPAREN) {
//		return nil
//	}
//
//	exp := p.parseExpression(LOWEST)
//
//	if p.peekTokenIs(token.COMMA) {
//		p.nextToken()
//		ret := p.parseTupleExpression(curToken, exp)
//		return ret
//	}
//
//	if !p.expectPeek(token.RPAREN) {
//		return nil
//	}
//
//	return exp
//}

func (p *Parser) parseGroupedExpression() ast.Expression {
	curToken := p.curToken
	p.nextToken()

	// NOTE: if previous token is toke.LPAREN, and the current
	//       token is token.RPAREN, that is an empty parentheses,
	//       we need to return earlier.
	if curToken.Type == token.LPAREN && p.curTokenIs(token.RPAREN) {
		if p.peekTokenIs(token.FATARROW) { //e.g. '() => 5': this is a short function
			p.nextToken() //skip current token
			ret := p.parseFatArrowFunction(nil)
			return ret
		}

		//empty tuple, e.g. 'x = ()'
		return &ast.TupleLiteral{Token: curToken, Members: []ast.Expression{}, RParenToken: p.curToken}
	}

	exp := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.COMMA) {
		p.nextToken()
		ret := p.parseTupleExpression(curToken, exp)
		return ret
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseTupleExpression(tok token.Token, expr ast.Expression) ast.Expression {
	members := []ast.Expression{expr}

	oldToken := tok
	for {
		switch p.curToken.Type {
		case token.RPAREN:
			ret := &ast.TupleLiteral{Token: tok, Members: members, RParenToken: p.curToken}
			return ret
		case token.COMMA:
			p.nextToken()
			//For a 1-tuple: "(1,)", the trailing comma is necessary to distinguish it
			//from the parenthesized expression (1).
			if p.curTokenIs(token.RPAREN) { //e.g.  let x = (1,)
				ret := &ast.TupleLiteral{Token: tok, Members: members, RParenToken: p.curToken}
				return ret
			}
			members = append(members, p.parseExpression(LOWEST))
			oldToken = p.curToken
			p.nextToken()
		default:
			oldToken.Pos.Col = oldToken.Pos.Col + len(oldToken.Literal)
			msg := fmt.Sprintf("Syntax Error:%v- expected token to be ',' or ')', got %s instead", oldToken.Pos, p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
	}
}

func (p *Parser) parseTryStatement() ast.Expression {

	savedToken := p.curToken
	ts := &ast.TryStmt{Token: p.curToken}
	ts.Catches = []ast.Expression{}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	ts.Block = p.parseBlockStatement()
	p.nextToken()

	for {
		if !p.curTokenIs(token.CATCH) {
			break
		}
		p.nextToken()
		savedToken := p.curToken
		if !p.curTokenIs(token.LBRACE) {
			catchStmt := &ast.CatchStmt{Token: savedToken}

			if p.curToken.Type == token.STRING {
				catchStmt.Var = p.curToken.Literal
				catchStmt.VarType = 0
			} else if p.curToken.Type == token.IDENT {
				aVar := p.parseIdentifier()
				catchStmt.Var = aVar.(*ast.Identifier).Value
				catchStmt.VarType = 1
			} else {
				return nil
			}

			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			catchStmt.Block = p.parseBlockStatement()
			ts.Catches = append(ts.Catches, catchStmt)
		} else {
			if !p.curTokenIs(token.LBRACE) {
				return nil
			}
			catchAllStmt := &ast.CatchAllStmt{Token: savedToken}
			catchAllStmt.Block = p.parseBlockStatement()
			ts.Catches = append(ts.Catches, catchAllStmt)
		}

		if !p.peekTokenIs(token.CATCH) {
			break
		}
		p.nextToken()
	} //end for

	//finally
	if p.curTokenIs(token.FINALLY) || p.peekTokenIs(token.FINALLY) {
		if p.peekTokenIs(token.FINALLY) {
			p.nextToken()
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
		} else {
			p.nextToken()
		}

		ts.Finally = p.parseBlockStatement()
	}

	if len(ts.Catches) == 0 && ts.Finally == nil { //no catch and no finally
		msg := fmt.Sprintf("Syntax Error:%v- Try block should have at least one 'catch' or 'finally' block.", savedToken.Pos)
		p.errors = append(p.errors, msg)
		return nil
	}
	return ts
}

func (p *Parser) parseThrowStatement() *ast.ThrowStmt {
	stmt := &ast.ThrowStmt{Token: p.curToken}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}
	p.nextToken()
	stmt.Expr = p.parseExpressionStatement().Expression

	return stmt

}

func (p *Parser) parseStringLiteralExpression() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInterpolatedString() ast.Expression {
	is := &ast.InterpolatedString{Token: p.curToken, Value: p.curToken.Literal, ExprMap: make(map[byte]ast.Expression)}

	key := "0"[0]
	for {
		if p.curTokenIs(token.LBRACE) {
			p.nextToken()
			expr := p.parseExpression(LOWEST)
			is.ExprMap[key] = expr
			key++
		}
		p.nextInterpToken()
		if p.curTokenIs(token.ISTRING) {
			break
		}
	}

	return is
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken, ReturnValues: []ast.Expression{}}
	if p.peekTokenIs(token.SEMICOLON) { //e.g.{ return; }
		p.nextToken()
		return stmt
	}
	if p.peekTokenIs(token.RBRACE) { //e.g. { return }
		return stmt
	}

	p.nextToken()
	for {
		v := p.parseExpressionStatement().Expression
		stmt.ReturnValues = append(stmt.ReturnValues, v)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
		p.nextToken()
	}

	if len(stmt.ReturnValues) > 0 {
		stmt.ReturnValue = stmt.ReturnValues[0]
	}
	return stmt
}

func (p *Parser) parseDeferStatement() *ast.DeferStmt {
	stmt := &ast.DeferStmt{Token: p.curToken}

	p.nextToken()
	stmt.Call = p.parseExpressionStatement().Expression

	return stmt
}

func (p *Parser) parseBreakWithoutLoopContext() ast.Expression {
	msg := fmt.Sprintf("Syntax Error:%v- 'break' outside of loop context", p.curToken.Pos)
	p.errors = append(p.errors, msg)

	return p.parseBreakExpression()
}

func (p *Parser) parseBreakExpression() ast.Expression {
	return &ast.BreakExpression{Token: p.curToken}
}

func (p *Parser) parseContinueWithoutLoopContext() ast.Expression {
	msg := fmt.Sprintf("Syntax Error:%v- 'continue' outside of loop context", p.curToken.Pos)
	p.errors = append(p.errors, msg)

	return p.parseContinueExpression()
}

func (p *Parser) parseContinueExpression() ast.Expression {
	return &ast.ContinueExpression{Token: p.curToken}
}

//let a,b,c = 1,2,3 (with assignment)
//let a; (without assignment, 'a' is assumed to be 'nil')
func (p *Parser) parseLetStatement(inClass bool, nextFlag bool) *ast.LetStatement {
	var stmt *ast.LetStatement
	if !nextFlag {
		//construct a dummy 'let' token
		tok := token.Token{Pos: p.curToken.Pos, Type: token.LET, Literal: "let"}
		stmt = &ast.LetStatement{Token: tok, InClass: inClass}
	} else {
		stmt = &ast.LetStatement{Token: p.curToken, InClass: inClass}
	}

	stmt.Doc = p.lineComment

	//parse left hand side of the assignment
	for {
		if nextFlag {
			p.nextToken()
		}
		nextFlag = true

		if !p.curTokenIs(token.IDENT) && !p.curTokenIs(token.UNDERSCORE) {
			msg := fmt.Sprintf("Syntax Error:%v- expected token to be identifier|underscore, got %s instead.", p.curToken.Pos, p.curToken.Type)
			p.errors = append(p.errors, msg)
			return stmt
		}
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		stmt.Names = append(stmt.Names, name)

		if !p.peekTokenIs(token.ASSIGN) && !p.curTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.COMMA) {
			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}
			stmt.SrcEndToken = p.curToken
			return stmt

		}

		p.nextToken()
		if p.curTokenIs(token.ASSIGN) || p.curTokenIs(token.SEMICOLON) {
			break
		}
		if !p.curTokenIs(token.COMMA) {
			msg := fmt.Sprintf("Syntax Error:%v- expected token to be comma, got %s instead.", p.curToken.Pos, p.curToken.Type)
			p.errors = append(p.errors, msg)
			return stmt
		}
	}

	if p.curTokenIs(token.SEMICOLON) { //let x;
		stmt.SrcEndToken = p.curToken
		return stmt
	}

	i := 0
	p.nextToken()
	for {
		var v ast.Expression
		if p.curTokenIs(token.CLASS) { //e.g.  let cls = class { block }
			v = p.parseClassLiteral()
			if len(stmt.Names) >= i {
				v.(*ast.ClassLiteral).Name = stmt.Names[i].Value //get ClassLiteral's class Name
				p.classMap[stmt.Names[i].Value] = true
			}
		} else {
			v = p.parseExpressionStatement().Expression
		}
		stmt.Values = append(stmt.Values, v)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
		p.nextToken()

		i += 1
	}

	stmt.SrcEndToken = p.curToken

	for idx, name := range stmt.Names {
		if idx < len(stmt.Values) {
			if v, ok := stmt.Values[idx].(*ast.FunctionLiteral); ok {
				p.Functions[name.Value] = v
			}
		}
	} //end for

	return stmt
}

//const a = xxx
func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}
	stmt.Doc = p.lineComment

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	stmt.SrcEndToken = p.curToken
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	expression := &ast.BlockStatement{Token: p.curToken}
	expression.Statements = []ast.Statement{}
	p.nextToken() //skip '{'
	for !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			expression.Statements = append(expression.Statements, stmt)
		}
		if p.peekTokenIs(token.EOF) {
			break
		}
		p.nextToken()
	}

	/*  LONG HIDDEN BUG!
	NOTE: If we got 'EOF' and current token is not '}', then that means that the block is not ended with a '}', like below:

		//if.my
	   if (10 > 2) {
	       println("10>2")

	Above 'if' expression has no '}', if we do not check below condition, it will evaluate correctly and no problem occurred.
	*/
	if p.peekTokenIs(token.EOF) && !p.curTokenIs(token.RBRACE) {
		pos := p.peekToken.Pos
		pos.Col += 1
		msg := fmt.Sprintf("Syntax Error:%v- expected next token to be '}', got EOF instead. Block should end with '}'.", pos)
		p.errors = append(p.errors, msg)
	}

	expression.RBraceToken = p.curToken
	return expression
}

func (p *Parser) parseAssignExpression(name ast.Expression) ast.Expression {
	e := &ast.AssignExpression{Token: p.curToken}

	// if n, ok := name.(*ast.Identifier); ok { //e.g. a = 10
	// 	e.Name = n
	// } else if call, ok := name.(*ast.MethodCallExpression); ok { //might be 'includeModule.a = xxx' or 'aHashObj.key = value'
	// 	e.Name = call
	// } else if indexExp, ok := name.(*ast.IndexExpression); ok {
	// 	//e.g.
	// 	//    doc = {"one": {"two": {"three": [1, 2, 3,] }}}
	// 	//    doc["one"]["two"]["three"][2] = 44
	// 	e.Name = indexExp
	// }

	e.Name = name
	p.nextToken()
	e.Value = p.parseExpression(LOWEST)
	switch v := e.Value.(type) {
	case *ast.InfixExpression:
		/*
		   e.g. In a '*.mp' file, you only have below line:
		        c = 1 -
		   and nothing more, then we assume it's a syntax error
		*/
		if v.Right == nil {
			msg := fmt.Sprintf("Syntax Error:%v- No right part of infix-expression", p.curToken.Pos)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	return e
}

func (p *Parser) parseIncludeStatement() *ast.IncludeStatement {
	stmt := &ast.IncludeStatement{Token: p.curToken}

	p.nextToken()
	if p.curToken.Type != token.STRING && p.curToken.Type != token.IDENT {
		msg := fmt.Sprintf("Syntax Error:%v- expected token to be STRING|IDENTIFIER, got %s instead", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return stmt
	}

	oldToken := p.curToken
	stmt.IncludePath = p.parseExpressionStatement().Expression
	includePath := strings.TrimSpace(stmt.IncludePath.String())
	if oldToken.Type == token.STRING { //if token type is STRING, we need to extract the basename of the path.
		path := stmt.IncludePath.(*ast.StringLiteral).Value
		includePath = path
		baseName := filepath.Base(path)
		oldToken.Literal = baseName
		stmt.IncludePath = &ast.StringLiteral{Token: oldToken, Value: baseName}
	}

	program, err := p.getIncludedStatements(includePath)
	if err != nil {
		p.errors = append(p.errors, err.Error())
		return stmt
	}
	stmt.Program = program
	return stmt
}

func (p *Parser) getIncludedStatements(importpath string) (*ast.Program, error) {
	path := p.path

	fn := filepath.Join(path, importpath+".my")
	f, err := ioutil.ReadFile(fn)
	if err != nil { //error occurred, maybe the file do not exists.
		// Check for 'MAGPIE_ROOT' environment variable
		includeRoot := os.Getenv("MAGPIE_ROOT")
		if len(includeRoot) == 0 { //'MAGPIE_ROOT' environment variable is not set
			return nil, fmt.Errorf("Syntax Error:%v- no file or directory: %s.my, %s", p.curToken.Pos, importpath, path)
		} else {
			fn = filepath.Join(includeRoot, importpath+".my")
			e, err := ioutil.ReadFile(fn)
			if err != nil {
				return nil, fmt.Errorf("Syntax Error:%v- no file or directory: %s.my, %s", p.curToken.Pos, importpath, includeRoot)
			}
			f = e
		}
	}

	l := lexer.New(fn, string(f))
	var ps *Parser
	if p.mode&ParseComments == 0 {
		ps = New(l, path)
	} else {
		ps = NewWithDoc(l, path)
	}
	parsed := ps.ParseProgram()
	if len(ps.errors) != 0 {
		p.errors = append(p.errors, ps.errors...)
	}
	return parsed, nil
}

func (p *Parser) parseDoLoopExpression() ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.DoLoop{Token: p.curToken}

	p.expectPeek(token.LBRACE)
	loop.Block = p.parseBlockStatement()

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

func (p *Parser) parseWhileLoopExpression() ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.WhileLoop{Token: p.curToken}

	p.nextToken()
	loop.Condition = p.parseExpressionStatement().Expression

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
	}

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		loop.Block = p.parseBlockStatement()
	} else if p.peekTokenIs(token.FATARROW) {
		p.nextToken() //skip current token
		p.nextToken() //skip '=>' token
		loop.Block = p.parseExpressionStatement().Expression
	} else {
		msg := fmt.Sprintf("Syntax Error:%v- for loop must be followed by a '{' or '=>'.", p.curToken.Pos)
		p.errors = append(p.errors, msg)
		return nil
	}

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

func (p *Parser) parseForLoopExpression() ast.Expression {
	curToken := p.curToken //save current token

	if p.peekTokenIs(token.LBRACE) {
		return p.parseForEverLoopExpression(curToken)
	}

	if p.peekTokenIs(token.LPAREN) {
		return p.parseCForLoopExpression(curToken)
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	variable := p.curToken.Literal //save current identifier

	if p.peekTokenIs(token.COMMA) {
		return p.parseForEachMapExpression(curToken, variable)
	}

	ret := p.parseForEachArrayOrRangeExpression(curToken, variable)
	return ret
}

//for (init; condition; update) {}
//for (; condition; update) {}  --- init is empty
//for (; condition;;) {}  --- init & update both empty
// for (;;;) {} --- init/condition/update all empty
func (p *Parser) parseCForLoopExpression(curToken token.Token) ast.Expression {
	var result ast.Expression
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	var init ast.Expression
	var cond ast.Expression
	var update ast.Expression

	p.nextToken()
	if !p.curTokenIs(token.SEMICOLON) {
		init = p.parseExpression(LOWEST)
		p.nextToken()
	}

	p.nextToken() //skip ';'
	if !p.curTokenIs(token.SEMICOLON) {
		cond = p.parseExpression(LOWEST)
		p.nextToken()
	}

	p.nextToken()
	if !p.curTokenIs(token.SEMICOLON) {
		update = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.peekTokenIs(token.LBRACE) && !p.peekTokenIs(token.FATARROW) {
		msg := fmt.Sprintf("Syntax Error:%v- for loop must be followed by a '{' or '=>'.", p.curToken.Pos)
		p.errors = append(p.errors, msg)
		return nil
	}

	var isBlock bool = false
	if p.peekTokenIs(token.LBRACE) {
		isBlock = true
	} else {
		p.nextToken()
	}

	p.nextToken()

	if init == nil && cond == nil && update == nil {
		loop := &ast.ForEverLoop{Token: curToken}
		if isBlock {
			loop.Block = p.parseBlockStatement()
		} else {
			msg := fmt.Sprintf("Syntax Error:%v- Never end loop must use block statment.", p.curToken.Pos)
			p.errors = append(p.errors, msg)
			return nil
		}
		result = loop
	} else {
		loop := &ast.ForLoop{Token: curToken, Init: init, Cond: cond, Update: update}
		if isBlock {
			loop.Block = p.parseBlockStatement()
		} else {
			loop.Block = p.parseExpressionStatement().Expression
		}
		result = loop
	}

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return result
}

//for item in array <where cond> {}
//for item in start..end <where cond> {}
func (p *Parser) parseForEachArrayOrRangeExpression(curToken token.Token, variable string) ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	var isRange bool = false
	//loop := &ast.ForEachArrayLoop{Token: curToken, Var:variable}

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	/*
		Note: Here we use precedence 'FATARROW', not 'LOWEST'.
		if we use 'LOWEST', below code will report error:
			c = for i in a => i + 1
		    (error message:  for loop must be followed by a '{' or '=>'.)
	*/
	aValue1 := p.parseExpression(FATARROW)

	var aValue2 ast.Expression
	if p.peekTokenIs(token.DOTDOT) {
		isRange = true
		p.nextToken()
		p.nextToken()
		aValue2 = p.parseExpression(DOTDOT)
	}

	var aCond ast.Expression
	if p.peekTokenIs(token.WHERE) {
		p.nextToken() //skip current token
		p.nextToken() //skip 'where' token
		aCond = p.parseExpression(FATARROW)
	}

	var aBlock ast.Node
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		aBlock = p.parseBlockStatement()
	} else if p.peekTokenIs(token.FATARROW) {
		p.nextToken() //skip current token
		p.nextToken() //skip '=>'
		aBlock = p.parseExpressionStatement().Expression
	} else {
		msg := fmt.Sprintf("Syntax Error:%v- for loop must be followed by a '{' or '=>'.", p.curToken.Pos)
		p.errors = append(p.errors, msg)
		return nil
	}

	var result ast.Expression
	if !isRange {
		result = &ast.ForEachArrayLoop{Token: curToken, Var: variable, Value: aValue1, Cond: aCond, Block: aBlock}
	} else {
		result = &ast.ForEachDotRange{Token: curToken, Var: variable, StartIdx: aValue1, EndIdx: aValue2, Cond: aCond, Block: aBlock}
	}

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return result
}

//for key, value in hash {}
func (p *Parser) parseForEachMapExpression(curToken token.Token, variable string) ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.ForEachMapLoop{Token: curToken}
	loop.Key = variable

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	loop.Value = p.curToken.Literal

	if !p.expectPeek(token.IN) {
		return nil
	}

	p.nextToken()
	loop.X = p.parseExpression(FATARROW)

	if p.peekTokenIs(token.WHERE) {
		p.nextToken()
		p.nextToken()
		loop.Cond = p.parseExpression(FATARROW)
	}

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		loop.Block = p.parseBlockStatement()
	} else if p.peekTokenIs(token.FATARROW) {
		p.nextToken() //skip current token
		p.nextToken() //skip '=>' token
		loop.Block = p.parseExpressionStatement().Expression
	} else {
		msg := fmt.Sprintf("Syntax Error:%v- for loop must be followed by a '{' or '=>'.", p.curToken.Pos)
		p.errors = append(p.errors, msg)
		return nil
	}

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

//Almost same with parseDoLoopExpression()
func (p *Parser) parseForEverLoopExpression(curToken token.Token) ast.Expression {
	p.registerPrefix(token.BREAK, p.parseBreakExpression)
	p.registerPrefix(token.CONTINUE, p.parseContinueExpression)

	loop := &ast.ForEverLoop{Token: curToken}

	p.expectPeek(token.LBRACE)
	loop.Block = p.parseBlockStatement()

	p.registerPrefix(token.BREAK, p.parseBreakWithoutLoopContext)
	p.registerPrefix(token.CONTINUE, p.parseContinueWithoutLoopContext)

	return loop
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	var value int64
	var err error

	p.curToken.Literal = convertNum(p.curToken.Literal)
	if strings.HasPrefix(p.curToken.Literal, "0b") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 2, 64)
	} else if strings.HasPrefix(p.curToken.Literal, "0x") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 16, 64)
	} else if strings.HasPrefix(p.curToken.Literal, "0o") {
		value, err = strconv.ParseInt(p.curToken.Literal[2:], 8, 64)
	} else {
		value, err = strconv.ParseInt(p.curToken.Literal, 10, 64)
	}

	if err != nil {
		msg := fmt.Sprintf("Syntax Error:%v- could not parse %q as integer", p.curToken.Pos, p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseUIntegerLiteral() ast.Expression {
	lit := &ast.UIntegerLiteral{Token: p.curToken}

	var value uint64
	var err error

	p.curToken.Literal = convertNum(p.curToken.Literal)
	if strings.HasPrefix(p.curToken.Literal, "0b") {
		value, err = strconv.ParseUint(p.curToken.Literal[2:], 2, 64)
	} else if strings.HasPrefix(p.curToken.Literal, "0x") {
		value, err = strconv.ParseUint(p.curToken.Literal[2:], 16, 64)
	} else if strings.HasPrefix(p.curToken.Literal, "0o") {
		value, err = strconv.ParseUint(p.curToken.Literal[2:], 8, 64)
	} else {
		value, err = strconv.ParseUint(p.curToken.Literal, 10, 64)
	}

	if err != nil {
		msg := fmt.Sprintf("Syntax Error:%v- could not parse %q as unsigned integer", p.curToken.Pos, p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	p.curToken.Literal = convertNum(p.curToken.Literal)
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("Syntax Error:%v- could not parse %q as float", p.curToken.Pos, p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseRegExLiteralExpression() ast.Expression {
	return &ast.RegExLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseGrepExpression() ast.Expression {
	if p.peekTokenIs(token.LBRACE) {
		return p.parseGrepBlockExpression(p.curToken)
	}
	return p.parseGrepExprExpression(p.curToken)
}

//grep { BLOCK } LIST
func (p *Parser) parseGrepBlockExpression(tok token.Token) ast.Expression {
	gl := &ast.GrepExpr{Token: p.curToken, Var: "$_"}

	p.nextToken()
	gl.Block = p.parseBlockStatement()

	p.nextToken()
	gl.Value = p.parseExpression(LOWEST)

	return gl
}

//grep EXPR, LIST
func (p *Parser) parseGrepExprExpression(tok token.Token) ast.Expression {
	ge := &ast.GrepExpr{Token: p.curToken, Var: "$_"}

	p.nextToken()
	ge.Expr = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	p.nextToken()
	ge.Value = p.parseExpression(LOWEST)

	return ge
}

func (p *Parser) parseMapExpression() ast.Expression {
	if p.peekTokenIs(token.LBRACE) {
		return p.parseMapBlockExpression(p.curToken)
	}
	return p.parseMapExprExpression(p.curToken)
}

//map { BLOCK } LIST
func (p *Parser) parseMapBlockExpression(tok token.Token) ast.Expression {
	me := &ast.MapExpr{Token: p.curToken, Var: "$_"}

	p.nextToken()
	me.Block = p.parseBlockStatement()

	p.nextToken()
	me.Value = p.parseExpression(LOWEST)

	return me
}

//map EXPR, LIST
func (p *Parser) parseMapExprExpression(tok token.Token) ast.Expression {
	me := &ast.MapExpr{Token: p.curToken, Var: "$_"}

	p.nextToken()
	me.Expr = p.parseExpression(LOWEST)

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	p.nextToken()
	me.Value = p.parseExpression(LOWEST)

	return me
}

//func (p *Parser) parseIfExpression() ast.Expression {
//	expression := &ast.IfExpression{Token: p.curToken}
//
//	if p.peekTokenIs(token.LPAREN) {
//		p.nextToken()
//	}
//	p.nextToken()
//	expression.Condition = p.parseExpression(LOWEST)
//
//	if p.peekTokenIs(token.RPAREN) {
//		p.nextToken()
//	}
//
//	if !p.expectPeek(token.LBRACE) {
//		return nil
//	}
//
//	expression.Consequence = p.parseBlockStatement()
//	if p.peekTokenIs(token.ELSE) {
//		p.nextToken()
//		if p.expectPeek(token.LBRACE) {
//			expression.Alternative = p.parseBlockStatement()
//		}
//	}
//
//	return expression
//}

func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{Token: p.curToken}
	// parse if/else-if expressions
	ie.Conditions = p.parseConditionalExpressions(ie)
	return ie
}

func (p *Parser) parseConditionalExpressions(ie *ast.IfExpression) []*ast.IfConditionExpr {
	// if part
	ic := []*ast.IfConditionExpr{p.parseConditionalExpression()}

	//else-if
	for p.peekTokenIs(token.ELSEIF) || p.peekTokenIs(token.ELSIF) || p.peekTokenIs(token.ELIF) || p.peekTokenIs(token.ELSE) { //could be 'elseif' or 'elsif' or 'elif', or 'else'
		p.nextToken()

		if p.curTokenIs(token.ELSE) {
			if !p.peekTokenIs(token.IF) {
				if p.peekTokenIs(token.LBRACE) { //block statement. e.g. 'else {'
					p.nextToken()
					ie.Alternative = p.parseBlockStatement()
				} else { //single expression, e.g. 'else println(xxx)'
					p.nextToken()
					ie.Alternative = p.parseExpressionStatement().Expression
				}
				break
			} else { //'else if'
				p.nextToken()
				ic = append(ic, p.parseConditionalExpression())
			}
		} else {
			ic = append(ic, p.parseConditionalExpression())
		}
	}

	return ic
}

func (p *Parser) parseConditionalExpression() *ast.IfConditionExpr {
	ic := &ast.IfConditionExpr{Token: p.curToken}
	p.nextToken()

	ic.Cond = p.parseExpressionStatement().Expression

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken() //skip current token
	}

	if !p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		ic.Body = p.parseExpressionStatement().Expression
	} else {
		p.nextToken()
		ic.Body = p.parseBlockStatement()
	}

	return ic
}

func (p *Parser) parseUnlessExpression() ast.Expression {
	expression := &ast.UnlessExpression{Token: p.curToken}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseSliceExpression(start ast.Expression) ast.Expression {
	slice := &ast.SliceExpression{Token: p.curToken, StartIndex: start}
	if p.peekTokenIs(token.RBRACKET) {
		slice.EndIndex = nil
	} else {
		p.nextToken()
		slice.EndIndex = p.parseExpression(LOWEST)
	}

	return slice
}

func (p *Parser) parseIndexExpression(arr ast.Expression) ast.Expression {
	var index ast.Expression
	var parameters []ast.Expression
	indexExp := &ast.IndexExpression{Token: p.curToken, Left: arr}
	if p.peekTokenIs(token.COLON) {
		indexTok := token.Token{Type: token.INT, Literal: "0"}
		prefix := &ast.IntegerLiteral{Token: indexTok, Value: int64(0)}
		p.nextToken()
		index = p.parseSliceExpression(prefix)
	} else {
		p.nextToken()
		oldToken := p.curToken
		index = p.parseExpression(LOWEST)
		if p.peekTokenIs(token.COMMA) { //class's index parameter. e.g. 'animalObj[x,y]'
			parameters = append(parameters, index)
			for p.peekTokenIs(token.COMMA) {
				p.nextToken()
				p.nextToken()
				parameters = append(parameters, p.parseExpression(LOWEST))
			}
			index = &ast.ClassIndexerExpression{Token: oldToken, Parameters: parameters}
		}
	}
	indexExp.Index = index
	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
	} else {
		pos := p.fixPosCol()
		msg := fmt.Sprintf("Syntax Error:%v- expected next token to be ']', got %s instead", pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
	}

	return indexExp
}

func (p *Parser) parseHashExpression() ast.Expression {
	curToken := p.curToken //save current token

	if p.peekTokenIs(token.RBRACE) { //empty hash
		p.nextToken()
		hash := &ast.HashLiteral{Token: curToken, RBraceToken: p.curToken, Order: []ast.Expression{}}
		hash.Pairs = make(map[ast.Expression]ast.Expression)
		return hash
	}

	p.nextToken()                       //skip the '{'
	keyExpr := p.parseExpression(SLICE) //note the precedence,if is LOWEST, then it will be parsed as sliceExpression

	if p.peekTokenIs(token.COLON) { //a hash comprehension
		p.nextToken() //skip current token
		p.nextToken() //skip the ':'

		valueExpr := p.parseExpression(LOWEST)

		if p.peekTokenIs(token.COMMA) || p.peekTokenIs(token.RBRACE) { // {k1:v1, ...} or {k1:v1}
			hash := &ast.HashLiteral{Token: curToken, Order: []ast.Expression{}}
			hash.Pairs = make(map[ast.Expression]ast.Expression)

			hash.Pairs[keyExpr] = valueExpr
			hash.Order = append(hash.Order, keyExpr)

			p.nextToken() //skip current token
			if p.curTokenIs(token.RBRACE) {
				hash.RBraceToken = p.curToken
				return hash
			}
			p.nextToken() //skip the ','

			for !p.curTokenIs(token.RBRACE) {
				key := p.parseExpression(SLICE)
				if k, ok := key.(*ast.Identifier); ok {
					//convert *ast.Identifier to *ast.StringLiteral
					key = &ast.StringLiteral{Token: k.Token, Value: k.Value}
				}

				if !p.expectPeek(token.COLON) {
					return nil
				}

				p.nextToken() //skip the ':'
				hash.Pairs[key] = p.parseExpression(LOWEST)
				hash.Order = append(hash.Order, key)
				p.nextToken() // skip the current token'
				if p.curTokenIs(token.RBRACE) {
					break
				}
				if p.curTokenIs(token.COMMA) && p.peekTokenIs(token.RBRACE) { //allow for the last comma symbol
					p.nextToken()
					break
				}
				p.nextToken()
			}
			hash.RBraceToken = p.curToken
			return hash
		}

		if !p.expectPeek(token.FOR) {
			return nil
		}
		if !p.expectPeek(token.IDENT) { //must be an identifier
			return nil
		}

		//{ k:k+1 for k in arr }     -----> k is a variable in an array
		//{ k:k+1 for k,v in hash }  -----> k is a key in a hash
		keyOrVariable := p.curToken.Literal

		if p.peekTokenIs(token.COMMA) { //hash map comprehension
			return p.parseHashMapComprehension(curToken, keyOrVariable, keyExpr, valueExpr, token.RBRACE)
		}

		// hash list comprehension
		return p.parseHashListComprehension(curToken, keyOrVariable, keyExpr, valueExpr, token.RBRACE)
	} else {
		pos := p.fixPosCol()
		msg := fmt.Sprintf("Syntax Error:%v- expected next token to be ':', got %s instead", pos, p.peekToken.Type)
		p.errors = append(p.errors, msg)
	}

	return nil
}

//func (p *Parser) parseHashExpression() ast.Expression {
//	curToken := p.curToken //save current token
//
//	if p.peekTokenIs(token.RBRACE) { //empty hash
//		p.nextToken()
//		hash := &ast.HashLiteral{Token: curToken, RBraceToken:p.curToken}
//		hash.Pairs = make(map[ast.Expression]ast.Expression)
//
//		return hash
//	}
//
//	p.nextToken() //skip the '{'
//	keyExpr := p.parseExpression(SLICE) //note the precedence,if is LOWEST, then it will be parsed as sliceExpression
//	if p.peekTokenIs(token.COLON) { //a hash comprehension
//		p.nextToken() //skip current token
//		p.nextToken() //skip the ':'
//
//		valueExpr := p.parseExpression(LOWEST)
//		if !p.expectPeek(token.FOR) {
//			return nil
//		}
//
//		if !p.expectPeek(token.IDENT) { //must be an identifier
//			return nil
//		}
//
//		//{ k:k+1 for k in arr }     -----> k is a variable in an array
//		//{ k:k+1 for k,v in hash }  -----> k is a key in a hash
//		keyOrVariable := p.curToken.Literal
//
//		if p.peekTokenIs(token.COMMA) { //hash map comprehension
//			return p.parseHashMapComprehension(curToken, keyOrVariable, keyExpr, valueExpr, token.RBRACE)
//		}
//
//		// hash list comprehension
//		return p.parseHashListComprehension(curToken, keyOrVariable, keyExpr, valueExpr, token.RBRACE)
//
//	} else if p.peekTokenIs(token.FATARROW) { //a hash
//		hash := &ast.HashLiteral{Token: curToken}
//		hash.Pairs = make(map[ast.Expression]ast.Expression)
//
//		p.nextToken() //skip current token
//		p.nextToken() //skip the '=>'
//
//		hash.Pairs[keyExpr] = p.parseExpression(LOWEST)
//		p.nextToken() //skip current token
//		for !p.curTokenIs(token.RBRACE) {
//			p.nextToken() //skip the ','
//
//			key := p.parseExpression(SLICE)
//			if !p.expectPeek(token.FATARROW) {
//				return nil
//			}
//
//			p.nextToken() //skip the '=>'
//			hash.Pairs[key] = p.parseExpression(LOWEST)
//			p.nextToken()
//			if p.curTokenIs(token.COMMA) && p.peekTokenIs(token.RBRACE) { //allow for the last comma symbol
//				p.nextToken()
//				break
//			}
//		}
//		hash.RBraceToken = p.curToken
//		return hash
//	} else {
//		pos := p.fixPosCol()
//		msg := fmt.Sprintf("Syntax Error:%v- expected next token to be ':' or '=>', got %s instead", pos, p.peekToken.Type)
//		p.errors = append(p.errors, msg)
//	}
//
//	return nil
//}

func (p *Parser) parseHashMapComprehension(curToken token.Token, key string, keyExpr ast.Expression, valueExpr ast.Expression, closure token.TokenType) ast.Expression {
	if !p.expectPeek(token.COMMA) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	value := p.curToken.Literal

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	X := p.parseExpression(LOWEST)

	var aCond ast.Expression
	if p.peekTokenIs(token.WHERE) {
		p.nextToken()
		p.nextToken()
		aCond = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(closure) {
		return nil
	}

	result := &ast.HashMapComprehension{Token: curToken, Key: key, Value: value, X: X, Cond: aCond, KeyExpr: keyExpr, ValExpr: valueExpr}
	return result
}

func (p *Parser) parseHashListComprehension(curToken token.Token, variable string, keyExpr ast.Expression, valueExpr ast.Expression, closure token.TokenType) ast.Expression {
	var isRange bool = false

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	aValue1 := p.parseExpression(LOWEST)

	var aValue2 ast.Expression
	if p.peekTokenIs(token.DOTDOT) {
		isRange = true
		p.nextToken()
		p.nextToken()
		aValue2 = p.parseExpression(DOTDOT)
	}

	var aCond ast.Expression
	if p.peekTokenIs(token.WHERE) {
		p.nextToken()
		p.nextToken()
		aCond = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(closure) {
		return nil
	}

	var result ast.Expression
	if !isRange {
		result = &ast.HashComprehension{Token: curToken, Var: variable, Value: aValue1, Cond: aCond, KeyExpr: keyExpr, ValExpr: valueExpr}
	} else {
		result = &ast.HashRangeComprehension{Token: curToken, Var: variable, StartIdx: aValue1, EndIdx: aValue2, Cond: aCond, KeyExpr: keyExpr, ValExpr: valueExpr}
	}

	return result
}

func (p *Parser) parseStructExpression() ast.Expression {
	s := &ast.StructLiteral{Token: p.curToken}
	s.Pairs = make(map[ast.Expression]ast.Expression)
	p.expectPeek(token.LBRACE)
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return s
	}
	for !p.curTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.FATARROW) {
			return nil
		}
		p.nextToken()
		s.Pairs[key] = p.parseExpression(LOWEST)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
			continue
		}
		p.nextToken()
	}
	s.RBraceToken = p.curToken
	return s
}

//func (p *Parser) parseArrayExpression() ast.Expression {
//	array := &ast.ArrayLiteral{Token: p.curToken}
//	array.Members = p.parseExpressionArray(array.Members, token.RBRACKET)
//	return array
//}

func (p *Parser) parseArrayExpression() ast.Expression {
	curToken := p.curToken
	temp, b, creationCount := p.parseExpressionArrayEx([]ast.Expression{}, token.RBRACKET)
	if b { //list comprehension or map comprehension
		p.nextToken()                   //skip 'for'
		if !p.expectPeek(token.IDENT) { //must be an identifier
			return nil
		}

		variable := p.curToken.Literal

		if p.peekTokenIs(token.COMMA) { //map comprehension
			return p.parseListMapComprehension(curToken, temp[0], variable, token.RBRACKET) //here 'variable' is the key of the map
		}

		//list comprehension
		return p.parseListComprehension(curToken, temp[0], variable, token.RBRACKET)
	}

	array := &ast.ArrayLiteral{Token: curToken, CreationCount: creationCount}
	array.Members = temp
	return array
}

func (p *Parser) parseListComprehension(curToken token.Token, expr ast.Expression, variable string, closure token.TokenType) ast.Expression {
	var isRange bool = false

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	aValue1 := p.parseExpression(LOWEST)

	var aValue2 ast.Expression
	if p.peekTokenIs(token.DOTDOT) {
		isRange = true
		p.nextToken()
		p.nextToken()
		aValue2 = p.parseExpression(DOTDOT)
	}

	var aCond ast.Expression
	if p.peekTokenIs(token.WHERE) {
		p.nextToken()
		p.nextToken()
		aCond = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(closure) {
		return nil
	}

	var result ast.Expression
	if !isRange {
		result = &ast.ListComprehension{Token: curToken, Var: variable, Value: aValue1, Cond: aCond, Expr: expr}
	} else {
		result = &ast.ListRangeComprehension{Token: curToken, Var: variable, StartIdx: aValue1, EndIdx: aValue2, Cond: aCond, Expr: expr}
	}

	return result
}

func (p *Parser) parseListMapComprehension(curToken token.Token, expr ast.Expression, variable string, closure token.TokenType) ast.Expression {

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	Value := p.curToken.Literal

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	X := p.parseExpression(LOWEST)

	var aCond ast.Expression
	if p.peekTokenIs(token.WHERE) {
		p.nextToken()
		p.nextToken()
		aCond = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(closure) {
		return nil
	}

	result := &ast.ListMapComprehension{Token: curToken, Key: variable, Value: Value, X: X, Cond: aCond, Expr: expr}
	return result
}

func (p *Parser) parseExpressionArrayEx(a []ast.Expression, closure token.TokenType) ([]ast.Expression, bool, *ast.IntegerLiteral) {
	if p.peekTokenIs(closure) {
		p.nextToken()
		if p.peekTokenIs(token.INT) {
			p.nextToken()
			creationCount := p.parseIntegerLiteral()
			return a, false, creationCount.(*ast.IntegerLiteral)
		}
		return a, false, nil
	}

	p.nextToken()
	v := p.parseExpressionStatement().Expression

	a = append(a, v)
	if p.peekTokenIs(token.FOR) { //list comprehension
		return a, true, nil
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if p.peekTokenIs(closure) {
			break
		}
		p.nextToken()
		a = append(a, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(closure) {
		return nil, false, nil
	}

	return a, false, nil
}

// case expr in {
//    expr,expr { expr }
//    expr { expr }
//    else { expr }
// }
func (p *Parser) parseCaseExpression() ast.Expression {

	ce := &ast.CaseExpr{Token: p.curToken}
	ce.Matches = []ast.Expression{}

	p.nextToken()

	ce.Expr = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.IN) {
		ce.IsWholeMatch = false
	} else if p.peekTokenIs(token.IS) {
		ce.IsWholeMatch = true
	} else {
		return nil
	}
	p.nextToken()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		if p.curTokenIs(token.ELSE) {
			aElse := &ast.CaseElseExpr{Token: p.curToken}
			ce.Matches = append(ce.Matches, aElse)
			p.nextToken()
			aElse.Block = p.parseBlockStatement()
			p.nextToken() //skip the '}'
		} else {
			var aMatches []*ast.CaseMatchExpr
			for !p.curTokenIs(token.LBRACE) {
				aMatch := &ast.CaseMatchExpr{Token: p.curToken}
				aMatch.Expr = p.parseExpression(LOWEST)
				aMatches = append(aMatches, aMatch)

				if !p.peekTokenIs(token.COMMA) {
					p.nextToken()
					break
				}
				p.nextToken()
				p.nextToken()

			} //end for

			if !p.curTokenIs(token.LBRACE) {
				msg := fmt.Sprintf("Syntax Error:%v- expected token to be '{', got %s instead", p.curToken.Pos, p.curToken.Type)
				p.errors = append(p.errors, msg)
			}

			aMatchBlock := p.parseBlockStatement()
			for i := 0; i < len(aMatches); i++ {
				aMatches[i].Block = aMatchBlock
			}

			for i := 0; i < len(aMatches); i++ {
				ce.Matches = append(ce.Matches, aMatches[i])
			}

			p.nextToken() //skip the '}'
		}
	} //end for

	if !p.curTokenIs(token.RBRACE) {
		return nil
	}

	return ce
}

//fn name(paramenters)
func (p *Parser) parseFunctionStatement() ast.Statement {
	FnStmt := &ast.FunctionStatement{Token: p.curToken}
	FnStmt.Doc = p.lineComment

	/* why below 'if'? please see below code:
		1. fn add(x, y) { x + y }
		2. fn(x, y) { x + y }(2, 3)
	for the second one, we have no identifier, so we need
	not advance to the next token.
	*/
	if !p.peekTokenIs(token.LPAREN) {
		p.nextToken()
	}

	FnStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	FnStmt.FunctionLiteral = p.parseFunctionLiteral().(*ast.FunctionLiteral)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	FnStmt.SrcEndToken = p.curToken
	p.Functions[FnStmt.Name.Value] = FnStmt.FunctionLiteral
	return FnStmt
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken, Variadic: false}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.parseFuncExpressionArray(fn, token.RPAREN)

	if p.expectPeek(token.LBRACE) {
		fn.Body = p.parseBlockStatement()
	}
	return fn
}

func (p *Parser) parseFuncExpressionArray(fn *ast.FunctionLiteral, closure token.TokenType) {
	if p.peekTokenIs(closure) {
		p.nextToken()
		return
	}

	var hasDefParamValue bool = false
	for {
		p.nextToken()
		if !p.curTokenIs(token.IDENT) {
			msg := fmt.Sprintf("Syntax Error:%v- Function parameter not identifier, GOT(%s)!", p.curToken.Pos, p.curToken.Literal)
			p.errors = append(p.errors, msg)
			return
		}
		key := p.curToken.Literal
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		fn.Parameters = append(fn.Parameters, name)

		if p.peekTokenIs(token.ASSIGN) {
			hasDefParamValue = true
			p.nextToken()
			p.nextToken()
			v := p.parseExpressionStatement().Expression

			if fn.Values == nil {
				fn.Values = make(map[string]ast.Expression)
			}
			fn.Values[key] = v
		} else if !p.peekTokenIs(token.ELLIPSIS) {
			if hasDefParamValue && !fn.Variadic {
				msg := fmt.Sprintf("Syntax Error:%v- Function's default parameter order not correct!", p.curToken.Pos.Sline())
				p.errors = append(p.errors, msg)
				return
			}
		}

		if p.peekTokenIs(token.COMMA) {
			if fn.Variadic {
				msg := fmt.Sprintf("Syntax Error:%v- Variadic argument in function should be last!", p.curToken.Pos.Sline())
				p.errors = append(p.errors, msg)
				return
			}
			p.nextToken()
		}

		if p.peekTokenIs(token.ELLIPSIS) { //Variadic function
			if fn.Variadic {
				msg := fmt.Sprintf("Syntax Error:%v- Only 1 variadic argument is allowed in function!", p.curToken.Pos.Sline())
				p.errors = append(p.errors, msg)
				return
			}
			fn.Variadic = true

			p.nextToken()
			if !p.peekTokenIs(closure) {
				msg := fmt.Sprintf("Syntax Error:%v- Variadic argument in function should be last!", p.curToken.Pos.Sline())
				p.errors = append(p.errors, msg)
				return
			}
		}

		if p.peekTokenIs(closure) {
			p.nextToken()
			break
		}
	}

	return
}

func (p *Parser) parseCallExpressions(f ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.curToken, Function: f}
	call.Arguments = p.parseExpressionArray(call.Arguments, token.RPAREN)
	return call
}

func (p *Parser) parseExpressionArray(a []ast.Expression, closure token.TokenType) []ast.Expression {
	if p.peekTokenIs(closure) {
		p.nextToken()
		return a
	}
	p.nextToken()
	a = append(a, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		a = append(a, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(closure) {
		return nil
	}
	return a
}

func (p *Parser) parseMethodCallExpression(obj ast.Expression) ast.Expression {
	methodCall := &ast.MethodCallExpression{Token: p.curToken, Object: obj}
	p.nextToken()

	name := p.parseIdentifier()
	if !p.peekTokenIs(token.LPAREN) {
		//methodCall.Call = p.parseExpression(LOWEST)
		//Note: here the precedence should not be `LOWEST`, or else when parsing below line:
		//     logger.LDATE + 1 ==> logger.(LDATE + 1)
		methodCall.Call = p.parseExpression(CALL)
	} else {
		p.nextToken()
		methodCall.Call = p.parseCallExpressions(name)
	}

	return methodCall
}

func (p *Parser) parseRangeLiteralExpression(startIdx ast.Expression) ast.Expression {
	expression := &ast.RangeLiteral{
		Token:    p.curToken,
		StartIdx: startIdx,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.EndIdx = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseTernaryExpression(condition ast.Expression) ast.Expression {
	expression := &ast.TernaryExpression{
		Token:     p.curToken,
		Condition: condition,
	}
	precedence := SLICE
	p.nextToken() //skip the '?'
	expression.IfTrue = p.parseExpression(precedence)

	if !p.expectPeek(token.COLON) { //skip the ":"
		return nil
	}

	// Get to next token, then parse the else part
	p.nextToken()
	expression.IfFalse = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseSpawnStatement() *ast.SpawnStmt {
	stmt := &ast.SpawnStmt{Token: p.curToken}

	p.nextToken()
	stmt.Call = p.parseExpressionStatement().Expression

	return stmt
}

func (p *Parser) parseNilExpression() ast.Expression {
	return &ast.NilLiteral{Token: p.curToken}
}

func (p *Parser) parseEnumStatement() ast.Statement {
	oldToken := p.curToken
	enumStmt := &ast.EnumStatement{Token: p.curToken}
	enumStmt.Doc = p.lineComment

	p.nextToken()
	enumStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	enumStmt.EnumLiteral = p.parseEnumExpression().(*ast.EnumLiteral)
	enumStmt.EnumLiteral.Token = oldToken

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	enumStmt.SrcEndToken = p.curToken
	return enumStmt
}

func (p *Parser) parseEnumExpression() ast.Expression {
	var autoInt int64 = 0 //autoIncrement

	e := &ast.EnumLiteral{Token: p.curToken}
	e.Pairs = make(map[ast.Expression]ast.Expression)
	idPair := make(map[string]ast.Expression)

	if !p.expectPeek(token.LBRACE) {
		return e
	}

	for {
		//check for empty `enum`
		if p.peekTokenIs(token.RBRACE) {
			p.nextToken()
			e.RBraceToken = p.curToken
			return e
		}

		// identifier is mandatory here
		if !p.expectPeek(token.IDENT) {
			return e
		}
		enum_id := p.parseIdentifier()

		// peek next that can be only '=' or ',' or '}'
		if !p.peekTokenIs(token.ASSIGN) && !p.peekTokenIs(token.COMMA) && !p.peekTokenIs(token.RBRACE) {
			msg := fmt.Sprintf("Syntax Error:%v- Token %s not allowed here.", p.peekToken.Pos, p.peekToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}

		// check for optional default value (optional only in `INT` case)
		var enum_value ast.Expression
		if p.peekTokenIs(token.ASSIGN) {
			p.nextToken()
			p.nextToken()

			enum_value = p.parseExpressionStatement().Expression
		}

		if enum_value != nil {
			if _, ok := enum_value.(*ast.IntegerLiteral); ok {
				intLiteral := enum_value.(*ast.IntegerLiteral)
				autoInt = intLiteral.Value + 1
			}
		} else {
			//create a new INT token with 'autoInt' as it's value
			tok := token.Token{Type: token.INT, Literal: strconv.Itoa(int(autoInt))}
			enum_value = &ast.IntegerLiteral{Token: tok, Value: autoInt}
			autoInt++
		}

		str_enum_id := enum_id.(*ast.Identifier).Value
		if _, ok := idPair[str_enum_id]; ok { //is identifier redeclared?
			msg := fmt.Sprintf("Syntax Error:%v- Identifier %s redeclared.", p.curToken.Pos, str_enum_id)
			p.errors = append(p.errors, msg)
			return nil
		} else {
			e.Pairs[enum_id] = enum_value
			idPair[str_enum_id] = enum_value
		}

		if !p.peekTokenIs(token.COMMA) {
			p.nextToken()
			break
		}
		p.nextToken()
	}

	e.RBraceToken = p.curToken
	return e
}

//qw(xx, xx, xx, xx)
func (p *Parser) parseQWExpression() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Members = p.parseStrExpressionArray(array.Members)
	return array
}

func isClassStmtToken(t token.Token) bool {
	//only allow below statements:
	//1. let xxx; or let xxx=yyy
	//2. property xxx { }
	//3. fn xxx(parameters) {}
	tt := t.Type //tt: token type
	if tt == token.LET || tt == token.PROPERTY || tt == token.FUNCTION {
		return true
	}

	if tt == token.PUBLIC || tt == token.PROTECTED || tt == token.PRIVATE { //modifier
		return true
	}

	if tt == token.STATIC { //static
		return true
	}

	if tt == token.ASYNC { //async
		return true
	}

	return false
}

//parse annotation class
func (p *Parser) parseClassLiteralForAnno() ast.Expression {
	cls := &ast.ClassLiteral{
		Token:      p.curToken,
		Properties: make(map[string]*ast.PropertyDeclStmt),
	}

	p.nextToken()
	if p.curTokenIs(token.COLON) {
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		cls.Parent = p.curToken.Literal
		p.nextToken()
	}
	if !p.curTokenIs(token.LBRACE) {
		msg := fmt.Sprintf("Syntax Error:%v- expected token to be '{', got %s instead", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	//why not calling parseBlockStatement()?
	//Because we need to parse 'public', 'private' modifiers, also 'get' and 'set'.
	//parseBlockStatement() function do not handling these.
	//cls.Block = p.parseBlockStatement()
	cls.Block = p.parseClassBody(true)
	for _, statement := range cls.Block.Statements {
		//fmt.Printf("In parseClassLiteral, stmt=%s\n", statement.String()) //debugging purpose
		switch s := statement.(type) {
		case *ast.PropertyDeclStmt: //properties
			cls.Properties[s.Name.Value] = s
		default:
			msg := fmt.Sprintf("Syntax Error:%v- Only 'property' statement is allow in class annotation.", s.Pos())
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	return cls
}

// class : parentClass { block }.
//e.g. let classname = class : parentClass { block }
func (p *Parser) parseClassLiteral() ast.Expression {
	cls := &ast.ClassLiteral{
		Token:      p.curToken,
		Members:    make([]*ast.LetStatement, 0),
		Properties: make(map[string]*ast.PropertyDeclStmt),
		Methods:    make(map[string]*ast.FunctionStatement),
	}

	p.nextToken()
	if p.curTokenIs(token.COLON) {
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		cls.Parent = p.curToken.Literal
		p.nextToken()
	}
	if !p.curTokenIs(token.LBRACE) {
		msg := fmt.Sprintf("Syntax Error:%v- expected token to be '{', got %s instead", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	//why not calling parseBlockStatement()?
	//Because we need to parse 'public', 'private' modifiers, also 'get' and 'set'.
	//parseBlockStatement() function do not handling these.
	//cls.Block = p.parseBlockStatement()
	cls.Block = p.parseClassBody(false)
	for _, statement := range cls.Block.Statements {
		//fmt.Printf("In parseClassLiteral, stmt=%s\n", statement.String()) //debugging purpose

		switch s := statement.(type) {
		case *ast.LetStatement: //class fields
			for _, value := range s.Values {
				switch value.(type) {
				case *ast.FunctionLiteral:
					msg := fmt.Sprintf("Syntax Error:%v- Function literal is not allowed in 'let' statement of class.", s.Pos())
					p.errors = append(p.errors, msg)
					return nil
				default:
					cls.Members = append(cls.Members, s)
				}
			} //end for

			for i := len(s.Values); i < len(s.Names); i++ { //For remaining names which have no values.
				cls.Members = append(cls.Members, s)
			}
		case *ast.FunctionStatement: //class methods
			cls.Methods[s.Name.Value] = s
		case *ast.PropertyDeclStmt: //properties
			cls.Properties[s.Name.Value] = s
		default:
			msg := fmt.Sprintf("Syntax Error:%v- Only 'let' statement, 'function' statement and 'property' statement is allow in class definition.", s.Pos())
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	return cls
}

func (p *Parser) parseClassBody(processAnnoClass bool) *ast.BlockStatement {
	stmts := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

	p.nextToken() //skip '{'
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseClassStmt(processAnnoClass)
		if stmt != nil {
			stmts.Statements = append(stmts.Statements, stmt)
		}

		p.nextToken()
	}

	if p.peekTokenIs(token.EOF) && !p.curTokenIs(token.RBRACE) {
		pos := p.peekToken.Pos
		pos.Col += 1
		msg := fmt.Sprintf("Syntax Error:%v- expected next token to be '}', got EOF instead. Block should end with '}'.", pos)
		p.errors = append(p.errors, msg)
	}

	return stmts
}

func (p *Parser) parseClassStmt(processAnnoClass bool) ast.Statement {
	var annos []*ast.AnnotationStmt

LABEL:
	//parse Annotation
	for p.curTokenIs(token.AT) {
		var tokenIsLParen bool
		anno := &ast.AnnotationStmt{Token: p.curToken, Attributes: map[string]ast.Expression{}}

		if !p.expectPeek(token.IDENT) {
			return nil
		}
		anno.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if p.peekTokenIs(token.LPAREN) {
			tokenIsLParen = true
			p.nextToken()
		} else if p.peekTokenIs(token.LBRACE) {
			tokenIsLParen = false
			p.nextToken()
		} else { //marker annotation, e.g. @Demo
			//only 'property' and 'function' can have annotations
			if !p.peekTokenIs(token.FUNCTION) && !p.peekTokenIs(token.PROPERTY) && !p.peekTokenIs(token.AT) && !p.peekTokenIs(token.STATIC) {
				msg := fmt.Sprintf("Syntax Error:%v- expected token to be 'fn'| 'property'|'static', or another annotation, got '%s' instead", p.peekToken.Pos, p.peekToken.Type)
				p.errors = append(p.errors, msg)
				return nil
			}
			tokenIsLParen = false
			p.nextToken()
			annos = append(annos, anno)
			goto LABEL
		}

		for {
			if !p.expectPeek(token.IDENT) {
				return nil
			}
			key := p.curToken.Literal

			if !p.expectPeek(token.ASSIGN) {
				return nil
			}
			p.nextToken()
			value := p.parseExpression(LOWEST)
			anno.Attributes[key] = value
			p.nextToken()
			if !p.curTokenIs(token.COMMA) {
				break
			}
		}

		if tokenIsLParen {
			if !p.curTokenIs(token.RPAREN) {
				msg := fmt.Sprintf("Syntax Error:%v- expected token to be ')', got '%s' instead", p.curToken.Pos, p.curToken.Type)
				p.errors = append(p.errors, msg)
				return nil
			}
		} else if !p.curTokenIs(token.RBRACE) {
			msg := fmt.Sprintf("Syntax Error:%v- expected token to be '}', got '%s' instead", p.curToken.Pos, p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
		p.nextToken()
		annos = append(annos, anno)
	} //end for

	if !isClassStmtToken(p.curToken) {
		msg := fmt.Sprintf("Syntax Error:%v- expected token to be 'let'|'property'|'fn'|'async'|'public'|'protected'|'private'|'static', got %s instead.", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	modifierLevel := ast.ModifierDefault
	tt := p.curToken.Type

	//check if it is a async function
	var Async bool
	if tt == token.ASYNC {
		Async = true
		p.nextToken() //skip the 'async'
	}
	if tt == token.PUBLIC || tt == token.PROTECTED || tt == token.PRIVATE { //modifier
		p.nextToken() //skip the modifier

		switch tt {
		case token.PUBLIC:
			modifierLevel = ast.ModifierPublic
		case token.PROTECTED:
			modifierLevel = ast.ModifierProtected
		case token.PRIVATE:
			modifierLevel = ast.ModifierPrivate
		}
	}

	var staticFlag bool
	if p.curToken.Type == token.STATIC { //static
		p.nextToken() //skip the 'static' keyword
		staticFlag = true
	}

	return p.parseClassSubStmt(modifierLevel, staticFlag, annos, processAnnoClass, Async)
}

func (p *Parser) parseClassSubStmt(modifierLevel ast.ModifierLevel, staticFlag bool, annos []*ast.AnnotationStmt, processAnnoClass bool, Async bool) ast.Statement {
	var r ast.Statement

	if processAnnoClass { //parse annotation class
		if p.curToken.Type != token.PROPERTY {
			msg := fmt.Sprintf("Syntax Error:%v- expected token to be 'property'.Only 'property' statement is allowed in class annotation.", p.curToken.Pos)
			p.errors = append(p.errors, msg)
			return nil
		}
		r = p.parsePropertyDeclStmt(processAnnoClass)
	} else {
		switch p.curToken.Type {
		case token.LET:
			r = p.parseLetStatement(true, true)
		case token.PROPERTY:
			r = p.parsePropertyDeclStmt(processAnnoClass)
		case token.FUNCTION:
			r = p.parseFunctionStatement()
		}

	}

	switch o := r.(type) {
	case *ast.LetStatement:
		o.ModifierLevel = modifierLevel
		o.StaticFlag = staticFlag
		o.Annotations = annos
	case *ast.PropertyDeclStmt:
		o.ModifierLevel = modifierLevel
		o.StaticFlag = staticFlag
		o.Annotations = annos
	case *ast.FunctionStatement:
		o.FunctionLiteral.ModifierLevel = modifierLevel
		o.FunctionLiteral.StaticFlag = staticFlag
		o.FunctionLiteral.Async = Async
		o.Annotations = annos
	}

	return r
}

//property xxx { get; set; }
//property xxx { get; } or
//property xxx { set; }
//property xxx { get {xxx} }
//property xxx { set {xxx} }
//property xxx { get {xxx} set {xxx} }
//property this[x] { get {xxx} set {xxx} }
//property xxx default defaultValue
//property xxx = defaultValue
func (p *Parser) parsePropertyDeclStmt(processAnnoClass bool) *ast.PropertyDeclStmt {
	stmt := &ast.PropertyDeclStmt{Token: p.curToken}
	stmt.Doc = p.lineComment

	if !p.expectPeek(token.IDENT) { //must be an identifier, it's property name
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if processAnnoClass || p.peekTokenIs(token.SEMICOLON) { //annotation class' property defaults to have both getter and setter.
		getterToken := token.Token{Pos: p.curToken.Pos, Type: token.GET, Literal: "get"}
		stmt.Getter = &ast.GetterStmt{Token: getterToken}
		stmt.Getter.Body = &ast.BlockStatement{Statements: []ast.Statement{}}

		setterToken := token.Token{Pos: p.curToken.Pos, Type: token.SET, Literal: "set"}
		stmt.Setter = &ast.SetterStmt{Token: setterToken}
		stmt.Setter.Body = &ast.BlockStatement{Statements: []ast.Statement{}}

		if processAnnoClass {
			if p.peekTokenIs(token.DEFAULT) || p.peekTokenIs(token.ASSIGN) {
				p.nextToken() //skip current token
				p.nextToken() //skip 'default' keyword
				stmt.Default = p.parseExpression(LOWEST)
			}
		} else {
			if p.peekTokenIs(token.SEMICOLON) { //e.g. 'property xxx；'
				p.nextToken()
			}
		}
		stmt.SrcEndToken = p.curToken
		return stmt
	}

	foundIndexer := false
	thisToken := p.curToken           //save the current token for later use
	if p.curToken.Literal == "this" { //assume it is a indexer declaration
		foundIndexer = true
		if !p.expectPeek(token.LBRACKET) { //e.g. 'property this[index]'
			return nil
		}

		if !p.expectPeek(token.IDENT) { //must be an identifier, it's the index
			return nil
		}
		stmt.Indexes = append(stmt.Indexes, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

		if p.peekTokenIs(token.COMMA) {
			for p.peekTokenIs(token.COMMA) {
				p.nextToken()
				p.nextToken()
				stmt.Indexes = append(stmt.Indexes, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
			}
		}

		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
	}

	if foundIndexer {
		//get property name, Note HERE, because we support multiple indexers:
		//    property this[x]   { get {xxx} set {xxx} }
		//    property this[x,y] { get {xxx} set {xxx} }
		//so we need to change the first 'this' to 'this1', and the second to 'this2'
		stmt.Name = &ast.Identifier{Token: thisToken, Value: fmt.Sprintf("this%d", len(stmt.Indexes))}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()
	if p.curTokenIs(token.GET) {
		stmt.Getter = p.parseGetter()
		p.nextToken()
	}

	if p.curTokenIs(token.SET) {
		stmt.Setter = p.parseSetter()
		p.nextToken()
	}

	if !p.curTokenIs(token.RBRACE) {
		return nil
	}

	stmt.SrcEndToken = p.curToken
	return stmt
}

func (p *Parser) parseGetter() *ast.GetterStmt {
	stmt := &ast.GetterStmt{Token: p.curToken}
	p.nextToken() //skip the 'get' keyword

	switch p.curToken.Type {
	case token.SEMICOLON:
		stmt.Body = &ast.BlockStatement{Statements: []ast.Statement{}}
	case token.LBRACE:
		stmt.Body = p.parseBlockStatement()
	}

	return stmt

}

func (p *Parser) parseSetter() *ast.SetterStmt {
	stmt := &ast.SetterStmt{Token: p.curToken}
	p.nextToken() //skip the set keyword

	switch p.curToken.Type {
	case token.SEMICOLON:
		stmt.Body = &ast.BlockStatement{Statements: []ast.Statement{}}
	case token.LBRACE:
		stmt.Body = p.parseBlockStatement()
	}
	return stmt
}

//new classname(xx, xx, xx, xx)
func (p *Parser) parseNewExpression() ast.Expression {
	newExp := &ast.NewExpression{Token: p.curToken}

	p.nextToken()
	exp := p.parseExpression(LOWEST)

	call, ok := exp.(*ast.CallExpression)
	if !ok {
		pos := p.fixPosCol()
		msg := fmt.Sprintf("Syntax Error:%v- Invalid object construction for 'new'. maybe you want 'new xxx()'", pos)
		p.errors = append(p.errors, msg)
		return nil
	}

	clsName := call.Function.(*ast.Identifier).Value
	if _, ok := p.classMap[clsName]; !ok {
		pos := p.fixPosCol()
		msg := fmt.Sprintf("Syntax Error:%v- 'new' should follow a 'class' name.", pos)
		p.errors = append(p.errors, msg)
		return nil
	}
	newExp.Class = call.Function
	newExp.Arguments = call.Arguments

	return newExp
}

func (p *Parser) parseStrExpressionArray(a []ast.Expression) []ast.Expression {
	var allowedPair = map[string]token.TokenType{
		"(": token.RPAREN,
		"<": token.GT,
		"{": token.RBRACE,
	}

	p.nextToken() //skip the 'qw'
	openPair := p.curToken.Literal
	if p.curTokenIs(allowedPair[openPair]) {
		p.nextToken()
		return a
	}

	p.nextToken()
	if !p.curTokenIs(token.IDENT) && !p.curTokenIs(token.INT) && !p.curTokenIs(token.FLOAT) {
		msg := fmt.Sprintf("Syntax Error:%v- expected token to be 'IDENT|INT|FLOAT', got %s instead", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
	a = append(a, &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		if !p.curTokenIs(token.IDENT) && !p.curTokenIs(token.INT) && !p.curTokenIs(token.FLOAT) {
			msg := fmt.Sprintf("Syntax Error:%v- expected token to be 'IDENT|INT|FLOAT', got %s instead", p.curToken.Pos, p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}
		a = append(a, &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal})
	}
	if !p.expectPeek(allowedPair[openPair]) {
		return nil
	}
	return a
}

// IDENT() |> IDENT()
func (p *Parser) parsePipeExpression(left ast.Expression) ast.Expression {
	expression := &ast.Pipe{
		Token: p.curToken,
		Left:  left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

// EXPRESSION => EXPRESSION
//(x, y) => x + y + 5      left expression is *TupleLiteral
//(x) => x + 5             left expression is *Identifier
// x  => x + 5             left expression is *Identifier
//()  => 5 + 5             left expression is nil
func (p *Parser) parseFatArrowFunction(left ast.Expression) ast.Expression {
	tok := token.Token{Type: token.FUNCTION, Literal: "fn"}
	fn := &ast.FunctionLiteral{Token: tok, Variadic: false}
	switch exprType := left.(type) {
	case nil:
		//no argument.
	case *ast.Identifier:
		// single argument.
		fn.Parameters = append(fn.Parameters, exprType)
	case *ast.TupleLiteral:
		// a list of arguments(maybe one element tuple, or multiple elements tuple).
		for _, v := range exprType.Members {
			switch param := v.(type) {
			case *ast.Identifier:
				fn.Parameters = append(fn.Parameters, param)
			default:
				msg := fmt.Sprintf("Syntax Error:%v- Arrow function expects a list of identifiers as arguments", param.Pos())
				p.errors = append(p.errors, msg)
				return nil
			}
		}
	default:
		msg := fmt.Sprintf("Syntax Error:%v- Arrow function expects identifiers as arguments", exprType.Pos())
		p.errors = append(p.errors, msg)
		return nil
	}

	p.nextToken()
	if p.curTokenIs(token.LBRACE) { //if it's block, we use parseBlockStatement
		fn.Body = p.parseBlockStatement()
	} else { //not block, we use parseStatement
		/* Note here, if we use parseExpressionStatement, then below is not correct:
		    (x) => return x  //error: no prefix parse functions for 'RETURN' found
		so we need to use parseStatement() here
		*/
		fn.Body = &ast.BlockStatement{
			Statements: []ast.Statement{
				p.parseStatement(),
			},
		}
	}
	return fn
}

func (p *Parser) parseUsingStatement() *ast.UsingStmt {
	usingStmt := &ast.UsingStmt{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	expr := p.parseExpression(LOWEST)
	if _, ok := expr.(*ast.AssignExpression); !ok {
		msg := fmt.Sprintf("Syntax Error:%v- Using should be followed by an assignment expression", p.curToken.Pos)
		p.errors = append(p.errors, msg)
		return nil
	}
	usingStmt.Expr = expr.(*ast.AssignExpression)

	if !p.curTokenIs(token.RPAREN) {
		msg := fmt.Sprintf("Syntax Error:%v- expected token to be ')', got %s instead.", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	usingStmt.Block = p.parseBlockStatement()

	return usingStmt
}

/*
	query_expression : from_clause query_body
	from_clause : FROM identifier IN expression
	query_body : query_body_clause* select_or_group_clause query_continuation?
	query_body_clause: from_clause| let_clause | where_clause | combined_join_clause | orderby_clause
	where_clause : WHERE expression
	combined_join_clause : JOIN  identifier IN expression ON expression EQUALS expression (INTO identifier)?
	orderby_clause : ORDERBY ordering (','  ordering)*
	ordering : expression (ASCENDING | DESCENDING)?
	select_or_group_clause : SELECT expression | GROUP expression BY expression
	query_continuation : INTO identifier query_body
*/

// query_expression : from_clause query_body
func (p *Parser) parseLinqExpression() ast.Expression {
	queryExpr := &ast.QueryExpr{Token: p.curToken}

	//from_clause
	queryExpr.From = p.parseFromExpression()
	p.nextToken()

	//query_body
	queryExpr.QueryBody = p.parseQueryBodyExpression()

	if queryExpr.QueryBody.(*ast.QueryBodyExpr).Expr == nil {
		msg := fmt.Sprintf("Syntax Error:%v- Linq query must be ended with 'select' or 'group'.", p.curToken.Pos)
		p.errors = append(p.errors, msg)
		return nil
	}

	return queryExpr
}

//from_clause : FROM identifier IN expression
func (p *Parser) parseFromExpression() ast.Expression {
	fromExpr := &ast.FromExpr{}

	//FROM
	fromExpr.Token = p.curToken

	//identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	fromExpr.Var = p.curToken.Literal

	//IN
	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	//expression
	fromExpr.Expr = p.parseExpression(LOWEST)

	return fromExpr
}

//query_body : query_body_clause* select_or_group_clause query_continuation?
func (p *Parser) parseQueryBodyExpression() ast.Expression {
	queryBodyExpr := &ast.QueryBodyExpr{}

	//------------------------
	// query_body_clause*
	//------------------------
	//query_body_clause: from_clause | let_clause | where_clause | combined_join_clause | orderby_clause
	queryBodyExpr.QueryBody = []ast.Expression{}
	for p.curTokenIs(token.FROM) || p.curTokenIs(token.LET) || p.curTokenIs(token.WHERE) || p.curTokenIs(token.JOIN) || p.curTokenIs(token.ORDERBY) {
		queryBodyExpr.QueryBody = append(queryBodyExpr.QueryBody, p.parseQueryBodyClauseExpr())
		p.nextToken()
	}

	//select_or_group_clause
	if p.curTokenIs(token.SELECT) {
		queryBodyExpr.Expr = p.parseSelectExpression()
	} else if p.curTokenIs(token.GROUP) {
		queryBodyExpr.Expr = p.parseGroupExpression()
	}

	//------------------------
	// query_continuation?
	//------------------------
	//query_continuation : INTO identifier query_body
	if p.peekTokenIs(token.INTO) {
		p.nextToken()
		queryBodyExpr.QueryContinuation = p.parseQueryContinuationExpression()
	}

	return queryBodyExpr
}

//query_body_clause: from_clause| let_clause | where_clause | combined_join_clause | orderby_clause
func (p *Parser) parseQueryBodyClauseExpr() ast.Expression {
	queryBodyClauseExpr := &ast.QueryBodyClauseExpr{}

	switch p.curToken.Type {
	case token.FROM:
		//from_clause
		queryBodyClauseExpr.Expr = p.parseFromExpression()
	case token.LET:
		//let_clause
		expr := &ast.AssignExpression{}
		expr.Token = token.Token{Type: token.ASSIGN, Literal: "="}
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		expr.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		if !p.expectPeek(token.ASSIGN) {
			return nil
		}
		p.nextToken()
		expr.Value = p.parseExpression(LOWEST)
		queryBodyClauseExpr.Expr = expr
	case token.WHERE:
		//where_clause
		queryBodyClauseExpr.Expr = p.parseWhereExpression()
	case token.JOIN:
		//combined_join_clause
		queryBodyClauseExpr.Expr = p.parseJoinExpression()
	case token.ORDERBY:
		//orderby_clause
		queryBodyClauseExpr.Expr = p.parseOrderByExpression()
	default:
		return nil
	}

	return queryBodyClauseExpr
}

//query_continuation : INTO identifier query_body
func (p *Parser) parseQueryContinuationExpression() ast.Expression {
	exp := &ast.QueryContinuationExpr{}

	//INTO
	exp.Token = p.curToken //'into'

	//identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	exp.Var = p.curToken.Literal
	p.nextToken()

	//query_body
	exp.Expr = p.parseQueryBodyExpression()

	return exp
}

//SELECT expression
func (p *Parser) parseSelectExpression() ast.Expression {
	exp := &ast.SelectExpr{}

	//SELECT
	exp.Token = p.curToken
	p.nextToken()

	//expression
	exp.Expr = p.parseExpression(LOWEST)

	return exp
}

//GROUP expression BY expression
func (p *Parser) parseGroupExpression() ast.Expression {
	exp := &ast.GroupExpr{}

	//GROUP
	exp.Token = p.curToken
	p.nextToken()

	//expression
	exp.GrpExpr = p.parseExpression(LOWEST)

	//BY
	if !p.expectPeek(token.BY) {
		return nil
	}
	p.nextToken()

	//expression
	exp.ByExpr = p.parseExpression(LOWEST)

	return exp
}

//where_clause : WHERE expression
func (p *Parser) parseWhereExpression() ast.Expression {
	exp := &ast.WhereExpr{}

	//WHERE
	exp.Token = p.curToken
	p.nextToken()

	//expression
	exp.Expr = p.parseExpression(LOWEST)

	return exp
}

//combined_join_clause : JOIN identifier IN expression ON expression EQUALS expression (INTO identifier)?
func (p *Parser) parseJoinExpression() ast.Expression {
	exp := &ast.JoinExpr{}

	//JOIN
	exp.Token = p.curToken

	//identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	exp.JoinVar = p.curToken.Literal

	//IN
	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()

	//expression
	exp.InExpr = p.parseExpression(LOWEST)
	p.nextToken()

	//ON
	if !p.curTokenIs(token.ON) {
		return nil
	}
	p.nextToken()

	exp.OnExpr = p.parseExpression(LOWEST)
	p.nextToken()

	//EQUALS
	if !p.curTokenIs(token.EQUALS) {
		return nil
	}
	p.nextToken()

	//expression
	exp.EqualExpr = p.parseExpression(LOWEST)
	p.nextToken()

	//(INTO identifier)?
	if p.curTokenIs(token.INTO) {
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		exp.IntoVar = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	return exp
}

//orderby_clause : ORDERBY ordering (','  ordering)*
func (p *Parser) parseOrderByExpression() ast.Expression {
	exp := &ast.OrderExpr{}

	//ORDERBY
	exp.Token = p.curToken // 'orderby'
	p.nextToken()

	//ordering (','  ordering)*
	exp.Ordering = []ast.Expression{}
	exp.Ordering = append(exp.Ordering, p.parseOrderingExpr())

	if !p.peekTokenIs(token.COMMA) {
		return exp
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		exp.Ordering = append(exp.Ordering, p.parseOrderingExpr())
	}

	return exp
}

//ordering : expression (ASCENDING | DESCENDING)?
func (p *Parser) parseOrderingExpr() ast.Expression {
	exp := &ast.OrderingExpr{IsAscending: true, HasSortOrder: false, Var: p.curToken.Literal}

	//expression
	exp.Expr = p.parseExpression(LOWEST)

	//(ASCENDING | DESCENDING)?
	if p.peekTokenIs(token.ASCENDING) {
		exp.HasSortOrder = true
		exp.IsAscending = true

		p.nextToken()
		exp.OrderToken = p.curToken
	} else if p.peekTokenIs(token.DESCENDING) {
		exp.HasSortOrder = true
		exp.IsAscending = false

		p.nextToken()
		exp.OrderToken = p.curToken
	}

	return exp
}

//let add = async fn(a, b)  { a + b }
//let add = async (a, b) => { a + b }
func (p *Parser) parseAsyncLiteral() ast.Expression {
	p.nextToken() //skip 'async'
	if !p.curTokenIs(token.FUNCTION) && !p.curTokenIs(token.LPAREN) {
		msg := fmt.Sprintf("Syntax Error:%v- async should be followed by a function or lambda, got %s instead.", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	fnLiteral := p.parseExpression(LOWEST).(*ast.FunctionLiteral)
	fnLiteral.Async = true

	return fnLiteral
}

//async fn add(a, b) { a + b }
func (p *Parser) parseAsyncStatement() ast.Statement {
	p.nextToken() //skip 'async'
	if !p.curTokenIs(token.FUNCTION) {
		msg := fmt.Sprintf("Syntax Error:%v- async should be followed by a function, got %s instead.", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	FnStmt := p.parseFunctionStatement().(*ast.FunctionStatement)
	FnStmt.FunctionLiteral.Async = true

	return FnStmt
}

//await add(1, 2)
//await obj.xxx(params)
func (p *Parser) parseAwaitExpression() ast.Expression {
	expr := &ast.AwaitExpr{Token: p.curToken}

	p.nextToken()
	expr.Call = p.parseExpressionStatement().Expression

	//check call type
	switch expr.Call.(type) {
	case *ast.CallExpression:
		return expr
	case *ast.MethodCallExpression:
		return expr
	default:
		msg := fmt.Sprintf("Syntax Error:%v- await keyword can only be used on function/method calls!", p.curToken.Pos)
		p.errors = append(p.errors, msg)
		return nil
	}

	return expr
}

// dt/2018-01-01T12:01:00/, dt//, dt/2018-01-01 12:01:00/, ...
func (p *Parser)parseDateTime() ast.Expression {
	expr := &ast.DateTimeExpr{Token: p.curToken}

	if p.curToken.Literal == "" { //e.g. dtime = dt//
		return expr
	}

	is := &ast.InterpolatedString{Token: p.curToken, Value: p.curToken.Literal, ExprMap: make(map[byte]ast.Expression)}

	key := "0"[0]
	for {
		if p.curTokenIs(token.LBRACE) {
			p.nextToken()
			expr := p.parseExpression(LOWEST)
			is.ExprMap[key] = expr
			key++
		}

		p.nextInterpToken2()
		if p.curTokenIs(token.DATETIME) {
			break
		}
	}

	expr.Pattern = is
	return expr
}

//service name on "addrs" { block }
func (p *Parser) parseServiceStatement() *ast.ServiceStatement {
	stmt := &ast.ServiceStatement{
		Token:   p.curToken,
		Methods: make(map[string]*ast.FunctionStatement),
	}

	stmt.Doc = p.lineComment

	if !p.expectPeek(token.IDENT) { //service name
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ON) { //service name
		return nil
	}

	if !p.expectPeek(token.STRING) { //service address
		return nil
	}

	stmt.Debug = strings.HasSuffix(p.curToken.Literal, ":debug")
	if stmt.Debug {
		stmt.Addr = strings.TrimSuffix(p.curToken.Literal, ":debug")
	} else {
		stmt.Addr = p.curToken.Literal
	}

	p.nextToken()

	stmt.Block = p.parseServiceBody(stmt)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	stmt.SrcEndToken = p.curToken
	//fmt.Printf("service statement=%s\n", stmt.String())

	for k, v := range stmt.Methods {
		annoLen := len(v.Annotations)
		if annoLen != 1 {
			msg := fmt.Sprintf("Syntax Error:%v- function(%s)'s annotation count not one", p.curToken.Pos, k)
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	return stmt
}

func (p *Parser) parseServiceBody(s *ast.ServiceStatement) *ast.BlockStatement {
	stmts := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

	p.nextToken() //skip '{'
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseServiceStmt(s)
		if stmt != nil {
			stmts.Statements = append(stmts.Statements, stmt)
		}

		p.nextToken()
	}

	if p.peekTokenIs(token.EOF) && !p.curTokenIs(token.RBRACE) {
		pos := p.peekToken.Pos
		pos.Col += 1
		msg := fmt.Sprintf("Syntax Error:%v- expected next token to be '}', got EOF instead. Block should end with '}'.", pos)
		p.errors = append(p.errors, msg)
	}

	return stmts
}

func (p *Parser) parseServiceStmt(s *ast.ServiceStatement) ast.Statement {
	var annos []*ast.AnnotationStmt

	//parse Annotation
	for p.curTokenIs(token.AT) {
		anno := &ast.AnnotationStmt{Token: p.curToken, Attributes: map[string]ast.Expression{}}

		if !p.expectPeek(token.IDENT) {
			return nil
		}
		if p.curToken.Literal != "route" {
			msg := fmt.Sprintf("Syntax Error:%v- expected token to be 'route', got '%s' instead", p.curToken.Pos, p.curToken.Literal)
			p.errors = append(p.errors, msg)
			return nil
		}
		anno.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if p.peekTokenIs(token.LPAREN) {
			p.nextToken()
		}

		for {
			if !p.expectPeek(token.IDENT) {
				return nil
			}
			key := p.curToken.Literal

			if !p.expectPeek(token.ASSIGN) {
				return nil
			}
			p.nextToken()
			value := p.parseExpression(LOWEST)
			anno.Attributes[key] = value
			p.nextToken()
			if !p.curTokenIs(token.COMMA) {
				break
			}
		}

		if !p.curTokenIs(token.RPAREN) {
			msg := fmt.Sprintf("Syntax Error:%v- expected token to be ')', got '%s' instead", p.curToken.Pos, p.curToken.Type)
			p.errors = append(p.errors, msg)
			return nil
		}

		p.nextToken()
		annos = append(annos, anno)
	} //end for

	if p.curToken.Type != token.FUNCTION {
		msg := fmt.Sprintf("Syntax Error:%v- expected token to be 'fn', got %s instead.", p.curToken.Pos, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	r := p.parseFunctionStatement().(*ast.FunctionStatement)
	r.Annotations = annos
	r.IsServiceAnno = true
	s.Methods[r.Name.Value] = r

	return r
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	if t != token.EOF {
		msg := fmt.Sprintf("Syntax Error:%v- no prefix parse functions for '%s' found", p.curToken.Pos, t)
		p.errors = append(p.errors, msg)
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) nextToken() {
	p.lineComment = nil

	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()

	var list []*ast.Comment
	for p.curToken.Type == token.COMMENT {
		//if p.curToken.Literal[0] != '#' {
		if p.isDocLine(p.curToken.Pos.Line) {
			comment := &ast.Comment{Token: p.curToken, Text: p.curToken.Literal}
			list = append(list, comment)
		}
		p.curToken = p.peekToken
		p.peekToken = p.l.NextToken()
	}
	if list != nil {
		p.lineComment = &ast.CommentGroup{List: list}
	}
}

func (p *Parser) nextInterpToken() {
	p.curToken = p.l.NextInterpToken()
	p.peekToken = p.l.NextToken()
}

// for date-time literal use
func (p *Parser) nextInterpToken2() {
	p.curToken = p.l.NextInterpToken2()
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	pos := p.fixPosCol()
	msg := fmt.Sprintf("Syntax Error:%v- expected next token to be %s, got %s instead", pos, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

//Is the line document line or not
func (p *Parser) isDocLine(lineNo int) bool {
	if len(FileLines) == 0 {
		return false
	}

	lineSlice := FileLines[lineNo-1 : lineNo]
	lineStr := strings.TrimLeft(lineSlice[0], "\t ")
	if len(lineStr) > 0 {
		if lineStr[0] == '#' || lineStr[0] == '/' {
			return true
		}
	}
	return false
}

//fix position column(for error report)
func (p *Parser) fixPosCol() token.Position {
	pos := p.curToken.Pos
	if p.curToken.Type == token.STRING || p.curToken.Type == token.ISTRING {
		pos.Col = pos.Col + len(p.curToken.Literal) + 2 //2: two double/single quote(s)
	} else {
		pos.Col = pos.Col + len(p.curToken.Literal)
	}

	return pos
}

//DEBUG ONLY
func (p *Parser) debugToken(message string) {
	fmt.Printf("%s, curToken = %s, peekToken = %s\n", message, p.curToken.Literal, p.peekToken.Literal)
}

func (p *Parser) debugNode(message string, node ast.Node) {
	fmt.Printf("%s, Node = %s\n", message, node.String())
}

//stupid method to convert 'some'(not all) unicode number to ascii number
func convertNum(numStr string) string {
	var out bytes.Buffer
	for _, c := range numStr {
		if v, ok := numMap[c]; ok {
			out.WriteRune(v)
		} else {
			out.WriteRune(c)
		}
	}
	return out.String()
}
