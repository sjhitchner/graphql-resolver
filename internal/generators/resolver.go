package generators

import (
	"fmt"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

const ResolverModule = "resolvers"

type ResolverTemplate struct {
	Imports      []string
	ResolverName string
	ObjectName   string
	Description  string
	Methods      []ResolverMethod
}

type ResolverMethod struct {
	Name         string
	Description  string
	Template     string
	Relationship string
	Return       string
	Args         []domain.Arg
}

type ResolverGenerator struct {
	path string
}

func NewResolverGenerator(path string) *ResolverGenerator {
	return &ResolverGenerator{path}
}

func (t *ResolverGenerator) Generate(config *config.Config) error {

	if !config.ShouldGenerate(ResolverModule) {
		return nil
	}

	models, _, imports := domain.ProcessConfig(config)

	imports = append(imports, "context")

	for _, model := range models {
		resolverName := fmt.Sprintf("%sResolver", model.Name)
		objectName := model.Name

		methods := make([]ResolverMethod, 0, len(model.Fields))
		for _, field := range model.Fields {

			methods = append(methods, ResolverMethod{
				Name: field.Name,
				Template: func() string {
					if field.Relationship != "" {
						return "relationship"
					}
					if field.Type == "id" {
						return "id"
					}
					return field.Type
				}(),
				Relationship: domain.CamelCasef(field.Relationship),
				Return:       field.Type,
				Args:         []domain.Arg{},
			})
		}

		//if err := GenerateGoFile(
		if err := GenerateFile(
			t.Filename(model.Name),
			"resolver.tmpl",
			ResolverTemplate{
				Imports:      imports,
				ResolverName: resolverName,
				ObjectName:   objectName,
				Description:  model.Description,
				Methods:      methods,
			}); err != nil {
			return err
		}
	}
	return nil
}

func (t *ResolverGenerator) Filename(name string) string {
	return TemplatePath(t.path, "resolvers", name)
}
