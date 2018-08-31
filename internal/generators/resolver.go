package generators

import (
	. "github.com/sjhitchner/graphql-resolver/domain"
)

type ResolverTemplate struct {
	Imports []string
	Model   Model
}

type ResolverGenerator struct {
	path string
}

func NewResolverGenerator(path string) *ResolverGenerator {
	return &ResolverGenerator{path}
}

func (t *ResolverGenerator) Generate(models ...Model) error {
	imports := []string{
		"context",
		"github.com/graph-gophers/graphql-go",
	}

	for _, model := range models {
		//if err := GenerateGoFile(
		if err := GenerateFile(
			t.Filename(model.Name),
			"resolver.tmpl",
			ResolverTemplate{
				Imports: imports,
				Model:   model,
			}); err != nil {
			return err
		}
	}
	return nil
}

func (t *ResolverGenerator) Filename(name string) string {
	return TemplatePath(t.path, "resolvers", name)
}
