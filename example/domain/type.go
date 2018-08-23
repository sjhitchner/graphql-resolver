package domain

type Lipid struct {
	ID             string
	Name           string
	Description    string
	ScientificName string
	NaOH           float64
	KOH            float64
	Iodine         int64
	Ins            int64
	Lauric         int64
	Myristic       int64
	Palmitic       float64
	Stearic        float64
	Ricinoleic     int64
	Oleic          float64
	Linoleic       float64
	Linolenic      float64
	Hardness       int64
	Cleansing      int64
	Condition      int64
	Bubbly         int64
	Creamy         int64
}

type LipidRepo interface {
	GetLipidById(ctx context.Context, id string) (*domain.Lipid, error)
	ListLipids(ctx context.Context, first *int32, after *string) ([]*domain.Lipid, error)
	SearchLipid(ctx context.Context, prefix string) ([]*domain.Lipid, error)
}

type Recipe struct {
	ID                 string
	Units              string
	LyeType            string
	LipidWeight        float64
	WaterLipidRatio    float64
	SuperFatPercentage float64
	FragranceRatio     float64
	Lipids             []RecipeLipid
}

type RecipeRepo interface {
	GetRecipeById(ctx context.Context, id string) (*domain.Recipe, error)
	ListRecipes(ctx context.Context, first *int32, after *string) ([]*domain.Recipe, error)

	CreateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error)
	UpdateRecipe(ctx context.Context, recipe *domain.Recipe) (*domain.Recipe, error)
}

type RecipeLipid struct {
	ID         string
	Name       string
	SAP        float64
	Weight     int64
	Percentage float64
}

const Lipids = []Lipid{
	Lipid{
		ID:             "olive-oil",
		Name:           "Olive Oil",
		Description:    "Olive Oil",
		ScientificName: "",
		NaOH:           0.135,
		Koh:            0.19,
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
		ID:              "coconut-oil",
		Name:            "Coconut Oil, 76 deg",
		Description:     "Coconut Oil, 76 deg",
		Scientific_name: "",
		naoh:            0.183,
		Koh:             0.257,
		Iodine:          10,
		Ins:             258,
		Lauric:          0.48,
		Myristic:        0.19,
		Palmitic:        0.09,
		Stearic:         0.03,
		Ricinoleic:      0,
		Oleic:           0.08,
		Linoleic:        0.02,
		Linolenic:       0,
		Hardness:        0,
		Cleansing:       0,
		Condition:       0,
		Bubbly:          0,
		Creamy:          0,
	},
	Lipid{
		ID:              "palm-oil",
		Name:            "Palm Oil",
		Description:     "Palm Oil",
		Scientific_name: "",
		naoh:            0.142,
		Koh:             0.199,
		Iodine:          53,
		Ins:             145,
		Lauric:          0,
		Myristic:        0.01,
		Palmitic:        0.44,
		Stearic:         0.05,
		Ricinoleic:      0,
		Oleic:           0.39,
		Linoleic:        0.1,
		Linolenic:       0,
		Hardness:        50,
		Cleansing:       1,
		Condition:       49,
		Bubbly:          1,
		Creamy:          49,
	},
	Lipid{
		ID:              "castor-oil",
		Name:            "Castor Oil",
		Description:     "Castor Oil",
		Scientific_name: "",
		naoh:            0.128,
		Koh:             0.18,
		Iodine:          86,
		Ins:             95,
		Lauric:          0,
		Myristic:        0,
		Palmitic:        0,
		Stearic:         0,
		Ricinoleic:      0.9,
		Oleic:           0.04,
		Linoleic:        0.04,
		Linolenic:       0,
		Hardness:        0,
		Cleansing:       0,
		Condition:       98,
		Bubbly:          90,
		Creamy:          90,
	},
	Lipid{
		ID:              "hemp-oil",
		Name:            "Hemp Oil",
		Description:     "Hemp Oil",
		Scientific_name: "",
		naoh:            0.138,
		Koh:             0.193,
		Iodine:          165,
		Ins:             39,
		Lauric:          0,
		Myristic:        0,
		Palmitic:        0.06,
		Stearic:         0.02,
		Ricinoleic:      0,
		Oleic:           0.12,
		Linoleic:        0.57,
		Linolenic:       0.21,
		Hardness:        8,
		Cleansing:       0,
		Condition:       90,
		Bubbly:          0,
		Creamy:          8,
	},
}
