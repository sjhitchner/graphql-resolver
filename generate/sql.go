package generate

import (
	"github.com/graph-gophers/graphql-go/introspection"
)

type SQLTemplate struct {
	Imports []string
}

type SQLResolver struct {
}

func NewSQLResolver() *SQLResolver {

}

func (t *SQLResolver) Generate(schema *introspection.Schema) error {

}
