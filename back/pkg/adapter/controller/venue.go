package controller

import (
	"back/graph/model"
	"back/pkg/usecase"
	"context"
)

type VenueController interface {
	List(ctx context.Context) ([]*model.Venue, error)
	FindVenue(ctx context.Context, name string) ([]*model.Venue, error)
}

type venueController struct {
	venueUsecase usecase.VenueUsecase
}

func NewVenueController(gu usecase.VenueUsecase) VenueController {
	return &venueController{
		venueUsecase: gu,
	}
}

func (gc venueController) FindVenue(ctx context.Context, name string) ([]*model.Venue, error) {

	return gc.venueUsecase.FindVenue(ctx, name)
}

func (gc venueController) List(ctx context.Context) ([]*model.Venue, error) {

	return gc.venueUsecase.List(ctx)
}
