package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/StevenACoffman/pg-gql-todo/assets"
	"github.com/StevenACoffman/pg-gql-todo/generated/gql"
	"github.com/StevenACoffman/pg-gql-todo/resolvers"
	"github.com/StevenACoffman/pg-gql-todo/sqldb"
)

const defaultPort = "3000"

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbInfo := sqldb.NewDBInfo("postgres", "", "127.0.0.1", "stevetest", "public")
	fmt.Println("Running migration")
	pool, err := sqldb.NewDBPool(ctx, dbInfo, true)
	if err != nil {
		log.Fatal(err)
	}
	srv := handler.NewDefaultServer(
		gql.NewExecutableSchema(gql.Config{Resolvers: &resolvers.Resolver{DBPool: pool}}),
	)
	realFrontend, err := fs.Sub(assets.EmbeddedFiles, "static")
	if err != nil {
		panic("Error getting frontend/build from embedded FS: " + err.Error())
	}
	fileHandler := http.FileServer(http.FS(realFrontend))
	http.Handle("/", fileHandler)
	http.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
