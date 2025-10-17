package repository

import (
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/utils/timeToDateString"
	"database/sql"
	"fmt"
	"time"
)

// CreateAviso insere um aviso simples.
func CreateAviso(a *entity.Aviso) error {
	if id, ok, err := getAvisoID(a.Tipo, a.ReferenciaID); err != nil {
		return fmt.Errorf("erro ao consultar aviso existente: %w", err)
	} else if ok {
		const up = `UPDATE aviso SET mensagem = ?, criadoEm = ?, ativo = TRUE WHERE avisoID = ?`
		_, err := DB.Exec(up, a.Mensagem, timeToDateString.TimeToDateString(a.CriadoEm), id)
		if err != nil {
			return fmt.Errorf("erro ao atualizar aviso existente: %w", err)
		}
		a.ID = id
		return nil
	}

	// 2) Existe inativo? (reativar)
	const findInactive = `SELECT avisoID FROM aviso WHERE tipo = ? AND referenciaID <=> ? AND ativo = FALSE LIMIT 1`
	var inactiveID int64
	err := DB.QueryRow(findInactive, a.Tipo, a.ReferenciaID).Scan(&inactiveID)
	if err == nil {
		const reactivate = `UPDATE aviso SET mensagem = ?, criadoEm = ?, ativo = TRUE WHERE avisoID = ?`
		_, err2 := DB.Exec(reactivate, a.Mensagem, timeToDateString.TimeToDateString(a.CriadoEm), inactiveID)
		if err2 != nil {
			return fmt.Errorf("erro ao reativar aviso: %w", err2)
		}
		a.ID = inactiveID
		return nil
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("erro ao consultar aviso inativo: %w", err)
	}

	// 3) Inserir novo
	const ins = `INSERT INTO aviso (tipo, mensagem, referenciaID, criadoEm, ativo) VALUES (?, ?, ?, ?, TRUE)`
	res, err := DB.Exec(ins, a.Tipo, a.Mensagem, a.ReferenciaID, timeToDateString.TimeToDateString(a.CriadoEm))
	if err != nil {
		return fmt.Errorf("erro ao criar aviso: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("erro ao obter id do aviso: %w", err)
	}
	a.ID = id
	return nil
}

func getAvisoID(tipo string, referenciaID *int64) (int64, bool, error) {
	const q = `SELECT avisoID FROM aviso WHERE tipo = ? AND referenciaID <=> ? AND ativo = TRUE LIMIT 1`
	var id int64
	err := DB.QueryRow(q, tipo, referenciaID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, err
	}
	return id, true, nil
}

// ListAvisos retorna os avisos ativos (mais recentes primeiro).
func ListAvisos() ([]entity.Aviso, error) {
	const q = `SELECT avisoID, tipo, mensagem, referenciaID, criadoEm, ativo
	           FROM aviso WHERE ativo = TRUE
	           ORDER BY criadoEm DESC, avisoID DESC`

	rows, err := DB.Query(q)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar avisos: %w", err)
	}
	defer rows.Close()

	var out []entity.Aviso
	for rows.Next() {
		var a entity.Aviso
		var criadoStr string
		var refID sql.NullInt64
		if err := rows.Scan(&a.ID, &a.Tipo, &a.Mensagem, &refID, &criadoStr, &a.Ativo); err != nil {
			return nil, err
		}
		if refID.Valid {
			a.ReferenciaID = &refID.Int64
		}
		// criadoEm: DATETIME (YYYY-MM-DD HH:MM:SS)
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", criadoStr, time.Local); err == nil {
			a.CriadoEm = t
		}
		out = append(out, a)
	}
	return out, nil
}

// DeleteAvisoByTypeAndRef remove aviso específico por tipo + referência.
func DeleteAvisoByTypeAndRef(tipo string, referenciaID int64) error {
	const q = `DELETE FROM aviso WHERE tipo = ? AND referenciaID = ?`
	_, err := DB.Exec(q, tipo, referenciaID)
	return err
}

// DeleteAvisosByType remove todos avisos de um tipo (usado para limpar pendentes já resolvidos).
func DeleteAvisosByType(tipo string) error {
	const q = `DELETE FROM aviso WHERE tipo = ?`
	_, err := DB.Exec(q, tipo)
	return err
}

// DeactivateAvisosVencidosAntes apenas desativa avisos muito antigos (failsafe).
func DeactivateAvisosVencidosAntes(de time.Time) error {
	const q = `UPDATE aviso SET ativo = FALSE WHERE ativo = TRUE AND criadoEm < ?`
	_, err := DB.Exec(q, timeToDateString.TimeToDateString(de))
	return err
}
