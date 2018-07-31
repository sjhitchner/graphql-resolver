package generators

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"
)

const (
	Object = "OBJECT"
	Enum   = "ENUM"
)

func Join(values ...interface{}) ([]string, error) {
	list := make([]string, 0, len(values))
	for _, value := range values {
		switch v := value.(type) {
		case []string:
			list = append(list, v...)
		case string:
			list = append(list, v)
		case int, int64:
			list = append(list, fmt.Sprintf("%d", v))
		case float32, float64:
			list = append(list, fmt.Sprintf("%f", v))
		case bool:
			list = append(list, fmt.Sprintf("%t", v))
		}
	}
	return list, nil
}

func Unique(values ...interface{}) ([]string, error) {
	list := make(map[string]struct{}, 0, len(values))
	for _, value := range values {
		switch v := value.(type) {
		case []string:
			for _, s := range v {
				list[s] = struct{}{}
			}
		case string:
			list[v] = struct{}{}
		case int, int64:
			list[fmt.Sprintf("%d", v)] = struct{}{}
		case float32, float64:
			list[fmt.Sprintf("%f", v)] = struct{}{}
		case bool:
			list[fmt.Sprintf("%t", v)] = struct{}{}
		}
	}
	return list, nil
}

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
	return strings.HasPrefix(*str, s)
}

func SafeHasSuffix(str *string, s string) bool {
	if str == nil {
		return false
	}
	return strings.HasSuffix(*str, s)
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

	value := reflect.ValueOf(values[0])
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		return value.String(), nil

	case reflect.Float32, reflect.Float64:
		return value.Float(), nil

	case reflect.Bool:
		return value.Bool(), nil

	case reflect.Int, reflect.Int64:
		return value.Int(), nil

	default:
		return value.Interface(), nil
	}
}

func IsQuery(values ...interface{}) bool {
	t, ok := values[0].(*introspection.Type)
	return ok && SafeString(t.Name()) == "Query"
}

func IsPageInfo(values ...interface{}) bool {
	t, ok := values[0].(*introspection.Type)
	return ok && SafeString(t.Name()) == "PageInfo"
}

func IsEdge(values ...interface{}) bool {
	t, ok := values[0].(*introspection.Type)
	return ok && SafeHasSuffix(t.Name(), "Edge")
}

func IsConnection(values ...interface{}) bool {
	t, ok := values[0].(*introspection.Type)
	return ok && SafeHasSuffix(t.Name(), "Connection")
}

func AllFields(values ...interface{}) (interface{}, error) {
	t, ok := values[0].(*introspection.Type)
	if !ok {
		return "", errors.New("Not an introspection type")
	}

	fields := t.Fields(&struct{ IncludeDeprecated bool }{true})
	if fields == nil {
		return "", errors.New("No fields")
	}

	return *fields, nil
}

func TypeName(values ...interface{}) (interface{}, error) {
	t, ok := values[0].(*introspection.Type)
	if !ok {
		return "", errors.New("Not an introspection type")
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
		return *t.Name() + "Resolver"

	case "SCALAR":
		switch *t.Name() {
		case "ID":
			return "graphql.ID"
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

/*
func TypeType(values ...interface{}) (interface{}, error) {
	t, ok := values[0].(*introspection.Type)
	if !ok {
		return "", errors.New("Not an introspection type")
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
*/
