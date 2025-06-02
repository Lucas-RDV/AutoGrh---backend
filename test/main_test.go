package test

import (
	"AutoGRH/pkg/repository"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	repository.ConnectDB()

	// Limpar tabelas antes dos testes
	repository.DB.Exec("DELETE FROM falta")
	repository.DB.Exec("DELETE FROM funcionario")
	repository.DB.Exec("DELETE FROM usuario")
	// Adicione mais conforme necessário

	code := m.Run()

	// Limpar novamente, se quiser
	repository.DB.Exec("DELETE FROM falta")
	repository.DB.Exec("DELETE FROM funcionario")
	repository.DB.Exec("DELETE FROM usuario")

	os.Exit(code)
}
