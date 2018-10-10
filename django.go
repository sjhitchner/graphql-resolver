package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v2"
)

type Graphs struct {
	Graphs []Graph
}

func (t Graphs) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

type Graph struct {
	False  bool    `json:"False"`
	Name   string  `json:"name"`
	Models []Model `json:"models"`
}

type Model struct {
	Name      string     `json:"name"`
	AppName   string     `json:"app_name"`
	Abstracts []string   `json:"abstracts"`
	Fields    []Field    `json:"fields"`
	Label     string     `json:"label"`
	Relations []Relation `json:"relations"`
}

type Relation struct {
	NeedsNode bool   `json:"needs_node"`
	Target    string `json:"target"`
	Arrows    string `json:"arrows"`
	Label     string `json:"label"`
	TargetApp string `json:"target_app"`
	Type      string `json:"type"`
	Name      string `json:"name"`
}

type Field struct {
	Name       string `json:"name"`
	Abstract   bool   `json:"abstract"`
	Label      string `json:"label"`
	Relation   bool   `json:"relation"`
	blank      bool   `json:"blank"`
	Type       string `json:"type"`
	PrimaryKey bool   `json:"primary_key"`
}

func ConvertDjangoJSONToYAML(djangoJSON string) error {
	f, err := os.Open(djangoJSON)
	if err != nil {
		return err
	}

	var graphs Graphs
	if err := json.NewDecoder(f).Decode(&graphs); err != nil {
		return err
	}

	fmt.Println(graphs)

	models := make([]config.Model, 0, 20)

	for _, graph := range graphs.Graphs {
		if strings.Contains(graph.Name, "django") {
			continue
		}

		baseRelations := make(map[string][]Relation)
		for _, m := range graph.Models {
			if len(m.Abstracts) > 0 {
				continue
			}
			baseRelations[m.Name] = m.Relations
		}

		for _, m := range graph.Models {
			if len(m.Abstracts) == 0 {
				continue
			}
			fmt.Println("MODEL", m.Name)

			if len(m.Abstracts) > 0 {
				for _, a := range m.Abstracts {
					m.Relations = append(m.Relations, baseRelations[a]...)
				}
			}

			fields := Fields(m)

			model := config.Model{
				Name:   strcase.SnakeCase(m.Name),
				Fields: fields,
			}

			models = append(models, model)
		}
	}

	c := config.Config{
		Models: models,
	}
	if err := yaml.NewEncoder(os.Stdout).Encode(c); err != nil {
		return err
	}

	return nil
}

var FieldLookup = map[string]string{
	"AutoField":     "id",
	"BooleanField":  "boolean",
	"DateTimeField": "timestamp",
	"CharField":     "string",
	"TextField":     "string",
	"EmailField":    "email",
	"IntegerField":  "integer",
	"FloatField":    "float",
}

//"OneToOneField (id)": "",
//"ForeignKey (id)":    "",

func Fields(m Model) []config.Field {
	fields := make([]config.Field, 0, 20)

	for _, f := range m.Fields {

		field := config.Field{
			Name: f.Name,
		}

		t, found := FieldLookup[f.Type]
		if !found {
			switch {
			case strings.HasPrefix(f.Type, "OneToOneField"):
				r := FindRelation(m, f.Name)
				field.Relationship = &config.Relationship{
					To:   strcase.SnakeCase(r.Target),
					Type: config.One2One,
				}

			case strings.HasPrefix(f.Type, "ManyToManyField"):
				r := FindRelation(m, f.Name)
				field.Relationship = &config.Relationship{
					To:      strcase.SnakeCase(r.Target),
					Field:   "id",
					Through: "",
					Type:    config.Many2Many,
				}

			case strings.HasPrefix(f.Type, "ForeignKey"):
				r := FindRelation(m, f.Name)
				field.Relationship = &config.Relationship{
					To:    strcase.SnakeCase(r.Target),
					Field: "id",
					Type:  config.One2Many,
				}
			default:
			}
		} else {
			field.Type = t
		}

		fields = append(fields, field)
	}

	return fields
}

func FindRelation(m Model, name string) Relation {
	for _, r := range m.Relations {
		if r.Name == name {
			return r
		}
	}
	panic("no relation")
}

func Contains(list []string, str string) bool {
	for _, l := range list {
		if l == str {
			return true
		}
	}
	return false
}
