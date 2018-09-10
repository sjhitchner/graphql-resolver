package generators

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type GraphQLTemplate struct {
	Models []Model
}

type Model struct {
	Name    string
	Methods []domain.Method
}

type GraphQLGenerator struct {
	path string
}

func NewGraphQLGenerator(path string) *GraphQLGenerator {
	return &GraphQLGenerator{path}
}

func (t *GraphQLGenerator) Generate(config *config.Config) error {

	/*
		if !config.ShouldGenerate(QueryModule) {
			return nil
		}
	*/

	models, relationships, _, _ := domain.ProcessConfig(config)

	fmt.Println(relationships)

	computedModels := make([]Model, 0, len(models))
	for _, model := range models {
		methods := make([]domain.Method, 0, len(model.Fields))

		for _, field := range model.Fields {
			returnType := GraphQLTypeInternal(field.Type, field.Primative)
			if field.Relationship != "" {
				if field.Type == "manytomany" {
					returnType = fmt.Sprintf("[%s!]", strcase.UpperCamelCase(field.Relationship))
				} else {
					returnType = strcase.UpperCamelCase(field.Relationship)
				}
			}

			methods = append(methods, domain.Method{
				Name:         field.Name,
				ReturnType:   returnType,
				ReturnPrefix: "",
			})
		}

		for _, r := range relationships {
			if r.Models[1].Name == model.Name {
				name := strcase.UpperCamelCase(r.Models[0].Name + "s")
				returnType := fmt.Sprintf("[%s!]", strcase.UpperCamelCase(r.Models[0].Name))

				methods = append(methods, domain.Method{
					Name:         name,
					ReturnType:   returnType,
					ReturnPrefix: "",
				})
			}
		}

		computedModels = append(computedModels, Model{
			Name:    model.Name,
			Methods: methods,
		})
	}

	if err := GenerateFile(
		t.Filename("schema"),
		"schema.tmpl",
		GraphQLTemplate{
			Models: computedModels,
		}); err != nil {
		return errors.Wrapf(err, "Error generating graphql schema")
	}

	return nil
}

func (t *GraphQLGenerator) Filename(name string) string {
	return SchemaPath(t.path, "", name)
}
