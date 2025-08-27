package controller

import (
	"AutoGRH/pkg/controller/httpjson"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DocumentoController struct {
	documentoService *service.DocumentoService
}

func NewDocumentoController(documentoService *service.DocumentoService) *DocumentoController {
	return &DocumentoController{documentoService: documentoService}
}

// CreateDocumento - upload multipart
func (c *DocumentoController) CreateDocumento(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	funcionarioIDStr := chi.URLParam(r, "id")
	funcionarioID, err := strconv.ParseInt(funcionarioIDStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "funcionarioID inválido")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		httpjson.BadRequest(w, "arquivo não enviado (esperado campo 'file')")
		return
	}
	defer file.Close()

	doc, err := c.documentoService.SalvarDocumento(r.Context(), claims, funcionarioID, file, header.Filename)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, doc)
}

// GetDocumentosByFuncionarioID
func (c *DocumentoController) GetDocumentosByFuncionarioID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	funcionarioIDStr := chi.URLParam(r, "id")
	funcionarioID, err := strconv.ParseInt(funcionarioIDStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "funcionarioID inválido")
		return
	}

	docs, err := c.documentoService.GetDocumentosByFuncionarioID(r.Context(), claims, funcionarioID)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, docs)
}

// ListDocumentos - lista todos
func (c *DocumentoController) ListDocumentos(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	docs, err := c.documentoService.ListDocumentos(r.Context(), claims)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, docs)
}

// DownloadDocumento - baixa arquivo físico
func (c *DocumentoController) DownloadDocumento(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	fullPath, err := c.documentoService.GetDocumentoPath(r.Context(), claims, id)
	if err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	http.ServeFile(w, r, fullPath)
}

// DeleteDocumento - remove do banco e do disco
func (c *DocumentoController) DeleteDocumento(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		httpjson.Unauthorized(w, "UNAUTHORIZED", "usuário não autenticado")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpjson.BadRequest(w, "id inválido")
		return
	}

	if err := c.documentoService.DeleteDocumento(r.Context(), claims, id); err != nil {
		httpjson.Internal(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
