package eval

import (
	"magpie/ast"
)

const (
	OPTIONAL_OBJ  = "OPTIONAL_OBJ"
	optional_name = "optional"
)

type Optional struct {
	Value Object
}

func NewOptionalObj() Object {
	ret := &Optional{Value: NIL}
	SetGlobalObj(optional_name, ret)

	return ret
}

func (o *Optional) Inspect() string {
	if o.Value != NIL {
		return "Optional[" + o.Value.Inspect() + "]"
	}
	return "Optional.empty"
}

func (o *Optional) Type() ObjectType { return OPTIONAL_OBJ }

func (o *Optional) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "empty":
		return o.Empty(line, args...)
	case "of":
		return o.Of(line, args...)
	case "ofNullable":
		return o.OfNullable(line, args...)
	case "get":
		return o.Get(line, args...)
	case "isPresent":
		return o.IsPresent(line, args...)
	case "or":
		return o.Or(line, scope, args...)
	case "orElse":
		return o.OrElse(line, args...)
	case "orElseGet":
		return o.OrElseGet(line, scope, args...)
	case "orElseThrow":
		return o.OrElseThrow(line, args...)
	case "ifPresent":
		return o.IfPresent(line, scope, args...)
	case "ifPresentOrElse":
		return o.IfPresentOrElse(line, scope, args...)
	case "filter":
		return o.Filter(line, scope, args...)
	case "map":
		return o.Map(line, scope, args...)
	case "flatMap":
		return o.FlatMap(line, scope, args...)
	}

	return NewError(line, NOMETHODERROR, method, o.Type())
}

// Returns an empty Optional instance. No value is present for this Optional.
func (o *Optional) Empty(line string, args ...Object) Object {
	if len(args) != 0 {
		return NewError(line, ARGUMENTERROR, "0", len(args))
	}

	return EMPTY
}

// Returns an Optional describing the given no-nil value.
func (o *Optional) Of(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if args[0] == NIL {
		return NewError(line, GENERICERROR, "of()'s parameter value must not be nil.")
	}

	//returns a new Optional
	return &Optional{Value: args[0]}
}

//Returns an Optional describing the given value, if non-nil, otherwise returns an empty Optional.
func (o *Optional) OfNullable(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if args[0] == NIL {
		return EMPTY
	}
	return o.Of(line, args...)
}

// If a value is present, returns the value, otherwise panic
func (o *Optional) Get(line string, args ...Object) Object {
	if len(args) != 0 {
		return NewError(line, ARGUMENTERROR, "0", len(args))
	}

	if o.Value == NIL {
		return NewError(line, GENERICERROR, "Option's value not present.")
	}
	return o.Value
}

// If a value is present, returns true, otherwise false.
func (o *Optional) IsPresent(line string, args ...Object) Object {
	if len(args) != 0 {
		return NewError(line, ARGUMENTERROR, "0", len(args))
	}

	if o.Value != NIL {
		return TRUE
	}
	return FALSE
}

// If a value is present, returns an Optional describing the value,
// otherwise returns an Optional produced by the supplying function.
func (o *Optional) Or(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if o.IsPresent(line) == TRUE {
		return o
	}

	supplier, ok := args[0].(*Function)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "or", "*Function", args[0].Type())
	}

	s := NewScope(scope, nil)
	opt := Eval(supplier.Literal.Body, s) // run the function
	if opt.Type() != OPTIONAL_OBJ {       // the supplier function must return an Optional
		return NewError(line, GENERICERROR, "The supplier function must return an optional.")
	}
	return opt
}

// If a value is present, returns the value, otherwise returns other.
func (o *Optional) OrElse(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if o.Value != NIL {
		return o.Value
	}
	return args[0]
}

// If a value is present, returns the value, otherwise returns the result
// produced by the supplying function.
func (o *Optional) OrElseGet(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if o.Value != NIL {
		return o.Value
	}

	supplier, ok := args[0].(*Function)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "orElseGet", "*Function", args[0].Type())
	}

	s := NewScope(scope, nil)
	ret := Eval(supplier.Literal.Body, s) // run the function
	return ret
}

