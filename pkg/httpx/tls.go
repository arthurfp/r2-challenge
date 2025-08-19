package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// ServerLifecycleTLS starts HTTP server using TLS when cert/key are provided.
func ServerLifecycleTLS(e *echo.Echo, s *http.Server, timeout time.Duration, certFile, keyFile string) fx.Hook {
	useTLS := certFile != "" && keyFile != ""

	return fx.Hook{
		OnStart: func(ctx context.Context) error {
			if useTLS {
				go func() { _ = s.ListenAndServeTLS(certFile, keyFile) }()
				return nil
			}

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
