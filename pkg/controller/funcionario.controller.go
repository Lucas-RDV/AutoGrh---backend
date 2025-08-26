package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	mw "AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/service"
	"AutoGRH/pkg/utils/dateStringToTime"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FuncionarioController struct {
	funcionarioService *service.FuncionarioService
}

func NewFuncionarioController(f *service.FuncionarioService) *FuncionarioController {
	return &FuncionarioController{funcionarioService: f}
}

// struct auxiliar para request
// permite usar camelCase e datas no formato YYYY-MM-DD
type funcionarioRequest struct {
	PessoaID          int64   `json:"pessoaID"`
	PIS               string  `json:"pis"`
	CTPF              string  `json:"ctpf"`
	Nascimento        string  `json:"nascimento"`
	Admissao          string  `json:"admissao"`
	Cargo             string  `json:"cargo"`
	SalarioInicial    float64 `json:"salarioInicial"`
	FeriasDisponiveis int     `json:"feriasDisponiveis"`
}

func (r *funcionarioRequest) ToEntity() (*entity.Funcionario, error) {
	var f entity.Funcionario
	f.PessoaID = r.PessoaID
	f.PIS = r.PIS
	f.CTPF = r.CTPF
	f.Cargo = r.Cargo
	f.SalarioInicial = r.SalarioInicial
	f.FeriasDisponiveis = r.FeriasDisponiveis

	if r.Nascimento != "" {
		d, err := dateStringToTime.DateStringToTime(r.Nascimento)
		if err != nil {
			return nil, fmt.Errorf("nascimento inválido: %w", err)
		}
		f.Nascimento = d
	}
	if r.Admissao != "" {
		d, err := dateStringToTime.DateStringToTime(r.Admissao)
		if err != nil {
			return nil, fmt.Errorf("admissão inválida: %w", err)
		}
		f.Admissao = d
	}

	return &f, nil
}

// CreateFuncionario cria um novo funcionário
func (c *FuncionarioController) CreateFuncionario(w http.ResponseWriter, r *http.Request) {
	var input funcionarioRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	f, err := input.ToEntity()
	if err != nil {
		httpjson.BadRequest(w, err.Error())
		return
	}

	if err := c.funcionarioService.CreateFuncionario(r.Context(), claims, f); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, f)
}

// GetFuncionarioByID retorna um funcionário pelo ID
func (c *FuncionarioController) GetFuncionarioByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	f, err := c.funcionarioService.GetFuncionarioByID(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if f == nil {
		httpjson.WriteJSON(w, http.StatusNotFound, httpjson.ErrorResponse{Error: "Funcionário não encontrado", Code: "NOT_FOUND"})
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, f)
}

// UpdateFuncionario atualiza um funcionário existente
func (c *FuncionarioController) UpdateFuncionario(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	var input funcionarioRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	f, err := input.ToEntity()
	if err != nil {
		httpjson.BadRequest(w, err.Error())
		return
	}
	f.ID = id

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.funcionarioService.UpdateFuncionario(r.Context(), claims, f); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, f)
}

// DeleteFuncionario remove um funcionário
func (c *FuncionarioController) DeleteFuncionario(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.funcionarioService.DeleteFuncionario(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListFuncionariosAtivos retorna todos os funcionários ativos
func (c *FuncionarioController) ListFuncionariosAtivos(w http.ResponseWriter, r *http.Request) {
	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	list, err := c.funcionarioService.ListFuncionariosAtivos(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, list)
}

// ListFuncionariosInativos retorna todos os funcionários inativos
func (c *FuncionarioController) ListFuncionariosInativos(w http.ResponseWriter, r *http.Request) {
	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	list, err := c.funcionarioService.ListFuncionariosInativos(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, list)
}

// ListTodosFuncionarios retorna todos os funcionários, ativos e inativos
func (c *FuncionarioController) ListTodosFuncionarios(w http.ResponseWriter, r *http.Request) {
	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	list, err := c.funcionarioService.ListTodosFuncionarios(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, list)
}
