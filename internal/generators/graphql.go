package generators

import ()

type GraphQLTemplate struct {
	Imports []string
}

type GraphQLGenerator struct {
	path string
}

func NewGraphQLGenerator(path string) *GraphQLGenerator {
	return &GraphQLGenerator{path}
}
