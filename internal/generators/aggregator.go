package generators

import (
	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type AggregatorTemplate struct {
	Imports []string
	Models  []domain.Model
}

type AggregatorGenerator struct {
	path string
}

func NewAggregatorGenerator(path string) *AggregatorGenerator {
	return &AggregatorGenerator{path}
}

func (t *AggregatorGenerator) Generate(cfg *config.Config) error {
	models := domain.BuildModels(cfg)

	imports := []string{
		"github.com/sjhitchner/graphql-resolver/generated/domain",
	}

	if err := GenerateGoFile(
		//	if err := GenerateFile(
		t.Filename("aggregator"),
		"aggregator.tmpl",
		AggregatorTemplate{
			Imports: imports,
			Models:  models,
		}); err != nil {
		return err
	}
	return nil
}

func (t *AggregatorGenerator) Filename(name string) string {
	return TemplatePath(t.path, "usecases/aggregator", name)
}
