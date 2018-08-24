package db

import (
	"context"
	"fmt"
	"log"

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
	, water_to_lipid_ratio
	, superfat_percentage
	, fragrance_ratio   
FROM recipes
`

const SelectRecipeById = SelectRecipe + `WHERE id = $1`
const SelectRecipes = SelectRecipe + `
WHERE id > $1 
ORDER BY id
LIMIT 100
OFFSET $2`

const CreateRecipe = `
INSERT INTO recipes (
	  name
	, units            
	, lye_type         
	, lipid_weight    
	, water_to_lipid_ratio
	, superfat_percentage
	, fragrance_ratio 
) VALUES (
	  $1
	, $2
	, $3
	, $4
	, $5
	, $6
	, $7
)
`

const UpdateRecipe = `
UPDATE recipes SET 
	name = $2
	, units = $3 
	, lye_type = $4         
	, lipid_weight = $5    
	, water_to_lipid_ratio = $6
	, superfat_percentage = $7
	, fragrance_ratio = $8
WHERE id = $1
`

type RecipeDB struct {
	db db.DBHandler
}

func NewRecipeDB(db db.DBHandler) *RecipeDB {
	return &RecipeDB{db}
}

func (t *RecipeDB) GetRecipeById(ctx context.Context, id string) (*domain.Recipe, error) {
	var obj domain.Recipe
	err := t.db.GetById(ctx, &obj, SelectRecipeById, id)
	return &obj, errors.Wrapf(err, "error geting recipe '%s'", id)
}

func (t *RecipeDB) ListRecipes(ctx context.Context, first int32, after string) ([]*domain.Recipe, error) {
	var list []*domain.Recipe
	err := t.db.Select(ctx, &list, SelectRecipes, after, first)
	return list, errors.Wrapf(err, "err selecting recipes")
}

func (t *RecipeDB) CreateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
	log.Println("CREATE RECIPE", recipe)

	id, err := t.db.InsertWithId(
		ctx,
		CreateRecipe,
		recipe.Name,
		recipe.Units,
		recipe.LyeType,
		recipe.LipidWeight,
		recipe.WaterLipidRatio,
		recipe.SuperFatPercentage,
		recipe.FragranceRatio,
	)
	recipe.ID = fmt.Sprintf("%d", id)
	return recipe, errors.Wrapf(err, "error inserting recipe")
}

func (t *RecipeDB) UpdateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
	log.Println("UPDATE RECIPE", recipe)

	_, err := t.db.Update(
		ctx,
		CreateRecipe,
		recipe.ID,
		recipe.Name,
		recipe.Units,
		recipe.LyeType,
		recipe.LipidWeight,
		recipe.WaterLipidRatio,
		recipe.SuperFatPercentage,
		recipe.FragranceRatio,
	)
	return recipe, errors.Wrapf(err, "error updating recipe")
}
