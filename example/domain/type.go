package domain

import (
	"context"
)

type Lipid struct {
	ID             string  `db:"id"`
	Name           string  `db:"name"`
	Description    string  `db:"description"`
	ScientificName string  `db:"scientific_name"`
	NaOH           float64 `db:"naoh"`
	KOH            float64 `db:"koh"`
	Iodine         int     `db:"iodine"`
	Ins            int     `db:"ins"`
	Lauric         float64 `db:"lauric"`
	Myristic       float64 `db:"myristic"`
	Palmitic       float64 `db:"palmitic"`
	Stearic        float64 `db:"stearic"`
	Ricinoleic     float64 `db:"ricinoleic"`
	Oleic          float64 `db:"oleic"`
	Linoleic       float64 `db:"linoleic"`
	Linolenic      float64 `db:"linolenic"`
	Hardness       int     `db:"hardness"`
	Cleansing      int     `db:"cleansing"`
	Condition      int     `db:"condition"`
	Bubbly         int     `db:"buddly"`
	Creamy         int     `db:"creamy"`
}

type LipidRepo interface {
	GetLipidById(ctx context.Context, id string) (*Lipid, error)
	ListLipids(ctx context.Context, first int32, after string) ([]*Lipid, error)
	SearchLipid(ctx context.Context, prefix string) ([]*Lipid, error)
	//LipidCount(ctx context.Context) (int64, error)
}

type Recipe struct {
	ID                 string  `db:"id"`
	Units              string  `db:"units"`
	LyeType            string  `db:"lye_type"`
	LipidWeight        float64 `db:"lipid_weight"`
	WaterLipidRatio    float64 `db:"water_lipid_ratio"`
	SuperFatPercentage float64 `db:"super_fat_percentage"`
	FragranceRatio     float64 `db:"fragrance_ratio"`
}

type RecipeRepo interface {
	GetRecipeById(ctx context.Context, id string) (*Recipe, error)
	ListRecipes(ctx context.Context, first int32, after string) ([]*Recipe, error)
	//RecipeCount(ctx context.Context) (int64, error)

	CreateRecipe(ctx context.Context, recipe *Recipe) (*Recipe, error)
	UpdateRecipe(ctx context.Context, recipe *Recipe) (*Recipe, error)
}

type RecipeLipid struct {
	ID         string
	RecipeID   string
	LipidID    string
	SAP        float64
	Weight     int
	Percentage float64
}

type Aggregator interface {
	LipidRepo
	RecipeRepo
}

var Lipids = []Lipid{
	Lipid{
		ID:             "olive-oil",
		Name:           "Olive Oil",
		Description:    "Olive Oil",
		ScientificName: "",
		NaOH:           0.135,
		KOH:            0.19,
		Iodine:         85,
		Ins:            105,
		Lauric:         0,
		Myristic:       0,
		Palmitic:       0.14,
		Stearic:        0.03,
		Ricinoleic:     0,
		Oleic:          0.69,
		Linoleic:       0.12,
		Linolenic:      0.01,
		Hardness:       17,
		Cleansing:      0,
		Condition:      82,
		Bubbly:         0,
		Creamy:         17,
	},
	Lipid{
		ID:             "coconut-oil",
		Name:           "Coconut Oil, 76 deg",
		Description:    "Coconut Oil, 76 deg",
		ScientificName: "",
		NaOH:           0.183,
		KOH:            0.257,
		Iodine:         10,
		Ins:            258,
		Lauric:         0.48,
		Myristic:       0.19,
		Palmitic:       0.09,
		Stearic:        0.03,
		Ricinoleic:     0,
		Oleic:          0.08,
		Linoleic:       0.02,
		Linolenic:      0,
		Hardness:       0,
		Cleansing:      0,
		Condition:      0,
		Bubbly:         0,
		Creamy:         0,
	},
	Lipid{
		ID:             "palm-oil",
		Name:           "Palm Oil",
		Description:    "Palm Oil",
		ScientificName: "",
		NaOH:           0.142,
		KOH:            0.199,
		Iodine:         53,
		Ins:            145,
		Lauric:         0,
		Myristic:       0.01,
		Palmitic:       0.44,
		Stearic:        0.05,
		Ricinoleic:     0,
		Oleic:          0.39,
		Linoleic:       0.1,
		Linolenic:      0,
		Hardness:       50,
		Cleansing:      1,
		Condition:      49,
		Bubbly:         1,
		Creamy:         49,
	},
	Lipid{
		ID:             "castor-oil",
		Name:           "Castor Oil",
		Description:    "Castor Oil",
		ScientificName: "",
		NaOH:           0.128,
		KOH:            0.18,
		Iodine:         86,
		Ins:            95,
		Lauric:         0,
		Myristic:       0,
		Palmitic:       0,
		Stearic:        0,
		Ricinoleic:     0.9,
		Oleic:          0.04,
		Linoleic:       0.04,
		Linolenic:      0,
		Hardness:       0,
		Cleansing:      0,
		Condition:      98,
		Bubbly:         90,
		Creamy:         90,
	},
	Lipid{
		ID:             "hemp-oil",
		Name:           "Hemp Oil",
		Description:    "Hemp Oil",
		ScientificName: "",
		NaOH:           0.138,
		KOH:            0.193,
		Iodine:         165,
		Ins:            39,
		Lauric:         0,
		Myristic:       0,
		Palmitic:       0.06,
		Stearic:        0.02,
		Ricinoleic:     0,
		Oleic:          0.12,
		Linoleic:       0.57,
		Linolenic:      0.21,
		Hardness:       8,
		Cleansing:      0,
		Condition:      90,
		Bubbly:         0,
		Creamy:         8,
	},
}
