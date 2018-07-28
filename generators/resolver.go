package generators

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/graph-gophers/graphql-go/introspection"
	"github.com/stoewer/go-strcase"
)

type ResolverTemplate struct {
	Imports []string
	Models  []*introspection.Type
}

type ResolverGenerator struct {
	path string
}

func NewResolverGenerator(path string) *ResolverGenerator {
	return &ResolverGenerator{path}
}

func (t *ResolverGenerator) Generate(schema *intropsection.Schema) error {
	for _, t := range schema.Types() {
		if t.Kind() != Object || SafeHasPrefix(t.Name(), "__") {
			continue
		}

		typeName := SafeString(t.Name())

		f, err := os.Create(t.Filename(typeName))
		if err != nil {
			return err
		}

		if err := ExecuteTemplate(f, "resolver.tmpl", ResolverTemplate{
			Model: model,
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
