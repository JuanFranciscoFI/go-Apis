package main

import (
	"Apis/internal/users"
	"Apis/pkg/bootstrap"
	"Apis/pkg/handler"
	"context"
	"log"
	"net/http"
)

func main() {
	server := http.NewServeMux() // Create a new server

	db := bootstrap.NewDB()
	logger := bootstrap.NewLogger()

	repo := users.NewRepo(db, logger)
	service := users.NewService(logger, repo)
	ctx := context.Background()

	handler.NewUserHTTPServer(ctx, server, users.MakeEndPoints(ctx, service))

	log.Fatal(http.ListenAndServe(":8080", server)) // Start the server on port 8080
}
