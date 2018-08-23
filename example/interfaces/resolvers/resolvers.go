package graphql

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	//"github.com/sjhitchner/graphql-resolver/example/domain"
)

// https://github.com/OscarYuen/go-graphql-starter

type Resolver struct {
}

func (t *Resolver) Lipids(ctx context.Context, args struct {
	prefix string
	First  *int32
	After  *string
}) (*LipidConnectionResolver, error) {
	return &LipidConnectionResolver{}, nil
}

func (t *Resolver) Lipid(ctx context.Context, args struct {
	id graphql.ID
}) (*LipidResolver, error) {
	return &LipidResolver{}, nil
}

func (t *Resolver) Recipes(ctx context.Context, args struct {
	id graphql.ID
}) (*LipidResolver, error) {
	return &LipidResolver{}, nil
}

func (t *Resolver) Recipe(ctx context.Context, args struct {
	id graphql.ID
}) (*RecipeResolver, error) {
	return &RecipeResolver{}, nil
}
