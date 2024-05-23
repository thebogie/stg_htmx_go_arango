package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.46

import (
	"back/graph/model"
	"context"
)

// FindVenue is the resolver for the FindVenue field.
func (r *queryResolver) FindVenue(ctx context.Context, name string) ([]*model.Venue, error) {
	return r.Venue.FindVenue(ctx, name)
}
