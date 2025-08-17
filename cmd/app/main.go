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
			userqry.NewService,
			userhttp.NewRegisterHandler,
			userhttp.NewLoginHandler,
			userhttp.NewGetUserHandler,
			userhttp.NewListUsersHandler,

			orderdb.NewDBRepository,
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
	place orderhttp.PlaceOrderHandler,
	getOrder orderhttp.GetOrderHandler,
	listOrders orderhttp.ListUserOrdersHandler,
	updateOrderStatus orderhttp.UpdateStatusHandler,
) error {
	e := httpx.NewServer(tracer)

	v1 := e.Group("/v1")

	ttl, _ := time.ParseDuration(envs.JWTExpire)
	tm := auth.NewTokenManager(envs.JWTSecret, envs.JWTIssuer, ttl)
	e.Use(auth.JWTMiddleware(tm, func(method, path string) bool {
		_, ok := publicRoutes[routeKey{Method: method, Path: path}]
		return ok
	}))

	// Product routes
	v1.POST("/products", create.Handle)
	v1.GET("/products/:id", get.Handle)
	v1.GET("/products", list.Handle)
	v1.PUT("/products/:id", update.Handle)
	v1.DELETE("/products/:id", deleteH.Handle)

	// Auth / Users
	authg := v1.Group("/auth")
	authg.POST("/register", register.Handle)
	authg.POST("/login", login.Handle)

	v1.GET("/users/:id", getUser.Handle)
	v1.GET("/users", listUsers.Handle)

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

	lc.Append(httpx.ServerLifecycle(e, server, httpTimeout))
	return nil
}

