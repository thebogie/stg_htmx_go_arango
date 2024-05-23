package repository

import (
	"back/graph/model"
	"context"
	"github.com/arangodb/go-driver"
	"log"
)

type VenueRepository interface {
	List(ctx context.Context) ([]*model.Venue, error)
	FindVenue(ctx context.Context, name string) ([]*model.Venue, error)
}

type venuerepository struct {
	db driver.Database
}

func NewVenueRepository(db driver.Database) VenueRepository {
	return &venuerepository{
		db: db,
	}
}

func (gr *venuerepository) FindVenue(ctx context.Context, name string) ([]*model.Venue, error) {
	query := "FOR venue IN venue FILTER LOWER(venue.address) LIKE CONCAT('%', LOWER('" + name + "'), '%') RETURN venue"
	//TODO: pass context from top level?

	cursor, err := gr.db.Query(ctx, query, nil)
	if err != nil {
		log.Fatal("Error login Query to db")
	}
	defer cursor.Close()

	//_, err = cursor.ReadDocument(ctx, &retuser)
	var results []*model.Venue // Replace with your struct type

	// Iterate over the cursor and append results to the array
	for {
		var doc model.Venue // Replace with your struct type
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break // No more documents in the cursor
		} else if err != nil {
			log.Fatal(err)
		}
		results = append(results, &doc)
	}

	return results, nil
}

func (gr *venuerepository) List(ctx context.Context) ([]*model.Venue, error) {
	query := "FOR doc IN venue RETURN doc"
	//TODO: pass context from top level?

	cursor, err := gr.db.Query(ctx, query, nil)
	if err != nil {
		log.Fatal("Error login Query to db")
	}
	defer cursor.Close()

	//_, err = cursor.ReadDocument(ctx, &retuser)
	var results []*model.Venue // Replace with your struct type

	// Iterate over the cursor and append results to the array
	for {
		var doc model.Venue // Replace with your struct type
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break // No more documents in the cursor
		} else if err != nil {
			log.Fatal(err)
		}
		results = append(results, &doc)
	}

	return results, nil
}
