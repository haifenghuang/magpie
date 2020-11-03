package eval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"magpie/ast"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	errEOS = errors.New("End of slice")
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Order []HashKey
	Pairs map[HashKey]HashPair
}

func NewHash() *Hash {
	return &Hash{Order: []HashKey{}, Pairs: make(map[HashKey]HashPair)}
}

func (h *Hash) iter() bool { return true }

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}

	//Iterate keys.
	for _, hk := range h.Order { //hk:hash key
		var key, val string

		pair, _ := h.Pairs[hk]
		if pair.Key.Type() == STRING_OBJ {
			key = "\"" + pair.Key.Inspect() + "\""
		} else {
			key = pair.Key.Inspect()
		}

		if pair.Value.Type() == STRING_OBJ {
			val = "\"" + pair.Value.Inspect() + "\""
		} else {
			val = pair.Value.Inspect()
		}
		pairs = append(pairs, fmt.Sprintf("%s : %s", key, val))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

func (h *Hash) CallMethod(line string, scope *Scope, method string, args ...Object) Object {
	switch method {
	case "clear":
		return h.Clear(line, args...)
	case "exists", "has":
		return h.Exists(line, args...)
	case "find", "index":
		return h.Index(line, args...)
	case "filter":
		return h.Filter(line, scope, args...)
	case "get":
		return h.Get(line, args...)
	case "getPath":
		return h.GetPath(line, args...)
	case "keys":
		return h.Keys(line, args...)
	case "len":
		return h.Len(line, args...)
	case "map":
		return h.Map(line, scope, args...)
	case "merge":
		return h.Merge(line, args...)
	case "pop", "delete", "remove":
		return h.Pop(line, args...)
	case "push", "set":
		return h.Push(line, args...)
	case "values":
		return h.Values(line, args...)
	}

	return NewError(line, NOMETHODERROR, method, h.Type())
}

func (h *Hash) Clear(line string, args ...Object) Object {
	if len(args) != 0 {
		return NewError(line, ARGUMENTERROR, "0", len(args))
	}

	h.Order = nil
	h.Pairs = nil

	return NIL
}

func (h *Hash) Exists(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}

	if len(h.Order) == 0 {
		return FALSE
	}
	hashable, ok := args[0].(Hashable)
	if !ok {
		return NewError(line, KEYERROR, args[0].Type())
	}

	_, ok = h.Pairs[hashable.HashKey()]
	if ok {
		return TRUE
	}
	return FALSE
}

func (h *Hash) Filter(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}
	block, ok := args[0].(*Function)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "filter", "*Function", args[0].Type())
	}
	if len(block.Literal.Parameters) != 2 {
		return NewError(line, ARGUMENTERROR, "2", len(block.Literal.Parameters))
	}
	hash := NewHash()
	s := NewScope(scope, nil)
	for _, hk := range h.Order { //hk:hash key
		argument, _ := h.Pairs[hk]
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument.Key)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, argument.Value)
		cond := Eval(block.Literal.Body, s)
		if IsTrue(cond) {
			hash.Push(line, argument.Key, argument.Value)
		}
	}
	return hash
}

func (h *Hash) Index(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}
	hashable, ok := args[0].(Hashable)
	if !ok {
		return NewError(line, KEYERROR, args[0].Type())
	}

	hk := hashable.HashKey()
	for idx, k := range h.Order {
		r := reflect.DeepEqual(hk, k)
		if r {
			return NewInteger(int64(idx))
		}
	}

	return NewInteger(-1)
}

func (h *Hash) Get(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}
	hashable, ok := args[0].(Hashable)
	if !ok {
		return NewError(line, KEYERROR, args[0].Type())
	}
	if hashPair, ok := h.Pairs[hashable.HashKey()]; ok {
		return hashPair.Value
	}
	return NIL
}

