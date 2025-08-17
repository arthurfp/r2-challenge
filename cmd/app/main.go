package main

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"

	"r2-challenge/cmd/envs"
	"r2-challenge/pkg/httpx"
	"r2-challenge/pkg/db"
	"r2-challenge/pkg/logger"
	"r2-challenge/pkg/observability"
	"r2-challenge/pkg/validator"

	producthttp "r2-challenge/internal/product/adapters/http"
	productdb "r2-challenge/internal/product/adapters/db"
	productcmd "r2-challenge/internal/product/services/command"
	productqry "r2-challenge/internal/product/services/query"
)

func main() {
	app := fx.New(
		fx.Provide(
			envs.NewEnvs,
			logger.Setup,
			observability.SetupTracer,
			validator.Setup,
			db.Setup,
		),

		fx.Provide(
			productdb.NewDBRepository,
			productcmd.NewCreateService,
			productcmd.NewUpdateService,
			productcmd.NewDeleteService,
			productqry.NewService,
			producthttp.NewCreateHandler,
			producthttp.NewUpdateHandler,
			producthttp.NewDeleteHandler,
			producthttp.NewGetHandler,
			producthttp.NewListHandler,
		),

		fx.Invoke(runHTTPServer),
	)

	app.Run()
}

func runHTTPServer(
	lc fx.Lifecycle,
	envs envs.Envs,
	tracer observability.Tracer,
	_ *db.Database,
	create producthttp.CreateHandler,
	update producthttp.UpdateHandler,
	deleteH producthttp.DeleteHandler,
	get producthttp.GetHandler,
	list producthttp.ListHandler,
) error {
	e := httpx.NewServer(tracer)

	v1 := e.Group("/v1")
	v1.POST("/products", create.Handle)
	v1.GET("/products/:id", get.Handle)
	v1.GET("/products", list.Handle)
	v1.PUT("/products/:id", update.Handle)
	v1.DELETE("/products/:id", deleteH.Handle)

	readHeaderTimeout, _ := time.ParseDuration(envs.ReadHeaderTimeout)
	httpTimeout, _ := time.ParseDuration(envs.HTTPTimeout)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", envs.HTTPPort),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           e,
	}

	lc.Append(httpx.ServerLifecycle(e, server, httpTimeout))
	return nil
}

