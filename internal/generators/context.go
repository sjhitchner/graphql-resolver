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
	pkg  string
}

func NewContextGenerator(path string) *ContextGenerator {
	return &ContextGenerator{
		path: path,
		pkg:  "interfaces/helpers",
	}
}

func (t *ContextGenerator) Generate(cfg *config.Config) error {
	imports := []string{
		"github.com/sjhitchner/graphql-resolver/generated/domain",
	}

	if err := GenerateGoFile(
		TemplatePath(t.path, t.pkg, "context"),
		"context.tmpl",
		ContextTemplate{
			Imports: imports,
		}); err != nil {
		return errors.Wrapf(err, "Error generating context helper functions")
	}

	return nil
}
