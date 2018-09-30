package main

import (
	"flag"
	//"fmt"
	//"os"

	//"github.com/graph-gophers/graphql-go"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	//"github.com/sjhitchner/graphql-resolver/internal/domain"
	"github.com/sjhitchner/graphql-resolver/internal/generators"
)

// GraphQL Schema
// Resolvers
// DB
// Aggregator
// Interactor
// Generate Aggregator
// Generate Resolvers
// Library for various functions

// Data Structures

// DB Models
// RepoMethod

var (
	configPath string
	outputPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.StringVar(&outputPath, "path", "generated", "Path to output directory")
}

func main() {
	flag.Parse()

	config, err := config.LoadConfigFile(configPath)
	CheckError(err)

	//b, err := schema.ToJSON()
	//fmt.Println(config)

	generators := []generators.Generator{
		generators.NewTypesGenerator(outputPath),
		generators.NewDomainGenerator(outputPath),
		generators.NewResolverGenerator(outputPath),
		generators.NewSQLGenerator(outputPath),
		generators.NewInteractorGenerator(outputPath),
		generators.NewAggregatorGenerator(outputPath),
		//generators.NewQueryGenerator(outputPath),
		//generators.NewGraphQLGenerator(outputPath),
	}

	for _, generator := range generators {
		err = generator.Generate(config)
		CheckError(err)
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

/*
func (s *GraphQLSuite) Test_GraphQL(c *C) {
	schema, err := graphql.ParseSchema(Schema, nil)
	c.Assert(err, IsNil)

	b, err := schema.ToJSON()
	c.Assert(err, IsNil)
	fmt.Println(string(b))

	inspect := schema.Inspect()

	fmt.Println()
	for _, d := range inspect.Directives() {
		fmt.Println("Name:", d.Name())
		fmt.Println("Description:", *d.Description())
		fmt.Println("Locations:", d.Locations())
		for _, a := range d.Args() {
			fmt.Println("\tName:", a.Name())
			fmt.Println("\tDescription:", *a.Description())
			fmt.Println("\tType:", a.Type().Kind())
			fmt.Println("\tDefault:", a.DefaultValue())
		}
	}

	fmt.Println()
	fmt.Println()
	for _, t := range inspect.Types() {
		fmt.Println("Name:", *t.Name())
		//fmt.Println("Description:", t.Description())
		if t.InputFields() != nil {
			fmt.Println("Input Fields")
			for _, a := range *t.InputFields() {
				fmt.Println("\tName:", a.Name())
				//fmt.Println("\tDescription:", *a.Description())
				fmt.Println("\tType:", a.Type().Kind())
				fmt.Println("\tDefault:", a.DefaultValue())
			}
		}
		if fields := t.Fields(&struct {
			IncludeDeprecated bool
		}{
			IncludeDeprecated: true,
		}); fields != nil {
			fmt.Println("Fields")
			for _, f := range *fields {
				fmt.Println("\tName:", f.Name())
				//fmt.Println("\tDescription:", *f.Description())
				fmt.Println("\tType:", f.Type().Kind())
				fmt.Println("\tDeprecated:", f.IsDeprecated(), f.DeprecationReason())
				for _, a := range f.Args() {
					fmt.Println("\t\tName:", a.Name())
					//fmt.Println("\t\tDescription:", *a.Description())
					fmt.Println("\t\tType:", a.Type().Kind())
					fmt.Println("\t\tDefault:", a.DefaultValue())
				}
			}
		}
	}
}
*/
