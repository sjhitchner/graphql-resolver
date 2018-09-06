package generators

import (
	"fmt"
	//"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type SQLTemplate struct {
	Imports []string
	Model   domain.Model
}

type SQLGenerator struct {
	path string
}

func NewSQLGenerator(path string) *SQLGenerator {
	return &SQLGenerator{path}
}

func (t *SQLGenerator) Generate(config *config.Config) error {

	/* TODO
	if !model.ShouldGenerate(SQL) {
		return nil
	}

		sql := model.SQL
		if sql == nil {
			return errors.Errorf("Model '%s' set to generate SQL but no sql block configured", model.Name)
		}

	imports := []string{
		"context",
		"github.com/graph-gophers/graphql-go",
	}
	*/

	models, _, imports := domain.ProcessConfig(config)

	fmt.Println(t.path)

	for _, model := range models {
		if err := GenerateGoFile(
			//if err := GenerateFile(
			t.Filename(model.Name),
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

func (t *SQLGenerator) Filename(name string) string {
	return TemplatePath(t.path, "interfaces/db", name)
}
