package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/codecrafter404/bubble/config"
	"github.com/codecrafter404/bubble/ent"
	"github.com/codecrafter404/bubble/ent/migrate"
	"github.com/codecrafter404/bubble/gql"
	"github.com/codecrafter404/bubble/server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	config := config.Config{
		DbPath: "file:" + "bubbles.db" + "?_foreign_keys=on",
		OrderConfig: config.OrderConfig{
			MaxNotificationQueue: 1000,
			MaximalIdentifiers:   100,
		},
		ServerConfig: config.ServerConfig{
			ServerPort:        8080,
			PlaygroundEnabled: true,
			CorsConfig: config.CorsConfig{
				AccessControlAllowCredentials: true,
				AccessControlAllowHeaders:     []string{"*"},
				AccessControlAllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
				AccessControlAllowOrigin:      []string{"*"},
			},
		},
	}
	log.Info().Int("port", config.ServerConfig.ServerPort).Msg("Configuration loaded")
	client, err := ent.Open("sqlite3", config.DbPath)
	if err != nil {
		log.Fatal().Str("path", config.DbPath).Err(err).Msg("Failed to open database")
	}
	defer client.Close()

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	mux := http.NewServeMux()

	if config.ServerConfig.PlaygroundEnabled {
		mux.Handle("/playground",
			playground.Handler("bubbles", "/query"),
		)
	}

	srv := handler.NewDefaultServer(gql.NewSchema(client, &config))
	srv.Use(entgql.Transactioner{TxOpener: client})

	mux.Handle("/query", srv)

	log.Info().Int("port", config.ServerConfig.ServerPort).Msg("Listening for incoming connections")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServerConfig.ServerPort), server.CorsMiddleware(&config, mux)); err != nil {
		log.Fatal().Err(err).Msg("The server did terminate.")
	}
}
