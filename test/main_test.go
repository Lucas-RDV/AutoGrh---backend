package test

import (
	"AutoGRH/pkg/Repository"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado. Variáveis devem estar no ambiente.")
	}

	Repository.ConnectDB()

	// Limpar tabelas antes dos testes
	Repository.DB.Exec("DELETE FROM falta")
	Repository.DB.Exec("DELETE FROM funcionario")
	Repository.DB.Exec("DELETE FROM usuario")
	// Adicione mais conforme necessário

	code := m.Run()

	// Limpar novamente, se quiser
	Repository.DB.Exec("DELETE FROM falta")
	Repository.DB.Exec("DELETE FROM funcionario")
	Repository.DB.Exec("DELETE FROM usuario")

	os.Exit(code)
}
