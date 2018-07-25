package generate

import (
	"fmt"
	"os"
	"strings"

	//"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/introspection"
)

type Type int

const (
	Null Type = iota
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
	NonNull
	UserType
)

func ToType(t *introspection.Type) Type {
	if t == nil {
		return Null
	}

	switch t.Kind() {
	case "ID":
		return ID
	case "LIST":
		return List
	case "OBJECT":
		return Object
	case "INTERFACE":
		return Interface
	case "INTEGER":
		return Integer
	case "FLOAT":
		return Float
	case "STRING":
		return String
	case "BOOLEAN":
		return Boolean
	case "NON_NULL":
		return NonNull
	default:
		return UserType
	}
}

func (t Type) String() string {
	switch t {
	case Integer:
		return "Int"
	case Boolean:
		return "Boolean"
	case String:
		return "String"
	case Float:
		return "Float"
	case Time:
		return "Time"
	case DateTime:
		return "DateTime"
	case ID:
		return "ID"
	default:
		return fmt.Sprintf("Invalid GRAPHQL Type (%d)", t)
	}
}

type ResolverTemplate struct {
	Imports   []string
	Queries   []Query
	Mutations []string
	Models    []Model
}

type Query struct {
	Name         string
	Args         []Argument
	ResolverName string
	IsList       bool
}

func (t Query) Return() string {
	str := fmt.Sprintf("*%s", t.ResolverName)
	if t.IsList {
		return "[]" + str
	}
	return str
}

type Argument struct {
	Name    string
	Type    Type
	Default string
}

func (t Argument) String() string {
	return fmt.Sprintf("%s %s%s", t.Name, t.Type, func() string {
		if t.Default != "" {
			switch t.Type {
			case String:
				return fmt.Sprintf(` = "%s"`, t.Default)
			default:
				return fmt.Sprintf(` = %s`, t.Default)
			}
		}
		return ""
	}())
}

type Model struct {
	Name        string
	Description string
	Implements  string
	Fields      []Field
}

type Field struct {
	Name        string
	Description string
	Args        []Argument
	Deprecated  string
	ReturnType  string
}

type ResolverGenerator struct {
}

func NewResolverGenerator() *ResolverGenerator {
	return &ResolverGenerator{}
}

func (t *ResolverGenerator) Generate(schema *introspection.Schema) error {

	models := Models(schema)

	return ExecuteTemplate(os.Stdout, "resolver.tmpl", ResolverTemplate{
		Models: models,
	})
}

func Models(schema *introspection.Schema) []Model {
	models := make([]Model, 0, len(schema.Types()))
	for _, t := range schema.Types() {

		if ToType(t) != Object {
			continue
		}
		if strings.HasPrefix(*t.Name(), "__") {
			continue
		}

		models = append(models, Model{
			Name:        StringRef(t.Name()),
			Fields:      Fields(t),
			Description: Description(t.Description()),
		})
	}
	return models
}

func Fields(t *introspection.Type) []Field {
	f := t.Fields(&struct{ IncludeDeprecated bool }{true})

	if f == nil {
		return nil
	}

	fields := make([]Field, 0, len(*f))

	for _, f := range *f {
		fields = append(fields, Field{
			Name:        f.Name(),
			Description: Description(f.Description()),
			Deprecated: func() string {
				if f.IsDeprecated() {
					return "DECPRECATED " + StringRef(f.DeprecationReason())
				}
				return ""
			}(),
			ReturnType: TypeBuilder(f.Type()),
		})
	}

	return fields
}

func StringRef(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Description(s *string) string {
	if s == nil {
		return ""
	}
	return "// " + strings.Replace(*s, "\n", "\n// ", -1)
}

func TypeBuilder(obj *introspection.Type) string {
	if obj == nil {
		return "NULL"
	}

	t := ToType(obj)
	switch t {
	case List:
		return "[]" + TypeBuilder(obj.OfType())
	case NonNull:
		return TypeBuilder(obj.OfType())
	case Object:
		fallthrough
	case Interface:
		return StringRef(obj.Name())
	default:
		return t.String()
	}
}
