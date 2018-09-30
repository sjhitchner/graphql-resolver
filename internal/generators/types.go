package generators

import (
	//"github.com/stoewer/go-strcase"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type TypesTemplate struct {
	Imports []string
	Types   []domain.Type
}

type TypesGenerator struct {
	path string
}

func NewTypesGenerator(path string) *TypesGenerator {
	return &TypesGenerator{path}
}

func (t *TypesGenerator) Generate(cfg *config.Config) error {

	//_, _, types, imports := domain.ProcessConfig(config)
	types, imports := domain.BuildTypes(cfg)

	if len(types) > 0 {
		if err := GenerateGoFile(
			//if err := GenerateFile(
			t.Filename("types"),
			"types.tmpl",
			TypesTemplate{
				Imports: imports.AsSlice(),
				Types:   types,
			}); err != nil {
			return err
		}
	}
	return nil
}

func (t *TypesGenerator) Filename(name string) string {
	return TemplatePath(t.path, "domain", name)
}
