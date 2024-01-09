package graph

import (
	"context"
	"log"

	"github.com/RSO-project-Prepih/AI-service/handlers"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver1 struct{}

func (r *queryResolver) FamousPlaces(ctx context.Context) ([]*handlers.Place, error) {
	places, err := handlers.FetchFamousPlaces()
	if err != nil {
		log.Fatal(err)
	}

	placesPtr := make([]*handlers.Place, len(places))
	for i := range places {
		placesPtr[i] = &places[i]
	}
	return placesPtr, nil
}
