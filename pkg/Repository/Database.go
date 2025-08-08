package Repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	createDatabaseIfNotExists()
	connectWithDatabase()
	createTables()
	seedDefaultData()
}

func createDatabaseIfNotExists() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	noDBdsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)
	db, err := sql.Open("mysql", noDBdsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão inicial: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}

	mustExec(db, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))

	log.Println("Banco de dados verificado/criado.")
}

func connectWithDatabase() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, password, host, port, dbname)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar com banco %s: %v", dbname, err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Erro ao pingar banco %s: %v", dbname, err)
	}

	log.Printf("Conexão final com banco %s OK.", dbname)
}

func createTables() {
	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS usuario (
			usuarioID BIGINT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(15) UNIQUE,
			password VARCHAR(255),
			isAdmin BOOLEAN
		);`,

		`CREATE TABLE IF NOT EXISTS evento (
			eventoID BIGINT AUTO_INCREMENT PRIMARY KEY,
			tipo VARCHAR(20)
		);`,

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

		`CREATE TABLE IF NOT EXISTS tipo_pagamento (
			tipoID BIGINT AUTO_INCREMENT PRIMARY KEY,
			tipo VARCHAR(20) UNIQUE
		);`,

		`CREATE TABLE IF NOT EXISTS vale (
			valeID BIGINT AUTO_INCREMENT PRIMARY KEY,
			funcionarioID BIGINT,
			valor FLOAT,
			data DATE,
			aprovado BOOLEAN,
			pago BOOLEAN,
			FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
		);`,

		`CREATE TABLE IF NOT EXISTS folha_pagamento (
			folhaID BIGINT AUTO_INCREMENT PRIMARY KEY,
			data DATE
		);`,

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

		`CREATE TABLE IF NOT EXISTS salario (
			salarioID BIGINT AUTO_INCREMENT PRIMARY KEY,
			funcionarioID BIGINT,
			inicio DATE,
			fim DATE DEFAULT NULL,
			valor FLOAT,
			FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
		);`,

		`CREATE TABLE IF NOT EXISTS falta (
			faltaID BIGINT AUTO_INCREMENT PRIMARY KEY,
			funcionarioID BIGINT,
			quantidade INTEGER,
			data DATE,
			FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
		);`,

		`CREATE TABLE IF NOT EXISTS documento (
			documentoID BIGINT AUTO_INCREMENT PRIMARY KEY,
			funcionarioID BIGINT,
			documento BLOB,
			FOREIGN KEY (funcionarioID) REFERENCES funcionario(funcionarioID)
		);`,

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

	for _, query := range tableQueries {
		mustExec(DB, query)
	}

	log.Println("Todas as tabelas foram criadas/verificadas com sucesso.")
}

func seedDefaultData() {
	tipos := []string{"salario", "vale", "outros"}
	for _, tipo := range tipos {
		insert := `
			INSERT INTO tipo_pagamento (tipo)
			SELECT * FROM (SELECT ?) AS tmp
			WHERE NOT EXISTS (
				SELECT tipo FROM tipo_pagamento WHERE tipo = ?
			) LIMIT 1;
		`
		mustExec(DB, insert, tipo, tipo)
	}

	eventos := []string{"LOGIN", "LOGOUT", "CRIAR", "ATUALIZAR", "DELETAR", "APROVAR", "NEGAR"}
	for _, evento := range eventos {
		query := `
			INSERT INTO evento (tipo)
			SELECT ? FROM DUAL
			WHERE NOT EXISTS (
				SELECT 1 FROM evento WHERE tipo = ?
			);
		`
		mustExec(DB, query, evento, evento)
	}

	log.Println("Dados padrão foram inseridos/verificados com sucesso.")
}

func mustExec(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Fatalf("Erro ao executar SQL:\n%v\nQuery: %s", err, query)
	}
}
