package generators

import (
	//"fmt"

	//"github.com/pkg/errors"

	//"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type QueryTemplate struct {
	Imports []string
	Types   []domain.Type
	Models  []domain.Model
}

type QueryGenerator struct {
	path string
}

func NewQueryGenerator(path string) *QueryGenerator {
	return &QueryGenerator{path}
}

/*
func (t *QueryGenerator) Generate(config *config.Config) error {

	/*
		if !config.ShouldGenerate(QueryModule) {
			return nil
		}
	*

	models, _, types, imports := domain.ProcessConfig(config)

	//if err := GenerateGoFile(
	if err := GenerateFile(
		t.Filename("query"),
		"query.tmpl",
		QueryTemplate{
			Imports: imports,
			Types:   types,
			Models:  models,
		}); err != nil {
		return errors.Wrapf(err, "Error generating query resolver")
	}

	return nil
}
*/

func (t *QueryGenerator) Filename(name string) string {
	return TemplatePath(t.path, "interfaces/resolvers", name)
}
