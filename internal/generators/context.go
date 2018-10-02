package generators

import (
	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

const ContextModule = "resolvers"

type ContextTemplate struct {
	Imports []string
}

type ContextGenerator struct {
	path string
}

func NewContextGenerator(path string) *ContextGenerator {
	return &ContextGenerator{path}
}

func (t *ContextGenerator) Generate(cfg *config.Config) error {
	imports := []string{
		"github.com/sjhitchner/graphql-resolver/generated/domain",
	}

	if err := GenerateGoFile(
		//if err := GenerateFile(
		t.Filename("context"),
		"context.tmpl",
		ContextTemplate{
			Imports: imports,
		}); err != nil {
		return errors.Wrapf(err, "Error generating context helper functions")
	}

	return nil
}

func (t *ContextGenerator) Filename(name string) string {
	return TemplatePath(t.path, "interfaces/helpers", name)
}
