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
	"github.com/codecrafter404/bubble/utils"
	"github.com/google/uuid"
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

	config, err := config.LoadConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	log.Logger = log.Logger.With().Logger().Level(config.LogLevel)

	client, err := ent.Open("sqlite3", config.DbPath)
	if err != nil {
		log.Fatal().Str("path", config.DbPath).Err(err).Msg("Failed to open database")
	}
	defer client.Close()

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		log.Fatal().Err(err).Str("connection_string", config.DbPath).Msg("Failed to run database migrations")
	}

	mux := http.NewServeMux()

	if config.ServerConfig.PlaygroundEnabled {
		mux.Handle("/playground",
			playground.Handler("bubbles", "/query"),
		)
	}

	disconnect := make(chan uuid.UUID, 1024)
	subscribers := make([]*utils.Subscriber, 0)
	notifier := make(chan gql.NotificationType, config.OrderConfig.MaxNotificationQueue)

	ctx := context.Background()

	go server.OrderBalancer(ctx, client, notifier, &subscribers, disconnect, &config)

	srv := handler.NewDefaultServer(gql.NewSchema(client, &config, &subscribers, disconnect, notifier))
	srv.Use(entgql.Transactioner{TxOpener: client})

	mux.Handle("/query", srv)

	log.Info().Int("port", config.ServerConfig.ServerPort).Msg("Listening for incoming connections")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServerConfig.ServerPort), server.CorsMiddleware(&config, mux)); err != nil {
		log.Fatal().Err(err).Msg("The server did terminate.")
	}
}
