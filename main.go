package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/RSO-project-Prepih/AI-service/graph"
	"github.com/RSO-project-Prepih/AI-service/handlers"
)

const defaultPort = "8080"

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	popularPlaces, err := handlers.FetchFamousPlaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, place := range popularPlaces {
		fmt.Printf("Famous Place: %s, Rating: %d, Kinds: %s, Type: %s, Coordinates: [%f, %f]\n",
			place.Properties.Name,
			place.Properties.Rate,
			place.Properties.Kinds,
			place.Geometry.Type,
			place.Geometry.Coordinates[0], // Assuming the first element is latitude
			place.Geometry.Coordinates[1], // Assuming the second element is longitude
		)
	}

	log.Println(popularPlaces)

	userID := "d8082c4f-85ef-41f7-943b-bc041d08fa2f"
	imageID := "600023cc-ad35-4b2c-a112-77503e6d1f05"
	handlers.PostColorEnhancementPhoto(userID, imageID)

}
