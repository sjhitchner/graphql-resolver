package generators

// "encoding/base64"

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
	globalImports := []string{
		"github.com/sjhitchner/graphql-resolver/generated/domain",
		"github.com/sjhitchner/graphql-resolver/generated/interfaces/helpers",
	}

	for _, model := range models {
		model.Imports.Add(globalImports...)

		if err := GenerateGoFile(
			//if err := GenerateFile(
			t.Filename(model.Name),
			"resolver.tmpl",
			ResolverTemplate{
				Imports: model.Imports.AsSlice(),
				Model:   model,
			}); err != nil {
			return errors.Wrapf(err, "Error generating resolver %s", model.Name)
		}
	}

	types, imports := domain.BuildTypes(cfg)
	imports.Add(globalImports[0])
	idType := domain.Type{
		Name: "id",
		Type: "integer",
	}
	for _, t := range types {
		if t.Name == "id" {
			idType = t
		}
	}

	if err := GenerateGoFile(
		//if err := GenerateFile(
		t.Filename("common"),
		"resolver_common.tmpl",
		struct {
			Id      domain.Type
			Imports []string
		}{
			Id:      idType,
			Imports: imports.AsSlice(),
		}); err != nil {
		return errors.Wrapf(err, "Error generating common resolver functions")
	}

	if err := GenerateGoFile(
		//if err := GenerateFile(
		t.Filename("query"),
		"query.tmpl",
		struct {
			Imports []string
			Models  []domain.Model
		}{
			Imports: []string{}, //imports.AsSlice(),
			Models:  models,
		}); err != nil {
		return errors.Wrapf(err, "Error generating query resolver functions")
	}

	return nil
}

func (t *ResolverGenerator) Filename(name string) string {
	return TemplatePath(t.path, "interfaces/resolvers", name)
}
