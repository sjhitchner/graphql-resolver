package generators

import (
	"path/filepath"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type AggregatorTemplate struct {
	Imports []string
	Models  []domain.Model
}

type AggregatorGenerator struct {
	path string
	pkg  string
}

func NewAggregatorGenerator(path string) *AggregatorGenerator {
	return &AggregatorGenerator{
		path: path,
		pkg:  "usecases/aggregator",
	}
}

func (t *AggregatorGenerator) Generate(cfg *config.Config) error {
	models := domain.BuildModels(cfg)

	imports := []string{
		"domainx:" + filepath.Join(cfg.BaseImport, "domain"),
	}

	if err := GenerateGoFile(
		TemplatePath(t.path, t.pkg, "aggregator"),
		"aggregator.tmpl",
		AggregatorTemplate{
			Imports: imports,
			Models:  models,
		}); err != nil {
		return err
	}

	return nil
}
