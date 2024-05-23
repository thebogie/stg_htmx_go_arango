package registry

import (
	"back/pkg/adapter/controller"
	"back/pkg/adapter/repository"
	"back/pkg/usecase"
)

func (r *registry) NewVenueController() controller.VenueController {
	gr := repository.NewVenueRepository(r.db)
	gc := usecase.NewVenueUsecase(gr)
	return controller.NewVenueController(gc)
}