// If a value is present, returns the value, otherwise throws an exception
// produced by the exception supplying function.
func (o *Optional) OrElseThrow(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if o.Value != NIL {
		return o.Value
	}

	exceptStr, ok := args[0].(*String)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "orElseThrow", "*String", args[0].Type())
	}

	//just like 'evalThrowStatement's return.
	return &Error{Kind: THROWNOTHANDLED, Message: exceptStr.String}
}

// If a value is present, performs the given action with the value,
// otherwise does nothing.
func (o *Optional) IfPresent(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if o.Value == NIL {
		return NIL // do nothing
	}

	action, ok := args[0].(*Function)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "ifPresent", "*Function", args[0].Type())
	}

	s := NewScope(scope, nil)
	s.Set(action.Literal.Parameters[0].(*ast.Identifier).Value, o.Value)
	ret := Eval(action.Literal.Body, s) // run the function

	return ret
}

// If a value is present, performs the given action with the value,
// otherwise performs the given empty-based action.
func (o *Optional) IfPresentOrElse(line string, scope *Scope, args ...Object) Object {
	if len(args) != 2 {
		return NewError(line, ARGUMENTERROR, "2", len(args))
	}

	if o.Value != NIL {
		action, ok := args[0].(*Function)
		if !ok {
			return NewError(line, PARAMTYPEERROR, "first", "ifPresentOrElse", "*Function", args[0].Type())
		}

		s := NewScope(scope, nil)
		s.Set(action.Literal.Parameters[0].(*ast.Identifier).Value, o.Value)
		ret := Eval(action.Literal.Body, s) // run the function
		return ret
	} else {
		emptyAction, ok := args[1].(*Function)
		if !ok {
			return NewError(line, PARAMTYPEERROR, "second", "ifPresentOrElse", "*Function", args[1].Type())
		}

		s := NewScope(scope, nil)
		ret := Eval(emptyAction.Literal.Body, s) // run the function
		return ret
	}
}

// If a value is present, and the value matches the given predicate,
// returns an Optional describing the value, otherwise returns an
// empty Optional.
func (o *Optional) Filter(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if o.IsPresent(line) == FALSE {
		return o
	}

	predicate, ok := args[0].(*Function)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "filter", "*Function", args[0].Type())
	}

	s := NewScope(scope, nil)
	s.Set(predicate.Literal.Parameters[0].(*ast.Identifier).Value, o.Value)
	cond := Eval(predicate.Literal.Body, s) // run the function
	if IsTrue(cond) {
		return o
	}
	return EMPTY
}

// If a value is present, returns an Optional describing (as if by
// ofNullable) the result of applying the given mapping function to
// the value, otherwise returns an empty Optional.
//
// If the mapping function returns a nil result then this method
// returns an empty Optional.
func (o *Optional) Map(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if o.IsPresent(line) == FALSE {
		return EMPTY
	}

	mapper, ok := args[0].(*Function)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "map", "*Function", args[0].Type())
	}

	s := NewScope(scope, nil)
	s.Set(mapper.Literal.Parameters[0].(*ast.Identifier).Value, o.Value)
	r := Eval(mapper.Literal.Body, s) // run the function
	if obj, ok := r.(*ReturnValue); ok {
		r = obj.Value
	}

	return o.OfNullable(line, r)
}

// If a value is present, returns the result of applying the given
// Optional-bearing mapping function to the value, otherwise returns
// an empty Optional.
//
// This method is similar to the 'map' function, but the mapping
// function is one whose result is already an Optional, and if
// invoked, flatMap does not wrap it within an additional Optional.

func (o *Optional) FlatMap(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if o.IsPresent(line) == FALSE {
		return EMPTY
	}

	mapper, ok := args[0].(*Function)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "flatMap", "*Function", args[0].Type())
	}

	s := NewScope(scope, nil)
	s.Set(mapper.Literal.Parameters[0].(*ast.Identifier).Value, o.Value)
	r := Eval(mapper.Literal.Body, s) // run the function
	if obj, ok := r.(*ReturnValue); ok {
		r = obj.Value
	}

	if r.Type() != OPTIONAL_OBJ {
		return NewError(line, GENERICERROR, "flatmap() function's return value not an optional.")
	}
	return r
}
