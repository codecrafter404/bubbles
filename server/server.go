package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/codecrafter404/bubble/graph"
	"github.com/codecrafter404/bubble/resolvers"

	// "github.com/codecrafter404/bubble/utils"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

	db, err := gorm.Open(sqlite.Open(db_path))

	if err != nil {
		log.Fatalf("Failed to open database: %+v\n", err)
	}

	// migrate db

	db.AutoMigrate(&graph.Order{})
	db.AutoMigrate(&graph.OrderCustomItem{})
	db.AutoMigrate(&graph.OrderItem{})

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &resolvers.Resolver{Db: db, EventChannel: []chan *graph.UpdateEvent{}}}))

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

	router.Handle("/config", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	app_fs := http.FileServer(http.Dir("./app/"))
	router.Handle("/*", app_fs)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
