////////////////////////////////////////////////
//                 POSITION
////////////////////////////////////////////////

//Position is the location of a code point in the source
fn newPosition(filename, offset, line, column) {
	//fffset: relative to entire file
	//column: relative to each line

	return {"Filename": filename, "Offset": offset, "Line": line, "Column": column}
}

fn invalidPos() {
	return newPosition("", 0, 0, 0)
}

//convert position to string
fn positionStr(position) {
	let msg = ""
	if position == nil {
		return "NoPos"
	}

	if position["Filename"] == nil || position["Filename"] == "" {
		msg = sprintf(" <%d:%d> ", position["Line"], position["Column"])
	} else {
		msg = sprintf(" <%s: %d:%d> ", position["Filename"], position["Line"], position["Column"])
	}

	return msg
}

////////////////////////////////////////////////
//                 TOKEN
////////////////////////////////////////////////
const (
	TOKEN_ILLEGAL,
	TOKEN_EOF,
	
	// Identifiers + literals
	TOKEN_IDENT,   // add, foobar, x, y, ...
	TOKEN_INT,     // 1343456
	TOKEN_FLOAT,   // 13.43456
	TOKEN_STRING,
	
	// Operators
	TOKEN_ASSIGN,      // "="
	TOKEN_PLUS,        // "+"
	TOKEN_MINUS,       // "-"
	TOKEN_BANG,        // "!"
	TOKEN_ASTERISK,    // "*"
	TOKEN_SLASH,       // "/"
	TOKEN_POWER,       // "**"
	TOKEN_INCREMENT,   // "++"
	TOKEN_DECREMENT,   // "--"
	
	TOKEN_LT,  // "<"
	TOKEN_LE,  // "<="
	TOKEN_GT,  // ">"
	TOKEN_GE,  // ">="
	
	TOKEN_EQ,     // "=="
	TOKEN_NOT_EQ, // "!="

	TOKEN_SHIFT_L,   // << (SHIFT LEFT)
	TOKEN_SHIFT_R,   // >> (SHIFT RIGHT)

	TOKEN_REGEXP,  // "/xxx/"
	TOKEN_MATCH,   // "=~""
	TOKEN_UNMATCH, // "!~""

	// Delimiters
	TOKEN_COMMA,     // ","
	TOKEN_SEMICOLON, // ";"
	TOKEN_LPAREN,    // "("
	TOKEN_RPAREN,    // ")"
	TOKEN_LBRACE,    // "{"
	TOKEN_RBRACE,    // "}"
	TOKEN_LBRAKET,   // "["
	TOKEN_RBRAKET,   // "]"
	TOKEN_COLON,     // ":"
	TOKEN_DOT,       // "."
	TOKEN_UNDERSCORE,// "_"(PlaceHolder)

	// Keywords
	TOKEN_FUNCTION, // "FUNCTION"
	TOKEN_LET,      // "LET"
	TOKEN_CONST,    // "CONST"
	TOKEN_TRUE,     // "TRUE"
	TOKEN_FALSE,    // "FALSE"
	TOKEN_IF,       // "IF"
	TOKEN_ELSE,     // "ELSE"
	TOKEN_RETURN,   // "RETURN"
	TOKEN_FOR,      // "FOR"
	TOKEN_IN,       // "IN"
	TOKEN_WHILE,    // "WHILE"
	TOKEN_BREAK,    // "BREAK"
	TOKEN_CONTINUE, // "CONTINUE"
	TOKEN_AND,      // "AND"
	TOKEN_OR,       // "OR"
	TOKEN_NIL,      // "NIL"
)

keywords = {
	"fn":       TOKEN_FUNCTION,
	"let":      TOKEN_LET,
	"const":    TOKEN_CONST,
	"true":     TOKEN_TRUE,
	"false":    TOKEN_FALSE,
	"if":       TOKEN_IF,
	"else":     TOKEN_ELSE,
	"return":   TOKEN_RETURN,
	"for":      TOKEN_FOR,
	"in":       TOKEN_IN,
	"while":    TOKEN_WHILE,
	"break":    TOKEN_BREAK,
	"continue": TOKEN_CONTINUE,
	"and":      TOKEN_AND,
	"or":       TOKEN_OR,
	"nil":      TOKEN_NIL,
}

fn lookupIdent(ident) {
	if keywords.exists(ident) {
		return keywords[ident]
	}
	return TOKEN_IDENT
}

fn newToken(t, ch, position) {
	return {"type": t, "literal": ch, "position": position}
}

fn tokenString(tok) {
	case tok in {
		TOKEN_ILLEGAL   { return "ILLEGAL" }
		TOKEN_EOF       { return "EOF" }
		TOKEN_IDENT     { return "IDENT" }
		TOKEN_INT       { return "INT" }
		TOKEN_FLOAT     { return "FLOAT" }
		TOKEN_STRING    { return "STRING" }
		TOKEN_ASSIGN    { return "=" }
		TOKEN_PLUS      { return "+" }
		TOKEN_MINUS     { return "-" }
		TOKEN_BANG      { return "!" }
		TOKEN_ASTERISK  { return "*" }
		TOKEN_SLASH     { return "/" }
		TOKEN_POWER     { return "**" }
		TOKEN_INCREMENT { return "++" }
		TOKEN_DECREMENT { return "--" }
		TOKEN_LT        { return "<" }
		TOKEN_LE        { return "<=" }
		TOKEN_GT        { return ">" }
		TOKEN_GE        { return ">=" }
		TOKEN_SHIFT_L   { return ">>" }
		TOKEN_SHIFT_R   { return "<<" }
		TOKEN_EQ        { return "==" }
		TOKEN_NOT_EQ    { return "!=" }
		TOKEN_REGEXP    { return "/REGEX/" }
		TOKEN_MATCH     { return "=~" }
		TOKEN_UNMATCH   { return "!~" }
		TOKEN_COMMA     { return "," }
		TOKEN_SEMICOLON { return ";" }
		TOKEN_LPAREN    { return "(" }
		TOKEN_RPAREN    { return ")" }
		TOKEN_LBRACE    { return "{" }
		TOKEN_RBRACE    { return "}" }
		TOKEN_LBRAKET   { return "[" }
		TOKEN_RBRAKET   { return "]" }
		TOKEN_COLON     { return ":" }
		TOKEN_DOT       { return "." }
		TOKEN_UNDERSCORE{ return "_" }
		TOKEN_FUNCTION  { return "FUNCTION" }
		TOKEN_LET       { return "LET" }
		TOKEN_CONST     { return "CONST" }
		TOKEN_TRUE      { return "TRUE" }
		TOKEN_FALSE     { return "FALSE" }
		TOKEN_IF        { return "IF" }
		TOKEN_ELSE      { return "ELSE" }
		TOKEN_RETURN    { return "RETURN" }
		TOKEN_FOR       { return "FOR" }
		TOKEN_IN        { return "IN" }
		TOKEN_WHILE     { return "WHILE" }
		TOKEN_BREAK     { return "BREAK" }
		TOKEN_CONTINUE  { return "CONTINUE" }
		TOKEN_AND       { return "AND" }
		TOKEN_OR        { return "OR" }
		TOKEN_NIL       { return "NIL" }
	}
}

////////////////////////////////////////////////
//                 LEXER
////////////////////////////////////////////////
lex_filename = ""
lex_input = ""        // input 
lex_position = 0     // current position in input (points to current char)
lex_readPosition = 0 // current reading position in input (after current char)
lex_ch = " "         // current char under examination
lex_line = 1
lex_col = 0
let lex_prevToken

//create a new lexer
fn NewLexer(input, filename) {
	lex_filename = filepath.base(filename)
	lex_input = input
	readChar()
}

fn isDigit(c) {
	return "0" <= c && c <= "9"
}

fn isLetter(c) {
	return "a" <= c && c <= "z" || "A" <= c && c <= "Z" || c == "_" || c == "$"
}

fn readChar() {
	if lex_readPosition >= len(lex_input) {
		lex_ch = 0
	} else {
		lex_ch = lex_input[lex_readPosition]
		if lex_ch == "\n" {
			lex_col = 0
			lex_line += 1
		} else {
			lex_col += 1
		}
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
	while (lex_ch != 0 && (isLetter(lex_ch) || isDigit(lex_ch))) {
		readChar()
	}

	return lex_input[pos : lex_position]
}

fn readNumber() {
	let ret = lex_ch
	readChar()

	while (lex_ch != 0 && (isDigit(lex_ch) || lex_ch == "." || lex_ch == "_")) { //support 1_000_000 format
		if lex_ch == "_" {
			readChar()
			continue
		}

		ret += lex_ch
		readChar()
	}

	return ret
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
		if lex_ch == "\"" || lex_ch == 0 {
			break
		}
	}

	return lex_input[position : lex_position]
}

fn skipComment() {
	while lex_ch != "\n" && lex_ch != 0 {
		readChar()
	}
	if lex_ch == "\n" {
		readChar()
	}
}

fn readRegExLiteral() {
	position = lex_position

	/* read until closing slash */
	for {
		readChar()
		if lex_ch == "\\" {
			// Skip escape sequence
			readChar()
		} elif lex_ch == "/" {
			// This is the closing
			literal = lex_input[position + 1 : lex_position]
			readChar() //skip the '/'

			return literal
		}
	}
}

