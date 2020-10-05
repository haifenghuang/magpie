////////////////////////////////////////////////
//                 TOKEN
////////////////////////////////////////////////
const TOKEN_ILLEGAL = "ILLEGAL"
const TOKEN_EOF     = "EOF"

// Identifiers + literals
const TOKEN_IDENT  = "IDENT" // add, foobar, x, y, ...
const TOKEN_INT    = "INT"   // 1343456
const TOKEN_STRING = "STRING"

// Operators
const TOKEN_ASSIGN   = "="
const TOKEN_PLUS     = "+"
const TOKEN_MINUS    = "-"
const TOKEN_BANG     = "!"
const TOKEN_ASTERISK = "*"
const TOKEN_SLASH    = "/"
const TOKEN_INCREMENT= "++"
const TOKEN_DECREMENT= "--"

const TOKEN_LT = "<"
const TOKEN_GT = ">"

const TOKEN_EQ     = "=="
const TOKEN_NOT_EQ = "!="

// Delimiters
const TOKEN_COMMA     = ","
const TOKEN_SEMICOLON = ";"
const TOKEN_LPAREN    = "("
const TOKEN_RPAREN    = ")"
const TOKEN_LBRACE    = "{"
const TOKEN_RBRACE    = "}"
const TOKEN_LBRAKET   = "["
const TOKEN_RBRAKET   = "]"
const TOKEN_COLON     = ":"
const TOKEN_DOT       = "."

// Keywords
const TOKEN_FUNCTION = "FUNCTION"
const TOKEN_LET      = "LET"
const TOKEN_TRUE     = "TRUE"
const TOKEN_FALSE    = "FALSE"
const TOKEN_IF       = "IF"
const TOKEN_ELSE     = "ELSE"
const TOKEN_RETURN   = "RETURN"
const TOKEN_WHILE    = "WHILE"
const TOKEN_BREAK    = "BREAK"
const TOKEN_CONTINUE = "CONTINUE"

keywords = {
	"fn":       TOKEN_FUNCTION,
	"let":      TOKEN_LET,
	"true":     TOKEN_TRUE,
	"false":    TOKEN_FALSE,
	"if":       TOKEN_IF,
	"else":     TOKEN_ELSE,
	"return":   TOKEN_RETURN,
	"while":    TOKEN_WHILE,
	"break":    TOKEN_BREAK,
	"continue": TOKEN_CONTINUE
}

