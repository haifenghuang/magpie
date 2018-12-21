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
	ret := &Optional{Value:NIL}
	SetGlobalObj(optional_name, ret)

	return ret
}

func (o *Optional) Inspect() string {
	if o.Value != NIL {
		return "Optional["+ o.Value.Inspect() + "]"
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
	case "ifPresent":
		return o.IfPresent(line, scope, args...)
	}

	panic(NewError(line, NOMETHODERROR, method, o.Type()))
}

// Returns an empty Optional instance. No value is present for this Optional.
func (o *Optional) Empty(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	return EMPTY
}

// Returns an Optional describing the given no-nil value.
func (o *Optional) Of(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if args[0] == NIL {
		panic(NewError(line, GENERICERROR, "of()'s parameter value must not be nil."))
	}

	//returns a new Optional
	return &Optional{Value:args[0]}
}

//Returns an Optional describing the given value, if non-nil, otherwise returns an empty Optional.
func (o *Optional) OfNullable(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if o.Value == NIL {
		return EMPTY
	}
	return o.Of(line, args...)
}

// If a value is present, returns the value, otherwise panic
func (o *Optional) Get(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
	}

	if o.Value == NIL {
		panic(NewError(line, GENERICERROR, "Option's value not present."))
	}
	return o.Value
}

// If a value is present, returns true, otherwise false.
func (o *Optional) IsPresent(line string, args ...Object) Object {
	if len(args) != 0 {
		panic(NewError(line, ARGUMENTERROR, "0", len(args)))
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
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if o.IsPresent(line, args...) == TRUE {
		return o
	}

	supplier, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "or", "*Function", args[0].Type()))
	}

	s := NewScope(scope)
	opt := Eval(supplier.Literal.Body, s) // run the function
	if opt.Type() != OPTIONAL_OBJ { // the supplier function must return an Optional
		panic(NewError(line, GENERICERROR, "The supplier function must return an optional."))
	}
	return opt
}

// If a value is present, returns the value, otherwise returns other.
func (o *Optional) OrElse(line string, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
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
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if o.Value != NIL {
		return o.Value
	}

	supplier, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "orElseGet", "*Function", args[0].Type()))
	}

	s := NewScope(scope)
	ret := Eval(supplier.Literal.Body, s) // run the function
	return ret
}


// If a value is present, performs the given action with the value,
// otherwise does nothing.
func (o *Optional) IfPresent(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		panic(NewError(line, ARGUMENTERROR, "1", len(args)))
	}

	if o.Value == NIL {
		return NIL // do nothing
	}

	action, ok := args[0].(*Function)
	if !ok {
		panic(NewError(line, PARAMTYPEERROR, "first", "ifPresent", "*Function", args[0].Type()))
	}

	s := NewScope(scope)
	s.Set(action.Literal.Parameters[0].(*ast.Identifier).Value, o.Value)
	ret := Eval(action.Literal.Body, s) // run the function
	
	return ret
}
