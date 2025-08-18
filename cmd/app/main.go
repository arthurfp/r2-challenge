package main

// @title           R2 Challenge API
// @version         1.0
// @description     REST API for e-commerce (products, users, orders)
// @BasePath        /
// @schemes         http https

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/fx"

	"r2-challenge/cmd/envs"
	"r2-challenge/pkg/httpx"
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

	"github.com/labstack/echo/v4"
)

// @title R2 Ecoomerce API

// @BasePath /v1
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
			userqry.NewGetByIDService,
			userqry.NewGetByEmailService,
			userqry.NewListUsersService,
			userhttp.NewRegisterHandler,
			userhttp.NewLoginHandler,
			userhttp.NewGetUserHandler,
			userhttp.NewListUsersHandler,
			userhttp.NewUpdateProfileHandler,

			orderdb.NewDBRepository,
			payment.NewNoopProcessor,
			notification.NewNoopSender,
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

	// Swagger UI serving (Option A): serve generated YAML and a minimal UI page
	e.File("/swagger.yaml", "/app/swagger.yaml")
	e.GET("/swagger", func(c echo.Context) error {
		html := `<!doctype html><html><head><meta charset="utf-8"/><title>R2 Challenge API</title>
	<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css"></head>
	<body><div id="header" style="padding:8px;border-bottom:1px solid #eee;">
	  <input id="token" placeholder="Bearer <token>" style="width:60%;padding:6px;" />
	  <button onclick="setToken()" style="padding:6px 12px;">Set Token</button>
	</div>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
	<script>
	  let authToken = localStorage.getItem('authToken') || '';
	  function setToken(){
	    let t = document.getElementById('token').value.trim();
	    if(t && !/^Bearer\s+/i.test(t)) { t = 'Bearer ' + t; }
	    authToken = t;
	    localStorage.setItem('authToken', authToken);
	  }
	  // Prefill input from storage
	  window.addEventListener('DOMContentLoaded', function(){
	    if(authToken){ document.getElementById('token').value = authToken; }
	  });
	  window.ui = SwaggerUIBundle({
	    url:'/swagger.yaml', dom_id:'#swagger-ui',
	    requestInterceptor: (req) => { if(authToken){ req.headers['Authorization'] = authToken; } return req; },
	    responseInterceptor: (res) => {
	      try {
	        if(res && res.url && res.status === 200 && /\/v1\/auth\/login$/.test(res.url)){
	          const data = typeof res.data === 'string' ? JSON.parse(res.data) : res.data;
	          if(data && data.access_token){
	            authToken = 'Bearer ' + data.access_token;
	            localStorage.setItem('authToken', authToken);
	            const el = document.getElementById('token'); if(el){ el.value = authToken; }
	          }
	        }
	      } catch(e){}
	      return res;
	    }
	  });
	</script>
	</body></html>`
		return c.HTML(http.StatusOK, html)
	})

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
		// Always allow preflight and Swagger/metrics endpoints
		if method == "OPTIONS" || path == "/swagger" || path == "/swagger.yaml" || path == envs.MetricsPath {
			return true
		}
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

	v1.GET("/users/:id", auth.RequireSelfOrRoles("id", "admin")(getUser.Handle))
	v1.GET("/users", auth.RequireRoles("admin")(listUsers.Handle))
	v1.PUT("/users/me", updateProfile.Handle)

	// Orders
	v1.POST("/orders", place.Handle)
	v1.GET("/orders/:id", auth.RequireRoles("admin")(getOrder.Handle))
	v1.GET("/users/:id/orders", listOrders.Handle)
	v1.PUT("/orders/:id/status", auth.RequireRoles("admin")(updateOrderStatus.Handle))

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

