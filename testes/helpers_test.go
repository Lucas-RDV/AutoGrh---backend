package testes

import (
	"AutoGRH/pkg/repository"
	"fmt"
)

var tables = []string{
	// apague filhos antes dos pais, se houver FK
	"documento",
	"descanso",
	"ferias",
	"pagamento",
	"folha_pagamento",
	"vale",
	"salario_real",
	"salario",
	"funcionario",
	"pessoa",
	"log",
	"usuario",
	"evento",
	"aviso",
}

func truncateAll() error {
	if repository.DB == nil {
		return fmt.Errorf("DB nil")
	}

	stmts := []string{
		"SET FOREIGN_KEY_CHECKS = 0",

		// Filhas primeiro:
		"TRUNCATE TABLE descanso",
		"TRUNCATE TABLE falta",
		"TRUNCATE TABLE documento",
		"TRUNCATE TABLE salario_real",
		"TRUNCATE TABLE salario",
		"TRUNCATE TABLE pagamento",
		"TRUNCATE TABLE folha_pagamento",
		"TRUNCATE TABLE vale", // se existir

		// Depois as pais:
		"TRUNCATE TABLE ferias",
		"TRUNCATE TABLE funcionario",
		"TRUNCATE TABLE pessoa",
		"TRUNCATE TABLE usuario",
		"TRUNCATE TABLE log", // seu log/auditoria, se usado

		"SET FOREIGN_KEY_CHECKS = 1",
	}

	for _, s := range stmts {
		if _, err := repository.DB.Exec(s); err != nil {
			return fmt.Errorf("truncateAll step %q: %w", s, err)
		}
	}
	return nil
}
