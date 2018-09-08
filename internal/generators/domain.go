package generators

import (
	"encoding/json"
	"fmt"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/internal/domain"
)

type DomainTemplate struct {
	Imports []string
	Model   domain.Model
}

type DomainGenerator struct {
	path string
}

func NewDomainGenerator(path string) *DomainGenerator {
	return &DomainGenerator{path}
}

func (t *DomainGenerator) Generate(config *config.Config) error {

	models, _, imports := domain.ProcessConfig(config)
	imports = append(imports, "context")

	for _, model := range models {
		b, _ := json.MarshalIndent(model, "", "  ")
		fmt.Println(string(b))

		if err := GenerateGoFile(
			//if err := GenerateFile(
			t.Filename(model.Name),
			"domain.tmpl",
			DomainTemplate{
				Imports: model.Imports,
				Model:   model,
			}); err != nil {
			return err
		}
	}
	return nil
}

func (t *DomainGenerator) Filename(name string) string {
	return TemplatePath(t.path, "domain", name)
}