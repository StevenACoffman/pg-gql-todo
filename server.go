package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/StevenACoffman/gqlgen-todos/generated/gql"
	"github.com/StevenACoffman/gqlgen-todos/resolvers"
	"github.com/StevenACoffman/gqlgen-todos/sqldb"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbInfo := sqldb.NewDBInfo("postgres", "", "127.0.0.1", "stevetest", "public")
	fmt.Println("Running migration")
	pool, err := sqldb.NewDBPool(dbInfo, true)
	if err != nil {
		log.Fatal(err)
	}
	srv := handler.NewDefaultServer(
		gql.NewExecutableSchema(gql.Config{Resolvers: &resolvers.Resolver{DBPool: pool}}),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
