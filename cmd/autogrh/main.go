package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"AutoGRH/pkg/Bootstrap"
	"AutoGRH/pkg/HTTP/router"
	middleware "AutoGRH/pkg/controller/middleware"
)

func main() {
	_ = godotenv.Load()

	app := Bootstrap.Load()
	if err := Bootstrap.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	auth := Bootstrap.BuildAuth(app)
	pessoaSvc := Bootstrap.BuildPessoaService(auth)
	funcSvc := Bootstrap.BuildFuncionarioService(auth)
	documentoSvc := Bootstrap.BuildDocumentoService(auth)
	faltaSvc := Bootstrap.BuildFaltaService(auth)
	feriasSvc := Bootstrap.BuildFeriasService(auth)
	descansoSvc := Bootstrap.BuildDescansoService(auth)
	salarioSvc := Bootstrap.BuildSalarioService(auth)
	salarioRealSvc := Bootstrap.BuildSalarioRealService(auth)
	valeCtl := Bootstrap.BuildValeService(auth)
	folhaCtl := Bootstrap.BuildFolhaPagamentoService(auth)
	pagamentoCtl := Bootstrap.BuildPagamentoService(auth)
	avisoSvc := Bootstrap.BuildAvisoService(auth)

	// Inicializar workers
	Bootstrap.InitWorkers(feriasSvc, descansoSvc, salarioRealSvc, funcSvc, faltaSvc, folhaCtl, avisoSvc)

	routes := router.New(auth, pessoaSvc, funcSvc, documentoSvc, faltaSvc, feriasSvc, descansoSvc, salarioSvc, salarioRealSvc, valeCtl, folhaCtl, pagamentoCtl, avisoSvc)

	cors := middleware.NewCORS(middleware.CORSConfig{

		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   nil,
		AllowedHeaders:   []string{"Content-Type", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	http.ListenAndServe(":8080", cors(routes))
}
