package generators

import (
	"path/filepath"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type InteractorTemplate struct {
	Imports []string
	Models  []domain.Model
}

type InteractorGenerator struct {
	path string
	pkg  string
}

func NewInteractorGenerator(path string) *InteractorGenerator {
	return &InteractorGenerator{
		path: path,
		pkg:  "usecases/interactor",
	}
}

func (t *InteractorGenerator) Generate(cfg *config.Config) error {
	models := domain.BuildModels(cfg)

	imports := []string{
		filepath.Join(cfg.BaseImport, "domain"),
	}

	if err := GenerateGoFile(
		TemplatePath(t.path, t.pkg, "interactor"),
		"interactor.tmpl",
		InteractorTemplate{
			Imports: imports,
			Models:  models,
		}); err != nil {
		return err
	}
	return nil
}
