package generators

import (
	"os"
	"path/filepath"

	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/stoewer/go-strcase"
)

type ResolverTemplate struct {
	Imports []string
	Model   *introspection.Type
}

type ResolverGenerator struct {
	path string
}

func NewResolverGenerator(path string) *ResolverGenerator {
	return &ResolverGenerator{path}
}

func (t *ResolverGenerator) Generate(schema *introspection.Schema) error {
	imports := []string{
		"context",
		"github.com/graph-gophers/graphql-go",
	}

	for _, model := range schema.Types() {
		if model.Kind() != Object || SafeHasPrefix(model.Name(), "__") {
			continue
		}

		modeleName := SafeString(model.Name())

		f, err := os.Create(t.Filename(modeleName))
		if err != nil {
			return err
		}

		if err := ExecuteTemplate(f, "resolver.tmpl", ResolverTemplate{
			Imports: imports,
			Model:   model,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (t *ResolverGenerator) Filename(name string) string {
	if err := os.MkdirAll(filepath.Join(t.path, "resolvers"), os.ModePerm); err != nil {
		panic(err)
	}
	return filepath.Join(t.path, "resolvers", strcase.SnakeCase(name)+".go")
}
