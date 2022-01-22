package ast

import (
	"bytes"
	"magpie/token"
	"strings"
	"unicode/utf8"
)

//Source interface is used in documentation for printing source code.
type Source interface {
	SrcStart() token.Position
	SrcEnd() token.Position
}

type Node interface {
	Pos() token.Position // position of first character belonging to the node
	End() token.Position // position of first character immediately after the node

	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
	Imports    map[string]*ImportStatement
}

func (p *Program) Pos() token.Position {
	if len(p.Statements) > 0 {
		return p.Statements[0].Pos()
	}
	return token.Position{}
}

func (p *Program) End() token.Position {
	aLen := len(p.Statements)
	if aLen > 0 {
		return p.Statements[aLen-1].End()
	}
	return token.Position{}
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token       token.Token
	Statements  []Statement
	RBraceToken token.Token
}

func (bs *BlockStatement) Pos() token.Position {
	return bs.Token.Pos

}

//func (bs *BlockStatement) End() token.Position {
//	aLen := len(bs.Statements)
//	if aLen > 0 {
//		return bs.Statements[aLen-1].End()
//	}
//	return bs.Token.Pos
//}

func (bs *BlockStatement) End() token.Position {
	return token.Position{Filename: bs.Token.Pos.Filename, Line: bs.RBraceToken.Pos.Line, Col: bs.RBraceToken.Pos.Col + 1}
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		str := s.String()

		out.WriteString(str)
		if str[len(str)-1:] != ";" {
			out.WriteString(";")
		}
	}

	return out.String()
}

///////////////////////////////////////////////////////////
//                        FOR LOOP                       //
///////////////////////////////////////////////////////////
type ForLoop struct {
	Token  token.Token
	Init   Expression
	Cond   Expression
	Update Expression
	Block  Node //BlockStatement or single expression
}

func (fl *ForLoop) Pos() token.Position {
	return fl.Token.Pos
}

func (fl *ForLoop) End() token.Position {
	return fl.Block.End()
}

func (fl *ForLoop) expressionNode()      {}
func (fl *ForLoop) TokenLiteral() string { return fl.Token.Literal }

func (fl *ForLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for")
	out.WriteString(" ( ")
	out.WriteString(fl.Init.String())
	out.WriteString(" ; ")
	out.WriteString(fl.Cond.String())
	out.WriteString(" ; ")
	out.WriteString(fl.Update.String())
	out.WriteString(" ) ")
	out.WriteString(" { ")
	out.WriteString(fl.Block.String())
	out.WriteString(" }")

	return out.String()
}

type ForEachArrayLoop struct {
	Token token.Token
	Var   string
	Value Expression //value to range over
	Cond  Expression //conditional clause(nil if there is no 'WHERE' clause)
	Block Node       //BlockStatement or single expression
}

func (fal *ForEachArrayLoop) Pos() token.Position {
	return fal.Token.Pos
}

func (fal *ForEachArrayLoop) End() token.Position {
	return fal.Block.End()
}

func (fal *ForEachArrayLoop) expressionNode()      {}
func (fal *ForEachArrayLoop) TokenLiteral() string { return fal.Token.Literal }

func (fal *ForEachArrayLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fal.Var)
	out.WriteString(" in ")
	out.WriteString(fal.Value.String())
	if fal.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(fal.Cond.String())
	}
	out.WriteString(" { ")
	out.WriteString(fal.Block.String())
	out.WriteString(" }")

	return out.String()
}

type ForEachMapLoop struct {
	Token token.Token
	Key   string
	Value string
	X     Expression //value to range over
	Cond  Expression //Conditional clause(nil if there is no 'WHERE' clause)
	Block Node       //BlockStatement or single expression
}

func (fml *ForEachMapLoop) Pos() token.Position {
	return fml.Token.Pos
}

func (fml *ForEachMapLoop) End() token.Position {
	return fml.Block.End()
}

func (fml *ForEachMapLoop) expressionNode()      {}
func (fml *ForEachMapLoop) TokenLiteral() string { return fml.Token.Literal }

func (fml *ForEachMapLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fml.Key + ", " + fml.Value)
	out.WriteString(" in ")
	out.WriteString(fml.X.String())
	if fml.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(fml.Cond.String())
	}
	out.WriteString(" { ")
	out.WriteString(fml.Block.String())
	out.WriteString(" }")

	return out.String()
}

type ForEverLoop struct {
	Token token.Token
	Block *BlockStatement
}

func (fel *ForEverLoop) Pos() token.Position {
	return fel.Token.Pos
}

func (fel *ForEverLoop) End() token.Position {
	return fel.Block.End()
}

func (fel *ForEverLoop) expressionNode()      {}
func (fel *ForEverLoop) TokenLiteral() string { return fel.Token.Literal }

func (fel *ForEverLoop) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(" { ")
	out.WriteString(fel.Block.String())
	out.WriteString(" }")

	return out.String()
}

//for i in start..end <where cond> { }
type ForEachDotRange struct {
	Token    token.Token
	Var      string
	StartIdx Expression
	EndIdx   Expression
	Cond     Expression //conditional clause(nil if there is no 'WHERE' clause)
	Block    Node       //BlockStatement or single expression
}

func (fdr *ForEachDotRange) Pos() token.Position {
	return fdr.Token.Pos
}

func (fdr *ForEachDotRange) End() token.Position {
	return fdr.Block.End()
}

func (fdr *ForEachDotRange) expressionNode()      {}
func (fdr *ForEachDotRange) TokenLiteral() string { return fdr.Token.Literal }

func (fdr *ForEachDotRange) String() string {
	var out bytes.Buffer

	out.WriteString("for ")
	out.WriteString(fdr.Var)
	out.WriteString(" in ")
	out.WriteString(fdr.StartIdx.String())
	out.WriteString(" .. ")
	out.WriteString(fdr.EndIdx.String())
	if fdr.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(fdr.Cond.String())
	}
	out.WriteString(" { ")
	out.WriteString(fdr.Block.String())
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                        WHILE LOOP                     //
///////////////////////////////////////////////////////////
type WhileLoop struct {
	Token     token.Token
	Condition Expression
	Block     Node //BlockStatement or single expression
}

func (wl *WhileLoop) Pos() token.Position {
	return wl.Token.Pos
}

func (wl *WhileLoop) End() token.Position {
	return wl.Block.End()
}

func (wl *WhileLoop) expressionNode()      {}
func (wl *WhileLoop) TokenLiteral() string { return wl.Token.Literal }

func (wl *WhileLoop) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(wl.Condition.String())
	out.WriteString("{")
	out.WriteString(wl.Block.String())
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         DO LOOP                       //
///////////////////////////////////////////////////////////
type DoLoop struct {
	Token token.Token
	Block *BlockStatement
}

func (dl *DoLoop) Pos() token.Position {
	return dl.Token.Pos
}

func (dl *DoLoop) End() token.Position {
	return dl.Block.End()
}

func (dl *DoLoop) expressionNode()      {}
func (dl *DoLoop) TokenLiteral() string { return dl.Token.Literal }

func (dl *DoLoop) String() string {
	var out bytes.Buffer

	out.WriteString("do")
	out.WriteString(" { ")
	out.WriteString(dl.Block.String())
	out.WriteString(" }")
	return out.String()
}

///////////////////////////////////////////////////////////
//                        IDENTIFIER                     //
///////////////////////////////////////////////////////////
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) Pos() token.Position {
	return i.Token.Pos
}

func (i *Identifier) End() token.Position {
	length := utf8.RuneCountInString(i.Value)
	return token.Position{Filename: i.Token.Pos.Filename, Line: i.Token.Pos.Line, Col: i.Token.Pos.Col + length}
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

///////////////////////////////////////////////////////////
//                IFELSE MACRO  STATEMENT                //
///////////////////////////////////////////////////////////
type IfMacroStatement struct {
	Token        token.Token
	Condition    bool
	ConditionStr string
	Consequence  *BlockStatement
	Alternative  *BlockStatement
}

func (ifex *IfMacroStatement) Pos() token.Position {
	return ifex.Token.Pos
}

func (ifex *IfMacroStatement) End() token.Position {
	if ifex.Alternative != nil {
		return ifex.Alternative.End()
	}
	return ifex.Consequence.End()
}

func (ifex *IfMacroStatement) statementNode()       {}
func (ifex *IfMacroStatement) TokenLiteral() string { return ifex.Token.Literal }

func (ifex *IfMacroStatement) String() string {
	var out bytes.Buffer

	out.WriteString("#ifdef ")
	out.WriteString(ifex.ConditionStr + " ")
	out.WriteString(ifex.Consequence.String())
	if ifex.Alternative != nil {
		out.WriteString(" #else ")
		out.WriteString(ifex.Alternative.String())
	}
	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Conditions  []*IfConditionExpr //if or elif part
	Alternative Node               //else part(BlockStatement or single ExpressionStatement)
}

func (ifex *IfExpression) Pos() token.Position {
	return ifex.Token.Pos
}

func (ifex *IfExpression) End() token.Position {
	if ifex.Alternative != nil {
		return ifex.Alternative.End()
	}

	aLen := len(ifex.Conditions)
	return ifex.Conditions[aLen-1].End()
}

func (ifex *IfExpression) expressionNode()      {}
func (ifex *IfExpression) TokenLiteral() string { return ifex.Token.Literal }

func (ifex *IfExpression) String() string {
	var out bytes.Buffer

	for i, c := range ifex.Conditions {
		if i == 0 {
			out.WriteString("if ")
		} else {
			out.WriteString("elif ")
		}
		out.WriteString(c.String())
	}

	if ifex.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(" { ")
		out.WriteString(ifex.Alternative.String())
		out.WriteString(" }")
	}

	return out.String()
}

//if/else-if condition
type IfConditionExpr struct {
	Token token.Token
	Cond  Expression //condition
	Body  Node       //body(BlockStatement or single ExpressionStatement)
}

func (ic *IfConditionExpr) Pos() token.Position {
	return ic.Token.Pos
}

func (ic *IfConditionExpr) End() token.Position {
	return ic.Body.End()
}

func (ic *IfConditionExpr) expressionNode()      {}
func (ic *IfConditionExpr) TokenLiteral() string { return ic.Token.Literal }

func (ic *IfConditionExpr) String() string {
	var out bytes.Buffer

	out.WriteString(ic.Cond.String())
	out.WriteString(" { ")
	out.WriteString(ic.Body.String())
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                    UNLESS-ELSE                        //
///////////////////////////////////////////////////////////
type UnlessExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ul *UnlessExpression) Pos() token.Position {
	return ul.Token.Pos
}

