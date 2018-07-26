package domain

import (
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
			Description: Description(typ.Description()),
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
			Description: Description(f.Description()),
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

		models = append(models, Model{
			Name:        t.Name,
			Fields:      Fields(typ),
			Description: Description(typ.Description()),
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
			ReturnType: ToType(f.Type()),
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

/*
func TypeBuilder(obj *introspection.Type) string {
	if obj == nil {
		return "NULL"
	}

	return string(ToType(obj))
}
*/
