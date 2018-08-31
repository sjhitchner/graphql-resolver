package generators

import ()

type GraphQLTemplate struct {
	Imports []string
}

func NewGraphQLGenerator(path string) *GraphQLGenerator {
	return &GraphQLGenerator{path}
}
