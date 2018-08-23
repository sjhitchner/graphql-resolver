package graphql

import (
	"github.com/graph-gophers/graphql-go"
)

type PageInfoResolver struct {
	startCursor graphql.ID
	endCursor   graphql.ID
	hasNextPage bool
}

func NewPageInfoResolver(start, end *string, hasNext bool) *PageInfoResolver {
	return &PageInfoResolver{
		startCursor: EncodeCursor(start),
		endCursor:   EncodeCursor(end),
		hasNextPage: false,
	}
}

func (t *PageInfoResolver) StartCursor() *graphql.ID {
	return &t.startCursor
}

func (t *PageInfoResolver) EndCursor() *graphql.ID {
	return &t.endCursor
}

func (t *PageInfoResolver) HasNextPage() bool {
	return t.hasNextPage
}
