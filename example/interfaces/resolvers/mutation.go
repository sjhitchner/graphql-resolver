package resolvers

import (
	"context"

	//"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/example/domain"
)

type Mutation struct {
}

type CreateRecipeInput struct {
	Name               string
	Units              string
	LyeType            string
	LipidWeight        float64
	WaterLipidRatio    float64
	SuperFatPercentage float64
	FragranceRatio     float64
}

type CreateRecipePayloadResolver struct {
	recipe *RecipeResolver
}

func NewCreateRecipePayloadResolver(recipe *domain.Recipe) *CreateRecipePayloadResolver {
	return &CreateRecipePayloadResolver{&RecipeResolver{recipe}}
}

func (t *CreateRecipePayloadResolver) Recipe(ctx context.Context) *RecipeResolver {
	return t.recipe
}

func (t *Resolver) CreateRecipe(ctx context.Context, args struct {
	Input CreateRecipeInput
}) (*CreateRecipePayloadResolver, error) {

	recipe := &domain.Recipe{
		Name:               args.Input.Name,
		Units:              args.Input.Units,
		LyeType:            args.Input.LyeType,
		LipidWeight:        args.Input.LipidWeight,
		WaterLipidRatio:    args.Input.WaterLipidRatio,
		SuperFatPercentage: args.Input.SuperFatPercentage,
		FragranceRatio:     args.Input.FragranceRatio,
	}

	recipe, err := Aggregator(ctx).CreateRecipe(ctx, recipe)
	return NewCreateRecipePayloadResolver(recipe),
		errors.Wrapf(err, "error creating recipe")
}

type UpdateRecipeInput struct {
	ID                 string
	Name               string
	Units              string
	LyeType            string
	LipidWeight        float64
	WaterLipidRatio    float64
	SuperFatPercentage float64
	FragranceRatio     float64
}

type UpdateRecipePayloadResolver struct {
	recipe *RecipeResolver
}

func (t *UpdateRecipePayloadResolver) Recipe(ctx context.Context) *RecipeResolver {
	return t.recipe
}

func NewUpdateRecipePayloadResolver(recipe *domain.Recipe) *UpdateRecipePayloadResolver {
	return &UpdateRecipePayloadResolver{&RecipeResolver{recipe}}
}

func (t *Resolver) UpdateReceipe(ctx context.Context, args struct {
	Input UpdateRecipeInput
}) (*UpdateRecipePayloadResolver, error) {

	recipe := &domain.Recipe{
		ID:                 args.Input.ID,
		Name:               args.Input.Name,
		Units:              args.Input.Units,
		LyeType:            args.Input.LyeType,
		LipidWeight:        args.Input.LipidWeight,
		WaterLipidRatio:    args.Input.WaterLipidRatio,
		SuperFatPercentage: args.Input.SuperFatPercentage,
		FragranceRatio:     args.Input.FragranceRatio,
	}

	recipe, err := Aggregator(ctx).UpdateRecipe(ctx, recipe)
	return NewUpdateRecipePayloadResolver(recipe),
		errors.Wrapf(err, "error updating recipe")
}

type AddRecipeLipidInput struct{}
type AddRecipeLipidPayloadResolver struct{}

func (t *Resolver) AddLipid(ctx context.Context, args struct {
	Input AddRecipeLipidInput
}) (*AddRecipeLipidPayloadResolver, error) {
	return nil, nil
}

type UpdateRecipeLipidInput struct{}
type UpdateRecipeLipidPayloadResolver struct{}

func (t *Resolver) UpdateLipid(ctx context.Context, args struct {
	Input UpdateRecipeLipidInput
}) (*UpdateRecipeLipidPayloadResolver, error) {
	return nil, nil
}

type DeleteRecipeLipidInput struct{}
type DeleteRecipeLipidPayloadResolver struct{}

func (t *Resolver) DeleteLipid(ctx context.Context, args struct {
	Input DeleteRecipeLipidInput
}) (*DeleteRecipeLipidPayloadResolver, error) {
	return nil, nil
}

/*
func (t *Resolver) CreateUser(ctx context.Context, args *struct {
	Email    string
	Password string
}) (*userResolver, error) {
	user := &model.User{
		Email:     args.Email,
		Password:  args.Password,
		IPAddress: *ctx.Value("requester_ip").(*string),
	}

	user, err := ctx.Value("userService").(*service.UserService).CreateUser(user)
	if err != nil {
		ctx.Value("log").(*logging.Logger).Errorf("Graphql error : %v", err)
		return nil, err
	}
	ctx.Value("log").(*logging.Logger).Debugf("Created user : %v", *user)
	return &userResolver{user}, nil
}
*/
