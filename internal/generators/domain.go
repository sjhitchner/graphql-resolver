package generators

import (
	"encoding/json"
	"fmt"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type DomainTemplate struct {
	Imports []string
	Model   domain.Model
}

type TypesTemplate struct {
	Imports []string
	Types   []domain.Type
	Models  []domain.Model
}

type DomainGenerator struct {
	path string
	pkg  string
}

func NewDomainGenerator(path string) *DomainGenerator {
	return &DomainGenerator{
		path: path,
		pkg:  "domain",
	}
}

func (t *DomainGenerator) Generate(cfg *config.Config) error {
	models := domain.BuildModels(cfg)
	types, imports := domain.BuildTypes(cfg)

	for _, model := range models {
		b, _ := json.MarshalIndent(model, "", "  ")
		fmt.Println(string(b))

		if err := GenerateGoFile(
			TemplatePath(t.path, t.pkg, model.Name),
			"domain.tmpl",
			DomainTemplate{
				Imports: model.Imports.AsSlice(),
				Model:   model,
			}); err != nil {
			return err
		}
	}

	if len(types) > 0 {
		if err := GenerateGoFile(
			TemplatePath(t.path, t.pkg, "types"),
			"types.tmpl",
			TypesTemplate{
				Imports: imports.AsSlice(),
				Types:   types,
				Models:  models,
			}); err != nil {
			return err
		}
	}

	return nil
}