/*
    Code copied from "go-json-map" with modifications

    document = {
        one: {
            "two": {
                "three": [1, 2, 3,]
            },
            "four": {
                "five": [ 11, 22, 33,]
            },
        },
    }
	prop = document.getPath("one.two.three[0]")
	println(prop) // prop => 1

	prop= document.getPath("one.two.three")
	println(prop) // prop => [1, 2, 3]
*/
func (h *Hash) GetPath(line string, args ...Object) Object {
	if len(args) != 1 && len(args) != 2 {
		return NewError(line, ARGUMENTERROR, "1|2", len(args))
	}

	pathObj, ok := args[0].(*String)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "getPath", "*String", args[0].Type())
	}

	separator := "."
	if len(args) == 2 {
		sepObj, ok := args[1].(*String)
		if !ok {
			return NewError(line, PARAMTYPEERROR, "second", "getPath", "*String", args[1].Type())
		}
		separator = sepObj.String
	}

	data := make(map[HashKey]Object)
	for _, hk := range h.Order { //hk:hash key
		pair, _ := h.Pairs[hk]
		data[hk] = pair.Value
	}

	path := pathObj.String
	err := fmt.Errorf("Property '%s' does not exist", path)
	if len(path) == 0 {
		path = separator
	}

	levels := strings.Split(path, separator)
	if len(levels) > 0 && path != separator {
		path_level_one := levels[0]
		// If we have a level in path_level_one

		re := regexp.MustCompile(`(\w+[\_]?[\-]?)+\[\d+\]{1}`)
		if matched := re.FindString(path_level_one); len(matched) > 0 {
			property_re := regexp.MustCompile(`(\w+[\_]?[\-]?)+`)
			index_re := regexp.MustCompile(`\[\d+\]{1}`)

			// Get property
			property := property_re.FindString(path_level_one)

			// Get index
			index_found := index_re.FindString(path_level_one)

			// If index > 0 - check if this property is array
			if len(index_found) > 0 {
				if len(property) > 0 {
					path_level_one = property
				}
				index_found = strings.Trim(index_found, "[]")
				if index, err := strconv.Atoi(index_found); err == nil {
					hk := NewString(property).HashKey()
					if v, ok := data[hk]; ok {
						if v.Type() == ARRAY_OBJ {
							arrLenObj := v.(*Array).Len(line)
							arrLen := arrLenObj.(*Integer).Int64
							if index >= 0 && int64(index) < arrLen {
								value := v.(*Array).Get(line, NewInteger(int64(index)))
								data[hk] = value
							} else {
								err = fmt.Errorf("%s: Index out of bounds(0, %d). got index %d", property, arrLen-1, index)
								return NewNil(err.Error())
							}
						} else {
							err = fmt.Errorf("%s: is not an array", property)
							return NewNil(err.Error())
						}
					} else {
						err = fmt.Errorf("Property '%s' does not exist", property)
						return NewNil(err.Error())
					}
				} else {
					err = fmt.Errorf("%s must be of type %s", fmt.Sprintf("%s[%s]", property, index_found), "number")
					return NewNil(err.Error())
				}
			}
		}

		if len(levels[1:]) >= 1 {
			hk := NewString(path_level_one).HashKey()
			if level_one_value, ok := data[hk]; ok {
				if level_one_value != nil {
					switch level_one_value.Type() {
					case HASH_OBJ:
						tmpHash := level_one_value.(*Hash)
						p := strings.Join(levels[1:], separator)
						return tmpHash.GetPath(line, NewString(p), NewString(separator))
					default:
						// pass
					}
				}
			} else {
				err := fmt.Errorf("Property '%s' does not exist", path_level_one)
				return NewNil(err.Error())
			}
		} else {
			hk := NewString(path_level_one).HashKey()
			if v, ok := data[hk]; ok {
				return v
			}
		}
	} else if path == separator {
		return h
	}

	return NewNil(err.Error())
}

func (h *Hash) Keys(line string, args ...Object) Object {
	if len(args) != 0 {
		return NewError(line, ARGUMENTERROR, "0", len(args))
	}
	keys := &Array{}
	for _, hk := range h.Order { //hk:hash key
		pair, _ := h.Pairs[hk]
		keys.Members = append(keys.Members, pair.Key)
	}
	return keys
}

func (h *Hash) Len(line string, args ...Object) Object {
	if len(args) != 0 {
		return NewError(line, ARGUMENTERROR, "0", len(args))
	}

	return NewInteger(int64(len(h.Order)))
}

func (h *Hash) Map(line string, scope *Scope, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}
	block, ok := args[0].(*Function)
	if !ok {
		return NewError(line, PARAMTYPEERROR, "first", "map", "*Function", args[0].Type())
	}
	if len(block.Literal.Parameters) != 2 {
		return NewError(line, ARGUMENTERROR, "2", len(block.Literal.Parameters))
	}
	hash := NewHash()
	s := NewScope(scope, nil)
	for _, hk := range h.Order { //hk:hash key
		argument, _ := h.Pairs[hk]
		s.Set(block.Literal.Parameters[0].(*ast.Identifier).Value, argument.Key)
		s.Set(block.Literal.Parameters[1].(*ast.Identifier).Value, argument.Value)
		r := Eval(block.Literal.Body, s)
		if obj, ok := r.(*ReturnValue); ok {
			r = obj.Value
		}
		rh, ok := r.(*Hash)
		if !ok {
			NewError(line, RTERROR, HASH_OBJ)
		}
		for _, hk := range rh.Order { //hk:hash key
			v, _ := rh.Pairs[hk]
			hash.Push(line, v.Key, v.Value)
		}
	}
	return hash
}

func (h *Hash) Merge(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}
	m, ok := args[0].(*Hash)
	if !ok {
		return NewError(line, ARGUMENTERROR, args[0].Type(), "hash.merge")
	}

	hash := NewHash()
	for _, hk := range h.Order {
		pair, _ := h.Pairs[hk]
		hash.Push(line, pair.Key, pair.Value)
	}

	for _, hk := range m.Order {
		pair, _ := m.Pairs[hk]
		hash.Push(line, pair.Key, pair.Value)
	}

	return hash
}

