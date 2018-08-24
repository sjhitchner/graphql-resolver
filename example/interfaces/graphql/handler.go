package graphql

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/example/starwars"
	//"github.com/graph-gophers/graphql-go/relay"
)

type Handler struct {
	schema *graphql.Schema
}

func NewHandler(schema string, resolver interface{}) *Handler {
	return &Handler{
		schema: graphql.MustParseSchema(schema, resolver),
	}
}

func (t *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serve GraphQL Request")

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Params", params.Query)

	//ctx := h.Loaders.Attach(r.Context())

	response := t.schema.Exec(
		r.Context(),
		params.Query,
		params.OperationName,
		params.Variables,
	)

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

type GraphiQLHandler struct {
	schema *graphql.Schema
}

func NewGraphiQLHandler(schema string) *GraphiQLHandler {
	return &GraphiQLHandler{
		schema: graphql.MustParseSchema(starwars.Schema, nil),
	}
}

func (t *GraphiQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(GraphiQLPage)
}

var GraphiQLPage = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css" rel="stylesheet" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/es6-promise/4.1.1/es6-promise.auto.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/graphql", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
