package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func apiID(ctx context.Context) string {
	opCtxt := graphql.GetOperationContext(ctx)
	return opCtxt.Variables["id"].(string)
}
