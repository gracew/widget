// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package graph

import (
	"github.com/gracew/widget/graph/generated"
)

func (r *Resolver) ActionDefinition() generated.ActionDefinitionResolver {
	return &actionDefinitionResolver{r}
}
func (r *Resolver) CreateDefinition() generated.CreateDefinitionResolver {
	return &createDefinitionResolver{r}
}
func (r *Resolver) DeleteDefinition() generated.DeleteDefinitionResolver {
	return &deleteDefinitionResolver{r}
}
func (r *Resolver) ListDefinition() generated.ListDefinitionResolver {
	return &listDefinitionResolver{r}
}
func (r *Resolver) ReadDefinition() generated.ReadDefinitionResolver {
	return &readDefinitionResolver{r}
}

type actionDefinitionResolver struct{ *Resolver }
type createDefinitionResolver struct{ *Resolver }
type deleteDefinitionResolver struct{ *Resolver }
type listDefinitionResolver struct{ *Resolver }
type readDefinitionResolver struct{ *Resolver }
