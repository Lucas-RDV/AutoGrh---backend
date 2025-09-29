package Adapter

import (
	"AutoGRH/pkg/entity"
)

// Interface usada pelo service
type FolhaPagamentoRepository interface {
	Create(f *entity.FolhaPagamentos) error
	GetByID(id int64) (*entity.FolhaPagamentos, error)
	GetByMesAnoTipo(mes, ano int, tipo string) (*entity.FolhaPagamentos, error)
	Update(f *entity.FolhaPagamentos) error
	Delete(id int64) error
	List() ([]entity.FolhaPagamentos, error)
	MarcarComoPaga(id int64) error
}

// Adapter que implementa a interface acima
type FolhaPagamentoRepositoryAdapter struct {
	create          func(*entity.FolhaPagamentos) error
	getByID         func(int64) (*entity.FolhaPagamentos, error)
	getByMesAnoTipo func(int, int, string) (*entity.FolhaPagamentos, error)
	update          func(*entity.FolhaPagamentos) error
	delete          func(int64) error
	list            func() ([]entity.FolhaPagamentos, error)
	marcarComoPaga  func(int64) error
}

// Construtor do adapter
func NewFolhaPagamentoRepositoryAdapter(
	create func(*entity.FolhaPagamentos) error,
	getByID func(int64) (*entity.FolhaPagamentos, error),
	getByMesAnoTipo func(int, int, string) (*entity.FolhaPagamentos, error),
	update func(*entity.FolhaPagamentos) error,
	delete func(int64) error,
	list func() ([]entity.FolhaPagamentos, error),
	marcarComoPaga func(int64) error,
) *FolhaPagamentoRepositoryAdapter {
	return &FolhaPagamentoRepositoryAdapter{
		create:          create,
		getByID:         getByID,
		getByMesAnoTipo: getByMesAnoTipo,
		update:          update,
		delete:          delete,
		list:            list,
		marcarComoPaga:  marcarComoPaga,
	}
}

// Implementações da interface
func (a *FolhaPagamentoRepositoryAdapter) Create(f *entity.FolhaPagamentos) error {
	return a.create(f)
}

func (a *FolhaPagamentoRepositoryAdapter) GetByID(id int64) (*entity.FolhaPagamentos, error) {
	return a.getByID(id)
}

func (a *FolhaPagamentoRepositoryAdapter) GetByMesAnoTipo(mes, ano int, tipo string) (*entity.FolhaPagamentos, error) {
	return a.getByMesAnoTipo(mes, ano, tipo)
}

func (a *FolhaPagamentoRepositoryAdapter) Update(f *entity.FolhaPagamentos) error {
	return a.update(f)
}

func (a *FolhaPagamentoRepositoryAdapter) Delete(id int64) error {
	return a.delete(id)
}

func (a *FolhaPagamentoRepositoryAdapter) List() ([]entity.FolhaPagamentos, error) {
	return a.list()
}

func (a *FolhaPagamentoRepositoryAdapter) MarcarComoPaga(id int64) error {
	return a.marcarComoPaga(id)
}
