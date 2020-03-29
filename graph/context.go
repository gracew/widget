package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
const apiIDCtxKey = "apiID"

func apiID(ctx context.Context) *string {
	opCtxt := graphql.GetOperationContext(ctx)
	return opCtxt.Variables["id"].(*string)
}