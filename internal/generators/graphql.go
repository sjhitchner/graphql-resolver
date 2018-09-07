package generators

import (
	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type GraphQLTemplate struct {
	Models []domain.Model
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

	models, _, _ := domain.ProcessConfig(config)

	if err := GenerateFile(
		t.Filename("schema"),
		"schema.tmpl",
		GraphQLTemplate{
			Models: models,
		}); err != nil {
		return errors.Wrapf(err, "Error generating graphql schema")
	}

	return nil
}

func (t *GraphQLGenerator) Filename(name string) string {
	return TemplatePath(t.path, "", name)
}
