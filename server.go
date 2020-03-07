package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/gracew/widget/graph"
	"github.com/gracew/widget/graph/generated"
	"github.com/gracew/widget/graph/model"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	db := pg.Connect(&pg.Options{User: "postgres"})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	// TODO(gracew): remove cors later
	http.Handle("/query", cors.Default().Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		(*model.API)(nil),
		(*model.Deploy)(nil),
		(*model.Auth)(nil),
		(*model.TestToken)(nil),
		(*model.CustomLogic)(nil),
	} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type createdBy struct {
	CreatedBy string `json:"createdBy"`
}

type unauthorized struct {
	Message string `json:"message"`
}