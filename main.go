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
	if PORT == "" {
		PORT = "8080"
	}

	dbHandle := store.NewDBHandle("profile.db")
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
