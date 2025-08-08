package main

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Variáveis devem estar no ambiente.")
	}
}

func main() {
	
}
