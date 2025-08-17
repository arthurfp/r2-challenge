package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"r2-challenge/pkg/observability"
)

type echoTracer struct{ tracer observability.Tracer }

func (et echoTracer) Middleware(skipPaths ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, end := et.tracer.StartSpan(c.Request().Context(), c.Request().Method+" "+c.Path())
			defer end()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func NewServer(tracer observability.Tracer) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(echoTracer{tracer: tracer}.Middleware())
	return e
}

func ServerLifecycle(e *echo.Echo, s *http.Server, timeout time.Duration) fx.Hook {
	return fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() { _ = s.ListenAndServe() }()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			_ = s.Shutdown(ctx)
			return nil
		},
	}
}
