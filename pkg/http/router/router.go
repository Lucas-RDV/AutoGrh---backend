package router

import (
	"net/http"

	"AutoGRH/pkg/controller"
	"AutoGRH/pkg/controller/middleware"
	"AutoGRH/pkg/service"

	"github.com/go-chi/chi/v5"
)

func New(
	auth *service.AuthService,
	pessoaSvc *service.PessoaService,
	funcionarioSvc *service.FuncionarioService,
	documentoSvc *service.DocumentoService,
	faltaSvc *service.FaltaService,
	feriasSvc *service.FeriasService,
) http.Handler {
	r := chi.NewRouter()

	authCtl := controller.NewAuthController(auth)
	users := service.NewUsuarioService()
	adminCtl := controller.NewAdminController(users)
	pessoaCtl := controller.NewPessoaController(pessoaSvc)
	funcionarioCtl := controller.NewFuncionarioController(funcionarioSvc)
	documentoCtl := controller.NewDocumentoController(documentoSvc)
	faltaCtl := controller.NewFaltaController(faltaSvc)
	feriasCtl := controller.NewFeriasController(feriasSvc)

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

		// Faltas dentro de funcionário
		r.With(middleware.RequirePerm(auth, "falta:list")).Get("/{id}/faltas", faltaCtl.GetFaltasByFuncionarioID)
		r.With(middleware.RequirePerm(auth, "falta:create")).Post("/{id}/faltas", faltaCtl.CreateFalta)

		// Ferias dentro de funcionário
		r.With(middleware.RequirePerm(auth, "ferias:list")).Get("/{id}/ferias", feriasCtl.GetFeriasByFuncionarioID)
	})

	// Rotas diretas de Documentos
	r.Route("/documentos", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "documento:list")).Get("/", documentoCtl.ListDocumentos)
		r.With(middleware.RequirePerm(auth, "documento:list")).Get("/{id}/download", documentoCtl.DownloadDocumento)
		r.With(middleware.RequirePerm(auth, "documento:delete")).Delete("/{id}", documentoCtl.DeleteDocumento)
	})

	// Rotas diretas de Faltas
	r.Route("/faltas", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "falta:list")).Get("/", faltaCtl.ListAllFaltas)
		r.With(middleware.RequirePerm(auth, "falta:read")).Get("/{id}", faltaCtl.GetFaltaByID)
		r.With(middleware.RequirePerm(auth, "falta:update")).Put("/{id}", faltaCtl.UpdateFalta)
		r.With(middleware.RequirePerm(auth, "falta:delete")).Delete("/{id}", faltaCtl.DeleteFalta)
	})

	//rotas diretas de Ferias
	r.Route("/ferias", func(r chi.Router) {
		r.With(middleware.RequirePerm(auth, "ferias:list")).Get("/", feriasCtl.ListFerias)
		r.With(middleware.RequirePerm(auth, "ferias:read")).Get("/{id}", feriasCtl.GetFeriasByID)
		r.With(middleware.RequirePerm(auth, "ferias:update")).Put("/{id}/vencida", feriasCtl.MarcarComoVencidas)
		r.With(middleware.RequirePerm(auth, "ferias:update")).Put("/{id}/terco-pago", feriasCtl.MarcarTercoComoPago)
		r.With(middleware.RequirePerm(auth, "ferias:read")).Get("/{id}/saldo", feriasCtl.GetSaldoFerias)
	})
	return r
}
