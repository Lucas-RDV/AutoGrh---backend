package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"AutoGRH/pkg/Bootstrap"
	"AutoGRH/pkg/HTTP/router"
)

func main() {
	_ = godotenv.Load()

	app := Bootstrap.Load()
	if err := Bootstrap.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	auth := Bootstrap.BuildAuth(app)
	mux := router.New(auth)

	http.ListenAndServe(":8080", mux)
}
