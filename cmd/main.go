package main

import (
	"Apis/internal/users"
	"Apis/pkg/bootstrap"
	"Apis/pkg/handler"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	_ = godotenv.Load()

	db, err := bootstrap.NewDB()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	logger := bootstrap.NewLogger()

	repo := users.NewRepo(db, logger)
	service := users.NewService(logger, repo)
	ctx := context.Background()

	h := handler.NewUserHTTPServer(users.MakeEndPoints(ctx, service))

	port := os.Getenv("PORT")

	address := fmt.Sprintf("127.0.0.1:%s", port)

	server := &http.Server{
		Addr:    address,
		Handler: accessControl(h),
	}

	log.Fatal(server.ListenAndServe())
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control, Origin, X-Requested-With, Content-Type, Accept, Authorization, Content-Length, Accept-Encoding, X-CSRF-Token, X-Api-Key, X-Auth-Token")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
