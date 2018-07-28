package generators

import (
	"reflect"
	"strings"

	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/stoewer/go-strcase"
)

const (
	Object = "OBJECT"
	Enum   = "ENUM"
)

func SafeString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func SafeHasPrefix(str *string, s string) bool {
	if str == nil {
		return false
	}
	return strings.HasPrefix(*str, "__")
}

func Args(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func SnakeCase(values ...interface{}) (string, error) {
	s, ok := values[0].(string)
	if !ok {
		return "", errors.Errorf("Invalud argument '%s'", values[0])
	}
	return strcase.SnakeCase(s), nil
}

func CamelCase(values ...interface{}) (string, error) {
	s, ok := values[0].(string)
	if !ok {
		return "", errors.Errorf("Invalud argument '%s'", values[0])
	}
	return strcase.UpperCamelCase(strcase.SnakeCase(s)), nil
}

func LowerCamelCase(values ...interface{}) (string, error) {
	s, ok := values[0].(string)
	if !ok {
		return "", errors.Errorf("Invalud argument '%s'", values[0])
	}
	return strcase.LowerCamelCase(strcase.SnakeCase(s)), nil
}

func Comment(values ...interface{}) (string, error) {
	s, ok := values[0].(string)
	if !ok {
		return "", nil
	}
	if s == "" {
		return "", nil
	}
	return "// " + strings.Replace(s, "\n", "\n// ", -1), nil
}

func Safe(values ...interface{}) (interface{}, error) {

	/*
		t:=reflect.TypeOf(values[0])
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		switch v := values[0].(type) {
		case string:
	*/
	return "", nil
}

/*
	"isQuery": func(values ...interface{}) bool {
		b, ok := values[0].(ModelType)
		return ok && b == Query
	},
	"isResolver": func(values ...interface{}) bool {
		b, ok := values[0].(ModelType)
		return ok && b == Resolver
	},
	"isPageInfo": func(values ...interface{}) bool {
		b, ok := values[0].(ModelType)
		return ok && b == PageInfo
	},
	"isEdge": func(values ...interface{}) bool {
		b, ok := values[0].(ModelType)
		return ok && b == Edge
	},
	"isConnection": func(values ...interface{}) bool {
		b, ok := values[0].(ModelType)
		return ok && b == Connection
	},
*/

func TypeName(values ...interface{}) (string, error) {
	t, ok := values[0].(*introspection.Type)
	if !ok {
		return "", err
	}

	return nameRecurse(t), nil
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
		return "*" + nameRecurse(t.OfType())
	}

	panic("Invalid Type " + t.Kind())
}

func TypeType(values ...interface{}) (string, error) {
	t, ok := values[0].(*introspection.Type)
	if !ok {
		return "", err
	}

	return typeRecurse(t), nil
}

func typeRecurse(t *introspection.Type) string {
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
