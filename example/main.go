package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "github.com/graph-gophers/graphql-go"
	"github.com/sjhitchner/graphql-resolver/example/domain"
	"github.com/sjhitchner/graphql-resolver/example/interfaces/db"
	"github.com/sjhitchner/graphql-resolver/example/interfaces/graphql"
	"github.com/sjhitchner/graphql-resolver/example/interfaces/resolvers"
	libdb "github.com/sjhitchner/graphql-resolver/lib/db"
	//libsql "github.com/sjhitchner/graphql-resolver/lib/db/psql"
	libsql "github.com/sjhitchner/graphql-resolver/lib/db/sqlite"
)

// https://github.com/graph-gophers/dataloader
var (
	initializeDB bool

	sqlitePath string
	schemaPath string
)

func init() {
	flag.BoolVar(&initializeDB, "initialize", false, "Initialize the DB")
	flag.StringVar(&sqlitePath, "sqlite", ":memory:", "Path to sqlite db")
	flag.StringVar(&schemaPath, "schema", "schema.gql", "Path to graphql schema")
}

func main() {
	flag.Parse()

	//dbh, err := libsql.NewPSQLDBHandler("localhost", "
	dbh, err := libsql.NewSQLiteDBHandler(sqlitePath)
	CheckError(err)

	if initializeDB {
		CheckError(InitializeDBSchema(dbh))
		os.Exit(0)
	}

	schema, err := ioutil.ReadFile(schemaPath)
	CheckError(err)

	aggregator := struct {
		domain.LipidRepo
		domain.RecipeRepo
		domain.RecipeLipidRepo
	}{
		db.NewLipidDB(dbh),
		db.NewRecipeDB(dbh),
		db.NewRecipeLipidDB(dbh),
	}

	handler := graphql.NewHandler(string(schema), &resolvers.Resolver{})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "agg", aggregator)

	http.Handle("/graphql", graphql.WrapContext(ctx, handler))
	http.Handle("/", graphql.NewGraphiQLHandler(string(schema)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func InitializeDBSchema(dbh libdb.DBHandler) error {

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
	if _, err := dbh.DB().Exec(lipidSchema); err != nil {
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
	if _, err := dbh.DB().Exec(recipeSchema); err != nil {
		return err
	}

	recipeLipidSchema := `
CREATE TABLE recipe_lipid (
	id VARCHAR(32) 
	, recipe_id VARCHAR(32) 
	, lipid_id VARCHAR(32) 
	, sap FLOAT   
	, weight INTEGER
	, percentage FLOAT
);`
	if _, err := dbh.DB().Exec(recipeLipidSchema); err != nil {
		return err
	}

	return nil
}
