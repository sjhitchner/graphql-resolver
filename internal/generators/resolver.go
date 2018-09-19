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

func (t *ResolverGenerator) Generate(cfg *config.Config) error {

	if !cfg.ShouldGenerate(ResolverModule) {
		return nil
	}

	models := domain.BuildModels(cfg)
	//models, _, _, imports := domain.ProcessConfig(config)

	for _, model := range models {
		if err := GenerateGoFile(
			//if err := GenerateFile(
			t.Filename(model.Name),
			"resolver.tmpl",
			ResolverTemplate{
				Imports: []string{}, //models.Imports,
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
