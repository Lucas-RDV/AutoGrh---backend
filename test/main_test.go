package test

import (
	"AutoGRH/pkg/repository"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func resetDB() {
	repository.DB.Exec("SET FOREIGN_KEY_CHECKS=0")
	repository.DB.Exec("TRUNCATE TABLE descanso")
	repository.DB.Exec("TRUNCATE TABLE pagamento")
	repository.DB.Exec("TRUNCATE TABLE vale")
	repository.DB.Exec("TRUNCATE TABLE salario")
	repository.DB.Exec("TRUNCATE TABLE falta")
	repository.DB.Exec("TRUNCATE TABLE documento")
	repository.DB.Exec("TRUNCATE TABLE ferias")
	repository.DB.Exec("TRUNCATE TABLE folha_pagamento")
	repository.DB.Exec("TRUNCATE TABLE log")
	repository.DB.Exec("TRUNCATE TABLE funcionario")
	repository.DB.Exec("TRUNCATE TABLE usuario")
	repository.DB.Exec("SET FOREIGN_KEY_CHECKS=1")
}

func TestMain(m *testing.M) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Variáveis devem estar no ambiente.")
	}

	repository.ConnectDB()

	resetDB()

	SeedOneOfEach()

	code := m.Run()

	resetDB()

	os.Exit(code)
}
