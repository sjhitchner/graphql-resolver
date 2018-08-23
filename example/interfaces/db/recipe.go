package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/sjhitchner/graphql-resolver/example/domain"
)

type LipidRepo interface {
	GetLipidById(ctx context.Context, id string) (*domain.Lipid, error)
	ListLipids(ctx context.Context, first *int32, after *string) ([]*domain.Lipid, error)
	SearchLipid(ctx context.Context, prefix string) ([]*domain.Lipid, error)
}

type RecipeRepo interface {
	GetRecipeById(ctx context.Context, id string) (*domain.Recipe, error)
	ListRecipes(ctx context.Context, first *int32, after *string) ([]*domain.Recipe, error)

	CreateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error)
	UpdateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error)
}

const SelectLipid = `
SELECT
	id
	, name
	, description
	, scientific_name
	, naoh
	, koh
	, iodine
	, ins
	, lauric
	, myristic
	, palmitic
	, stearic
	, ricinoleic
	, oleic
	, linoleic
	, linolenic
	, hardness
	, cleansing
	, condition
	, bubbly
	, creamy
FROM lipid
`

const SelectLipidById = SelectLipid + `WHERE id = $1`
const SelectLipids = SelectLipid
const SearchLipids = SelectLipid + `WHERE name ILIKE $1 LIMIT 10`

type LipidDB struct {
}

func (t *LipidDB) GetLipidById(ctx context.Context, id string) (*domain.Lipid, error) {
}

func (t *LipidDB) ListLipids(ctx context.Context, first *int32, after *string) ([]*domain.Lipid, error) {
	var lipids []*domain.Lipid
}

func (t *LipidDB) SearchLipid(ctx context.Context, prefix string) ([]*domain.Lipid, error) {

}

type RecipeDB struct {
}
