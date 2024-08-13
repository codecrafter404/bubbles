package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/codecrafter404/bubble/graph"
	"github.com/codecrafter404/bubble/graph/model"
	"github.com/codecrafter404/bubble/utils"
	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	db_path := os.Getenv("DB_PATH")
	if db_path == "" {
		db_path = "bubbles.db"
	}

	connection, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", db_path))

	if err != nil {
		fmt.Printf("Failed to open db connection: %s\n", err.Error())
		return
	}

	defer connection.Close()

	err = utils.MigrateDb(connection)
	if err != nil {
		fmt.Println("Failed to setup:", err)
		return
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Db: connection, EventChannel: make(chan *model.UpdateEvent)}}))

	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
