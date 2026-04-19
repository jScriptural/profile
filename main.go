package main

import (
	"log"
	"net/http"
	"os"
	"profile/internal/api"
	"profile/internal/service"
	mw "profile/middleware"
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
		Handler: mw.CORS(mux,"*"),
	}

	log.Printf("Server is starting at %v: db->%v\n",server.Addr,DATABASE_URL);
	err := server.ListenAndServe()
	log.Fatal(err)
}
