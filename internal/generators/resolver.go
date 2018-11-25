package generators

import (
	"path/filepath"

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
	pkg  string
}

func NewResolverGenerator(path string) *ResolverGenerator {
	return &ResolverGenerator{
		path: path,
		pkg:  "interfaces/resolvers",
	}
}

func (t *ResolverGenerator) Generate(cfg *config.Config) error {

	if !cfg.ShouldGenerate(ResolverModule) {
		return nil
	}

	models := domain.BuildModels(cfg)
	globalImports := []string{
		"domainx:" + filepath.Join(cfg.BaseImport, "domain"),
	}

	for _, model := range models {
		model.Imports.Add(globalImports...)

		if err := GenerateGoFile(
			TemplatePath(t.path, t.pkg, model.Name),
			"resolver.tmpl",
			ResolverTemplate{
				Imports: model.Imports.AsSlice(),
				Model:   model,
			}); err != nil {
			return errors.Wrapf(err, "Error generating resolver %s", model.Name)
		}

		if err := GenerateGoFile(
			TemplatePath(t.path, t.pkg, model.Name+"_mutation"),
			"mutation.tmpl",
			ResolverTemplate{
				Imports: model.Imports.AsSlice(),
				Model:   model,
			}); err != nil {
			return errors.Wrapf(err, "Error generating mutation resolver %s", model.Name)
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
		TemplatePath(t.path, t.pkg, "common"),
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
		TemplatePath(t.path, t.pkg, "resolver"),
		"query.tmpl",
		struct {
			Imports []string
			Models  []domain.Model
		}{
			Imports: imports.AsSlice(),
			Models:  models,
		}); err != nil {
		return errors.Wrapf(err, "Error generating query resolver functions")
	}

	return nil
}
