package domain

import (
	"fmt"
	//"strings"

	"github.com/sjhitchner/graphql-resolver/internal/config"
)

func BuildModels(config *config.Config) []Model {

	models := make([]Model, 0, len(config.Models))

	for _, m := range config.Models {
		models = append(models, BuildModel(config, m))
	}

	return models
}

func BuildModel(cfg *config.Config, model config.Model) Model {

	//indexes := make([]Index, 0, len(model.Fields))
	imports := make([]string, 0, len(model.Fields))

	/*
		for _, f := range model.Fields {
			fields = append(fields, Field{
				Name:
				Type:
				Primative:
				Internal:
			})

			methods = append(methods, Method{
				Name:
				Args: func() []Arg{
					args := make([]Arg, 0, 5)


					args = append(args, Arg{
					})
				}
				Return: Return{
					Type:
					Multi:
				},
			})
		}
	*/

	return Model{
		Name:        model.Name,
		Description: model.Description,
		Fields:      BuildFields(cfg, model, imports),
		Repo: Repo{
			Name:    fmt.Sprintf("%s_repo", model.Name),
			Methods: BuildMethods(cfg, model),
		},
		Imports: imports,
	}
}

func BuildFields(cfg *config.Config, model config.Model, imports []string) []Field {
	fields := make([]Field, 0, len(model.Fields))

	for _, f := range model.Fields {

		if f.Relationship != nil {
			if f.Relationship.Type == config.Many2Many {
				continue
			}
		}

		primative, impt := cfg.Primative(f.Type)
		if impt != "" {
			imports = append(imports, impt)
		}

		fields = append(fields, Field{
			Name:      f.Name,
			Type:      f.Type,
			Primative: primative,
			Internal:  f.Internal,
		})
	}

	return fields
}

func BuildMethods(cfg *config.Config, model config.Model) []Method {
	methods := make([]Method, 0, len(model.Fields))

	for _, idx := range model.Indexes() {

		methods = append(methods, Method{
			Name: func() string {
				if idx.Type == config.MultiIndex {
					return fmt.Sprintf("list_%s_by_%s", model.Name, idx.Name)
				}
				return fmt.Sprintf("get_%s_by_%s", model.Name, idx.Name)
			}(),
			Args: func() []Arg {
				args := make([]Arg, 0, len(idx.Fields))
				for _, fieldName := range idx.Fields {
					field := model.FindFieldByName(fieldName)
					//typ, _ := cfg.Primative(field.Type)
					args = append(args, Arg{
						Name: fieldName,
						Type: field.Type,
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

	methods = append(methods, Method{
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

	return methods
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

func ProcessTypes(config *config.Config, ct []config.Type, imports *[]string) []Type {
	types := make([]Type, 0, len(ct))
	for _, t := range ct {
		typ, impt := config.TypePrimative(t.Primative)

		types = append(types, Type{
			Name: t.Name,
			Type: typ,
		})

		if impt != "" {
			*imports = append(*imports, impt)
		}
	}
	return types
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