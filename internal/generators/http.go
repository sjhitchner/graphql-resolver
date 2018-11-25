package generators

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

type HttpTemplate struct {
	Imports []string
}

type HttpGenerator struct {
	path string
	pkg  string
}

func NewHttpGenerator(path string) *HttpGenerator {
	return &HttpGenerator{
		path: path,
		pkg:  "interfaces/helpers",
	}
}

func (t *HttpGenerator) Generate(cfg *config.Config) error {
	imports := []string{
		"domainx:" + filepath.Join(cfg.BaseImport, "domain"),
	}

	if err := GenerateGoFile(
		TemplatePath(t.path, t.pkg, "auth"),
		"http_auth.tmpl",
		HttpTemplate{
			Imports: imports,
		}); err != nil {
		return errors.Wrapf(err, "Error generating http functions")
	}

	return nil
}
