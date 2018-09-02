package generators

import (
	//. "github.com/sjhitchner/graphql-resolver/domain"
	"github.com/sjhitchner/graphql-resolver/internal/config"
)

const ResolverModule = "resolvers"

type ResolverTemplate struct {
	Imports []string
	Model   config.Model
}

type ResolverGenerator struct {
	path string
}

func NewResolverGenerator(path string) *ResolverGenerator {
	return &ResolverGenerator{path}
}

func (t *ResolverGenerator) Generate(config *config.Config) error {

	if !config.ShouldGenerate(ResolverModule) {
		return nil
	}

	imports := []string{
		"context",
		"github.com/graph-gophers/graphql-go",
	}

	for _, model := range config.Models {
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