func (ul *UnlessExpression) End() token.Position {
	if ul.Alternative != nil {
		return ul.Alternative.End()
	}
	return ul.Consequence.End()
}

func (ul *UnlessExpression) expressionNode()      {}
func (ul *UnlessExpression) TokenLiteral() string { return ul.Token.Literal }

func (ul *UnlessExpression) String() string {
	var out bytes.Buffer

	out.WriteString("unless ")
	out.WriteString("(")
	out.WriteString(ul.Condition.String())
	out.WriteString(")")
	out.WriteString(" { ")
	out.WriteString(ul.Consequence.String())
	out.WriteString(" }")
	if ul.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(" { ")
		out.WriteString(ul.Alternative.String())
		out.WriteString(" }")
	}

	return out.String()
}

///////////////////////////////////////////////////////////
//                         HASH LITERAL                  //
///////////////////////////////////////////////////////////
type HashLiteral struct {
	Token       token.Token
	Order       []Expression //For keeping the order of the hash key
	Pairs       map[Expression]Expression
	RBraceToken token.Token
}

func (h *HashLiteral) Pos() token.Position {
	return h.Token.Pos
}

//func (h *HashLiteral) End() token.Position {
//	maxLineMap := make(map[int]Expression)
//
//	for _, value := range h.Pairs {
//		v := value.(Expression)
//		maxLineMap[v.End().Line] = v
//	}
//
//	maxLine := 0
//	for line, _ := range maxLineMap {
//		if line > maxLine {
//			maxLine = line
//		}
//	}
//
//	ret := maxLineMap[maxLine].(Expression)
//	return ret.End()
//}

func (h *HashLiteral) End() token.Position {
	return token.Position{Filename: h.Token.Pos.Filename, Line: h.RBraceToken.Pos.Line, Col: h.RBraceToken.Pos.Col + 1}
}

