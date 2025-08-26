package router

import (
	"net/http"

	"AutoGRH/pkg/controller"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"

	"github.com/go-chi/chi/v5"
)

func New(auth *service.AuthService, pessoaSvc *service.PessoaService, funcionarioSvc *service.FuncionarioService, documentoSvc *service.DocumentoService) http.Handler {
	r := chi.NewRouter()

	authCtl := controller.NewAuthController(auth)
	users := service.NewUsuarioService()
	adminCtl := controller.NewAdminController(users)
	pessoaCtl := controller.NewPessoaController(pessoaSvc)
	funcionarioCtl := controller.NewFuncionarioController(funcionarioSvc)
	documentoCtl := controller.NewDocumentoController(documentoSvc)

	// Rota pública
	r.Post("/auth/login", authCtl.Login)

	// Rota autenticada básica
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(auth))
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			if claims, ok := middleware.GetClaims(r.Context()); ok {
				w.Write([]byte("ola " + claims.Nome))
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
		})
	})

	// Rotas com permissão para gerenciar usuários
	r.Route("/admin/usuarios", func(r chi.Router) {
		r.Use(middleware.RequirePerm(auth, "usuario:list"))
		r.Get("/", adminCtl.UsuariosList)

		r.With(middleware.RequirePerm(auth, "usuario:create")).
			Post("/", adminCtl.CreateUsuario)

		r.With(middleware.RequirePerm(auth, "usuario:update")).
			Put("/{id}", adminCtl.UpdateUsuario)

		r.With(middleware.RequirePerm(auth, "usuario:delete")).
			Delete("/{id}", adminCtl.DeleteUsuario)
	})

	// Rotas de Pessoas
	r.Route("/pessoas", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "pessoa:list")).Get("/", pessoaCtl.ListPessoas)
		r.With(middleware.RequirePerm(auth, "pessoa:create")).Post("/", pessoaCtl.CreatePessoa)
		r.With(middleware.RequirePerm(auth, "pessoa:update")).Put("/{id}", pessoaCtl.UpdatePessoa)
		r.With(middleware.RequirePerm(auth, "pessoa:delete")).Delete("/{id}", pessoaCtl.DeletePessoa)
		r.With(middleware.RequirePerm(auth, "pessoa:read")).Get("/{id}", pessoaCtl.GetPessoaByID)
	})

	// Rotas de Funcionários
	r.Route("/funcionarios", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "funcionario:list")).Get("/", funcionarioCtl.ListTodosFuncionarios)
		r.With(middleware.RequirePerm(auth, "funcionario:list")).Get("/ativos", funcionarioCtl.ListFuncionariosAtivos)
		r.With(middleware.RequirePerm(auth, "funcionario:list")).Get("/inativos", funcionarioCtl.ListFuncionariosInativos)

		r.With(middleware.RequirePerm(auth, "funcionario:create")).Post("/", funcionarioCtl.CreateFuncionario)
		r.With(middleware.RequirePerm(auth, "funcionario:update")).Put("/{id}", funcionarioCtl.UpdateFuncionario)
		r.With(middleware.RequirePerm(auth, "funcionario:delete")).Delete("/{id}", funcionarioCtl.DeleteFuncionario)
		r.With(middleware.RequirePerm(auth, "funcionario:read")).Get("/{id}", funcionarioCtl.GetFuncionarioByID)

		// Documentos dentro de funcionário
		r.With(middleware.RequirePerm(auth, "documento:create")).Post("/{id}/documentos", documentoCtl.CreateDocumento)
		r.With(middleware.RequirePerm(auth, "documento:list")).Get("/{id}/documentos", documentoCtl.GetDocumentosByFuncionarioID)
	})

	// Rotas diretas de Documentos
	r.Route("/documentos", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "documento:list")).Get("/", documentoCtl.ListDocumentos)
		r.With(middleware.RequirePerm(auth, "documento:delete")).Delete("/{id}", documentoCtl.DeleteDocumento)
	})

	return r
}