fn lex_nextToken() {
	let tok
	skipWhitespace()

	let tok_pos = newPosition(lex_filename, lex_position, lex_line, lex_col)
	tok_pos["Column"] -= 1

	if lex_ch == "#" {
		skipComment()
		return lex_nextToken()
	}

	if lex_ch == "=" {
		if peekChar() == "=" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_EQ, ch + lex_ch, tok_pos)
		} elif peekChar() == "~" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_MATCH, ch + lex_ch, tok_pos)
		} else {
			tok = newToken(TOKEN_ASSIGN, lex_ch, tok_pos)
		}
	} elif lex_ch == "+" {
		if peekChar() == "+" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_INCREMENT, ch + lex_ch, tok_pos)
		} else {
			tok = newToken(TOKEN_PLUS, lex_ch, tok_pos)
		}
	} elif lex_ch == "-" {
		if peekChar() == "-" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_DECREMENT, ch + lex_ch, tok_pos)
		} else {
			tok = newToken(TOKEN_MINUS, lex_ch, tok_pos)				
		}
	} elif lex_ch == "!" {
		if peekChar() == "=" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_NOT_EQ, ch + lex_ch, tok_pos)
		} elif peekChar() == "~" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_UNMATCH, ch + lex_ch, tok_pos)
		} else {
			tok = newToken(TOKEN_BANG, lex_ch, tok_pos)
		}
	} elif lex_ch == "/" {
		if lex_prevToken["type"] == TOKEN_RBRACE || // impossible?
			lex_prevToken["type"] == TOKEN_RPAREN || // (a+c) / b
			lex_prevToken["type"] == TOKEN_RBRACKET || // a[3] / b
			lex_prevToken["type"] == TOKEN_IDENT || // a / b
			lex_prevToken["type"] == TOKEN_INT || // 3 / b
			lex_prevToken["type"] == TOKEN_FLOAT || // 3.5 / b
			lex_prevToken["type"] == TOKEN_FUNCTION { // e.g. fn /() - operator overloading
				tok = newToken(TOKEN_SLASH, lex_ch, tok_pos)
		} else { //regexp
			lit = readRegExLiteral()
			return newToken(TOKEN_REGEXP, lit, tok_pos)
		}

	} elif lex_ch == "*" {
		if peekChar() == "*" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_POWER, ch + lex_ch, tok_pos)
		} else {
			tok = newToken(TOKEN_ASTERISK, lex_ch, tok_pos)			
		}
	} elif lex_ch == "<" {
		if peekChar() == "=" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_LE, ch + lex_ch, tok_pos)
		} elif peekChar() == "<" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_SHIFT_L, ch + lex_ch, tok_pos) 
		} else {
			tok = newToken(TOKEN_LT, lex_ch, tok_pos)
		}
	} elif lex_ch == ">" {
		if peekChar() == "=" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_GE, ch + lex_ch, tok_pos)
		} elif peekChar() == ">" {
			ch = lex_ch
			readChar()
			tok = newToken(TOKEN_SHIFT_R, ch + lex_ch, tok_pos)
		} else {
			tok = newToken(TOKEN_GT, lex_ch, tok_pos)
		}
	} elif lex_ch == ";" {
		tok = newToken(TOKEN_SEMICOLON, lex_ch, tok_pos)
	} elif lex_ch == "(" {
		tok = newToken(TOKEN_LPAREN, lex_ch, tok_pos)
	} elif lex_ch == ")" {
		tok = newToken(TOKEN_RPAREN, lex_ch, tok_pos)
	} elif lex_ch == "," {
		tok = newToken(TOKEN_COMMA, lex_ch, tok_pos)
	} elif lex_ch == "{" {
		tok = newToken(TOKEN_LBRACE, lex_ch, tok_pos)
	} elif lex_ch == "}" {
		tok = newToken(TOKEN_RBRACE, lex_ch, tok_pos)
	} elif lex_ch == "[" {
		tok = newToken(TOKEN_LBRAKET, lex_ch, tok_pos)
	} elif lex_ch == "]" {
		tok = newToken(TOKEN_RBRAKET, lex_ch, tok_pos)
	} elif lex_ch == ":" {
		tok = newToken(TOKEN_COLON, lex_ch, tok_pos)
	} elif lex_ch == "\"" {
		tok = newToken(TOKEN_STRING, readString(), tok_pos)
	} elif lex_ch == "." {
		tok = newToken(TOKEN_DOT, lex_ch, tok_pos)
	} elif lex_ch == "_" {
		tok = newToken(TOKEN_UNDERSCORE, lex_ch, tok_pos)
	} elif lex_ch == 0 {
		tok = newToken(TOKEN_EOF, "", tok_pos)
	} else {
		if isLetter(lex_ch) {
			lit = readIdentifier()
			typ = lookupIdent(lit)
			tok = newToken(typ, lit, tok_pos, tok_pos)
			lex_prevToken = tok
			return tok
		} else if isDigit(lex_ch) {
			let num_literal = readNumber()
			if num_literal.contains(".") {
				tok = newToken(TOKEN_FLOAT, num_literal, tok_pos)
				lex_prevToken = tok
				return tok
			} else {
				tok = newToken(TOKEN_INT, num_literal, tok_pos)
				lex_prevToken = tok
				return tok
			}
		} else {
			tok = newToken(TOKEN_ILLEGAL, lex_ch)
		}
	}

	readChar()
	lex_prevToken = tok
	return tok
}



////////////////////////////////////////////////
//          AST(Abstract Syntax Tree)
////////////////////////////////////////////////
/* All AST node must have below properties:
    "Type":            Node Type
    "TokenLiteral":    Token Literal
    "String":          Mainly for debug
    "StartPos":        position of first character belonging to the AST node
	"EndPos":          position of first character immediately after the AST node
*/

//AST node type
const (
	ProgramNode,
	IdentifierNode,
	LetStatementNode,
	ConstStatementNode,
	ReturnStatementNode,
	ExpressionStatementNode,
	IntegerLiteralNode,
	FloatLiteralNode,
	PrefixExpressionNode,
	InfixExpressionNode,
	BooleanNode,
	IfExpressionNode,
	BlockStatementNode,
	FunctionStatementNode,
	FunctionLiteralNode,
	CallExpressionNode,
	MethodCallExpressionNode,
	StringLiteralNode,
	ArrayLiteralNode,
	TupleLiteralNode,
	IndexExpressionNode,
	HashLiteralNode,
	ForArrayExpressionNode,
	ForMapExpressionNode,
	ForEverLoopExpressionNode,
	WhileExpressionNode,
	BreakExpressionNode,
	ContinueExpressionNode,
	AssignExpressionNode,
	PostfixExpressionNode,
	NilExpressionNode,
	RegExLiteralNode,
)

fn Ast_Program(statements) {
	let self = {"Statements": statements, "Type": ProgramNode}

	self["TokenLiteral"] = ""
	if len(statements) > 0 {
		let stmt = statements[0]
		self["TokenLiteral"] = stmt["TokenLiteral"]
	}

	let str = ""
	for stmt in statements {
		str += stmt["String"] + "\n"
	}
	self["String"] = str

	//start & end position
	let stmtLen = len(statements)
	if stmtLen > 0 {
		self["StartPos"] = statements[0]["StartPos"]
		self["EndPos"] = statements[stmtLen-1]["EndPos"]
	} else {
		self["StartPos"] = invalidPos()
		self["EndPos"] = invalidPos()
	}

	return self
}

fn Ast_Identifier(tok, value) {
	let self = {"Token" : tok, "Value": value, "Type": IdentifierNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = value

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let valLen = len(value)
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + valLen, tokPos["Line"], tokPos["Column"] + valLen)
	self["EndPos"] = endPos

	return self
}

fn Ast_LetStatement(tok, names, values) {
	let self = {"Token": tok, "Names": names, "Values": values, "Type": LetStatementNode}

	self["TokenLiteral"] = tok["literal"]

	let str = ""
	str += tok["literal"] + " "

	//names
	let n = []
	for name in names {
		n.push(name["String"])
	}
	str += strings.join(n, ", ")

	if len(values) == 0 { // no values
		str += ";"
		self["String"] = str
		return self
	}

	str += " = "

	//values
	let v = []
	for value in values {
		v.push(value["String"])
	}
	str += strings.join(v, ", ")

	str += ";"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]

	let valLen = len(values)
	if valLen > 0 {
		self["EndPos"] = values[valLen - 1]["EndPos"]
	} else {
		self["EndPos"] = names[0]["EndPos"]
	}

	return self
}

fn Ast_ConstStatement(tok, names, values) {
	let self = {"Token": tok, "Names": names, "Values": values, "Type": ConstStatementNode}

	self["TokenLiteral"] = tok["literal"]

	let str = ""
	str += tok["literal"] + " ("

	for idx, name in names {
		str += name["TokenLiteral"]
		str += " = "
		let v = values[idx]
		if v != nil {
			str += v["String"]
		}
		str += ", "
	}

	str += ")"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]

	let valLen = len(values)
	if valLen > 0 {
		self["EndPos"] = values[valLen - 1]["EndPos"]
	} else {
		self["EndPos"] = names[0]["EndPos"]
	}

	return self
}

//Now because we support return multiple values, so here the "ReturnValue" is reserved only
//for campatibility with old code.
fn Ast_ReturnStatement(tok, rv, rvs) {
	let self = {"Token": tok, "ReturnValue": rv, "ReturnValues": rvs, "Type": ReturnStatementNode}

	self["TokenLiteral"] = tok["literal"]

	let str = tok["literal"] + " "

	let values = []
	for value in rvs {
		values.push(value["String"])
	}
	str += strings.join(values, ", ")

	str += ";"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]

	let rvsLen = len(rvs)
	if rvsLen > 0 {
		self["EndPos"] = rvs[rvsLen-1]["EndPos"]
	} else {
		let literalLen = len(tok["literal"])
		let tokPos = tok["position"]
		let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + literalLen, tokPos["Line"], tokPos["Column"] + literalLen)
		self["EndPos"] = endPos
	}

	return self
}

fn Ast_ExpressionStatement(tok, expr) {
	let self = {"Token": tok, "Expression": expr, "Type": ExpressionStatementNode}

	self["TokenLiteral"] = tok["literal"]

	self["String"] = ""
	if expr != nil {
		self["String"] = expr["String"]
	}

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = tok["position"]
	if expr != nil {
		self["EndPos"] = expr["EndPos"]
	}

	return self
}

fn Ast_IntegerLiteral(tok, value) {
	let self = {"Token": tok, "Value": value, "Type": IntegerLiteralNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = tok["literal"]

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let valLen = len(tok["literal"])
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + valLen, tokPos["Line"], tokPos["Column"] + valLen)
	self["EndPos"] = endPos

	return self
}

fn Ast_FloatLiteral(tok, value) {
	let self = {"Token": tok, "Value": value, "Type": FloatLiteralNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = tok["literal"]

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let valLen = len(tok["literal"])
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + valLen, tokPos["Line"], tokPos["Column"] + valLen)
	self["EndPos"] = endPos

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

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = right["EndPos"]

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

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = right["EndPos"]

	return self
}

fn Ast_Boolean(tok, value) {
	let self = {"Token": tok, "Value": value, "Type": BooleanNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = tok["literal"]

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let valLen = len(tok["literal"])
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + valLen, tokPos["Line"], tokPos["Column"] + valLen)
	self["EndPos"] = endPos

	return self
}

fn Ast_IfExpression(tok, condition, consequence, alternative) {
	let self = {"Token": tok, "Condition": condition, "Consequence": consequence, "Alternative": alternative}

	self["Type"] = IfExpressionNode
	self["TokenLiteral"] = tok["literal"]

	let str = "if "
	str += condition["String"]
	str += " {\n"
	str += "\t" + consequence["String"] + "}"
	if alternative != nil {
		str += "else {\n"
		str += "\t" + alternative["String"] + "}\n"
	}
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	if alternative != nil {
		self["EndPos"] = alternative["EndPos"]
	} else {
		self["EndPos"] = consequence["EndPos"]
	}

	return self
}

//rbraceToken: right brace token
fn Ast_BlockStatement(tok, statements, rbraceToken) {
	let self = {"Token": tok, "Statements": statements, "Type": BlockStatementNode}

	self["TokenLiteral"] = tok["literal"]

	let str = ""
	for stmt in statements {
		str += stmt["String"] + "\n"
	}
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = rbraceToken["position"]

	return self
}

