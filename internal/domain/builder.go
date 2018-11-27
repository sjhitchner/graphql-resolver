package domain

import (
	"fmt"
	"sort"
	"strings"

	//"github.com/pkg/errors"
	"github.com/sjhitchner/graphql-resolver/internal/config"
	//"github.com/stoewer/go-strcase"
)

func BuildModels(cfg *config.Config) []Model {
	repoMap := BuildRepoMethods(cfg)

	models := make([]Model, 0, len(cfg.Models))

	for _, m := range cfg.Models {
		model := BuildModel(cfg, m)

		sort.Sort(ByMethodName(repoMap[model.Name]))

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
		Type:        model.Type,
		Plural:      model.Plural,
		Description: model.Description,
		Fields:      fields,
		Internal:    model.Internal,
		Imports:     imports,
		Mutations:   BuildMutations(cfg, model),
	}
}

func BuildMutations(cfg *config.Config, model config.Model) []Mutation {
	mutations := make([]Mutation, 0, len(model.Mutations))
	for _, m := range model.Mutations {
		mutations = append(mutations, Mutation{
			Name: fmt.Sprintf("%s_%s", m.Name, model.Name),
			Type: m.Type,
			Fields: func() []Field {
				fields := make([]Field, 0, len(model.Fields))

				for _, name := range m.Fields {
					field, _ := BuildField(cfg, model.FindFieldByName(name))
					fields = append(fields, field)
				}
				return fields
			}(),
			Key: m.Key,
		})
	}
	return mutations
}

func BuildFields(cfg *config.Config, model config.Model) ([]Field, Imports) {
	fields := make([]Field, 0, len(model.Fields))
	imports := NewImports()

	for _, field := range model.Fields {
		field, imprt := BuildField(cfg, field)
		imports.Add(imprt)
		fields = append(fields, field)
	}

	return fields, imports
}

func BuildField(cfg *config.Config, f config.Field) (Field, string) {
	primative, imprt := cfg.Primative(f.Type)

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
		ShouldExpose: *f.Expose,
	}

	if f.Relationship != nil {
		to := cfg.FindModelByName(f.Relationship.To)

		field.Relationship = &Relationship{
			To: NameInternal{
				Name:     to.Name,
				Internal: to.Internal,
			},
			Through: func() NameInternal {
				if f.Relationship.Through == "" {
					return NameInternal{}
				}
				through := cfg.FindModelByName(f.Relationship.Through)
				return NameInternal{
					Name:     through.Plural,
					Internal: through.Internal,
				}
			}(),
			Field: f.Relationship.Field,
			Type:  f.Relationship.Type,
		}
	}

	return field, imprt
}

func BuildRepoMethods(cfg *config.Config) map[string][]Method {
	methodMap := make(map[string][]Method, len(cfg.Models))
	for _, model := range cfg.Models {
		buildRepoMethods(cfg, model, methodMap)
	}
	return methodMap
}

func buildRepoMethods(cfg *config.Config, model config.Model, methodMap map[string][]Method) {
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
						field := model.FindFieldByInternal(fieldName)
						args = append(args, Arg{
							Parent: NameInternal{
								Name:     model.Name,
								Internal: model.Internal,
							},
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

	for _, f := range model.Fields {
		if f.Relationship == nil {
			continue
		}

		switch f.Relationship.Type {
		case config.One2Many:
			to := cfg.FindModelByName(f.Relationship.To)
			methodMap[f.Relationship.To] = append(
				methodMap[f.Relationship.To],
				Method{
					Type: "list",
					Name: fmt.Sprintf("list_%s_by_%s", to.Plural, f.Relationship.Field),
					Args: []Arg{
						Arg{
							Name: f.Relationship.Field,
							Parent: NameInternal{
								Name:     to.Name,
								Internal: to.Internal,
							},
							Type: config.ID,
						},
					},
					Relationship: &Relationship{
						To: NameInternal{
							Name:     model.Name,
							Internal: model.Internal,
						},
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
			to := cfg.FindModelByName(f.Relationship.To)
			through := cfg.FindModelByName(f.Relationship.Through)

			methodMap[f.Relationship.To] = append(
				methodMap[f.Relationship.To],
				Method{
					Type: "list",
					Name: fmt.Sprintf("list_%s_by_%s", to.Plural, f.Relationship.Field),
					Args: []Arg{
						Arg{
							Name: f.Relationship.Field,
							Parent: NameInternal{
								Name:     through.Name,
								Internal: through.Internal,
							},
							Type: config.ID,
						},
					},
					Relationship: &Relationship{
						To: NameInternal{
							Name:     model.Name,
							Internal: model.Internal,
						},
						Through: NameInternal{
							Name:     through.Name,
							Internal: through.Internal,
						},
						Field: f.Relationship.Field,
						Type:  f.Relationship.Type,
					},
					Return: Return{
						Type:  f.Relationship.To,
						Multi: true,
					},
				},
			)
		}
	}

	methodMap[model.Name] = append(
		methodMap[model.Name],
		Method{
			Type: "list",
			Name: "list_" + model.Plural,
			Args: []Arg{},
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
