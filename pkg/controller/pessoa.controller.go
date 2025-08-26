package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	mw "AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PessoaController struct {
	pessoaService *service.PessoaService
}

func NewPessoaController(p *service.PessoaService) *PessoaController {
	return &PessoaController{pessoaService: p}
}

// struct auxiliar para request em camelCase
type pessoaRequest struct {
	Nome              string `json:"nome"`
	CPF               string `json:"cpf"`
	RG                string `json:"rg"`
	Endereco          string `json:"endereco"`
	Contato           string `json:"contato"`
	ContatoEmergencia string `json:"contatoEmergencia"`
}

func (r *pessoaRequest) ToEntity() *entity.Pessoa {
	return &entity.Pessoa{
		Nome:              r.Nome,
		CPF:               r.CPF,
		RG:                r.RG,
		Endereco:          r.Endereco,
		Contato:           r.Contato,
		ContatoEmergencia: r.ContatoEmergencia,
	}
}

// CreatePessoa cria uma nova pessoa
func (c *PessoaController) CreatePessoa(w http.ResponseWriter, r *http.Request) {
	var input pessoaRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	p := input.ToEntity()
	if err := c.pessoaService.CreatePessoa(r.Context(), claims, p); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, p)
}

// GetPessoaByID retorna uma pessoa pelo ID
func (c *PessoaController) GetPessoaByID(w http.ResponseWriter, r *http.Request) {
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

	p, err := c.pessoaService.GetPessoaByID(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	if p == nil {
		httpjson.WriteJSON(w, http.StatusNotFound, httpjson.ErrorResponse{Error: "Pessoa não encontrada", Code: "NOT_FOUND"})
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, p)
}

// UpdatePessoa atualiza os dados de uma pessoa
func (c *PessoaController) UpdatePessoa(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	var input pessoaRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpjson.BadRequest(w, "JSON inválido")
		return
	}

	p := input.ToEntity()
	p.ID = id

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.pessoaService.UpdatePessoa(r.Context(), claims, p); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, p)
}

// DeletePessoa remove uma pessoa
func (c *PessoaController) DeletePessoa(w http.ResponseWriter, r *http.Request) {
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

	if err := c.pessoaService.DeletePessoa(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPessoas retorna todas as pessoas
func (c *PessoaController) ListPessoas(w http.ResponseWriter, r *http.Request) {
	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	list, err := c.pessoaService.ListPessoas(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, list)
}
