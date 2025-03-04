package main

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"helix.io/helix"
	"helix.io/helix/ent"
	"helix.io/helix/ent/migrate"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger, _ := zap.NewProduction()
	router := chi.NewRouter()
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)
	// Create ent.Client and run the schema migration.
	client, err := ent.Open("mysql", "user:password@tcp(127.0.0.1:3306)/mydb?parseTime=True")
	if err != nil {
		logger.Fatal("open ent.Client failed", zap.Error(err))
	}
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		logger.Fatal("ent.Client create schema failed", zap.Error(err))
	}

	// Configure the server and start listening on :8081.
	srv := handler.NewDefaultServer(helix.NewSchema(client))
	router.Handle("/", playground.Handler("Catalyst", "/graphql"))
	router.Handle("/graphql", srv)
	logger.Info("the api service is severing on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		logger.Fatal("the api service is terminated", zap.Error(err))
	}
}
