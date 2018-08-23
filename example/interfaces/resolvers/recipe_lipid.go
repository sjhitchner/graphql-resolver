package graphql

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/sjhitchner/graphql-resolver/example/domain"
	gqllib "github.com/sjhitchner/graphql-resolver/lib/graphql"
)

type RecipeLipidResolver struct {
	rl *domain.RecipeLipid
}

func (t *RecipeLipidResolver) ID() graphql.ID {
	return graphql.ID(t.rl.ID)
}

func (t *RecipeLipidResolver) RecipeID() graphql.ID {
	return graphql.ID(t.rl.RecipeID)
}

func (t *RecipeLipidResolver) Name() string {
	return t.rl.Name
}

func (t *RecipeLipidResolver) SAP() float64 {
	return t.rl.SAP
}

func (t *RecipeLipidResolver) Weight() int64 {
	return t.rl.Weight
}

func (t *RecipeLipidResolver) Percentage() float64 {
	return t.rl.Percentage
}