func (h *Hash) Pop(line string, args ...Object) Object {
	if len(args) != 1 {
		return NewError(line, ARGUMENTERROR, "1", len(args))
	}
	hashable, ok := args[0].(Hashable)
	if !ok {
		return NewError(line, KEYERROR, args[0].Type())
	}

	hk := hashable.HashKey()
	if hashPair, ok := h.Pairs[hk]; ok {
		// remove the 'key' from 'Order' array of Hash.
		for idx, k := range h.Order {
			r := reflect.DeepEqual(hk, k)
			if r {
				h.Order = append(h.Order[:idx], h.Order[idx+1:]...)
				break
			}
		}

		// remove from the pair, and return the deleted value.
		delete(h.Pairs, hk)
		return hashPair.Value
	}
	return NIL
}

func (h *Hash) Push(line string, args ...Object) Object {
	if len(args) != 2 {
		return NewError(line, ARGUMENTERROR, "2", len(args))
	}
	if hashable, ok := args[0].(Hashable); ok {
		hk := hashable.HashKey()
		// Check if it already exists, if not, add to the 'Order' array of the OrderedHash.
		_, exists := h.Pairs[hk]
		if !exists {
			h.Order = append(h.Order, hk)
		}

		h.Pairs[hk] = HashPair{Key: args[0], Value: args[1]}
	} else {
		return NewError(line, KEYERROR, args[0].Type())
	}
	return h
}

func (h *Hash) Values(line string, args ...Object) Object {
	if len(args) != 0 {
		return NewError(line, ARGUMENTERROR, "0", len(args))
	}

	values := &Array{}
	//Iterate keys.
	for _, hk := range h.Order { //hk:hash key
		pair, _ := h.Pairs[hk]
		values.Members = append(values.Members, pair.Value)
	}

	return values
}

//Json marshal handling
func (h *Hash) MarshalJSON() ([]byte, error) {
	if len(h.Pairs) == 0 {
		return []byte("{}"), nil
	}

	var out bytes.Buffer

	var first = true
	out.WriteString("{")

	for _, hk := range h.Order {
		pair, _ := h.Pairs[hk]

		if first {
			first = false
		} else {
			out.WriteString(",")
		}

		var res bytes.Buffer
		var err error
		res, err = marshalJsonObject(pair.Key)
		if err != nil {
			return []byte{}, err
		}
		out.WriteString(res.String())

		out.WriteString(":")

		res, err = marshalJsonObject(pair.Value)
		if err != nil {
			return []byte{}, err
		}
		out.WriteString(res.String())

	} //end for
	out.WriteString("}")

	return out.Bytes(), nil
}

func (h *Hash) UnmarshalJSON(b []byte) error {
	//Using Decoder to parse the bytes.
	in := bytes.TrimSpace(b)
	dec := json.NewDecoder(bytes.NewReader(in))

	t, err := dec.Token()
	if err != nil {
		return err
	}

	// must open with a delim token '{'
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expect JSON object open with '{'")
	}

	h.unmarshalJSON(dec)

	t, err = dec.Token() //'}'
	if err != nil {
		return err
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expect JSON object close with '}'")
	}

	return nil
}

func (h *Hash) unmarshalJSON(dec *json.Decoder) error {
	for dec.More() { // Loop until it has no more tokens
		t, err := dec.Token()
		if err != nil {
			return err
		}

		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("key must be a string, got %T\n", t)
		}

		val, err := parseObject(dec)
		if err != nil {
			return err
		}
		h.Push("", NewString(key), val)
	}
	return nil
}

func parseObject(dec *json.Decoder) (Object, error) {
	t, err := dec.Token()
	if err != nil {
		return NIL, err
	}

	switch tok := t.(type) {
	case json.Delim:
		switch tok {
		case '[': // If it's an array
			return parseArray(dec)
		case '{': // If it's a map
			h := NewHash()
			err := h.unmarshalJSON(dec)
			if err != nil {
				return NIL, err
			}
			_, err = dec.Token() // }
			if err != nil {
				return NIL, err
			}
			return h, nil
		case ']':
			return NIL, errEOS
		case '}':
			return NIL, errors.New("unexpected '}'")
		default:
			return NIL, fmt.Errorf("Unexpected delimiter: %q", tok)
		}
	default:
		var ret Object
		// Check the type of the token, and convert it to magpie object
		switch tok.(type) {
		case float64:
			ret = NewFloat(tok.(float64))
		case bool:
			b := tok.(bool)
			if b {
				ret = TRUE
			} else {
				ret = FALSE
			}
		case string:
			ret = NewString(tok.(string))
		case nil:
			ret = NIL
		}

		return ret, nil
	}
}

func parseArray(dec *json.Decoder) (Object, error) {
	arr := &Array{}
	for {
		v, err := parseObject(dec)
		if err == errEOS {
			return arr, nil
		}
		if err != nil {
			return NIL, err
		}
		arr.Members = append(arr.Members, v)
	}
}
