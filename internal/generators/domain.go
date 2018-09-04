package generators

import (
	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type DomainTemplate struct {
	Imports []string
	Model   domain.Model
}

type DomainGenerator struct {
	path string
}

func NewDomainGenerator(path string) *DomainGenerator {
	return &DomainGenerator{path}
}

func (t *DomainGenerator) Generate(config *config.Config) error {

	imports := []string{
		"context",
		"github.com/graph-gophers/graphql-go",
	}

	for _, model := range domain.GenerateModels(config) {
		if err := GenerateGoFile(
			t.Filename(model.Name),
			"domain.tmpl",
			DomainTemplate{
				Imports: imports,
				Model:   model,
			}); err != nil {
			return err
		}
	}
	return nil
}

func (t *DomainGenerator) Filename(name string) string {
	return TemplatePath(t.path, "domain", name)
}
