package generators

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

//const HttpModule = "resolvers"

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
		pkg:  "interfaces/http",
	}
}

func (t *HttpGenerator) Generate(cfg *config.Config) error {
	imports := []string{
		filepath.Join(cfg.BaseImport, "domain"),
		filepath.Join(cfg.BaseImport, "interfaces/helpers"),
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
