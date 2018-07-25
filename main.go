package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/graph-gophers/graphql-go"

	"github.com/sjhitchner/graphql-resolver/generate"
	//"io"
)

var (
	schemaPath string
	outputPath string
)

func init() {
	flag.StringVar(&schemaPath, "schema", "", "Path to schema")
	flag.StringVar(&outputPath, "path", "", "Path to output directory")
}

func main() {
	flag.Parse()

	schema, err := LoadSchema(schemaPath)
	CheckError(err)

	b, err := schema.ToJSON()
	fmt.Println(string(b))

	resolver := generate.NewResolverGenerator()

	// Generate Aggregator
	// Generate Resolvers
	// Library for various functions
	err = resolver.Generate(schema.Inspect())
	CheckError(err)
}

func LoadSchema(schemaPath string) (*graphql.Schema, error) {
	b, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	schema, err := graphql.ParseSchema(string(b), nil)
	if err != nil {
		return nil, err
	}

	return schema, nil
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
