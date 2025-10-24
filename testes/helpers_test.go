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
	for _, t := range tables {
		if _, err := repository.DB.Exec(fmt.Sprintf("DELETE FROM %s", t)); err != nil {
			return fmt.Errorf("truncate %s: %w", t, err)
		}
		_, _ = repository.DB.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", t))
	}
	return nil
}
