package main

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"

	"r2-challenge/cmd/envs"
	"r2-challenge/pkg/httpx"
	"github.com/labstack/echo/v4"
	"r2-challenge/pkg/metrics"
	"r2-challenge/pkg/db"
	"r2-challenge/pkg/logger"
	"r2-challenge/pkg/observability"
	"r2-challenge/pkg/validator"

	"r2-challenge/pkg/auth"

	producthttp "r2-challenge/internal/product/adapters/http"
	productdb "r2-challenge/internal/product/adapters/db"
	productcmd "r2-challenge/internal/product/services/command"
	productqry "r2-challenge/internal/product/services/query"

	userhttp "r2-challenge/internal/user/adapters/http"
	userdb "r2-challenge/internal/user/adapters/db"
	usercmd "r2-challenge/internal/user/services/command"
	userqry "r2-challenge/internal/user/services/query"

	orderhttp "r2-challenge/internal/order/adapters/http"
	orderdb "r2-challenge/internal/order/adapters/db"
	ordercmd "r2-challenge/internal/order/services/command"
	orderqry "r2-challenge/internal/order/services/query"
	payment "r2-challenge/internal/order/adapters/payment"
	notification "r2-challenge/internal/order/adapters/notification"
	pmtdb "r2-challenge/internal/payment/adapters/db"
	pmtcmd "r2-challenge/internal/payment/services/command"
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

			userdb.NewDBRepository,
			usercmd.NewRegisterService,
			usercmd.NewUpdateProfileService,
			userqry.NewService,
			userhttp.NewRegisterHandler,
			userhttp.NewLoginHandler,
			userhttp.NewGetUserHandler,
			userhttp.NewListUsersHandler,
			userhttp.NewUpdateProfileHandler,

			orderdb.NewDBRepository,
			payment.NewMockProcessor,
			notification.NewMockSender,
			pmtdb.NewDBRepository,
			pmtcmd.NewService,
			ordercmd.NewPlaceOrderService,
			ordercmd.NewUpdateStatusService,
			orderqry.NewService,
			orderhttp.NewPlaceOrderHandler,
			orderhttp.NewGetOrderHandler,
			orderhttp.NewListUserOrdersHandler,
			orderhttp.NewUpdateStatusHandler,
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
	register userhttp.RegisterHandler,
	login userhttp.LoginHandler,
	getUser userhttp.GetUserHandler,
	listUsers userhttp.ListUsersHandler,
	updateProfile userhttp.UpdateProfileHandler,
	place orderhttp.PlaceOrderHandler,
	getOrder orderhttp.GetOrderHandler,
	listOrders orderhttp.ListUserOrdersHandler,
	updateOrderStatus orderhttp.UpdateStatusHandler,
) error {
	e := httpx.NewServer(tracer)

	// Metrics
	if envs.MetricsEnabled {
		provider, handler, err := metrics.Setup()
		if err == nil {
			_ = provider // keep provider alive
			e.Use(httpx.MetricsMiddleware(provider.Meter("r2-challenge")))
			if envs.MetricsPort != "" {
				go func() { _ = http.ListenAndServe(fmt.Sprintf(":%s", envs.MetricsPort), handler) }()
			} else {
				e.GET(envs.MetricsPath, echo.WrapHandler(handler))
			}
		}
	}

	// Rate limit per IP (RPM configurable)
	e.Use(httpx.RateLimitMiddleware(envs.RateLimitRPM))

	v1 := e.Group("/v1")

	ttl, _ := time.ParseDuration(envs.JWTExpire)
	tm := auth.NewTokenManager(envs.JWTSecret, envs.JWTIssuer, ttl)
	e.Use(auth.JWTMiddleware(tm, func(method, path string) bool {
		_, ok := publicRoutes[routeKey{Method: method, Path: path}]
		return ok
	}))

	// Product routes (admin-only for mutations)
	v1.POST("/products", auth.RequireRoles("admin")(create.Handle))
	v1.PUT("/products/:id", auth.RequireRoles("admin")(update.Handle))
	v1.DELETE("/products/:id", auth.RequireRoles("admin")(deleteH.Handle))
	v1.GET("/products/:id", get.Handle)
	v1.GET("/products", list.Handle)

	// Auth / Users
	authg := v1.Group("/auth")
	authg.POST("/register", register.Handle)
	authg.POST("/login", login.Handle)

	v1.GET("/users/:id", getUser.Handle)
	v1.GET("/users", listUsers.Handle)
	v1.PUT("/users/me", updateProfile.Handle)

	// Orders
	v1.POST("/orders", place.Handle)
	v1.GET("/orders/:id", getOrder.Handle)
	v1.GET("/orders", listOrders.Handle)
	v1.PUT("/orders/:id/status", updateOrderStatus.Handle)

	readHeaderTimeout, _ := time.ParseDuration(envs.ReadHeaderTimeout)
	httpTimeout, _ := time.ParseDuration(envs.HTTPTimeout)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", envs.HTTPPort),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           e,
	}

	// TLS if cert/key provided via envs
	lc.Append(httpx.ServerLifecycleTLS(e, server, httpTimeout, envs.TLSCertFile, envs.TLSKeyFile))
	return nil
}

