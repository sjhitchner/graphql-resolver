package generate

import (
	"os"

	. "github.com/sjhitchner/graphql-resolver/domain"

	//"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/introspection"
)

type ResolverTemplate struct {
	Imports   []string
	Mutations []string
	Models    []Model
	Enums     []Enum
}

type ResolverGenerator struct {
}

func NewResolverGenerator() *ResolverGenerator {
	return &ResolverGenerator{}
}

func (t *ResolverGenerator) Generate(schema *introspection.Schema) error {

	models := Models(schema)
	enums := Enums(schema)

	return ExecuteTemplate(os.Stdout, "resolver.tmpl", ResolverTemplate{
		Models: models,
		Enums:  enums,
	})
}
