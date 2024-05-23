package usecase

import (
	"back/graph/model"
	"back/pkg/adapter/repository"
	"context"
)

type venueUsecase struct {
	venueRepository repository.VenueRepository
}

type VenueUsecase interface {
	List(ctx context.Context) ([]*model.Venue, error)
	FindVenue(ctx context.Context, name string) ([]*model.Venue, error)
}

func NewVenueUsecase(ur repository.VenueRepository) VenueUsecase {
	return &venueUsecase{
		venueRepository: ur,
	}
}

func (gu venueUsecase) FindVenue(ctx context.Context, name string) ([]*model.Venue, error) {

	return gu.venueRepository.FindVenue(ctx, name)
	//tu.todoRepository.List(ctx)
}

func (gu venueUsecase) List(ctx context.Context) ([]*model.Venue, error) {

	return gu.venueRepository.List(ctx)
	//tu.todoRepository.List(ctx)
}
