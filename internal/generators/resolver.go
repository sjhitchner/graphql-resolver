package generators

import (
	//"fmt"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

const ResolverModule = "resolvers"

type ResolverTemplate struct {
	Imports []string
	Model   domain.Model
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

	models, _, imports := domain.ProcessConfig(config)

	for _, model := range models {
		if err := GenerateGoFile(
			//if err := GenerateFile(
			t.Filename(model.Name),
			"resolver.tmpl",
			ResolverTemplate{
				Imports: imports,
				Model:   model,
			}); err != nil {
			return errors.Wrapf(err, "Error generating resolver %s", model.Name)
		}
	}

	return nil
}

func (t *ResolverGenerator) Filename(name string) string {
	return TemplatePath(t.path, "interfaces/resolvers", name)
}
