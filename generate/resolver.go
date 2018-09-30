package generate

import (
	"os"
	"path/filepath"

	"github.com/stoewer/go-strcase"

	. "github.com/sjhitchner/graphql-resolver/domain"
)

type ResolverTemplate struct {
	Imports []string
	Model   Model
}

type ResolverGenerator struct {
	path string
}

func NewResolverGenerator(path string) *ResolverGenerator {
	return &ResolverGenerator{path}
}

func (t *ResolverGenerator) Generate(schema *ParsedSchema) error {
	for _, model := range schema.Models {
		f, err := os.Create(t.Filename(model.Name))
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
