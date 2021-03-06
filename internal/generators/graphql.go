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

func (t *GraphQLGenerator) Generate(cfg *config.Config) error {
	/*
		if !config.ShouldGenerate(QueryModule) {
			return nil
		}
	*/

	models := domain.BuildModels(cfg)

	if err := GenerateFile(
		SchemaPath(t.path, "", "schema"),
		"schema.tmpl",
		GraphQLTemplate{
			Models: models,
		}); err != nil {
		return errors.Wrapf(err, "Error generating graphql schema")
	}

	if err := GenerateFile(
		ReactPath(t.path, "react", "mutations"),
		"react_mutations.tmpl",
		GraphQLTemplate{
			Models: models,
		}); err != nil {
		return errors.Wrapf(err, "Error generating react mutations")
	}

	return nil
}
