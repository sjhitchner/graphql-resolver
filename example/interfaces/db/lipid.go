package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/example/domain"
	"github.com/sjhitchner/graphql-resolver/lib/db"
)

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
FROM lipids
`

const SelectLipidById = SelectLipid + `WHERE id = $1`
const SelectLipids = SelectLipid + `
WHERE id > $1 
ORDER BY id
LIMIT 100
OFFSET $2`

const SearchLipids = SelectLipid + `WHERE name ILIKE $1 LIMIT 10`

type LipidDB struct {
	db db.DBHandler
}

func NewLipidDB(db db.DBHandler) *LipidDB {
	return &LipidDB{db}
}

func (t *LipidDB) GetLipidById(ctx context.Context, id string) (*domain.Lipid, error) {
	var obj domain.Lipid
	err := t.db.GetById(ctx, &obj, SelectRecipeById, id)
	return &obj, errors.Wrapf(err, "error geting recipe '%s'", id)
}

func (t *LipidDB) ListLipids(ctx context.Context, first int32, after string) ([]*domain.Lipid, error) {
	var list []*domain.Lipid
	err := t.db.Select(ctx, &list, SelectLipids, after, first)
	return list, errors.Wrapf(err, "err selecting recipes")
}

func (t *LipidDB) SearchLipid(ctx context.Context, prefix string) ([]*domain.Lipid, error) {
	return nil, errors.New("Not implemented")
}
