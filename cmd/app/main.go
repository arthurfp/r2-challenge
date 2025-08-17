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
		fx.Invoke(runHTTPServer),
	)
	app.Run()
}

func runHTTPServer(
	lc fx.Lifecycle,
	envs envs.Envs,
	tracer observability.Tracer,
	_ *db.Database,
) error {
	e := httpx.NewServer(tracer)

	v1 := e.Group("/v1")
	_ = v1

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

