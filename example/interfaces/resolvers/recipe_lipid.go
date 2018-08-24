package resolvers

import (
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

func (t *RecipeLipidResolver) Lipid() *LipidResolver {
	return &LipidResolver{}
}

func (t *RecipeLipidResolver) Name() string {
	return t.rl.Name
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
