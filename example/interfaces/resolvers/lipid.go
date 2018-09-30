package resolvers

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/sjhitchner/graphql-resolver/example/domain"
	gqllib "github.com/sjhitchner/graphql-resolver/lib/graphql"
)

type LipidResolver struct {
	lipid *domain.Lipid
}

func (t *LipidResolver) ID() graphql.ID {
	return graphql.ID(t.lipid.ID)
}

func (t *LipidResolver) Name() string {
	return t.lipid.Name
}

func (t *LipidResolver) Description() string {
	return t.lipid.Description
}

func (t *LipidResolver) ScientificName() string {
	return t.lipid.ScientificName
}

func (t *LipidResolver) NaOH() float64 {
	return t.lipid.NaOH
}

func (t *LipidResolver) KOH() float64 {
	return t.lipid.KOH
}

func (t *LipidResolver) Iodine() int32 {
	return int32(t.lipid.Iodine)
}

func (t *LipidResolver) Ins() int32 {
	return int32(t.lipid.Ins)
}

func (t *LipidResolver) Lauric() float64 {
	return t.lipid.Lauric
}

func (t *LipidResolver) Myristic() float64 {
	return t.lipid.Myristic
}

func (t *LipidResolver) Palmitic() float64 {
	return t.lipid.Palmitic
}

func (t *LipidResolver) Stearic() float64 {
	return t.lipid.Stearic
}

func (t *LipidResolver) Ricinoleic() float64 {
	return t.lipid.Ricinoleic
}

func (t *LipidResolver) Oleic() float64 {
	return t.lipid.Oleic
}

func (t *LipidResolver) Linoleic() float64 {
	return t.lipid.Linoleic
}

func (t *LipidResolver) Linolenic() float64 {
	return t.lipid.Linolenic
}

func (t *LipidResolver) Hardness() int32 {
	return int32(t.lipid.Hardness)
}

func (t *LipidResolver) Cleansing() int32 {
	return int32(t.lipid.Cleansing)
}

func (t *LipidResolver) Condition() int32 {
	return int32(t.lipid.Condition)
}

func (t *LipidResolver) Bubbly() int32 {
	return int32(t.lipid.Bubbly)
}

func (t *LipidResolver) Creamy() int32 {
	return int32(t.lipid.Creamy)
}

type LipidConnectionResolver struct {
	lipids     []*domain.Lipid
	totalCount int
	from       *string
	to         *string
}

func (t *LipidConnectionResolver) TotalCount() int32 {
	return int32(t.totalCount)
}

func (t *LipidConnectionResolver) Edges() *[]*LipidEdgeResolver {
	l := make([]*LipidEdgeResolver, len(t.lipids))
	for i := range l {
		l[i] = &LipidEdgeResolver{
			// EncodeCursor
			cursor: graphql.ID(t.lipids[i].ID),
			model:  t.lipids[i],
		}
	}
	return &l
}

func (t *LipidConnectionResolver) PageInfo() *gqllib.PageInfoResolver {
	return gqllib.NewPageInfoResolver(
		t.from,
		t.to,
		false,
	)
}

type LipidEdgeResolver struct {
	cursor graphql.ID
	model  *domain.Lipid
}

func (t *LipidEdgeResolver) Cursor() graphql.ID {
	return t.cursor
}

func (t *LipidEdgeResolver) Node() *LipidResolver {
	return &LipidResolver{t.model}
}
