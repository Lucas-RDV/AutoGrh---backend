package Adapter

import (
	"AutoGRH/pkg/entity"
)

type ValeRepositoryAdapter struct {
	create                func(v *entity.Vale) error
	getByID               func(id int64) (*entity.Vale, error)
	getByFuncID           func(funcionarioID int64) ([]entity.Vale, error)
	update                func(v *entity.Vale) error
	softDelete            func(id int64) error
	delete                func(id int64) error
	listPendentes         func() ([]entity.Vale, error)
	listAprovadosNaoPagos func() ([]entity.Vale, error)
}

// Construtor
func NewValeRepositoryAdapter(
	create func(v *entity.Vale) error,
	getByID func(id int64) (*entity.Vale, error),
	getByFuncID func(funcionarioID int64) ([]entity.Vale, error),
	update func(v *entity.Vale) error,
	softDelete func(id int64) error,
	delete func(id int64) error,
	listPendentes func() ([]entity.Vale, error),
	listAprovadosNaoPagos func() ([]entity.Vale, error),
) *ValeRepositoryAdapter {
	return &ValeRepositoryAdapter{
		create:                create,
		getByID:               getByID,
		getByFuncID:           getByFuncID,
		update:                update,
		softDelete:            softDelete,
		delete:                delete,
		listPendentes:         listPendentes,
		listAprovadosNaoPagos: listAprovadosNaoPagos,
	}
}

// Implementações para o ValeService
func (a *ValeRepositoryAdapter) Create(v *entity.Vale) error {
	return a.create(v)
}

func (a *ValeRepositoryAdapter) GetByID(id int64) (*entity.Vale, error) {
	return a.getByID(id)
}

func (a *ValeRepositoryAdapter) GetValesByFuncionarioID(funcionarioID int64) ([]entity.Vale, error) {
	return a.getByFuncID(funcionarioID)
}

func (a *ValeRepositoryAdapter) Update(v *entity.Vale) error {
	return a.update(v)
}

func (a *ValeRepositoryAdapter) SoftDelete(id int64) error {
	return a.softDelete(id)
}

func (a *ValeRepositoryAdapter) Delete(id int64) error {
	return a.delete(id)
}

func (a *ValeRepositoryAdapter) ListPendentes() ([]entity.Vale, error) {
	return a.listPendentes()
}

func (a *ValeRepositoryAdapter) ListAprovadosNaoPagos() ([]entity.Vale, error) {
	return a.listAprovadosNaoPagos()
}
