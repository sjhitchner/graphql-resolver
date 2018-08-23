package graphql

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/sjhitchner/graphql-resolver/example/domain"
	gqllib "github.com/sjhitchner/graphql-resolver/lib/graphql"
)

type RecipeResolver struct {
	recipe *domain.Recipe
}

func (t *RecipeResolver) ID() graphql.ID {
	return graphql.ID(t.recipe.ID)
}

func (t *RecipeResolver) Units() string {
	return t.recipe.Units
}

func (t *RecipeResolver) LyeType() string {
	return t.recipe.LyeType
}

func (t *RecipeResolver) LipidWeight() float64 {
	return t.recipe.LipidWeight
}

func (t *RecipeResolver) WaterLipidRatio() float64 {
	return t.recipe.WaterLipidRatio
}

func (t *RecipeResolver) SuperFatPercentage() float64 {
	return t.recipe.SuperFatPercentage
}

func (t *RecipeResolver) FragranceRatio() float64 {
	return t.recipe.FragranceRatio
}

func (t *RecipeResolver) Lipid() []*RecipeLipidResolver {
	//TODO connection
	return []*RecipeLipidResolver{} //t.recipe.Lipids
}

type RecipeConnectionResolver struct {
	recipes    []*domain.Recipe
	totalCount int
	from       *string
	to         *string
}

func (t *RecipeConnectionResolver) TotalCount() int32 {
	return int32(t.totalCount)
}

func (t *RecipeConnectionResolver) Edges() *[]*RecipeEdgeResolver {
	l := make([]*RecipeEdgeResolver, len(t.recipes))
	for i := range l {
		l[i] = &RecipeEdgeResolver{
			// EncodeCursor
			cursor: gqllib.EncodeCursor(t.recipes[i].ID),
			model:  t.recipes[i],
		}
	}
	return &l
}

func (t *RecipeConnectionResolver) PageInfo() *gqllib.PageInfoResolver {
	return gqllib.NewPageInfoResolver(
		t.from,
		t.to,
		false,
	)
}

type RecipeEdgeResolver struct {
	cursor graphql.ID
	model  *domain.Recipe
}

func (t *RecipeEdgeResolver) Cursor() graphql.ID {
	return t.cursor
}

func (t *RecipeEdgeResolver) Node() *RecipeResolver {
	return &RecipeResolver{t.model}
}
