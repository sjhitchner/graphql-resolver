package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/example/domain"
	"github.com/sjhitchner/graphql-resolver/lib/db"
)

const SelectRecipe = `
SELECT
	id                
	, units            
	, lye_type         
	, lipid_weight    
	, water_lipid_ratio
	, super_fat_percentage
	, fragrance_ratio   
FROM recipe
`

const SelectRecipeById = SelectRecipe + `WHERE id = $1`
const SelectRecipes = SelectRecipe + `
WHERE id > $1 
ORDER BY id
LIMIT 100
OFFSET $2`

const CreateRecipe = `
INSERT INTO recipe (
	id
	, units            
	, lye_type         
	, lipid_weight    
	, water_lipid_ratio
	, super_fat_percentage
	, fragrance_ratio 
) VALUES (
	 $1
	, $2
	, $3
	, $4
	, $5
	, $6
)
`

const UpdateRecipe = `
UPDATE recipe SET 
	units = $2 
	, lye_type = $3         
	, lipid_weight = $4    
	, water_lipid_ratio = $5
	, super_fat_percentage = $6
	, fragrance_ratio = $7
WHERE id = $1
`

type RecipeDB struct {
	db db.DBHandler
}

func (t *RecipeDB) GetRecipeById(ctx context.Context, id string) (*domain.Recipe, error) {
	var obj domain.Recipe
	err := t.db.GetById(ctx, &obj, SelectRecipeById, id)
	return &obj, errors.Wrapf(err, "error geting recipe '%s'", id)
}

func (t *RecipeDB) ListRecipes(ctx context.Context, first *int32, after *string) ([]*domain.Recipe, error) {
	var list []*domain.Recipe
	err := t.db.Select(ctx, &list, SelectRecipes, after, first)
	return list, errors.Wrapf(err, "err selecting recipes")
}

func (t *RecipeDB) CreateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
	return nil, errors.New("Not implemented")
}

func (t *RecipeDB) UpdateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
	return nil, errors.New("Not implemented")
}
