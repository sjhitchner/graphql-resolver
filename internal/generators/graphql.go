package generators

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stoewer/go-strcase"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type GraphQLTemplate struct {
	Model         domain.Model
	Relationships []Relationship
}

type Relationship struct {
	Name string
	Type string
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

	for _, model := range models {

		rs := make([]Relationship, 0, len(models))
		for _, r := range relationships {
			if r.Models[0].Name == model.Name {
			} else if r.Models[1].Name == model.Name {
				rs = append(rs, Relationship{
					Name: r.Models[0].Name + "s",
					Type: fmt.Sprintf("[%s!]!", strcase.UpperCamelCase(r.Models[1].Name)),
				})
			}
		}

		if err := GenerateFile(
			t.Filename("schema_"+model.Name),
			"schema.tmpl",
			GraphQLTemplate{
				Model:         model,
				Relationships: rs,
			}); err != nil {
			return errors.Wrapf(err, "Error generating graphql schema")
		}
	}

	return nil
}

func (t *GraphQLGenerator) Filename(name string) string {
	return SchemaPath(t.path, "", name)
}
