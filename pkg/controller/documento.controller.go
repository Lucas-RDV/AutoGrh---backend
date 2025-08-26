package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	mw "AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/entity"
	"AutoGRH/pkg/service"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DocumentoController struct {
	documentoService *service.DocumentoService
}

func NewDocumentoController(s *service.DocumentoService) *DocumentoController {
	return &DocumentoController{documentoService: s}
}

// documentoRequest para receber uploads em JSON base64 ou binário simples
type documentoRequest struct {
	Doc string `json:"doc"` // conteúdo base64 ou string
}

// CreateDocumento cria um novo documento vinculado a um funcionário
func (c *DocumentoController) CreateDocumento(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	funcionarioID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || funcionarioID <= 0 {
		httpjson.BadRequest(w, "funcionarioID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	// aceita tanto JSON quanto upload direto no corpo (binário)
	var d entity.Documento
	d.FuncionarioID = funcionarioID

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		var req documentoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpjson.BadRequest(w, "JSON inválido")
			return
		}
		d.Doc = []byte(req.Doc)
	} else {
		data, err := io.ReadAll(r.Body)
		if err != nil || len(data) == 0 {
			httpjson.BadRequest(w, "conteúdo do documento inválido")
			return
		}
		d.Doc = data
	}

	if err := c.documentoService.CreateDocumento(r.Context(), claims, &d); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, d)
}

// GetDocumentosByFuncionarioID retorna documentos de um funcionário
func (c *DocumentoController) GetDocumentosByFuncionarioID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	funcionarioID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || funcionarioID <= 0 {
		httpjson.BadRequest(w, "funcionarioID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	docs, err := c.documentoService.GetDocumentosByFuncionarioID(r.Context(), claims, funcionarioID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, docs)
}

// ListDocumentos retorna todos os documentos
func (c *DocumentoController) ListDocumentos(w http.ResponseWriter, r *http.Request) {
	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	docs, err := c.documentoService.ListDocumentos(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, docs)
}

// DeleteDocumento remove um documento (somente admin)
func (c *DocumentoController) DeleteDocumento(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	docID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || docID <= 0 {
		httpjson.BadRequest(w, "ID inválido")
		return
	}

	claims, ok := mw.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "não autenticado")
		return
	}

	if err := c.documentoService.DeleteDocumento(r.Context(), claims, docID); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
