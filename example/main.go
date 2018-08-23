package main

import (
	"io/ioutil"
	"net/http"

	// "github.com/graph-gophers/graphql-go"
	"github.com/sjhitchner/graphql-resolver/example/interfaces/graphql"
	"github.com/sjhitchner/graphql-resolver/example/interfaces/resolvers"
)

// https://github.com/graph-gophers/dataloader

func init() {

}

func main() {

	schema, err := ioutil.ReadFile("schema.gql")
	CheckError(err)

	handler := graphql.NewHandler(string(schema), &resolvers.Resolver{})

	http.Handle("/graphql", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
