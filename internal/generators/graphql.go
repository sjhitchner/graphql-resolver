package generators

import (
	//"fmt"
	//"github.com/pkg/errors"
	//"github.com/stoewer/go-strcase"

	//"github.com/sjhitchner/graphql-resolver/internal/config"
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

/*
func (t *GraphQLGenerator) Generate(config *config.Config) error {

	/*
		if !config.ShouldGenerate(QueryModule) {
			return nil
		}
	*

	models, relationships, _, _ := domain.ProcessConfig(config)

	fmt.Println(relationships)

	query := Model{
		Name: "query",
		Methods: []domain.Method{
			Name: "ping",
			ReturnType: "String",
		},
	}

	computedModels := make([]Model, 0, len(models))
	for _, model := range models {
		methods := make([]domain.Method, 0, len(model.Fields))

		for _, field := range model.Fields {
			name := field.Name
			returnType := GraphQLTypeInternal(field.Type, field.Primative)
			if field.Relationship != "" {
				name = field.Name + "_list"
				if field.Type == "manytomany" {
					returnType = fmt.Sprintf("[%s!]", strcase.UpperCamelCase(field.Relationship))
				} else {
					returnType = strcase.UpperCamelCase(field.Relationship)
				}
			}

			methods = append(methods, domain.Method{
				Name:         name,
				ReturnType:   returnType,
				ReturnPrefix: "",
			})
		}

		for _, r := range relationships {
			if r.Models[1].Name == model.Name {
				methods = append(methods, domain.Method{
					Name:         r.Models[0].Name + "_list",
					ReturnType:   fmt.Sprintf("[%s!]", strcase.UpperCamelCase(r.Models[0].Name)),
					ReturnPrefix: "",
				})
			}
		}

		computedModels = append(computedModels, Model{
			Name:    model.Name,
			Methods: methods,
		})

		query.Methods = append(query.Methods, domain.Method{
			Name: model.Name + "_list",
			ReturnType: fmt.Sprintf("[%s!]", strcase.UpperCamelCase(
		})

	}
	computedModels = append(computedModels, query)

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
*/
