package test

import (
	"AutoGRH/pkg/Entity"
	"AutoGRH/pkg/repository"
	"fmt"
	"log"
	"time"
)

func SeedOneOfEach() {
	// IDs dependentes dos seeds padrão criados pelo ConnectDB()
	var tipoSalarioID int64
	if err := repository.DB.QueryRow("SELECT tipoID FROM tipo_pagamento WHERE tipo = 'salario'").Scan(&tipoSalarioID); err != nil {
		log.Fatal(fmt.Errorf("seed tipo_pagamento 'salario': %w", err))
	}
	var eventoCriarID int64
	if err := repository.DB.QueryRow("SELECT eventoID FROM evento WHERE tipo = 'CRIAR'").Scan(&eventoCriarID); err != nil {
		log.Fatal(fmt.Errorf("seed evento 'CRIAR': %w", err))
	}

	// 1) Usuario
	u := Entity.NewUsuario("seed_user", "seed_pass", true)
	if err := repository.CreateUsuario(u); err != nil {
		log.Fatal(fmt.Errorf("seed usuario: %w", err))
	}

	// 2) Funcionario
	f := Entity.NewFuncionario(
		"Funcionario Seed", "RG000", "CPF000", "PIS000", "CTPF000",
		"Rua Seed, 123", "(11)1111-1111", "(11)2222-2222", "Analista",
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		3500.00,
	)
	if err := repository.CreateFuncionario(f); err != nil {
		log.Fatal(fmt.Errorf("seed funcionario: %w", err))
	}

	// 3) Salario (vigente)
	s := Entity.NewSalario(f.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 4000.00)
	if err := repository.CreateSalario(s); err != nil {
		log.Fatal(fmt.Errorf("seed salario: %w", err))
	}

	// 4) Documento
	doc := Entity.NewDocumento([]byte("documento de teste"), f.ID)
	if err := repository.CreateDocumento(doc); err != nil {
		log.Fatal(fmt.Errorf("seed documento: %w", err))
	}

	// 5) Férias
	fer := &Entity.Ferias{
		FuncionarioID: f.ID,
		Dias:          30,
		Inicio:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Vencimento:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Vencido:       false,
		Valor:         0,
	}
	if err := repository.CreateFerias(fer); err != nil {
		log.Fatal(fmt.Errorf("seed ferias: %w", err))
	}

	// 6) Descanso (parcial das férias)
	des := Entity.NewDescanso(
		time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
		fer.ID,
	)
	if err := repository.CreateDescanso(des); err != nil {
		log.Fatal(fmt.Errorf("seed descanso: %w", err))
	}

	// 7) Falta (1 no mês)
	fal := Entity.NewFalta(1, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), f.ID)
	if err := repository.CreateFalta(fal); err != nil {
		log.Fatal(fmt.Errorf("seed falta: %w", err))
	}

	// 8) Folha de Pagamentos
	folha := Entity.NewFolhaPagamentos(time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC))
	if err := repository.CreateFolha(folha); err != nil {
		log.Fatal(fmt.Errorf("seed folha_pagamento: %w", err))
	}

	// 9) Pagamento (tipo salário)
	pay := Entity.NewPagamento(tipoSalarioID, time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC), 1000.00)
	pay.FuncionarioID = f.ID
	pay.FolhaID = folha.ID
	if err := repository.CreatePagamento(pay); err != nil {
		log.Fatal(fmt.Errorf("seed pagamento: %w", err))
	}

	// 10) Vale
	vale := Entity.NewVale(f.ID, 500.00, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))
	if err := repository.CreateVale(vale); err != nil {
		log.Fatal(fmt.Errorf("seed vale: %w", err))
	}

	// 11) Log (evento CRIAR)
	lg := Entity.NewLog(u.ID, eventoCriarID, "seed inicial")
	if err := repository.CreateLog(lg); err != nil {
		log.Fatal(fmt.Errorf("seed log: %w", err))
	}
}
