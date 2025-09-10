package Adapter

import (
	"AutoGRH/pkg/entity"
)

type SalarioRealRepositoryAdapter struct {
	create     func(s *entity.SalarioReal) error
	listByFunc func(funcionarioID int64) ([]*entity.SalarioReal, error)
	getAtual   func(funcionarioID int64) (*entity.SalarioReal, error)
	update     func(s *entity.SalarioReal) error
	delete     func(id int64) error
}

func NewSalarioRealRepositoryAdapter(
	create func(s *entity.SalarioReal) error,
	listByFunc func(funcionarioID int64) ([]*entity.SalarioReal, error),
	getAtual func(funcionarioID int64) (*entity.SalarioReal, error),
	update func(s *entity.SalarioReal) error,
	delete func(id int64) error,
) *SalarioRealRepositoryAdapter {
	return &SalarioRealRepositoryAdapter{
		create:     create,
		listByFunc: listByFunc,
		getAtual:   getAtual,
		update:     update,
		delete:     delete,
	}
}

// ---- Implementações exigidas pela interface SalarioRealRepository ----

func (a *SalarioRealRepositoryAdapter) Create(s *entity.SalarioReal) error {
	return a.create(s)
}

// nome EXATO que a interface pede
func (a *SalarioRealRepositoryAdapter) GetByFuncionarioID(funcionarioID int64) ([]*entity.SalarioReal, error) {
	return a.listByFunc(funcionarioID)
}

// nome EXATO que a interface pede
func (a *SalarioRealRepositoryAdapter) GetAtual(funcionarioID int64) (*entity.SalarioReal, error) {
	return a.getAtual(funcionarioID)
}

func (a *SalarioRealRepositoryAdapter) Update(s *entity.SalarioReal) error {
	return a.update(s)
}

func (a *SalarioRealRepositoryAdapter) Delete(id int64) error {
	return a.delete(id)
}
