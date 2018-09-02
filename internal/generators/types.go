package generators

import (
	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

const TypesModule = "resolvers"

type TypesTemplate struct {
	Imports []string
	Model   domain.Model
}

type TypesGenerator struct {
	path string
}

func NewTypesGenerator(path string) *TypesGenerator {
	return &TypesGenerator{path}
}

func (t *TypesGenerator) Generate(config *config.Config) error {

	if !config.ShouldGenerate(TypesModule) {
		return nil
	}

	imports := []string{
		"context",
		"github.com/graph-gophers/graphql-go",
	}

	for _, model := range domain.GenerateModels(config) {
		//if err := GenerateGoFile(
		if err := GenerateFile(
			t.Filename(model.Name),
			"types.tmpl",
			TypesTemplate{
				Imports: imports,
				Model:   model,
			}); err != nil {
			return err
		}
	}
	return nil
}

func (t *TypesGenerator) Filename(name string) string {
	return TemplatePath(t.path, "domain", name)
}
