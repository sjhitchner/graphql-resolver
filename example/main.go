package main

import (
	"io/ioutil"
	"net/http"

	// "github.com/graph-gophers/graphql-go"
	"github.com/sjhitchner/graphql-resolver/example/interfaces/graphql"
	"github.com/sjhitchner/graphql-resolver/example/interfaces/resolvers"
	//"github.com/sjhitchner/graphql-resolver/lib/db/psql"
	"github.com/sjhitchner/graphql-resolver/lib/db"
	"github.com/sjhitchner/graphql-resolver/lib/db/sqlite"
)

// https://github.com/graph-gophers/dataloader

func init() {

}

func main() {

	//dbh, err := psql.NewPSQLDBHandler("localhost", "
	dbh, err := psql.NewSQLiteDBHandler(":memory:")

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

func DbSetup(dbh db.DBHandler) error {

	lipidSchema := `
CREATE TABLE lipid (
	id VARCHAR(32)
	, name VARCHAR(32) 
	, description TEXT
	, scientific_name VARCHAR(100) 
	, naoh FLOAT
	, koh FLOAT
	, iodine INTEGER 
	, ins INTEGER 
	, lauric FLOAT 
	, myristic FLOAT 
	, palmitic FLOAT 
	, stearic  FLOAT
	, ricinoleic FLOAT 
	, oleic FLOAT 
	, linoleic FLOAT 
	, linolenic FLOAT 
	, hardness INTEGER 
	, cleansing INTEGER 
	, condition INTEGER 
	, bubbly INTEGER 
	, creamy INTEGER
);`
	result, err := dbh.DB().Exec(lipidSchema)
	if err != nil {
		return err
	}

	recipeSchema := `
CREATE TABLE recipe (
	id VARCHAR(32) 
	, units VARCHAR(3) 
	, lye_type VARCHAR(10) 
	, lipid_weight FLOAT   
	, water_lipid_ratio FLOAT
	, super_fat_percentage FLOAT
	, fragrance_ratio FLOAT
);`
	result, err := dbh.DB().Exec(recipeSchema)
	if err != nil {
		return err
	}

	return nil
}
