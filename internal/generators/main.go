package generators

import (
	"github.com/pkg/errors"
	"path/filepath"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type MainTemplate struct {
	Models  []domain.Model
	Imports []string
}

type MainGenerator struct {
	path string
}

func NewMainGenerator(path string) *MainGenerator {
	return &MainGenerator{path}
}

func (t *MainGenerator) Generate(cfg *config.Config) error {

	models := domain.BuildModels(cfg)

	imports := []string{
		filepath.Join(cfg.BaseImport, "interfaces/db"),
		filepath.Join(cfg.BaseImport, "interfaces/helpers"),
		filepath.Join(cfg.BaseImport, "interfaces/resolvers"),
		filepath.Join(cfg.BaseImport, "usecases/aggregator"),
		filepath.Join(cfg.BaseImport, "usecases/interactor"),
		"libdb:github.com/sjhitchner/graphql-resolver/lib/db",
		"github.com/sjhitchner/graphql-resolver/lib/graphql",
		"libsql:github.com/sjhitchner/graphql-resolver/lib/db/sqlite",
	}

	if err := GenerateGoFile(
		TemplatePath(t.path, "", "main"),
		"main.tmpl",
		MainTemplate{
			Models:  models,
			Imports: imports,
		}); err != nil {
		return errors.Wrapf(err, "Error generating main")
	}

	return nil
}