fn lookupIdent(ident) {
	if keywords.exists(ident) {
		return keywords[ident]
	}
	return TOKEN_IDENT
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

fn lex_nextToken() {
	let tok
	skipWhitespace()

	if lex_ch == "=" {
		if peekChar() == "=" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_EQ, ch + lex_ch)
		} else {
			tok = newToken(TOKEN_ASSIGN, lex_ch)
		}
	} elif lex_ch == "+" {
		if peekChar() == "+" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_INCREMENT, ch + lex_ch)
		} else {
			tok = newToken(TOKEN_PLUS, lex_ch)
		}
	} elif lex_ch == "-" {
		if peekChar() == "-" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_DECREMENT, ch + lex_ch)
		} else {
			tok = newToken(TOKEN_MINUS, lex_ch)				
		}
	} elif lex_ch == "!" {
		if peekChar() == "=" {
			c = lex_ch
			readChar()
			tok = newToken(TOKEN_NOT_EQ, c + lex_ch)
		} else {
			tok = newToken(TOKEN_BANG, lex_ch)
		}
	} elif lex_ch == "/" {
		tok = newToken(TOKEN_SLASH, lex_ch)
	} elif lex_ch == "*" {
		tok = newToken(TOKEN_ASTERISK, lex_ch)
	} elif lex_ch == "<" {
		tok = newToken(TOKEN_LT, lex_ch)
	} elif lex_ch == ">" {
		tok = newToken(TOKEN_GT, lex_ch)
	} elif lex_ch == ";" {
		tok = newToken(TOKEN_SEMICOLON, lex_ch)
	} elif lex_ch == "(" {
		tok = newToken(TOKEN_LPAREN, lex_ch)
	} elif lex_ch == ")" {
		tok = newToken(TOKEN_RPAREN, lex_ch)
	} elif lex_ch == "," {
		tok = newToken(TOKEN_COMMA, lex_ch)
	} elif lex_ch == "{" {
		tok = newToken(TOKEN_LBRACE, lex_ch)
	} elif lex_ch == "}" {
		tok = newToken(TOKEN_RBRACE, lex_ch)
	} elif lex_ch == "[" {
		tok = newToken(TOKEN_LBRAKET, lex_ch)
	} elif lex_ch == "]" {
		tok = newToken(TOKEN_RBRAKET, lex_ch)
	} elif lex_ch == ":" {
		tok = newToken(TOKEN_COLON, lex_ch)
	} elif lex_ch == "\"" {
		tok = newToken(TOKEN_STRING, readString())
	} elif lex_ch == "." {
		tok = newToken(TOKEN_DOT, lex_ch)
	} elif lex_ch == 0 {
		tok = newToken(TOKEN_EOF, "")
	} else {
		if isLetter(lex_ch) {
			lit = readIdentifier()
			typ = lookupIdent(lit)
			return newToken(typ, lit)
		} else if isDigit(lex_ch) {
			return newToken(TOKEN_INT, readNumber())
		} else {
			tok = newToken(TOKEN_ILLEGAL, lex_ch)
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
const WhileExpressionNode      = "WhileExpression"
const BreakExpressionNode      = "BreakExpression"
const ContinueExpressionNode   = "ContinueExpression"
const AssignExpressionNode     = "AssignExpression"
const PostfixExpressionNode    = "PostfixExpression"

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

fn Ast_WhileExpression(tok, condition, block) {
	let self = {"Token": tok, "Condition":  condition, "Block": block, "Type": WhileExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "while "
	str += condition["String"]
	str += "{ "
	str += block["String"]
	str += " }"
	self["String"] = str

	return self
}

fn Ast_BreakExpression(tok) {
	return {"Token": tok, "TokenLiteral": tok["literal"], "String": tok["literal"], "Type": BreakExpressionNode}
}


fn Ast_ContinueExpression(tok) {
	return {"Token": tok, "TokenLiteral": tok["literal"], "String": tok["literal"], "Type": ContinueExpressionNode}
}

fn Ast_AssignExpression(tok, name, value) {
	let self = {"Token": tok, "Name":  name, "Value": value, "Type": AssignExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = name["String"]
	str += " = "
	str += value["String"]
	self["String"] = str

	return self
}

fn Ast_PostfixExpression(tok, left, operator) {
	let self = {"Token": tok, "Operator": operator, "Left": left,"Type": PostfixExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "( "
	str += left["String"]
	str += " " + operator + " "
	str += ")"
	self["String"] = str
	return self
}




////////////////////////////////////////////////
//                 PARSER
////////////////////////////////////////////////

const NONE         = 0
const LOWEST       = 1
const ASSIGN       = 2  // =
const EQUALS       = 3  // ==
const LESSGREATER  = 4  // > or <
const SUM          = 5  // +
const PRODUCT      = 6  // *
const PREFIX       = 7  // -X or !X
const CALL         = 8  // myFunction(X)
const INDEX        = 9  // array[index]
const INCREMENT    = 10 // ++, --

precedences = {
	TOKEN_ASSIGN:   ASSIGN,
	TOKEN_EQ:       EQUALS,
	TOKEN_NOT_EQ:   EQUALS,
	TOKEN_LT:       LESSGREATER,
	TOKEN_GT:       LESSGREATER,
	TOKEN_PLUS:     SUM,
	TOKEN_MINUS:    SUM,
	TOKEN_SLASH:    PRODUCT,
	TOKEN_ASTERISK: PRODUCT,
	TOKEN_LPAREN:   CALL,
	TOKEN_DOT:      CALL,
	TOKEN_LBRAKET:  INDEX,
	TOKEN_INCREMENT:INCREMENT,
	TOKEN_DECREMENT:INCREMENT,
}

let parser_Errors = []
let parser_curToken = newToken(TOKEN_ILLEGAL, "")
let parser_peekToken = newToken(TOKEN_ILLEGAL, "")

let isInLoop = false

fn Parser_Errors() {
	return parser_Errors
}

fn debugToken(msg) {
	printf("%s, curToken.type=%s, curToken.literal=%s, peekToken.type=%s, peekToken.literal=%s\n", msg, parser_curToken["type"], parser_curToken["literal"], parser_peekToken["type"], parser_peekToken["literal"])
}

fn parser_nextToken() {
	parser_curToken = parser_peekToken
	parser_peekToken = lex_nextToken()
}

/* Parser entry point */
fn ParseProgram() {
	parser_nextToken()
	parser_nextToken()

	let statements = []

	while parser_curToken["type"] != TOKEN_EOF {
		let stmt = parseStatement()
		if stmt != nil {
			statements.push(stmt)
		}
		parser_nextToken()
	}

	return Ast_Program(statements)
}

fn parseStatement() {
	if parser_curToken["type"] == TOKEN_LET {
		return parseLetStatement()
	} elif parser_curToken["type"] == TOKEN_RETURN {
		return parseReturnStatement()
	} else {
		return parseExpressionStatement()
	}
}

fn parseLetStatement() {
	let curToken = parser_curToken
	if !expectPeek(TOKEN_IDENT) {
		return nil
	}

	let name = Ast_Identifier(parser_curToken, parser_curToken["literal"])

	if !expectPeek(TOKEN_ASSIGN) {
		return nil
	}

	parser_nextToken()

	let value = parseExpression(LOWEST)
	if peekTokenIs(TOKEN_SEMICOLON) {
		parser_nextToken()
	}

	return Ast_LetStatement(curToken, name, value)
}

fn parseReturnStatement() {
	let curToken = parser_curToken

	parser_nextToken()

	let returnValue = parseExpression(LOWEST)

	if peekTokenIs(TOKEN_SEMICOLON) {
		parser_nextToken()
	}

	return Ast_ReturnStatement(curToken, returnValue)
}

fn parseExpressionStatement() {
	let curToken = parser_curToken

	let expr = parseExpression(LOWEST)

	if peekTokenIs(TOKEN_SEMICOLON) {
		parser_nextToken()
	}

	return Ast_ExpressionStatement(curToken, expr)
}

/* most important function of the Parser */
fn parseExpression(precedence) {
	let leftExp

	//prefix expression
	if parser_curToken["type"] == TOKEN_IDENT {
		leftExp = parseIdentifier()
	} elif parser_curToken["type"] == TOKEN_INT {
		leftExp = parseIntegerLiteral()
	} elif parser_curToken["type"] == TOKEN_BANG {
		leftExp = parsePrefixExpression()
	} elif parser_curToken["type"] == TOKEN_MINUS {
		leftExp = parsePrefixExpression()
	} elif parser_curToken["type"] == TOKEN_INCREMENT {
		leftExp = parsePrefixExpression()
	} elif parser_curToken["type"] == TOKEN_DECREMENT {
		leftExp = parsePrefixExpression()
	} elif parser_curToken["type"] == TOKEN_TRUE {
		leftExp = parseBoolean()
	} elif parser_curToken["type"] == TOKEN_FALSE {
		leftExp = parseBoolean()
	} elif parser_curToken["type"] == TOKEN_LPAREN {
		leftExp = parseGroupedExpression()
	} elif parser_curToken["type"] == TOKEN_IF {
		leftExp = parseIfExpression()
	} elif parser_curToken["type"] == TOKEN_FUNCTION {
		leftExp = parseFunctionLiteral()
	} elif parser_curToken["type"] == TOKEN_STRING {
		leftExp = parseStringLiteral()
	} elif parser_curToken["type"] == TOKEN_LBRAKET {
		leftExp = parseArrayLiteral()
	} elif parser_curToken["type"] == TOKEN_LBRACE {
		leftExp = parseHashLiteral()
	} elif parser_curToken["type"] == TOKEN_WHILE {
		leftExp = parseWhileExpression()
	} elif parser_curToken["type"] == TOKEN_BREAK {
		if isInLoop {
			leftExp = parseBreakExpression()	
		} else {
			leftExp = parseBreakWithoutLoopContext()
		}
	} elif parser_curToken["type"] == TOKEN_CONTINUE {
		if isInLoop {
			leftExp = parseContinueExpression()
		} else {
			leftExp = parseContinueWithoutLoopContext()
		}
	} else {
		noPrefixParseFnError(parser_curToken["type"])
		return nil
	}

	while !peekTokenIs(SEMICOLON) && precedence < peekPrecedence() {
		//infix expression
		if parser_peekToken["type"] == TOKEN_PLUS || parser_peekToken["type"] == TOKEN_MINUS ||
		   parser_peekToken["type"] == TOKEN_SLASH || parser_peekToken["type"] == TOKEN_ASTERISK ||
		   parser_peekToken["type"] == TOKEN_EQ || parser_peekToken["type"] == TOKEN_NOT_EQ ||
		   parser_peekToken["type"] == TOKEN_LT || parser_peekToken["type"] == TOKEN_GT {

			parser_nextToken()
			leftExp = parseInfixExpression(leftExp)
		} elif parser_peekToken["type"] == TOKEN_LPAREN {
			parser_nextToken()
			leftExp = parseCallExpression(leftExp)
		} elif parser_peekToken["type"] == TOKEN_DOT {
			parser_nextToken()
			leftExp = parseMethodCallExpression(leftExp)
		} elif parser_peekToken["type"] == TOKEN_LBRAKET {
			parser_nextToken()
			leftExp = parseIndexExpression(leftExp)
		} elif parser_peekToken["type"] == TOKEN_ASSIGN {
			parser_nextToken()
			leftExp = parseAssignExpression(leftExp)
		} elif parser_peekToken["type"] == TOKEN_INCREMENT {
			parser_nextToken()
			leftExp = parsePostfixExpression(leftExp)
		} elif parser_peekToken["type"] == TOKEN_DECREMENT {
			parser_nextToken()
			leftExp = parsePostfixExpression(leftExp)
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
	let val = curTokenIs(TOKEN_TRUE)
	return Ast_Boolean(parser_curToken, val)
}

fn parseGroupedExpression() {
	parser_nextToken()
	let exp = parseExpression(LOWEST)

	if !expectPeek(TOKEN_RPAREN) {
		return nil
	}

	return exp
}

fn parseIfExpression() {
	let curToken = parser_curToken

	if !expectPeek(TOKEN_LPAREN) {
		return nil
	}

	parser_nextToken()
	let condition = parseExpression(LOWEST)

	if !expectPeek(TOKEN_RPAREN) {
		return nil
	}

	if !expectPeek(TOKEN_LBRACE) {
		return nil
	}

	let consequence = parseBlockStatement()

	let alternative = nil
	if peekTokenIs(TOKEN_ELSE) {
		parser_nextToken()

		if !expectPeek(TOKEN_LBRACE) {
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
	while !curTokenIs(TOKEN_RBRACE) {
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

	if !expectPeek(TOKEN_LPAREN) {
		return nil
	}

	let parameters = parseFunctionParameters()

	if !expectPeek(TOKEN_LBRACE) {
		return nil
	}

	let body = parseBlockStatement()

	return Ast_FunctionLiteral(curToken, parameters, body)
}

fn parseFunctionParameters() {
	let identifiers = []

	if peekTokenIs(TOKEN_RPAREN) {
		parser_nextToken()
		return identifiers
	}

	parser_nextToken()
	let ident = Ast_Identifier(parser_curToken, parser_curToken["literal"])
	identifiers.push(ident)

	while peekTokenIs(TOKEN_COMMA) {
		parser_nextToken()
		parser_nextToken()
		let ident2 = Ast_Identifier(parser_curToken, parser_curToken["literal"])
		identifiers.push(ident2)
	}

	if !expectPeek(TOKEN_RPAREN) {
		return nil
	}

	return identifiers
}

fn parseCallExpression(function) {
	let curToken = parser_curToken

	let arguments = parseExpressionList(TOKEN_RPAREN)

	return Ast_CallExpression(curToken, function, arguments)
}

fn parseMethodCallExpression(obj) {
	let curToken = parser_curToken

	parser_nextToken()
	let name = Ast_Identifier(parser_curToken, parser_curToken["literal"])

	let call
	if !peekTokenIs(TOKEN_LPAREN) {
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

	let elements = parseExpressionList(TOKEN_RBRAKET)

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

	while peekTokenIs(TOKEN_COMMA) {
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

	if !expectPeek(TOKEN_RBRAKET) {
		return nil
	}

	return Ast_IndexExpression(curToken, left, index)
}

fn parseHashLiteral() {
	let curToken = parser_curToken

	let pairs = {}
	while !peekTokenIs(TOKEN_RBRACE) {
		parser_nextToken()

		let key = parseExpression(LOWEST)
		if !expectPeek(TOKEN_COLON) {
			return nil
		}

		parser_nextToken()
		let value = parseExpression(LOWEST)

		let keyStr = key["String"]
		pairs[keyStr] = value

		if peekTokenIs(TOKEN_RBRACE) {
			break
		}

		if !expectPeek(TOKEN_COMMA) {
			return nil
		}
	}

	if !expectPeek(TOKEN_RBRACE) {
		return nil
	}

	return Ast_HashLiteral(curToken, pairs)
}

fn parseWhileExpression() {
	let curToken = parser_curToken

	isInLoop = true

	parser_nextToken()
	let exprStmt = parseExpressionStatement()
	let condition = exprStmt["Expression"]

	if peekTokenIs(TOKEN_RPAREN) {
		parser_nextToken()
	}

	let block = nil
	if peekTokenIs(TOKEN_LBRACE) {
		parser_nextToken()
		block = parseBlockStatement()
	} else {
		let msg = sprintf("Syntax Error: while loop must be followed by a '{'.")
		parser_Errors.push(msg)
		return nil
	}

	isInLoop = false

	return Ast_WhileExpression(curToken, condition, block)
}

fn parseBreakWithoutLoopContext() {
	let msg = sprintf("Syntax Error: 'break' outside of loop context")
	parser_Errors.push(msg)

	return parseBreakExpression()
}

fn parseBreakExpression() {
	return Ast_BreakExpression(parser_curToken)
}

fn parseContinueWithoutLoopContext() {
	let msg = sprintf("Syntax Error: 'continue' outside of loop context")
	p.errors = append(p.errors, msg)

	return parseContinueExpression()
}

fn parseContinueExpression() {
	return Ast_ContinueExpression(parser_curToken)
}

fn parseAssignExpression(left) {
	let curToken = parser_curToken

	let name = left

	parser_nextToken()
	let value = parseExpression(LOWEST)

	return Ast_AssignExpression(curToken, name, value)
}

fn parsePostfixExpression(left) {
	return Ast_PostfixExpression(parser_curToken, left, parser_curToken["literal"])
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


fn ResetScope(scope, name, val) {
	let store = scope["store"]
	let obj = store[name]
	if obj != nil {
		store[name] = val
	}

	if obj == nil {
		let outer = scope["outer"]
		if outer != nil {
			return ResetScope(outer, name, val)
		}
	}

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
const BREAK_OBJ        = "BREAK"
const CONTINUE_OBJ     = "CONTINUE"

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

	Array_Inspect(self)

	self["CallMethod"] = fn (scope, methodName, methodArgs) {
		if methodName == "push" {
			return Array_Push(self, scope, methodArgs)
		} elif methodName == "pop" {
			return Array_Pop(self, scope)
		} else {
			let msg = sprintf("undefined method '%s' for object %s", methodName, ARRAY_OBJ)
			return ErrorObj(msg)
		}
	}

	return self
}

fn Array_Inspect(self) {
	let elems = self["Elements"]

	let str = ""

	let items = []
	for elem in elems {
		items.push(elem["Inspect"])
	}

	str += "["
	str += strings.join(items, ", ")
	str += "]"
	self["Inspect"] = str
}


fn Array_Push(self, scope, methodArgs) {
	let arg = methodArgs[0]

	let elements = self["Elements"]
	elements.push(arg)

	//recalculate 'Inspect'
	Array_Inspect(self)

	return self
}

fn Array_Pop(self, scope) {
	let elements = self["Elements"]
	elements.pop()

	//recalculate 'Inspect'
	Array_Inspect(self)

	return self
}

fn HashPair(key, value) {
	return {"Key": key, "Value": value}
}

fn HashObj(pairs) {
	//Pairs map[HashKey]HashPair
	//Here I use obj["Inspect"] as the hash key.

	let self = {"Type": HASH_OBJ, "Pairs": pairs}

	Hash_Inspect(self)

	self["CallMethod"] = fn (scope, methodName, methodArgs) {
		if methodName == "exists" || methodName == "has" {
			return Hash_Exists(self, scope, methodArgs)
		} elif methodName == "push" {
			return Hash_Push(self, scope, methodArgs)
		} elif methodName == "pop" {
			return Hash_Pop(self, scope, methodArgs)
		} elif methodName == "keys" {
			return Hash_Keys(self, scope)
		} elif methodName == "values" {
			return Hash_Values(self, scope)
		} else {
			let msg = sprintf("undefined method '%s' for object %s", methodName, HASH_OBJ)
			return ErrorObj(msg)
		}
	}

	return self
}

fn Hash_Inspect(self) {
	let pairs = self["Pairs"]

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
}

fn Hash_Exists(self, scope, methodArgs) {
	let arg0 = methodArgs[0]
	let pairs = self["Pairs"]

	for notUsed, pair in pairs {
		let k = pair["Key"]
		if k["Inspect"] == arg0["Inspect"] {
			return TRUE
		}
 	}
	return FALSE
}


fn Hash_Push(self, scope, methodArgs) {
	let pairs = self["Pairs"]

	let arg0 = methodArgs[0]
	let arg1 = methodArgs[1]

	let key =  arg0["Inspect"]
	pairs[key] = HashPair(arg0, arg1)

	//recalculate 'Inspect'
	Hash_Inspect(self)

	return self
}

fn Hash_Pop(self, scope, methodArgs) {
	let key = methodArgs[0]

	let pairs = self["Pairs"]
	for notUsed, pair in pairs {
		let k = pair["Key"]
		if k["Inspect"] == key["Inspect"] {
			pairs.pop(k["Inspect"])

			//recalculate 'Inspect'
			Hash_Inspect(self)

			return pair["Value"]
		}
 	}

	return NilObj()
}

fn Hash_Keys(self, scope) {
	let keys = []

	let pairs = self["Pairs"]
	for notUsed, pair in pairs {
		keys.push(pair["Key"])
	}
	return ArrayObj(keys)
}

fn Hash_Values(self, scope) {
	let values = []

	let pairs = self["Pairs"]
	for notUsed, pair in pairs {
		values.push(pair["Value"])
	}
	return ArrayObj(values)
}


fn BreakObj() {
	return {"Type": BREAK_OBJ, "Inspect": "break", "CallMethod": nil}
}

fn ContinueObj() {
	return {"Type": CONTINUE_OBJ, "Inspect": "continue", "CallMethod": nil}
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
const BREAK = BreakObj()
const CONTINUE = ContinueObj()

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
	} elif nodeType == WhileExpressionNode {
		return evalWhileExpression(node, scope)
	} elif nodeType == BreakExpressionNode {
		return BREAK
	} elif nodeType == ContinueExpressionNode {
		return CONTINUE
	} elif nodeType == AssignExpressionNode {
		return evalAssignExpression(node, scope)
	} elif nodeType == PostfixExpressionNode {
		let p_left = Eval(node["Left"], scope)
		if p_left["Type"] == ERROR_OBJ {
			return left
		}
		return evalPostfixExpression(node["Operator"], p_left)
	}

	return nil
}

// Private functions from here

fn evalProgram(program, scope) {
	let result = NIL

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

			if rt == BREAK_OBJ || rt == CONTINUE_OBJ {
				return
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
	} elif operator == "++" {
		return evalIncrementPrefixOperatorExpression(right)
	} elif operator == "--" {
		return evalDecrementPrefixOperatorExpression(right)
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

fn evalIncrementPrefixOperatorExpression(right) {
	if right["Type"] != INTEGER_OBJ {
		let msg = sprintf("unsupported operator for prefix expression:'++' and type: %s", right["Type"])
		return ErrorObj(msg)
	}

	let value = right["Value"]
	right["Value"] = value + 1
	right["Inspect"] = sprintf("%d", value + 1)

	return IntegerObj(value + 1)
}

fn evalDecrementPrefixOperatorExpression(right) {
	if right["Type"] != INTEGER_OBJ {
		let msg = sprintf("unsupported operator for prefix expression:'--' and type: %s", right["Type"])
		return ErrorObj(msg)
	}

	let value = right["Value"]
	right["Value"] = value - 1
	right["Inspect"] = sprintf("%d", value - 1)

	return IntegerObj(value - 1)
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
		return nativeBoolToBooleanObject(leftVal < rightVal)
	} elif operator == ">" {
		return nativeBoolToBooleanObject(leftVal > rightVal)
	} elif operator == "==" {
		return nativeBoolToBooleanObject(leftVal == rightVal)
	} elif operator == "!=" {
		return nativeBoolToBooleanObject(leftVal != rightVal)
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

fn evalWhileExpression(node, scope) {
	let innerScope = NewEnclosedScope(scope)

	let condition = Eval(node["Condition"], innerScope)
	if condition["Type"] == ERROR_OBJ {
		return condition
	}

	let result = NIL
	while isTruthy(condition) {
		result = Eval(node["Block"], innerScope)
		if result["Type"] == ERROR_OBJ {
			return result
		}

		if result["Type"] == BREAK_OBJ {
			break
		}

		if result["Type"] == CONTINUE_OBJ {
			continue
		}

		if result["Type"] == RETURN_VALUE_OBJ {
			if result["Value"] != nil {
				return result
			}
			break
		}

		condition = Eval(node["Condition"], innerScope)
		if condition["Type"] == ERROR_OBJ {
			return condition
		}
	}

	if result == nil || result["Type"] == BREAK_OBJ || result["Type"] == CONTINUE_OBJ {
		return NIL
	}

	return result
}

fn evalAssignExpression(node, scope) {
	let val = Eval(node["Value"], scope)
	if val["Type"] != ERROR_OBJ {
		let name = node["Name"]
		let v = ResetScope(scope, name["String"], val)
		if v != nil {
			return v
		}

		let msg = sprintf("unknown identifier: '%s' is not defined", name["String"])
		return  ErrorObj(msg)
	}
	return val
}


fn evalPostfixExpression(operator, left) {
	if operator == "++" {
		return evalIncrementPostfixOperatorExpression(left)
	} elif operator == "--" {
		return evalDecrementPostfixOperatorExpression(left)
	}

	let msg = sprintf("unsupported operator for postfix expression:'%s' and type: %s", operator, left["Type"])
	return ErrorObj(msg)

}

fn evalIncrementPostfixOperatorExpression(left) {
	if left["Type"] != INTEGER_OBJ {
		let msg = sprintf("unsupported operator for post expression:'++' and type: %s", left["Type"])
		return ErrorObj(msg)
	}

	let left_value = left["Value"]
	let returnVal = IntegerObj(left_value)
	left["Value"] = left_value + 1
	left["Inspect"] = sprintf("%d", left_value + 1)
	return returnVal
}

fn evalDecrementPostfixOperatorExpression(left) {
	if left["Type"] != INTEGER_OBJ {
		let msg = sprintf("unsupported operator for post expression:'--' and type: %s", left["Type"])
		return ErrorObj(msg)
	}

	let left_value = left["Value"]
	let returnVal = IntegerObj(left_value)
	left["Value"] = left_value - 1
	left["Inspect"] = sprintf("%d", left_value - 1)
	return returnVal
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
	println(h.keys())
	println(h.values())
	println(h.exists("Name"))

	h.push("Height", "164cm")
	println(h)
	println(h.pop("Name"))
	println(h)

	let strHello = "Hello World"
	println(strHello.lower())
	println(strHello.upper())

	let i = 10
	println("i = ", i)
	while (i > 3) {
    	i--
    	if (i == 7) { continue }
    	if (i == 5) { break }
    	println(i)
	}
	``

	/* lexer */
	NewLexer(source)

/* for debugging : print each tokens */
//	do {
//		n = nextToken()
//		if n["type"] == TOKEN_EOF {
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

