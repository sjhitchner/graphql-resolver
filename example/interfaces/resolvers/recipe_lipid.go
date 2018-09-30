package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"github.com/sjhitchner/graphql-resolver/example/domain"
	//gqllib "github.com/sjhitchner/graphql-resolver/lib/graphql"
)

type RecipeLipidResolver struct {
	agg domain.Aggregator
	rl  *domain.RecipeLipid
}

func (t *RecipeLipidResolver) ID() graphql.ID {
	return graphql.ID(t.rl.ID)
}

func (t *RecipeLipidResolver) Recipe(ctx context.Context) (*RecipeResolver, error) {
	obj, err := Aggregator(ctx).GetRecipeById(ctx, t.rl.RecipeID)
	return &RecipeResolver{obj}, err
}

func (t *RecipeLipidResolver) Lipid(ctx context.Context) (*LipidResolver, error) {
	obj, err := Aggregator(ctx).GetLipidById(ctx, t.rl.LipidID)
	return &LipidResolver{obj}, err
}

func (t *RecipeLipidResolver) SAP() float64 {
	return t.rl.SAP
}

func (t *RecipeLipidResolver) Weight() int32 {
	return int32(t.rl.Weight)
}

func (t *RecipeLipidResolver) Percentage() float64 {
	return t.rl.Percentage
}