fn Ast_FunctionStatement(tok, name, functionLiteral) {
	let self = {"Token": tok, "Name": name, "FunctionLiteral": functionLiteral, "Type": FunctionStatementNode}

	self["TokenLiteral"] = tok["literal"]


	let str = "fn "
	str += name["String"]

	let params = []
	for param in functionLiteral["Parameters"] {
		params.push(param["String"])
	}

	str += " ("
	str += strings.join(params, ", ")
	str += ") "
	str += "{\n"

	str += "\t" + functionLiteral["Body"]["String"]
	str += "}\n"

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = functionLiteral["Body"]["EndPos"]

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
	str += "{\n"
	str += "\t" + body["String"]
	str += "}\n"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = body["EndPos"]

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

	//start & end position
	let funLlen = len(function["String"])
	let tokPos = tok["position"]
	let startPos = newPosition(tokPos["Filename"], tokPos["Offset"] - funLlen, tokPos["Line"], tokPos["Column"] - funLlen)
	self["StartPos"] = startPos

	let aLen = len(arguments)
	if aLen > 0 {
		self["EndPos"] = arguments[aLen - 1]["EndPos"]
	} else {
		self["EndPos"] = function["EndPos"]
	}

	return self
}

fn Ast_MethodCallExpression(tok, obj, call) {
	let self = {"Token": tok, "Object": obj, "Call": call, "Type": MethodCallExpressionNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = obj["String"] + "." + call["String"]

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = call["EndPos"]

	return self
}

fn Ast_StringLiteral(tok, value) {
	let self = {"Token": tok, "Value": value, "Type": StringLiteralNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = tok["literal"]

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let valLen = len(value)
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + valLen, tokPos["Line"], tokPos["Column"] + valLen)
	self["EndPos"] = endPos

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

	//start & end position
	self["StartPos"] = tok["position"]

	let aLen = len(elements)
	if aLen > 0 {
		self["EndPos"] = elements[aLen-1]["EndPos"]
	} else {
		let tokPos = tok["position"]
		tokPos["Column"] = tokPos["Column"] + 1
		self["EndPos"] = tokPos
	}


	return self
}

//rParenToken: right paren token
fn Ast_TupleLiteral(tok, elements, rParenToken) {
	let self = {"Token": tok, "Elements": elements, "Type": TupleLiteralNode}

	self["TokenLiteral"] = tok["literal"]

	let str = ""
	let elems = []
	for elem in elements {
		elems.push(elem["String"])
	}
	str += "("
	str += strings.join(elems, ", ")
	str += ")"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = rParenToken["position"]

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

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = index["EndPos"]

	return self
}

//rbrackToken: right brace token
fn Ast_HashLiteral(tok, pairs, rbraceToken) {
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

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = rbraceToken["position"]

	return self
}

//for variable in iterable { block }
fn Ast_ForArrayExpression(tok, variable, iterable, block) {
	let self = {"Token": tok, "Var": variable, "Value": iterable, "Block": block, "Type": ForArrayExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "for "
	str += variable
	str += " in "
	str += iterable["String"]
	str += " {\n "
	str += "\t" + block["String"]
	str += "}\n"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = block["EndPos"]

	return self
}

//for key, value in hashObj { block }
fn Ast_ForMapExpression(tok, key, value, h, block) {
	let self = {"Token": tok, "Key": key, "Value": value, "Hash": h, "Block": block, "Type": ForMapExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "for "
	str += key + ", " + value
	str += " in "
	str += h["String"]
	str += " {\n"
	str += "\t" + block["String"]
	str += "}\n"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = block["EndPos"]

	return self
}

//for { block }
fn Ast_ForEverLoopExpression(tok, block) {
	let self = {"Token": tok, "Block": block, "Type": ForEverLoopExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "for "
	str += " {\n"
	str += "\t" + block["String"]
	str += "}\n"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = block["EndPos"]

	return self
}

//while (condition) { block }
fn Ast_WhileExpression(tok, condition, block) {
	let self = {"Token": tok, "Condition":  condition, "Block": block, "Type": WhileExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = "while "
	str += condition["String"]
	str += " {\n"
	str += "\t" + block["String"]
	str += "}\n"
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = block["EndPos"]

	return self
}

fn Ast_BreakExpression(tok) {
	let self = {"Token": tok, "TokenLiteral": tok["literal"], "String": tok["literal"], "Type": BreakExpressionNode}

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let literalLen = len(tok["literal"])
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + literalLen, tokPos["Line"], tokPos["Column"] + literalLen)
	self["EndPos"] = endPos

	return self
}


fn Ast_ContinueExpression(tok) {
	let self = {"Token": tok, "TokenLiteral": tok["literal"], "String": tok["literal"], "Type": ContinueExpressionNode}

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let literalLen = len(tok["literal"])
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + literalLen, tokPos["Line"], tokPos["Column"] + literalLen)
	self["EndPos"] = endPos

	return self
}

fn Ast_AssignExpression(tok, name, value) {
	let self = {"Token": tok, "Name":  name, "Value": value, "Type": AssignExpressionNode}

	self["TokenLiteral"] = tok["literal"]

	let str = name["String"]
	str += " = "
	str += value["String"]
	self["String"] = str

	//start & end position
	self["StartPos"] = tok["position"]
	self["EndPos"] = value["EndPos"]

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

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let operatorLen = len(operator)
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + operatorLen, tokPos["Line"], tokPos["Column"] + operatorLen)
	self["EndPos"] = endPos

	return self
}

fn Ast_NilLiteral(tok) {
	let self = {"Token": tok, "TokenLiteral": tok["literal"], "String": tok["literal"], "Type": NilExpressionNode}

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let literalLen = len(tok["literal"])
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + literalLen, tokPos["Line"], tokPos["Column"] + literalLen)
	self["EndPos"] = endPos

	return self
}

fn Ast_RegExLiteral(tok, value) {
	let self = {"Token": tok, "Value": value, "Type": RegExLiteralNode}

	self["TokenLiteral"] = tok["literal"]
	self["String"] = tok["literal"]

	//start & end position
	let tokPos = tok["position"]
	self["StartPos"] = tokPos

	let valLen = len(tok["literal"])
	let endPos = newPosition(tokPos["Filename"], tokPos["Offset"] + valLen, tokPos["Line"], tokPos["Column"] + valLen)
	self["EndPos"] = endPos

	return self
}

////////////////////////////////////////////////
//                 PARSER
////////////////////////////////////////////////
const (
	NONE  = 0,
	LOWEST,
	ASSIGN,      // =
	COND_OR,     // or
	COND_AND,    // and
	EQUALS,      // ==
	LESSGREATER, // >, <, >=, <=
	SHIFTS,      // >> or <<
	SUM,         // +
	PRODUCT,     // *
	PREFIX,      // -X or !X
	MATCHING,    // =~, !~
	CALL,        // myFunction(X)
	INDEX,       // array[index]
	INCREMENT,   // ++, --	
)


precedences = {
	TOKEN_ASSIGN:   ASSIGN,
	TOKEN_OR:       COND_OR,
	TOKEN_AND:      COND_AND,
	TOKEN_EQ:       EQUALS,
	TOKEN_NOT_EQ:   EQUALS,
	TOKEN_LT:       LESSGREATER,
	TOKEN_LE:       LESSGREATER,
	TOKEN_GT:       LESSGREATER,
	TOKEN_GE:       LESSGREATER,
	TOKEN_IN:       LESSGREATER,
	TOKEN_SHIFT_L:  SHIFTS,
	TOKEN_SHIFT_R:  SHIFTS,
	TOKEN_PLUS:     SUM,
	TOKEN_MINUS:    SUM,
	TOKEN_SLASH:    PRODUCT,
	TOKEN_ASTERISK: PRODUCT,
	TOKEN_POWER:    PRODUCT,
	TOKEN_MATCH:    MATCHING,
	TOKEN_UNMATCH:  MATCHING,
	TOKEN_LPAREN:   CALL,
	TOKEN_DOT:      CALL,
	TOKEN_LBRAKET:  INDEX,
	TOKEN_INCREMENT:INCREMENT,
	TOKEN_DECREMENT:INCREMENT,
}

let parser_Errors = []
let parser_curToken = newToken(TOKEN_ILLEGAL, "", nil)
let parser_peekToken = newToken(TOKEN_ILLEGAL, "", nil)

let isInLoop = false

fn debugToken(msg) {
	printf("%s, curToken.type=%s, curToken.literal=%s, curTok.Pos=%s, peekToken.type=%s, peekToken.literal=%s, peekToken.Pos=%s\n", 
			msg, tokenString(parser_curToken["type"]), parser_curToken["literal"], positionStr(parser_curToken["position"]),
			tokenString(parser_peekToken["type"]), parser_peekToken["literal"], positionStr(parser_peekToken["position"]))
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
		return parseLetStatement(true)
	} if parser_curToken["type"] == TOKEN_CONST {
		return parseConstStatement()
	} elif parser_curToken["type"] == TOKEN_RETURN {
		return parseReturnStatement()
	} elif parser_curToken["type"] == TOKEN_FUNCTION {
		return parseFunctionStatement()
	} elif parser_curToken["type"] == TOKEN_IDENT {
		//if the current token is an 'identifier' and next token is a ',', 
		//then we think it's a multiple assignment, but we treat it as a 'let' statement.
		//otherwise, we just call parseExpressionStatement.
		if peekTokenIs(TOKEN_COMMA) {
			return parseLetStatement(false)
		}
		return parseExpressionStatement()
	} else {
		return parseExpressionStatement()
	}
}

fn parseLetStatement(nextFlag) {
	let curToken = parser_curToken

	if !nextFlag {
		//construct a dummy 'let' token
		curToken = newToken(TOKEN_LET, "let", curToken["position"])
	}

	//parse left hand side of the assignment
	let names = []
	for {
		if nextFlag {
			parser_nextToken()
		}
		nextFlag = true

		if !curTokenIs(TOKEN_IDENT) && !curTokenIs(TOKEN_UNDERSCORE) {
			let msg = sprintf("Syntax Error%s: expected token to be identifier|underscore, got %s instead.", positionStr(parser_curToken["position"]), tokenString(parser_curToken["type"]))
			parser_Errors.push(msg)
			return nil
		}

		let name = Ast_Identifier(parser_curToken, parser_curToken["literal"])
		names.push(name)

		if !peekTokenIs(TOKEN_ASSIGN) && !curTokenIs(TOKEN_SEMICOLON) && !peekTokenIs(TOKEN_COMMA) {
			if peekTokenIs(TOKEN_SEMICOLON) {
				parser_nextToken()
			}

			return Ast_LetStatement(curToken, names, nil)
		}

		parser_nextToken()
		if curTokenIs(TOKEN_ASSIGN) || curTokenIs(TOKEN_SEMICOLON) {
			break
		}
		if !curTokenIs(TOKEN_COMMA) {
			let msg = sprintf("Syntax Error%s: expected token to be comma, got %s instead.", positionStr(parser_curToken["position"]), tokenString(parser_curToken["type"]))
			parser_Errors.push(msg)
			return nil
		}
	}

	if curTokenIs(TOKEN_SEMICOLON) { //let x;
		return Ast_LetStatement(curToken, names, nil)
	}

	//parse right hand side of the assignment
	let values = []
	let i = 0
	parser_nextToken()
	for {
		let tmp = parseExpressionStatement()
		values.push(tmp["Expression"])

		if !peekTokenIs(TOKEN_COMMA) {
			break
		}
		parser_nextToken()
		parser_nextToken()

		i += 1
	}

	return Ast_LetStatement(curToken, names, values)
}

