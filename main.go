package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/generators"
)

var (
	goPath     string
	configPath string
	outputPath string

	dbGenerate string
	dbHost     string
	dbPort     int
	dbName     string
	dbPassword string
)

func init() {
	goPath = os.Getenv("GOPATH")

	flag.StringVar(&configPath, "config", "", "Path to config")
	flag.StringVar(&outputPath, "path", "generated", "Path to output directory")

	flag.StringVar(&dbGenerate, "db-generate", "", "Generate config file from SQL db")

	flag.StringVar(&dbHost, "db-host", "localhost", "DB Host")
	flag.IntVar(&dbPort, "db-port", 5432, "DB Port")
	flag.StringVar(&dbName, "db-name", "", "DB Name")
	flag.StringVar(&dbUser, "db-user", "", "DB User")
	flag.StringVar(&dbPass, "db-pass", "", "DB Password")
}

func main() {
	flag.Parse()

	if dbGenerate != "" {
		if err := GenerateConfigFromDb(); err != nil {
			CheckError(err)
		}
		os.Exit(0)
	}

	config, err := config.LoadConfigFile(configPath)
	CheckError(err)

	pkgPath, err := packagePath(goPath, outputPath)
	CheckError(err)

	config.BaseImport = pkgPath

	generators := []generators.Generator{
		generators.NewDomainGenerator(outputPath),
		generators.NewResolverGenerator(outputPath),
		generators.NewSQLGenerator(outputPath),
		generators.NewInteractorGenerator(outputPath),
		generators.NewAggregatorGenerator(outputPath),
		generators.NewGraphQLGenerator(outputPath),
		generators.NewContextGenerator(outputPath),
		generators.NewMainGenerator(outputPath),
	}

	for _, generator := range generators {
		err = generator.Generate(config)
		CheckError(err)
	}
}

func packagePath(goPath, path string) (string, error) {
	//fullPath, err := filepath.Abs(filepath.Dir(outputPath))
	fullPath, err := filepath.Abs(outputPath)
	if err != nil {
		return "", err
	}

	goPath = filepath.Clean(goPath + "/src")
	fullPath = filepath.Clean(fullPath)

	if !strings.HasPrefix(fullPath, goPath) {
		return "", errors.Errorf("Generation path '%s' not in GOPATH '%s'", fullPath, goPath)
	}

	return filepath.Rel(goPath, fullPath)
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
