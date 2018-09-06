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

func (t *TypesGenerator) Generate(config *config.Config) error {

	imports := []string{}

	_, types := domain.ProcessConfig(config)
	if len(types) > 0 {
		if err := GenerateGoFile(
			//if err := GenerateFile(
			t.Filename("common"),
			"types.tmpl",
			TypesTemplate{
				Imports: imports,
				Types:   types,
			}); err != nil {
			return err
		}
	}
	return nil
	/*
		args := make([]domain.Arg, 0, len(config.Types))
		for _, typ := range config.Types {
			args = append(args, domain.Arg{
				Name: strcase.UpperCamelCase(typ.Name),
				Type: domain.Type{
					Type:      "",
					Primative: config.TypePrimative(typ.Primative),
				},
			})
		}

		if len(args) > 0 {
			if err := GenerateGoFile(
				//if err := GenerateFile(
				t.Filename("common"),
				"types.tmpl",
				TypesTemplate{
					Imports: imports,
					Types:   args,
				}); err != nil {
				return err
			}
		}
		return nil
	*/
}

func (t *TypesGenerator) Filename(name string) string {
	return TemplatePath(t.path, "domain", name)
}
