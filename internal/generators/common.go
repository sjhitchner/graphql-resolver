package generators

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
	"github.com/stoewer/go-strcase"
)

const (
//Object = "OBJECT"
//Enum   = "ENUM"
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
	mlist := make(map[string]struct{})
	for _, value := range values {
		switch v := value.(type) {
		case []string:
			for _, s := range v {
				mlist[s] = struct{}{}
			}
		case string:
			mlist[v] = struct{}{}
		case int, int64:
			mlist[fmt.Sprintf("%d", v)] = struct{}{}
		case float32, float64:
			mlist[fmt.Sprintf("%f", v)] = struct{}{}
		case bool:
			mlist[fmt.Sprintf("%t", v)] = struct{}{}
		}
	}

	list := make([]string, 0, len(mlist))
	for k, _ := range mlist {
		list = append(list, k)
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

func GoType(values ...interface{}) (string, error) {
	s, ok := values[0].(string)
	if !ok {
		return "", errors.Errorf("Invalud argument '%s'", values[0])
	}
	switch s {
	case "integer":
		return "int64", nil
	case "string":
		return "string", nil
	case "float":
		return "float64", nil
	case "boolean":
		return "bool", nil
	case "timestamp":
		return "time.Time", nil
	default:
		return strcase.UpperCamelCase(s), nil
	}
}

func GraphQLType(values ...interface{}) (string, error) {
	f, ok := values[0].(domain.Field)
	if !ok {
		return "", errors.Errorf("Invalud argument '%s'", values[0])
	}
	return GraphQLTypeInternal(f.Type, f.Primative), nil
}

func GraphQLTypeInternal(typ, primative string) string {
	switch typ {
	case "id":
		return "ID"
	case "integer":
		return "Int"
	case "string":
		return "String"
	case "float":
		return "Float"
	case "boolean":
		return "Boolean"
	case "timestamp":
		return "String"
	default:
		return strcase.UpperCamelCase(primative)
	}
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

func Add(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() + int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) + bv.Float(), nil
		default:
			return nil, errors.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() + bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) + bv.Float(), nil
		default:
			return nil, errors.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() + float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() + float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() + bv.Float(), nil
		default:
			return nil, errors.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, errors.Errorf("add: unknown type for %q (%T)", av, a)
	}
}

// subtract returns the difference of b from a.
func Subtract(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() - int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) - bv.Float(), nil
		default:
			return nil, errors.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() - bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) - bv.Float(), nil
		default:
			return nil, errors.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() - float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() - float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() - bv.Float(), nil
		default:
			return nil, errors.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, errors.Errorf("subtract: unknown type for %q (%T)", av, a)
	}
}

// multiply returns the product of a and b.
func Multiply(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() * bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() * int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) * bv.Float(), nil
		default:
			return nil, errors.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) * bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() * bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) * bv.Float(), nil
		default:
			return nil, errors.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() * float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() * float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() * bv.Float(), nil
		default:
			return nil, errors.Errorf("multiply: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, errors.Errorf("multiply: unknown type for %q (%T)", av, a)
	}
}

// divide returns the division of b from a.
func Divide(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() / bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() / int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) / bv.Float(), nil
		default:
			return nil, errors.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) / bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() / bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) / bv.Float(), nil
		default:
			return nil, errors.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() / float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() / float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() / bv.Float(), nil
		default:
			return nil, errors.Errorf("divide: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, errors.Errorf("divide: unknown type for %q (%T)", av, a)
	}
}

func Now() (interface{}, error) {
	format := time.RFC3339
	return time.Now().UTC().Format(format), nil
}

func Find(values ...interface{}) (interface{}, error) {
	if len(values) != 2 {
		return "", errors.Errorf("Invalid argument '%s'", values[0])
	}

	name := values[1].(string)

	switch v := values[0].(type) {
	case []domain.Model:
		for _, m := range v {
			if m.Name == name {
				return m, nil
			}
		}

	//case []domain.Relationship:
	//	for _, m := range v {
	//		if m.Name == name {
	//			return m, nil
	//		}
	//	}

	case []domain.Field:
		for _, m := range v {
			if m.Name == name {
				return m, nil
			}
		}
	}
	panic("Invalid find type")
}

func MethodReturn(value interface{}) (interface{}, error) {
	methodReturn, ok := value.(domain.Return)
	if !ok {
		return "", errors.Errorf("Invalid argument type")
	}

	str := strcase.UpperCamelCase(methodReturn.Type)

	if methodReturn.Multi {
		return fmt.Sprintf("[]*%s", str), nil
	}

	return fmt.Sprintf("*%s", str), nil
}

func IsMany2Many(value interface{}) (bool, error) {
	relationship, ok := value.(*domain.Relationship)
	if !ok {
		return false, errors.Errorf("Invalid argument type")
	}

	if relationship == nil {
		return false, nil
	}

	if relationship.Type == config.Many2Many {
		return true, nil
	}

	return false, nil
}

/*
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
