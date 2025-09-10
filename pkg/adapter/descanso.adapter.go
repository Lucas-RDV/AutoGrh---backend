package Adapter

import (
	"AutoGRH/pkg/entity"
)

// DescansoRepositoryAdapter conecta o pacote repository com a interface usada no service
type DescansoRepositoryAdapter struct {
	create             func(d *entity.Descanso) error
	getByID            func(id int64) (*entity.Descanso, error)
	update             func(d *entity.Descanso) error
	delete             func(id int64) error
	getByFeriasID      func(feriasID int64) ([]*entity.Descanso, error)
	getByFuncionarioID func(funcionarioID int64) ([]*entity.Descanso, error)
	getAprovados       func() ([]*entity.Descanso, error)
	getPendentes       func() ([]*entity.Descanso, error)
}

// Construtor do adapter
func NewDescansoRepositoryAdapter(
	create func(d *entity.Descanso) error,
	getByID func(id int64) (*entity.Descanso, error),
	update func(d *entity.Descanso) error,
	delete func(id int64) error,
	getByFeriasID func(feriasID int64) ([]*entity.Descanso, error),
	getByFuncionarioID func(funcionarioID int64) ([]*entity.Descanso, error),
	getAprovados func() ([]*entity.Descanso, error),
	getPendentes func() ([]*entity.Descanso, error),
) *DescansoRepositoryAdapter {
	return &DescansoRepositoryAdapter{
		create:             create,
		getByID:            getByID,
		update:             update,
		delete:             delete,
		getByFeriasID:      getByFeriasID,
		getByFuncionarioID: getByFuncionarioID,
		getAprovados:       getAprovados,
		getPendentes:       getPendentes,
	}
}

// Implementação da interface service.DescansoRepository
func (a *DescansoRepositoryAdapter) Create(d *entity.Descanso) error {
	return a.create(d)
}
func (a *DescansoRepositoryAdapter) GetDescansoByID(id int64) (*entity.Descanso, error) {
	return a.getByID(id)
}
func (a *DescansoRepositoryAdapter) Update(d *entity.Descanso) error {
	return a.update(d)
}
func (a *DescansoRepositoryAdapter) Delete(id int64) error {
	return a.delete(id)
}
func (a *DescansoRepositoryAdapter) GetDescansosByFeriasID(feriasID int64) ([]*entity.Descanso, error) {
	return a.getByFeriasID(feriasID)
}
func (a *DescansoRepositoryAdapter) GetDescansosByFuncionarioID(funcionarioID int64) ([]*entity.Descanso, error) {
	return a.getByFuncionarioID(funcionarioID)
}
func (a *DescansoRepositoryAdapter) GetDescansosAprovados() ([]*entity.Descanso, error) {
	return a.getAprovados()
}
func (a *DescansoRepositoryAdapter) GetDescansosPendentes() ([]*entity.Descanso, error) {
	return a.getPendentes()
}
