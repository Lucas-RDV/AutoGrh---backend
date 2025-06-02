package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func ConnectDB() {

	noDBdsn := "root:C27<jP^@3Adn@tcp(127.0.0.1:3306)/"
	dsn := "root:C27<jP^@3Adn@tcp(127.0.0.1:3306)/autogrh?charset=utf8"

	db, err := sql.Open("mysql", noDBdsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão inicial: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS autogrh")
	if err != nil {
		log.Fatalf("Erro ao criar banco de dados: %v", err)
	}
	log.Println("Banco de dados verificado/criado.")

	err = db.Close()
	if err != nil {
		log.Fatalf("Erro ao fechar banco de dados: %v", err)
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar com banco autogrh: %v", err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatalf("Erro ao pingar banco autogrh: %v", err)
	}

	log.Println("Conexão final com banco autogrh OK.")

	err = createTables()
	if err != nil {
		log.Fatalf("Erro ao criar tablas: %v", err)
	}
}

func createTables() error {
	queries := []string{
		// Tabela Usuario
		`CREATE TABLE IF NOT EXISTS usuario (
            usuarioID BIGINT AUTO_INCREMENT PRIMARY KEY,
            username VARCHAR(15) UNIQUE,
            password VARCHAR(255),
    		isAdmin BOOLEAN
        );`,

		// Tabela Evento
		`CREATE TABLE IF NOT EXISTS evento (
            eventoID BIGINT AUTO_INCREMENT PRIMARY KEY,
            tipo VARCHAR(20)
        );`,

		// Tabela Funcionario
		`CREATE TABLE IF NOT EXISTS funcionario (
            funcionarioID BIGINT AUTO_INCREMENT PRIMARY KEY,
            nome VARCHAR(100),
            rg VARCHAR(20),
            cpf VARCHAR(20),
            pis VARCHAR(20),
            ctpf VARCHAR(20),
            endereco TEXT,
            contato VARCHAR(100),
            contatoEmergencia VARCHAR(100),
            nascimento DATE,
            admissao DATE,
            demissao DATE NULL,
            cargo VARCHAR(50),
            salarioInicial FLOAT,
            feriasDisponiveis INT
        );`,

		// Tabela Tipo_Pagamento
		`CREATE TABLE IF NOT EXISTS tipo_pagamento (
            tipoID BIGINT AUTO_INCREMENT PRIMARY KEY,
            tipo VARCHAR(20) UNIQUE
        );`,

		//Tabela Vale
		`CREATE TABLE IF NOT EXISTS vale (
    valeID BIGINT AUTO_INCREMENT PRIMARY KEY,
    funcionarioID BIGINT,
    valor FLOAT,
    data DATE,
    aprovado BOOLEAN,
    pago BOOLEAN,
    FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
);`,
		// Tabela Folha_Pagamento
		`CREATE TABLE IF NOT EXISTS folha_pagamento (
            folhaID BIGINT AUTO_INCREMENT PRIMARY KEY,
            data DATE
        );`,

		// Tabela Pagamento
		`CREATE TABLE IF NOT EXISTS pagamento (
            pagamentoID BIGINT AUTO_INCREMENT PRIMARY KEY,
            funcionarioID BIGINT,
            folhaID BIGINT,
            tipoID BIGINT,
            valor FLOAT,
            data DATE,
            FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID),
            FOREIGN KEY (folhaID) REFERENCES folha_pagamento(folhaID),
            FOREIGN KEY (tipoID) REFERENCES tipo_pagamento(tipoID)
        );`,

		// Tabela Salario
		`CREATE TABLE IF NOT EXISTS salario (
            salarioID BIGINT AUTO_INCREMENT PRIMARY KEY,
            funcionarioID BIGINT,
            inicio DATE,
            fim DATE DEFAULT NULL,
            valor FLOAT,
            FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
        );`,

		// Tabela Falta
		`CREATE TABLE IF NOT EXISTS falta (
            faltaID BIGINT AUTO_INCREMENT PRIMARY KEY,
            funcionarioID BIGINT,
            quantidade INTEGER,
            data DATE,
            FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
        );`,

		// Tabela Documento
		`CREATE TABLE IF NOT EXISTS documento (
            documentoID BIGINT AUTO_INCREMENT PRIMARY KEY,
            funcionarioID BIGINT,
            documento BLOB,
            FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
        );`,

		// Tabela Ferias
		`CREATE TABLE IF NOT EXISTS ferias (
    		feriasID BIGINT AUTO_INCREMENT PRIMARY KEY,
    		funcionarioID BIGINT,
    		dias INT,
    		inicio DATE,
   	 		vencimento DATE,
   	 		vencido BOOLEAN,
   	 		valor FLOAT,
    		FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
		);`,

		// Tabela Descanso
		`CREATE TABLE IF NOT EXISTS descanso (
            descansoID BIGINT AUTO_INCREMENT PRIMARY KEY,
            feriasID BIGINT,
            inicio DATE,
            fim DATE,
            valor FLOAT,
            pago BOOLEAN,
            aprovado BOOLEAN,
            FOREIGN KEY (feriasID) REFERENCES ferias(feriasID)
        );`,

		// Tabela Log
		`CREATE TABLE IF NOT EXISTS log (
    		logID BIGINT AUTO_INCREMENT PRIMARY KEY,
			usuarioID BIGINT,
   			eventoID BIGINT,
    		data TIMESTAMP,
    		action VARCHAR(200),
    		FOREIGN KEY (usuarioID) REFERENCES usuario(usuarioID),
    		FOREIGN KEY (eventoID) REFERENCES evento(eventoID)
		);`,
	}

	var err error
	for _, query := range queries {
		_, err = DB.Exec(query)
		if err != nil {
			log.Fatalf("Erro ao criar tabela: %v\nQuery: %s", err, query)
		}
	}

	// Inserir tipos de pagamento padrão, se não existirem
	tipoPagamentos := []string{"salario", "vale", "outros"}
	for _, tipo := range tipoPagamentos {
		insert := `
        INSERT INTO tipo_pagamento (tipo)
        SELECT * FROM (SELECT ?) AS tmp
        WHERE NOT EXISTS (
            SELECT tipo FROM tipo_pagamento WHERE tipo = ?
        ) LIMIT 1;
        `
		_, err = DB.Exec(insert, tipo, tipo)
		if err != nil {
			log.Fatalf("Erro ao inserir tipo_pagamento padrão (%s): %v", tipo, err)
		}
	}
	// inserir tipos de eventos de log, se nao existirem
	eventos := []string{"LOGIN", "LOGOUT", "CRIAR", "ATUALIZAR", "DELETAR", "APROVAR", "NEGAR"}

	for _, evento := range eventos {
		query := `
		INSERT INTO evento (tipo)
		SELECT ? FROM DUAL
		WHERE NOT EXISTS (
			SELECT 1 FROM evento WHERE tipo = ?
		);
	`
		_, err := DB.Exec(query, evento, evento)
		if err != nil {
			log.Fatalf("Erro ao inserir tipo de evento padrão (%s): %v", evento, err)
		}
	}

	log.Println("Todas as tabelas foram criadas/verificadas com sucesso.")
	return err
}
