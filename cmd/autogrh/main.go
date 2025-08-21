package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"AutoGRH/pkg/Bootstrap"
	"AutoGRH/pkg/HTTP/router"
	middleware "AutoGRH/pkg/controller/middleware"
)

func main() {
	_ = godotenv.Load()

	app := Bootstrap.Load()
	if err := Bootstrap.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	auth := Bootstrap.BuildAuth(app)
	routes := router.New(auth)

	cors := middleware.NewCORS(middleware.CORSConfig{

		AllowedOrigins:   []string{"*"}, // aceita todos para teste. mudar depois
		AllowedMethods:   nil,           // default
		AllowedHeaders:   nil,           // default
		AllowCredentials: false,
	})

	http.ListenAndServe(":8080", cors(routes))
}