fn parseConstStatement() {
	let curToken = parser_curToken

	let names = []
	let values = []

	if peekTokenIs(TOKEN_LPAREN) {
		parser_nextToken()

		let autoInt = 0 //autoIncrement

		let idPair = {}

		for {
			if peekTokenIs(TOKEN_RPAREN) {
				parser_nextToken()
				return Ast_ConstStatement(curToken, names, values)
			}

			if !expectPeek(TOKEN_IDENT) { // identifier is mandatory here
				return nil
			}

			let id = parseIdentifier()
			names.push(id)

			// peek next that can be only '=' or ',' or ')'
			if !peekTokenIs(TOKEN_ASSIGN) && !peekTokenIs(TOKEN_COMMA) && !peekTokenIs(TOKEN_RPAREN) {
				let msg = sprintf("Syntax Error%s: Token %s not allowed here.", positionStr(parser_peekToken["position"]), tokenString(parser_peekToken["type"]))
				parser_Errors.push(msg)
				return nil
			}

			// check for optional default value (optional only in `INT` case)
			let value = nil
			if peekTokenIs(TOKEN_ASSIGN) {
				parser_nextToken()
				parser_nextToken()

				let expStmt = parseExpressionStatement()
				value = expStmt["Expression"]
			}

			if value != nil {
				if value["Type"] == IntegerLiteralNode {
					autoInt = value["Value"] + 1
				}
			} else {
				//create a new INT token with 'autoInt' as it's value
				let tok = newToken(TOKEN_INT, str(autoInt), parser_curToken["position"])
				value = Ast_IntegerLiteral(tok, autoInt)
				autoInt++
			}

			let str_id = id["Value"]
			if idPair.exists(str_id) { //is identifier redeclared?
				let msg = sprintf("Syntax Error%s: Identifier %s redeclared.", positionStr(parser_curToken["position"]), str_id)
				parser_Errors.push(msg)
				return nil
			} else {
				idPair[str_id] = value
				values.push(value)
			}

			if !peekTokenIs(TOKEN_COMMA) {
				parser_nextToken()
				break
			}
			parser_nextToken()
		}

	} else {
		if !expectPeek(TOKEN_IDENT) {
			return nil
		}

		let name = parseIdentifier()
		names.push(name)

		if !expectPeek(TOKEN_ASSIGN) {
			return nil
		}

		parser_nextToken()
		let value = parseExpression(LOWEST)
		values.push(value)
	}

	if peekTokenIs(TOKEN_SEMICOLON) {
		parser_nextToken()
	}

	return Ast_ConstStatement(curToken, names, values)
}

