package main

import (
	"log"
	"net/http"
	"os"

	db "github.com/AngelPwG/devprofile/internal/db"
	handler "github.com/AngelPwG/devprofile/internal/handler"
	"github.com/AngelPwG/devprofile/internal/router"
)

func main() {
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "devprofile.db"
	}
	db, err := db.NewDB(path)
	if err != nil {
		log.Fatal(err)
	}
	handler := handler.NewHandler(db)
	r := router.NewRouter(handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
