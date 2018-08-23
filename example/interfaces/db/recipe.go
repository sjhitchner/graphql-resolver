package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/sjhitchner/graphql-resolver/example/domain"
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
const SelectRecipes = SelectLipid

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
}

func (t *RecipeDB) GetRecipeById(ctx context.Context, id string) (*domain.Recipe, error) {
}

func (t *RecipeDB) ListRecipes(ctx context.Context, first *int32, after *string) ([]*domain.Recipe, error) {
}

func (t *RecipeDB) CreateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
}

func (t *RecipeDB) UpdateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error) {
}
