package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/example/domain"
	"github.com/sjhitchner/graphql-resolver/lib/db"
)

const SelectRecipeLipid = `
SELECT
	 id
	 , recipe_id
	 , lipid_id
	 , sap
	 , weight
	 , percentage
FROM recipe_to_lipid_link
`

const SelectRecipeLipidById = SelectRecipeLipid + `
WHERE id = $1
`

const SelectRecipeLipidsByRecipeId = SelectRecipeLipid + `
WHERE recipe_id = $1
`

type RecipeLipidDB struct {
	db db.DBHandler
}

func NewRecipeLipidDB(db db.DBHandler) *RecipeLipidDB {
	return &RecipeLipidDB{db}
}

func (t *RecipeLipidDB) GetRecipeLipidById(ctx context.Context, id string) (*domain.RecipeLipid, error) {
	var obj domain.RecipeLipid
	err := t.db.GetById(ctx, &obj, SelectRecipeLipidById, id)
	return &obj, errors.Wrapf(err, "error geting recipe lipid '%s'", id)
}

func (t *RecipeLipidDB) ListRecipeLipidsByRecipeId(ctx context.Context, recipeID string) ([]*domain.RecipeLipid, error) {
	var list []*domain.RecipeLipid
	err := t.db.Select(ctx, &list, SelectRecipeLipidsByRecipeId, recipeID)
	return list, errors.Wrapf(err, "err selecting recipes")
}
