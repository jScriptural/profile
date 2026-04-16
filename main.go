package main

import (
	"log"
	"net/http"
	"os"
	"profile/internal/api"
	"profile/internal/service"
	"profile/middleware"
	"profile/store"
)


func main() {
	PORT := os.Getenv("PORT");
	DATABASE_URL := os.Getenv("DATABASE_URL")
	if PORT == "" {
		PORT = "8080"
	}

	if DATABASE_URL == "" {
		DATABASE_URL = "profile.db"
	}

	dbHandle := store.NewDBHandle(DATABASE_URL)
	service := service.NewService(dbHandle)
	handler := api.NewHandler(service)

	mux := handler.Routes()
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: middleware.EnableCORS(mux),
	}

	log.Printf("Server is starting at %v\n",server.Addr);
	err := server.ListenAndServe()
	log.Fatal(err)
}