func (h *HashLiteral) expressionNode()      {}
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	//for key, value := range h.Pairs {
	//	pairs = append(pairs, key.String()+": "+value.String())
	//}
	for _, key := range h.Order {
		value, _ := h.Pairs[key]
		pairs = append(pairs, key.String()+": "+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     NIL LITERAL                   //
///////////////////////////////////////////////////////////
type NilLiteral struct {
	Token token.Token
}

func (n *NilLiteral) Pos() token.Position {
	return n.Token.Pos
}

func (n *NilLiteral) End() token.Position {
	length := len(n.Token.Literal)
	pos := n.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (n *NilLiteral) expressionNode()      {}
func (n *NilLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NilLiteral) String() string       { return n.Token.Literal }

///////////////////////////////////////////////////////////
//                     INTEGER LITERAL                   //
///////////////////////////////////////////////////////////
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) Pos() token.Position {
	return il.Token.Pos
}

func (il *IntegerLiteral) End() token.Position {
	length := utf8.RuneCountInString(il.Token.Literal)
	pos := il.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

///////////////////////////////////////////////////////////
//               UNSIGNED INTEGER LITERAL                //
///////////////////////////////////////////////////////////
type UIntegerLiteral struct { //U: Unsigned
	Token token.Token
	Value uint64
}

func (il *UIntegerLiteral) Pos() token.Position {
	return il.Token.Pos
}

func (il *UIntegerLiteral) End() token.Position {
	length := utf8.RuneCountInString(il.Token.Literal)
	pos := il.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (il *UIntegerLiteral) expressionNode()      {}
func (il *UIntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *UIntegerLiteral) String() string       { return il.Token.Literal }

///////////////////////////////////////////////////////////
//                     FLOAT LITERAL                     //
///////////////////////////////////////////////////////////
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) Pos() token.Position {
	return fl.Token.Pos
}

func (fl *FloatLiteral) End() token.Position {
	length := utf8.RuneCountInString(fl.Token.Literal)
	pos := fl.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

///////////////////////////////////////////////////////////
//                     BOOLEAN LITERAL                   //
///////////////////////////////////////////////////////////
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) Pos() token.Position {
	return b.Token.Pos
}

func (b *Boolean) End() token.Position {
	length := utf8.RuneCountInString(b.Token.Literal)
	pos := b.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

///////////////////////////////////////////////////////////
//                      REGEX LITERAL                    //
///////////////////////////////////////////////////////////
type RegExLiteral struct {
	Token token.Token
	Value string
}

func (rel *RegExLiteral) Pos() token.Position {
	return rel.Token.Pos
}

//func (rel *RegExLiteral) End() token.Position {
//	return rel.Token.Pos
//}

func (rel *RegExLiteral) End() token.Position {
	length := utf8.RuneCountInString(rel.Token.Literal)
	pos := rel.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}

	return rel.Token.Pos
}

func (rel *RegExLiteral) expressionNode()      {}
func (rel *RegExLiteral) TokenLiteral() string { return rel.Token.Literal }
func (rel *RegExLiteral) String() string       { return rel.Value }

///////////////////////////////////////////////////////////
//                      ARRAY LITERAL                    //
///////////////////////////////////////////////////////////
type ArrayLiteral struct {
	Token         token.Token
	Members       []Expression
	CreationCount *IntegerLiteral
}

func (a *ArrayLiteral) Pos() token.Position {
	return a.Token.Pos
}

func (a *ArrayLiteral) End() token.Position {
	aLen := len(a.Members)
	if aLen > 0 {
		return a.Members[aLen-1].End()
	}
	if a.CreationCount == nil {
		ret := a.Token.Pos
		ret.Col = ret.Col + 1
		return ret
	}
	return a.CreationCount.End()
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	members := []string{}
	for _, m := range a.Members {
		members = append(members, m.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(members, ", "))
	out.WriteString("]")
	if a.CreationCount != nil {
		out.WriteString(a.CreationCount.String())
	}
	return out.String()
}

/////////////////////////////////////////////////////////
//                     RANGE LITERAL(..)                //
/////////////////////////////////////////////////////////
type RangeLiteral struct {
	Token    token.Token
	StartIdx Expression
	EndIdx   Expression
}

func (r *RangeLiteral) Pos() token.Position {
	return r.Token.Pos
}

func (r *RangeLiteral) End() token.Position {
	return r.EndIdx.End()
}

func (r *RangeLiteral) expressionNode()      {}
func (r *RangeLiteral) TokenLiteral() string { return r.Token.Literal }
func (r *RangeLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(r.StartIdx.String())
	out.WriteString(" .. ")
	out.WriteString(r.EndIdx.String())
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     FUNCTION LITERAL                  //
///////////////////////////////////////////////////////////
type FunctionLiteral struct {
	Token      token.Token
	Parameters []Expression
	Body       *BlockStatement

	//Default values
	Values map[string]Expression

	Variadic bool

	StaticFlag    bool
	ModifierLevel ModifierLevel //for 'class' use

	//If the function is async or not
	Async bool
}

func (fl *FunctionLiteral) Pos() token.Position {
	return fl.Token.Pos
}

func (fl *FunctionLiteral) End() token.Position {
	return fl.Body.End()
}

// For debugger use
func (fl *FunctionLiteral) StmtPos() token.Position {
	if len(fl.Body.Statements) > 0 {
		return fl.Body.Statements[0].Pos()
	}
	return fl.Token.Pos
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	if fl.Async {
		out.WriteString("async ")
	}

	out.WriteString(fl.ModifierLevel.String())
	if fl.StaticFlag {
		out.WriteString("static ")
	}

	out.WriteString(fl.TokenLiteral())
	params := []string{}
	for i, p := range fl.Parameters {
		param := p.String()
		if fl.Variadic && i == len(fl.Parameters)-1 {
			param = "..." + param
		}

		params = append(params, p.String())

	}
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString("{ ")
	out.WriteString(fl.Body.String())
	out.WriteString(" }")
	return out.String()
}

///////////////////////////////////////////////////////////
//                  FUNCTION STATEMENT                   //
///////////////////////////////////////////////////////////
type FunctionStatement struct {
	Token           token.Token
	Name            *Identifier
	FunctionLiteral *FunctionLiteral
	Annotations     []*AnnotationStmt
	IsServiceAnno   bool //service annotation(@route) is processed differently
	//Doc related
	Doc         *CommentGroup // associated documentation; or nil
	SrcEndToken token.Token   //used for printing source code
}

func (f *FunctionStatement) Pos() token.Position {
	return f.Token.Pos
}

func (f *FunctionStatement) End() token.Position {
	return f.FunctionLiteral.Body.End()
}

//Below two methods implements 'Source' interface.
func (f *FunctionStatement) SrcStart() token.Position {
	return f.Pos()
}

func (f *FunctionStatement) SrcEnd() token.Position {
	ret := f.SrcEndToken.Pos
	length := utf8.RuneCountInString(f.SrcEndToken.Literal)
	ret.Offset += length
	return ret
}

func (f *FunctionStatement) statementNode()       {}
func (f *FunctionStatement) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionStatement) String() string {
	var out bytes.Buffer

	for _, anno := range f.Annotations { //for each annotation
		out.WriteString(anno.String())
	}

	out.WriteString(f.FunctionLiteral.ModifierLevel.String())

	out.WriteString(" fn ")
	out.WriteString(f.Name.String())

	params := []string{}
	for i, p := range f.FunctionLiteral.Parameters {
		param := p.String()
		if f.FunctionLiteral.Variadic && i == len(f.FunctionLiteral.Parameters)-1 {
			param = "..." + param
		}

		params = append(params, p.String())

	}
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString("{ ")
	out.WriteString(f.FunctionLiteral.Body.String())
	out.WriteString(" }")

	return out.String()
}

func (f *FunctionStatement) Docs() string {
	var out bytes.Buffer

	out.WriteString(f.FunctionLiteral.ModifierLevel.String())

	out.WriteString("fn ")
	out.WriteString(f.Name.String())

	params := []string{}
	for i, p := range f.FunctionLiteral.Parameters {
		param := p.String()
		if f.FunctionLiteral.Variadic && i == len(f.FunctionLiteral.Parameters)-1 {
			param = "..." + param
		}

		params = append(params, p.String())

	}
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	return out.String()
}

///////////////////////////////////////////////////////////
//                      STRING LITERAL                   //
///////////////////////////////////////////////////////////
type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) Pos() token.Position {
	return s.Token.Pos
}

func (s *StringLiteral) End() token.Position {
	length := utf8.RuneCountInString(s.Value)
	return token.Position{Filename: s.Token.Pos.Filename, Line: s.Token.Pos.Line, Col: s.Token.Pos.Col + length}
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return s.Token.Literal }

///////////////////////////////////////////////////////////
//                  INTERPOLATED STRING                  //
///////////////////////////////////////////////////////////
type InterpolatedString struct {
	Token   token.Token
	Value   string
	ExprMap map[byte]Expression
}

func (is *InterpolatedString) Pos() token.Position {
	return is.Token.Pos
}

func (is *InterpolatedString) End() token.Position {
	length := utf8.RuneCountInString(is.Value)
	return token.Position{Filename: is.Token.Pos.Filename, Line: is.Token.Pos.Line, Col: is.Token.Pos.Col + length}
}

func (is *InterpolatedString) expressionNode()      {}
func (is *InterpolatedString) TokenLiteral() string { return is.Token.Literal }
func (is *InterpolatedString) String() string       { return is.Token.Literal }

///////////////////////////////////////////////////////////
//                    TRY/CATCH/FINALLY                  //
///////////////////////////////////////////////////////////
//TryStmt provide "try/catch/finally" statement.
type TryStmt struct {
	Token   token.Token
	Try     *BlockStatement
	Var     string
	Catch   *BlockStatement
	Finally *BlockStatement
}

func (t *TryStmt) Pos() token.Position {
	return t.Token.Pos
}

func (t *TryStmt) End() token.Position {
	if t.Finally != nil {
		return t.Finally.End()
	}

	return t.Catch.End()
}

func (t *TryStmt) statementNode()       {}
func (t *TryStmt) TokenLiteral() string { return t.Token.Literal }

func (t *TryStmt) String() string {
	var out bytes.Buffer

	out.WriteString("try { ")
	out.WriteString(t.Try.String())
	out.WriteString(" }")

	if t.Catch != nil {
		if len(t.Var) > 0 {
			out.WriteString(" catch " + t.Var + " { ")
		} else {
			out.WriteString(" catch { ")
		}
		out.WriteString(t.Catch.String())
		out.WriteString(" }")
	}

	if t.Finally != nil {
		out.WriteString(" finally { ")
		out.WriteString(t.Finally.String())
		out.WriteString(" }")
	}

	return out.String()
}

//throw <expression>
type ThrowStmt struct {
	Token token.Token
	Expr  Expression
}

func (ts *ThrowStmt) Pos() token.Position {
	return ts.Token.Pos
}

func (ts *ThrowStmt) End() token.Position {
	return ts.Expr.End()
}

func (ts *ThrowStmt) statementNode()       {}
func (ts *ThrowStmt) TokenLiteral() string { return ts.Token.Literal }

func (ts *ThrowStmt) String() string {
	var out bytes.Buffer

	out.WriteString("throw ")
	out.WriteString(ts.Expr.String())
	out.WriteString(";")

	return out.String()
}

///////////////////////////////////////////////////////////
//                      STRUCT LITERAL                   //
///////////////////////////////////////////////////////////
type StructLiteral struct {
	Token       token.Token
	Pairs       map[Expression]Expression
	RBraceToken token.Token
}

func (s *StructLiteral) Pos() token.Position {
	return s.Token.Pos
}

//func (s *StructLiteral) End() token.Position {
//	maxLineMap := make(map[int]Expression)
//
//	for _, value := range s.Pairs {
//		v := value.(Expression)
//		maxLineMap[v.End().Line] = v
//	}
//
//	maxLine := 0
//	for line, _ := range maxLineMap {
//		if line > maxLine {
//			maxLine = line
//		}
//	}
//
//	ret := maxLineMap[maxLine].(Expression)
//	return ret.End()
//}

func (s *StructLiteral) End() token.Position {
	return token.Position{Filename: s.Token.Pos.Filename, Line: s.RBraceToken.Pos.Line, Col: s.RBraceToken.Pos.Col + 1}
}

func (s *StructLiteral) expressionNode()      {}
func (s *StructLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StructLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range s.Pairs {
		pairs = append(pairs, key.String()+"=>"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     DEFER STATEMENT                   //
///////////////////////////////////////////////////////////
type DeferStmt struct {
	Token token.Token
	Call  Expression
}

func (ds *DeferStmt) Pos() token.Position {
	return ds.Token.Pos
}

func (ds *DeferStmt) End() token.Position {
	return ds.Call.End()
}

func (ds *DeferStmt) statementNode()       {}
func (ds *DeferStmt) TokenLiteral() string { return ds.Token.Literal }

func (ds *DeferStmt) String() string {
	var out bytes.Buffer

	out.WriteString(ds.TokenLiteral() + " ")
	out.WriteString(ds.Call.String())
	out.WriteString("; ")

	return out.String()
}

///////////////////////////////////////////////////////////
//                    RETURN STATEMENT                   //
///////////////////////////////////////////////////////////
type ReturnStatement struct {
	Token        token.Token
	ReturnValue  Expression //for old campatibility
	ReturnValues []Expression
}

func (rs *ReturnStatement) Pos() token.Position {
	return rs.Token.Pos
}

func (rs *ReturnStatement) End() token.Position {
	aLen := len(rs.ReturnValues)
	if aLen > 0 {
		return rs.ReturnValues[aLen-1].End()
	}

	return token.Position{Filename: rs.Token.Pos.Filename, Line: rs.Token.Pos.Line, Col: rs.Token.Pos.Col + len(rs.Token.Literal)}
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	//	if rs.ReturnValue != nil {
	//		out.WriteString(rs.ReturnValue.String())
	//	}

	values := []string{}
	for _, value := range rs.ReturnValues {
		values = append(values, value.String())
	}
	out.WriteString(strings.Join(values, ", "))

	out.WriteString(";")

	return out.String()
}

///////////////////////////////////////////////////////////
//                  EXPRESSION STATEMENT                 //
///////////////////////////////////////////////////////////
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) Pos() token.Position {
	return es.Token.Pos
}

func (es *ExpressionStatement) End() token.Position {
	return es.Expression.End()
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		str := es.Expression.String()
		if str[len(str)-1:] != ";" {
			str = str + ";"
		}
		return str
	}
	return ""
}

///////////////////////////////////////////////////////////
//                      LET STATEMENT                    //
///////////////////////////////////////////////////////////
type LetStatement struct {
	Token  token.Token
	Names  []*Identifier
	Values []Expression

	StaticFlag    bool
	ModifierLevel ModifierLevel //used in 'class'
	Annotations   []*AnnotationStmt

	//Doc related
	Doc         *CommentGroup // associated documentation; or nil
	SrcEndToken token.Token

	//destructuring assigment flag
	DestructingFlag bool

	//For debugger use, If the LetStatement is in a class declaration,
	//we do not want the debugger to stop at it.
	InClass bool //true if the LetStatement is in a Class declaration
}

func (ls *LetStatement) Pos() token.Position {
	return ls.Token.Pos
}

func (ls *LetStatement) End() token.Position {
	aLen := len(ls.Values)
	if aLen > 0 {
		return ls.Values[aLen-1].End()
	}

	return ls.Names[0].End()
}

//Below two methods implements 'Source' interface.
func (ls *LetStatement) SrcStart() token.Position {
	return ls.Pos()
}

func (ls *LetStatement) SrcEnd() token.Position {
	ret := ls.SrcEndToken.Pos
	length := utf8.RuneCountInString(ls.SrcEndToken.Literal)
	ret.Offset += length
	return ret
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.ModifierLevel.String())
	if ls.StaticFlag {
		out.WriteString("static ")
	}

	out.WriteString(ls.TokenLiteral() + " ")
	if ls.DestructingFlag {
		out.WriteString("(")
	}

	names := []string{}
	for _, name := range ls.Names {
		names = append(names, name.String())
	}
	out.WriteString(strings.Join(names, ", "))

	if ls.DestructingFlag {
		out.WriteString(")")
	}

	if len(ls.Values) == 0 { //e.g. 'let x'
		out.WriteString(";")
		return out.String()
	}

	out.WriteString(" = ")

	values := []string{}
	for _, value := range ls.Values {
		values = append(values, value.String())
	}
	out.WriteString(strings.Join(values, ", "))

	return out.String()
}

func (ls *LetStatement) Docs() string {
	return ls.String()
}

///////////////////////////////////////////////////////////
//                    CONST STATEMENT                    //
///////////////////////////////////////////////////////////
type ConstStatement struct {
	Token token.Token
	Name  []*Identifier
	Value []Expression

	StaticFlag    bool
	ModifierLevel ModifierLevel //used in 'class'
	Annotations   []*AnnotationStmt

	//Doc related
	Doc         *CommentGroup // associated documentation; or nil
	SrcEndToken token.Token
}

func (cs *ConstStatement) Pos() token.Position {
	if len(cs.Name) > 0 {
		return cs.Name[0].Pos()
	}

	return cs.Token.Pos
}

func (cs *ConstStatement) End() token.Position {
	aLen := len(cs.Value)
	if aLen > 0 {
		return cs.Value[aLen-1].End()
	}
	return token.Position{}
}

//Below two methods implements 'Source' interface.
func (cs *ConstStatement) SrcStart() token.Position {
	return cs.Pos()
}

func (cs *ConstStatement) SrcEnd() token.Position {
	ret := cs.SrcEndToken.Pos
	length := utf8.RuneCountInString(cs.SrcEndToken.Literal)
	ret.Offset += length
	return ret
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Literal }

func (cs *ConstStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.ModifierLevel.String())
	if cs.StaticFlag {
		out.WriteString("static ")
	}

	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(" ( ")

	for idx, name := range cs.Name {
		out.WriteString(name.TokenLiteral())
		out.WriteString(" = ")
		if cs.Value[idx] != nil {
			out.WriteString(cs.Value[idx].String())
		}
		out.WriteString(",")
	}

	out.WriteString(" )")
	return out.String()
}

func (cs *ConstStatement) Docs() string {
	return cs.String()
}

///////////////////////////////////////////////////////////
//                      IMPORT STATEMENT                //
///////////////////////////////////////////////////////////
type ImportStatement struct {
	Token      token.Token
	ImportPath string
	Program    *Program
	Functions  map[string]*FunctionLiteral //for debugger usage
}

func (is *ImportStatement) Pos() token.Position {
	return is.Token.Pos
}

func (is *ImportStatement) End() token.Position {
	length := utf8.RuneCountInString(is.ImportPath)
	return token.Position{Filename: is.Token.Pos.Filename, Line: is.Token.Pos.Line, Col: is.Token.Pos.Col + length}
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(is.ImportPath)

	return out.String()
}

///////////////////////////////////////////////////////////
//                         BREAK                         //
///////////////////////////////////////////////////////////
type BreakExpression struct {
	Token token.Token
}

func (be *BreakExpression) Pos() token.Position {
	return be.Token.Pos
}

func (be *BreakExpression) End() token.Position {
	length := utf8.RuneCountInString(be.Token.Literal)
	pos := be.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (be *BreakExpression) expressionNode()      {}
func (be *BreakExpression) TokenLiteral() string { return be.Token.Literal }

func (be *BreakExpression) String() string { return be.Token.Literal }

///////////////////////////////////////////////////////////
//                         CONTINUE                      //
///////////////////////////////////////////////////////////
type ContinueExpression struct {
	Token token.Token
}

func (ce *ContinueExpression) Pos() token.Position {
	return ce.Token.Pos
}

func (ce *ContinueExpression) End() token.Position {
	length := utf8.RuneCountInString(ce.Token.Literal)
	pos := ce.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (ce *ContinueExpression) expressionNode()      {}
func (ce *ContinueExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *ContinueExpression) String() string { return ce.Token.Literal }

///////////////////////////////////////////////////////////
//                         ASSIGN                        //
///////////////////////////////////////////////////////////
type AssignExpression struct {
	Token token.Token
	Name  Expression
	Value Expression
}

func (ae *AssignExpression) Pos() token.Position {
	//return ae.Token.Pos
	return ae.Name.Pos()
}

func (ae *AssignExpression) End() token.Position {
	return ae.Value.End()
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }

func (ae *AssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Name.String())
	//out.WriteString(" = ")
	out.WriteString(ae.Token.Literal)
	out.WriteString(ae.Value.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                         GREP                          //
///////////////////////////////////////////////////////////
type GrepExpr struct {
	Token token.Token
	Var   string          //Name is "$_"
	Value Expression      //value to range over
	Block *BlockStatement //Grep Block, may be nil
	Expr  Expression      //Grep Expr, may be nil
}

func (ge *GrepExpr) Pos() token.Position {
	return ge.Token.Pos
}

func (ge *GrepExpr) End() token.Position {
	if ge.Block == nil {
		return ge.Expr.End()
	}
	if ge.Expr == nil {
		return ge.Block.End()
	}
	return ge.Token.Pos //should never happen
}

func (ge *GrepExpr) expressionNode()      {}
func (ge *GrepExpr) TokenLiteral() string { return ge.Token.Literal }

func (ge *GrepExpr) String() string {
	var out bytes.Buffer

	out.WriteString("grep ")
	if ge.Block != nil {
		out.WriteString(" { ")
		out.WriteString(ge.Block.String())
		out.WriteString(" } ")
	} else {
		out.WriteString(ge.Expr.String())
		out.WriteString(" , ")
	}

	out.WriteString(ge.Value.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                         MAP                           //
///////////////////////////////////////////////////////////
type MapExpr struct {
	Token token.Token
	Var   string          //Name is "$_"
	Value Expression      //value to range over
	Block *BlockStatement //Grep Block, may be nil
	Expr  Expression      //Grep Expr, may be nil
}

func (me *MapExpr) Pos() token.Position {
	return me.Token.Pos
}

func (me *MapExpr) End() token.Position {
	if me.Block == nil {
		return me.Expr.End()
	}
	if me.Expr == nil {
		return me.Block.End()
	}
	return me.Token.Pos //should never happen
}

func (me *MapExpr) expressionNode()      {}
func (me *MapExpr) TokenLiteral() string { return me.Token.Literal }

func (me *MapExpr) String() string {
	var out bytes.Buffer

	out.WriteString("map ")
	if me.Block != nil {
		out.WriteString(" { ")
		out.WriteString(me.Block.String())
		out.WriteString(" } ")
	} else {
		out.WriteString(me.Expr.String())
		out.WriteString(" , ")
	}

	out.WriteString(me.Value.String())
	return out.String()

}

///////////////////////////////////////////////////////////
//                         INFIX                         //
///////////////////////////////////////////////////////////
type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) Pos() token.Position {
	return ie.Token.Pos
}

func (ie *InfixExpression) End() token.Position {
	return ie.Right.End()
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         PREFIX                        //
///////////////////////////////////////////////////////////
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) Pos() token.Position {
	return pe.Token.Pos
}

func (pe *PrefixExpression) End() token.Position {
	return pe.Right.End()
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         POSTFIX                       //
///////////////////////////////////////////////////////////
type PostfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) Pos() token.Position {
	return pe.Token.Pos
}

func (pe *PostfixExpression) End() token.Position {
	ret := pe.Left.End()
	ret.Col = ret.Col + len(pe.Operator)
	return ret
}

func (pe *PostfixExpression) expressionNode() {}

func (pe *PostfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PostfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         TERNARY                         //
///////////////////////////////////////////////////////////
type TernaryExpression struct {
	Token     token.Token
	Condition Expression
	IfTrue    Expression
	IfFalse   Expression
}

func (te *TernaryExpression) Pos() token.Position {
	return te.Token.Pos
}

func (te *TernaryExpression) End() token.Position {
	return te.IfFalse.End()
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TernaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(te.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(te.IfTrue.String())
	out.WriteString(" : ")
	out.WriteString(te.IfFalse.String())
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                          CALL                         //
///////////////////////////////////////////////////////////
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
	Awaited   bool // if it is an awaited call
}

func (ce *CallExpression) Pos() token.Position {
	length := utf8.RuneCountInString(ce.Function.String())
	return token.Position{Filename: ce.Token.Pos.Filename, Line: ce.Token.Pos.Line, Col: ce.Token.Pos.Col - length}
}

func (ce *CallExpression) End() token.Position {
	aLen := len(ce.Arguments)
	if aLen > 0 {
		return ce.Arguments[aLen-1].End()
	}
	return ce.Function.End()
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     METHOD  CALL                      //
///////////////////////////////////////////////////////////
type MethodCallExpression struct {
	Token  token.Token
	Object Expression
	Call   Expression
}

func (mc *MethodCallExpression) Pos() token.Position {
	return mc.Token.Pos
}

func (mc *MethodCallExpression) End() token.Position {
	return mc.Call.End()
}

func (mc *MethodCallExpression) expressionNode()      {}
func (mc *MethodCallExpression) TokenLiteral() string { return mc.Token.Literal }
func (mc *MethodCallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(mc.Object.String())
	out.WriteString(".")
	out.WriteString(mc.Call.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                       CASE/ESLE                       //
///////////////////////////////////////////////////////////
type CaseExpr struct {
	Token        token.Token
	IsWholeMatch bool
	Expr         Expression
	Matches      []Expression
}

func (c *CaseExpr) Pos() token.Position {
	return c.Token.Pos
}

func (c *CaseExpr) End() token.Position {
	aLen := len(c.Matches)
	if aLen > 0 {
		return c.Matches[aLen-1].End()
	}
	return c.Expr.End()
}

func (c *CaseExpr) expressionNode()      {}
func (c *CaseExpr) TokenLiteral() string { return c.Token.Literal }

func (c *CaseExpr) String() string {
	var out bytes.Buffer

	out.WriteString("case ")
	out.WriteString(c.Expr.String())
	if c.IsWholeMatch {
		out.WriteString(" is ")
	} else {
		out.WriteString(" in ")
	}
	out.WriteString(" { ")

	matches := []string{}
	for _, m := range c.Matches {
		matches = append(matches, m.String())
	}

	out.WriteString(strings.Join(matches, " "))
	out.WriteString(" }")
	return out.String()
}

type CaseMatchExpr struct {
	Token token.Token
	Expr  Expression
	Block *BlockStatement
}

func (cm *CaseMatchExpr) Pos() token.Position {
	return cm.Token.Pos
}

func (cm *CaseMatchExpr) End() token.Position {
	return cm.Block.End()
}

func (cm *CaseMatchExpr) expressionNode()      {}
func (cm *CaseMatchExpr) TokenLiteral() string { return cm.Token.Literal }

func (cm *CaseMatchExpr) String() string {
	var out bytes.Buffer

	out.WriteString(cm.Expr.String())
	out.WriteString(" { ")
	out.WriteString(cm.Block.String())
	out.WriteString(" }")

	return out.String()
}

type CaseElseExpr struct {
	Token token.Token
	Block *BlockStatement
}

func (ce *CaseElseExpr) Pos() token.Position {
	return ce.Token.Pos
}

func (ce *CaseElseExpr) End() token.Position {
	return ce.Block.End()
}

func (ce *CaseElseExpr) expressionNode()      {}
func (ce *CaseElseExpr) TokenLiteral() string { return ce.Token.Literal }

func (ce *CaseElseExpr) String() string {
	var out bytes.Buffer

	out.WriteString("else ")
	out.WriteString(" { ")
	out.WriteString(ce.Block.String())
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                       SLICE/INDEX                     //
///////////////////////////////////////////////////////////
type SliceExpression struct {
	Token      token.Token
	StartIndex Expression
	EndIndex   Expression
}

func (se *SliceExpression) Pos() token.Position {
	return se.Token.Pos
}

func (se *SliceExpression) End() token.Position {
	if se.EndIndex != nil {
		return se.EndIndex.End()
	}

	return se.StartIndex.End()
}

func (se *SliceExpression) expressionNode()      {}
func (se *SliceExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SliceExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	if se.StartIndex != nil {
		out.WriteString(se.StartIndex.String())
	}
	out.WriteString(":")
	if se.EndIndex != nil {
		out.WriteString(se.EndIndex.String())
	}
	out.WriteString(")")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) Pos() token.Position {
	return ie.Token.Pos
}

func (ie *IndexExpression) End() token.Position {
	return ie.Index.End()
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")
	return out.String()
}

///////////////////////////////////////////////////////////
//                     SPAWN STATEMENT                   //
///////////////////////////////////////////////////////////
type SpawnStmt struct {
	Token token.Token
	Call  Expression
}

func (ss *SpawnStmt) Pos() token.Position {
	return ss.Token.Pos
}

func (ss *SpawnStmt) End() token.Position {
	return ss.Call.End()
}

func (ss *SpawnStmt) statementNode()       {}
func (ss *SpawnStmt) TokenLiteral() string { return ss.Token.Literal }

func (ss *SpawnStmt) String() string {
	var out bytes.Buffer

	out.WriteString(ss.TokenLiteral() + " ")
	out.WriteString(ss.Call.String())
	out.WriteString("; ")

	return out.String()
}

///////////////////////////////////////////////////////////
//                  PIPE OPERATOR                        //
///////////////////////////////////////////////////////////
// Pipe operator.
type Pipe struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (p *Pipe) Pos() token.Position {
	return p.Token.Pos
}

func (p *Pipe) End() token.Position {
	return p.Right.End()
}

func (p *Pipe) expressionNode()      {}
func (p *Pipe) TokenLiteral() string { return p.Token.Literal }
func (p *Pipe) String() string {
	var out bytes.Buffer

	out.WriteString(p.Left.String())
	out.WriteString(" |> ")
	out.WriteString(p.Right.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                   ENUM Literal                        //
///////////////////////////////////////////////////////////
type EnumLiteral struct {
	Token       token.Token
	Pairs       map[Expression]Expression
	RBraceToken token.Token
}

func (e *EnumLiteral) Pos() token.Position {
	return e.Token.Pos
}

//func (e *EnumLiteral) End() token.Position {
//	maxLineMap := make(map[int]Expression)
//
//	for _, value := range e.Pairs {
//		v := value.(Expression)
//		maxLineMap[v.End().Line] = v
//	}
//
//	maxLine := 0
//	for line, _ := range maxLineMap {
//		if line > maxLine {
//			maxLine = line
//		}
//	}
//
//	ret := maxLineMap[maxLine].(Expression)
//	return ret.End()
//}

func (e *EnumLiteral) End() token.Position {
	ret := e.RBraceToken.Pos
	ret.Col = ret.Col + 1
	return ret
}

func (e *EnumLiteral) expressionNode()      {}
func (e *EnumLiteral) TokenLiteral() string { return e.Token.Literal }

func (e *EnumLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("enum ")
	out.WriteString("{")

	pairs := []string{}
	for k, v := range e.Pairs {
		pairs = append(pairs, k.String()+" = "+v.String())
	}
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (e *EnumLiteral) Text(name string) string {
	var out bytes.Buffer

	out.WriteString("enum " + name)
	out.WriteString("{\n")

	pairs := []string{}
	for k, v := range e.Pairs {
		pairs = append(pairs, "\t"+k.String()+" = "+v.String())
	}
	out.WriteString(strings.Join(pairs, ", \n"))
	out.WriteString("\n}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                     ENUM STATEMENT                    //
///////////////////////////////////////////////////////////
type EnumStatement struct {
	Token       token.Token
	Name        *Identifier
	EnumLiteral *EnumLiteral

	//Doc related
	Doc         *CommentGroup // associated documentation; or nil
	SrcEndToken token.Token
}

func (e *EnumStatement) Pos() token.Position {
	return e.Token.Pos
}

func (e *EnumStatement) End() token.Position {
	return e.EnumLiteral.End()
}

//Below two methods implements 'Source' interface.
func (e *EnumStatement) SrcStart() token.Position {
	return e.Pos()
}

func (e *EnumStatement) SrcEnd() token.Position {
	ret := e.SrcEndToken.Pos
	length := utf8.RuneCountInString(e.SrcEndToken.Literal)
	ret.Offset += length
	return ret
}

func (e *EnumStatement) statementNode()       {}
func (e *EnumStatement) TokenLiteral() string { return e.Token.Literal }

func (e *EnumStatement) String() string {
	return e.EnumLiteral.Text(e.Name.Value)
}

func (e *EnumStatement) Docs() string {
	return e.String()
}

///////////////////////////////////////////////////////////
//        List Comprehension(for array & string)         //
///////////////////////////////////////////////////////////
// [ Expr for Var in Value <where Cond> ] ---> Value could be array or string
type ListComprehension struct {
	Token token.Token
	Var   string
	Value Expression //value(array or string) to range over
	Cond  Expression //conditional clause(nil if there is no 'WHERE' clause)
	Expr  Expression //the result expression
}

func (lc *ListComprehension) Pos() token.Position {
	return lc.Token.Pos
}

func (lc *ListComprehension) End() token.Position {
	if lc.Cond != nil {
		return lc.Cond.End()
	}
	return lc.Value.End()
}

func (lc *ListComprehension) expressionNode()      {}
func (lc *ListComprehension) TokenLiteral() string { return lc.Token.Literal }

func (lc *ListComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("[ ")
	out.WriteString(lc.Expr.String())
	out.WriteString(" for ")
	out.WriteString(lc.Var)
	out.WriteString(" in ")
	out.WriteString(lc.Value.String())
	if lc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(lc.Cond.String())
	}
	out.WriteString(" ]")

	return out.String()
}

///////////////////////////////////////////////////////////
//             List Comprehension(for range)             //
///////////////////////////////////////////////////////////
//[Expr for Var in StartIdx..EndIdx <where Cond>]
type ListRangeComprehension struct {
	Token    token.Token
	Var      string
	StartIdx Expression
	EndIdx   Expression
	Cond     Expression //conditional clause(nil if there is no 'WHERE' clause)
	Expr     Expression //the result expression
}

func (lc *ListRangeComprehension) Pos() token.Position {
	return lc.Token.Pos
}

func (lc *ListRangeComprehension) End() token.Position {
	if lc.Cond != nil {
		return lc.Cond.End()
	}
	return lc.EndIdx.End()
}

func (lc *ListRangeComprehension) expressionNode()      {}
func (lc *ListRangeComprehension) TokenLiteral() string { return lc.Token.Literal }

func (lc *ListRangeComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("[ ")
	out.WriteString(lc.Expr.String())
	out.WriteString(" for ")
	out.WriteString(lc.Var)
	out.WriteString(" in ")
	out.WriteString(lc.StartIdx.String())
	out.WriteString("..")
	out.WriteString(lc.EndIdx.String())
	if lc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(lc.Cond.String())
	}
	out.WriteString(" ]")

	return out.String()
}

///////////////////////////////////////////////////////////
//                LIST Map Comprehension                 //
///////////////////////////////////////////////////////////
//[ Expr for Key,Value in X <where Cond>]
type ListMapComprehension struct {
	Token token.Token
	Key   string
	Value string
	X     Expression //value(hash) to range over
	Cond  Expression //Conditional clause(nil if there is no 'WHERE' clause)
	Expr  Expression //the result expression
}

func (mc *ListMapComprehension) Pos() token.Position {
	return mc.Token.Pos
}

func (mc *ListMapComprehension) End() token.Position {
	if mc.Cond != nil {
		return mc.Cond.End()
	}
	return mc.Expr.End()
}

func (mc *ListMapComprehension) expressionNode()      {}
func (mc *ListMapComprehension) TokenLiteral() string { return mc.Token.Literal }

func (mc *ListMapComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("[ ")
	out.WriteString(mc.Expr.String())
	out.WriteString(" for ")
	out.WriteString(mc.Key + ", " + mc.Value)
	out.WriteString(" in ")
	out.WriteString(mc.X.String())
	if mc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(mc.Cond.String())
	}
	out.WriteString(" ]")

	return out.String()
}

///////////////////////////////////////////////////////////
//        Hash Comprehension(for array & string)         //
///////////////////////////////////////////////////////////
//{ KeyExpr:ValExpr for Var in Value <where Cond> }  -->Value could be array or string
type HashComprehension struct {
	Token   token.Token
	Var     string
	Value   Expression //value(array or string) to range over
	Cond    Expression //conditional clause(nil if there is no 'WHERE' clause)
	KeyExpr Expression //the result Key expression
	ValExpr Expression //the result Value expression
}

func (hc *HashComprehension) Pos() token.Position {
	return hc.Token.Pos
}

func (hc *HashComprehension) End() token.Position {
	if hc.Cond != nil {
		return hc.Cond.End()
	}
	return hc.Value.End()
}

func (hc *HashComprehension) expressionNode()      {}
func (hc *HashComprehension) TokenLiteral() string { return hc.Token.Literal }

func (hc *HashComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	out.WriteString(hc.KeyExpr.String())
	out.WriteString(" : ")
	out.WriteString(hc.ValExpr.String())
	out.WriteString(" for ")
	out.WriteString(hc.Var)
	out.WriteString(" in ")
	out.WriteString(hc.Value.String())
	if hc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(hc.Cond.String())
	}
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//             Hash Comprehension(for range)             //
///////////////////////////////////////////////////////////
//{ KeyExp:ValExp for Var in StartIdx..EndIdx <where Cond> }
type HashRangeComprehension struct {
	Token    token.Token
	Var      string
	StartIdx Expression
	EndIdx   Expression
	Cond     Expression //conditional clause(nil if there is no 'WHERE' clause)
	KeyExpr  Expression //the result Key expression
	ValExpr  Expression //the result Value expression
}

func (hc *HashRangeComprehension) Pos() token.Position {
	return hc.Token.Pos
}

func (hc *HashRangeComprehension) End() token.Position {
	if hc.Cond != nil {
		return hc.Cond.End()
	}
	return hc.EndIdx.End()
}

func (hc *HashRangeComprehension) expressionNode()      {}
func (hc *HashRangeComprehension) TokenLiteral() string { return hc.Token.Literal }

func (hc *HashRangeComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	out.WriteString(hc.KeyExpr.String())
	out.WriteString(" : ")
	out.WriteString(hc.ValExpr.String())
	out.WriteString(" for ")
	out.WriteString(hc.Var)
	out.WriteString(" in ")
	out.WriteString(hc.StartIdx.String())
	out.WriteString("..")
	out.WriteString(hc.EndIdx.String())
	if hc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(hc.Cond.String())
	}
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                Hash Map Comprehension                 //
///////////////////////////////////////////////////////////
//{ KeyExpr:ValExpr for Key,Value in X <where Cond> }
type HashMapComprehension struct {
	Token   token.Token
	Key     string
	Value   string
	X       Expression //value(hash) to range over
	Cond    Expression //Conditional clause(nil if there is no 'WHERE' clause)
	KeyExpr Expression //the result Key expression
	ValExpr Expression //the result Value expression
}

func (mc *HashMapComprehension) Pos() token.Position {
	return mc.Token.Pos
}

func (mc *HashMapComprehension) End() token.Position {
	if mc.Cond != nil {
		return mc.Cond.End()
	}
	return mc.X.End()
}

func (mc *HashMapComprehension) expressionNode()      {}
func (mc *HashMapComprehension) TokenLiteral() string { return mc.Token.Literal }

func (mc *HashMapComprehension) String() string {
	var out bytes.Buffer

	out.WriteString("{ ")
	out.WriteString(mc.KeyExpr.String())
	out.WriteString(" : ")
	out.WriteString(mc.ValExpr.String())
	out.WriteString(" for ")
	out.WriteString(mc.Key + ", " + mc.Value)
	out.WriteString(" in ")
	out.WriteString(mc.X.String())
	if mc.Cond != nil {
		out.WriteString(" where ")
		out.WriteString(mc.Cond.String())
	}
	out.WriteString(" }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                      Tuple LITERAL                    //
///////////////////////////////////////////////////////////
type TupleLiteral struct {
	Token       token.Token
	Members     []Expression
	RParenToken token.Token
}

func (t *TupleLiteral) Pos() token.Position {
	return t.Token.Pos
}

//func (t *TupleLiteral) End() token.Position {
//	aLen := len(t.Members)
//	if aLen > 0 {
//		return t.Members[aLen-1].End()
//	}
//	return t.Token.Pos
//}

func (t *TupleLiteral) End() token.Position {
	return t.RParenToken.Pos
}

func (t *TupleLiteral) expressionNode()      {}
func (t *TupleLiteral) TokenLiteral() string { return t.Token.Literal }
func (t *TupleLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("(")

	members := []string{}
	for _, m := range t.Members {
		members = append(members, m.String())
	}

	out.WriteString(strings.Join(members, ", "))
	out.WriteString(")")

	return out.String()
}

//class's method modifier
type ModifierLevel int

const (
	ModifierDefault ModifierLevel = iota
	ModifierPrivate
	ModifierProtected
	ModifierPublic
)

//for debug purpose
func (m ModifierLevel) String() string {
	switch m {
	case ModifierPrivate:
		return "private "
	case ModifierProtected:
		return "protected "
	case ModifierPublic:
		return "public "
	}

	return ""
}

///////////////////////////////////////////////////////////
//                      CLASS LITERAL                    //
///////////////////////////////////////////////////////////
// class : parentClass { block }
type ClassLiteral struct {
	Token      token.Token
	Name       string
	Parent     string
	Members    []*LetStatement               //class's fields
	Properties map[string]*PropertyDeclStmt  //class's properties
	Methods    map[string]*FunctionStatement //class's methods
	Block      *BlockStatement               //mainly used for debugging purpose
	Modifier   ModifierLevel                 //NOT IMPLEMENTED
}

func (c *ClassLiteral) Pos() token.Position {
	return c.Token.Pos
}

func (c *ClassLiteral) End() token.Position {
	return c.Block.End()
}

func (c *ClassLiteral) expressionNode()      {}
func (c *ClassLiteral) TokenLiteral() string { return c.Token.Literal }

func (c *ClassLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(c.TokenLiteral() + " ")
	out.WriteString(c.Name)
	if len(c.Parent) != 0 {
		out.WriteString(" : " + c.Parent + " ")
	}

	out.WriteString("{ ")
	out.WriteString(c.Block.String())
	out.WriteString("} ")

	return out.String()
}

//class classname : parentClass { block }
//class @classname: parentClass { block } //Annotation
///////////////////////////////////////////////////////////
//                     CLASS STATEMENT                   //
///////////////////////////////////////////////////////////
type ClassStatement struct {
	Token        token.Token
	Name         *Identifier //Class name
	CategoryName *Identifier
	ClassLiteral *ClassLiteral
	IsAnnotation bool //class is a annotation class

	//Doc related
	Doc         *CommentGroup // associated documentation; or nil
	SrcEndToken token.Token
}

func (c *ClassStatement) Pos() token.Position {
	return c.Token.Pos
}

func (c *ClassStatement) End() token.Position {
	return c.ClassLiteral.Block.End()
}

//Below two methods implements 'Source' interface.
func (c *ClassStatement) SrcStart() token.Position {
	return c.Pos()
}

func (c *ClassStatement) SrcEnd() token.Position {
	ret := c.SrcEndToken.Pos
	length := utf8.RuneCountInString(c.SrcEndToken.Literal)
	ret.Offset += length
	return ret
}

func (c *ClassStatement) statementNode()       {}
func (c *ClassStatement) TokenLiteral() string { return c.Token.Literal }
func (c *ClassStatement) String() string {
	var out bytes.Buffer

	out.WriteString(c.Token.Literal + " ")
	if c.IsAnnotation {
		out.WriteString("@")
	}
	out.WriteString(c.Name.String())

	if c.CategoryName != nil {
		out.WriteString("(")
		out.WriteString(c.CategoryName.String())
		out.WriteString(") ")
	} else {
		if len(c.ClassLiteral.Parent) > 0 {
			out.WriteString(" : " + c.ClassLiteral.Parent)
		}
	}

	out.WriteString("{ ")
	out.WriteString(c.ClassLiteral.Block.String())
	out.WriteString(" }")

	return out.String()
}

func (c *ClassStatement) Docs() string {
	var out bytes.Buffer

	out.WriteString(c.Token.Literal + " ")
	if c.IsAnnotation {
		out.WriteString("@")
	}
	out.WriteString(c.Name.String())

	if c.CategoryName != nil {
		out.WriteString("(")
		out.WriteString(c.CategoryName.String())
		out.WriteString(") ")
	} else {
		if len(c.ClassLiteral.Parent) > 0 {
			out.WriteString(" : " + c.ClassLiteral.Parent)
		}
	}

	out.WriteString("{ ... }")
	return out.String()
}

///////////////////////////////////////////////////////////
//                   NEW EXPRESSION                      //
///////////////////////////////////////////////////////////

type NewExpression struct {
	Token     token.Token
	Class     Expression
	Arguments []Expression
}

func (c *NewExpression) Pos() token.Position {
	return c.Token.Pos
}

func (c *NewExpression) End() token.Position {
	aLen := len(c.Arguments)
	if aLen > 0 {
		return c.Arguments[aLen-1].End()
	}
	return c.Class.End()
}

func (n *NewExpression) expressionNode()      {}
func (n *NewExpression) TokenLiteral() string { return n.Token.Literal }
func (n *NewExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range n.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(n.TokenLiteral() + " ")
	out.WriteString(n.Class.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(") ")

	return out.String()
}

//class's property declaration
type PropertyDeclStmt struct {
	Token         token.Token
	Name          *Identifier   //property name
	Getter        *GetterStmt   //getter
	Setter        *SetterStmt   //setter
	Indexes       []*Identifier //only used in class's indexer
	StaticFlag    bool
	ModifierLevel ModifierLevel //property's modifier
	Annotations   []*AnnotationStmt
	Default       Expression

	//Doc related
	Doc         *CommentGroup // associated documentation; or nil
	SrcEndToken token.Token
}

func (p *PropertyDeclStmt) Pos() token.Position {
	return p.Token.Pos
}

func (p *PropertyDeclStmt) End() token.Position {
	if p.Getter == nil {
		return p.Setter.End()
	}
	return p.Getter.End()
}

//Below two methods implements 'Source' interface.
func (p *PropertyDeclStmt) SrcStart() token.Position {
	return p.Pos()
}

func (p *PropertyDeclStmt) SrcEnd() token.Position {
	ret := p.SrcEndToken.Pos
	length := utf8.RuneCountInString(p.SrcEndToken.Literal)
	ret.Offset += length
	return ret
}

func (p *PropertyDeclStmt) statementNode()       {}
func (p *PropertyDeclStmt) TokenLiteral() string { return p.Token.Literal }

func (p *PropertyDeclStmt) String() string {
	var out bytes.Buffer

	out.WriteString(p.ModifierLevel.String())
	if p.StaticFlag {
		out.WriteString("static ")
	}

	out.WriteString("property ")
	if strings.HasPrefix(p.Name.String(), "this") {
		out.WriteString("this")
	} else {
		out.WriteString(p.Name.String())
	}

	if p.Indexes != nil {
		parameters := []string{}
		for _, idx := range p.Indexes {
			parameters = append(parameters, idx.String())
		}

		out.WriteString("[")
		out.WriteString(strings.Join(parameters, ", "))
		out.WriteString("]")
	} else {
	}

	if p.Default != nil { //must be an annotation class
		out.WriteString(" default ")
		out.WriteString(p.Default.String())
		return out.String()
	}

	out.WriteString(" { ")

	if p.Getter != nil {
		out.WriteString(p.Getter.String())
	}

	if p.Setter != nil {
		out.WriteString(p.Setter.String())
	}

	out.WriteString("} ")
	return out.String()
}

func (p *PropertyDeclStmt) Docs() string {
	var out bytes.Buffer

	out.WriteString(p.ModifierLevel.String())
	if p.StaticFlag {
		out.WriteString("static ")
	}

	out.WriteString("property ")
	if strings.HasPrefix(p.Name.String(), "this") {
		out.WriteString("this")
	} else {
		out.WriteString(p.Name.String())
	}

	if p.Indexes != nil {
		parameters := []string{}
		for _, idx := range p.Indexes {
			parameters = append(parameters, idx.String())
		}

		out.WriteString("[")
		out.WriteString(strings.Join(parameters, ", "))
		out.WriteString("]")
	} else {
	}

	if p.Default != nil { //must be an annotation class
		out.WriteString(" default ")
		out.WriteString(p.Default.String())
		return out.String()
	}

	out.WriteString(" { ")

	if p.Getter != nil {
		out.WriteString("get; ")
	}

	if p.Setter != nil {
		out.WriteString("set; ")
	}

	out.WriteString("} ")
	return out.String()
}

//property's getter statement
type GetterStmt struct {
	Token token.Token
	Body  *BlockStatement
}

func (g *GetterStmt) Pos() token.Position {
	return g.Token.Pos
}

func (g *GetterStmt) End() token.Position {
	return g.Body.End()
}

func (g *GetterStmt) statementNode()       {}
func (g *GetterStmt) TokenLiteral() string { return g.Token.Literal }

func (g *GetterStmt) String() string {
	var out bytes.Buffer

	out.WriteString("get")
	if len(g.Body.Statements) == 0 {
		out.WriteString("; ")
	} else {
		out.WriteString("{")
		out.WriteString(g.Body.String())
		out.WriteString("} ")
	}

	return out.String()
}

//property's setter statement
//setter variable is always 'value' like c#
type SetterStmt struct {
	Token token.Token
	Body  *BlockStatement
}

func (s *SetterStmt) Pos() token.Position {
	return s.Token.Pos
}

func (s *SetterStmt) End() token.Position {
	return s.Body.End()
}

func (s *SetterStmt) statementNode()       {}
func (s *SetterStmt) TokenLiteral() string { return s.Token.Literal }

func (s *SetterStmt) String() string {
	var out bytes.Buffer

	out.WriteString("set")
	if len(s.Body.Statements) == 0 {
		out.WriteString("; ")
	} else {
		out.WriteString("{")
		out.WriteString(s.Body.String())
		out.WriteString("} ")
	}

	return out.String()
}

///////////////////////////////////////////////////////////
//                     CLASS/INDEXER                     //
///////////////////////////////////////////////////////////
type ClassIndexerExpression struct {
	Token      token.Token
	Parameters []Expression //indexer's parameters
}

func (ci *ClassIndexerExpression) Pos() token.Position {
	return ci.Token.Pos
}

func (ci *ClassIndexerExpression) End() token.Position {
	aLen := len(ci.Parameters)
	if aLen > 0 {
		return ci.Parameters[aLen-1].End()
	}
	return ci.Token.Pos
}

func (ci *ClassIndexerExpression) expressionNode()      {}
func (ci *ClassIndexerExpression) TokenLiteral() string { return ci.Token.Literal }
func (ci *ClassIndexerExpression) String() string {
	var out bytes.Buffer

	parameters := []string{}
	for _, p := range ci.Parameters {
		parameters = append(parameters, p.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString("]")

	return out.String()
}

///////////////////////////////////////////////////////////
//                      ANNOTATIONS                      //
///////////////////////////////////////////////////////////
type AnnotationStmt struct {
	Token      token.Token
	Name       *Identifier
	Attributes map[string]Expression
}

func (anno *AnnotationStmt) Pos() token.Position {
	return anno.Token.Pos
}

func (anno *AnnotationStmt) End() token.Position {
	return anno.Token.Pos
}

func (anno *AnnotationStmt) statementNode()       {}
func (anno *AnnotationStmt) TokenLiteral() string { return anno.Token.Literal }
func (anno *AnnotationStmt) String() string {
	var out bytes.Buffer

	Attrs := []string{}
	for _, attr := range anno.Attributes {
		Attrs = append(Attrs, attr.String())
	}

	out.WriteString("@")
	out.WriteString(anno.Name.String())
	out.WriteString("{")
	out.WriteString(strings.Join(Attrs, ", "))
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         USING                         //
///////////////////////////////////////////////////////////
type UsingStmt struct {
	Token token.Token
	Expr  *AssignExpression
	Block *BlockStatement
}

func (u *UsingStmt) Pos() token.Position {
	return u.Token.Pos
}

func (u *UsingStmt) End() token.Position {
	return u.Block.End()
}

func (u *UsingStmt) statementNode()       {}
func (u *UsingStmt) TokenLiteral() string { return u.Token.Literal }
func (u *UsingStmt) String() string {
	var out bytes.Buffer

	out.WriteString("using (")
	out.WriteString(u.Expr.String())
	out.WriteString(") ")
	out.WriteString("{")
	out.WriteString(u.Block.String())
	out.WriteString("}")

	return out.String()
}

///////////////////////////////////////////////////////////
//                        COMMAND                        //
///////////////////////////////////////////////////////////
type CmdExpression struct {
	Token token.Token
	Value string
}

func (c *CmdExpression) Pos() token.Position {
	return c.Token.Pos
}

func (c *CmdExpression) End() token.Position {
	length := utf8.RuneCountInString(c.Value)
	return token.Position{Filename: c.Token.Pos.Filename, Line: c.Token.Pos.Line, Col: c.Token.Pos.Col + length}
}

func (c *CmdExpression) expressionNode()      {}
func (c *CmdExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CmdExpression) String() string       { return c.Value }

///////////////////////////////////////////////////////////
//                     LINQ QUERY                        //
///////////////////////////////////////////////////////////

/* Syntax:(From Antlr)
query_expression : from_clause query_body
from_clause : FROM identifier IN expression
query_body : query_body_clause* select_or_group_clause query_continuation?
query_body_clause: from_clause | let_clause | where_clause | combined_join_clause | orderby_clause
where_clause : WHERE expression
combined_join_clause : JOIN identifier IN expression ON expression EQUALS expression (INTO identifier)?
orderby_clause : ORDERBY ordering (','  ordering)*
ordering : expression (ASCENDING | DESCENDING)?
select_or_group_clause : SELECT expression | GROUP expression BY expression
query_continuation : INTO identifier query_body
*/

//query_expression : from_clause query_body
type QueryExpr struct {
	Token     token.Token //'from'
	From      Expression  //FromExpr
	QueryBody Expression  //QueryBodyExpr
}

func (q *QueryExpr) Pos() token.Position {
	return q.Token.Pos
}

func (q *QueryExpr) End() token.Position {
	return q.QueryBody.End()
}

func (q *QueryExpr) expressionNode()      {}
func (q *QueryExpr) TokenLiteral() string { return q.Token.Literal }
func (q *QueryExpr) String() string {
	var out bytes.Buffer

	out.WriteString(q.From.String())
	out.WriteString(" ")
	out.WriteString(q.QueryBody.String())

	s := out.String()
	return strings.Join(strings.Fields(s), " ") // remove extra spaces
}

//from_clause : FROM identifier IN expression
type FromExpr struct {
	Token token.Token //from
	Var   string      //identifier
	Expr  Expression
}

func (f *FromExpr) Pos() token.Position {
	return f.Token.Pos
}

func (f *FromExpr) End() token.Position {
	return f.Expr.End()
}

func (f *FromExpr) expressionNode()      {}
func (f *FromExpr) TokenLiteral() string { return f.Token.Literal }
func (f *FromExpr) String() string {
	var out bytes.Buffer

	out.WriteString("from ")
	out.WriteString(f.Var)
	out.WriteString(" in ")
	out.WriteString(f.Expr.String())

	return out.String()
}

//query_body : query_body_clause* select_or_group_clause query_continuation?
type QueryBodyExpr struct {
	QueryBody         []Expression //QueryBodyClauseExpr
	Expr              Expression   //SelectExpr or GroupExpr
	QueryContinuation Expression   //QueryContinuationExpr
}

func (q *QueryBodyExpr) Pos() token.Position {
	if len(q.QueryBody) == 0 {
		return q.Expr.Pos()
	}
	return q.QueryBody[0].Pos()
}

func (q *QueryBodyExpr) End() token.Position {
	if q.QueryContinuation == nil {
		return q.Expr.End()
	}
	return q.QueryContinuation.End()
}

func (q *QueryBodyExpr) expressionNode()      {}
func (q *QueryBodyExpr) TokenLiteral() string { return "query_body_expr" }
func (q *QueryBodyExpr) String() string {
	var out bytes.Buffer

	queryBody := []string{}
	for _, qb := range q.QueryBody {
		queryBody = append(queryBody, qb.String())
	}
	out.WriteString(strings.Join(queryBody, " "))

	out.WriteString(q.Expr.String())

	if q.QueryContinuation != nil {
		out.WriteString(q.QueryContinuation.String())
	}

	return out.String()
}

//query_body_clause: from_clause | let_clause | where_clause | combined_join_clause | orderby_clause
type QueryBodyClauseExpr struct {
	Expr Expression
}

func (q *QueryBodyClauseExpr) Pos() token.Position {
	return q.Expr.Pos()
}

func (q *QueryBodyClauseExpr) End() token.Position {
	return q.Expr.End()
}

func (q *QueryBodyClauseExpr) expressionNode()      {}
func (q *QueryBodyClauseExpr) TokenLiteral() string { return "query_body_clause_expr" }
func (q *QueryBodyClauseExpr) String() string {
	var out bytes.Buffer

	out.WriteString(q.Expr.String())
	return out.String()
}

//where_clause : WHERE expression
type WhereExpr struct {
	Token token.Token //'where'
	Expr  Expression
}

func (w *WhereExpr) Pos() token.Position {
	return w.Token.Pos
}

func (w *WhereExpr) End() token.Position {
	return w.Expr.End()
}

func (w *WhereExpr) expressionNode()      {}
func (w *WhereExpr) TokenLiteral() string { return w.Token.Literal }
func (w *WhereExpr) String() string {
	var out bytes.Buffer

	out.WriteString(" where ")
	out.WriteString(w.Expr.String())

	return out.String()
}

//combined_join_clause : JOIN identifier IN expression ON expression EQUALS expression (INTO identifier)?
type JoinExpr struct {
	Token     token.Token //'join'
	JoinVar   string      //identifier
	InExpr    Expression
	OnExpr    Expression
	EqualExpr Expression
	IntoVar   *Identifier //why IntoVar's type is '*Identifier', not 'string'? because we need it in 'End()' function.
}

func (j *JoinExpr) Pos() token.Position {
	return j.Token.Pos
}

func (j *JoinExpr) End() token.Position {
	if j.IntoVar == nil {
		return j.EqualExpr.End()
	}
	return j.IntoVar.End()
}

func (j *JoinExpr) expressionNode()      {}
func (j *JoinExpr) TokenLiteral() string { return j.Token.Literal }
func (j *JoinExpr) String() string {
	var out bytes.Buffer
	out.WriteString(" join ")
	out.WriteString(j.JoinVar)
	out.WriteString(" in ")
	out.WriteString(j.InExpr.String())
	out.WriteString(" on ")
	out.WriteString(j.OnExpr.String())
	out.WriteString(" equals ")
	out.WriteString(j.EqualExpr.String())

	if j.IntoVar != nil {
		out.WriteString(" into ")
		out.WriteString(j.IntoVar.String())
	}

	return out.String()
}

//orderby_clause : ORDERBY ordering (','  ordering)*
type OrderExpr struct {
	Token    token.Token  //'orderby'
	Ordering []Expression //[]*OrderingExpr
}

func (o *OrderExpr) Pos() token.Position {
	return o.Token.Pos
}

func (o *OrderExpr) End() token.Position {
	return o.Ordering[len(o.Ordering)-1].End()
}

func (o *OrderExpr) expressionNode()      {}
func (o *OrderExpr) TokenLiteral() string { return o.Token.Literal }
func (o *OrderExpr) String() string {
	var out bytes.Buffer

	out.WriteString(" orderby ")
	ordering := []string{}
	for _, order := range o.Ordering {
		ordering = append(ordering, order.String())
	}
	out.WriteString(strings.Join(ordering, ", "))

	return out.String()
}

//ordering : expression (ASCENDING | DESCENDING)?
type OrderingExpr struct {
	Expr         Expression
	IsAscending  bool // if there is no 'ASCENDING or 'DESCENDING', it's default to 'ASCENDING'
	HasSortOrder bool
	OrderToken   token.Token //'ascending' or 'descending'
	Var          string
}

func (o *OrderingExpr) Pos() token.Position {
	return o.Expr.Pos()
}

func (o *OrderingExpr) End() token.Position {
	if o.HasSortOrder {
		length := utf8.RuneCountInString(o.OrderToken.Literal)
		return token.Position{Filename: o.OrderToken.Pos.Filename, Line: o.OrderToken.Pos.Line, Col: o.OrderToken.Pos.Col + length}
	}

	return o.Expr.End()
}

func (o *OrderingExpr) expressionNode()      {}
func (o *OrderingExpr) TokenLiteral() string { return "ordering_expr" }
func (o *OrderingExpr) String() string {
	var out bytes.Buffer

	out.WriteString(o.Expr.String())
	if o.HasSortOrder {
		out.WriteString(" ")
		out.WriteString(o.OrderToken.Literal)
	}
	return out.String()
}

//select_or_group_clause : SELECT expression | GROUP expression BY expression
//SELECT expression
type SelectExpr struct {
	Token token.Token //'select'
	Expr  Expression
}

func (s *SelectExpr) Pos() token.Position {
	return s.Token.Pos
}

func (s *SelectExpr) End() token.Position {
	return s.Expr.End()
}

func (s *SelectExpr) expressionNode()      {}
func (s *SelectExpr) TokenLiteral() string { return s.Token.Literal }
func (s *SelectExpr) String() string {
	var out bytes.Buffer

	out.WriteString(" select ")
	out.WriteString(s.Expr.String())

	return out.String()
}

//GROUP expression BY expression
type GroupExpr struct {
	Token   token.Token //'group'
	GrpExpr Expression
	ByExpr  Expression
}

func (g *GroupExpr) Pos() token.Position {
	return g.Token.Pos
}

func (g *GroupExpr) End() token.Position {
	return g.ByExpr.End()
}

func (g *GroupExpr) expressionNode()      {}
func (g *GroupExpr) TokenLiteral() string { return g.Token.Literal }
func (g *GroupExpr) String() string {
	var out bytes.Buffer

	out.WriteString(" group ")
	out.WriteString(g.GrpExpr.String())
	out.WriteString(" by ")
	out.WriteString(g.ByExpr.String())

	return out.String()
}

//query_continuation : INTO identifier query_body
type QueryContinuationExpr struct {
	Token token.Token // 'into'
	Var   string
	Expr  Expression //QueryBodyExpr
}

func (q *QueryContinuationExpr) Pos() token.Position {
	return q.Token.Pos
}

func (q *QueryContinuationExpr) End() token.Position {
	return q.Expr.End()
}

func (q *QueryContinuationExpr) expressionNode()      {}
func (q *QueryContinuationExpr) TokenLiteral() string { return q.Token.Literal }
func (q *QueryContinuationExpr) String() string {
	var out bytes.Buffer

	out.WriteString(" into ")
	out.WriteString(q.Var)
	out.WriteString(q.Expr.String())

	return out.String()
}

///////////////////////////////////////////////////////////
//                     AWAIT EXPRESSION                  //
///////////////////////////////////////////////////////////

//result = await add(1, 2)
//await hello()
type AwaitExpr struct {
	Token token.Token
	Call  Expression
}

func (aw *AwaitExpr) Pos() token.Position {
	return aw.Token.Pos
}

func (aw *AwaitExpr) End() token.Position {
	return aw.Call.End()
}

func (aw *AwaitExpr) expressionNode()      {}
func (aw *AwaitExpr) TokenLiteral() string { return aw.Token.Literal }

func (aw *AwaitExpr) String() string {
	var out bytes.Buffer

	out.WriteString(aw.TokenLiteral() + " ")
	out.WriteString(aw.Call.String())
	out.WriteString("; ")

	return out.String()
}

//service servicename on addrs { block }
///////////////////////////////////////////////////////////
//                     Service STATEMENT                 //
///////////////////////////////////////////////////////////
type ServiceStatement struct {
	Token   token.Token
	Name    *Identifier //Service name
	Addr    string
	Debug   bool
	Methods map[string]*FunctionStatement //service's methods
	Block   *BlockStatement               //mainly used for debugging purpose

	//Doc related
	Doc         *CommentGroup // associated documentation; or nil
	SrcEndToken token.Token
}

func (s *ServiceStatement) Pos() token.Position {
	return s.Token.Pos
}

func (s *ServiceStatement) End() token.Position {
	return s.Block.End()
}

//Below two methods implements 'Source' interface.
func (s *ServiceStatement) SrcStart() token.Position {
	return s.Pos()
}

func (s *ServiceStatement) SrcEnd() token.Position {
	ret := s.SrcEndToken.Pos
	length := utf8.RuneCountInString(s.SrcEndToken.Literal)
	ret.Offset += length
	return ret
}

func (s *ServiceStatement) statementNode()       {}
func (s *ServiceStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ServiceStatement) String() string {
	var out bytes.Buffer

	out.WriteString(s.Token.Literal + " ")
	out.WriteString(s.Name.String())
	out.WriteString(" on '")
	out.WriteString(s.Addr)
	out.WriteString("' { ")
	out.WriteString(s.Block.String())
	out.WriteString(" }")

	return out.String()
}

func (s *ServiceStatement) Docs() string {
	var out bytes.Buffer

	out.WriteString(s.Token.Literal + " ")
	out.WriteString(s.Name.String())
	out.WriteString(" on '")
	out.WriteString(s.Addr)
	out.WriteString("'{ ... }")

	return out.String()
}

///////////////////////////////////////////////////////////
//                   DateTime Expression                 //
///////////////////////////////////////////////////////////
type DateTimeExpr struct {
	Token   token.Token
	Pattern *InterpolatedString // pattern string
}

func (dt *DateTimeExpr) Pos() token.Position {
	return dt.Token.Pos
}

func (dt *DateTimeExpr) End() token.Position {
	return dt.Pattern.End()
	// length := len(dt.Pattern)
	// pos := dt.Token.Pos
	// return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (dt *DateTimeExpr) expressionNode()      {}
func (dt *DateTimeExpr) TokenLiteral() string { return dt.Token.Literal }

func (dt *DateTimeExpr) String() string {
	var out bytes.Buffer

	out.WriteString(dt.TokenLiteral() + "/")
	out.WriteString(dt.Pattern.String())
	out.WriteString("/")

	return out.String()
}

///////////////////////////////////////////////////////////
//                         Diamond                       //
///////////////////////////////////////////////////////////
type DiamondExpr struct {
	Token token.Token
	Value string
}

func (d *DiamondExpr) Pos() token.Position {
	return d.Token.Pos
}

func (d *DiamondExpr) End() token.Position {
	length := utf8.RuneCountInString(d.Value)
	return token.Position{Filename: d.Token.Pos.Filename, Line: d.Token.Pos.Line, Col: d.Token.Pos.Col + length + 1}
}

func (d *DiamondExpr) expressionNode()      {}
func (d *DiamondExpr) TokenLiteral() string { return d.Token.Literal }
func (d *DiamondExpr) String() string       { return "<$" + d.Value + ">" }

///////////////////////////////////////////////////////////
//                       COMMENTS                        //
///////////////////////////////////////////////////////////

// A Comment node represents a single //-style or /*-style comment.
type Comment struct {
	Token token.Token
	Text  string // comment text
}

func (c *Comment) Pos() token.Position { return c.Token.Pos }
func (c *Comment) End() token.Position {
	tokLen := utf8.RuneCountInString(c.Token.Literal)
	textLen := utf8.RuneCountInString(c.Text)
	pos := c.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + tokLen + textLen - 1}
}

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
//
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() token.Position { return g.List[0].Pos() }
func (g *CommentGroup) End() token.Position { return g.List[len(g.List)-1].End() }

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func stripTrailingWhitespace(s string) string {
	i := len(s)
	for i > 0 && isWhitespace(s[i-1]) {
		i--
	}
	return s[0:i]
}

// Text returns the text of the comment.
// Comment markers (//, /*, and */), the first space of a line comment, and
// leading and trailing empty lines are removed. Multiple empty lines are
// reduced to one, and trailing space on lines is trimmed. Unless the result
// is empty, it is newline-terminated.
//
func (g *CommentGroup) Text() string {
	if g == nil {
		return ""
	}
	comments := make([]string, len(g.List))
	for i, c := range g.List {
		comments[i] = c.Text
	}

	lines := make([]string, 0, 10) // most comments are less than 10 lines
	for _, c := range comments {
		// Remove comment markers.
		// The parser has given us exactly the comment text.
		if c[0] == '#' {
			c = c[1:]
		} else {
			switch c[1] {
			case '/':
				//-style comment (no newline at the end)
				c = c[2:]
				// strip first space - required for Example tests
				if len(c) > 0 && c[0] == ' ' {
					c = c[1:]
				}
			case '*':
				/*-style comment */
				c = c[2 : len(c)-2]
			}
		}

		// Split on newlines.
		cl := strings.Split(c, "\n")

		// Walk lines, stripping trailing white space and adding to list, also strip the line which begins with '*'
		for _, l := range cl {
			tmpline := l

			if len(tmpline) != 0 {
				tmpline = strings.TrimLeftFunc(l, func(r rune) bool {
					return r == ' ' || r == '\t'
				})

				if len(tmpline) > 0 && tmpline[0] == '*' {
					// strip first '*'
					tmpline = tmpline[1:]
				}
				if len(tmpline) > 0 && tmpline[0] == ' ' {
					//strip first space
					tmpline = tmpline[1:]
				}
			}
			lines = append(lines, stripTrailingWhitespace(tmpline))
		}
	}

	// Remove leading blank lines; convert runs of
	// interior blank lines to a single blank line.
	n := 0
	for _, line := range lines {
		if line != "" || n > 0 && lines[n-1] != "" {
			lines[n] = line
			n++
		}
	}
	lines = lines[0:n]

	// Add final "" entry to get trailing newline from Join.
	if n > 0 && lines[n-1] != "" {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}
