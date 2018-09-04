package domain

import (
	"fmt"
	"strings"

	"github.com/stoewer/go-strcase"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

type Model struct {
	Name        string
	Internal    string
	Description string
	Fields      []Field
	Deprecated  string
	Methods     []Method
}

type Field struct {
	Name         string
	Internal     string
	Description  string
	Type         string
	Indexes      []string
	Deprecated   string
	Relationship string
}

type Method struct {
	Name   string
	Args   []Arg
	Return string
}

type Arg struct {
	Name string
	Type string
}

type Index struct {
	Name   string
	Fields []Field
}

func (t Index) MethodName() string {
	return func() string {
		names := make([]string, 0, len(t.Fields))
		for _, f := range t.Fields {
			names = append(names, f.Internal)
		}
		return strings.Join(names, "")
	}()
}

func (t Model) Indexes() []Index {
	indexes := make(map[string][]Field)
	for _, field := range t.Fields {
		for _, index := range field.Indexes {
			_, found := indexes[index]
			if !found {
				indexes[index] = make([]Field, 0, 10)
			}
			indexes[index] = append(indexes[index], field)
		}
	}

	i := make([]Index, 0, len(indexes))
	for name, fields := range indexes {
		i = append(i, Index{
			Name:   name,
			Fields: fields,
		})
	}
	return i
}

func GenerateModels(config *config.Config) []Model {

	models := make([]Model, 0, len(config.Models))
	for _, m := range config.Models {

		fields := make([]Field, 0, len(m.Fields))
		for _, f := range m.Fields {

			if f.Type == "manytomany" {
				models = append(models, Model{
					Name:        CamelCasef("%s_%s", m.Name, f.Name),
					Description: fmt.Sprintf("Link %s %s", m.Name, f.Relationship),
					Fields: []Field{
						Field{
							Name:     "Id",
							Internal: "id",
							Type:     config.TypeMapping("id"),
						},
						Field{
							Name:     CamelCasef("%s_id", m.Name),
							Internal: fmt.Sprintf("%s_id", m.Name),
							Type:     config.TypeMapping("id"),
						},
						Field{
							Name:     CamelCasef("%s_id", f.Relationship),
							Internal: fmt.Sprintf("%s_id", f.Relationship),
							Type:     config.TypeMapping("id"),
						},
					},
				})
				continue
			}

			fields = append(fields, Field{
				Name: strcase.UpperCamelCase(f.Name),
				Internal: func() string {
					if f.Internal == "" {
						return strcase.SnakeCase(f.Name)
					}
					return f.Internal
				}(),
				Description: f.Description,
				Deprecated:  f.Deprecated,
				Type:        config.TypeMapping(f.Type),
				Indexes:     f.Indexes,
			})
		}

		models = append(models, Model{
			Name: strcase.UpperCamelCase(m.Name),
			Internal: func() string {
				if m.Internal == "" {
					return strcase.SnakeCase(m.Name)
				}
				return m.Internal
			}(),
			Description: m.Description,
			Deprecated:  m.Deprecated,
			Fields:      fields,
		})
	}

	return PopulateMethods(models)
}

func PopulateMethods(models []Model) []Model {

	for i, m := range models {
		methods := make([]Method, 0, len(m.Fields))
		for _, index := range m.Indexes() {
			verb := "List"
			if index.Name == "primary" || strings.Contains(index.Name, "unique") {
				verb = "Get"
			}

			args := make([]Arg, 0, len(m.Fields))
			for _, f := range index.Fields {
				args = append(args, Arg{
					Name: LowerCamelCasef(f.Internal),
					Type: f.Type,
				})
			}

			methods = append(methods, Method{
				Name: CamelCasef(
					"%s_%s%s_by_%s",
					verb,
					m.Name,
					func() string {
						if verb == "List" {
							return "s"
						}
						return ""
					}(),
					index.MethodName(),
				),
				Args: args,
				Return: func() string {
					prefix := "*"
					if verb == "List" {
						prefix = "[]*"
					}
					return prefix + CamelCasef("%s", m.Name)
				}(),
			})
		}

		models[i].Methods = methods
	}

	return models
}

func CamelCasef(f string, args ...interface{}) string {
	return strcase.UpperCamelCase(fmt.Sprintf(f, args...))
}

func LowerCamelCasef(f string, args ...interface{}) string {
	return strcase.LowerCamelCase(fmt.Sprintf(f, args...))
}
