package generators

import (
	. "github.com/sjhitchner/graphql-resolver/domain"
)

type SQLTemplate struct {
	Imports []string
	Model   Model
}

type SQLGenerator struct {
	path string
}

func NewSQLGenerator(path string) *SQLGenerator {
	return &SQLGenerator{path}
}

func (t *SQLGenerator) Generate(models ...Model) error {
	imports := []string{
		"context",
		"github.com/graph-gophers/graphql-go",
	}

	for _, model := range models {
		//if err := GenerateGoFile(
		if err := GenerateFile(
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
