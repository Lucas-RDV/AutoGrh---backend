package Adapter

import (
	"AutoGRH/pkg/entity"
)

// Interface que o service consome
type PagamentoRepository interface {
	Create(p *entity.Pagamento) error
	Update(p *entity.Pagamento) error
	GetPagamentosByFolhaID(folhaID int64) ([]entity.Pagamento, error)
	DeletePagamentosByFolhaID(folhaID int64) error
	GetPagamentoByID(id int64) (*entity.Pagamento, error)
	ListPagamentosByFuncionarioID(funcionarioID int64) ([]entity.Pagamento, error)
}

// Adapter para conectar o repository ao service
type PagamentoRepositoryAdapter struct {
	create                        func(p *entity.Pagamento) error
	update                        func(p *entity.Pagamento) error
	getPagamentosByFolhaID        func(folhaID int64) ([]entity.Pagamento, error)
	deletePagamentosByFolhaID     func(folhaID int64) error
	getPagamentoByID              func(id int64) (*entity.Pagamento, error)
	listPagamentosByFuncionarioID func(funcionarioID int64) ([]entity.Pagamento, error)
}

func NewPagamentoRepositoryAdapter(
	create func(p *entity.Pagamento) error,
	update func(p *entity.Pagamento) error,
	getPagamentosByFolhaID func(folhaID int64) ([]entity.Pagamento, error),
	deletePagamentosByFolhaID func(folhaID int64) error,
	getPagamentoByID func(id int64) (*entity.Pagamento, error),
	listPagamentosByFuncionarioID func(funcionarioID int64) ([]entity.Pagamento, error),
) PagamentoRepository {
	return &PagamentoRepositoryAdapter{
		create:                        create,
		update:                        update,
		getPagamentosByFolhaID:        getPagamentosByFolhaID,
		deletePagamentosByFolhaID:     deletePagamentosByFolhaID,
		getPagamentoByID:              getPagamentoByID,
		listPagamentosByFuncionarioID: listPagamentosByFuncionarioID,
	}
}

// Métodos implementando a interface
func (a *PagamentoRepositoryAdapter) Create(p *entity.Pagamento) error {
	return a.create(p)
}

func (a *PagamentoRepositoryAdapter) Update(p *entity.Pagamento) error {
	return a.update(p)
}

func (a *PagamentoRepositoryAdapter) GetPagamentosByFolhaID(folhaID int64) ([]entity.Pagamento, error) {
	return a.getPagamentosByFolhaID(folhaID)
}

func (a *PagamentoRepositoryAdapter) DeletePagamentosByFolhaID(folhaID int64) error {
	return a.deletePagamentosByFolhaID(folhaID)
}

func (a *PagamentoRepositoryAdapter) GetPagamentoByID(id int64) (*entity.Pagamento, error) {
	return a.getPagamentoByID(id)
}

func (a *PagamentoRepositoryAdapter) ListPagamentosByFuncionarioID(funcionarioID int64) ([]entity.Pagamento, error) {
	return a.listPagamentosByFuncionarioID(funcionarioID)
}
