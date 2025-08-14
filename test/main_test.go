package test

import (
	"AutoGRH/pkg/Repository"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func resetDB() {
	Repository.DB.Exec("SET FOREIGN_KEY_CHECKS=0")
	Repository.DB.Exec("TRUNCATE TABLE descanso")
	Repository.DB.Exec("TRUNCATE TABLE pagamento")
	Repository.DB.Exec("TRUNCATE TABLE vale")
	Repository.DB.Exec("TRUNCATE TABLE salario")
	Repository.DB.Exec("TRUNCATE TABLE falta")
	Repository.DB.Exec("TRUNCATE TABLE documento")
	Repository.DB.Exec("TRUNCATE TABLE ferias")
	Repository.DB.Exec("TRUNCATE TABLE folha_pagamento")
	Repository.DB.Exec("TRUNCATE TABLE log")
	Repository.DB.Exec("TRUNCATE TABLE funcionario")
	Repository.DB.Exec("TRUNCATE TABLE usuario")
	Repository.DB.Exec("SET FOREIGN_KEY_CHECKS=1")
}

func TestMain(m *testing.M) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Variáveis devem estar no ambiente.")
	}

	Repository.ConnectDB()

	resetDB()

	SeedOneOfEach()

	code := m.Run()

	resetDB()

	os.Exit(code)
}
