package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/example/domain"
)

// https://github.com/OscarYuen/go-graphql-starter

func Aggregator(ctx context.Context) domain.Aggregator {
	return ctx.Value("agg").(domain.Aggregator)
}

type Resolver struct {
	*Mutation
}

func (t *Resolver) Lipids(ctx context.Context, args struct {
	Prefix *string
	//First  *int32
	//After  *string
}) ([]*LipidResolver, error) {
	//return &LipidConnectionResolver{}, nil

	lipids, err := Aggregator(ctx).ListLipids(ctx, 0, "0")
	if err != nil {
		return nil, errors.Wrapf(err, "error getting lipids")
	}

	out := make([]*LipidResolver, len(lipids))
	for i := range out {
		out[i] = &LipidResolver{lipids[i]}
	}

	return out, nil
}

func (t *Resolver) Lipid(ctx context.Context, args struct {
	ID graphql.ID
}) (*LipidResolver, error) {

	if args.ID == "" {
		return nil, errors.New("No lipid id specified")
	}

	lipid, err := Aggregator(ctx).GetLipidById(ctx, string(args.ID))
	return &LipidResolver{lipid}, err
}

func (t *Resolver) Recipes(ctx context.Context) ([]*RecipeResolver, error) {
	//}) (*RecipeConnectionResolver, error) {
	//return &RecipeConnectionResolver{}, nil
	return []*RecipeResolver{}, nil
}

func (t *Resolver) Recipe(ctx context.Context, args struct {
	ID graphql.ID
}) (*RecipeResolver, error) {

	if args.ID == "" {
		return nil, errors.New("No recipe id specified")
	}

	recipe, err := Aggregator(ctx).GetRecipeById(ctx, string(args.ID))
	return &RecipeResolver{recipe}, err
}
