package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/codecrafter404/bubble/graph"
	"github.com/codecrafter404/bubble/graph/model"
	"github.com/codecrafter404/bubble/utils"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
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
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Db: connection, EventChannel: []chan *model.UpdateEvent{}}}))

	srv.AddTransport(transport.SSE{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	// router.Handle("/config", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	app_fs := http.FileServer(http.Dir("./app/"))
	router.Handle("/*", app_fs)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