fn parseReturnStatement() {
	let curToken = parser_curToken

	if peekTokenIs(TOKEN_SEMICOLON) { // e.g. { return; }
		parser_nextToken()
		return Ast_ReturnStatement(curToken, nil, [])
	}

	if peekTokenIs(TOKEN_RBRACE) { //e.g. { return }
		return Ast_ReturnStatement(curToken, nil, [])
	}

	let retValues = []
	parser_nextToken()
	for {
		let expr = parseExpressionStatement()
		retValues.push(expr["Expression"])

		if !peekTokenIs(TOKEN_COMMA) {
			break
		}
		parser_nextToken()
		parser_nextToken()
	}

	if len(retValues) > 0 {
		return Ast_ReturnStatement(curToken, retValues[0], retValues)
	} else {
		return Ast_ReturnStatement(curToken, nil, [])
	}
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
	} elif parser_curToken["type"] == TOKEN_FLOAT {
		leftExp = parseFloatLiteral()
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
	}  elif parser_curToken["type"] == TOKEN_FOR {
		leftExp = parseForExpression()
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
	} elif parser_curToken["type"] == TOKEN_NIL {
		leftExp = parseNilExpression()
	} elif parser_curToken["type"] == TOKEN_REGEXP {
		leftExp = parseRegExLiteralExpression()
	} else {
		noPrefixParseFnError(parser_curToken["type"])
		return nil
	}

	while !peekTokenIs(SEMICOLON) && precedence < peekPrecedence() {
		//infix expression
		if parser_peekToken["type"] == TOKEN_PLUS || parser_peekToken["type"] == TOKEN_MINUS ||
		   parser_peekToken["type"] == TOKEN_SLASH || parser_peekToken["type"] == TOKEN_ASTERISK || parser_peekToken["type"] == TOKEN_POWER ||
		   parser_peekToken["type"] == TOKEN_EQ || parser_peekToken["type"] == TOKEN_NOT_EQ ||
		   parser_peekToken["type"] == TOKEN_LT || parser_peekToken["type"] == TOKEN_GT ||
		   parser_peekToken["type"] == TOKEN_LE || parser_peekToken["type"] == TOKEN_GE ||
		   parser_peekToken["type"] == TOKEN_AND || parser_peekToken["type"] == TOKEN_OR ||
		   parser_peekToken["type"] == TOKEN_SHIFT_L || parser_peekToken["type"] == TOKEN_SHIFT_R ||
		   parser_peekToken["type"] == TOKEN_MATCH || parser_peekToken["type"] == TOKEN_UNMATCH || 
		   parser_peekToken["type"] == TOKEN_IN {
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

fn parseFloatLiteral() {
	let literal = parser_curToken["literal"]
	
	return Ast_FloatLiteral(parser_curToken, literal.parseFloat())
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

	// if the token is '**', we process it specially. e.g. 3 ** 2 ** 3 = 3 ** (2 ** 3)
	// i.e. Exponent operator '**'' has right-to-left associativity
	if curTokenIs(TOKEN_POWER) {
		precedence = precedence - 1
	}

	parser_nextToken()
	let right = parseExpression(precedence)

	return Ast_InfixExpression(curToken, left, curToken["literal"], right)
}

fn parseBoolean() {
	let val = curTokenIs(TOKEN_TRUE)
	return Ast_Boolean(parser_curToken, val)
}

fn parseGroupedExpression() {
	let curToken = parser_curToken
	parser_nextToken()

	if parser_curToken["type"] == TOKEN_LPAREN && curTokenIs(TOKEN_RPAREN) {
		//empty tuple, e.g. 'x = ()'
		return Ast_TupleLiteral(curToken, [], parser_curToken)
	}

	let exp = parseExpression(LOWEST)

	if peekTokenIs(TOKEN_COMMA) {
		parser_nextToken()
		return parseTupleExpression(curToken, exp)
	}

	if !expectPeek(TOKEN_RPAREN) {
		return nil
	}

	return exp
}

fn parseTupleExpression(tok, expr) {
	let elements = [expr]

	for {
		if parser_curToken["type"] == TOKEN_RPAREN {
			return Ast_TupleLiteral(tok, elements, parser_curToken)
		} elif parser_curToken["type"] == TOKEN_COMMA {
			parser_nextToken()
			//For a 1-tuple: "(1,)", the trailing comma is necessary to distinguish it
			//from the parenthesized expression (1).
			if curTokenIs(TOKEN_RPAREN) { //e.g.  let x = (1,)
				return Ast_TupleLiteral(tok, elements, parser_curToken)
			}
			elements.push(parseExpression(LOWEST))
			parser_nextToken()
		} else {
			let msg = sprintf("Syntax Error%s: expected token to be ',' or ')', got %s instead", positionStr(parser_curToken["position"]), tokenString(parser_curToken["type"]))
			parser_Errors.push(msg)
			return nil
		}
	}
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

	return Ast_BlockStatement(curToken, statements, parser_curToken)
}

fn parseFunctionStatement() {
	let curToken = parser_curToken

	/* why below 'if'? please see below code:
		1. fn add(x, y) { x + y }
		2. fn(x, y) { x + y }(2, 3)
	for the second one, we have no identifier, so we need
	not advance to the next token.
	*/
	if !peekTokenIs(TOKEN_LPAREN) {
		parser_nextToken()
	}

	let name = Ast_Identifier(parser_curToken, parser_curToken["literal"])

	let functionLiteral = parseFunctionLiteral()

	if peekTokenIs(TOKEN_SEMICOLON) {
		parser_nextToken
	}

	return Ast_FunctionStatement(curToken, name, functionLiteral)
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

	return Ast_HashLiteral(curToken, pairs, parser_curToken)
}

fn parseForExpression() {
	let curToken = parser_curToken

	if peekTokenIs(TOKEN_LBRACE) {
		return parseForEverLoopExpression(curToken)
	}

	if !expectPeek(TOKEN_IDENT) {
		return nil
	}

	let variable = parser_curToken["literal"]

	if peekTokenIs(TOKEN_COMMA) {
		return parseForEachMap(curToken, variable)
	}

	return parseForEachArray(curToken, variable)
}

fn parseForEverLoopExpression(tok) {
	isInLoop = true

	expectPeek(TOKEN_LBRACE)
	let block = parseBlockStatement()

	isInLoop = false

	return Ast_ForEverLoopExpression(tok, block) 
}

//for key, value in hashObj { block }
fn parseForEachMap(tok, variable) {
	isInLoop = true

	let key = variable

	if !expectPeek(TOKEN_COMMA) {
		return nil
	}

	if !expectPeek(TOKEN_IDENT) {
		return nil
	}

	let value = parser_curToken["literal"]

	if !expectPeek(TOKEN_IN) {
		return nil
	}

	parser_nextToken()
	let h = parseExpression(LOWEST)

	let block = nil
	if peekTokenIs(TOKEN_LBRACE) {
		parser_nextToken()
		block = parseBlockStatement()
	} else {
		let msg = sprintf("Syntax Error%s: for loop must be followed by a '{'.", positionStr(parser_curToken["position"]))
		parser_Errors.push(msg)
		return nil
	}

	isInLoop = false

	return Ast_ForMapExpression(tok, key, value, h, block) 
}

//for item in array { block }
fn parseForEachArray(tok, variable) {
	isInLoop = true

	if !expectPeek(TOKEN_IN) {
		return nil
	}
	parser_nextToken()

	let iterable = parseExpression(LOWEST)

	let block = nil
	if !expectPeek(TOKEN_LBRACE) {
		return nil
	}
	block = parseBlockStatement()

	isInLoop = false

	return Ast_ForArrayExpression(tok, variable, iterable, block)
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
	if !expectPeek(TOKEN_LBRACE) {
		return nil
	}
	block = parseBlockStatement()

	isInLoop = false

	return Ast_WhileExpression(curToken, condition, block)
}

fn parseBreakWithoutLoopContext() {
	let msg = sprintf("Syntax Error%s: 'break' outside of loop context", positionStr(parser_curToken["position"]))
	parser_Errors.push(msg)

	return parseBreakExpression()
}

fn parseBreakExpression() {
	return Ast_BreakExpression(parser_curToken)
}

fn parseContinueWithoutLoopContext() {
	let msg = sprintf("Syntax Error%s: 'continue' outside of loop context", positionStr(parser_curToken["position"]))
	parser_Errors.push(msg)

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

fn parseNilExpression() {
	return Ast_NilLiteral(parser_curToken)
}

fn parseRegExLiteralExpression() {
	return Ast_RegExLiteral(parser_curToken, parser_curToken["literal"])
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
	let msg = sprintf("Syntax Error%s: expected next token to be %s, got %s instead", 
		positionStr(parser_peekToken["position"]), tokenString(t), tokenString(parser_peekToken["type"]))
	parser_Errors.push(msg)
}

fn noPrefixParseFnError(t) {
	let msg = sprintf("Syntax Error%s: no prefix parse function for %s found", positionStr(parser_curToken["position"]), tokenString(t))
	parser_Errors.push(msg)
}



////////////////////////////////////////////////
//                 SCOPE
////////////////////////////////////////////////
fn NewScope() {
	return {"store": {}, "readonly": {}, "outer": nil}
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
	let ro = scope["readonly"]

	if store.exists(name) && ro.exists(name) {
		printf("Const variable '%s' cannot be modified.\n", name)
		os.exit(3)
	}

	store[name] = val
	return val
}

fn SetIntoConstScope(scope, name, val) {
	let store = scope["store"]
	let ro = scope["readonly"]

	store[name] = val
	
	ro[name] = true //mark it as readonly

	return val

}

fn ResetScope(scope, name, val) {
	let store = scope["store"]
	let ro = scope["readonly"]

	let obj = store[name]
	if obj != nil {
		if ro.exists(name) {
			printf("Const variable '%s' cannot be modified.\n", name)
			os.exit(3)
		}
		store[name] = val
	}

	// if we not found in the current, we search the parent scope
	if obj == nil {
		let outer = scope["outer"]
		if outer != nil {
			obj = ResetScope(outer, name, val)
		}
	}

	// still could not found in the parent scope, then we set it in the current scope
	if obj == nil {
		store[name] = val
	}

	return val
}

//global scope used for some builtin object
let globalScopes = {}

fn GetGlobalObj(name) {
	if globalScopes.exists(name) {
		return globalScopes[name]
	}
	return NIL
}

fn SetGlobalObj(name, obj) {
	globalScopes[name] = obj
}


////////////////////////////////////////////////
//                 OBJECT
////////////////////////////////////////////////
/* Object Type */
const INTEGER_OBJ      = "INTEGER"
const FLOAT_OBJ        = "FLOAT"
const BOOLEAN_OBJ      = "BOOLEAN"
const NIL_OBJ          = "NIL"
const RETURN_VALUE_OBJ = "RETURN_VALUE"
const ERROR_OBJ        = "ERROR"
const FUNCTION_OBJ     = "FUNCTION"
const STRING_OBJ       = "STRING"
const BUILTIN_OBJ      = "BUILTIN"
const ARRAY_OBJ        = "ARRAY"
const TUPLE_OBJ        = "TUPLE"
const HASH_OBJ         = "HASH"
const BREAK_OBJ        = "BREAK"
const CONTINUE_OBJ     = "CONTINUE"
const FILE_OBJ         = "FILE"
const REGEX_OBJ        = "REGEX"

/* Every Object must have below properties:
    "Type":       Object Type
    "Inspect":    Mainly for debug
    "CallMethod": Object method which has the following prototype
         (scope, line, methodName, methodArgs) -> Object"

   And may be other properties, e.g.
    "Iterable": means the object could be used in a 'for-in' expression. e.g. STRING_OBJ/ARRAY_OBJ/HASH_OBJ
    "Writable": means the object is writable, e.g. FILE_OBJ
*/

//lhsV : left-hand side Value
fn objEqual(lhsV, rhsV) {
	if lhsV == nil && rhsV == nil {
		return true
	}
	if (lhsV == nil && rhsV != nil) || (lhsV != nil && rhsV == nil) {
		return false
	}

	if lhsV["Type"] != rhsV["Type"] {
		return false
	}

	if lhsV["Type"] == NIL_OBJ {
		if rhsV["Type"] == NIL_OBJ {
			return true
		}
	}

	return lhsV["Inspect"] == rhsV["Inspect"]
}


fn IntegerObj(value) {
	let self = {"Type": INTEGER_OBJ, "Inspect": sprintf("%d", value), "Value": value, "IsNumber": true}

	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
		if methodName == "next" {
			return Integer_Next(value)

		} elif methodName == "prev" {
			return Integer_Prev(value)
		} else {
			return undefinedMethodError(line, methodName, INTEGER_OBJ)
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

fn FloatObj(value) {
	let self = {"Type": FLOAT_OBJ, "Inspect": sprintf("%g", value), "Value": value, "IsNumber": true}

	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
		if methodName == "ceil" {
			return Float_Ceil(value)

		} elif methodName == "floor" {
			return Float_Floor(value)
		} else {
			return undefinedMethodError(line, methodName, FLOAT_OBJ)
		}
	}

	return self
}

fn Float_Ceil(value) {
	return FloatObj(value.ceil())
}

fn Float_Floor(value) {
	return FloatObj(value.floor())
}

fn BooleanObj(value) {
	let self = {"Type": BOOLEAN_OBJ, "Inspect": sprintf("%t", value), "Value": value}

	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
		if methodName == "toYesNo" {
			return Boolean_ToYesNo(value)

		} elif methodName == "toTrueFalse" {
			return Boolean_ToTrueFalse(value)
		} else {
			return undefinedMethodError(line, methodName, BOOLEAN_OBJ)
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

//Now because we support return multiple values, so here the "Value" is reserved only
//for campatibility with old code.
fn ReturnObj(value, values) {
	let self = {"Type": RETURN_VALUE_OBJ, "Value": value, "Values": values, "CallMethod": nil}

	let retValues = []
	for v in values {
		retValues.push(v["Inspect"])
	}

	self["Inspect"] = strings.join(retValues, ", ")

	return self
}

//Generic Error Object
fn ErrorObj(line, message) {
	let msg = sprintf("Runtime Error at line %d: %s", line, message)
	return {"Type": ERROR_OBJ, "Inspect": msg, "CallMethod": nil}
}

//Specific error object for 'undefined method'
fn undefinedMethodError(line, methodName, obj) {
	let msg = sprintf("undefined method '%s' for object %s", methodName, obj)
	return ErrorObj(line, msg)
}

//Specific error object for 'wrong number of arguemnts' of function
fn argumentError(line, got, want) {
	let msg = sprintf("wrong number of arguments. got=%d, want=%s", got, want)
	return ErrorObj(line, msg)
}

//Specific error object for 'prefix' operator
fn prefixOperatorError(line, operator, rightType) {
	let msg = sprintf("unsupported operator for prefix expression:'%s' and type: %s", operator, rightType)
	return ErrorObj(line, msg)
}

//Specific error object for 'infix' operator
fn infixOperatorError(line, leftType, operator, rightType) {
	let msg = sprintf("unsupported operator for infix expression: %s '%s' %s", leftType, operator, rightType)
	return ErrorObj(line, msg)
}

//Specific error object for 'postfix' operator
fn postfixOperatorError(line, operator, leftType) {
	let msg = sprintf("unsupported operator for postfix expression:'%s' and type: %s", operator, leftType)
	return ErrorObj(line, msg)
}

//Specifif error object for 'unknown identifier'
fn unknownIdentifierError(line, name) {
	let msg = sprintf("unknown identifier: '%s' is not defined", name)
	return ErrorObj(line, msg)
}

//Specific error object for 'index out-of-range'
fn indexError(line, idx) {
	let msg = sprintf("index error: '%d' out of range", idx)
	return ErrorObj(line, msg)
}

//Specific error object for 'index type'
fn indexTypeError(line, leftType) {
	let msg = sprintf("index error: type %s is not indexable", leftType)
	return ErrorObj(line, msg)
}

//Specific error object for 'not iterable'
fn notIterableError(line) {
	return ErrorObj(line, "foreach's operating type must be iterable")
}

//Specific error for 'parameter type'
fn parameterTypeError(line, ordinal, name, expectType, gotType) {
	let msg = sprintf("%s argument for '%s' should be type '%s'. got '%s'", ordinal, name, expectType, gotType)
	return ErrorObj(line, msg)
}

//Specific error for 'not callable'
fn callableError(line, funcName) {
	let msg = sprintf("'%s' is not callable", funcName)
	return  ErrorObj(line, msg)
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
	let self = {"Type": STRING_OBJ, "Inspect": value, "Value": value, "Iterable": true}

	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
		if methodName == "upper" {
			return String_Upper(value)
		} elif methodName == "lower" {
			return String_Lower(value)
		} elif methodName == "trim" {
			return String_Trim(value)
		} elif methodName == "set" {
			return String_Set(self, scope, methodArgs)
		} elif methodName == "replace" {
			return String_Replace(self, scope, methodArgs)
		} elif methodName == "find" {
			return String_Find(self, scope, methodArgs)
		} elif methodName == "index" {
			return String_Find(self, scope, methodArgs)
		} elif methodName == "substr" {
			return String_Substr(self, scope, methodArgs)
		} elif methodName == "contains" {
			return String_Contains(self, scope, methodArgs)
		} else {
			return undefinedMethodError(line, methodName, STRING_OBJ)
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

fn String_Set(self, scope, methodArgs) {
	let idx = methodArgs[0]
	let replaceStr = methodArgs[1]
	let value = self["Value"]
	value.set(idx["Value"], replaceStr["Value"])
	self["Inspect"] = value
}

fn String_Replace(self, scope, methodArgs) {
	let srcStr = methodArgs[0]
	let replaceStr = methodArgs[1]

	let value = self["Value"] 
	let ret = value.replace(srcStr["Value"], replaceStr["Value"])
	return StringObj(ret)
}

fn String_Find(self, scope, methodArgs) {
	let subStr = methodArgs[0]

	let value = self["Value"] 
	let idx = value.find(subStr["Value"])
	return IntegerObj(idx)
}

fn String_Substr(self, scope, methodArgs) {
	let pos = methodArgs[0]
	let length = methodArgs[1]

	let value = self["Value"]
	let ret = value.substr(pos["Value"], length["Value"])
	return StringObj(ret)
}

fn String_Contains(self, scope, methodArgs) {
	let arg0 = methodArgs[0]

	let value = self["Value"]
	if value.contains(arg0["Inspect"]) {
		return TRUE
	}

	return FALSE
}


fn BuiltinObj(builtinFn) {
	return {"Type": BUILTIN_OBJ, "Inspect": "<builtin>", "BuiltinFn": builtinFn, "CallMethod": nil}
}

fn ArrayObj(elements) {
	let self = {"Type": ARRAY_OBJ, "Elements": elements, "Iterable": true}

	Array_Inspect(self)

	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
		if methodName == "push" {
			return Array_Push(self, scope, methodArgs)
		} elif methodName == "pop" {
			return Array_Pop(self, scope)
		} elif methodName == "set" {
			return Array_Set(self, scope, methodArgs)
		} elif methodName == "get" {
			return Array_Get(self, scope, methodArgs)
		} elif methodName == "include" {
			return Array_Include(self, scope, methodArgs)
		} else {
			return undefinedMethodError(line, methodName, ARRAY_OBJ)
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

fn Array_Include(self, scope, args) {
	let elements = self["Elements"]

	let obj = args[0]

	for item in elements {
		if objEqual(obj, item) {
			return TRUE
		}
	}
	return FALSE
}

fn Array_Get(self, scope, args) {
	let elements = self["Elements"]

	let idxObj = args[0]
	let value = idxObj["Value"]
	if value < 0 || value >= len(elements) {
		return NIL
	}
	return elements[value]
}

fn Array_Set(self, scope, args) {
	let elements = self["Elements"]

	idxObj = args[0]

	let value = idxObj["Value"]

	if value < 0 || value >= len(elements) {
		let oldLen = len(elements)
		let i = 0;
		for (i = oldLen; i <= value; i++) {
			elements.push(NIL)
		}
	}

	elements[value] = args[1]

	//recalculate 'Inspect'
	Array_Inspect(self)

	return NIL
}

/* 'IsMulti' flag is used in function return values.
   if a function returns multiple values, they will be wrapped into a tuple,
   the flag will be set to true
*/
fn TupleObj(elements, isMulti = false) {
	let self = {"Type": TUPLE_OBJ, "Elements": elements, "Iterable": true, "IsMulti": isMulti}

	let items = []
	for elem in elements {
		items.push(elem["Inspect"])
	}

	let str = ""
	str += "("
	str += strings.join(items, ", ")
	if len(elements) == 1 {
		str += ","
	}
	str += ")"
	self["Inspect"] = str


	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
		if methodName == "include" {
			return Tuple_Include(self, scope, methodArgs)
		} else {
			return undefinedMethodError(line, methodName, TUPLE_OBJ)
		}
	}

	return self
}

fn Tuple_Inspect(self) {

}

fn Tuple_Include(self, scope, args) {
	let elements = self["Elements"]

	let obj = args[0]

	for item in elements {
		if objEqual(obj, item) {
			return TRUE
		}
	}
	return FALSE
}

fn HashPair(key, value) {
	return {"Key": key, "Value": value}
}

fn HashObj(pairs) {
	//Pairs map[HashKey]HashPair
	//Here I use obj["Inspect"] as the hash key.

	let self = {"Type": HASH_OBJ, "Pairs": pairs, "Iterable": true}

	Hash_Inspect(self)

	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
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
			return undefinedMethodError(line, methodName, HASH_OBJ)
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

fn FileObj(filehandle, filename) {
	let self = {"Type": FILE_OBJ, "File": filehandle, "Inspect": "<file object: " + filename + ">", "Writable": true}

	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
		if methodName == "close" {
			return File_Close(self, scope, methodArgs)
		} elif methodName == "read" {
			return File_Read(self, scope, methodArgs)
		} elif methodName == "readLine" {
			return File_ReadLine(self, scope, methodArgs)
		} elif methodName == "write" {
			return File_Write(self, scope, methodArgs)
		} elif methodName == "writeLine" {
			return File_WriteLine(self, scope, methodArgs)
		} elif methodName == "writeString" {
			return File_WriteString(self, scope, methodArgs)
		} elif methodName == "writeString" {
			return File_WriteString(self, scope, methodArgs)
		} else {
			return undefinedMethodError(line, methodName, FILE_OBJ)
		}
	}

	return self
}

//Because it's just a simple demostration. 
//so below code has no error check.
fn File_Close(self, scope, methodArgs) {
	let fh = self["File"]
	fh.close()
	return NIL
}

fn File_Read(self, scope, methodArgs) {
	let fh = self["File"]

	let size = 1024 //default buffer
	if len(methodArgs) >0 {
		size = methodArgs[0]["Value"]
	}
	let s = fh.read(size)
	return StringObj(s)
}

fn File_ReadLine(self, scope, methodArgs) {
	let fh = self["File"]
	let s = fh.readLine()

	if s == nil {
		return NIL
	}
	return StringObj(s)
}

fn File_Write(self, scope, methodArgs) {
	let fh = self["File"]
	let s = methodArgs[0]

	let n = fh.write(s["Value"])
	return IntegerObj(n)
}

fn File_WriteLine(self, scope, methodArgs) {
	let fh = self["File"]
	let s = methodArgs[0]
	let n = fh.writeLine(s["Value"])
	return IntegerObj(n)
}

fn File_WriteString(self, scope, methodArgs) {
	let fh = self["File"]
	let s = methodArgs[0]

	let n = fh.writeString(s["Value"])
	return IntegerObj(n)
}

fn RegExObj(regex, value) {
	let self = {"Type": REGEX_OBJ, "Inspect": value, "RegEx":regex, "Value": value}

	self["CallMethod"] = fn (scope, line, methodName, methodArgs) {
		if methodName == "match" {
			return RegEx_Match(self, methodArgs)
		} elif methodName == "replace" {
			return RegEx_Replace(self, methodArgs)
		} elif methodName == "split" {
			return RegEx_Split(self, methodArgs)
		} else {
			return undefinedMethodError(line, methodName, BOOLEAN_OBJ)
		}
	}

	return self
}

fn RegEx_Match(self, methodArgs) {
	let regex = self["RegEx"]

	let matched = regex.match(methodArgs[0])
	if matched {
		return TRUE
	}

	return FALSE
}

fn RegEx_Replace(self, methodArgs) {
	let regex = self["RegEx"]

	let source = methodArgs[0]
	let repl = methodArgs[1]

	let result = regex.replace(source["Value"], repl["Value"])
	return StringObj(result)
}

fn RegEx_Split(self, methodArgs) {
	let regex = self["RegEx"]

	let source = methodArgs[0]

	strArr = regex.Split(source, -1)
	let arr = []
	for v in strArr {
		arr.push(v)
	}

	return ArrayObj(arr) 

}



////////////////////////////////////////////////
//                 BUILTINS
////////////////////////////////////////////////
/* builtins = map[string]BuiltinObj */

fn openBuiltin(line, args) {
	if len(args) != 1 && len(args) != 2 {
		return argumentError(line, len(args), "1|2")
	}

	let filename = args[0]["Value"]
	let mode = "r" // default is 'read'
	if len(args) == 2 {
		mode = args[1]["Value"]
	}

	let f = newFile(filename, mode)
	return FileObj(f, filename)
}

fn lenBuiltin(line, args) {
	if len(args) != 1 {
		return argumentError(line, len(args), "1")
	}

	let arg0 = args[0]
	let arg0Type = arg0["Type"]
	if arg0Type == ARRAY_OBJ || arg0Type == TUPLE_OBJ {
		return IntegerObj(len(arg0["Elements"]))
	} elif arg0Type == STRING_OBJ {
			return IntegerObj(len(arg0["Value"]))
	}

	return parameterTypeError(line, "first", "len", "ARRAY|TUPLE|STRING", arg0Type)
}

fn printBuiltin(line, args) {
	for arg in args {
		print(arg["Inspect"])
	}

	return NIL
}

fn printlnBuiltin(line, args) {
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
	"open": BuiltinObj(openBuiltin),
}

//Predefined global object
globalScopes["stdin"]  = FileObj(stdin, "<stdin>")
globalScopes["stdout"] = FileObj(stdout, "<stdout>")
globalScopes["stderr"] = FileObj(stderr, "<stderr>")
globalScopes["endl"] = StringObj("\n")

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
	} elif nodeType == FloatLiteralNode {
		return FloatObj(node["Value"])
	} elif nodeType == BooleanNode {
		return nativeBoolToBooleanObject(node["Value"])
	} elif nodeType == PrefixExpressionNode {
		let p_right = Eval(node["Right"], scope)
		if isError(p_right) {
			return p_right
		}
		return evalPrefixExpression(node, p_right)
	} elif nodeType == InfixExpressionNode {
		let i_left = Eval(node["Left"], scope)
		if isError(i_left) {
			return i_left
		}

		let i_right = Eval(node["Right"], scope)
		if isError(i_right) {
			return i_right
		}

		return evalInfixExpression(node, i_left, i_right, scope)
	} elif nodeType == BlockStatementNode {
		return evalBlockStatement(node, scope)
	} elif nodeType == IfExpressionNode {
		return evalIfExpression(node, scope)
	} elif nodeType == ReturnStatementNode {
		let retValues = []
		for value in node["ReturnValues"] {
			let val = Eval(value, scope)
			if isError(val) {
				return val
			}
			retValues.push(val)
		}

		let retValue = NIL
		if len(retValues) > 0 {
			retValue = retValues[0] // for old campatibility
		}

		return ReturnObj(retValue, retValues)
	} elif nodeType == LetStatementNode {
		return evalLetStatement(node, scope)
//		let l_val = Eval(node["Value"], scope)
//		if isError(l_val) {
//			return l_val
//		}
//
//		let l_name = node["Name"]
//		let name_value = l_name["Value"]
//		SetIntoScope(scope, name_value, l_val)
	} elif nodeType == ConstStatementNode {
		return evalConstStatement(node, scope)
	} elif nodeType == IdentifierNode {
		return evalIdentifier(node, scope)
	} elif nodeType == FunctionLiteralNode {
		return evalFunctionLiteral(node, scope)
	} elif nodeType == FunctionStatementNode {
		return evalFunctionStatement(node, scope)
	} elif nodeType == CallExpressionNode {
		let call_fn = Eval(node["Function"], scope)
		if isError(call_fn) {
			return call_fn
		}

		let call_args = evalExpressions(node["Arguments"], scope)
		if len(call_args) == 1 && isError(call_args[0]) {
			return call_args[0]
		}

		return applyFunction(node, call_fn, call_args)
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
	} elif nodeType == TupleLiteralNode {
		let tup_elements = evalExpressions(node["Elements"], scope)
		if len(tup_elements) == 1 && isError(tup_elements[0]) {
			return tup_elements[0]
		}
		return TupleObj(tup_elements)
	} elif nodeType == IndexExpressionNode {
		let idx_left = Eval(node["Left"], scope)
		if isError(idx_left) {
			return idx_left
		}
		let idx_index = Eval(node["Index"], scope)
		if isError(idx_index) {
			return idx_index
		}
		return evalIndexExpression(node, idx_left, idx_index)
	} elif nodeType == HashLiteralNode {
		return evalHashLiteral(node, scope)
	} elif nodeType == ForEverLoopExpressionNode {
		return evalForEverLoopExpression(node, scope)
	} elif nodeType == ForArrayExpressionNode {
		return evalForArrayExpression(node, scope)
	} elif nodeType == ForMapExpressionNode {
		return evalForMapExpression(node, scope)
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
		if isError(p_left) {
			return p_left
		}
		return evalPostfixExpression(node, p_left)
	} elif nodeType == NilExpressionNode {
		return NIL
	} elif nodeType == RegExLiteralNode {
		return evalRegExLiteral(node, scope)
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
				return result["Values"]
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

fn objectToNativeBoolean(o) {
	let obj = o

	if obj["Type"] == RETURN_VALUE_OBJ {
		obj = unwrapReturnValue(obj)
	}

	if obj["Type"] == INTEGER_OBJ {
		if obj["Value"] == 0 {
			return false
		}
		return true
	} elif obj["Type"] == FLOAT_OBJ {
		if obj["Value"] == 0.0 {
			return false
		}
		return true
	} elif obj["Type"] == BOOLEAN_OBJ {
		return obj["Value"]
	} elif obj["Type"] == ARRAY_OBJ {
		if len(obj["Elements"]) == 0 {
			return false
		}
		return true
	} elif obj["Type"] == TUPLE_OBJ {
		if len(obj["Elements"]) == 0 {
			return false
		}
		return true
	} elif obj["Type"] == HASH_OBJ {
		let h = obj["Pairs"]
		if h.len() == 0 {
			return false
		}
		return true
	}

	return true
}

fn evalPrefixExpression(node, right) {
	let operator = node["Operator"]
	if operator == "!" {
		return evalBangOperatorExpression(right)
	} elif operator =="-" {
		return evalMinusPrefixOperatorExpression(node, right)
	} elif operator == "++" {
		return evalIncrementPrefixOperatorExpression(node, right)
	} elif operator == "--" {
		return evalDecrementPrefixOperatorExpression(node, right)
	}

	return prefixOperatorError(node["StartPos"]["Line"], operator, right["Type"])
}

fn evalBangOperatorExpression(right) {
	return nativeBoolToBooleanObject(!isTruthy(right))
}

fn evalMinusPrefixOperatorExpression(node, right) {
	if right["Type"] == INTEGER_OBJ {
		let value = right["Value"]
		return IntegerObj(-value)
	} elif right["Type"] == FLOAT_OBJ {
		let value = right["Value"]
		return FloatObj(-value)
	}

	return prefixOperatorError(node["StartPos"]["Line"], "-", right["Type"])
}

fn evalIncrementPrefixOperatorExpression(node, right) {
	if right["Type"] == INTEGER_OBJ {
		let value = right["Value"]
		right["Value"] = value + 1
		right["Inspect"] = sprintf("%d", value + 1)

		return IntegerObj(value + 1)

	} elif right["Type"] == FLOAT_OBJ {
		let value = right["Value"]
		right["Value"] = value + 1
		right["Inspect"] = sprintf("%g", value + 1)

		return FloatObj(value + 1)
	}

	return prefixOperatorError(node["StartPos"]["Line"], "++", right["Type"])
}

fn evalDecrementPrefixOperatorExpression(node, right) {
	if right["Type"] == INTEGER_OBJ {
		let value = right["Value"]
		right["Value"] = value - 1
		right["Inspect"] = sprintf("%d", value - 1)

		return IntegerObj(value - 1)

	} elif right["Type"] == FLOAT_OBJ {
		let value = right["Value"]
		right["Value"] = value - 1
		right["Inspect"] = sprintf("%g", value - 1)

		return FloatObj(value - 1)
	}

	return prefixOperatorError(node["StartPos"]["Line"], "--", right["Type"])
}

fn evalInfixExpression(node, left, right, scope) {
	let operator = node["Operator"]

	//check if it is 'Writable'
	if left["Writable"] != nil && left["Writable"] {
		if operator == "<<" {
			if left["Type"] == FILE_OBJ {
				let strObj = StringObj(right["Inspect"])
				File_Write(left, scope, [strObj])
				//Here we return left, so we can chain multiple '<<'.
				// e.g.
				//     stdout << "hello " << "world!"
				return left
			}
		} elif operator == ">>" {
			// let a;
			// stdin >> a;

			//right should be an identifier
			if node["Right"]["Type"] != IdentifierNode {
				return infixOperatorError(node["StartPos"]["Line"], left["Type"], operator, right["Type"])
			}
			ret = File_ReadLine(left, scope, [])
			SetIntoScope(scope, node["Right"]["String"], ret)
			return ret
		}
	}

	if operator == "in" {
		if right["Type"] == ARRAY_OBJ {
			return Array_Include(right, scope, [left])
		} elif right["Type"] == TUPLE_OBJ {
			return Tuple_Include(right, scope, [left])
		} elif right["Type"] == HASH_OBJ {
			return Hash_Exists(right, scope, [left])
		} elif right["Type"] == STRING_OBJ {
			return String_Contains(right, scope, [left])
		}
		return FALSE
	}

	if operator == "and" || operator == "or" {
		if operator == "and" {
			let leftCond = objectToNativeBoolean(left)
			if  leftCond == false {
				return FALSE
			}
			let rightCond = objectToNativeBoolean(right)
			return nativeBoolToBooleanObject(leftCond && rightCond)
		} elif operator == "or" {
			let leftCond = objectToNativeBoolean(left)
			if leftCond == true {
				return TRUE
			}

			let rightCond = objectToNativeBoolean(right)
			return nativeBoolToBooleanObject(leftCond || rightCond)
		}
	} elif left["IsNumber"] && right["IsNumber"] {
		return evalNumberInfixExpression(node, left, right)
	} elif operator == "==" {
		if left["Type"] == BOOLEAN_OBJ && right["Type"] == BOOLEAN_OBJ {
			if l["Value"] == r["Value"] {
				return TRUE
			}
			return FALSE
		}

		if left["Type"] == NIL_OBJ && right["Type"] == NIL_OBJ { //(s == nil) should return true if s is nil
			return TRUE
		}

		return nativeBoolToBooleanObject(left == right)
	} elif operator == "!=" {
		if left["Type"] != right["Type"] {
			return TRUE
		}

		if left["Type"]== BOOLEAN_OBJ && right["Type"] == BOOLEAN_OBJ {
			if l["Value"] != r["Value"] {
				return TRUE
			}
			return FALSE
		}

		return nativeBoolToBooleanObject(left != right)
	} elif left["Type"] == STRING_OBJ && right["Type"] == STRING_OBJ {
		return evalStringInfixExpression(node, left, right)
	} elif left["Type"] != right["Type"] {
		return infixOperatorError(node["StartPos"]["Line"], left["Type"], operator, right["Type"])
	}

	return infixOperatorError(node["StartPos"]["Line"], left["Type"], operator, right["Type"])
}

fn evalNumberInfixExpression(node, left, right) {
	let operator = node["Operator"]
	let leftVal = left["Value"]
	let rightVal = right["Value"]

	let isInt = left["Type"] == INTEGER_OBJ && right["Type"] == INTEGER_OBJ

	if operator == "+" {
		if isInt {
			return IntegerObj(leftVal + rightVal)
		} else {
			return FloatObj(leftVal + rightVal)
		}
	} elif operator == "-" {
		if isInt {
			return IntegerObj(leftVal - rightVal)
		} else {
			return FloatObj(leftVal - rightVal)
		}
	} elif operator == "*" {
		if isInt {
			return IntegerObj(leftVal * rightVal)
		} else {
			return FloatObj(leftVal * rightVal)
		} 
	} elif operator == "/" {
		if isInt {
			return IntegerObj(leftVal / rightVal)
		} else {
			return FloatObj(leftVal / rightVal)
		}
	} elif operator == "**" {
		let v = math.pow(leftVal, rightVal)
		if isInt {
			return IntegerObj(int(v))
		} else {
			return FloatObj(v)
		}
	} elif operator == "<" {
		return nativeBoolToBooleanObject(leftVal < rightVal)
	} elif operator == "<=" {
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	} elif operator == ">" {
		return nativeBoolToBooleanObject(leftVal > rightVal)
	} elif operator == ">=" {
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	} elif operator == "==" {
		return nativeBoolToBooleanObject(leftVal == rightVal)
	} elif operator == "!=" {
		return nativeBoolToBooleanObject(leftVal != rightVal)
	} elif operator == ">>" {
		let aRes = leftVal >> rightVal
		return IntegerObj(aRes)
	} elif operator == "<<" {
		let aRes = leftVal << rightVal
		return IntegerObj(aRes)
	}

	return infixOperatorError(node["StartPos"]["Line"], left["Type"], operator, right["Type"])
}

fn evalIfExpression(ie, scope) {
	let condition = Eval(ie["Condition"], scope)
	if isError(condition) {
		return condition
	}

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

fn evalLetStatement(node, scope) {
	let values = []
	let valuesLen = 0
	for value in node["Values"] {
		let val = Eval(value, scope)
		if isError(val) {
			return val
		}

		if val["Type"] == TUPLE_OBJ {
			let tupleObj = val
			if tupleObj["IsMulti"] {
				valuesLen += len(tupleObj["Elements"])
				for tupleItem in tupleObj["Elements"] {
					values.push(tupleItem)
				}
			} else {
				valuesLen += 1
				values.push(tupleObj)
			}
		} else {
			valuesLen += 1
			values.push(val)
		}
	}

	for idx, item in node["Names"] {
		let tok = item["Token"]
		if idx >= valuesLen { //There are more Names than Values
			let tok = item["Token"]
			if tok["Type"] != TOKEN_UNDERSCORE {
				SetIntoScope(scope, item["String"], NIL)
			}
		} else {
			if tok["Type"] == TOKEN_UNDERSCORE {
				continue
			}
			val = values[idx]
			if val["Type"] != ERROR_OBJ {
				SetIntoScope(scope, item["String"], val)
			} else {
				return
			}
		}
	}
}

fn evalConstStatement(node, scope) {
	for idx, name in node["Names"] {
		let v = node["Values"]
		val = Eval(v[idx], scope)
		if isError(val) {
			return val
		}
		SetIntoConstScope(scope, name["Value"], val)
	}

	return NIL
}

fn evalIdentifier(node, scope) {
	//Get from global scope first
	let obj = GetGlobalObj(node["String"])
	if obj != NIL {
		return obj
	}

	let nodeVal = node["Value"]
	let val = GetFromScope(scope, nodeVal)
	if val != nil {
		return val
	}

	let built_in = builtins[nodeVal]
	if built_in != nil {
		return built_in
	}

	return unknownIdentifierError(node["StartPos"]["Line"], nodeVal)
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

fn evalStringInfixExpression(node, left, right) {
	let operator = node["Operator"]
	let lv = left["Value"]
	let rv = right["Value"]

	if operator == "=~" {
		regObj = regexp.compile(rv)
		let matched = regObj.match(lv)
		if matched {
			return TRUE
		}
		return FALSE
	} elif operator == "!~" {
		regObj = regexp.compile(rv)
		let matched = regObj.match(lv)
		if matched {
			return FALSE
		}
		return TRUE
	} elif operator == "==" {
		return nativeBoolToBooleanObject(lv == rv)
	} elif operator == "!=" {
		return nativeBoolToBooleanObject(lv != rv)
	} elif operator == "+" {
		return StringObj(lv + rv)
	} elif operator == "<" {
		return nativeBoolToBooleanObject(lv < rv)
	} elif operator == "<=" {
		return nativeBoolToBooleanObject(lv <= rv)
	} elif operator == ">" {
		return nativeBoolToBooleanObject(lv > rv)
	} elif operator == ">=" {
		return nativeBoolToBooleanObject(lv >= rv)
	} else {
		return infixOperatorError(node["StartPos"]["Line"], left["Type"], operator, right["Type"])
	}

	return StringObj(lv + rv)
}

fn evalIndexExpression(node, left, index) {
	if (left["Type"] == ARRAY_OBJ || left["Type"] == TUPLE_OBJ) && index["Type"] == INTEGER_OBJ {
		return evalArrayIndexExpression(node, left, index)
	} elif left["Type"] == HASH_OBJ {
		return evalHashIndexExpression(node, left, index)
	} elif left["Type"] == STRING_OBJ && index["Type"] == INTEGER_OBJ {
		return evalStringIndexExpression(node, left, index)
	}

	return indexTypeError(node["StartPos"]["Line"], left["Type"])
}

fn evalStringIndexExpression(node, left, index) {
	let l = left["Value"]
	let length = len(l)
	let idx = index["Value"]

	if idx > length - 1 {
		return indexError(node["StartPos"]["Line"], idx)
	}

	if idx < 0 {
		idx = length + idx
		if idx > length - 1 || idx < 0 {
			return indexError(node["StartPos"]["Line"], idx)
		}
	}

	return StringObj(l[idx])
}

fn evalArrayIndexExpression(node, arr, index) {
	let length = len(arr["Elements"])
	let idx = index["Value"]

	if idx > length - 1 {
		return indexError(node["StartPos"]["Line"], idx)
	}

	if idx < 0 {
		idx = length + idx
		if idx > length - 1 || idx < 0 {
			return indexError(node["StartPos"]["Line"], idx)
		}
	}

	let item = arr["Elements"]
	return item[idx]
}

fn evalHashIndexExpression(node, h, index) {
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
	let obj = NIL

	let tmpObj = GetGlobalObj(node["Object"]["String"])
	if tmpObj != NIL {
		obj = tmpObj
	} else {
		obj = Eval(node["Object"], scope)
		if isError(obj) {
			return obj
		}
	}

	let call = node["Call"]
	let methName = call["Function"]["String"]

	if obj["CallMethod"] == nil {
		let startPos = node["StartPos"]
		return undefinedMethodError(startPos["Line"], methName, obj["Type"])
	}

	//arguments
	let args = []
	for v in call["Arguments"] {
		let obj = Eval(v, scope)
		if isError(obj) {
			return obj
		}
		args.push(obj)
	}

	let startPos = node["StartPos"]
	return obj["CallMethod"](scope, startPos["Line"], methName, args)
}

//for { block }
fn evalForEverLoopExpression(node, scope) {
	let e = NIL
	let innerScope = NewEnclosedScope(scope)
	for {
		e = Eval(node["Block"], innerScope)
		if isError(e) {
			return e
		}

		if e["Type"] == BREAK_OBJ {
			break
		}
		if e["Type"] == CONTINUE_OBJ {
			continue
		}
		if e["Type"] == RETURN_VALUE_OBJ {
			return e
		}
	}

	if e == NIL || e["Type"] == BREAK_OBJ || e["Type"] == CONTINUE_OBJ {
		return NIL
	}
	return e
}

//for var in iterable
fn evalForArrayExpression(node, scope) {
	let innerScope = NewEnclosedScope(scope)

	let value = Eval(node["Value"], innerScope)
	let valueType = value["Type"]

	if valueType == ERROR_OBJ {
		return value
	}

	//first check if it's a Nil object
	if valueType == NIL_OBJ {
		//return an empty array object
		return ArrayObj([])
	}

	//The object must be 'iterable'
	if value["Iterable"] == nil || value["Iterable"] == false {
		return notIterableError(node["StartPos"]["Line"])
	}

	let members = []
	if valueType == STRING_OBJ {
		for c in value["Value"] {
			members.push(StringObj(c))
		}
	} elif valueType == ARRAY_OBJ {
		members = value["Elements"]
	} elif valueType == TUPLE_OBJ {
		members = value["Elements"]
	}

	let elements = []
	let ret = ArrayObj(elements)
	let result = nil

	for idx, member in members {
		let newSubScope = NewEnclosedScope(innerScope)
		SetIntoScope(newSubScope, "$_", IntegerObj(idx)) // "$_" is idx
		SetIntoScope(newSubScope, node["Var"], member)

		result = Eval(node["Block"], newSubScope)
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
			return result
		} else {
			elements.push(result)
		}
	}

	if result == nil || result["Type"] == BREAK_OBJ || result["Type"] == CONTINUE_OBJ {
		return ret
	}

	return ret
}

fn evalForMapExpression(node, scope) {
	let innerScope = NewEnclosedScope(scope)

	let aValue = Eval(node["Hash"], innerScope)

	if aValue["Type"] == ERROR_OBJ {
		return aValue
	}

	if aValue["Type"] == NIL_OBJ {
		//return an empty array object
		return ArrayObj([])
	}

	//The object must be 'iterable'
	if aValue["Iterable"] == nil || aValue["Iterable"] == false {
		return notIterableError(node["StartPos"]["Line"])
	}

	let elements = []
	let ret = ArrayObj(elements)
	let result = nil

	let pairs = aValue["Pairs"]
	for notUsed, pair in pairs {
		let newSubScope = NewEnclosedScope(innerScope)
		SetIntoScope(newSubScope, node["Key"], pair["Key"])
		SetIntoScope(newSubScope, node["Value"], pair["Value"])

		result = Eval(node["Block"], newSubScope)
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
			return result
		} else {
			elements.push(result)
		}
	}

	if result == nil || result["Type"] == BREAK_OBJ || result["Type"] == CONTINUE_OBJ {
		return ret
	}

	return ret
}

fn evalWhileExpression(node, scope) {
	let innerScope = NewEnclosedScope(scope)

	let result = NIL
	for  {
		let condition = Eval(node["Condition"], innerScope)
		if condition["Type"] == ERROR_OBJ {
			return condition
		}

		if !isTruthy(condition) {
			return NIL
		}

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
			return result
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

		if name["Type"] == IndexExpressionNode {
			let leftVal = Eval(name["Left"], scope)
			if isError(leftVal) {
				return leftVal
			}

			let indexVal = Eval(name["Index"], scope)
			if isError(indexVal) {
				return indexVal
			}

			if leftVal["Type"] == ARRAY_OBJ {
				Array_Set(leftVal, scope, [indexVal, val])
			} elif leftVal["Type"] == TUPLE_OBJ {
				let s = sprintf("%s[IDX]", TUPLE_OBJ)
				return infixOperatorError(node["StartPos"]["Line"], s, node["Token"]["literal"], val["Type"])
			} elif leftVal["Type"] == HASH_OBJ {
				Hash_Push(leftVal, scope, [indexVal, val])
			} elif leftVal["Type"] == STRING_OBJ {
				String_Set(leftVal, scope, [indexVal, val])
			}
		}

		let v = ResetScope(scope, name["String"], val)
		if v != nil {
			return v
		}

		return unknownIdentifierError(node["StartPos"]["Line"], name["String"])
	}
	return val
}


fn evalPostfixExpression(node, left) {
	let operator = node["Operator"]
	if operator == "++" {
		return evalIncrementPostfixOperatorExpression(node, left)
	} elif operator == "--" {
		return evalDecrementPostfixOperatorExpression(node, left)
	}

	return postfixOperatorError(node["StartPos"]["Line"], operator, left["Type"])
}

fn evalIncrementPostfixOperatorExpression(node, left) {
	if left["Type"] == INTEGER_OBJ {
		let left_value = left["Value"]
		let returnVal = IntegerObj(left_value)
		left["Value"] = left_value + 1
		left["Inspect"] = sprintf("%d", left_value + 1)
		return returnVal
	} elif left["Type"] == FLOAT_OBJ {
		let left_value = left["Value"]
		let returnVal = FloatObj(left_value)
		left["Value"] = left_value + 1
		left["Inspect"] = sprintf("%g", left_value + 1)
		return returnVal
	} 

	return postfixOperatorError(node["StartPos"]["Line"], "++", left["Type"])
}

fn evalDecrementPostfixOperatorExpression(node, left) {
	if left["Type"] == INTEGER_OBJ {
		let left_value = left["Value"]
		let returnVal = IntegerObj(left_value)
		left["Value"] = left_value - 1
		left["Inspect"] = sprintf("%d", left["Value"])
		return returnVal
	} elif left["Type"] == FLOAT_OBJ {
		let left_value = left["Value"]
		let returnVal = FloatObj(left_value)
		left["Value"] = left_value - 1
		left["Inspect"] = sprintf("%g", left_value - 1)

		return returnVal
	}

	return postfixOperatorError(node["StartPos"]["Line"], "--", left["Type"])
}

fn evalFunctionLiteral(node, scope) {
	let fn_params = node["Parameters"]
	let fn_body   = node["Body"]
	let fn_name   = node["Token"]["literal"]
	return FunctionObj(fn_params, fn_body, fn_name, scope)
}

fn evalFunctionStatement(node, scope) {
	let fnObj = evalFunctionLiteral(node["FunctionLiteral"], scope)

	SetIntoScope(scope, node["Name"]["Value"], fnObj) //save to scope

	return fnObj
}

fn evalRegExLiteral(node, scope) {
	regObj = regexp.compile(node["Value"])
	return RegExObj(regObj, node["Value"])
}

fn isTruthy(obj) {
	if isError(obj) {
		return false
	}

	if obj == NIL {
		return false
	} elif obj == TRUE {
		return true
	} elif obj == FALSE {
		return false
	}

	if obj["Type"] == BOOLEAN_OBJ {
		return obj["Value"]
	}

	if obj["Type"] == INTEGER_OBJ {
		if obj["Value"] == 0 {
			return false
		}
		return true
	}

	if obj["Type"] == FLOAT_OBJ {
		if obj["Value"] == 0.0 {
			return false
		}
		return true
	}

	return true
}

fn isError(obj) {
	if obj != nil {
		return obj["Type"] == ERROR_OBJ
	}

	return false
}

fn applyFunction(node, func, args) {
	let funcType = func["Type"]
	if funcType == FUNCTION_OBJ {
		let extendedScope = extendFunctionScope(func, args)
		let evaluated = Eval(func["Body"], extendedScope)
		if isError(evaluated) {
			return evaluated
		}

		return unwrapReturnValue(evaluated)
	} elif funcType == BUILTIN_OBJ {
		let builtin_fn = func["BuiltinFn"]
		let startPos = node["StartPos"]
		return builtin_fn(startPos["Line"], args)
	}

	return callableError(node["StartPos"]["Line"], node["Function"]["Token"]["literal"])
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
		if len(obj["Values"]) > 1 {
			//if the function returns multiple values, then we wrap it in a tuple object.
			return TupleObj(obj["Values"], true)
		}
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
	let filename = "examples/simple_interpreter/example.txt"
	let source = ioutil.readFile(filename)
	/* lexer */
	NewLexer(source, filename)

/* for debugging : print each tokens */
//	do {
//		n = lex_nextToken()
//		printf("{tokentype: \"%s\", literal: \"%s\", pos : %s}\n", tokenString(n["type"]), n["literal"], positionStr(n["position"]))
//		if n["type"] == TOKEN_EOF {
//			break
//		}
//	}

	/* parser */
	let p = ParseProgram()
	if len(parser_Errors) != 0 {
		for errMsg in parser_Errors {
			println("\t" + errMsg + "\n")
		}
		os.exit(1)
	}

	//for debugging
//	println("==========================================")
//	println(p["String"]) //print the 'program' node
//	println("==========================================")

	/* evaluator */
	let root_scope = NewScope()
	obj = Eval(p, root_scope)
	if obj["Type"] == ERROR_OBJ {
		println(obj["Inspect"])
	}
}

