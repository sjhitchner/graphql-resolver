package domain

import (
	"fmt"
	"strings"

	"github.com/stoewer/go-strcase"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

type Repo struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name         string
	Args         []Arg
	ReturnType   string
	ReturnPrefix string
}

type Arg struct {
	Name string
	Type string
}

type Model struct {
	Name        string
	Description string
	Fields      []Field
	Indexes     []Index
	RepoName    string
	Methods     []Method
	SubModels   []Model
}

type Index struct {
	Name   string
	Type   string
	Fields []string
}

type Field struct {
	Name         string
	Type         string
	Primative    string
	Internal     string // Snake Case
	Relationship string
}

type Type struct {
	Name string
	Type string
}

func ProcessConfig(config *config.Config) ([]Model, []Type) {
	models := make([]Model, 0, len(config.Models))

	for _, m := range config.Models {
		models = append(
			models,
			ProcessModel(config, m),
		)
	}

	for i, m := range models {
		for _, f := range m.Fields {
			if f.Type == "manytomany" {
				models[i].SubModels = append(models[i].SubModels, Model{
					Name:        m.Name + "_" + f.Name,
					Description: fmt.Sprintf("Linking %s to %s", m.Name, f.Relationship),
					Fields: []Field{
						Field{
							Name:      "id",
							Type:      "id",
							Primative: "id",
							Internal:  "id",
						},
						Field{
							Name:      m.Name + "_id",
							Type:      "id",
							Primative: "id",
							Internal:  m.Name + "_id",
						},
						Field{
							Name:      f.Relationship + "_id",
							Type:      "id",
							Primative: "id",
							Internal:  f.Relationship + "_id",
						},
					},
					Indexes: nil,
					Methods: nil,
				})
			}
		}
	}
	return models, ProcessTypes(config, config.Types)
}

func ProcessTypes(config *config.Config, ct []config.Type) []Type {
	types := make([]Type, 0, len(ct))
	for _, t := range ct {
		types = append(types, Type{
			Name: t.Name,
			Type: config.TypePrimative(t.Primative),
		})
	}
	return types
}

func ProcessModel(config *config.Config, model config.Model) Model {
	fields := ProcessFields(config, model)
	indexes := ProcessIndexes(config, model)
	methods := ProcessMethods(config, model, indexes)

	return Model{
		Name:        model.Name,
		Description: model.Description,
		Fields:      fields,
		Indexes:     indexes,
		RepoName:    model.Name + "_" + "repo",
		Methods:     methods,
	}
}

func ProcessMethods(config *config.Config, model config.Model, indexes []Index) []Method {
	methods := make([]Method, 0, len(model.Fields))
	for _, index := range indexes {
		args := make([]Arg, 0, len(index.Fields))
		for _, f := range index.Fields {
			field := model.FindFieldByInternal(f)
			args = append(args, Arg{
				Name: field.Internal,
				Type: field.Type,
			})
		}

		methods = append(methods, Method{
			Name: fmt.Sprintf(
				"%s_%s_by_%s",
				func() string {
					if index.Type == "index" {
						return "List"
					}
					return "Get"
				}(),
				func() string {
					if index.Type == "index" {
						return model.Name + "s"
					}
					return model.Name
				}(),
				strings.Join(index.Fields, "_"),
			),
			Args:       args,
			ReturnType: model.Name,
			ReturnPrefix: func() string {
				if index.Type == "index" {
					return "[]*"
				}
				return "*"
			}(),
		})
	}
	return methods
}

func ProcessFields(config *config.Config, model config.Model) []Field {
	fields := make([]Field, 0, len(model.Fields))
	for _, field := range model.Fields {
		fields = append(fields, Field{
			Name:      field.Name,
			Type:      field.Type,
			Primative: config.TypePrimative(field.Type),
			Internal: func() string {
				if field.Internal == "" {
					return field.Name
				}
				return field.Internal
			}(),
			Relationship: field.Relationship,
		})
	}
	return fields
}

func ProcessIndexes(config *config.Config, model config.Model) []Index {
	indexMap := make(map[string][]string)
	for _, field := range model.Fields {
		for _, index := range field.Indexes {
			_, found := indexMap[index]
			if !found {
				indexMap[index] = make([]string, 0, 10)
			}
			indexMap[index] = append(indexMap[index], field.Internal)
		}
	}

	indexes := make([]Index, 0, len(indexMap))
	for name, fields := range indexMap {
		indexes = append(indexes, Index{
			Name: name,
			Type: Validate(func() string {
				s := strings.Split(name, "_")
				if len(s) == 2 {
					return s[1]
				}
				return name
			}(),
				"index",
				"unique",
				"primary",
			),
			Fields: fields,
		})
	}

	return indexes
}

/*
type Model struct {
	Internal   string
	Deprecated string
	Methods    []Method
}

type Field struct {
	Description  string
	Type         Type
	Deprecated   string
	Relationship string
}

type Type struct {
	Type      string
	Primative string
}

func (t Type) String() string {
	return t.Type
}

type Method struct {
	Name   string
	Args   []Arg
	Return string
}

type Arg struct {
	Name string
	Type Type
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
							Internal: "Id",
							Type: Type{
								Type:      config.TypeMapping("id"),
								Primative: config.TypePrimative("id"),
							},
						},
						Field{
							Name:     CamelCasef("%s_id", m.Name),
							Internal: fmt.Sprintf("%s_id", m.Name),
							Type: Type{
								Type:      config.TypeMapping("id"),
								Primative: config.TypePrimative("id"),
							},
						},
						Field{
							Name:     CamelCasef("%s_id", f.Relationship),
							Internal: fmt.Sprintf("%s_id", f.Relationship),
							Type: Type{
								Type:      config.TypeMapping("id"),
								Primative: config.TypePrimative("id"),
							},
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
				Type: Type{
					Type:      config.TypeMapping(f.Type),
					Primative: config.TypePrimative(f.Type),
				},
				Relationship: f.Relationship,
				Indexes:      f.Indexes,
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
*/

func CamelCasef(f string, args ...interface{}) string {
	return strcase.UpperCamelCase(fmt.Sprintf(f, args...))
}

func LowerCamelCasef(f string, args ...interface{}) string {
	return strcase.LowerCamelCase(fmt.Sprintf(f, args...))
}

func Validate(input string, valid ...string) string {
	for _, v := range valid {
		if v == input {
			return input
		}
	}
	panic(fmt.Sprintf("Invalid Input Not Valid %v", valid))
}
