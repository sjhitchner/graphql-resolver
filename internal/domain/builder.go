package domain

import (
	"fmt"
	"strings"

	//"github.com/pkg/errors"
	"github.com/sjhitchner/graphql-resolver/internal/config"
	//"github.com/stoewer/go-strcase"
)

/*
func Parse(cfg *config.Config) ([]*Model, []Type, error) {
	models, err := ParseModels(cfg)

	return models, nil, err
}

func ParseModels(cfg *config.Config) ([]*Model, error) {
	modelMap := make(map[string]*Model)

	for _, m := range cfg.Models {
		model, found := modelMap[m.Name]
		if !found {
			model = BuildModel(cfg, m)
		}

		for _, f := range m.Fields {
			if f.Relationship != nil && f.Relationship.Type == config.Many2Many {
				fmt.Println("here", f.Relationship.Type, f.Relationship.Through)

				mmodel, found := modelMap[f.Relationship.Through]
				if !found {
					mmodel = &Model{
						Name:   f.Relationship.Through,
						Plural: f.Relationship.Through,
					}
				}
				mmodel.Fields = append(mmodel.Fields, Field{
					Name:         f.Relationship.Field,
					Type:         config.ID,
					Primative:    config.ID,
					ShouldExpose: true,
				})
				modelMap[f.Relationship.Through] = mmodel
			}
		}

		modelMap[m.Name] = model
	}

	fmt.Println(modelMap)

	for _, m := range cfg.Models {
		fields, imports := BuildFields(cfg, m)

		model, found := modelMap[m.Name]
		if !found {
			return nil, errors.Errorf("Model not found during fields %s", m.Name)
		}
		model.Fields = fields
		model.Imports = imports
	}

	models := make([]*Model, 0, len(modelMap))
	for _, m := range modelMap {
		models = append(models, m)
	}
	return models, nil
}
*/

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
	fields, imports := BuildFields(cfg, model)

	return Model{
		Name:        model.Name,
		Plural:      model.Plural,
		Description: model.Description,
		Fields:      fields,
		Imports:     imports,
		Mutations: func() []Mutation {
			mutations := make([]Mutation, 0, len(model.Mutations))
			for _, m := range model.Mutations {
				mutations = append(mutations, Mutation{
					Name:  fmt.Sprintf("%s_%s", m, model.Name),
					Type:  m,
					Field: fields,
				})
			}
			return mutations
		}(),
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
				Type: func() string {
					if idx.Type == config.MultiIndex {
						return "list"
					}
					return "get"
				}(),
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
					Type: "list",
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
				},
			)

		case config.Many2Many:
			to := cfg.FindModel(f.Relationship.To)
			methodMap[f.Relationship.To] = append(
				methodMap[f.Relationship.To],
				Method{
					Type: "list",
					Name: fmt.Sprintf("list_%s_by_%s", to.Plural, f.Relationship.Field),
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
				},
				/*
					Method{
						Name: fmt.Sprintf("add_%s_to_%s", to.Name, model.Name),
						Args: []Arg{
							Arg{
								Name:   f.Relationship.Field,
								Parent: f.Relationship.Through,
								Type:   config.ID,
							},
						},
					},
				*/
			)
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

	methodMap[model.Name] = append(
		methodMap[model.Name],
		Method{
			Name: "create_" + model.Name,
			Args: []Arg{
				Arg{
					Name:  model.Name,
					Type:  model.Name,
					Deref: true,
				},
			},
			Return: Return{
				Type:  model.Name,
				Multi: false,
			},
		},
		Method{
			Name: "update_" + model.Name,
			Args: []Arg{
				Arg{
					Name:  model.Name,
					Type:  model.Name,
					Deref: true,
				},
			},
			Return: Return{
				Type:  model.Name,
				Multi: false,
			},
		},
		Method{
			Name: "delete_" + model.Name,
			Args: []Arg{
				Arg{
					Name: "id",
					Type: config.ID,
				},
			},
			Return: Return{
				Type:  model.Name,
				Multi: false,
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
