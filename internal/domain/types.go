package domain

import (
	"fmt"

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
}

func GenerateModels(config *config.Config) []Model {

	models := make([]Model, 0, len(config.Models))
	for _, m := range config.Models {

		/*
			methods := make([]Method, 0, len(m.Fields))
			for _, d := range m.Fields {
				methods = append(methods, Method{

				})
			}
		*/

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
			//Methods:     methods,
		})
	}

	return models
}

func CamelCasef(f string, args ...interface{}) string {
	return strcase.UpperCamelCase(fmt.Sprintf(f, args...))
}
