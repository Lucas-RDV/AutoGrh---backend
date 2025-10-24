package testes

import (
	"log"
	"os"
	"testing"

	"AutoGRH/pkg/repository"
)

const testDBName = "autogrh_test"

func TestMain(m *testing.M) {
	// Blindagem para não rodar em DB de dev/prod por engano
	if os.Getenv("DB_NAME") == "autogrh" {
		log.Fatal("Não rode testes com DB_NAME=autogrh. Use um schema de testes, ex.: autogrh_test.")
	}

	// Defaults de teste (você pode sobrescrever via env antes do go test)
	mustSetEnvDefault("DB_USER", "root")
	mustSetEnvDefault("DB_PASSWORD", "C27<jP^@3Adn")
	mustSetEnvDefault("DB_HOST", "127.0.0.1")
	mustSetEnvDefault("DB_PORT", "3306")
	_ = os.Setenv("DB_NAME", testDBName)

	// Sobe DB e cria tabelas
	repository.ConnectDB() // cria schema e tabelas (reaproveita seu createTables)
	if repository.DB == nil {
		log.Fatal("DB não inicializado para testes")
	}

	// Limpa tabelas antes da suíte
	if err := truncateAll(); err != nil {
		log.Fatalf("erro ao limpar tabelas antes: %v", err)
	}

	code := m.Run()

	// Limpa tabelas ao final (opcional)
	if err := truncateAll(); err != nil {
		log.Printf("erro ao limpar tabelas depois: %v", err)
	}

	os.Exit(code)
}
