package test

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/repository"
	"testing"
)

var logUsuarioId int64
var logEventoId int64 = 1 // ID de evento genérico para teste
var logEntity *entity.Log

func createLogUsuario(t *testing.T) int64 {
	tempUser := entity.NewUsuario("logtester", "1234", false)
	err := repository.CreateUsuario(tempUser)
	if err != nil {
		// Se já existe, tenta buscar um existente com mesmo nome e usar seu ID
		rows, errQuery := repository.DB.Query("SELECT usuarioID FROM usuario WHERE username = ?", "logtester")
		if errQuery != nil {
			t.Fatalf("erro ao buscar usuário de teste: %v", errQuery)
		}
		defer rows.Close()
		if rows.Next() {
			var existingID int64
			err = rows.Scan(&existingID)
			if err != nil {
				t.Fatalf("erro ao ler ID do usuário existente: %v", err)
			}
			return existingID
		}
		t.Fatalf("usuário de teste já existe mas não foi possível recuperar o ID")
	}
	return tempUser.Id
}

func TestCreateLog(t *testing.T) {
	logUsuarioId = createLogUsuario(t)
	log := entity.NewLog(logUsuarioId, logEventoId, "Usuário testou a criação de log")
	err := repository.CreateLog(log)
	if err != nil {
		t.Fatalf("erro ao criar log: %v", err)
	}
	if log.Id == 0 {
		t.Error("ID do log não foi atribuído")
	}
	logEntity = log
}

func TestGetLogByID(t *testing.T) {
	if logEntity == nil {
		t.Skip("log de teste não criado")
	}
	log, err := repository.GetLogByID(logEntity.Id)
	if err != nil {
		t.Fatalf("erro ao buscar log: %v", err)
	}
	if log == nil {
		t.Error("log não encontrado")
	} else if log.Message != logEntity.Message {
		t.Errorf("mensagem do log incorreta: esperada %q, obtida %q", logEntity.Message, log.Message)
	}
}

func TestGetLogsByUsuarioID(t *testing.T) {
	logUsuarioId = createLogUsuario(t)
	logs, err := repository.GetLogsByUsuarioID(logUsuarioId)
	if err != nil {
		t.Fatalf("erro ao buscar logs do usuário: %v", err)
	}
	if len(logs) == 0 {
		t.Error("nenhum log retornado para o usuário")
	}
}

func TestListAllLogs(t *testing.T) {
	logs, err := repository.ListAllLogs(10)
	if err != nil {
		t.Fatalf("erro ao listar logs: %v", err)
	}
	if len(logs) == 0 {
		t.Error("nenhum log retornado na listagem geral")
	}
}
