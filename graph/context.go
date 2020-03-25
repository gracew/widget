package graph

import "context"

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
const apiIDCtxKey = "apiID"

func apiID(ctx context.Context) *string {
	raw, _ := ctx.Value(apiIDCtxKey).(*string)
	return raw
}