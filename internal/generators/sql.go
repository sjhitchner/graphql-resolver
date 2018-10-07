package generators

import (
	"path/filepath"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type SQLTemplate struct {
	Imports []string
	Model   domain.Model
}

type SQLGenerator struct {
	path string
	pkg  string
}

func NewSQLGenerator(path string) *SQLGenerator {
	return &SQLGenerator{
		path: path,
		pkg:  "interfaces/db",
	}
}

func (t *SQLGenerator) Generate(cfg *config.Config) error {

	imports := []string{
		filepath.Join(cfg.BaseImport, "domain"),
	}

	/* TODO
	if !model.ShouldGenerate(SQL) {
		return nil
	}

		sql := model.SQL
		if sql == nil {
			return errors.Errorf("Model '%s' set to generate SQL but no sql block configured", model.Name)
		}

	*/

	models := domain.BuildModels(cfg)

	for _, model := range models {
		if err := GenerateGoFile(
			TemplatePath(t.path, t.pkg, model.Name),
			"sql.tmpl",
			SQLTemplate{
				Imports: imports,
				Model:   model,
			}); err != nil {
			return err
		}
	}

	return nil
}
