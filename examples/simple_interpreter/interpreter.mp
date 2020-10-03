////////////////////////////////////////////////
//                 TOKEN
////////////////////////////////////////////////
const ILLEGAL = "ILLEGAL"
const EOF     = "EOF"

// Identifiers + literals
const IDENT  = "IDENT" // add, foobar, x, y, ...
const INT    = "INT"   // 1343456
const STRING = "STRING"

// Operators
const ASSIGN   = "="
const PLUS     = "+"
const MINUS    = "-"
const BANG     = "!"
const ASTERISK = "*"
const SLASH    = "/"

const LT = "<"
const GT = ">"

const EQ     = "=="
const NOT_EQ = "!="

// Delimiters
const COMMA     = ","
const SEMICOLON = ";"
const LPAREN    = "("
const RPAREN    = ")"
const LBRACE    = "{"
const RBRACE    = "}"
const LBRAKET   = "["
const RBRAKET   = "]"
const COLON     = ":"
const DOT       = "."

// Keywords
const FUNCTION = "FUNCTION"
const LET      = "LET"
const TRUE     = "TRUE"
const FALSE    = "FALSE"
const IF       = "IF"
const ELSE     = "ELSE"
const RETURN   = "RETURN"

keywords = {
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

fn lookupIdent(ident) {
	if keywords.exists(ident) {
		return keywords[ident]
	}
	return IDENT
}

fn newToken(t, ch) {
	return {"type": t, "literal": ch}
}



////////////////////////////////////////////////
//                 LEXER
////////////////////////////////////////////////
lex_input = ""        // input 
lex_position = 0     // current position in input (points to current char)
lex_readPosition = 0 // current reading position in input (after current char)
lex_ch = " "         // current char under examination

//create a new lexer
fn NewLexer(input) {
	lex_input = input
	readChar()
}

fn isDigit(c) {
	return "0" <= c && c <= "9"
}

fn isLetter(c) {
	return "a" <= c && c <= "z" || "A" <= c && c <= "Z" || c == "_"
}

fn readChar() {
	if lex_readPosition >= len(lex_input) {
		lex_ch = 0
	} else {
		lex_ch = lex_input[lex_readPosition]
	}

	lex_position = lex_readPosition
	lex_readPosition = lex_readPosition + 1
}

fn peekChar() {
	if lex_readPosition >= len(lex_input) {
		return 0
	}

	return lex_input[lex_readPosition]
}

fn readIdentifier() {
	pos = lex_position
	while (isLetter(lex_ch)) {
		readChar()
	}

	return lex_input[pos : lex_position]
}

fn readNumber() {
	p = lex_position
	while isDigit(lex_ch) {
		readChar()
	}
	return lex_input[p : lex_position]
}

fn skipWhitespace() {
	while lex_ch == " " || lex_ch == "\t" || lex_ch == "\n" || lex_ch == "\r" {
		readChar()
	}
}

fn readString() {
	position = lex_position + 1
	for {
		readChar()
		if lex_ch == "\"" {
			break
		}
	}

	return lex_input[position : lex_position]
}

fn nextToken() {
	let tok
	skipWhitespace()

	if lex_ch == "=" {
		if peekChar() == "=" {
			ch = lex_ch
			readChar()
			tok = newToken(EQ, ch + lex_ch)
		} else {
			tok = newToken(ASSIGN, lex_ch)
		}
	} elif lex_ch == "+" {
		tok = newToken(PLUS, lex_ch)
	} elif lex_ch == "-" {
		tok = newToken(MINUS, lex_ch)
	} elif lex_ch == "!" {
		if peekChar() == "=" {
			c = lex_ch
			readChar()
			tok = newToken(NOT_EQ, c + lex_ch)
		} else {
			tok = newToken(BANG, lex_ch)
		}
	} elif lex_ch == "/" {
		tok = newToken(SLASH, lex_ch)
	} elif lex_ch == "*" {
		tok = newToken(ASTERISK, lex_ch)
	} elif lex_ch == "<" {
		tok = newToken(LT, lex_ch)
	} elif lex_ch == ">" {
		tok = newToken(GT, lex_ch)
	} elif lex_ch == ";" {
		tok = newToken(SEMICOLON, lex_ch)
	} elif lex_ch == "(" {
		tok = newToken(LPAREN, lex_ch)
	} elif lex_ch == ")" {
		tok = newToken(RPAREN, lex_ch)
	} elif lex_ch == "," {
		tok = newToken(COMMA, lex_ch)
	} elif lex_ch == "{" {
		tok = newToken(LBRACE, lex_ch)
	} elif lex_ch == "}" {
		tok = newToken(RBRACE, lex_ch)
	} elif lex_ch == "[" {
		tok = newToken(LBRAKET, lex_ch)
	} elif lex_ch == "]" {
		tok = newToken(RBRAKET, lex_ch)
	} elif lex_ch == ":" {
		tok = newToken(COLON, lex_ch)
	} elif lex_ch == "\"" {
		tok = newToken(STRING, readString())
	} elif lex_ch == "." {
		tok = newToken(DOT, lex_ch)
	} elif lex_ch == 0 {
		tok = newToken(EOF, "")
	} else {
		if isLetter(lex_ch) {
			lit = readIdentifier()
			typ = lookupIdent(lit)
			return newToken(typ, lit)
		} else if isDigit(lex_ch) {
			return newToken(INT, readNumber())
		} else {
			tok = newToken(ILLEGAL, lex_ch)
		}
	}

	readChar()
	return tok
}



////////////////////////////////////////////////
//                 AST
////////////////////////////////////////////////
/* All AST node should have below properties:
    "Type":            Node Type
    "TokenLiteral":    Token Literal
    "String":          Mainly for debug
*/

//AST node type
const ProgramNode              = "Program"
const IdentifierNode           = "Identifier"
const LetStatementNode         = "LetStatement"
const ReturnStatementNode      = "ReturnStatement"
const ExpressionStatementNode  = "ExpressionStatement"
const IntegerLiteralNode       = "IntegerLiteral"
const PrefixExpressionNode     = "PrefixExpression"
const InfixExpressionNode      = "InfixExpression"
const BooleanNode              = "Boolean"
const IfExpressionNode         = "IfExpression"
const BlockStatementNode       = "BlockStatement"
const FunctionLiteralNode      = "FunctionLiteral"
const CallExpressionNode       = "CallExpression"
const MethodCallExpressionNode = "MethodCallExpression"
const StringLiteralNode        = "StringLiteral"
const ArrayLiteralNode         = "ArrayLiteral"
const IndexExpressionNode      = "IndexExpression"
const HashLiteralNode          = "HashLiteral"

fn Ast_Program(statements) {
	let self = {"Statements": statements, "Type": ProgramNode}

	self["TokenLiteral"] = ""
	if len(statements) > 0 {
		let stmt = statements[0]
		self["TokenLiteral"] = stmt["TokenLiteral"]
	}

	let str = ""
	for stmt in statements {
		str += stmt["String"]
	}
	self["String"] = str

	return self
}

fn Ast_Identifier(tok, value) {
	let self = {"Token" : tok, "Value": value, "Type": IdentifierNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = value

	return self
}

fn Ast_LetStatement(tok, name, value) {
	let self = {"Token": tok, "Name": name, "Value": value, "Type": LetStatementNode}

	self["TokenLiteral"] = tok["literal"]

	let str = ""
	str += tok["literal"] + " "
	str += name["TokenLiteral"]
	str += " = "
	if value != nil {
		str += value["String"]
	}
	str += ";"
	self["String"] = str

	return self
}

fn Ast_ReturnStatement(tok, rv) {
	let self = {"Token": tok, "ReturnValue": rv, "Type": ReturnStatementNode}

	self["TokenLiteral"] = tok["literal"]

	let str = tok["literal"] + " "
	if rv != nil {
		str += rv["TokenLiteral"]
	}
	str += ";"
	self["String"] = str

	return self
}

fn Ast_ExpressionStatement(tok, expr) {
	let self = {"Token": tok, "Expression": expr, "Type": ExpressionStatementNode}

	self["TokenLiteral"] = tok["literal"]

	self["String"] = ""
	if expr != nil {
		self["String"] = expr["String"]
	}

	return self
}

fn Ast_IntegerLiteral(tok, value) {
	let self = {"Token": tok, "Value": value, "Type": IntegerLiteralNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = tok["literal"]

	return self
}

fn Ast_PrefixExpression(tok, operator, right) {
	let self = {"Token": tok, "Operator": operator, "Right": right, "Type": PrefixExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "("
	str += operator
	str += right["String"]
	str += ")"
	self["String"] = str

	return self
}

fn Ast_InfixExpression(tok, left, operator, right) {
	let self = {"Token": tok, "Operator": operator, "Left": left, "Right":right, "Type": InfixExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "("
	str += left["String"]
	str += " " + operator + " "
	str += right["String"]
	str += ")"
	self["String"] = str

	return self
}

fn Ast_Boolean(tok, value) {
	let self = {"Token": tok, "Value": value, "Type": BooleanNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = tok["literal"]

	return self
}

fn Ast_IfExpression(tok, condition, consequence, alternative) {
	let self = {"Token": tok, "Condition": condition, "Consequence": consequence, "Alternative": alternative}

	self["Type"] = IfExpressionNode
	self["TokenLiteral"] = tok["literal"]

	let str = "if"
	str += condition["String"]
	str += " "
	str += consequence["String"]
	if alternative != nil {
		str += "else"
		str += alternative["String"]
	}
	self["String"] = str

	return self
}

fn Ast_BlockStatement(tok, statements) {
	let self = {"Token": tok, "Statements": statements, "Type": BlockStatementNode}

	self["TokenLiteral"] = tok["literal"]

	let str = ""
	for stmt in statements {
		str += stmt["String"]
	}
	self["String"] = str

	return self
}

fn Ast_FunctionLiteral(tok, parameters, body) {
	let self = {"Token": tok, "Parameters": parameters, "Body": body, "Type": FunctionLiteralNode}

	self["TokenLiteral"] = tok["literal"]

	let params = []
	for param in parameters {
		params.push(param["String"])
	}

	let str = ""
	str += tok["literal"]
	str += "("
	str += strings.join(params, ", ")
	str += ") "
	str += body["String"]
	self["String"] = str

	return self
}

fn Ast_CallExpression(tok, function, arguments) {
	let self = {"Token": tok, "Function": function, "Arguments": arguments, "Type": CallExpressionNode}

	self["TokenLiteral"] = tok["literal"]


	let str = ""
	let args = []
	for arg in arguments {
		args.push(arg["String"])
	}

	str += function["String"]
	str += "("
	str += strings.join(args, ", ")
	str += ")"
	self["String"] = str

	return self
}

fn Ast_MethodCallExpression(tok, obj, call) {
	let self = {"Token": tok, "Object": obj, "Call": call, "Type": MethodCallExpressionNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = obj["String"] + "." + call["String"]

	return self
}

fn Ast_StringLiteral(tok, value) {
	let self = {"Token": tok, "Value": value, "Type": StringLiteralNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = tok["literal"]

	return self
}

fn Ast_ArrayLiteral(tok, elements) {
	let self = {"Token": tok, "Elements": elements, "Type": ArrayLiteralNode}

	self["TokenLiteral"] = tok["literal"]

	let str = ""
	let elems = []
	for elem in elements {
		elems.push(elem["String"])
	}
	str += "["
	str += strings.join(elems, ", ")
	str += "]"
	self["String"] = str

	return self
}

fn Ast_IndexExpression(tok, left, index) {
	let self = {"Token": tok, "Left":  left, "Index": index, "Type": IndexExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "("
	str += left["String"]
	str += "["
	str += index["String"]
	str += "])"
	self["String"] = str

	return self
}

fn Ast_HashLiteral(tok, pairs) {
	let self = {"Token": tok, "Pairs": pairs, "Type": HashLiteralNode}

	self["TokenLiteral"] = tok["literal"]

	let str = ""

	let kvPairs = []
	for key, value in pairs {
		kvPairs.push(key["String"] + ":" + value["String"])
	}
	str += "{"
	str += strings.join(kvPairs, ", ")
	str += "}"
	self["String"] = str

	return self
}



////////////////////////////////////////////////
//                 PARSER
////////////////////////////////////////////////

const NONE         = 0
const LOWEST       = 1
const EQUALS       = 2  // ==
const LESSGREATER  = 3  // > or <
const SUM          = 4  // +
const PRODUCT      = 5  // *
const PREFIX       = 6  // -X or !X
const CALL         = 7  // myFunction(X)
const INDEX        = 8  // array[index]

precedences = {
	EQ:       EQUALS,
	NOT_EQ:   EQUALS,
	LT:       LESSGREATER,
	GT:       LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
	LPAREN:   CALL,
	DOT:      CALL,
	LBRAKET:  INDEX,
}

let parser_Errors = []
let parser_curToken = newToken(ILLEGAL, "")
let parser_peekToken = newToken(ILLEGAL, "")

fn Parser_Errors() {
	return parser_Errors
}

fn debugToken(msg) {
	printf("%s, curToken.type=%s, curToken.literal=%s, peekToken.type=%s, peekToken.literal=%s\n", msg, parser_curToken["type"], parser_curToken["literal"], parser_peekToken["type"], parser_peekToken["literal"])
}

fn parser_nextToken() {
	parser_curToken = parser_peekToken
	parser_peekToken = nextToken()
}

/* Parser entry point */
fn ParseProgram() {
	parser_nextToken()
	parser_nextToken()

	let statements = []

	while parser_curToken["type"] != EOF {
		let stmt = parseStatement()
		if stmt != nil {
			statements.push(stmt)
		}
		parser_nextToken()
	}

	return Ast_Program(statements)
}

fn parseStatement() {
	if parser_curToken["type"] == LET {
		return parseLetStatement()
	} elif parser_curToken["type"] == RETURN {
		return parseReturnStatement()
	} else {
		return parseExpressionStatement()
	}
}

fn parseLetStatement() {
	let curToken = parser_curToken
	if !expectPeek(IDENT) {
		return nil
	}

	let name = Ast_Identifier(parser_curToken, parser_curToken["literal"])

	if !expectPeek(ASSIGN) {
		return nil
	}

	parser_nextToken()

	let value = parseExpression(LOWEST)
	if peekTokenIs(SEMICOLON) {
		parser_nextToken()
	}

	return Ast_LetStatement(curToken, name, value)
}

fn parseReturnStatement() {
	let curToken = parser_curToken

	parser_nextToken()

	let returnValue = parseExpression(LOWEST)

	if peekTokenIs(SEMICOLON) {
		parser_nextToken()
	}

	return Ast_ReturnStatement(curToken, returnValue)
}

fn parseExpressionStatement() {
	let curToken = parser_curToken

	let expr = parseExpression(LOWEST)

	if peekTokenIs(SEMICOLON) {
		parser_nextToken()
	}

	return Ast_ExpressionStatement(curToken, expr)
}

/* most important function of the Parser */
fn parseExpression(precedence) {
	let leftExp

	//prefix expression
	if parser_curToken["type"] == IDENT {
		leftExp = parseIdentifier()
	} elif parser_curToken["type"] == INT {
		leftExp = parseIntegerLiteral()
	} elif parser_curToken["type"] == BANG {
		leftExp = parsePrefixExpression()
	} elif parser_curToken["type"] == MINUS {
		leftExp = parsePrefixExpression()
	} elif parser_curToken["type"] == TRUE {
		leftExp = parseBoolean()
	} elif parser_curToken["type"] == FALSE {
		leftExp = parseBoolean()
	} elif parser_curToken["type"] == LPAREN {
		leftExp = parseGroupedExpression()
	} elif parser_curToken["type"] == IF {
		leftExp = parseIfExpression()
	} elif parser_curToken["type"] == FUNCTION {
		leftExp = parseFunctionLiteral()
	} elif parser_curToken["type"] == STRING {
		leftExp = parseStringLiteral()
	} elif parser_curToken["type"] == LBRAKET {
		leftExp = parseArrayLiteral()
	} elif parser_curToken["type"] == LBRACE {
		leftExp = parseHashLiteral()
	} 
	else {
		noPrefixParseFnError(parser_curToken["type"])
		return nil
	}

	while !peekTokenIs(SEMICOLON) && precedence < peekPrecedence() {
		//infix expression
		if parser_peekToken["type"] == PLUS || parser_peekToken["type"] == MINUS ||
		   parser_peekToken["type"] == SLASH || parser_peekToken["type"] == ASTERISK ||
		   parser_peekToken["type"] == EQ || parser_peekToken["type"] == NOT_EQ ||
		   parser_peekToken["type"] == LT || parser_peekToken["type"] == GT {

			parser_nextToken()
			leftExp = parseInfixExpression(leftExp)
		} elif parser_peekToken["type"] == LPAREN {
			parser_nextToken()
			leftExp = parseCallExpression(leftExp)
		} elif parser_peekToken["type"] == DOT {
			parser_nextToken()
			leftExp = parseMethodCallExpression(leftExp)
		} elif parser_peekToken["type"] == LBRAKET {
			parser_nextToken()
			leftExp = parseIndexExpression(leftExp)
		} else {
			return leftExp
		}
	}

	return leftExp
}

fn parseIdentifier() {
	return Ast_Identifier(parser_curToken, parser_curToken["literal"])
}

fn parseIntegerLiteral() {
	return Ast_IntegerLiteral(parser_curToken, int(parser_curToken["literal"]))
}

fn parsePrefixExpression() {
	let curToken = parser_curToken

	parser_nextToken()

	let right = parseExpression(PREFIX)

	return Ast_PrefixExpression(curToken, curToken["literal"], right)
}

fn parseInfixExpression(left) {
	let curToken = parser_curToken

	let precedence = curPrecedence()
	parser_nextToken()
	let right = parseExpression(precedence)

	return Ast_InfixExpression(curToken, left, curToken["literal"], right)
}

fn parseBoolean() {
	let val = curTokenIs(TRUE)
	return Ast_Boolean(parser_curToken, val)
}

fn parseGroupedExpression() {
	parser_nextToken()
	let exp = parseExpression(LOWEST)

	if !expectPeek(RPAREN) {
		return nil
	}

	return exp
}

fn parseIfExpression() {
	let curToken = parser_curToken

	if !expectPeek(LPAREN) {
		return nil
	}

	parser_nextToken()
	let condition = parseExpression(LOWEST)

	if !expectPeek(RPAREN) {
		return nil
	}

	if !expectPeek(LBRACE) {
		return nil
	}

	let consequence = parseBlockStatement()

	let alternative = nil
	if peekTokenIs(ELSE) {
		parser_nextToken()

		if !expectPeek(LBRACE) {
			return nil
		}

		alternative = parseBlockStatement()
	}

	return Ast_IfExpression(curToken, condition, consequence, alternative)
}

fn parseBlockStatement() {
	let curToken = parser_curToken
	parser_nextToken()

	let statements = []
	while !curTokenIs(RBRACE) {
		let stmt = parseStatement()
		if stmt != nil {
			statements.push(stmt)
		}
		parser_nextToken()
    }

	return Ast_BlockStatement(curToken, statements)
}

fn parseFunctionLiteral() {
	let curToken = parser_curToken

	if !expectPeek(LPAREN) {
		return nil
	}

	let parameters = parseFunctionParameters()

	if !expectPeek(LBRACE) {
		return nil
	}

	let body = parseBlockStatement()

	return Ast_FunctionLiteral(curToken, parameters, body)
}

fn parseFunctionParameters() {
	let identifiers = []

	if peekTokenIs(RPAREN) {
		parser_nextToken()
		return identifiers
	}

	parser_nextToken()
	let ident = Ast_Identifier(parser_curToken, parser_curToken["literal"])
	identifiers.push(ident)

	while peekTokenIs(COMMA) {
		parser_nextToken()
		parser_nextToken()
		let ident2 = Ast_Identifier(parser_curToken, parser_curToken["literal"])
		identifiers.push(ident2)
	}

	if !expectPeek(RPAREN) {
		return nil
	}

	return identifiers
}

fn parseCallExpression(function) {
	let curToken = parser_curToken

	let arguments = parseExpressionList(RPAREN)

	return Ast_CallExpression(curToken, function, arguments)
}

fn parseMethodCallExpression(obj) {
	let curToken = parser_curToken

	parser_nextToken()
	let name = Ast_Identifier(parser_curToken, parser_curToken["literal"])

	let call
	if !peekTokenIs(LPAREN) {
		call = parseExpression(LOWEST)
	} else {
		parser_nextToken()
		call = parseCallExpression(name)
	}
	return Ast_MethodCallExpression(curToken, obj, call)
}

fn parseStringLiteral() {
	return Ast_StringLiteral(parser_curToken, parser_curToken["literal"])
}

fn parseArrayLiteral() {
	let curToken = parser_curToken

	let elements = parseExpressionList(RBRAKET)

	return Ast_ArrayLiteral(curToken, elements)
}

fn parseExpressionList(end) {
	let lst = []

	if peekTokenIs(end) {
		parser_nextToken()
		return lst
	}

	parser_nextToken()
	let exp = parseExpression(LOWEST)
	lst.push(exp)

	while peekTokenIs(COMMA) {
		parser_nextToken()
		parser_nextToken()
		let exp2 = parseExpression(LOWEST)
		lst.push(exp2)
	}

	if !expectPeek(end) {
		return nil
	}

	return lst
}

fn parseIndexExpression(left) {
	let curToken = parser_curToken

	parser_nextToken()
	let index = parseExpression(LOWEST)

	if !expectPeek(RBRAKET) {
		return nil
	}

	return Ast_IndexExpression(curToken, left, index)
}

fn parseHashLiteral() {
	let curToken = parser_curToken

	let pairs = {}
	while !peekTokenIs(RBRACE) {
		parser_nextToken()

		let key = parseExpression(LOWEST)
		if !expectPeek(COLON) {
			return nil
		}

		parser_nextToken()
		let value = parseExpression(LOWEST)

		let keyStr = key["String"]
		pairs[keyStr] = value

		if peekTokenIs(RBRACE) {
			break
		}

		if !expectPeek(COMMA) {
			return nil
		}
	}

	if !expectPeek(RBRACE) {
		return nil
	}

	return Ast_HashLiteral(curToken, pairs)
}

fn curTokenIs(t) {
	return parser_curToken["type"] == t
}

fn peekTokenIs(t) {
	return parser_peekToken["type"] == t
}

fn peekPrecedence() {
	if precedences.has(parser_peekToken["type"]) {
		return precedences[parser_peekToken["type"]]
	}

	return LOWEST
}

fn curPrecedence() {
	if precedences.has(parser_curToken["type"]) {
		return precedences[parser_curToken["type"]]
	}

	return LOWEST
}

fn expectPeek(t) {
	if peekTokenIs(t) {
		parser_nextToken()
		return true
	}

	peekError(t)
	return false
}

fn peekError(t) {
	let msg = "expected next token to be " + parser_peekToken["type"] + ", got " + t + " instead"
	parser_Errors.push(msg)
}

fn noPrefixParseFnError(t) {
	let msg = "no prefix parse function for " + t +" found"
	parser_Errors.push(msg)
}



////////////////////////////////////////////////
//                 SCOPE
////////////////////////////////////////////////
fn NewScope() {
	return {"store": {}, "outer": nil}
}

fn NewEnclosedScope(outerScope) {
	let scope = NewScope()
	scope["outer"] = outerScope
	return scope
}

fn GetFromScope(scope, name) {
    let store = scope["store"]
	let obj = store[name]
	if obj == nil {
		let outer = scope["outer"]
		if outer != nil {
			return GetFromScope(outer, name)
		}
	}

	return obj
}

fn SetIntoScope(scope, name, val) {
	let store = scope["store"]
	store[name] = val
	return val
}


////////////////////////////////////////////////
//                 OBJECT
////////////////////////////////////////////////
/* Object Type */
const INTEGER_OBJ      = "INTEGER"
const BOOLEAN_OBJ      = "BOOLEAN"
const NIL_OBJ          = "NIL"
const RETURN_VALUE_OBJ = "RETURN_VALUE"
const ERROR_OBJ        = "ERROR"
const FUNCTION_OBJ     = "FUNCTION"
const STRING_OBJ       = "STRING"
const BUILTIN_OBJ      = "BUILTIN"
const ARRAY_OBJ        = "ARRAY"
const HASH_OBJ         = "HASH"

/* Every Object should have below properties:
    "Type":       Object Type
    "Inspect":    Mainly for debug
    "CallMethod": Object method which has the following prototype
         (scope, methodName, methodArgs) -> Object"
*/

fn IntegerObj(value) {
	let self = {"Type": INTEGER_OBJ, "Inspect": sprintf("%d", value), "Value": value}

	self["CallMethod"] = fn (scope, methodName, methodArgs) {
		if methodName == "next" {
			return Integer_Next(value)

		} elif methodName == "prev" {
			return Integer_Prev(value)
		} else {
			let msg = sprintf("undefined method '%s' for object %s", methodName, INTEGER_OBJ)
			return ErrorObj(msg)
		}
	}

	return self
}

fn Integer_Next(value) {
	return IntegerObj(value + 1)
}

fn Integer_Prev(value) {
	return IntegerObj(value - 1)
}

fn BooleanObj(value) {
	let self = {"Type": BOOLEAN_OBJ, "Inspect": sprintf("%t", value), "Value": value}

	self["CallMethod"] = fn (scope, methodName, methodArgs) {
		if methodName == "toYesNo" {
			return Boolean_ToYesNo(value)

		} elif methodName == "toTrueFalse" {
			return Boolean_ToTrueFalse(value)
		} else {
			let msg = sprintf("undefined method '%s' for object %s", methodName, BOOLEAN_OBJ)
			return ErrorObj(msg)
		}
	}

	return self
}

fn Boolean_ToYesNo(value) {
	if value {
		return StringObj("yes")
	}
	return StringObj("no")
}

fn Boolean_ToTrueFalse(value) {
	if value {
		return StringObj("true")
	}
	return StringObj("false")
}

fn NilObj() {
	return {"Type": NIL_OBJ, "Inspect": "nil", "CallMethod": nil}
}

fn ReturnObj(value) {
	return {"Type": RETURN_VALUE_OBJ, "Inspect": value["Inspect"], "Value": value, "CallMethod": nil}
}

fn ErrorObj(message) {
	return {"Type": ERROR_OBJ, "Inspect": "Runtime Error:" + message, "Message": message, "CallMethod": nil}
}

fn FunctionObj(parameters, body, scope) {
	let self = {"Type": FUNCTION_OBJ, "Parameters": parameters, "Body": body, "Scope": scope, "CallMethod": nil}

	let str = ""

	let params = []
	for p in parameters {
		params.push(p["String"])
	}

	str += "fn"
	str += "("
	str += strings.join(params, ", ")
	str += ") {\n"
	str += body["String"]
	str += "\n}"
	self["Inspect"] = str

	return self
}

fn StringObj(value) {
	let self = {"Type": STRING_OBJ, "Inspect": value, "Value": value}

	self["CallMethod"] = fn (scope, methodName, methodArgs) {
		if methodName == "upper" {
			return String_Upper(value)

		} elif methodName == "lower" {
			return String_Lower(value)
		} elif methodName == "trim" {
			return String_Trim(value)
		} else {
			let msg = sprintf("undefined method '%s' for object %s", methodName, BOOLEAN_OBJ)
			return ErrorObj(msg)
		}
	}

	return self
}

fn String_Lower(value) {
	return StringObj(value.lower())
}

fn String_Upper(value) {
	return StringObj(value.upper())
}

fn String_Trim(value) {
	return StringObj(value.trim())
}

fn BuiltinObj(builtinFn) {
	return {"Type": BUILTIN_OBJ, "Inspect": "<builtin>", "BuiltinFn": builtinFn, "CallMethod": nil}
}

fn ArrayObj(elements) {
	let self = {"Type": ARRAY_OBJ, "Elements": elements}

	//recalculate "Inspect"
	let str = ""

	let items = []
	for elem in self["Elements"] {
		items.push(elem["Inspect"])
	}

	str += "["
	str += strings.join(items, ", ")
	str += "]"
	self["Inspect"] = str

	self["CallMethod"] = fn (scope, methodName, methodArgs) {
		if methodName == "push" {
			let arg0 = methodArgs[0]
			Array_Push(self, scope, arg0)
		} elif methodName == "pop" {
			Array_Pop(self, scope)
		} else {
			let msg = sprintf("undefined method '%s' for object %s", methodName, ARRAY_OBJ)
			return ErrorObj(msg)
		}
	}

	return self
}

fn Array_Push(self, scope, arg) {
	let elements = self["Elements"]

	elements.push(arg)
	self["Elements"] = elements

	//recalculate "Inspect"
	let str = ""

	let items = []
	for elem in elements {
		items.push(elem["Inspect"])
	}

	str += "["
	str += strings.join(items, ", ")
	str += "]"
	self["Inspect"] = str
}

fn Array_Pop(self, scope) {
	let elements = self["Elements"]

	elements.pop()
	self["Elements"] = elements

	let str = ""

	let items = []
	for elem in elements {
		items.push(elem["Inspect"])
	}

	str += "["
	str += strings.join(items, ", ")
	str += "]"
	self["Inspect"] = str
}

fn HashPair(key, value) {
	return {"Key": key, "Value": value}
}

fn HashObj(pairs) {
	//Pairs map[HashKey]HashPair
	//Here I use obj["Inspect"] as the hash key.

	let self = {"Type": HASH_OBJ, "Pairs": pairs, "CallMethod": nil}

	let str = ""

	let pairsArr = []
	for notUsed, pair in pairs {
		let k = pair["Key"]
		let v = pair["Value"]
		let s = sprintf("%s: %s", k["Inspect"], v["Inspect"])
		pairsArr.push(s)
	}

	str += "{"
	str += strings.join(pairsArr, ", ")
	str += "}"
	self["Inspect"] = str

	return self
}


////////////////////////////////////////////////
//                 BUILTINS
////////////////////////////////////////////////
/* builtins = map[string]BuiltinObj */

fn lenBuiltin(args) {
	if len(args) != 1 {
		let msg = sprintf("wrong number of arguments. got=%d, want=1", len(args))
		return ErrorObj(msg)
	}

	let arg0 = args[0]
	let arg0Type = arg0["Type"]
	if arg0Type == ARRAY_OBJ {
		return IntegerObj(len(arg0["Elements"]))
	} elif arg0Type == STRING_OBJ {
			return IntegerObj(len(arg0["Value"]))
	}

	let msg2 = sprintf("argument to `len` not supported, got %s", arg0Type)
	return ErrorObj(msg2)
}

fn printBuiltin(args) {
	for arg in args {
		print(arg["Inspect"])
	}

	return NIL
}

fn printlnBuiltin(args) {
	for arg in args {
		print(arg["Inspect"])
	}
	println()

	return NIL
}

const builtins = {
	"len":  BuiltinObj(lenBuiltin),
	"print": BuiltinObj(printBuiltin),
	"println": BuiltinObj(printlnBuiltin),
}


////////////////////////////////////////////////
//                 EVALUATOR
////////////////////////////////////////////////
/* singleton object */
const NIL  = NilObj()
const TRUE  = BooleanObj(true)
const FALSE = BooleanObj(false)

/* Evaluator entry point */
fn Eval(node, scope) {
	let nodeType = node["Type"]

	//printf("debug: node.type=%s, node=%s\n", node["Type"], node["String"])

	if nodeType == ProgramNode {
		return evalProgram(node, scope)
	} elif nodeType == ExpressionStatementNode {
		return Eval(node["Expression"], scope)
	} elif nodeType == IntegerLiteralNode {
		return IntegerObj(node["Value"])
	} elif nodeType == BooleanNode {
		return nativeBoolToBooleanObject(node["Value"])
	} elif nodeType == PrefixExpressionNode {
		let p_right = Eval(node["Right"], scope)
		if isError(p_right) {
			return p_right
		}
		return evalPrefixExpression(node["Operator"], p_right)
	} elif nodeType == InfixExpressionNode {
		let i_left = Eval(node["Left"], scope)
		if isError(i_left) {
			return i_left
		}

		let i_right = Eval(node["Right"], scope)
		if isError(i_right) {
			return i_right
		}

		return evalInfixExpression(node["Operator"], i_left, i_right)
	} elif nodeType == BlockStatementNode {
		return evalBlockStatement(node, scope)
	} elif nodeType == IfExpressionNode {
		return evalIfExpression(node, scope)
	} elif nodeType == ReturnStatementNode {
		let r_val = Eval(node["ReturnValue"], scope)
		if isError(r_val) {
			return r_val
		}
		return ReturnValueObj(val)
	} elif nodeType == LetStatementNode {
		let l_val = Eval(node["Value"], scope)
		if isError(l_val) {
			return l_val
		}

		let l_name = node["Name"]
		let name_value = l_name["Value"]
		SetIntoScope(scope, name_value, l_val)

	} elif nodeType == IdentifierNode {
		return evalIdentifier(node, scope)
	} elif nodeType == FunctionLiteralNode {
		let fn_params = node["Parameters"]
		let fn_body   = node["Body"]
		return FunctionObj(fn_params, fn_body, scope)
	} elif nodeType == CallExpressionNode {
		let call_fn = Eval(node["Function"], scope)
		if isError(call_fn) {
			return call_fn
		}

		let call_args = evalExpressions(node["Arguments"], scope)
		if len(call_args) == 1 && isError(call_args[0]) {
			return call_args[0]
		}

		return applyFunction(call_fn, call_args)
	} elif nodeType == MethodCallExpressionNode {
		return evalMethodCallExpression(node, scope)
	} elif nodeType == StringLiteralNode {
		return StringObj(node["Value"])
	} elif nodeType == ArrayLiteralNode {
		let arr_elements = evalExpressions(node["Elements"], scope)
		if len(arr_elements) == 1 && isError(arr_elements[0]) {
			return arr_elements[0]
		}
		return ArrayObj(arr_elements)
	} elif nodeType == IndexExpressionNode {
		let idx_left = Eval(node["Left"], scope)
		if isError(idx_left) {
			return idx_left
		}
		let idx_index = Eval(node["Index"], scope)
		if isError(idx_index) {
			return idx_index
		}
		return evalIndexExpression(idx_left, idx_index)
	} elif nodeType == HashLiteralNode {
		return evalHashLiteral(node, scope)
	}

	return nil
}

// Private functions from here

fn evalProgram(program, scope) {
	let result

	for statement in program["Statements"] {
		result = Eval(statement, scope)
		if result != nil {
			if result["Type"] == RETURN_VALUE_OBJ {
				return result["Value"]
			} elif result["Type"] == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

fn evalBlockStatement(block, scope) {
	let result = NIL

	for statement in block["Statements"] {
		result = Eval(statement, scope)

		if result != nil {
			let rt = result["Type"]
			if rt == RETURN_VALUE_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

fn nativeBoolToBooleanObject(input) {
	if input {
		return TRUE
	}

	return FALSE
}

fn evalPrefixExpression(operator, right) {
	if operator == "!" {
		return evalBangOperatorExpression(right)
	} elif operator =="-" {
		return evalMinusPrefixOperatorExpression(right)
	}

	let msg = sprintf("unknown operator: %s%s", operator, right["Type"])
	return ErrorObj(msg)

}

fn evalBangOperatorExpression(right) {
	if right == TRUE {
		return FALSE
	} elif right == FALSE {
		return TRUE
	} elif right == NIL {
		return TRUE
	}

	return FALSE
}

fn evalMinusPrefixOperatorExpression(right) {
	if right["Type"] != INTEGER_OBJ {
		let msg = sprintf("unknown operator: -%s", right["Type"])
		return ErrorObj(msg)
	}

	let value = right["Value"]
	return IntegerObj(-value)
}

fn evalInfixExpression(operator, left, right) {
	if left["Type"] == INTEGER_OBJ && right["Type"] == INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	} elif operator == "==" {
		return nativeBoolToBooleanObject(left == right)
	} elif operator == "!=" {
		return nativeBoolToBooleanObject(left != right)
	} elif left["Type"] == STRING_OBJ && right["Type"] == STRING_OBJ {
		return evalStringInfixExpression(operator, left, right)
	} elif left["Type"] != right["Type"] {
		let msg = sprintf("type mismatch: %s %s %s", left["Type"], operator, right["Type"])
		return ErrorObj(msg)
	}

	let msg2 = sprintf("unknown operator: %s %s %s", left["Type"], operator, right["Type"])
	return  ErrorObj(msg2)
}

fn evalIntegerInfixExpression(operator, left, right) {
	let leftVal = left["Value"]
	let rightVal = right["Value"]
	if operator == "+" {
		return IntegerObj(leftVal + rightVal)
	} elif operator == "-" {
		return IntegerObj(leftVal - rightVal)
	} elif operator == "*" {
		return IntegerObj(leftVal * rightVal)
	} elif operator == "/" {
		return IntegerObj(leftVal / rightVal)
	} elif operator == "<" {
		return IntegerObj(leftVal < rightVal)
	} elif operator == ">" {
		return IntegerObj(leftVal > rightVal)
	} elif operator == "==" {
		return IntegerObj(leftVal == rightVal)
	} elif operator == "!=" {
		return IntegerObj(leftVal != rightVal)
	}

	let msg = sprintf("unknown operator: %s %s %s", left["Type"], operator, right["Type"])
	return  ErrorObj(msg)
}

fn evalIfExpression(ie, scope) {
	let condition = Eval(ie["Condition"], scope)
	let alternative = ie["Alternative"]
	let consequence = ie["Consequence"]

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(consequence, scope)
	} elif alternative != nil {
		return Eval(alternative, scope)
	} else {
		return NIL
	}
}

fn evalIdentifier(node, scope) {
	let nodeVal = node["Value"]
	let val = GetFromScope(scope, nodeVal)
	if val != nil {
		return val
	}

	let built_in = builtins[nodeVal]
	if built_in != nil {
		return built_in
	}

	let msg = sprintf("identifier not found: " + node.Value)
	return  ErrorObj(msg)
}

fn evalExpressions(exps, scope) {
	let result = []

	for e in exps {
		let evaluated = Eval(e, scope)
		if isError(evaluated) {
			result.push(evaluated)
			return result
		}
		result.push(evaluated)
	}

	return result
}

fn evalStringInfixExpression(operator, left, right) {
	if operator != "+" {
		let msg = sprintf("unknown operator: %s %s %s", left["Type"], operator, right["Type"])
		return  ErrorObj(msg)
	}

	let lv = left["Value"]
	let rv = right["Value"]
	return StringObj(lv + rv)
}

fn evalIndexExpression(left, index) {
	if left["Type"] == ARRAY_OBJ && index["Type"] == INTEGER_OBJ {
		return evalArrayIndexExpression(left, index)
	} elif left["Type"] == HASH_OBJ {
		return evalHashIndexExpression(left, index)
	}

	let msg = sprintf("index operator not supported: %s", left["Type"])
	return  ErrorObj(msg)
}

fn evalArrayIndexExpression(arr, index) {
	let max = len(arr["Elements"]) - 1
	let idx = index["Value"]

	if idx < 0 || idx > max {
		return NIL
	}

	let item = arr["Elements"]
	return item[idx]
}

fn evalHashIndexExpression(h, index) {
	let k = index["Inspect"]
	let pairs = h["Pairs"]
	for notUsed, pair in pairs {
		let kObj = pair["Key"]
		if k == kObj["Inspect"] {
			return pair["Value"]
		}
	}

	return NIL
}

fn evalHashLiteral(node, scope) {
	let pairs = {}

	for k, v in node["Pairs"] {
		let value = Eval(v, scope)
		if isError(value) {
			return value
		}

		let key =  StringObj(k)
		pairs[k] = HashPair(key, value)
	}

	return HashObj(pairs)
}

fn evalMethodCallExpression(node, scope) {
	let obj = Eval(node["Object"], scope)

	let call = node["Call"]
	let methName = call["Function"]["String"]

	if obj["CallMethod"] == nil {
		let msg = sprintf("undefined method '%s' for object %s", methName, obj["Type"])
		return  ErrorObj(msg)
	}

	//arguments
	let args = []
	for v in call["Arguments"] {
		args.push(Eval(v, scope))
	}
	return obj["CallMethod"](scope, methName, args)
}

fn isTruthy(obj) {
	if obj == NIL {
		return false
	} elif obj == TRUE {
		return true
	} elif obj == FALSE {
		return false
	}

	return true
}

fn isError(obj) {
	if obj != nil {
		return obj["Type"] == ERROR_OBJ
	}

	return false
}

fn applyFunction(func, args) {
	let funcType = func["Type"]
	if funcType == FUNCTION_OBJ {
		let extendedScope = extendFunctionScope(func, args)
		let evaluated = Eval(func["Body"], extendedScope)
		return unwrapReturnValue(evaluated)
	} elif funcType == BUILTIN_OBJ {
		let builtin_fn = func["BuiltinFn"]
		return builtin_fn(args)
	}

	let msg = sprintf("not a function: %s", func["Type"])
	return  ErrorObj(msg)
}

fn extendFunctionScope(func, args) {
	let s = NewEnclosedScope(func["Scope"])

	for i, param in func["Parameters"] {
		SetIntoScope(s, param["Value"], args[i])
	}

	return s
}

fn unwrapReturnValue(obj) {
	if obj["Type"] == RETURN_VALUE_OBJ {
		return obj["Value"]
	}

	return obj
}



////////////////////////////////////////////////
//                 TESTING
////////////////////////////////////////////////

#define DEBUG
#ifdef DEBUG
{
	# source = ``let five = 5
	# let ten = 10
	# println("five = ", five)
	# println("ten = ", ten)

    # let fifteen = five + ten
    # println("fifteen = ", fifteen)

	# let add = fn(x, y) {
	# 	x + y
	# }
	# let result = add(five, ten)
	# println("result = ", result)

	# let arr = [1, 2]
    # println(arr)
	# ``

	source = ``let five = 5
	let ten = 10
	println("five = ", five)
	println("ten = ", ten)
	println("six = ", five.next())
	println("nine = ", ten.prev())

	let fifteen = five + ten
	println("fifteen = ", fifteen)

	let add = fn(x, y) {
		x + y
	}
	let result = add(five, ten)
	println("add(", five, ",", ten, ") = ", result)

	let arr = [1, 2]
	println(arr)
	println("arr[1] = ", arr[1])

	arr.push(3)
	println("new arr = ", arr)
	println("new arr length = ", len(arr))
	arr.pop()
	println("new arr2 = ", arr)
	println("new arr2 length = ", len(arr))

	let h = {"Name": "HHF", "Age": 43}
	println(h)
	println("h[Age] = ", h["Age"])

	let strHello = "Hello World"
	println(strHello.lower())
	println(strHello.upper())
	``

	/* lexer */
	NewLexer(source)

/* for debugging : print each tokens */
//	do {
//		n = nextToken()
//		if n["type"] == EOF {
//			break
//		} else {
//			printf("{token.%s, \"%s\"}\n", n["type"], n["literal"])
//		}
//	}

	/* parser */
	let p = ParseProgram()
	if len(parser_Errors) != 0 {
		for errMsg in parser_Errors {
			println("\t" + errMsg + "\n")
		}
	}

	//for debugging
//	println("==========================================")
//	println(p["String"]) //print the 'program' node
//	println("==========================================")

	/* evaluator */
	let root_scope = NewScope()
	Eval(p, root_scope)
}

