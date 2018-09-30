package generators

import (
	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type InteractorTemplate struct {
	Imports []string
	Models  []domain.Model
}

type InteractorGenerator struct {
	path string
}

func NewInteractorGenerator(path string) *InteractorGenerator {
	return &InteractorGenerator{path}
}

func (t *InteractorGenerator) Generate(cfg *config.Config) error {
	models := domain.BuildModels(cfg)

	imports := []string{
		"github.com/sjhitchner/graphql-resolver/generated/domain",
	}

	if err := GenerateGoFile(
		//	if err := GenerateFile(
		t.Filename("interactor"),
		"interactor.tmpl",
		InteractorTemplate{
			Imports: imports,
			Models:  models,
		}); err != nil {
		return err
	}
	return nil
}

func (t *InteractorGenerator) Filename(name string) string {
	return TemplatePath(t.path, "usecases/interactor", name)
}
