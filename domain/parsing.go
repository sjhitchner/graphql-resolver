package domain

/*
import (
	"fmt"
	"strings"

	"github.com/graph-gophers/graphql-go/introspection"
)

func Enums(schema *introspection.Schema) []Enum {
	enums := make([]Enum, 0, len(schema.Types()))
	for _, typ := range schema.Types() {

		t := ToType(typ)
		if t.Base != ENUM {
			continue
		}
		if strings.HasPrefix(t.Name, "__") {
			continue
		}

		enums = append(enums, Enum{
			Name:        t.Name,
			Values:      EnumValues(typ),
			Description: StringRef(typ.Description()),
		})
	}
	return enums
}

func EnumValues(t *introspection.Type) []Field {
	f := t.EnumValues(&struct{ IncludeDeprecated bool }{true})

	if f == nil {
		return nil
	}

	fields := make([]Field, 0, len(*f))

	for _, f := range *f {
		fields = append(fields, Field{
			Name:        f.Name(),
			Description: StringRef(f.Description()),
			Deprecated: func() string {
				if f.IsDeprecated() {
					return "DECPRECATED " + StringRef(f.DeprecationReason())
				}
				return ""
			}(),
		})
	}

	return fields
}

func Models(schema *introspection.Schema) []Model {
	models := make([]Model, 0, len(schema.Types()))
	for _, typ := range schema.Types() {

		t := ToType(typ)
		if t.Base != Object {
			continue
		}
		if strings.HasPrefix(t.Name, "__") {
			continue
		}

		name, modelType := Name(t.Name)
		models = append(models, Model{
			Name:        name,
			Type:        modelType,
			Fields:      Fields(typ),
			Description: StringRef(typ.Description()),
		})
	}
	return models
}

func Name(name string) (string, ModelType) {
	if strings.HasSuffix(name, "Connection") {
		return name, Connection
	}

	if strings.HasSuffix(name, "Edge") {
		return name, Edge
	}

	if strings.HasSuffix(name, "PageInfo") {
		return name, PageInfo
	}

	return name, Resolver
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
			Description: StringRef(f.Description()),
			Deprecated: func() string {
				if f.IsDeprecated() {
					return "DECPRECATED " + StringRef(f.DeprecationReason())
				}
				return ""
			}(),
			Args: Args(f),
			Type: ToType(f.Type()),
		})
	}

	return fields
}

func Args(f *introspection.Field) []Argument {
	fargs := f.Args()
	if len(fargs) == 0 {
		return nil
	}

	args := make([]Argument, 0, len(fargs))
	for _, a := range fargs {
		args = append(args, Argument{
			Name:        a.Name(),
			Description: StringRef(a.Description()),
			Type:        ToType(a.Type()),
			Default:     StringRef(a.DefaultValue()),
		})
	}
	return args
}

func StringRef(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type ParsedSchema struct {
	Models []Model
	Enums  []Enum
}

func ParseSchema(schema *introspection.Schema) *ParsedSchema {

	for _, d := range schema.Directives() {
		fmt.Println("Name:", d.Name())
		fmt.Println("Description:", StringRef(d.Description()))
		for _, a := range d.Args() {
			fmt.Println("\tArg Name:", a.Name())
			fmt.Println("\tArg Description:", StringRef(a.Description()))
		}
		for _, l := range d.Locations() {
			fmt.Println("Location:", l)
		}
	}

	models := Models(schema)
	enums := Enums(schema)

	return &ParsedSchema{
		Models: models,
		Enums:  enums,
	}
}

func TypeBuilder(obj *introspection.Type) string {
	if obj == nil {
		return "NULL"
	}

	return string(ToType(obj))
}
*/
