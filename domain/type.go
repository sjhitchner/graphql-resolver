package domain

import (
	"github.com/graph-gophers/graphql-go/introspection"
)

type Type struct {
	Base BaseType
	Name string
}

type BaseType int

const (
	Null BaseType = iota
	ID
	Integer
	Float
	Boolean
	String
	Time
	DateTime
	List
	Object
	Interface
	ENUM
	Union
	InputObject
	NonNull
	UserType
)

func ToType(t *introspection.Type) Type {
	if t == nil {
		return Type{}
	}

	return Type{
		Base: typeRecurse(t),
		Name: nameRecurse(t),
	}
}

func nameRecurse(t *introspection.Type) string {
	switch t.Kind() {
	case "LIST":
		return "[]" + nameRecurse(t.OfType())

	case "OBJECT":
		fallthrough
	case "INTERFACE":
		fallthrough
	case "ENUM":
		fallthrough
	case "UNION":
		fallthrough
	case "INPUT_OBJECT":
		return *t.Name()

	case "SCALAR":
		switch *t.Name() {
		case "ID":
			return "*graphql.ID"
		case "Int":
			return "int64"
		case "Float":
			return "float64"
		case "String":
			return "string"
		case "Boolean":
			return "bool"
		default:
			panic("Invalid Type " + *t.Name())
		}

	case "NON_NULL":
		//return "*" + nameRecurse(t.OfType())
		return nameRecurse(t.OfType())
	}

	panic("Invalid Type " + t.Kind())
}

func typeRecurse(t *introspection.Type) BaseType {
	switch t.Kind() {
	case "LIST":
		return List
	case "OBJECT":
		return Object
	case "INTERFACE":
		return Interface
	case "NON_NULL":
		return typeRecurse(t.OfType())
	case "ENUM":
		return ENUM
	case "UNION":
		return Union
	case "INPUT_OBJECT":
		return InputObject
	case "SCALAR":
		switch *t.Name() {
		case "ID":
			return ID
		case "Int":
			return Integer
		case "Float":
			return Float
		case "String":
			return String
		case "Boolean":
			return Boolean
		default:
			panic("Invalid Type " + *t.Name())
		}
	default:
		panic("Invalid Type " + t.Kind())
	}
}

func (t Type) String() string {
	switch t.Base {
	case Object:
		fallthrough
	case Interface:
		fallthrough
	case ENUM:
		fallthrough
	case List:
		return t.Name

	case Integer:
		return "int64"
	case Boolean:
		return "bool"
	case String:
		return "string"
	case Float:
		return "float64"
	case Time:
		return "time.Time"
	case DateTime:
		return "time.Time"
	case ID:
		return "*graphql.ID"

	default:
		panic("Invalid GRAPHQL Type")
	}
}
