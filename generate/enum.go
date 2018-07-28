package generate

import (
	"os"
	"path/filepath"

	"github.com/stoewer/go-strcase"

	. "github.com/sjhitchner/graphql-resolver/domain"
)

type EnumTemplate struct {
	Imports []string
	Enums   []Enum
}

type EnumGenerator struct {
	path string
}

func NewEnumGenerator(path string) *EnumGenerator {
	return &EnumGenerator{path}
}

func (t *EnumGenerator) Generate(schema *ParsedSchema) error {
	f, err := os.Create(t.Filename("enums"))
	if err != nil {
		return err
	}

	return ExecuteTemplate(f, "enum.tmpl", EnumTemplate{
		Enums: schema.Enums,
	})
}

func (t *EnumGenerator) Filename(name string) string {
	if err := os.MkdirAll(filepath.Join(t.path, "resolvers"), os.ModePerm); err != nil {
		panic(err)
	}
	return filepath.Join(t.path, "resolvers", strcase.SnakeCase(name)+".go")
}
