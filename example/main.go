package main

import (
	"net/http"

	"github.com/graph-gophers/graphql-go"
)

// https://github.com/graph-gophers/dataloader

func init() {

}

func main() {


	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
}
