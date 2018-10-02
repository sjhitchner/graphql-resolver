package domain

import (
	"fmt"
	"strings"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

func BuildModels(cfg *config.Config) []Model {
	repoMap := BuildRepoMethods(cfg)

	models := make([]Model, 0, len(cfg.Models))

	for _, m := range cfg.Models {
		if m.Type == config.Link {
			continue
		}

		model := BuildModel(cfg, m)
		model.Repo = Repo{
			Name:    fmt.Sprintf("%s_repo", model.Name),
			Methods: repoMap[model.Name],
		}
		models = append(models, model)
	}

	return models
}

func BuildModel(cfg *config.Config, model config.Model) Model {
	//indexes := make([]Index, 0, len(model.Fields))
	fields, imports := BuildFields(cfg, model)

	return Model{
		Name:        model.Name,
		Plural:      model.Plural,
		Description: model.Description,
		Fields:      fields,
		Imports:     imports,
	}
}

func BuildFields(cfg *config.Config, model config.Model) ([]Field, Imports) {
	fields := make([]Field, 0, len(model.Fields))
	imports := NewImports()

	for _, f := range model.Fields {

		primative, impt := cfg.Primative(f.Type)
		imports.Add(impt)

		internal := f.Internal
		if f.Type == config.ID &&
			!strings.HasSuffix(f.Internal, "id") {
			internal = f.Internal + "_id"
		}

		field := Field{
			Name:         f.Name,
			Type:         f.Type,
			Primative:    primative,
			Internal:     internal,
			ShouldExpose: f.Expose,
		}

		if f.Relationship != nil {
			to := cfg.FindModel(f.Relationship.To)

			field.Relationship = &Relationship{
				To: to.Name,
				Through: func() string {
					if f.Relationship.Through == "" {
						return ""
					}
					through := cfg.FindModel(f.Relationship.Through)
					return through.Plural
				}(),
				Field: f.Relationship.Field,
				Type:  f.Relationship.Type,
			}
		}

		fields = append(fields, field)
	}

	return fields, imports
}

func BuildRepoMethods(cfg *config.Config) map[string][]Method {
	methodMap := make(map[string][]Method, len(cfg.Models))
	for _, model := range cfg.Models {
		buildRepoMethods(cfg, model, methodMap)
	}
	return methodMap
}

func buildRepoMethods(cfg *config.Config, model config.Model, methodMap map[string][]Method) {
	//methods := make([]Method, 0, len(model.Fields))
	for _, idx := range model.Indexes() {
		methodMap[model.Name] = append(
			methodMap[model.Name],
			Method{
				Name: func() string {
					if idx.Type == config.MultiIndex {
						return fmt.Sprintf("list_%s_by_%s", model.Plural, idx.NameWithIds())
					}
					return fmt.Sprintf("get_%s_by_%s", model.Name, idx.NameWithIds())
				}(),
				Args: func() []Arg {
					args := make([]Arg, 0, len(idx.Fields))
					for _, fieldName := range idx.Fields {
						field := model.FindFieldByName(fieldName)
						args = append(args, Arg{
							Parent: model.Internal,
							Name:   fieldName,
							Type:   field.Type,
						})
					}
					return args
				}(),
				Return: Return{
					Type:  model.Name,
					Multi: idx.Type == config.MultiIndex,
				},
			})
	}

	for _, f := range model.Fields {
		if f.Relationship == nil {
			continue
		}

		switch f.Relationship.Type {
		case config.One2Many:
			to := cfg.FindModel(f.Relationship.To)
			methodMap[f.Relationship.To] = append(
				methodMap[f.Relationship.To],
				Method{
					Name: fmt.Sprintf("list_%s_by_%s", to.Plural, f.Relationship.Field),
					Args: []Arg{
						Arg{
							Name:   f.Relationship.Field,
							Parent: f.Relationship.To,
							Type:   config.ID,
						},
					},
					Relationship: &Relationship{
						To:    model.Name,
						Field: f.Relationship.Field,
						Type:  f.Relationship.Type,
					},
					Return: Return{
						Type:  f.Relationship.To,
						Multi: true,
					},
				})

		case config.Many2Many:
			through := cfg.FindModel(f.Relationship.Through)
			methodMap[f.Relationship.To] = append(
				methodMap[f.Relationship.To],
				Method{
					Name: fmt.Sprintf("list_%s_by_%s", through.Plural, f.Relationship.Field),
					Args: []Arg{
						Arg{
							Name:   f.Relationship.Field,
							Parent: f.Relationship.Through,
							Type:   config.ID,
						},
					},
					Relationship: &Relationship{
						To:      model.Name,
						Through: f.Relationship.Through,
						Field:   f.Relationship.Field,
						Type:    f.Relationship.Type,
					},
					Return: Return{
						Type:  f.Relationship.To,
						Multi: true,
					},
				})
		}
	}

	methodMap[model.Name] = append(
		methodMap[model.Name],
		Method{
			Name: "list_" + model.Plural,
			Args: []Arg{
				Arg{
					Name: "first",
					Type: config.Integer,
				},
				Arg{
					Name: "after",
					Type: config.ID,
				},
			},
			Return: Return{
				Type:  model.Name,
				Multi: true,
			},
		})
}

func BuildTypes(cfg *config.Config) ([]Type, Imports) {
	imports := NewImports()
	types := make([]Type, 0, len(cfg.Types))
	for _, t := range cfg.Types {
		typ, impt := cfg.Primative(t.Primative)

		types = append(types, Type{
			Name: t.Name,
			Type: typ,
		})

		imports.Add(impt)
	}
	return types, imports
}

/*
// Return []Models, []Types, []Global Imports
func ProcessConfig(config *config.Config) ([]Model, []Relationship, []Type, []string) {
	models := make([]Model, 0, len(config.Models))
	relationships := make([]Relationship, 0, len(config.Models))
	imports := make([]string, 0, 10)

	for _, m := range config.Models {
		models = append(
			models,
			ProcessModel(config, m),
		)
	}

	for _, m := range models {
		for _, f := range m.Fields {
			if f.Relationship != "" {
				if f.Type == "manytomany" {
					thruModel := FindModel(models, f.Relationship)

					relationships = append(
						relationships,
						ManyToManyRelationship(m, f, thruModel),
					)
				}
			}
		}
	}

	types := ProcessTypes(config, config.Types, &imports)

	return models, relationships, types, imports
}

func ManyToManyRelationship(m Model, f Field, thruModel Model) Relationship {
	return Relationship{
		Name:        m.Name + "_" + f.Name,
		Description: fmt.Sprintf("Linking %s to %s", m.Name, f.Relationship),
		Through:     f.Name,
		Models: []Model{
			m,
			thruModel,
		},
		Fields: []Field{
			Field{
				Name:      "id",
				Type:      "id",
				Primative: "id",
				Internal:  "id",
			},
			Field{
				Name:      m.Name + "_id",
				Type:      "id",
				Primative: "id",
				Internal:  m.Name + "_id",
			},
			Field{
				Name:      f.Relationship + "_id",
				Type:      "id",
				Primative: "id",
				Internal:  f.Relationship + "_id",
			},
		},
		Indexes: nil,
		Methods: []Method{
			Method{
				Name: fmt.Sprintf(
					"%s_%s_%s_by_%s_id",
					"list",
					m.Name,
					f.Name,
					m.Name,
				),
				Args: []Arg{
					Arg{
						Name: m.Name + "_id",
						Type: "id",
					},
				},
				ReturnType:   m.Name,
				ReturnPrefix: "[]*",
			},
			Method{
				Name: fmt.Sprintf(
					"%s_%s_%s_by_%s_id",
					"list",
					m.Name,
					f.Name,
					f.Relationship,
				),
				Args: []Arg{
					Arg{
						Name: f.Relationship + "_id",
						Type: "id",
					},
				},
				ReturnType:   f.Relationship,
				ReturnPrefix: "[]*",
			},
		},
	}
}

func ProcessModel(config *config.Config, model config.Model) Model {
	imports := make([]string, 0, 10)

	indexes := ProcessIndexes(config, model)
	fields := ProcessFields(config, model, &imports)
	methods := ProcessMethods(config, model, indexes)

	return Model{
		Name:        model.Name,
		Description: model.Description,
		Fields:      fields,
		Indexes:     indexes,
		RepoName:    model.Name + "_" + "repo",
		Methods:     methods,
		Imports:     imports,
	}
}

func ProcessMethods(config *config.Config, model config.Model, indexes []Index) []Method {
	methods := make([]Method, 0, len(model.Fields))
	for _, index := range indexes {
		args := make([]Arg, 0, len(index.Fields))
		for _, f := range index.Fields {
			field := model.FindFieldByInternal(f)
			args = append(args, Arg{
				Name: field.Internal,
				Type: field.Type,
			})
		}

		methods = append(methods, Method{
			Name: fmt.Sprintf(
				"%s_%s_by_%s",
				func() string {
					if index.Type == "index" {
						return "list"
					}
					return "get"
				}(),
				func() string {
					if index.Type == "index" {
						return model.Name + "s"
					}
					return model.Name
				}(),
				strings.Join(index.Fields, "_"),
			),
			Args:       args,
			ReturnType: model.Name,
			ReturnPrefix: func() string {
				if index.Type == "index" {
					return "[]*"
				}
				return "*"
			}(),
		})
	}
	return methods
}

func ProcessFields(config *config.Config, model config.Model, imports *[]string) []Field {
	fields := make([]Field, 0, len(model.Fields))
	for _, field := range model.Fields {
		typ, impt := config.TypePrimative(field.Type)

		fields = append(fields, Field{
			Name:      field.Name,
			Type:      field.Type,
			Primative: typ,
			Internal: func() string {
				if field.Internal == "" {
					return field.Name
				}
				return field.Internal
			}(),
			Relationship: field.Relationship,
		})

		if impt != "" {
			*imports = append(*imports, impt)
		}
	}
	return fields
}

func ProcessIndexes(config *config.Config, model config.Model) []Index {
	indexMap := make(map[string][]string)
	for _, field := range model.Fields {
		for _, index := range field.Indexes {
			_, found := indexMap[index]
			if !found {
				indexMap[index] = make([]string, 0, 10)
			}
			indexMap[index] = append(indexMap[index], field.Internal)
		}
	}

	indexes := make([]Index, 0, len(indexMap))
	for name, fields := range indexMap {
		indexes = append(indexes, Index{
			Name: name,
			Type: Validate(func() string {
				s := strings.Split(name, "_")
				if len(s) == 2 {
					return s[1]
				}
				return name
			}(),
				"index",
				"unique",
				"primary",
			),
			Fields: fields,
		})
	}

	return indexes
}
*/
